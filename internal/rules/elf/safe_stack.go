package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// SafeStackRule checks for SafeStack protection
// Clang: https://clang.llvm.org/docs/SafeStack.html
// LLVM: https://llvm.org/docs/SafeStack.html
type SafeStackRule struct{}

func (r SafeStackRule) ID() string                     { return "safe-stack" }
func (r SafeStackRule) Name() string                   { return "SafeStack" }
func (r SafeStackRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r SafeStackRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r SafeStackRule) TargetArch() model.Architecture { return model.ArchAll }
func (r SafeStackRule) HasPerfImpact() bool            { return true }

func (r SafeStackRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 8}, Flag: "-fsanitize=safe-stack"},
		},
	}
}

func (r SafeStackRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasSafeStack := false

	symbols, err := f.Symbols()
	if err != nil {
		symbols = nil
	}

	dynsyms, err := f.DynamicSymbols()
	if err != nil {
		dynsyms = nil
	}

	for _, sym := range symbols {
		if strings.Contains(sym.Name, "__safestack") {
			hasSafeStack = true
			break
		}
	}

	if !hasSafeStack {
		for _, sym := range dynsyms {
			if strings.Contains(sym.Name, "__safestack") {
				hasSafeStack = true
				break
			}
		}
	}

	if hasSafeStack {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "SafeStack is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "SafeStack is NOT enabled",
	}
}
