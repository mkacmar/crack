package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

// KernelCFIRule checks for Kernel CFI (kCFI) protection
// LLVM: https://llvm.org/docs/LangRef.html#kcfi
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-fsanitize-kcfi
type KernelCFIRule struct{}

func (r KernelCFIRule) ID() string   { return "kernel-cfi" }
func (r KernelCFIRule) Name() string { return "Kernel CFI (kCFI)" }

func (r KernelCFIRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Arch: binary.ArchAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 12, Minor: 0}, Flag: "-fsanitize=kcfi"},
		},
	}
}

func (r KernelCFIRule) Execute(f *elf.File, info *binary.Parsed) rule.Result {
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
		return rule.Result{
			State:   rule.CheckStatePassed,
			Message: "Kernel CFI (kCFI) is enabled",
		}
	}
	return rule.Result{
		State:   rule.CheckStateFailed,
		Message: "Kernel CFI is NOT enabled",
	}
}
