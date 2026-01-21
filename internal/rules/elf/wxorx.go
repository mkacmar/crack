package elf

import (
	"debug/elf"
	"fmt"

	"github.com/mkacmar/crack/internal/model"
)

// WXorXRule checks for W^X (Write XOR Execute) policy
// GNU ld: https://sourceware.org/binutils/docs/ld/Options.html (-z noexecstack)
type WXorXRule struct{}

func (r WXorXRule) ID() string                     { return "wxorx" }
func (r WXorXRule) Name() string                   { return "W^X (Write XOR Execute)" }
func (r WXorXRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r WXorXRule) FlagType() model.FlagType       { return model.FlagTypeLink }
func (r WXorXRule) TargetArch() model.Architecture { return model.ArchAll }
func (r WXorXRule) HasPerfImpact() bool            { return false }

func (r WXorXRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 3, Minor: 0}, DefaultVersion: model.Version{Major: 3, Minor: 0}, Flag: "-z noexecstack"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, DefaultVersion: model.Version{Major: 3, Minor: 0}, Flag: "-z noexecstack"},
		},
	}
}

func (r WXorXRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	for _, prog := range f.Progs {
		// Check PT_LOAD segments for W+X
		if prog.Type == elf.PT_LOAD {
			if (prog.Flags&elf.PF_W) != 0 && (prog.Flags&elf.PF_X) != 0 {
				return model.RuleResult{
					State:   model.CheckStateFailed,
					Message: fmt.Sprintf("W^X violation: segment at offset 0x%x is both writable and executable", prog.Off),
				}
			}
		}
		// Check PT_GNU_STACK for executable stack
		if prog.Type == elf.PT_GNU_STACK && (prog.Flags&elf.PF_X) != 0 {
			return model.RuleResult{
				State:   model.CheckStateFailed,
				Message: "W^X violation: executable stack",
			}
		}
	}

	return model.RuleResult{
		State:   model.CheckStatePassed,
		Message: "All memory segments follow W^X policy (no segment is both writable and executable)",
	}
}
