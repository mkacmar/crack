package elf

import (
	"debug/elf"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
)

// PIERuleID is the rule ID for PIE.
const PIERuleID = "pie"

// PIERule checks if binary is compiled as Position Independent Executable.
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Code-Gen-Options.html#index-fPIE
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-fpie
type PIERule struct{}

func (r PIERule) ID() string   { return PIERuleID }
func (r PIERule) Name() string { return "Position Independent Executable" }
func (r PIERule) Description() string {
	return "Checks if the binary is compiled as a Position Independent Executable (PIE). PIE enables full ASLR by allowing the executable to be loaded at a random base address, making return-oriented programming (ROP) attacks significantly harder."
}

func (r PIERule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, DefaultVersion: toolchain.Version{Major: 6, Minor: 1}, Flag: "-fPIE -pie"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, DefaultVersion: toolchain.Version{Major: 4, Minor: 0}, Flag: "-fPIE -pie"},
			toolchain.Rustc: {MinVersion: toolchain.Version{Major: 1, Minor: 26}, DefaultVersion: toolchain.Version{Major: 1, Minor: 26}, Flag: "-C relocation-model=pie"},
		},
	}
}

func (r PIERule) Execute(bin *binary.ELFBinary) rule.Result {
	switch bin.Type {
	case elf.ET_EXEC:
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: "Not PIE",
		}
	case elf.ET_DYN:
	default:
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	if bin.HasDynFlag(elf.DT_FLAGS_1, uint64(elf.DF_1_PIE)) {
		return rule.Result{
			Status:  rule.StatusPassed,
			Message: "PIE enabled",
		}
	}

	for _, prog := range bin.Progs {
		if prog.Type == elf.PT_INTERP {
			return rule.Result{
				Status:  rule.StatusPassed,
				Message: "PIE enabled",
			}
		}
	}

	return rule.Result{
		Status:  rule.StatusSkipped,
		Message: "Shared library, PIE not applicable",
	}
}
