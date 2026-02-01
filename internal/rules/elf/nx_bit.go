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
func (r NXBitRule) Description() string {
	return "Checks if the stack is marked non-executable (NX bit). This prevents stack-based buffer overflow exploits from executing shellcode placed on the stack."
}

func (r NXBitRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 3, Minor: 4}, DefaultVersion: toolchain.Version{Major: 3, Minor: 4}, Flag: "-z noexecstack"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 1, Minor: 0}, DefaultVersion: toolchain.Version{Major: 1, Minor: 0}, Flag: "-z noexecstack"},
		},
	}
}

func (r NXBitRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	if bin.File.Type != elf.ET_EXEC && bin.File.Type != elf.ET_DYN {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	for _, prog := range bin.File.Progs {
		if prog.Type == elf.PT_GNU_STACK {
			if (prog.Flags & elf.PF_X) == 0 {
				return rule.ExecuteResult{
					Status:  rule.StatusPassed,
					Message: "NX enabled, stack non-executable",
				}
			}
			return rule.ExecuteResult{
				Status:  rule.StatusFailed,
				Message: "NX not enabled, stack executable",
			}
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "NX status unknown, no stack segment",
	}
}
