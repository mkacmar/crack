package elf

import (
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const SafeStackRuleID = "safe-stack"

// SafeStackRule checks for SafeStack protection
// Clang: https://clang.llvm.org/docs/SafeStack.html
// LLVM: https://llvm.org/docs/SafeStack.html
type SafeStackRule struct{}

func (r SafeStackRule) ID() string   { return SafeStackRuleID }
func (r SafeStackRule) Name() string { return "SafeStack" }

func (r SafeStackRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 7}, Flag: "-fsanitize=safe-stack"},
		},
	}
}

func (r SafeStackRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	for _, sym := range append(bin.Symbols, bin.DynSymbols...) {
		if strings.HasPrefix(sym.Name, "__safestack_") {
			return rule.ExecuteResult{
				Status:  rule.StatusPassed,
				Message: "SafeStack enabled",
			}
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "SafeStack not enabled",
	}
}
