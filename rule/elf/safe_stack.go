package elf

import (
	"strings"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/binary/elf"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
)

// SafeStackRuleID is the rule ID for SafeStack.
const SafeStackRuleID = "safe-stack"

// SafeStackRule checks for SafeStack protection.
//
// References:
//   - https://clang.llvm.org/docs/SafeStack.html
//   - https://llvm.org/docs/SafeStack.html
type SafeStackRule struct{}

func (r SafeStackRule) ID() string   { return SafeStackRuleID }
func (r SafeStackRule) Name() string { return "SafeStack" }
func (r SafeStackRule) Description() string {
	return "Checks for Clang SafeStack instrumentation. SafeStack separates the stack into a safe stack for return addresses and an unsafe stack for buffers, protecting control flow from stack buffer overflows."
}

func (r SafeStackRule) Applicability() rule.Applicability {
	return rule.Applicability{
		// LLVM does not support SafeStack on riscv64.
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 7}, Flag: "-fsanitize=safe-stack"},
		},
		LibC: binary.LibCAll,
	}
}

func (r SafeStackRule) Execute(bin elf.Binary) rule.Result {
	symbols, err := bin.Symbols()
	if err != nil {
		return rule.Skip("symbols unavailable", err)
	}
	dynSymbols, err := bin.DynSymbols()
	if err != nil {
		return rule.Skip("dynamic symbols unavailable", err)
	}
	for _, sym := range append(symbols, dynSymbols...) {
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
