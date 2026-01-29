package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const NoDLOpenRuleID = "no-dlopen"

// NoDLOpenRule checks if dlopen is disabled
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type NoDLOpenRule struct{}

func (r NoDLOpenRule) ID() string   { return NoDLOpenRuleID }
func (r NoDLOpenRule) Name() string { return "Disallow dlopen" }

func (r NoDLOpenRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,nodlopen"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,nodlopen"},
		},
	}
}

func (r NoDLOpenRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	if bin.File.Type != elf.ET_DYN {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Not a shared library",
		}
	}

	if HasDynFlag(bin.File, elf.DT_FLAGS_1, uint64(elf.DF_1_NOOPEN)) {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "dlopen disabled",
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "dlopen not disabled",
	}
}
