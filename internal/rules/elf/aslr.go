package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const ASLRRuleID = "aslr"

// ASLRRule checks if binary is ASLR compatible
// Linux Kernel: https://github.com/torvalds/linux/blob/master/Documentation/admin-guide/sysctl/kernel.rst
type ASLRRule struct{}

func (r ASLRRule) ID() string   { return ASLRRuleID }
func (r ASLRRule) Name() string { return "ASLR Compatibility" }
func (r ASLRRule) Description() string {
	return "Checks if the binary is compatible with Address Space Layout Randomization (ASLR). ASLR randomizes memory addresses at runtime, making it difficult for attackers to predict the location of code and data. This checks binary compatibility only, not system ASLR settings."
}

func (r ASLRRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC: {MinVersion: toolchain.Version{Major: 4, Minor: 1}, DefaultVersion: toolchain.Version{Major: 6, Minor: 1}, Flag: "-fPIE -pie -z noexecstack"},
			// -fPIE available since early Clang, -z is linker pass-through
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, DefaultVersion: toolchain.Version{Major: 4, Minor: 0}, Flag: "-fPIE -pie -z noexecstack"},
		},
	}
}

func (r ASLRRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	switch bin.File.Type {
	case elf.ET_EXEC:
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: "Not ASLR compatible, not PIE",
		}
	case elf.ET_DYN:
	default:
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	// ET_DYN can be PIE executable or shared library.
	isPIE := HasDynFlag(bin.File, elf.DT_FLAGS_1, uint64(elf.DF_1_PIE))
	if !isPIE {
		// Check for PT_INTERP as fallback for older binaries.
		for _, prog := range bin.File.Progs {
			if prog.Type == elf.PT_INTERP {
				isPIE = true
				break
			}
		}
	}

	if !isPIE {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Shared library, ASLR not applicable",
		}
	}

	hasNXStack := false
	for _, prog := range bin.File.Progs {
		if prog.Type == elf.PT_GNU_STACK {
			hasNXStack = (prog.Flags & elf.PF_X) == 0
			break
		}
	}

	if !hasNXStack {
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: "Not ASLR compatible, executable stack",
		}
	}

	// Check for text relocations (breaks ASLR).
	if HasDynTag(bin.File, elf.DT_TEXTREL) {
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: "Not ASLR compatible, text relocations present",
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusPassed,
		Message: "ASLR compatible",
	}
}
