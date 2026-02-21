package elf

import (
	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
)

// ARMPACRuleID is the rule ID for ARM PAC.
const ARMPACRuleID = "arm-pac"

// ARMPACRule checks for ARM Pointer Authentication Code.
// ARM: https://developer.arm.com/documentation/ddi0487/latest
// GCC: https://gcc.gnu.org/onlinedocs/gcc/AArch64-Options.html#index-mbranch-protection
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-mbranch-protection
type ARMPACRule struct{}

func (r ARMPACRule) ID() string   { return ARMPACRuleID }
func (r ARMPACRule) Name() string { return "ARM Pointer Authentication" }
func (r ARMPACRule) Description() string {
	return "Checks for ARM Pointer Authentication Code (PAC). PAC signs return addresses with a cryptographic key, detecting tampering when the signature is verified on function return. This prevents attackers from overwriting return addresses to hijack control flow."
}

func (r ARMPACRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformARM64v83,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 10, Minor: 1}, Flag: "-mbranch-protection=pac-ret"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 12, Minor: 0}, Flag: "-mbranch-protection=pac-ret"},
		},
	}
}

func (r ARMPACRule) Execute(bin *binary.ELFBinary) rule.Result {
	hasPAC := bin.HasGNUProperty(binary.GNU_PROPERTY_AARCH64_FEATURE_1_AND, binary.GNU_PROPERTY_AARCH64_FEATURE_1_PAC)

	if hasPAC {
		return rule.Result{
			Status:  rule.StatusPassed,
			Message: "ARM PAC enabled",
		}
	}
	return rule.Result{
		Status:  rule.StatusFailed,
		Message: "ARM PAC not enabled",
	}
}
