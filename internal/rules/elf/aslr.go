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
	isPIE := f.Type == elf.ET_DYN

	hasNXStack := false
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_GNU_STACK {
			hasNXStack = (prog.Flags & elf.PF_X) == 0
			break
		}
	}

	hasTextRel := false
	dynSec := f.Section(".dynamic")
	if dynSec != nil {
		data, err := dynSec.Data()
		if err == nil {
			tagSize := 8
			if f.Class == elf.ELFCLASS64 {
				tagSize = 16
			}
			for i := 0; i < len(data); i += tagSize {
				if i+tagSize > len(data) {
					break
				}
				var tag uint64
				if f.Class == elf.ELFCLASS64 {
					tag = f.ByteOrder.Uint64(data[i : i+8])
				} else {
					tag = uint64(f.ByteOrder.Uint32(data[i : i+4]))
				}
				if tag == uint64(elf.DT_TEXTREL) {
					hasTextRel = true
					break
				}
				if tag == uint64(elf.DT_NULL) {
					break
				}
			}
		}
	}

	isASLRCompatible := isPIE && hasNXStack && !hasTextRel

	if isASLRCompatible {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Binary is fully ASLR compatible (PIE + NX stack + no text relocations)",
		}
	}

	message := "Binary is NOT ASLR compatible"
	if !isPIE {
		message = "Binary is NOT ASLR compatible (not compiled as PIE)"
	} else if !hasNXStack {
		message = "Binary is NOT fully ASLR compatible (executable stack)"
	} else if hasTextRel {
		message = "Binary is NOT ASLR compatible (has text relocations)"
	}

	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: message,
	}
}
