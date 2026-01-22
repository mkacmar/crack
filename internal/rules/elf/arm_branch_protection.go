package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

// ARMBranchProtectionRule checks for ARM branch protection (PAC+BTI)
// ARM: https://developer.arm.com/documentation/ddi0487/latest
// GCC: https://gcc.gnu.org/onlinedocs/gcc/AArch64-Options.html#index-mbranch-protection
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-mbranch-protection
type ARMBranchProtectionRule struct{}

func (r ARMBranchProtectionRule) ID() string   { return "arm-branch-protection" }
func (r ARMBranchProtectionRule) Name() string { return "ARM Branch Protection" }

func (r ARMBranchProtectionRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Arch: binary.ArchARM64,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 9, Minor: 1}, Flag: "-mbranch-protection=standard"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 8, Minor: 0}, Flag: "-mbranch-protection=standard"},
		},
	}
}

func (r ARMBranchProtectionRule) Execute(f *elf.File, info *binary.Parsed) rule.Result {
	hasPAC := parseGNUProperty(f, GNU_PROPERTY_AARCH64_FEATURE_1_AND, GNU_PROPERTY_AARCH64_FEATURE_1_PAC)
	hasBTI := parseGNUProperty(f, GNU_PROPERTY_AARCH64_FEATURE_1_AND, GNU_PROPERTY_AARCH64_FEATURE_1_BTI)
	hasBranchProt := hasPAC && hasBTI

	if hasBranchProt {
		return rule.Result{
			State:   rule.CheckStatePassed,
			Message: "ARM branch protection (PAC+BTI) is fully enabled",
		}
	}

	message := "ARM branch protection is NOT enabled (requires ARMv8.3+ hardware)"
	if hasPAC && !hasBTI {
		message = "ARM branch protection is partial (PAC enabled, BTI missing, requires ARMv8.3+ hardware)"
	} else if !hasPAC && hasBTI {
		message = "ARM branch protection is partial (BTI enabled, PAC missing, requires ARMv8.3+ hardware)"
	}

	return rule.Result{
		State:   rule.CheckStateFailed,
		Message: message,
	}
}
