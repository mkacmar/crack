package elf

import (
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const CFIRuleID = "cfi"

// Cross-DSO CFI runtime symbols https://github.com/llvm/llvm-project/blob/main/compiler-rt/lib/cfi/cfi.cpp
var cfiCrossDSOSymbols = []string{
	"__cfi_check",
	"__cfi_slowpath",
	"__cfi_init",
}

// CFI jump table https://github.com/llvm/llvm-project/blob/main/llvm/lib/Transforms/IPO/LowerTypeTests.cpp
// Original function f() becomes a jump table entry, actual code is renamed to f.cfi
var cfiSuffixes = []string{
	".cfi",
}

// CFI type metadata https://github.com/llvm/llvm-project/blob/main/llvm/lib/Transforms/IPO/LowerTypeTests.cpp
// Encodes function signature types for validating indirect call targets
var cfiPrefixes = []string{
	"__typeid__",
}

// CFIRule checks for Clang Control Flow Integrity
// https://clang.llvm.org/docs/ControlFlowIntegrity.html
type CFIRule struct{}

func (r CFIRule) ID() string   { return CFIRuleID }
func (r CFIRule) Name() string { return "Control Flow Integrity" }

func (r CFIRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 7}, Flag: "-fsanitize=cfi -flto -fvisibility=hidden"},
		},
	}
}

func (r CFIRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	for _, sym := range bin.DynSymbols {
		for _, cfiSym := range cfiCrossDSOSymbols {
			if strings.Contains(sym.Name, cfiSym) {
				return rule.ExecuteResult{
					Status:  rule.StatusPassed,
					Message: "CFI is enabled (cross-DSO mode)",
				}
			}
		}
	}

	if len(bin.Symbols) == 0 {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Binary is stripped, cannot detect regular CFI",
		}
	}

	for _, sym := range bin.Symbols {
		for _, suffix := range cfiSuffixes {
			if strings.HasSuffix(sym.Name, suffix) {
				return rule.ExecuteResult{
					Status:  rule.StatusPassed,
					Message: "CFI is enabled",
				}
			}
		}
		for _, prefix := range cfiPrefixes {
			if strings.HasPrefix(sym.Name, prefix) {
				return rule.ExecuteResult{
					Status:  rule.StatusPassed,
					Message: "CFI is enabled",
				}
			}
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "CFI is NOT enabled",
	}
}
