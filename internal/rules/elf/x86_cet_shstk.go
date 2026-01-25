package elf

import (
	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const X86CETShadowStackRuleID = "x86-cet-shstk"

// X86CETShadowStackRule checks for CET Shadow Stack (Intel/AMD)
// GCC: https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html#index-fcf-protection
type X86CETShadowStackRule struct{}

func (r X86CETShadowStackRule) ID() string   { return X86CETShadowStackRuleID }
func (r X86CETShadowStackRule) Name() string { return "x86 CET - Shadow Stack" }

func (r X86CETShadowStackRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAllX86,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 8, Minor: 1}, Flag: "-fcf-protection=full"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 7, Minor: 0}, Flag: "-fcf-protection=full"},
		},
	}
}

func (r X86CETShadowStackRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	hasShadowStack := parseGNUProperty(bin.File, GNU_PROPERTY_X86_FEATURE_1_AND, GNU_PROPERTY_X86_FEATURE_1_SHSTK)

	if hasShadowStack {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "CET Shadow Stack is enabled",
		}
	}
	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "CET Shadow Stack is NOT enabled",
	}
}
