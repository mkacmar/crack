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

	// ET_DYN can be either a PIE executable or a shared library
	// Check DF_1_PIE flag first, fall back to PT_INTERP for older binaries
	if checkDF1PIE(f) {
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

func checkDF1PIE(f *elf.File) bool {
	dynSec := f.Section(".dynamic")
	if dynSec == nil {
		return false
	}

	data, err := dynSec.Data()
	if err != nil {
		return false
	}

	if f.Class == elf.ELFCLASS64 {
		for i := 0; i < len(data); i += 16 {
			if i+16 > len(data) {
				break
			}
			tag := f.ByteOrder.Uint64(data[i : i+8])
			val := f.ByteOrder.Uint64(data[i+8 : i+16])
			if tag == uint64(elf.DT_NULL) {
				break
			}
			if tag == uint64(elf.DT_FLAGS_1) && (val&DF_1_PIE) != 0 {
				return true
			}
		}
	} else {
		for i := 0; i < len(data); i += 8 {
			if i+8 > len(data) {
				break
			}
			tag := f.ByteOrder.Uint32(data[i : i+4])
			val := f.ByteOrder.Uint32(data[i+4 : i+8])
			if tag == uint32(elf.DT_NULL) {
				break
			}
			if tag == uint32(elf.DT_FLAGS_1) && (val&DF_1_PIE) != 0 {
				return true
			}
		}
	}

	return false
}
