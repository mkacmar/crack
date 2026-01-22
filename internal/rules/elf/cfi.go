package elf

import (
	"debug/elf"
	"fmt"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

var cfiSymbols = []string{
	"__cfi_check",
	"__cfi_check_fail",
	"__cfi_slowpath",
	"__cfi_slowpath_diag",
	"__ubsan_handle_cfi_check_fail",
	"__ubsan_handle_cfi_check_fail_abort",
}

// CFIRule checks for Clang Control Flow Integrity
// Clang: https://clang.llvm.org/docs/ControlFlowIntegrity.html
type CFIRule struct{}

func (r CFIRule) ID() string                 { return "cfi" }
func (r CFIRule) Name() string               { return "Control Flow Integrity" }
func (r CFIRule) Format() model.BinaryFormat { return model.FormatELF }
func (r CFIRule) FlagType() model.FlagType   { return model.FlagTypeCompile }
func (r CFIRule) HasPerfImpact() bool        { return true }

func (r CFIRule) Applicability() model.Applicability {
	return model.Applicability{
		Arch: model.ArchAll,
		Compilers: map[model.Compiler]model.CompilerRequirement{
			model.CompilerClang: {MinVersion: model.Version{Major: 3, Minor: 7}, Flag: "-fsanitize=cfi -flto -fvisibility=hidden"},
		},
	}
}

func (r CFIRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {

	symbols, _ := f.Symbols()
	dynsyms, _ := f.DynamicSymbols()

	allSymbols := make(map[string]struct{})
	for _, sym := range symbols {
		allSymbols[sym.Name] = struct{}{}
	}
	for _, sym := range dynsyms {
		allSymbols[sym.Name] = struct{}{}
	}

	var foundSymbols []string
	for _, cfiSym := range cfiSymbols {
		for symName := range allSymbols {
			if strings.Contains(symName, cfiSym) {
				foundSymbols = append(foundSymbols, cfiSym)
				break
			}
		}
	}

	if len(foundSymbols) > 0 {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: fmt.Sprintf("Clang CFI is enabled, found: %v", foundSymbols),
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Clang CFI is NOT enabled",
	}
}
