package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// ASLRRuleID is the rule ID for ASLR compatibility.
const ASLRRuleID = "aslr"

// ASLRRule checks if binary is ASLR compatible.
// Linux Kernel: https://github.com/torvalds/linux/blob/master/Documentation/admin-guide/sysctl/kernel.rst
type ASLRRule struct{}

func (r ASLRRule) ID() string   { return ASLRRuleID }
func (r ASLRRule) Name() string { return "ASLR Compatibility" }
func (r ASLRRule) Description() string {
	return "Checks if the binary is compatible with Address Space Layout Randomization (ASLR). ASLR randomizes memory addresses at runtime, making it difficult for attackers to predict the location of code and data. This checks binary compatibility only, not system ASLR settings."
}

func (r ASLRRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, DefaultVersion: toolchain.Version{Major: 6, Minor: 1}, Flag: "-fPIE -pie -z noexecstack"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, DefaultVersion: toolchain.Version{Major: 4, Minor: 0}, Flag: "-fPIE -pie -z noexecstack"},
			toolchain.Rustc: {MinVersion: toolchain.Version{Major: 1, Minor: 0}, DefaultVersion: toolchain.Version{Major: 1, Minor: 26}, Flag: "-C relocation-model=pie"},
		},
	}
}

func (r ASLRRule) Execute(bin *binary.ELFBinary) rule.Result {
	switch bin.Type {
	case elf.ET_EXEC:
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: "Not ASLR compatible, not PIE",
		}
	case elf.ET_DYN:
	default:
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	isPIE := bin.HasDynFlag(elf.DT_FLAGS_1, uint64(elf.DF_1_PIE))
	if !isPIE {
		for _, prog := range bin.Progs {
			if prog.Type == elf.PT_INTERP {
				isPIE = true
				break
			}
		}
	}

	if !isPIE {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Shared library, ASLR not applicable",
		}
	}

	hasNXStack := false
	for _, prog := range bin.Progs {
		if prog.Type == elf.PT_GNU_STACK {
			hasNXStack = (prog.Flags & elf.PF_X) == 0
			break
		}
	}

	if !hasNXStack {
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: "Not ASLR compatible, executable stack",
		}
	}

	if bin.HasDynTag(elf.DT_TEXTREL) {
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: "Not ASLR compatible, text relocations present",
		}
	}

	return rule.Result{
		Status:  rule.StatusPassed,
		Message: "ASLR compatible",
	}
}
