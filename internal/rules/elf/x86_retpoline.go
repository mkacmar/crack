package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const X86RetpolineRuleID = "x86-retpoline"

// X86RetpolineRule checks for Spectre v2 mitigation (retpoline)
// GCC: https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html#index-mindirect-branch
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-mretpoline
type X86RetpolineRule struct{}

func (r X86RetpolineRule) ID() string   { return X86RetpolineRuleID }
func (r X86RetpolineRule) Name() string { return "x86 Retpoline (Spectre v2)" }

func (r X86RetpolineRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAllX86,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 7, Minor: 3}, Flag: "-mindirect-branch=thunk -mfunction-return=thunk"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 5, Minor: 0}, Flag: "-mretpoline"},
		},
	}
}

func (r X86RetpolineRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	// CET-IBT and retpoline both mitigate Spectre v2 indirect branch attacks.
	// Hardware mitigations take precedence over software (retpoline).
	// See: https://www.kernel.org/doc/html/latest/admin-guide/hw-vuln/spectre.html
	hasCETIBT := parseGNUProperty(bin.File, GNU_PROPERTY_X86_FEATURE_1_AND, GNU_PROPERTY_X86_FEATURE_1_IBT)
	if hasCETIBT {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "CET-IBT enabled (hardware mitigation supersedes retpoline)",
		}
	}

	// Retpoline thunks are typically local symbols in .symtab.
	// Stripped binaries lose .symtab, making detection impossible.
	hasSymtab := false
	for _, sec := range bin.File.Sections {
		if sec.Type == elf.SHT_SYMTAB {
			hasSymtab = true
			break
		}
	}
	if !hasSymtab {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Stripped binary (retpoline symbols not available)",
		}
	}

	var hasGCCThunk, hasLLVMRetpoline bool
	for _, sym := range append(bin.Symbols, bin.DynSymbols...) {
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
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: msg,
		}
	}
	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "Retpoline is NOT enabled (x86 Spectre v2 mitigation missing)",
	}
}
