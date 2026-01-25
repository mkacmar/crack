package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const NoDLOpenRuleID = "no-dlopen"

const (
	DT_FLAGS_1    = 0x6ffffffb
	DF_1_NODLOPEN = 0x00000008
)

// NoDLOpenRule checks if dlopen is disabled
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type NoDLOpenRule struct{}

func (r NoDLOpenRule) ID() string   { return NoDLOpenRuleID }
func (r NoDLOpenRule) Name() string { return "Disallow dlopen" }

func (r NoDLOpenRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,nodlopen"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,nodlopen"},
		},
	}
}

func (r NoDLOpenRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	hasNoDLOpen := false

	dynSec := bin.File.Section(".dynamic")
	if dynSec != nil {
		data, err := dynSec.Data()
		if err == nil {
			var dynData []elf.Dyn64
			if bin.File.Class == elf.ELFCLASS64 {
				for i := 0; i < len(data); i += 16 {
					if i+16 > len(data) {
						break
					}
					tag := bin.File.ByteOrder.Uint64(data[i : i+8])
					val := bin.File.ByteOrder.Uint64(data[i+8 : i+16])
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
					tag := bin.File.ByteOrder.Uint32(data[i : i+4])
					val := bin.File.ByteOrder.Uint32(data[i+4 : i+8])
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
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "dlopen is disabled via linker flags",
		}
	}
	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "dlopen is NOT disabled (binary can load libraries dynamically)",
	}
}
