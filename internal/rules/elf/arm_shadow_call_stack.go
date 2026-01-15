package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// ARMShadowCallStackRule checks for shadow call stack protection on ARM64/RISC-V
// Clang: https://clang.llvm.org/docs/ShadowCallStack.html
// LLVM: https://llvm.org/docs/ShadowCallStack.html
type ARMShadowCallStackRule struct{}

func (r ARMShadowCallStackRule) ID() string                 { return "arm-shadow-call-stack" }
func (r ARMShadowCallStackRule) Name() string               { return "ARM Shadow Call Stack" }
func (r ARMShadowCallStackRule) Format() model.BinaryFormat { return model.FormatELF }
func (r ARMShadowCallStackRule) FlagType() model.FlagType   { return model.FlagTypeCompile }
func (r ARMShadowCallStackRule) TargetArch() model.Architecture {
	return model.ArchARM64 | model.ArchRISCV
}
func (r ARMShadowCallStackRule) HasPerfImpact() bool { return false }

func (r ARMShadowCallStackRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 7, Minor: 0}, Flag: "-fsanitize=shadow-call-stack"},
		},
	}
}

func (r ARMShadowCallStackRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	symbols, _ := f.Symbols()
	dynsyms, _ := f.DynamicSymbols()

	hasSCS := false
	for _, sym := range symbols {
		if strings.Contains(sym.Name, "shadow_call_stack") {
			hasSCS = true
			break
		}
	}
	if !hasSCS {
		for _, sym := range dynsyms {
			if strings.Contains(sym.Name, "shadow_call_stack") {
				hasSCS = true
				break
			}
		}
	}

	if hasSCS {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Shadow call stack is enabled (ARM64/RISC-V)",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Shadow call stack is NOT enabled",
	}
}
