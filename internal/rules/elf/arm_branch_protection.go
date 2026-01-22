package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// ARMBranchProtectionRule checks for ARM branch protection (PAC+BTI)
// ARM: https://developer.arm.com/documentation/ddi0487/latest
// GCC: https://gcc.gnu.org/onlinedocs/gcc/AArch64-Options.html#index-mbranch-protection
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-mbranch-protection
type ARMBranchProtectionRule struct{}

func (r ARMBranchProtectionRule) ID() string                 { return "arm-branch-protection" }
func (r ARMBranchProtectionRule) Name() string               { return "ARM Branch Protection" }
func (r ARMBranchProtectionRule) Format() model.BinaryFormat { return model.FormatELF }
func (r ARMBranchProtectionRule) FlagType() model.FlagType   { return model.FlagTypeCompile }
func (r ARMBranchProtectionRule) HasPerfImpact() bool        { return false }

func (r ARMBranchProtectionRule) Applicability() model.Applicability {
	return model.Applicability{
		Arch: model.ArchARM64,
		Compilers: map[model.Compiler]model.CompilerRequirement{
			model.CompilerGCC:   {MinVersion: model.Version{Major: 9, Minor: 1}, Flag: "-mbranch-protection=standard"},
			model.CompilerClang: {MinVersion: model.Version{Major: 8, Minor: 0}, Flag: "-mbranch-protection=standard"},
		},
	}
}

func (r ARMBranchProtectionRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasPAC := parseGNUProperty(f, GNU_PROPERTY_AARCH64_FEATURE_1_AND, GNU_PROPERTY_AARCH64_FEATURE_1_PAC)
	hasBTI := parseGNUProperty(f, GNU_PROPERTY_AARCH64_FEATURE_1_AND, GNU_PROPERTY_AARCH64_FEATURE_1_BTI)
	hasBranchProt := hasPAC && hasBTI

	if hasBranchProt {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "ARM branch protection (PAC+BTI) is fully enabled",
		}
	}

	message := "ARM branch protection is NOT enabled (requires ARMv8.3+ hardware)"
	if hasPAC && !hasBTI {
		message = "ARM branch protection is partial (PAC enabled, BTI missing, requires ARMv8.3+ hardware)"
	} else if !hasPAC && hasBTI {
		message = "ARM branch protection is partial (BTI enabled, PAC missing, requires ARMv8.3+ hardware)"
	}

	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: message,
	}
}
