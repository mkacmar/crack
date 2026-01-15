package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// NoInsecureRUNPATHRule checks for insecure RUNPATH values
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type NoInsecureRUNPATHRule struct{}

func (r NoInsecureRUNPATHRule) ID() string                     { return "no-insecure-runpath" }
func (r NoInsecureRUNPATHRule) Name() string                   { return "Secure RUNPATH" }
func (r NoInsecureRUNPATHRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r NoInsecureRUNPATHRule) FlagType() model.FlagType       { return model.FlagTypeLink }
func (r NoInsecureRUNPATHRule) TargetArch() model.Architecture { return model.ArchAll }
func (r NoInsecureRUNPATHRule) HasPerfImpact() bool            { return false }

func (r NoInsecureRUNPATHRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,--enable-new-dtags -Wl,-rpath,/absolute/path"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,--enable-new-dtags -Wl,-rpath,/absolute/path"},
		},
	}
}

func (r NoInsecureRUNPATHRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasInsecureRUNPATH := false

	dynSec := f.Section(".dynamic")
	if dynSec != nil {
		data, _ := dynSec.Data()
		if data != nil {
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
				if tag.Tag == int64(elf.DT_RUNPATH) {
					dynstr := f.Section(".dynstr")
					if dynstr != nil {
						strtab, _ := dynstr.Data()
						if strtab != nil && int(tag.Val) < len(strtab) {
							end := int(tag.Val)
							for end < len(strtab) && strtab[end] != 0 {
								end++
							}
							runpath := string(strtab[int(tag.Val):end])
							hasInsecureRUNPATH = isInsecurePath(runpath)
						}
					}
					break
				}
			}
		}
	}

	if !hasInsecureRUNPATH {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "RUNPATH is secure or not present",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Binary has INSECURE RUNPATH (contains relative or world-writable paths)",
	}
}
