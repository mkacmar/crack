package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// X86CETIBTRule checks for CET Indirect Branch Tracking (Intel/AMD)
// GCC: https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html#index-fcf-protection
type X86CETIBTRule struct{}

func (r X86CETIBTRule) ID() string                     { return "x86-cet-ibt" }
func (r X86CETIBTRule) Name() string                   { return "x86 CET - Indirect Branch Tracking" }
func (r X86CETIBTRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r X86CETIBTRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r X86CETIBTRule) TargetArch() model.Architecture { return model.ArchAllX86 }
func (r X86CETIBTRule) HasPerfImpact() bool            { return false }

func (r X86CETIBTRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 8, Minor: 1}, Flag: "-fcf-protection=full"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 7, Minor: 0}, Flag: "-fcf-protection=full"},
		},
	}
}

func (r X86CETIBTRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasIBT := parseGNUProperty(f, GNU_PROPERTY_X86_FEATURE_1_AND, GNU_PROPERTY_X86_FEATURE_1_IBT)

	if hasIBT {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "CET IBT (Indirect Branch Tracking) is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "CET IBT is NOT enabled",
	}
}
