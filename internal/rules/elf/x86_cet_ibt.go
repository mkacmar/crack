package elf

import (
	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const X86CETIBTRuleID = "x86-cet-ibt"

// X86CETIBTRule checks for CET Indirect Branch Tracking (Intel/AMD)
// GCC: https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html#index-fcf-protection
type X86CETIBTRule struct{}

func (r X86CETIBTRule) ID() string   { return X86CETIBTRuleID }
func (r X86CETIBTRule) Name() string { return "x86 CET - Indirect Branch Tracking" }
func (r X86CETIBTRule) Description() string {
	return "Checks for Intel Control-flow Enforcement Technology Indirect Branch Tracking (CET-IBT). IBT requires indirect branches to land on ENDBR instructions, preventing attackers from redirecting indirect calls and jumps to arbitrary code. Invalid branch targets trigger a control protection exception."
}

func (r X86CETIBTRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAllX86,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 8, Minor: 1}, Flag: "-fcf-protection=full"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 10, Minor: 0}, Flag: "-fcf-protection=full"},
		},
	}
}

func (r X86CETIBTRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	hasIBT := parseGNUProperty(bin.File, GNU_PROPERTY_X86_FEATURE_1_AND, GNU_PROPERTY_X86_FEATURE_1_IBT)

	if hasIBT {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "CET IBT enabled",
		}
	}
	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "CET IBT not enabled",
	}
}
