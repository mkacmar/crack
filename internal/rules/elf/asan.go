package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

// ASANRule checks for AddressSanitizer instrumentation
// Clang: https://clang.llvm.org/docs/AddressSanitizer.html
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-fsanitize=address
type ASANRule struct{}

func (r ASANRule) ID() string   { return "asan" }
func (r ASANRule) Name() string { return "Address Sanitizer" }

func (r ASANRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Arch: binary.ArchAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 8}, Flag: "-fsanitize=address"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 1}, Flag: "-fsanitize=address"},
		},
	}
}

func (r ASANRule) Execute(f *elf.File, info *binary.Parsed) rule.ExecuteResult {

	hasASan := false

	symbols, err := f.Symbols()
	if err != nil {
		symbols = nil
	}

	dynsyms, err := f.DynamicSymbols()
	if err != nil {
		dynsyms = nil
	}

	for _, sym := range symbols {
		if strings.HasPrefix(sym.Name, "__asan_") {
			hasASan = true
			break
		}
	}

	if !hasASan {
		for _, sym := range dynsyms {
			if strings.HasPrefix(sym.Name, "__asan_") {
				hasASan = true
				break
			}
		}
	}

	if hasASan {
		return rule.ExecuteResult{
			Status: rule.StatusPassed,
			Message: "AddressSanitizer is enabled",
		}
	}
	return rule.ExecuteResult{
		Status: rule.StatusFailed,
		Message: "AddressSanitizer is NOT enabled",
	}
}
