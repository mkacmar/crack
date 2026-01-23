package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

// StackCanaryRule checks for stack canary protection
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-fstack-protector
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-fstack-protector-strong
type StackCanaryRule struct{}

func (r StackCanaryRule) ID() string   { return "stack-canary" }
func (r StackCanaryRule) Name() string { return "Stack Canary Protection" }

func (r StackCanaryRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 9}, Flag: "-fstack-protector-strong"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-fstack-protector-strong"},
		},
	}
}

func (r StackCanaryRule) Execute(f *elf.File, info *binary.Parsed) rule.ExecuteResult {
	symbols, _ := f.Symbols()
	dynsyms, _ := f.DynamicSymbols()

	// Check both static and dynamic symbol tables for stack protection symbols.
	// Different compilers/platforms use different symbols:
	// - __stack_chk_fail: GCC/Clang on Linux
	// - __stack_smash_handler: older GCC
	// - __intel_security_cookie: Intel compiler
	for _, sym := range append(symbols, dynsyms...) {
		if strings.Contains(sym.Name, "__stack_chk_fail") ||
			strings.Contains(sym.Name, "__stack_smash_handler") ||
			strings.Contains(sym.Name, "__intel_security_cookie") {
			return rule.ExecuteResult{
				Status:  rule.StatusPassed,
				Message: "Stack canary protection is enabled",
			}
		}
	}

	// No stack protection symbols found.
	// While the compiler might omit these if no functions need protection, real-world binaries typically have stack buffers.
	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "Stack canary protection is NOT enabled",
	}
}
