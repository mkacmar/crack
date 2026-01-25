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

func (r ARMBranchProtectionRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformARM64v8_5,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 9, Minor: 1}, Flag: "-mbranch-protection=standard"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 8, Minor: 0}, Flag: "-mbranch-protection=standard"},
		},
	}
}

func (r ARMBranchProtectionRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	hasPAC := parseGNUProperty(bin.File, GNU_PROPERTY_AARCH64_FEATURE_1_AND, GNU_PROPERTY_AARCH64_FEATURE_1_PAC)
	hasBTI := parseGNUProperty(bin.File, GNU_PROPERTY_AARCH64_FEATURE_1_AND, GNU_PROPERTY_AARCH64_FEATURE_1_BTI)
	hasBranchProt := hasPAC && hasBTI

	if hasBranchProt {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "ARM branch protection (PAC+BTI) is fully enabled",
		}
	}

	// PAC property requires all linked objects (including libc) to have PAC support
	message := "ARM branch protection is NOT enabled (libc must also be compiled with PAC+BTI)"
	if hasPAC && !hasBTI {
		message = "ARM branch protection is partial (PAC enabled, BTI missing)"
	} else if !hasPAC && hasBTI {
		message = "ARM branch protection is partial (BTI enabled, PAC missing; libc may lack PAC support)"
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: message,
	}
}
