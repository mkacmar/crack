package elf

import (
	"strings"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// StackCanaryRuleID is the rule ID for stack canary.
const StackCanaryRuleID = "stack-canary"

// StackCanaryRule checks for stack canary protection.
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-fstack-protector
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-fstack-protector-strong
type StackCanaryRule struct{}

func (r StackCanaryRule) ID() string   { return StackCanaryRuleID }
func (r StackCanaryRule) Name() string { return "Stack Canary Protection" }
func (r StackCanaryRule) Description() string {
	return "Checks for stack canary (stack protector) instrumentation. Stack canaries detect buffer overflows by placing a guard value before the return address. If the canary is corrupted, the program terminates before exploitation can occur."
}

func (r StackCanaryRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 9}, DefaultVersion: toolchain.Version{Major: 4, Minor: 9}, Flag: "-fstack-protector-strong"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 5}, DefaultVersion: toolchain.Version{Major: 3, Minor: 5}, Flag: "-fstack-protector-strong"},
		},
	}
}

func (r StackCanaryRule) Execute(bin *binary.ELFBinary) rule.Result {
	for _, sym := range append(bin.Symbols, bin.DynSymbols...) {
		if strings.Contains(sym.Name, "__stack_chk_fail") ||
			strings.Contains(sym.Name, "__stack_smash_handler") ||
			strings.Contains(sym.Name, "__intel_security_cookie") {
			return rule.Result{
				Status:  rule.StatusPassed,
				Message: "Stack canary enabled",
			}
		}
	}

	return rule.Result{
		Status:  rule.StatusFailed,
		Message: "Stack canary not enabled",
	}
}
