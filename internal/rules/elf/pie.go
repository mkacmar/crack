package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

const DF_1_PIE = 0x08000000

// PIERule checks if binary is compiled as Position Independent Executable
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Code-Gen-Options.html#index-fPIE
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-fpie
type PIERule struct{}

func (r PIERule) ID() string                 { return "pie" }
func (r PIERule) Name() string               { return "Position Independent Executable" }

func (r PIERule) Applicability() model.Applicability {
	return model.Applicability{
		Arch: model.ArchAll,
		Compilers: map[model.Compiler]model.CompilerRequirement{
			model.CompilerGCC:   {MinVersion: model.Version{Major: 6, Minor: 0}, DefaultVersion: model.Version{Major: 6, Minor: 0}, Flag: "-fPIE -pie"},
			model.CompilerClang: {MinVersion: model.Version{Major: 3, Minor: 0}, DefaultVersion: model.Version{Major: 6, Minor: 0}, Flag: "-fPIE -pie"},
		},
	}
}

func (r PIERule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	if f.Type == elf.ET_EXEC {
		return model.RuleResult{
			State:   model.CheckStateFailed,
			Message: "Binary is NOT compiled as PIE (ASLR not possible)",
		}
	}

	if f.Type != elf.ET_DYN {
		return model.RuleResult{
			State:   model.CheckStateSkipped,
			Message: "Not an executable or shared library",
		}
	}

	// ET_DYN can be either a PIE executable or a shared library.
	// 1. DF_1_PIE flag in .dynamic section - set by modern linkers for PIE executables (including static-pie).
	// 2. PT_INTERP program header - present in dynamically linked executables but not in shared libraries.
	// static-pie binaries (-static-pie) have DF_1_PIE but no PT_INTERP, so the DF_1_PIE check must come first.
	if HasDynFlag(f, elf.DT_FLAGS_1, DF_1_PIE) {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Binary is compiled as PIE (enables ASLR when system supports it)",
		}
	}

	for _, prog := range f.Progs {
		if prog.Type == elf.PT_INTERP {
			return model.RuleResult{
				State:   model.CheckStatePassed,
				Message: "Binary is compiled as PIE (enables ASLR when system supports it)",
			}
		}
	}

	return model.RuleResult{
		State:   model.CheckStateSkipped,
		Message: "Shared library (PIE check not applicable)",
	}
}
