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
func (r NoDLOpenRule) Description() string {
	return "Checks if the shared library disallows being loaded via dlopen(). This prevents attackers from injecting the library into arbitrary processes."
}

func (r NoDLOpenRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, Flag: "-Wl,-z,nodlopen"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, Flag: "-Wl,-z,nodlopen"},
			toolchain.CompilerRustc: {MinVersion: toolchain.Version{Major: 1, Minor: 0}, Flag: "-C link-arg=-z -C link-arg=nodlopen"},
		},
	}
}

func (r NoDLOpenRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	if bin.File.Type != elf.ET_DYN {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Not a shared library, dlopen protection not applicable",
		}
	}

	if HasDynFlag(bin.File, elf.DT_FLAGS_1, uint64(elf.DF_1_PIE)) {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "PIE executable, dlopen protection not applicable",
		}
	}
	for _, prog := range bin.File.Progs {
		if prog.Type == elf.PT_INTERP {
			return rule.ExecuteResult{
				Status:  rule.StatusSkipped,
				Message: "PIE executable, dlopen protection not applicable",
			}
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
