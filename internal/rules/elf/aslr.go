package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// ASLRRule checks if binary is ASLR compatible
// Linux Kernel: https://github.com/torvalds/linux/blob/master/Documentation/admin-guide/sysctl/kernel.rst
type ASLRRule struct{}

func (r ASLRRule) ID() string                     { return "aslr" }
func (r ASLRRule) Name() string                   { return "ASLR Compatibility" }
func (r ASLRRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r ASLRRule) FlagType() model.FlagType       { return model.FlagTypeBoth }
func (r ASLRRule) TargetArch() model.Architecture { return model.ArchAll }
func (r ASLRRule) HasPerfImpact() bool            { return false }

func (r ASLRRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 6, Minor: 0}, DefaultVersion: model.Version{Major: 6, Minor: 0}, Flag: "-fPIE -pie -z noexecstack"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, DefaultVersion: model.Version{Major: 6, Minor: 0}, Flag: "-fPIE -pie -z noexecstack"},
		},
	}
}

func (r ASLRRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	if f.Type == elf.ET_EXEC {
		return model.RuleResult{
			State:   model.CheckStateFailed,
			Message: "Binary is NOT ASLR compatible (not compiled as PIE)",
		}
	}

	if f.Type != elf.ET_DYN {
		return model.RuleResult{
			State:   model.CheckStateSkipped,
			Message: "Unknown binary type",
		}
	}

	// ET_DYN can be PIE executable or shared library
	isPIE := HasDynFlag(f, elf.DT_FLAGS_1, DF_1_PIE)
	if !isPIE {
		// Check for PT_INTERP as fallback for older binaries
		for _, prog := range f.Progs {
			if prog.Type == elf.PT_INTERP {
				isPIE = true
				break
			}
		}
	}

	if !isPIE {
		return model.RuleResult{
			State:   model.CheckStateSkipped,
			Message: "Shared library (ASLR check not applicable)",
		}
	}

	hasNXStack := false
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_GNU_STACK {
			hasNXStack = (prog.Flags & elf.PF_X) == 0
			break
		}
	}

	if !hasNXStack {
		return model.RuleResult{
			State:   model.CheckStateFailed,
			Message: "Binary is NOT fully ASLR compatible (executable stack)",
		}
	}

	// Check for text relocations (breaks ASLR)
	if HasDynTag(f, elf.DT_TEXTREL) {
		return model.RuleResult{
			State:   model.CheckStateFailed,
			Message: "Binary is NOT ASLR compatible (has text relocations)",
		}
	}

	return model.RuleResult{
		State:   model.CheckStatePassed,
		Message: "Binary is fully ASLR compatible (PIE + NX stack + no text relocations)",
	}
}
