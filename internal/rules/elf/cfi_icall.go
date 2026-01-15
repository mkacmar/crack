package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// CFIICallRule checks for CFI indirect call protection
// Clang: https://clang.llvm.org/docs/ControlFlowIntegrity.html
type CFIICallRule struct{}

func (r CFIICallRule) ID() string                     { return "cfi-icall" }
func (r CFIICallRule) Name() string                   { return "CFI - Indirect Call Check" }
func (r CFIICallRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r CFIICallRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r CFIICallRule) TargetArch() model.Architecture { return model.ArchAll }
func (r CFIICallRule) HasPerfImpact() bool            { return true }

func (r CFIICallRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 9}, Flag: "-fsanitize=cfi-icall -flto"},
		},
	}
}

func (r CFIICallRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	symbols, _ := f.Symbols()
	dynsyms, _ := f.DynamicSymbols()

	hasCFIICall := false
	for _, sym := range symbols {
		if strings.Contains(sym.Name, "__cfi_check") {
			hasCFIICall = true
			break
		}
	}
	if !hasCFIICall {
		for _, sym := range dynsyms {
			if strings.Contains(sym.Name, "__cfi_check") {
				hasCFIICall = true
				break
			}
		}
	}

	if hasCFIICall {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "CFI for indirect calls is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "CFI for indirect calls is NOT enabled",
	}
}
