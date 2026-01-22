package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

// ARMBTIRule checks for ARM Branch Target Identification
type ARMBTIRule struct{}

func (r ARMBTIRule) ID() string   { return "arm-bti" }
func (r ARMBTIRule) Name() string { return "ARM Branch Target Identification" }

func (r ARMBTIRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Arch: binary.ArchARM64,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 9, Minor: 1}, Flag: "-mbranch-protection=bti"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 8, Minor: 0}, Flag: "-mbranch-protection=bti"},
		},
	}
}

func (r ARMBTIRule) Execute(f *elf.File, info *binary.Parsed) rule.Result {
	hasBTI := parseGNUProperty(f, GNU_PROPERTY_AARCH64_FEATURE_1_AND, GNU_PROPERTY_AARCH64_FEATURE_1_BTI)

	if hasBTI {
		return rule.Result{
			State:   rule.CheckStatePassed,
			Message: "ARM BTI (Branch Target Identification) is enabled",
		}
	}
	return rule.Result{
		State:   rule.CheckStateFailed,
		Message: "ARM BTI is NOT enabled (requires ARMv8.5+ hardware)",
	}
}
