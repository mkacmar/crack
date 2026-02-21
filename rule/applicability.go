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
)

// String returns a human-readable skip message for non-applicable results.
func (r ApplicabilityResult) String() string {
	switch r {
	case NotApplicableArchitecture:
		return "architecture not applicable"
	case NotApplicableCompiler:
		return "compiler not applicable"
	default:
		return ""
	}
}

// SkipMessage returns a detailed skip message including binary metadata.
func (r ApplicabilityResult) SkipMessage(info binary.Info) string {
	switch r {
	case NotApplicableArchitecture:
		return "rule not applicable to " + info.Architecture.String() + " architecture"
	case NotApplicableCompiler:
		return "rule not applicable to " + info.Build.Compiler.String() + " binaries"
	default:
		return ""
	}
}

// CheckApplicability determines whether a rule applies to the binary.
func CheckApplicability(app Applicability, info binary.Info) ApplicabilityResult {
	if !info.Architecture.Matches(app.Platform.Architecture) {
		return NotApplicableArchitecture
	}

	if info.Build.Compiler != toolchain.Unknown {
		hasCompiler := false
		for comp := range app.Compilers {
			if comp == info.Build.Compiler {
				hasCompiler = true
				break
			}
		}
		if len(app.Compilers) > 0 && !hasCompiler {
			return NotApplicableCompiler
		}
	}

	return Applicable
}
