package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// ASANRule checks for AddressSanitizer instrumentation
// Clang: https://clang.llvm.org/docs/AddressSanitizer.html
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-fsanitize=address
type ASANRule struct{}

func (r ASANRule) ID() string                     { return "asan" }
func (r ASANRule) Name() string                   { return "Address Sanitizer" }
func (r ASANRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r ASANRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r ASANRule) TargetArch() model.Architecture { return model.ArchAll }
func (r ASANRule) HasPerfImpact() bool            { return true }

func (r ASANRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 4, Minor: 8}, Flag: "-fsanitize=address"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 1}, Flag: "-fsanitize=address"},
		},
	}
}

func (r ASANRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	// AddressSanitizer is a C/C++ feature - not applicable to standard Go/Rust toolchains
	if info != nil {
		switch info.Language {
		case model.LangGo:
			return model.RuleResult{
				State:   model.CheckStateSkipped,
				Message: "Standard Go toolchain (ASAN not applicable)",
			}
		case model.LangRust:
			return model.RuleResult{
				State:   model.CheckStateSkipped,
				Message: "Standard Rust toolchain (has compile-time memory safety)",
			}
		}
	}

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
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "AddressSanitizer is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "AddressSanitizer is NOT enabled",
	}
}
