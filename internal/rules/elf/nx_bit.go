package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const NXBitRuleID = "nx-bit"

// NXBitRule checks for non-executable stack
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Link-Options.html#index-z
type NXBitRule struct{}

func (r NXBitRule) ID() string   { return NXBitRuleID }
func (r NXBitRule) Name() string { return "Non-Executable Stack" }

func (r NXBitRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 3, Minor: 0}, DefaultVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-z noexecstack"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 0}, DefaultVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-z noexecstack"},
		},
	}
}

func (r NXBitRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	for _, prog := range bin.File.Progs {
		if prog.Type == elf.PT_GNU_STACK {
			if (prog.Flags & elf.PF_X) == 0 {
				return rule.ExecuteResult{
					Status:  rule.StatusPassed,
					Message: "Stack is marked as non-executable (NX bit enabled)",
				}
			}
			return rule.ExecuteResult{
				Status:  rule.StatusFailed,
				Message: "Stack is EXECUTABLE (NX bit NOT enabled)",
			}
		}
	}

	// PT_GNU_STACK missing - this typically means executable stack on older systems.
	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "PT_GNU_STACK segment missing (stack executability depends on system defaults)",
	}
}
