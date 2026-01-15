package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

const DF_1_NODUMP = 0x00001000

// NoDumpRule checks if core dumps are disabled
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type NoDumpRule struct{}

func (r NoDumpRule) ID() string                     { return "no-dump" }
func (r NoDumpRule) Name() string                   { return "Core Dump Protection" }
func (r NoDumpRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r NoDumpRule) FlagType() model.FlagType       { return model.FlagTypeLink }
func (r NoDumpRule) TargetArch() model.Architecture { return model.ArchAll }
func (r NoDumpRule) HasPerfImpact() bool            { return false }

func (r NoDumpRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,nodump"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,nodump"},
		},
	}
}

func (r NoDumpRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasNoDump := false

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
					if (tag.Val & DF_1_NODUMP) != 0 {
						hasNoDump = true
						break
					}
				}
			}
		}
	}

	if hasNoDump {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Core dumps are disabled (DF_1_NODUMP set)",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Core dumps are NOT explicitly disabled",
	}
}
