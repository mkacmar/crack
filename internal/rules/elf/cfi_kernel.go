package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// KernelCFIRule checks for Kernel CFI (kCFI) protection
// LLVM: https://llvm.org/docs/LangRef.html#kcfi
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-fsanitize-kcfi
type KernelCFIRule struct{}

func (r KernelCFIRule) ID() string                 { return "kernel-cfi" }
func (r KernelCFIRule) Name() string               { return "Kernel CFI (kCFI)" }
func (r KernelCFIRule) Format() model.BinaryFormat { return model.FormatELF }

func (r KernelCFIRule) Applicability() model.Applicability {
	return model.Applicability{
		Arch: model.ArchAll,
		Compilers: map[model.Compiler]model.CompilerRequirement{
			model.CompilerClang: {MinVersion: model.Version{Major: 12, Minor: 0}, Flag: "-fsanitize=kcfi"},
		},
	}
}

func (r KernelCFIRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	symbols, _ := f.Symbols()
	dynsyms, _ := f.DynamicSymbols()

	hasKCFI := false
	for _, sym := range symbols {
		if strings.Contains(sym.Name, "__kcfi") {
			hasKCFI = true
			break
		}
	}
	if !hasKCFI {
		for _, sym := range dynsyms {
			if strings.Contains(sym.Name, "__kcfi") {
				hasKCFI = true
				break
			}
		}
	}

	if hasKCFI {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Kernel CFI (kCFI) is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Kernel CFI is NOT enabled",
	}
}
