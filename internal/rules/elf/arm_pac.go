package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// ARMPACRule checks for ARM Pointer Authentication Code
// ARM: https://developer.arm.com/documentation/ddi0487/latest
// GCC: https://gcc.gnu.org/onlinedocs/gcc/AArch64-Options.html#index-mbranch-protection
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-mbranch-protection
type ARMPACRule struct{}

func (r ARMPACRule) ID() string                 { return "arm-pac" }
func (r ARMPACRule) Name() string               { return "ARM Pointer Authentication" }

func (r ARMPACRule) Applicability() model.Applicability {
	return model.Applicability{
		Arch: model.ArchARM64,
		Compilers: map[model.Compiler]model.CompilerRequirement{
			model.CompilerGCC:   {MinVersion: model.Version{Major: 9, Minor: 1}, Flag: "-mbranch-protection=pac-ret"},
			model.CompilerClang: {MinVersion: model.Version{Major: 8, Minor: 0}, Flag: "-mbranch-protection=pac-ret"},
		},
	}
}

func (r ARMPACRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasPAC := parseGNUProperty(f, GNU_PROPERTY_AARCH64_FEATURE_1_AND, GNU_PROPERTY_AARCH64_FEATURE_1_PAC)

	if hasPAC {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "ARM PAC (Pointer Authentication Code) is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "ARM PAC is NOT enabled (requires ARMv8.3+ hardware)",
	}
}
