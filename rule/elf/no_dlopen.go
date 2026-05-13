package elf

import (
	stdelf "debug/elf"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/binary/elf"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
)

// NoDLOpenRuleID is the rule ID for no dlopen.
const NoDLOpenRuleID = "no-dlopen"

// NoDLOpenRule checks if dlopen is disabled.
//
// References:
//   - https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type NoDLOpenRule struct{}

func (r NoDLOpenRule) ID() string   { return NoDLOpenRuleID }
func (r NoDLOpenRule) Name() string { return "Disallow dlopen" }
func (r NoDLOpenRule) Description() string {
	return "Checks if the shared library has the DF_1_NOOPEN flag set to prevent loading via dlopen(3). This restricts an attacker's ability to load the library into arbitrary processes at runtime."
}

func (r NoDLOpenRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM | binary.ArchRISCV},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, Flag: "-Wl,-z,nodlopen"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, Flag: "-Wl,-z,nodlopen"},
		},
		LibC: binary.LibCAll,
	}
}

func (r NoDLOpenRule) Execute(bin elf.Binary) rule.Result {
	if bin.Type() != stdelf.ET_DYN {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not a shared library, dlopen protection not applicable",
		}
	}

	pie, err := elf.HasDynFlag(bin, stdelf.DT_FLAGS_1, uint64(stdelf.DF_1_PIE))
	if err != nil {
		return rule.Skip("failed to read dynamic section", err)
	}
	if pie {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "PIE executable, dlopen protection not applicable",
		}
	}
	for _, prog := range bin.Progs() {
		if prog.Type == stdelf.PT_INTERP {
			return rule.Result{
				Status:  rule.StatusSkipped,
				Message: "PIE executable, dlopen protection not applicable",
			}
		}
	}

	noopen, err := elf.HasDynFlag(bin, stdelf.DT_FLAGS_1, uint64(stdelf.DF_1_NOOPEN))
	if err != nil {
		return rule.Skip("failed to read dynamic section", err)
	}
	if noopen {
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
