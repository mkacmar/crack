package elf

import (
	"strings"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
)

// ASANRuleID is the rule ID for ASan.
const ASANRuleID = "asan"

// ASANRule checks for AddressSanitizer instrumentation.
// Clang: https://clang.llvm.org/docs/AddressSanitizer.html
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-fsanitize=address
type ASANRule struct{}

func (r ASANRule) ID() string   { return ASANRuleID }
func (r ASANRule) Name() string { return "Address Sanitizer" }
func (r ASANRule) Description() string {
	return "Checks for AddressSanitizer (ASan) instrumentation. ASan detects memory errors including buffer overflows, use-after-free, and memory leaks at runtime."
}

func (r ASANRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 5, Minor: 1}, Flag: "-fsanitize=address"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, Flag: "-fsanitize=address"},
		},
	}
}

func (r ASANRule) Execute(bin *binary.ELFBinary) rule.Result {
	for _, sym := range append(bin.Symbols, bin.DynSymbols...) {
		if strings.HasPrefix(sym.Name, "__asan_") {
			return rule.Result{
				Status:  rule.StatusPassed,
				Message: "ASan enabled",
			}
		}
	}

	return rule.Result{
		Status:  rule.StatusFailed,
		Message: "ASan not enabled",
	}
}
