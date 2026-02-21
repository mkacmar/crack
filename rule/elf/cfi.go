package elf

import (
	"strings"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
)

// CFIRuleID is the rule ID for CFI.
const CFIRuleID = "cfi"

// Cross-DSO CFI runtime symbols.
var cfiCrossDSOSymbols = []string{
	"__cfi_check",
	"__cfi_slowpath",
	"__cfi_init",
}

// CFI jump table suffixes.
var cfiSuffixes = []string{
	".cfi",
}

// CFI type metadata prefixes.
var cfiPrefixes = []string{
	"__typeid__",
}

// CFIRule checks for Clang Control Flow Integrity.
// https://clang.llvm.org/docs/ControlFlowIntegrity.html
type CFIRule struct{}

func (r CFIRule) ID() string   { return CFIRuleID }
func (r CFIRule) Name() string { return "Control Flow Integrity" }
func (r CFIRule) Description() string {
	return "Checks for Clang Control Flow Integrity (CFI) instrumentation. CFI validates that indirect calls and jumps target expected locations, preventing attackers from hijacking control flow through corrupted function pointers or vtables."
}

func (r CFIRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 6, Minor: 0}, Flag: "-fsanitize=cfi -flto -fvisibility=hidden"},
		},
	}
}

func (r CFIRule) Execute(bin *binary.ELFBinary) rule.Result {
	for _, sym := range bin.DynSymbols {
		for _, cfiSym := range cfiCrossDSOSymbols {
			if strings.Contains(sym.Name, cfiSym) {
				return rule.Result{
					Status:  rule.StatusPassed,
					Message: "CFI enabled (cross-DSO mode)",
				}
			}
		}
	}

	if len(bin.Symbols) == 0 {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Stripped binary, CFI detection limited",
		}
	}

	for _, sym := range bin.Symbols {
		for _, suffix := range cfiSuffixes {
			if strings.HasSuffix(sym.Name, suffix) {
				return rule.Result{
					Status:  rule.StatusPassed,
					Message: "CFI enabled",
				}
			}
		}
		for _, prefix := range cfiPrefixes {
			if strings.HasPrefix(sym.Name, prefix) {
				return rule.Result{
					Status:  rule.StatusPassed,
					Message: "CFI enabled",
				}
			}
		}
	}

	return rule.Result{
		Status:  rule.StatusFailed,
		Message: "CFI not enabled",
	}
}
