package elf

import (
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const ASANRuleID = "asan"

// ASANRule checks for AddressSanitizer instrumentation
// Clang: https://clang.llvm.org/docs/AddressSanitizer.html
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-fsanitize=address
type ASANRule struct{}

func (r ASANRule) ID() string   { return ASANRuleID }
func (r ASANRule) Name() string { return "Address Sanitizer" }

func (r ASANRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 8}, Flag: "-fsanitize=address"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 1}, Flag: "-fsanitize=address"},
		},
	}
}

func (r ASANRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	for _, sym := range append(bin.Symbols, bin.DynSymbols...) {
		if strings.HasPrefix(sym.Name, "__asan_") {
			return rule.ExecuteResult{
				Status:  rule.StatusPassed,
				Message: "AddressSanitizer is enabled",
			}
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "AddressSanitizer is NOT enabled",
	}
}
