package elf

import (
	"debug/elf"
	"fmt"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const NoInsecureRPATHRuleID = "no-insecure-rpath"

// NoInsecureRPATHRule checks for insecure RPATH values
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type NoInsecureRPATHRule struct{}

func (r NoInsecureRPATHRule) ID() string   { return NoInsecureRPATHRuleID }
func (r NoInsecureRPATHRule) Name() string { return "Secure RPATH" }
func (r NoInsecureRPATHRule) Description() string {
	return "Checks for insecure RPATH values that could enable library injection. RPATH takes precedence over system library paths, so relative paths or world-writable directories allow attackers to hijack library loading."
}

func (r NoInsecureRPATHRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, Flag: "-Wl,-rpath,/absolute/path"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, Flag: "-Wl,-rpath,/absolute/path"},
			toolchain.CompilerRustc: {MinVersion: toolchain.Version{Major: 1, Minor: 0}, Flag: "-C link-arg=-rpath -C link-arg=/absolute/path"},
		},
	}
}

func (r NoInsecureRPATHRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	if bin.File.Type != elf.ET_EXEC && bin.File.Type != elf.ET_DYN {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	rpath := GetDynString(bin.File, elf.DT_RPATH)
	if rpath == "" {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "No RPATH set",
		}
	}

	if insecure := findInsecurePaths(rpath); len(insecure) > 0 {
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: fmt.Sprintf("Insecure RPATH: %s", strings.Join(insecure, ", ")),
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusPassed,
		Message: "RPATH secure",
	}
}

func findInsecurePaths(rpath string) []string {
	var insecure []string
	for _, p := range strings.Split(rpath, ":") {
		if isInsecurePath(p) {
			if p == "" {
				insecure = append(insecure, "(empty)")
			} else {
				insecure = append(insecure, p)
			}
		}
	}
	return insecure
}

// CWE-426: Untrusted Search Path (https://cwe.mitre.org/data/definitions/426.html)
func isInsecurePath(p string) bool {
	if p == "" {
		return true // empty path component means current directory
	}

	// Relative paths resolve from CWD
	if !strings.HasPrefix(p, "/") && !strings.HasPrefix(p, "$") {
		return true
	}

	// Attacker can control binary location and use ../ to reach writable directories
	if strings.HasPrefix(p, "$ORIGIN") && strings.Contains(p, "..") {
		return true
	}

	worldWritable := []string{"/tmp", "/var/tmp", "/dev/shm"}
	for _, ww := range worldWritable {
		if p == ww || strings.HasPrefix(p, ww+"/") {
			return true
		}
	}

	return false
}
