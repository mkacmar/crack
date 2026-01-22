package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

// NXBitRule checks for non-executable stack
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Link-Options.html#index-z
type NXBitRule struct{}

func (r NXBitRule) ID() string   { return "nx-bit" }
func (r NXBitRule) Name() string { return "Non-Executable Stack" }

func (r NXBitRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Arch: binary.ArchAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 3, Minor: 0}, DefaultVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-z noexecstack"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 0}, DefaultVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-z noexecstack"},
		},
	}
}

func (r NXBitRule) Execute(f *elf.File, info *binary.Parsed) rule.Result {
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_GNU_STACK {
			if (prog.Flags & elf.PF_X) == 0 {
				return rule.Result{
					State:   rule.CheckStatePassed,
					Message: "Stack is marked as non-executable (NX bit enabled)",
				}
			}
			return rule.Result{
				State:   rule.CheckStateFailed,
				Message: "Stack is EXECUTABLE (NX bit NOT enabled)",
			}
		}
	}

	// PT_GNU_STACK missing - this typically means executable stack on older systems.
	return rule.Result{
		State:   rule.CheckStateFailed,
		Message: "PT_GNU_STACK segment missing (stack executability depends on system defaults)",
	}
}
