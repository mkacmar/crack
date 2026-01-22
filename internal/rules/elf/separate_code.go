package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// SeparateCodeRule checks if code and data are in separate pages
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type SeparateCodeRule struct{}

func (r SeparateCodeRule) ID() string                 { return "separate-code" }
func (r SeparateCodeRule) Name() string               { return "Separate Code Segments" }
func (r SeparateCodeRule) Format() model.BinaryFormat { return model.FormatELF }
func (r SeparateCodeRule) FlagType() model.FlagType   { return model.FlagTypeLink }
func (r SeparateCodeRule) HasPerfImpact() bool        { return false }

func (r SeparateCodeRule) Applicability() model.Applicability {
	return model.Applicability{
		Arch: model.ArchAll,
		Compilers: map[model.Compiler]model.CompilerRequirement{
			model.CompilerGCC:   {MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,separate-code"},
			model.CompilerClang: {MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,separate-code"},
		},
	}
}

func (r SeparateCodeRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	// Check file offsets at 4KB page granularity
	const pageSize uint64 = 4096

	var codePages, dataPages [][2]uint64 // [start, end) page ranges

	for _, prog := range f.Progs {
		if prog.Type != elf.PT_LOAD {
			continue
		}

		startPage := prog.Off / pageSize
		endPage := (prog.Off + prog.Filesz + pageSize - 1) / pageSize

		if (prog.Flags & elf.PF_X) != 0 {
			codePages = append(codePages, [2]uint64{startPage, endPage})
		}
		if (prog.Flags & elf.PF_W) != 0 {
			dataPages = append(dataPages, [2]uint64{startPage, endPage})
		}
	}

	if len(codePages) == 0 {
		return model.RuleResult{
			State:   model.CheckStateSkipped,
			Message: "No code segments found",
		}
	}

	// Check if any code page range overlaps with any data page range
	for _, code := range codePages {
		for _, data := range dataPages {
			if code[0] < data[1] && code[1] > data[0] {
				return model.RuleResult{
					State:   model.CheckStateFailed,
					Message: "Code and data segments share page boundary",
				}
			}
		}
	}

	return model.RuleResult{
		State:   model.CheckStatePassed,
		Message: "Code and data are in separate pages",
	}
}
