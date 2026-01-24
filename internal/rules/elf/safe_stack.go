package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

// SafeStackRule checks for SafeStack protection
// Clang: https://clang.llvm.org/docs/SafeStack.html
// LLVM: https://llvm.org/docs/SafeStack.html
type SafeStackRule struct{}

func (r SafeStackRule) ID() string   { return "safe-stack" }
func (r SafeStackRule) Name() string { return "SafeStack" }

func (r SafeStackRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 8}, Flag: "-fsanitize=safe-stack"},
		},
	}
}

func (r SafeStackRule) Execute(f *elf.File, info *binary.Parsed) rule.ExecuteResult {
	symbols, _ := f.Symbols()
	dynsyms, _ := f.DynamicSymbols()

	for _, sym := range append(symbols, dynsyms...) {
		if strings.HasPrefix(sym.Name, "__safestack_") {
			return rule.ExecuteResult{
				Status:  rule.StatusPassed,
				Message: "SafeStack is enabled",
			}
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "SafeStack is NOT enabled",
	}
}
