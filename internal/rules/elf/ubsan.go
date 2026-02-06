package elf

import (
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const UBSanRuleID = "ubsan"

// UBSanRule checks for Undefined Behavior Sanitizer
// Clang: https://clang.llvm.org/docs/UndefinedBehaviorSanitizer.html
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-fsanitize=undefined
type UBSanRule struct{}

func (r UBSanRule) ID() string   { return UBSanRuleID }
func (r UBSanRule) Name() string { return "Undefined Behavior Sanitizer" }
func (r UBSanRule) Description() string {
	return "Checks for Undefined Behavior Sanitizer (UBSan) instrumentation. UBSan detects undefined behavior such as integer overflows, null pointer dereferences, and misaligned accesses at runtime."
}

func (r UBSanRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			// rustc: different mechanism, compile-time checks rather than runtime sanitizer
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 5, Minor: 1}, Flag: "-fsanitize=undefined"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, Flag: "-fsanitize=undefined"},
		},
	}
}

func (r UBSanRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	for _, sym := range append(bin.Symbols, bin.DynSymbols...) {
		if strings.HasPrefix(sym.Name, "__ubsan_") {
			return rule.ExecuteResult{
				Status:  rule.StatusPassed,
				Message: "UBSan enabled",
			}
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "UBSan not enabled",
	}
}
