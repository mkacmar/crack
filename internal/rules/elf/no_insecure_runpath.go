package elf

import (
	"debug/elf"
	"fmt"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const NoInsecureRUNPATHRuleID = "no-insecure-runpath"

// NoInsecureRUNPATHRule checks for insecure RUNPATH values
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type NoInsecureRUNPATHRule struct{}

func (r NoInsecureRUNPATHRule) ID() string   { return NoInsecureRUNPATHRuleID }
func (r NoInsecureRUNPATHRule) Name() string { return "Secure RUNPATH" }

func (r NoInsecureRUNPATHRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, DefaultVersion: toolchain.Version{Major: 6, Minor: 1}, Flag: "-Wl,--enable-new-dtags -Wl,-rpath,/absolute/path"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, DefaultVersion: toolchain.Version{Major: 4, Minor: 0}, Flag: "-Wl,--enable-new-dtags -Wl,-rpath,/absolute/path"},
		},
	}
}

func (r NoInsecureRUNPATHRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	runpath := GetDynString(bin.File, elf.DT_RUNPATH)
	if runpath == "" {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "No RUNPATH set",
		}
	}

	if insecure := findInsecurePaths(runpath); len(insecure) > 0 {
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: fmt.Sprintf("Insecure RUNPATH: %s", strings.Join(insecure, ", ")),
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusPassed,
		Message: "RUNPATH secure",
	}
}
