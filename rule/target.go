package rule

import (
	"slices"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/toolchain"
)

// PlatformTarget specifies an architecture constraint for filtering rules.
type PlatformTarget struct {
	Architecture binary.Architecture
	MaxISA       *binary.ISA
}

// CompilerTarget specifies a compiler constraint for filtering rules.
type CompilerTarget struct {
	Compiler   toolchain.Compiler
	MaxVersion *toolchain.Version
}

// TargetFilter selects rules based on platform and compiler constraints.
type TargetFilter struct {
	Platforms []PlatformTarget
	Compilers []CompilerTarget
}

func (f *TargetFilter) isEmpty() bool {
	return len(f.Platforms) == 0 && len(f.Compilers) == 0
}

func (f *TargetFilter) matches(app Applicability) bool {
	return matchesAnyTarget(f.Platforms, app.matchesPlatform) &&
		matchesAnyTarget(f.Compilers, app.matchesCompiler)
}

func matchesAnyTarget[T any](targets []T, matches func(T) bool) bool {
	return len(targets) == 0 || slices.ContainsFunc(targets, matches)
}

func (app Applicability) matchesPlatform(pt PlatformTarget) bool {
	if !app.Platform.Architecture.Matches(pt.Architecture) {
		return false
	}
	return pt.MaxISA == nil || app.Platform.MinISA.Major == 0 || pt.MaxISA.IsAtLeast(app.Platform.MinISA)
}

func (app Applicability) matchesCompiler(ct CompilerTarget) bool {
	req, ok := app.Compilers[ct.Compiler]
	if !ok {
		return false
	}
	return ct.MaxVersion == nil || req.MinVersion.Major == 0 || ct.MaxVersion.IsAtLeast(req.MinVersion)
}

// FilterRules returns only rules matching the filter. Returns all rules if filter is nil.
func FilterRules[T Rule](rules []T, filter *TargetFilter) []T {
	if filter == nil || filter.isEmpty() {
		return rules
	}

	var filtered []T
	for _, r := range rules {
		if filter.matches(r.Applicability()) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
