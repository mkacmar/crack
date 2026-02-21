package rule

import (
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
	if f.isEmpty() {
		return true
	}

	for _, pt := range f.Platforms {
		if !app.Platform.Architecture.Matches(pt.Architecture) {
			return false
		}
		if pt.MaxISA != nil && app.Platform.MinISA.Major > 0 {
			if !pt.MaxISA.IsAtLeast(app.Platform.MinISA) {
				return false
			}
		}
	}

	for _, ct := range f.Compilers {
		hasCompiler := false
		for comp, req := range app.Compilers {
			if comp == ct.Compiler {
				hasCompiler = true
				if ct.MaxVersion != nil && req.MinVersion.Major > 0 {
					if !ct.MaxVersion.IsAtLeast(req.MinVersion) {
						return false
					}
				}
				break
			}
		}
		if !hasCompiler {
			return false
		}
	}

	return true
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
