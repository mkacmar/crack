package rule

import (
	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/toolchain"
)

// ApplicabilityResult indicates whether a rule applies to a binary.
type ApplicabilityResult int

const (
	Applicable ApplicabilityResult = iota
	NotApplicableArchitecture
	NotApplicableCompiler
	NotApplicableLibC
)

// String returns a human-readable skip message for non-applicable results.
func (r ApplicabilityResult) String() string {
	switch r {
	case NotApplicableArchitecture:
		return "architecture not applicable"
	case NotApplicableCompiler:
		return "compiler not applicable"
	case NotApplicableLibC:
		return "libc not applicable"
	default:
		return ""
	}
}

// Reason describes why the rule is not applicable to the given classification.
func (r ApplicabilityResult) Reason(profile binary.Profile) string {
	switch r {
	case NotApplicableArchitecture:
		return "rule not applicable to " + profile.Architecture.String() + " architecture"
	case NotApplicableCompiler:
		return "rule not applicable to " + profile.Toolchain.Compiler.String() + " binaries"
	case NotApplicableLibC:
		return "rule not applicable to " + profile.LibC.String() + " binaries"
	default:
		return ""
	}
}

// CheckApplicability determines whether a rule applies to the binary.
// When detection of an optional axis yields the Unknown sentinel (compiler, libc), the axis is skipped in the filter and the rule runs as best-effort.
// Architecture has no such bypass because ELF machine detection cannot fail in practice.
func CheckApplicability(app Applicability, profile binary.Profile) ApplicabilityResult {
	if !profile.Architecture.Matches(app.Platform.Architecture) {
		return NotApplicableArchitecture
	}

	if profile.Toolchain.Compiler != toolchain.Unknown {
		hasCompiler := false
		for comp := range app.Compilers {
			if comp == profile.Toolchain.Compiler {
				hasCompiler = true
				break
			}
		}
		if len(app.Compilers) > 0 && !hasCompiler {
			return NotApplicableCompiler
		}
	}

	if profile.LibC != binary.LibCUnknown && !app.LibC.Matches(profile.LibC) {
		return NotApplicableLibC
	}

	return Applicable
}
