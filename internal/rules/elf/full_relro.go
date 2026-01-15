package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

const (
	DF_BIND_NOW = 0x8
	DF_1_NOW    = 0x1
)

// FullRELRORule checks for full RELRO protection
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type FullRELRORule struct{}

func (r FullRELRORule) ID() string                     { return "full-relro" }
func (r FullRELRORule) Name() string                   { return "Full RELRO" }
func (r FullRELRORule) Format() model.BinaryFormat     { return model.FormatELF }
func (r FullRELRORule) FlagType() model.FlagType       { return model.FlagTypeLink }
func (r FullRELRORule) TargetArch() model.Architecture { return model.ArchAll }
func (r FullRELRORule) HasPerfImpact() bool            { return false }

func (r FullRELRORule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,relro,-z,now"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,relro,-z,now"},
		},
	}
}

func (r FullRELRORule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasRELRO := false
	isFullRELRO := false

	for _, prog := range f.Progs {
		if prog.Type == elf.PT_GNU_RELRO {
			hasRELRO = true
			break
		}
	}

	if hasRELRO {
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
					if tag.Tag == int64(elf.DT_BIND_NOW) {
						isFullRELRO = true
						break
					}
					if tag.Tag == int64(elf.DT_FLAGS) && (tag.Val&DF_BIND_NOW) != 0 {
						isFullRELRO = true
						break
					}
					if tag.Tag == int64(elf.DT_FLAGS_1) && (tag.Val&DF_1_NOW) != 0 {
						isFullRELRO = true
						break
					}
				}
			}
		}
	}

	if isFullRELRO {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Full RELRO is enabled (GOT read-only, lazy binding disabled)",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Full RELRO is NOT enabled (GOT may be writable)",
	}
}
