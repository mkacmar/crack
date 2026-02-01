package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const NoDumpRuleID = "no-dump"

// NoDumpRule checks if core dumps are disabled
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type NoDumpRule struct{}

func (r NoDumpRule) ID() string   { return NoDumpRuleID }
func (r NoDumpRule) Name() string { return "Core Dump Protection" }
func (r NoDumpRule) Description() string {
	return "Checks if the binary is excluded from core dumps. Disabling core dumps prevents sensitive data like cryptographic keys and passwords from being written to disk during crashes."
}

func (r NoDumpRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, Flag: "-Wl,-z,nodump"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, Flag: "-Wl,-z,nodump"},
		},
	}
}

func (r NoDumpRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	if bin.File.Type != elf.ET_EXEC && bin.File.Type != elf.ET_DYN {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	if HasDynFlag(bin.File, elf.DT_FLAGS_1, uint64(elf.DF_1_NODUMP)) {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "Core dumps disabled",
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "Core dumps not explicitly disabled",
	}
}
