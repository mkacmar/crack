package elf

import (
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const KernelCFIRuleID = "kernel-cfi"

// KernelCFIRule checks for Kernel CFI (kCFI) protection
// LLVM: https://llvm.org/docs/LangRef.html#kcfi
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-fsanitize-kcfi
type KernelCFIRule struct{}

func (r KernelCFIRule) ID() string   { return KernelCFIRuleID }
func (r KernelCFIRule) Name() string { return "Kernel CFI (kCFI)" }

func (r KernelCFIRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 12, Minor: 0}, Flag: "-fsanitize=kcfi"},
		},
	}
}

func (r KernelCFIRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	for _, sym := range append(bin.Symbols, bin.DynSymbols...) {
		if strings.Contains(sym.Name, "__kcfi") {
			return rule.ExecuteResult{
				Status:  rule.StatusPassed,
				Message: "Kernel CFI (kCFI) is enabled",
			}
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "Kernel CFI is NOT enabled",
	}
}
