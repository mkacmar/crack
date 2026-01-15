package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// PIERule checks if binary is compiled as Position Independent Executable
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Code-Gen-Options.html#index-fPIE
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-fpie
type PIERule struct{}

func (r PIERule) ID() string                     { return "pie" }
func (r PIERule) Name() string                   { return "Position Independent Executable" }
func (r PIERule) Format() model.BinaryFormat     { return model.FormatELF }
func (r PIERule) FlagType() model.FlagType       { return model.FlagTypeBoth }
func (r PIERule) TargetArch() model.Architecture { return model.ArchAll }
func (r PIERule) HasPerfImpact() bool            { return false }

func (r PIERule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 6, Minor: 0}, DefaultVersion: model.Version{Major: 6, Minor: 0}, Flag: "-fPIE -pie"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, DefaultVersion: model.Version{Major: 6, Minor: 0}, Flag: "-fPIE -pie"},
		},
	}
}

func (r PIERule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	isPIE := f.Type == elf.ET_DYN

	if isPIE {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Binary is compiled as PIE (enables ASLR when system supports it)",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Binary is NOT compiled as PIE (ASLR not possible)",
	}
}
