package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// StrippedRuleID is the rule ID for stripped binary.
const StrippedRuleID = "stripped"

// StrippedRule checks if binary is fully stripped.
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-_002d_002dstrip_002dall
type StrippedRule struct{}

func (r StrippedRule) ID() string   { return StrippedRuleID }
func (r StrippedRule) Name() string { return "Stripped Binary" }
func (r StrippedRule) Description() string {
	return "Checks if the binary has been stripped of symbol tables and debug information. Stripping removes metadata that could help attackers understand the binary's structure and identify vulnerabilities."
}

func (r StrippedRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-s"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 1, Minor: 0}, Flag: "-s"},
			toolchain.Rustc: {MinVersion: toolchain.Version{Major: 1, Minor: 59}, Flag: "-C strip=symbols"},
		},
	}
}

func (r StrippedRule) Execute(bin *binary.ELFBinary) rule.Result {
	hasSymbolTable := false
	hasDebugSections := false

	for _, section := range bin.Sections {
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
			Status:  rule.StatusPassed,
			Message: "Fully stripped",
		}
	}

	if hasSymbolTable && hasDebugSections {
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: "Not stripped, has symbols and debug info",
		}
	}

	if hasSymbolTable {
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: "Not stripped, has symbols",
		}
	}

	return rule.Result{
		Status:  rule.StatusFailed,
		Message: "Partially stripped, has debug info",
	}
}
