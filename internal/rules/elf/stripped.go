package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const StrippedRuleID = "stripped"

// StrippedRule checks if binary is fully stripped
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-_002d_002dstrip_002dall
type StrippedRule struct{}

func (r StrippedRule) ID() string   { return StrippedRuleID }
func (r StrippedRule) Name() string { return "Stripped Binary" }

func (r StrippedRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-s"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 1, Minor: 0}, Flag: "-s"},
		},
	}
}

func (r StrippedRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	hasSymbolTable := false
	hasDebugSections := false

	for _, section := range bin.File.Sections {
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
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "Fully stripped",
		}
	}

	if hasSymbolTable && hasDebugSections {
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: "Not stripped, has symbols and debug info",
		}
	}

	if hasSymbolTable {
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: "Not stripped, has symbols",
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "Partially stripped, has debug info",
	}
}
