package elf

import (
	"debug/elf"
	"fmt"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// NoInsecureRUNPATHRule checks for insecure RUNPATH values
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type NoInsecureRUNPATHRule struct{}

func (r NoInsecureRUNPATHRule) ID() string                     { return "no-insecure-runpath" }
func (r NoInsecureRUNPATHRule) Name() string                   { return "Secure RUNPATH" }
func (r NoInsecureRUNPATHRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r NoInsecureRUNPATHRule) FlagType() model.FlagType       { return model.FlagTypeLink }
func (r NoInsecureRUNPATHRule) TargetArch() model.Architecture { return model.ArchAll }
func (r NoInsecureRUNPATHRule) HasPerfImpact() bool            { return false }

func (r NoInsecureRUNPATHRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,--enable-new-dtags -Wl,-rpath,/absolute/path"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,--enable-new-dtags -Wl,-rpath,/absolute/path"},
		},
	}
}

func (r NoInsecureRUNPATHRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	runpath := GetDynString(f, elf.DT_RUNPATH)
	if runpath == "" {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "No RUNPATH set",
		}
	}

	if insecure := findInsecurePaths(runpath); len(insecure) > 0 {
		return model.RuleResult{
			State:   model.CheckStateFailed,
			Message: fmt.Sprintf("Insecure RUNPATH: %s", strings.Join(insecure, ", ")),
		}
	}

	return model.RuleResult{
		State:   model.CheckStatePassed,
		Message: "RUNPATH is secure",
	}
}
