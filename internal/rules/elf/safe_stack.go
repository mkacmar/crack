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
	// SafeStack is an LLVM feature for C/C++ - not applicable to standard Go/Rust toolchains
	if info != nil {
		switch info.Language {
		case model.LangGo:
			return model.RuleResult{
				State:   model.CheckStateSkipped,
				Message: "Standard Go toolchain (SafeStack not applicable)",
			}
		case model.LangRust:
			return model.RuleResult{
				State:   model.CheckStateSkipped,
				Message: "Standard Rust toolchain (has compile-time memory safety)",
			}
		}
	}

	symbols, _ := f.Symbols()
	dynsyms, _ := f.DynamicSymbols()

	for _, sym := range symbols {
		if strings.HasPrefix(sym.Name, "__safestack_") {
			return model.RuleResult{
				State:   model.CheckStatePassed,
				Message: "SafeStack is enabled",
			}
		}
	}

	for _, sym := range dynsyms {
		if strings.HasPrefix(sym.Name, "__safestack_") {
			return model.RuleResult{
				State:   model.CheckStatePassed,
				Message: "SafeStack is enabled",
			}
		}
	}

	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "SafeStack is NOT enabled",
	}
}
