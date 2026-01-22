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

func (r SafeStackRule) ID() string                 { return "safe-stack" }
func (r SafeStackRule) Name() string               { return "SafeStack" }
func (r SafeStackRule) Format() model.BinaryFormat { return model.FormatELF }
func (r SafeStackRule) FlagType() model.FlagType   { return model.FlagTypeCompile }
func (r SafeStackRule) HasPerfImpact() bool        { return true }

func (r SafeStackRule) Applicability() model.Applicability {
	return model.Applicability{
		Arch: model.ArchAll,
		Compilers: map[model.Compiler]model.CompilerRequirement{
			model.CompilerClang: {MinVersion: model.Version{Major: 3, Minor: 8}, Flag: "-fsanitize=safe-stack"},
		},
	}
}

func (r SafeStackRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {

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
