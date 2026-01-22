package elf

import (
	"debug/elf"
	"fmt"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
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

func (r CFIRule) ID() string   { return "cfi" }
func (r CFIRule) Name() string { return "Control Flow Integrity" }

func (r CFIRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Arch: binary.ArchAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 7}, Flag: "-fsanitize=cfi -flto -fvisibility=hidden"},
		},
	}
}

func (r CFIRule) Execute(f *elf.File, info *binary.Parsed) rule.ExecuteResult {

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
		return rule.ExecuteResult{
			Status: rule.StatusPassed,
			Message: fmt.Sprintf("Clang CFI is enabled, found: %v", foundSymbols),
		}
	}
	return rule.ExecuteResult{
		Status: rule.StatusFailed,
		Message: "Clang CFI is NOT enabled",
	}
}
