package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// ARMBTIRule checks for ARM Branch Target Identification
type ARMBTIRule struct{}

func (r ARMBTIRule) ID() string                 { return "arm-bti" }
func (r ARMBTIRule) Name() string               { return "ARM Branch Target Identification" }
func (r ARMBTIRule) Format() model.BinaryFormat { return model.FormatELF }

func (r ARMBTIRule) Applicability() model.Applicability {
	return model.Applicability{
		Arch: model.ArchARM64,
		Compilers: map[model.Compiler]model.CompilerRequirement{
			model.CompilerGCC:   {MinVersion: model.Version{Major: 9, Minor: 1}, Flag: "-mbranch-protection=bti"},
			model.CompilerClang: {MinVersion: model.Version{Major: 8, Minor: 0}, Flag: "-mbranch-protection=bti"},
		},
	}
}

func (r ARMBTIRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasBTI := parseGNUProperty(f, GNU_PROPERTY_AARCH64_FEATURE_1_AND, GNU_PROPERTY_AARCH64_FEATURE_1_BTI)

	if hasBTI {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "ARM BTI (Branch Target Identification) is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "ARM BTI is NOT enabled (requires ARMv8.5+ hardware)",
	}
}
