package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// WXorXRule checks for W^X (Write XOR Execute) policy
// OpenBSD: https://www.openbsd.org/papers/auug04/index.html
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
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 3, Minor: 0}, DefaultVersion: model.Version{Major: 3, Minor: 0}, Flag: "-fPIE -pie -z noexecstack"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, DefaultVersion: model.Version{Major: 3, Minor: 0}, Flag: "-fPIE -pie -z noexecstack"},
		},
	}
}

func (r WXorXRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasWriteXorExecute := true

	for _, prog := range f.Progs {
		if prog.Type == elf.PT_LOAD {
			if (prog.Flags&elf.PF_W) != 0 && (prog.Flags&elf.PF_X) != 0 {
				hasWriteXorExecute = false
				break
			}
		}
	}

	if hasWriteXorExecute {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "All memory segments follow W^X policy (no segment is both writable and executable)",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "W^X policy is VIOLATED - some segments may be both writable and executable",
	}
}
