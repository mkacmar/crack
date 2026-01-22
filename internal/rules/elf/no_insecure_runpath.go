package elf

import (
	"debug/elf"
	"fmt"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

// NoInsecureRUNPATHRule checks for insecure RUNPATH values
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type NoInsecureRUNPATHRule struct{}

func (r NoInsecureRUNPATHRule) ID() string   { return "no-insecure-runpath" }
func (r NoInsecureRUNPATHRule) Name() string { return "Secure RUNPATH" }

func (r NoInsecureRUNPATHRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Arch: binary.ArchAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-Wl,--enable-new-dtags -Wl,-rpath,/absolute/path"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-Wl,--enable-new-dtags -Wl,-rpath,/absolute/path"},
		},
	}
}

func (r NoInsecureRUNPATHRule) Execute(f *elf.File, info *binary.Parsed) rule.Result {
	runpath := GetDynString(f, elf.DT_RUNPATH)
	if runpath == "" {
		return rule.Result{
			State:   rule.CheckStatePassed,
			Message: "No RUNPATH set",
		}
	}

	if insecure := findInsecurePaths(runpath); len(insecure) > 0 {
		return rule.Result{
			State:   rule.CheckStateFailed,
			Message: fmt.Sprintf("Insecure RUNPATH: %s", strings.Join(insecure, ", ")),
		}
	}

	return rule.Result{
		State:   rule.CheckStatePassed,
		Message: "RUNPATH is secure",
	}
}
