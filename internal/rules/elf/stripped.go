package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// StrippedRule checks if binary is fully stripped
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-_002d_002dstrip_002dall
type StrippedRule struct{}

func (r StrippedRule) ID() string                     { return "stripped" }
func (r StrippedRule) Name() string                   { return "Stripped Binary" }
func (r StrippedRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r StrippedRule) FlagType() model.FlagType       { return model.FlagTypeLink }
func (r StrippedRule) TargetArch() model.Architecture { return model.ArchAll }
func (r StrippedRule) HasPerfImpact() bool            { return false }

func (r StrippedRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-s"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-s"},
		},
	}
}

func (r StrippedRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasSymbolTable := false
	hasDebugSections := false

	for _, section := range f.Sections {
		if section.Type == elf.SHT_SYMTAB {
			hasSymbolTable = true
		}
		if strings.HasPrefix(section.Name, ".debug_") ||
			strings.HasPrefix(section.Name, ".zdebug_") {
			hasDebugSections = true
		}
	}

	isFullyStripped := !hasSymbolTable && !hasDebugSections

	if isFullyStripped {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Binary is fully stripped (no symbol table or debug sections)",
		}
	}

	message := "Binary is NOT stripped"
	if hasSymbolTable && hasDebugSections {
		message = "Binary is NOT stripped (contains symbol table and debug sections)"
	} else if hasSymbolTable {
		message = "Binary is NOT stripped (contains symbol table)"
	} else if hasDebugSections {
		message = "Binary is partially stripped (no symbol table, but has debug sections)"
	}

	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: message,
	}
}
