package elf

import (
	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const ARMBranchProtectionRuleID = "arm-branch-protection"

// ARMBranchProtectionRule checks for ARM branch protection (PAC+BTI)
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
		Platform: binary.PlatformARM64v8_5,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			// rustc: nightly-only via -Z branch-protection=bti,pac-ret
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 10, Minor: 1}, Flag: "-mbranch-protection=standard"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 12, Minor: 0}, Flag: "-mbranch-protection=standard"},
		},
	}
}

func (r ARMBranchProtectionRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	hasPAC := parseGNUProperty(bin.File, GNU_PROPERTY_AARCH64_FEATURE_1_AND, GNU_PROPERTY_AARCH64_FEATURE_1_PAC)
	hasBTI := parseGNUProperty(bin.File, GNU_PROPERTY_AARCH64_FEATURE_1_AND, GNU_PROPERTY_AARCH64_FEATURE_1_BTI)

	switch {
	case hasPAC && hasBTI:
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "ARM branch protection enabled (PAC+BTI)",
		}
	case hasPAC:
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: "ARM branch protection partial, BTI missing",
		}
	case hasBTI:
		// PAC requires all linked objects (including libc) to have PAC support.
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: "ARM branch protection partial, PAC missing (libc may lack PAC support)",
		}
	default:
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: "ARM branch protection not enabled",
		}
	}
}
