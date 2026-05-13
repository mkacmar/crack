package elf

import (
	stdelf "debug/elf"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/binary/elf"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
)

// PIERuleID is the rule ID for PIE.
const PIERuleID = "pie"

// PIERule checks if binary is compiled as Position Independent Executable.
//
// References:
//   - https://gcc.gnu.org/onlinedocs/gcc/Code-Gen-Options.html#index-fPIE
//   - https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-fpie
type PIERule struct{}

func (r PIERule) ID() string   { return PIERuleID }
func (r PIERule) Name() string { return "Position Independent Executable" }
func (r PIERule) Description() string {
	return "Checks if the binary is compiled as a Position Independent Executable (PIE). PIE enables full ASLR by allowing the executable to be loaded at a random base address, making return-oriented programming (ROP) attacks significantly harder."
}

func (r PIERule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM | binary.ArchRISCV},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, DefaultVersion: toolchain.Version{Major: 6, Minor: 1}, Flag: "-fPIE -pie"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, DefaultVersion: toolchain.Version{Major: 4, Minor: 0}, Flag: "-fPIE -pie"},
		},
		LibC: binary.LibCAll,
	}
}

func (r PIERule) Execute(bin elf.Binary) rule.Result {
	switch bin.Type() {
	case stdelf.ET_EXEC:
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: "Not PIE",
		}
	case stdelf.ET_DYN:
	default:
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	pie, err := elf.HasDynFlag(bin, stdelf.DT_FLAGS_1, uint64(stdelf.DF_1_PIE))
	if err != nil {
		return rule.Skip("failed to read dynamic section", err)
	}
	if pie {
		return rule.Result{
			Status:  rule.StatusPassed,
			Message: "PIE enabled",
		}
	}

	for _, prog := range bin.Progs() {
		if prog.Type == stdelf.PT_INTERP {
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
