package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// CFIVCallRule checks for CFI virtual call protection
// Clang: https://clang.llvm.org/docs/ControlFlowIntegrity.html
type CFIVCallRule struct{}

func (r CFIVCallRule) ID() string                     { return "cfi-vcall" }
func (r CFIVCallRule) Name() string                   { return "CFI - Virtual Call Check" }
func (r CFIVCallRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r CFIVCallRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r CFIVCallRule) TargetArch() model.Architecture { return model.ArchAll }
func (r CFIVCallRule) HasPerfImpact() bool            { return true }

func (r CFIVCallRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 9}, Flag: "-fsanitize=cfi-vcall -flto"},
		},
	}
}

func (r CFIVCallRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	symbols, _ := f.Symbols()
	dynsyms, _ := f.DynamicSymbols()

	hasCFIVCall := false
	for _, sym := range symbols {
		if strings.Contains(sym.Name, "__cfi_check") {
			hasCFIVCall = true
			break
		}
	}
	if !hasCFIVCall {
		for _, sym := range dynsyms {
			if strings.Contains(sym.Name, "__cfi_check") {
				hasCFIVCall = true
				break
			}
		}
	}

	if hasCFIVCall {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "CFI for virtual calls is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "CFI for virtual calls is NOT enabled",
	}
}
