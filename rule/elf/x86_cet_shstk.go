package elf

import (
	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// X86CETShadowStackRuleID is the rule ID for CET Shadow Stack.
const X86CETShadowStackRuleID = "x86-cet-shstk"

// X86CETShadowStackRule checks for CET Shadow Stack (Intel/AMD).
// GCC: https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html#index-fcf-protection
type X86CETShadowStackRule struct{}

func (r X86CETShadowStackRule) ID() string   { return X86CETShadowStackRuleID }
func (r X86CETShadowStackRule) Name() string { return "x86 CET - Shadow Stack" }
func (r X86CETShadowStackRule) Description() string {
	return "Checks for Intel Control-flow Enforcement Technology Shadow Stack (CET-SS). Shadow Stack maintains a hardware-protected copy of return addresses, detecting ROP attacks when the shadow and regular stacks diverge on function return."
}

func (r X86CETShadowStackRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAllX86,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 8, Minor: 1}, Flag: "-fcf-protection=full"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 10, Minor: 0}, Flag: "-fcf-protection=full"},
		},
	}
}

func (r X86CETShadowStackRule) Execute(bin *binary.ELFBinary) rule.Result {
	hasShadowStack := bin.HasGNUProperty(binary.GNU_PROPERTY_X86_FEATURE_1_AND, binary.GNU_PROPERTY_X86_FEATURE_1_SHSTK)

	if hasShadowStack {
		return rule.Result{
			Status:  rule.StatusPassed,
			Message: "CET Shadow Stack enabled",
		}
	}
	return rule.Result{
		Status:  rule.StatusFailed,
		Message: "CET Shadow Stack not enabled",
	}
}
