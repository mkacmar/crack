package elf

import (
	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const ARMBTIRuleID = "arm-bti"

// ARMBTIRule checks for ARM Branch Target Identification
type ARMBTIRule struct{}

func (r ARMBTIRule) ID() string   { return ARMBTIRuleID }
func (r ARMBTIRule) Name() string { return "ARM Branch Target Identification" }
func (r ARMBTIRule) Description() string {
	return "Checks for ARM Branch Target Identification (BTI). BTI marks valid indirect branch targets with landing pad instructions, causing the CPU to fault if an indirect branch lands elsewhere. This prevents attackers from redirecting indirect calls and jumps to arbitrary code."
}

func (r ARMBTIRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformARM64v8_5,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			// rustc: nightly-only via -Z branch-protection=bti
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 10, Minor: 1}, Flag: "-mbranch-protection=bti"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 12, Minor: 0}, Flag: "-mbranch-protection=bti"},
		},
	}
}

func (r ARMBTIRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	hasBTI := parseGNUProperty(bin.File, GNU_PROPERTY_AARCH64_FEATURE_1_AND, GNU_PROPERTY_AARCH64_FEATURE_1_BTI)

	if hasBTI {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "ARM BTI enabled",
		}
	}
	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "ARM BTI not enabled",
	}
}
