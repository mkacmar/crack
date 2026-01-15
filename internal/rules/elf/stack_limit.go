package elf

import (
	"debug/elf"
	"fmt"

	"github.com/mkacmar/crack/internal/model"
)

// StackLimitRule checks for explicit stack size limit
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type StackLimitRule struct{}

func (r StackLimitRule) ID() string                     { return "stack-limit" }
func (r StackLimitRule) Name() string                   { return "Explicit Stack Size Limit" }
func (r StackLimitRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r StackLimitRule) FlagType() model.FlagType       { return model.FlagTypeLink }
func (r StackLimitRule) TargetArch() model.Architecture { return model.ArchAll }
func (r StackLimitRule) HasPerfImpact() bool            { return false }

func (r StackLimitRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,stack-size=8388608"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,stack-size=8388608"},
		},
	}
}

func (r StackLimitRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasExplicitStackLimit := false
	stackSize := uint64(0)
	foundGNUStack := false

	for _, prog := range f.Progs {
		if prog.Type == elf.PT_GNU_STACK {
			foundGNUStack = true
			stackSize = prog.Memsz
			hasExplicitStackLimit = stackSize > 0
			break
		}
	}

	if hasExplicitStackLimit {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: fmt.Sprintf("Stack has explicit size limit: %d bytes", stackSize),
		}
	}

	msg := "No PT_GNU_STACK segment found"
	if foundGNUStack {
		msg = "Stack uses system default size (no explicit limit set)"
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: msg,
	}
}
