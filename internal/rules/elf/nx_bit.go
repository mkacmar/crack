package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// NXBitRule checks for non-executable stack
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Link-Options.html#index-z
type NXBitRule struct{}

func (r NXBitRule) ID() string                     { return "nx-bit" }
func (r NXBitRule) Name() string                   { return "Non-Executable Stack" }
func (r NXBitRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r NXBitRule) FlagType() model.FlagType       { return model.FlagTypeLink }
func (r NXBitRule) TargetArch() model.Architecture { return model.ArchAll }
func (r NXBitRule) HasPerfImpact() bool            { return false }

func (r NXBitRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 3, Minor: 0}, DefaultVersion: model.Version{Major: 3, Minor: 0}, Flag: "-z noexecstack"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, DefaultVersion: model.Version{Major: 3, Minor: 0}, Flag: "-z noexecstack"},
		},
	}
}

func (r NXBitRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_GNU_STACK {
			if (prog.Flags & elf.PF_X) == 0 {
				return model.RuleResult{
					State:   model.CheckStatePassed,
					Message: "Stack is marked as non-executable (NX bit enabled)",
				}
			}
			return model.RuleResult{
				State:   model.CheckStateFailed,
				Message: "Stack is EXECUTABLE (NX bit NOT enabled)",
			}
		}
	}

	// PT_GNU_STACK missing - this typically means executable stack on older systems.
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "PT_GNU_STACK segment missing (stack executability depends on system defaults)",
	}
}
