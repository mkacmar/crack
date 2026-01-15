package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// IntegerOverflowRule checks for integer overflow protection
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-ftrapv
// Clang: https://clang.llvm.org/docs/UndefinedBehaviorSanitizer.html#available-checks
type IntegerOverflowRule struct{}

func (r IntegerOverflowRule) ID() string                     { return "integer-overflow" }
func (r IntegerOverflowRule) Name() string                   { return "Integer Overflow Protection" }
func (r IntegerOverflowRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r IntegerOverflowRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r IntegerOverflowRule) TargetArch() model.Architecture { return model.ArchAll }
func (r IntegerOverflowRule) HasPerfImpact() bool            { return true }

func (r IntegerOverflowRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 4, Minor: 9}, Flag: "-ftrapv"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 3}, Flag: "-fsanitize=signed-integer-overflow,unsigned-integer-overflow"},
		},
	}
}

func (r IntegerOverflowRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasIntOverflow := false
	symbols, err := f.Symbols()
	if err != nil {
		dynsyms, err := f.DynamicSymbols()
		if err == nil {
			symbols = dynsyms
		}
	}

	for _, sym := range symbols {
		if strings.Contains(sym.Name, "__ubsan_") {
			hasIntOverflow = true
			break
		}
	}

	if hasIntOverflow {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Integer overflow protection is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Integer overflow protection is NOT enabled",
	}
}
