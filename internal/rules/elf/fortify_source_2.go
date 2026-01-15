package elf

import (
	"debug/elf"
	"fmt"

	"github.com/mkacmar/crack/internal/model"
)

// FortifyLevel2Rule checks for FORTIFY_SOURCE level 2 or higher
// glibc: https://sourceware.org/glibc/wiki/FortifySourceLevel3
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-D_FORTIFY_SOURCE
type FortifyLevel2Rule struct{}

func (r FortifyLevel2Rule) ID() string                     { return "fortify-source-2" }
func (r FortifyLevel2Rule) Name() string                   { return "FORTIFY_SOURCE Level 2+" }
func (r FortifyLevel2Rule) Format() model.BinaryFormat     { return model.FormatELF }
func (r FortifyLevel2Rule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r FortifyLevel2Rule) TargetArch() model.Architecture { return model.ArchAll }
func (r FortifyLevel2Rule) HasPerfImpact() bool            { return false }

func (r FortifyLevel2Rule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 4, Minor: 0}, Flag: "-D_FORTIFY_SOURCE=2 -O1"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 5, Minor: 0}, Flag: "-D_FORTIFY_SOURCE=2 -O1"},
		},
	}
}

func (r FortifyLevel2Rule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	symbols, _ := f.Symbols()
	dynsyms, _ := f.DynamicSymbols()

	fortifyLevel := 0
	fortified2Count := 0
	fortified3Count := 0

	for _, sym := range symbols {
		if fortified3Count == 0 && (sym.Name == "__chk_fail" || sym.Name == "__fortify_fail") {
			fortified3Count++
		}
		if sym.Name == "__sprintf_chk" || sym.Name == "__snprintf_chk" {
			fortified2Count++
		}
	}
	for _, sym := range dynsyms {
		if fortified3Count == 0 && (sym.Name == "__chk_fail" || sym.Name == "__fortify_fail") {
			fortified3Count++
		}
		if sym.Name == "__sprintf_chk" || sym.Name == "__snprintf_chk" {
			fortified2Count++
		}
	}

	if fortified3Count > 0 {
		fortifyLevel = 3
	} else if fortified2Count > 0 {
		fortifyLevel = 2
	}

	if fortifyLevel >= 2 {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: fmt.Sprintf("FORTIFY_SOURCE level %d is enabled", fortifyLevel),
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "FORTIFY_SOURCE level 2+ is NOT enabled",
	}
}
