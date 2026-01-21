package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// X86CETShadowStackRule checks for CET Shadow Stack (Intel/AMD)
// GCC: https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html#index-fcf-protection
type X86CETShadowStackRule struct{}

func (r X86CETShadowStackRule) ID() string                     { return "x86-cet-shstk" }
func (r X86CETShadowStackRule) Name() string                   { return "x86 CET - Shadow Stack" }
func (r X86CETShadowStackRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r X86CETShadowStackRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r X86CETShadowStackRule) TargetArch() model.Architecture { return model.ArchAllX86 }
func (r X86CETShadowStackRule) HasPerfImpact() bool            { return false }

func (r X86CETShadowStackRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 8, Minor: 1}, Flag: "-fcf-protection=full"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 7, Minor: 0}, Flag: "-fcf-protection=full"},
		},
	}
}

func (r X86CETShadowStackRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasShadowStack := parseGNUProperty(f, GNU_PROPERTY_X86_FEATURE_1_AND, GNU_PROPERTY_X86_FEATURE_1_SHSTK)

	if hasShadowStack {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "CET Shadow Stack is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "CET Shadow Stack is NOT enabled",
	}
}
