package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// ARMMTERule checks for ARM Memory Tagging Extension
// ARM: https://developer.arm.com/documentation/ddi0487/latest
// Clang: https://clang.llvm.org/docs/MemTagSanitizer.html
type ARMMTERule struct{}

func (r ARMMTERule) ID() string                     { return "arm-mte" }
func (r ARMMTERule) Name() string                   { return "ARM Memory Tagging Extension" }
func (r ARMMTERule) Format() model.BinaryFormat     { return model.FormatELF }
func (r ARMMTERule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r ARMMTERule) TargetArch() model.Architecture { return model.ArchARM64 }
func (r ARMMTERule) HasPerfImpact() bool            { return true }

func (r ARMMTERule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 11, Minor: 0}, Flag: "-march=armv8.5-a+memtag -fsanitize=memtag"},
		},
	}
}

func (r ARMMTERule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasMTE := false
	for _, sec := range f.Sections {
		if sec.Name == ".note.memtag" {
			hasMTE = true
			break
		}
	}

	if hasMTE {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "ARM MTE (Memory Tagging Extension) is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "ARM MTE is NOT enabled (requires ARMv8.5+ hardware)",
	}
}
