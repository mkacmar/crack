package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// StackCanaryRule checks for stack canary protection
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-fstack-protector
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-fstack-protector-strong
type StackCanaryRule struct{}

func (r StackCanaryRule) ID() string                     { return "stack-canary" }
func (r StackCanaryRule) Name() string                   { return "Stack Canary Protection" }
func (r StackCanaryRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r StackCanaryRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r StackCanaryRule) TargetArch() model.Architecture { return model.ArchAll }
func (r StackCanaryRule) HasPerfImpact() bool            { return false }

func (r StackCanaryRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 4, Minor: 9}, Flag: "-fstack-protector-strong"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-fstack-protector-strong"},
		},
	}
}

func (r StackCanaryRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	symbols, _ := f.Symbols()
	dynsyms, _ := f.DynamicSymbols()

	for _, sym := range symbols {
		if strings.Contains(sym.Name, "__stack_chk_fail") ||
			strings.Contains(sym.Name, "__stack_smash_handler") ||
			strings.Contains(sym.Name, "__intel_security_cookie") {
			return model.RuleResult{
				State:   model.CheckStatePassed,
				Message: "Stack canary protection is enabled",
			}
		}
	}

	for _, sym := range dynsyms {
		if strings.Contains(sym.Name, "__stack_chk_fail") ||
			strings.Contains(sym.Name, "__stack_smash_handler") ||
			strings.Contains(sym.Name, "__intel_security_cookie") {
			return model.RuleResult{
				State:   model.CheckStatePassed,
				Message: "Stack canary protection is enabled",
			}
		}
	}

	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Stack canary protection is NOT enabled",
	}
}
