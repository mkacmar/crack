package elf

import (
	"debug/elf"
	"fmt"
	"path"
	"strings"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// NoInsecureRPATHRuleID is the rule ID for secure RPATH.
const NoInsecureRPATHRuleID = "no-insecure-rpath"

// NoInsecureRPATHRule checks for insecure RPATH values.
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type NoInsecureRPATHRule struct{}

func (r NoInsecureRPATHRule) ID() string   { return NoInsecureRPATHRuleID }
func (r NoInsecureRPATHRule) Name() string { return "Secure RPATH" }
func (r NoInsecureRPATHRule) Description() string {
	return "Checks for insecure RPATH values that could enable library injection. RPATH takes precedence over system library paths, so relative paths or world-writable directories allow attackers to hijack library loading."
}

func (r NoInsecureRPATHRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, Flag: "-Wl,-rpath,/absolute/path"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, Flag: "-Wl,-rpath,/absolute/path"},
			toolchain.Rustc: {MinVersion: toolchain.Version{Major: 1, Minor: 0}, Flag: "-C link-arg=-rpath -C link-arg=/absolute/path"},
		},
	}
}

func (r NoInsecureRPATHRule) Execute(bin *binary.ELFBinary) rule.Result {
	if bin.Type != elf.ET_EXEC && bin.Type != elf.ET_DYN {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	rpath := bin.DynString(elf.DT_RPATH)
	if rpath == "" {
		return rule.Result{
			Status:  rule.StatusPassed,
			Message: "No RPATH set",
		}
	}

	if insecure := findInsecurePaths(rpath); len(insecure) > 0 {
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: fmt.Sprintf("Insecure RPATH: %s", strings.Join(insecure, ", ")),
		}
	}

	return rule.Result{
		Status:  rule.StatusPassed,
		Message: "RPATH secure",
	}
}

func findInsecurePaths(rpath string) []string {
	var insecure []string
	hasEmpty := false
	for _, p := range strings.Split(rpath, ":") {
		if isInsecurePath(p) {
			if p == "" {
				if !hasEmpty {
					insecure = append(insecure, "(empty)")
					hasEmpty = true
				}
			} else {
				insecure = append(insecure, p)
			}
		}
	}
	return insecure
}

func isInsecurePath(p string) bool {
	if p == "" {
		return true
	}

	if !strings.HasPrefix(p, "/") && !strings.HasPrefix(p, "$") {
		return true
	}

	if strings.HasPrefix(p, "$ORIGIN") && strings.Contains(p, "..") {
		// Allow $ORIGIN/../lib and $ORIGIN/../lib64 — standard layout for co-installed binaries and libraries.
		base := path.Base(p)
		if base != "lib" && base != "lib64" {
			return true
		}
	}

	worldWritable := []string{"/tmp", "/var/tmp", "/dev/shm"}
	for _, ww := range worldWritable {
		if p == ww || strings.HasPrefix(p, ww+"/") {
			return true
		}
	}

	return false
}
