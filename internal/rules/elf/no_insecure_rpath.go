package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// NoInsecureRPATHRule checks for insecure RPATH values
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type NoInsecureRPATHRule struct{}

func (r NoInsecureRPATHRule) ID() string                     { return "no-insecure-rpath" }
func (r NoInsecureRPATHRule) Name() string                   { return "Secure RPATH" }
func (r NoInsecureRPATHRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r NoInsecureRPATHRule) FlagType() model.FlagType       { return model.FlagTypeLink }
func (r NoInsecureRPATHRule) TargetArch() model.Architecture { return model.ArchAll }
func (r NoInsecureRPATHRule) HasPerfImpact() bool            { return false }

func (r NoInsecureRPATHRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-rpath,/absolute/path"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-rpath,/absolute/path"},
		},
	}
}

func (r NoInsecureRPATHRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasInsecureRPATH := false

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
				if tag.Tag == int64(elf.DT_RPATH) {
					dynstr := f.Section(".dynstr")
					if dynstr != nil {
						strtab, _ := dynstr.Data()
						if strtab != nil && int(tag.Val) < len(strtab) {
							end := int(tag.Val)
							for end < len(strtab) && strtab[end] != 0 {
								end++
							}
							rpath := string(strtab[int(tag.Val):end])
							hasInsecureRPATH = isInsecurePath(rpath)
						}
					}
					break
				}
			}
		}
	}

	if !hasInsecureRPATH {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "RPATH is secure or not present",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Binary has INSECURE RPATH (contains relative or world-writable paths)",
	}
}

func isInsecurePath(path string) bool {
	if path == "" {
		return false
	}

	insecurePaths := []string{"/tmp", ".", ".."}
	paths := strings.Split(path, ":")
	for _, p := range paths {
		p = strings.TrimSpace(p)
		for _, insecure := range insecurePaths {
			if p == insecure || strings.HasPrefix(p, insecure+"/") {
				return true
			}
		}
		if strings.HasPrefix(p, "./") || strings.HasPrefix(p, "../") {
			return true
		}
	}
	return false
}
