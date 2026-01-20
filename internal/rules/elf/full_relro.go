package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

const (
	DF_BIND_NOW = 0x8
	DF_1_NOW    = 0x1
)

// FullRELRORule checks for full RELRO protection
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type FullRELRORule struct{}

func (r FullRELRORule) ID() string                     { return "full-relro" }
func (r FullRELRORule) Name() string                   { return "Full RELRO" }
func (r FullRELRORule) Format() model.BinaryFormat     { return model.FormatELF }
func (r FullRELRORule) FlagType() model.FlagType       { return model.FlagTypeLink }
func (r FullRELRORule) TargetArch() model.Architecture { return model.ArchAll }
func (r FullRELRORule) HasPerfImpact() bool            { return false }

func (r FullRELRORule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,relro,-z,now"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,relro,-z,now"},
		},
	}
}

func (r FullRELRORule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasRELRO := false
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_GNU_RELRO {
			hasRELRO = true
			break
		}
	}

	if !hasRELRO {
		return model.RuleResult{
			State:   model.CheckStateFailed,
			Message: "Full RELRO is NOT enabled (no PT_GNU_RELRO segment)",
		}
	}

	// Check for BIND_NOW flag which makes RELRO "full" by disabling lazy binding.
	// This can be indicated by DT_BIND_NOW, DT_FLAGS with DF_BIND_NOW, or
	// DT_FLAGS_1 with DF_1_NOW.
	if HasDynTag(f, elf.DT_BIND_NOW) ||
		HasDynFlag(f, elf.DT_FLAGS, DF_BIND_NOW) ||
		HasDynFlag(f, elf.DT_FLAGS_1, DF_1_NOW) {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Full RELRO is enabled (GOT read-only, lazy binding disabled)",
		}
	}

	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Full RELRO is NOT enabled (GOT may be writable)",
	}
}
