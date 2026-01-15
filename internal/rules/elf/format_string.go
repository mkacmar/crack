package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// FormatStringRule checks for format string protection
// glibc: https://sourceware.org/glibc/wiki/FortifySourceLevel3
type FormatStringRule struct{}

func (r FormatStringRule) ID() string                     { return "format-string" }
func (r FormatStringRule) Name() string                   { return "Format String Protection" }
func (r FormatStringRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r FormatStringRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r FormatStringRule) TargetArch() model.Architecture { return model.ArchAll }
func (r FormatStringRule) HasPerfImpact() bool            { return false }

func (r FormatStringRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 4, Minor: 0}, Flag: "-D_FORTIFY_SOURCE=2 -O1"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-D_FORTIFY_SOURCE=2 -O1"},
		},
	}
}

func (r FormatStringRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasFortify := false
	symbols, err := f.Symbols()
	if err != nil {
		dynsyms, err := f.DynamicSymbols()
		if err == nil {
			symbols = dynsyms
		}
	}

	for _, sym := range symbols {
		if strings.HasSuffix(sym.Name, "_chk") {
			hasFortify = true
			break
		}
	}

	if hasFortify {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Format string protection is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Format string protection is NOT enabled",
	}
}
