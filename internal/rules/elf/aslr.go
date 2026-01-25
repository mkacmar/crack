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

func (r ASLRRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 6, Minor: 0}, DefaultVersion: toolchain.Version{Major: 6, Minor: 0}, Flag: "-fPIE -pie -z noexecstack"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 0}, DefaultVersion: toolchain.Version{Major: 6, Minor: 0}, Flag: "-fPIE -pie -z noexecstack"},
		},
	}
}

func (r ASLRRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	switch bin.File.Type {
	case elf.ET_EXEC:
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: "Binary is NOT ASLR compatible (not compiled as PIE)",
		}
	case elf.ET_DYN:
		// Continue with analysis below
	default:
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	// ET_DYN can be PIE executable or shared library
	isPIE := HasDynFlag(bin.File, elf.DT_FLAGS_1, DF_1_PIE)
	if !isPIE {
		// Check for PT_INTERP as fallback for older binaries
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
			Message: "Shared library (ASLR check not applicable)",
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
			Message: "Binary is NOT fully ASLR compatible (executable stack)",
		}
	}

	// Check for text relocations (breaks ASLR)
	if HasDynTag(bin.File, elf.DT_TEXTREL) {
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: "Binary is NOT ASLR compatible (has text relocations)",
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusPassed,
		Message: "Binary is fully ASLR compatible (PIE + NX stack + no text relocations)",
	}
}
