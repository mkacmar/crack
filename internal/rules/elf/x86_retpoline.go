package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// RetpolineRule checks for Spectre v2 mitigation (retpoline)
// GCC: https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html#index-mindirect-branch
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-mretpoline
type RetpolineRule struct{}

func (r RetpolineRule) ID() string                     { return "x86-retpoline" }
func (r RetpolineRule) Name() string                   { return "x86 Retpoline (Spectre v2)" }
func (r RetpolineRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r RetpolineRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r RetpolineRule) TargetArch() model.Architecture { return model.ArchAllX86 }
func (r RetpolineRule) HasPerfImpact() bool            { return true }

func (r RetpolineRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 7, Minor: 3}, Flag: "-mindirect-branch=thunk-extern -mfunction-return=thunk-extern"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 5, Minor: 0}, Flag: "-mretpoline"},
		},
	}
}

func (r RetpolineRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {

	symbols, _ := f.Symbols()
	dynsyms, _ := f.DynamicSymbols()

	hasRetpoline := false
	hasIndirectBranchThunk := false
	hasFunctionReturnThunk := false

	for _, sym := range symbols {
		if strings.Contains(sym.Name, "__x86_indirect_thunk") {
			hasRetpoline = true
			hasIndirectBranchThunk = true
		}
		if strings.Contains(sym.Name, "__x86_return_thunk") {
			hasRetpoline = true
			hasFunctionReturnThunk = true
		}
	}
	for _, sym := range dynsyms {
		if strings.Contains(sym.Name, "__x86_indirect_thunk") {
			hasRetpoline = true
			hasIndirectBranchThunk = true
		}
		if strings.Contains(sym.Name, "__x86_return_thunk") {
			hasRetpoline = true
			hasFunctionReturnThunk = true
		}
	}

	if hasRetpoline {
		msg := "Retpoline enabled"
		if hasIndirectBranchThunk {
			msg = "Retpoline enabled (indirect branch thunks detected)"
		} else if hasFunctionReturnThunk {
			msg = "Return thunks detected (partial Spectre mitigation)"
		}
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: msg,
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Retpoline is NOT enabled (x86 Spectre v2 mitigation missing)",
	}
}
