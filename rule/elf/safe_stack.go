package elf

import (
	"strings"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// SafeStackRuleID is the rule ID for SafeStack.
const SafeStackRuleID = "safe-stack"

// SafeStackRule checks for SafeStack protection.
// Clang: https://clang.llvm.org/docs/SafeStack.html
// LLVM: https://llvm.org/docs/SafeStack.html
type SafeStackRule struct{}

func (r SafeStackRule) ID() string   { return SafeStackRuleID }
func (r SafeStackRule) Name() string { return "SafeStack" }
func (r SafeStackRule) Description() string {
	return "Checks for Clang SafeStack instrumentation. SafeStack separates the stack into a safe stack for return addresses and an unsafe stack for buffers, protecting control flow from stack buffer overflows."
}

func (r SafeStackRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 7}, Flag: "-fsanitize=safe-stack"},
		},
	}
}

func (r SafeStackRule) Execute(bin *binary.ELFBinary) rule.Result {
	for _, sym := range append(bin.Symbols, bin.DynSymbols...) {
		if strings.HasPrefix(sym.Name, "__safestack_") {
			return rule.Result{
				Status:  rule.StatusPassed,
				Message: "SafeStack enabled",
			}
		}
	}

	return rule.Result{
		Status:  rule.StatusFailed,
		Message: "SafeStack not enabled",
	}
}
