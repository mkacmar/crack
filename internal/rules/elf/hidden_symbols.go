package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

const (
	STV_DEFAULT   = 0
	STV_INTERNAL  = 1
	STV_HIDDEN    = 2
	STV_PROTECTED = 3
)

const (
	STT_NOTYPE  = 0
	STT_OBJECT  = 1
	STT_FUNC    = 2
	STT_SECTION = 3
	STT_FILE    = 4
)

// HiddenSymbolsRule checks for hidden symbol visibility
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Code-Gen-Options.html#index-fvisibility
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-fvisibility
type HiddenSymbolsRule struct{}

func (r HiddenSymbolsRule) ID() string                     { return "hidden-symbols" }
func (r HiddenSymbolsRule) Name() string                   { return "Hidden Symbol Visibility" }
func (r HiddenSymbolsRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r HiddenSymbolsRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r HiddenSymbolsRule) TargetArch() model.Architecture { return model.ArchAll }
func (r HiddenSymbolsRule) HasPerfImpact() bool            { return false }

func (r HiddenSymbolsRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 4, Minor: 0}, Flag: "-fvisibility=hidden"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-fvisibility=hidden"},
		},
	}
}

func (r HiddenSymbolsRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	symbols, _ := f.Symbols()
	dynsyms, _ := f.DynamicSymbols()

	visibleCount := 0
	totalCount := 0

	for _, sym := range symbols {
		symType := sym.Info & 0xf
		if symType == STT_OBJECT || symType == STT_FUNC {
			totalCount++
			visibility := sym.Other & 0x3
			if visibility == STV_DEFAULT {
				visibleCount++
			}
		}
	}

	for _, sym := range dynsyms {
		symType := sym.Info & 0xf
		if symType == STT_OBJECT || symType == STT_FUNC {
			totalCount++
			visibility := sym.Other & 0x3
			if visibility == STV_DEFAULT {
				visibleCount++
			}
		}
	}

	hasHiddenVisibility := totalCount > 0 && float64(visibleCount)/float64(totalCount) < 0.5

	if hasHiddenVisibility {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Symbols are mostly hidden by default (good for security)",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Most symbols are visible (consider using -fvisibility=hidden)",
	}
}
