package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// IntelCETShadowStackRule checks for Intel CET Shadow Stack
// Intel: https://www.intel.com/content/www/us/en/developer/articles/technical/technical-look-control-flow-enforcement-technology.html
// GCC: https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html#index-fcf-protection
type IntelCETShadowStackRule struct{}

func (r IntelCETShadowStackRule) ID() string                     { return "intel-cet-shstk" }
func (r IntelCETShadowStackRule) Name() string                   { return "Intel CET - Shadow Stack" }
func (r IntelCETShadowStackRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r IntelCETShadowStackRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r IntelCETShadowStackRule) TargetArch() model.Architecture { return model.ArchAllX86 }
func (r IntelCETShadowStackRule) HasPerfImpact() bool            { return false }

func (r IntelCETShadowStackRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 8, Minor: 1}, Flag: "-fcf-protection=full"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 7, Minor: 0}, Flag: "-fcf-protection=full"},
		},
	}
}

func (r IntelCETShadowStackRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasShadowStack := parseGNUPropertyForX86Feature(f, GNU_PROPERTY_X86_FEATURE_1_SHSTK)

	if hasShadowStack {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Intel CET Shadow Stack is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Intel CET Shadow Stack is NOT enabled (requires Intel CET-capable CPU)",
	}
}
