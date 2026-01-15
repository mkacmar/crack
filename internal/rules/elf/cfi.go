package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// CFIRule checks for Clang Control Flow Integrity
// Clang: https://clang.llvm.org/docs/ControlFlowIntegrity.html
type CFIRule struct{}

func (r CFIRule) ID() string                     { return "cfi" }
func (r CFIRule) Name() string                   { return "Control Flow Integrity" }
func (r CFIRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r CFIRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r CFIRule) TargetArch() model.Architecture { return model.ArchAll }
func (r CFIRule) HasPerfImpact() bool            { return true }

func (r CFIRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 7}, Flag: "-fsanitize=cfi -flto -fvisibility=hidden"},
		},
	}
}

func (r CFIRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasCFI := false

	symbols, err := f.Symbols()
	if err != nil {
		symbols = nil
	}

	dynsyms, err := f.DynamicSymbols()
	if err != nil {
		dynsyms = nil
	}

	for _, sym := range symbols {
		if strings.Contains(sym.Name, "__cfi_check") ||
			strings.Contains(sym.Name, "__cfi_slowpath") ||
			strings.Contains(sym.Name, "__ubsan_handle_cfi_check_fail") {
			hasCFI = true
			break
		}
	}

	if !hasCFI {
		for _, sym := range dynsyms {
			if strings.Contains(sym.Name, "__cfi_check") ||
				strings.Contains(sym.Name, "__cfi_slowpath") ||
				strings.Contains(sym.Name, "__ubsan_handle_cfi_check_fail") {
				hasCFI = true
				break
			}
		}
	}

	if hasCFI {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Clang CFI features are enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Clang CFI features are NOT enabled",
	}
}
