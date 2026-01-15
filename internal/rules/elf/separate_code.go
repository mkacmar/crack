package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

const PT_GNU_SEPARATE_CODE = 0x65041580

// SeparateCodeRule checks if code and data are in separate segments
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type SeparateCodeRule struct{}

func (r SeparateCodeRule) ID() string                     { return "separate-code" }
func (r SeparateCodeRule) Name() string                   { return "Separate Code Segments" }
func (r SeparateCodeRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r SeparateCodeRule) FlagType() model.FlagType       { return model.FlagTypeLink }
func (r SeparateCodeRule) TargetArch() model.Architecture { return model.ArchAll }
func (r SeparateCodeRule) HasPerfImpact() bool            { return false }

func (r SeparateCodeRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,separate-code"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,separate-code"},
		},
	}
}

func (r SeparateCodeRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasSeparateCode := false

	for _, prog := range f.Progs {
		if prog.Type == elf.ProgType(PT_GNU_SEPARATE_CODE) {
			hasSeparateCode = true
			break
		}
	}

	if hasSeparateCode {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Code and data are in separate segments",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Code and data segments are NOT properly separated",
	}
}
