package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// NoDumpRuleID is the rule ID for no dump.
const NoDumpRuleID = "no-dump"

// NoDumpRule checks if core dumps are disabled.
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type NoDumpRule struct{}

func (r NoDumpRule) ID() string   { return NoDumpRuleID }
func (r NoDumpRule) Name() string { return "Core Dump Protection" }
func (r NoDumpRule) Description() string {
	return "Checks if the binary has the DF_1_NODUMP flag set to prevent dldump(3) from copying the shared object. This restricts an attacker's ability to extract the mapped object from a running process."
}

func (r NoDumpRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, Flag: "-Wl,-z,nodump"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, Flag: "-Wl,-z,nodump"},
		},
	}
}

func (r NoDumpRule) Execute(bin *binary.ELFBinary) rule.Result {
	if bin.Type != elf.ET_EXEC && bin.Type != elf.ET_DYN {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	if bin.HasDynFlag(elf.DT_FLAGS_1, uint64(elf.DF_1_NODUMP)) {
		return rule.Result{
			Status:  rule.StatusPassed,
			Message: "Core dumps disabled",
		}
	}

	return rule.Result{
		Status:  rule.StatusFailed,
		Message: "Core dumps not explicitly disabled",
	}
}
