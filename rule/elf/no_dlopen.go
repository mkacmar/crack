package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// NoDLOpenRuleID is the rule ID for no dlopen.
const NoDLOpenRuleID = "no-dlopen"

// NoDLOpenRule checks if dlopen is disabled.
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type NoDLOpenRule struct{}

func (r NoDLOpenRule) ID() string   { return NoDLOpenRuleID }
func (r NoDLOpenRule) Name() string { return "Disallow dlopen" }
func (r NoDLOpenRule) Description() string {
	return "Checks if the shared library disallows being loaded via dlopen(). This prevents attackers from injecting the library into arbitrary processes."
}

func (r NoDLOpenRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, Flag: "-Wl,-z,nodlopen"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, Flag: "-Wl,-z,nodlopen"},
			toolchain.Rustc: {MinVersion: toolchain.Version{Major: 1, Minor: 0}, Flag: "-C link-arg=-z -C link-arg=nodlopen"},
		},
	}
}

func (r NoDLOpenRule) Execute(bin *binary.ELFBinary) rule.Result {
	if bin.Type != elf.ET_DYN {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not a shared library, dlopen protection not applicable",
		}
	}

	if bin.HasDynFlag(elf.DT_FLAGS_1, uint64(elf.DF_1_PIE)) {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "PIE executable, dlopen protection not applicable",
		}
	}
	for _, prog := range bin.Progs {
		if prog.Type == elf.PT_INTERP {
			return rule.Result{
				Status:  rule.StatusSkipped,
				Message: "PIE executable, dlopen protection not applicable",
			}
		}
	}

	if bin.HasDynFlag(elf.DT_FLAGS_1, uint64(elf.DF_1_NOOPEN)) {
		return rule.Result{
			Status:  rule.StatusPassed,
			Message: "dlopen disabled",
		}
	}

	return rule.Result{
		Status:  rule.StatusFailed,
		Message: "dlopen not disabled",
	}
}
