package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// X86RetpolineRule checks for Spectre v2 mitigation (retpoline)
// GCC: https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html#index-mindirect-branch
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-mretpoline
type X86RetpolineRule struct{}

func (r X86RetpolineRule) ID() string                     { return "x86-retpoline" }
func (r X86RetpolineRule) Name() string                   { return "x86 Retpoline (Spectre v2)" }
func (r X86RetpolineRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r X86RetpolineRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r X86RetpolineRule) TargetArch() model.Architecture { return model.ArchAllX86 }
func (r X86RetpolineRule) HasPerfImpact() bool            { return true }

func (r X86RetpolineRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 7, Minor: 3}, Flag: "-mindirect-branch=thunk -mfunction-return=thunk"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 5, Minor: 0}, Flag: "-mretpoline"},
		},
	}
}

func (r X86RetpolineRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	// CET-IBT and retpoline both mitigate Spectre v2 indirect branch attacks.
	// Hardware mitigations take precedence over software (retpoline).
	// See: https://www.kernel.org/doc/html/latest/admin-guide/hw-vuln/spectre.html
	hasCETIBT := parseGNUProperty(f, GNU_PROPERTY_X86_FEATURE_1_AND, GNU_PROPERTY_X86_FEATURE_1_IBT)
	if hasCETIBT {
		return model.RuleResult{
			State:   model.CheckStateSkipped,
			Message: "CET-IBT enabled (hardware mitigation supersedes retpoline)",
		}
	}

	// Retpoline thunks are typically local symbols in .symtab.
	// Stripped binaries lose .symtab, making detection impossible.
	hasSymtab := false
	for _, sec := range f.Sections {
		if sec.Type == elf.SHT_SYMTAB {
			hasSymtab = true
			break
		}
	}
	if !hasSymtab {
		return model.RuleResult{
			State:   model.CheckStateSkipped,
			Message: "Stripped binary (retpoline symbols not available)",
		}
	}

	symbols, _ := f.Symbols()
	dynsyms, _ := f.DynamicSymbols()

	var hasGCCThunk, hasLLVMRetpoline bool
	for _, sym := range append(symbols, dynsyms...) {
		switch {
		case strings.Contains(sym.Name, "__x86_indirect_thunk"),
			strings.Contains(sym.Name, "__x86_return_thunk"):
			hasGCCThunk = true
		case strings.Contains(sym.Name, "__llvm_retpoline"):
			hasLLVMRetpoline = true
		}
		if hasGCCThunk || hasLLVMRetpoline {
			break
		}
	}

	if hasGCCThunk || hasLLVMRetpoline {
		msg := "Retpoline enabled (GCC)"
		if hasLLVMRetpoline {
			msg = "Retpoline enabled (LLVM)"
		}
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: msg,
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Retpoline is NOT enabled (x86 Spectre v2 mitigation missing)",
	}
}
