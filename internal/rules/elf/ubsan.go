package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// UBSanRule checks for Undefined Behavior Sanitizer
// Clang: https://clang.llvm.org/docs/UndefinedBehaviorSanitizer.html
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-fsanitize=undefined
type UBSanRule struct{}

func (r UBSanRule) ID() string                     { return "ubsan" }
func (r UBSanRule) Name() string                   { return "Undefined Behavior Sanitizer" }
func (r UBSanRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r UBSanRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r UBSanRule) TargetArch() model.Architecture { return model.ArchAll }
func (r UBSanRule) HasPerfImpact() bool            { return true }

func (r UBSanRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 4, Minor: 9}, Flag: "-fsanitize=undefined"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 3}, Flag: "-fsanitize=undefined"},
		},
	}
}

func (r UBSanRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasUBSan := false

	symbols, err := f.Symbols()
	if err != nil {
		symbols = nil
	}

	dynsyms, err := f.DynamicSymbols()
	if err != nil {
		dynsyms = nil
	}

	for _, sym := range symbols {
		if strings.HasPrefix(sym.Name, "__ubsan_") {
			hasUBSan = true
			break
		}
	}

	if !hasUBSan {
		for _, sym := range dynsyms {
			if strings.HasPrefix(sym.Name, "__ubsan_") {
				hasUBSan = true
				break
			}
		}
	}

	if hasUBSan {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "UBSan is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "UBSan is NOT enabled",
	}
}
