package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

// StrippedRule checks if binary is fully stripped
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-_002d_002dstrip_002dall
type StrippedRule struct{}

func (r StrippedRule) ID() string   { return "stripped" }
func (r StrippedRule) Name() string { return "Stripped Binary" }

func (r StrippedRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Arch: binary.ArchAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-s"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-s"},
		},
	}
}

func (r StrippedRule) Execute(f *elf.File, info *binary.Parsed) rule.Result {
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
		if hasSymbolTable && hasDebugSections {
			break
		}
	}

	if !hasSymbolTable && !hasDebugSections {
		return rule.Result{
			State:   rule.CheckStatePassed,
			Message: "Binary is fully stripped (no symbol table or debug sections)",
		}
	}

	if hasSymbolTable && hasDebugSections {
		return rule.Result{
			State:   rule.CheckStateFailed,
			Message: "Binary is NOT stripped (contains symbol table and debug sections)",
		}
	}

	if hasSymbolTable {
		return rule.Result{
			State:   rule.CheckStateFailed,
			Message: "Binary is NOT stripped (contains symbol table)",
		}
	}

	return rule.Result{
		State:   rule.CheckStateFailed,
		Message: "Binary is partially stripped (no symbol table, but has debug sections)",
	}
}
