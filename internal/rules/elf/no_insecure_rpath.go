package elf

import (
	"debug/elf"
	"fmt"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

// NoInsecureRPATHRule checks for insecure RPATH values
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type NoInsecureRPATHRule struct{}

func (r NoInsecureRPATHRule) ID() string                 { return "no-insecure-rpath" }
func (r NoInsecureRPATHRule) Name() string               { return "Secure RPATH" }

func (r NoInsecureRPATHRule) Applicability() model.Applicability {
	return model.Applicability{
		Arch: model.ArchAll,
		Compilers: map[model.Compiler]model.CompilerRequirement{
			model.CompilerGCC:   {MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-rpath,/absolute/path"},
			model.CompilerClang: {MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-Wl,-rpath,/absolute/path"},
		},
	}
}

func (r NoInsecureRPATHRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	rpath := GetDynString(f, elf.DT_RPATH)
	if rpath == "" {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "No RPATH set",
		}
	}

	if insecure := findInsecurePaths(rpath); len(insecure) > 0 {
		return model.RuleResult{
			State:   model.CheckStateFailed,
			Message: fmt.Sprintf("Insecure RPATH: %s", strings.Join(insecure, ", ")),
		}
	}

	return model.RuleResult{
		State:   model.CheckStatePassed,
		Message: "RPATH is secure",
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

// See CWE-426: Untrusted Search Path (https://cwe.mitre.org/data/definitions/426.html)
func isInsecurePath(p string) bool {
	if p == "" {
		return true // empty path component means current directory
	}

	// Relative paths
	if p == "." || p == ".." || strings.HasPrefix(p, "./") || strings.HasPrefix(p, "../") {
		return true
	}

	// World-writable directories
	worldWritable := []string{"/tmp", "/var/tmp"}
	for _, ww := range worldWritable {
		if p == ww || strings.HasPrefix(p, ww+"/") {
			return true
		}
	}

	return false
}
