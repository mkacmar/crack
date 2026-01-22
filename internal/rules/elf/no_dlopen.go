package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

const (
	DT_FLAGS_1    = 0x6ffffffb
	DF_1_NODLOPEN = 0x00000008
)

// NoDLOpenRule checks if dlopen is disabled
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type NoDLOpenRule struct{}

func (r NoDLOpenRule) ID() string                 { return "no-dlopen" }
func (r NoDLOpenRule) Name() string               { return "Disallow dlopen" }
func (r NoDLOpenRule) Format() model.BinaryFormat { return model.FormatELF }

func (r NoDLOpenRule) Applicability() model.Applicability {
	return model.Applicability{
		Arch: model.ArchAll,
		Compilers: map[model.Compiler]model.CompilerRequirement{
			model.CompilerGCC:   {MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,nodlopen"},
			model.CompilerClang: {MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,nodlopen"},
		},
	}
}

func (r NoDLOpenRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasNoDLOpen := false

	dynSec := f.Section(".dynamic")
	if dynSec != nil {
		data, err := dynSec.Data()
		if err == nil {
			var dynData []elf.Dyn64
			if f.Class == elf.ELFCLASS64 {
				for i := 0; i < len(data); i += 16 {
					if i+16 > len(data) {
						break
					}
					tag := f.ByteOrder.Uint64(data[i : i+8])
					val := f.ByteOrder.Uint64(data[i+8 : i+16])
					dynData = append(dynData, elf.Dyn64{Tag: int64(tag), Val: val})
					if int64(tag) == int64(elf.DT_NULL) {
						break
					}
				}
			} else {
				for i := 0; i < len(data); i += 8 {
					if i+8 > len(data) {
						break
					}
					tag := f.ByteOrder.Uint32(data[i : i+4])
					val := f.ByteOrder.Uint32(data[i+4 : i+8])
					dynData = append(dynData, elf.Dyn64{Tag: int64(tag), Val: uint64(val)})
					if int64(tag) == int64(elf.DT_NULL) {
						break
					}
				}
			}

			for _, tag := range dynData {
				if tag.Tag == DT_FLAGS_1 {
					if (tag.Val & DF_1_NODLOPEN) != 0 {
						hasNoDLOpen = true
						break
					}
				}
			}
		}
	}

	if hasNoDLOpen {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "dlopen is disabled via linker flags",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "dlopen is NOT disabled (binary can load libraries dynamically)",
	}
}
