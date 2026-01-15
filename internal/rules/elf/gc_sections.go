package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// GCSectionsRule checks for function sections (dead code elimination)
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Optimize-Options.html#index-ffunction-sections
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-_002d_002dgc_002dsections
type GCSectionsRule struct{}

func (r GCSectionsRule) ID() string                     { return "gc-sections" }
func (r GCSectionsRule) Name() string                   { return "Function Sections" }
func (r GCSectionsRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r GCSectionsRule) FlagType() model.FlagType       { return model.FlagTypeBoth }
func (r GCSectionsRule) TargetArch() model.Architecture { return model.ArchAll }
func (r GCSectionsRule) HasPerfImpact() bool            { return false }

func (r GCSectionsRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 2, Minor: 0}, Flag: "-ffunction-sections -fdata-sections -Wl,--gc-sections"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-ffunction-sections -fdata-sections -Wl,--gc-sections"},
		},
	}
}

func (r GCSectionsRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasFunctionSections := false

	for _, sec := range f.Sections {
		if sec.Name == ".text.unlikely" ||
			sec.Name == ".text.hot" ||
			sec.Name == ".text.startup" ||
			sec.Name == ".text.exit" {
			hasFunctionSections = true
			break
		}
	}

	if hasFunctionSections {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Function sections detected (enables dead code elimination)",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Function sections NOT detected (consider -ffunction-sections -fdata-sections -Wl,--gc-sections)",
	}
}
