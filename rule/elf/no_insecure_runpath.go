package elf

import (
	"debug/elf"
	"fmt"
	"strings"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// NoInsecureRUNPATHRuleID is the rule ID for secure RUNPATH.
const NoInsecureRUNPATHRuleID = "no-insecure-runpath"

// NoInsecureRUNPATHRule checks for insecure RUNPATH values.
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type NoInsecureRUNPATHRule struct{}

func (r NoInsecureRUNPATHRule) ID() string   { return NoInsecureRUNPATHRuleID }
func (r NoInsecureRUNPATHRule) Name() string { return "Secure RUNPATH" }
func (r NoInsecureRUNPATHRule) Description() string {
	return "Checks for insecure RUNPATH values that could enable library injection. Relative paths, empty components, or world-writable directories in RUNPATH allow attackers to place malicious libraries that get loaded instead of legitimate ones."
}

func (r NoInsecureRUNPATHRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, DefaultVersion: toolchain.Version{Major: 6, Minor: 1}, Flag: "-Wl,--enable-new-dtags -Wl,-rpath,/absolute/path"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, DefaultVersion: toolchain.Version{Major: 4, Minor: 0}, Flag: "-Wl,--enable-new-dtags -Wl,-rpath,/absolute/path"},
			toolchain.Rustc: {MinVersion: toolchain.Version{Major: 1, Minor: 74}, Flag: "-C link-arg=--enable-new-dtags -C link-arg=-rpath -C link-arg=/absolute/path"},
		},
	}
}

func (r NoInsecureRUNPATHRule) Execute(bin *binary.ELFBinary) rule.Result {
	if bin.Type != elf.ET_EXEC && bin.Type != elf.ET_DYN {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	runpath := bin.DynString(elf.DT_RUNPATH)
	if runpath == "" {
		return rule.Result{
			Status:  rule.StatusPassed,
			Message: "No RUNPATH set",
		}
	}

	if insecure := findInsecurePaths(runpath); len(insecure) > 0 {
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: fmt.Sprintf("Insecure RUNPATH: %s", strings.Join(insecure, ", ")),
		}
	}

	return rule.Result{
		Status:  rule.StatusPassed,
		Message: "RUNPATH secure",
	}
}
