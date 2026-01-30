package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const PIERuleID = "pie"

// PIERule checks if binary is compiled as Position Independent Executable
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Code-Gen-Options.html#index-fPIE
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-fpie
type PIERule struct{}

func (r PIERule) ID() string   { return PIERuleID }
func (r PIERule) Name() string { return "Position Independent Executable" }

func (r PIERule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, DefaultVersion: toolchain.Version{Major: 6, Minor: 1}, Flag: "-fPIE -pie"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, DefaultVersion: toolchain.Version{Major: 4, Minor: 0}, Flag: "-fPIE -pie"},
		},
	}
}

func (r PIERule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	switch bin.File.Type {
	case elf.ET_EXEC:
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: "Not PIE",
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
	if HasDynFlag(bin.File, elf.DT_FLAGS_1, uint64(elf.DF_1_PIE)) {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "PIE enabled",
		}
	}

	for _, prog := range bin.File.Progs {
		if prog.Type == elf.PT_INTERP {
			return rule.ExecuteResult{
				Status:  rule.StatusPassed,
				Message: "PIE enabled",
			}
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusSkipped,
		Message: "Shared library, PIE not applicable",
	}
}
