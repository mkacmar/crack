package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// X86CETShadowStackRule checks for CET Shadow Stack (Intel/AMD)
// GCC: https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html#index-fcf-protection
type X86CETShadowStackRule struct{}

func (r X86CETShadowStackRule) ID() string                 { return "x86-cet-shstk" }
func (r X86CETShadowStackRule) Name() string               { return "x86 CET - Shadow Stack" }

func (r X86CETShadowStackRule) Applicability() model.Applicability {
	return model.Applicability{
		Arch: model.ArchAllX86,
		Compilers: map[model.Compiler]model.CompilerRequirement{
			model.CompilerGCC:   {MinVersion: model.Version{Major: 8, Minor: 1}, Flag: "-fcf-protection=full"},
			model.CompilerClang: {MinVersion: model.Version{Major: 7, Minor: 0}, Flag: "-fcf-protection=full"},
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
