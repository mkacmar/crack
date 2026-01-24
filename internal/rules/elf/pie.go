package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const DF_1_PIE = 0x08000000

// PIERule checks if binary is compiled as Position Independent Executable
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Code-Gen-Options.html#index-fPIE
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-fpie
type PIERule struct{}

func (r PIERule) ID() string   { return "pie" }
func (r PIERule) Name() string { return "Position Independent Executable" }

func (r PIERule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 6, Minor: 0}, DefaultVersion: toolchain.Version{Major: 6, Minor: 0}, Flag: "-fPIE -pie"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 0}, DefaultVersion: toolchain.Version{Major: 6, Minor: 0}, Flag: "-fPIE -pie"},
		},
	}
}

func (r PIERule) Execute(f *elf.File, info *binary.Parsed) rule.ExecuteResult {
	switch f.Type {
	case elf.ET_EXEC:
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: "Binary is NOT compiled as PIE (ASLR not possible)",
		}
	case elf.ET_DYN:
	default:
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	// ET_DYN can be either a PIE executable or a shared library.
	// 1. DF_1_PIE flag in .dynamic section - set by modern linkers for PIE executables (including static-pie).
	// 2. PT_INTERP program header - present in dynamically linked executables but not in shared libraries.
	// static-pie binaries (-static-pie) have DF_1_PIE but no PT_INTERP, so the DF_1_PIE check must come first.
	if HasDynFlag(f, elf.DT_FLAGS_1, DF_1_PIE) {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "Binary is compiled as PIE (enables ASLR when system supports it)",
		}
	}

	for _, prog := range f.Progs {
		if prog.Type == elf.PT_INTERP {
			return rule.ExecuteResult{
				Status:  rule.StatusPassed,
				Message: "Binary is compiled as PIE (enables ASLR when system supports it)",
			}
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusSkipped,
		Message: "Shared library (PIE check not applicable)",
	}
}
