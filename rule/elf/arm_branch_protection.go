package elf

import (
	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// ARMBranchProtectionRuleID is the rule ID for ARM branch protection.
const ARMBranchProtectionRuleID = "arm-branch-protection"

// ARMBranchProtectionRule checks for ARM branch protection (PAC+BTI).
// ARM: https://developer.arm.com/documentation/ddi0487/latest
// GCC: https://gcc.gnu.org/onlinedocs/gcc/AArch64-Options.html#index-mbranch-protection
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-mbranch-protection
type ARMBranchProtectionRule struct{}

func (r ARMBranchProtectionRule) ID() string   { return ARMBranchProtectionRuleID }
func (r ARMBranchProtectionRule) Name() string { return "ARM Branch Protection" }
func (r ARMBranchProtectionRule) Description() string {
	return "Checks for ARM branch protection (BTI + PAC combined). This enables both Branch Target Identification to validate indirect branch targets and Pointer Authentication to sign return addresses."
}

func (r ARMBranchProtectionRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformARM64v85,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 10, Minor: 1}, Flag: "-mbranch-protection=standard"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 12, Minor: 0}, Flag: "-mbranch-protection=standard"},
		},
	}
}

func (r ARMBranchProtectionRule) Execute(bin *binary.ELFBinary) rule.Result {
	hasPAC := bin.HasGNUProperty(binary.GNU_PROPERTY_AARCH64_FEATURE_1_AND, binary.GNU_PROPERTY_AARCH64_FEATURE_1_PAC)
	hasBTI := bin.HasGNUProperty(binary.GNU_PROPERTY_AARCH64_FEATURE_1_AND, binary.GNU_PROPERTY_AARCH64_FEATURE_1_BTI)

	switch {
	case hasPAC && hasBTI:
		return rule.Result{
			Status:  rule.StatusPassed,
			Message: "ARM branch protection enabled (PAC+BTI)",
		}
	case hasPAC:
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: "ARM branch protection partial, BTI missing",
		}
	case hasBTI:
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: "ARM branch protection partial, PAC missing (libc may lack PAC support)",
		}
	default:
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: "ARM branch protection not enabled",
		}
	}
}
