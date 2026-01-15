package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// RELRORule checks for partial RELRO
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type RELRORule struct{}

func (r RELRORule) ID() string                     { return "relro" }
func (r RELRORule) Name() string                   { return "Partial RELRO" }
func (r RELRORule) Format() model.BinaryFormat     { return model.FormatELF }
func (r RELRORule) FlagType() model.FlagType       { return model.FlagTypeLink }
func (r RELRORule) TargetArch() model.Architecture { return model.ArchAll }
func (r RELRORule) HasPerfImpact() bool            { return false }

func (r RELRORule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 3, Minor: 0}, DefaultVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,relro"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, DefaultVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,relro"},
		},
	}
}

func (r RELRORule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasRELRO := false
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_GNU_RELRO {
			hasRELRO = true
			break
		}
	}

	if hasRELRO {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Partial RELRO is enabled (some ELF sections read-only after load)",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Partial RELRO is NOT enabled (ELF sections remain writable)",
	}
}
