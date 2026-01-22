package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

// X86CETIBTRule checks for CET Indirect Branch Tracking (Intel/AMD)
// GCC: https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html#index-fcf-protection
type X86CETIBTRule struct{}

func (r X86CETIBTRule) ID() string   { return "x86-cet-ibt" }
func (r X86CETIBTRule) Name() string { return "x86 CET - Indirect Branch Tracking" }

func (r X86CETIBTRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Arch: binary.ArchAllX86,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 8, Minor: 1}, Flag: "-fcf-protection=full"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 7, Minor: 0}, Flag: "-fcf-protection=full"},
		},
	}
}

func (r X86CETIBTRule) Execute(f *elf.File, info *binary.Parsed) rule.Result {
	hasIBT := parseGNUProperty(f, GNU_PROPERTY_X86_FEATURE_1_AND, GNU_PROPERTY_X86_FEATURE_1_IBT)

	if hasIBT {
		return rule.Result{
			State:   rule.CheckStatePassed,
			Message: "CET IBT (Indirect Branch Tracking) is enabled",
		}
	}
	return rule.Result{
		State:   rule.CheckStateFailed,
		Message: "CET IBT is NOT enabled",
	}
}
