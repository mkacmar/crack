package elf

import (
	stdelf "debug/elf"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/binary/elf"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
)

// ASLRRuleID is the rule ID for ASLR compatibility.
const ASLRRuleID = "aslr"

// ASLRRule checks if binary is ASLR compatible.
//
// References:
//   - https://github.com/torvalds/linux/blob/master/Documentation/admin-guide/sysctl/kernel.rst
type ASLRRule struct{}

func (r ASLRRule) ID() string   { return ASLRRuleID }
func (r ASLRRule) Name() string { return "ASLR Compatibility" }
func (r ASLRRule) Description() string {
	return "Checks if the binary is compatible with Address Space Layout Randomization (ASLR). ASLR randomizes memory addresses at runtime, making it difficult for attackers to predict the location of code and data. This checks binary compatibility only, not system ASLR settings."
}

func (r ASLRRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM | binary.ArchRISCV},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, DefaultVersion: toolchain.Version{Major: 6, Minor: 1}, Flag: "-fPIE -pie -z noexecstack"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, DefaultVersion: toolchain.Version{Major: 4, Minor: 0}, Flag: "-fPIE -pie -z noexecstack"},
		},
		LibC: binary.LibCAll,
	}
}

func (r ASLRRule) Execute(bin elf.Binary) rule.Result {
	switch bin.Type() {
	case stdelf.ET_EXEC:
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: "Not ASLR compatible, not PIE",
		}
	case stdelf.ET_DYN:
	default:
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	isPIE, err := elf.HasDynFlag(bin, stdelf.DT_FLAGS_1, uint64(stdelf.DF_1_PIE))
	if err != nil {
		return rule.Skip("failed to read dynamic section", err)
	}
	if !isPIE {
		for _, prog := range bin.Progs() {
			if prog.Type == stdelf.PT_INTERP {
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
	for _, prog := range bin.Progs() {
		if prog.Type == stdelf.PT_GNU_STACK {
			hasNXStack = (prog.Flags & stdelf.PF_X) == 0
			break
		}
	}

	if !hasNXStack {
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: "Not ASLR compatible, executable stack",
		}
	}

	textrel, err := elf.HasDynTag(bin, stdelf.DT_TEXTREL)
	if err != nil {
		return rule.Skip("failed to read dynamic section", err)
	}
	if textrel {
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
