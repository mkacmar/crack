package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// SeparateCodeRule checks if code and data are in separate pages
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type SeparateCodeRule struct{}

func (r SeparateCodeRule) ID() string                     { return "separate-code" }
func (r SeparateCodeRule) Name() string                   { return "Separate Code Segments" }
func (r SeparateCodeRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r SeparateCodeRule) FlagType() model.FlagType       { return model.FlagTypeLink }
func (r SeparateCodeRule) TargetArch() model.Architecture { return model.ArchAll }
func (r SeparateCodeRule) HasPerfImpact() bool            { return false }

func (r SeparateCodeRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,separate-code"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,separate-code"},
		},
	}
}

func (r SeparateCodeRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	type segment struct {
		startPage, endPage uint64
		isCode             bool
		isData             bool
	}

	// Determine page size from PT_LOAD alignment (typically matches system page size)
	var pageSize uint64 = 4096
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_LOAD && prog.Align > pageSize {
			pageSize = prog.Align
		}
	}

	var segments []segment

	// Collect all PT_LOAD segments and convert addresses to page numbers
	for _, prog := range f.Progs {
		if prog.Type != elf.PT_LOAD {
			continue
		}

		isExecutable := (prog.Flags & elf.PF_X) != 0
		isWritable := (prog.Flags & elf.PF_W) != 0

		seg := segment{
			startPage: prog.Vaddr / pageSize,
			endPage:   (prog.Vaddr + prog.Memsz + pageSize - 1) / pageSize,
			isCode:    isExecutable,
			isData:    isWritable,
		}

		segments = append(segments, seg)
	}

	hasCode := false
	for _, seg := range segments {
		if seg.isCode {
			hasCode = true
			break
		}
	}

	if !hasCode {
		return model.RuleResult{
			State:   model.CheckStateSkipped,
			Message: "No code segments found",
		}
	}

	// Check if any code segment shares a page with any data segment
	// Without -z separate-code, code and data may end up on the same page,
	// requiring the page to be both writable and executable at runtime
	for _, code := range segments {
		if !code.isCode {
			continue
		}
		for _, data := range segments {
			if !data.isData {
				continue
			}
			if code.endPage >= data.startPage && code.startPage < data.endPage {
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
