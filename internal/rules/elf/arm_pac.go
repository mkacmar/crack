package elf

import (
	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const ARMPACRuleID = "arm-pac"

// ARMPACRule checks for ARM Pointer Authentication Code
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
		Platform: binary.PlatformARM64v8_3,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			// rustc: nightly-only via -Z branch-protection=pac-ret
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 10, Minor: 1}, Flag: "-mbranch-protection=pac-ret"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 12, Minor: 0}, Flag: "-mbranch-protection=pac-ret"},
		},
	}
}

func (r ARMPACRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	hasPAC := parseGNUProperty(bin.File, GNU_PROPERTY_AARCH64_FEATURE_1_AND, GNU_PROPERTY_AARCH64_FEATURE_1_PAC)

	if hasPAC {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "ARM PAC enabled",
		}
	}
	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "ARM PAC not enabled",
	}
}
