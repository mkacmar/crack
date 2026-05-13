package elf

import (
	stdelf "debug/elf"

	"fmt"
	"strings"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/binary/elf"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
)

// NoInsecureRUNPATHRuleID is the rule ID for secure RUNPATH.
const NoInsecureRUNPATHRuleID = "no-insecure-runpath"

// NoInsecureRUNPATHRule checks for insecure RUNPATH values.
//
// References:
//   - https://sourceware.org/binutils/docs/ld/Options.html
type NoInsecureRUNPATHRule struct{}

func (r NoInsecureRUNPATHRule) ID() string   { return NoInsecureRUNPATHRuleID }
func (r NoInsecureRUNPATHRule) Name() string { return "Secure RUNPATH" }
func (r NoInsecureRUNPATHRule) Description() string {
	return "Checks for insecure RUNPATH values that could enable library injection. Relative paths, empty components, or world-writable directories in RUNPATH allow attackers to place malicious libraries that get loaded instead of legitimate ones."
}

func (r NoInsecureRUNPATHRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM | binary.ArchRISCV},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, DefaultVersion: toolchain.Version{Major: 6, Minor: 1}, Flag: "-Wl,--enable-new-dtags -Wl,-rpath,/absolute/path"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, DefaultVersion: toolchain.Version{Major: 4, Minor: 0}, Flag: "-Wl,--enable-new-dtags -Wl,-rpath,/absolute/path"},
		},
		LibC: binary.LibCAll,
	}
}

func (r NoInsecureRUNPATHRule) Execute(bin elf.Binary) rule.Result {
	if bin.Type() != stdelf.ET_EXEC && bin.Type() != stdelf.ET_DYN {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	runpath, err := elf.DynString(bin, stdelf.DT_RUNPATH)
	if err != nil {
		return rule.Skip("failed to read dynamic section", err)
	}
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
