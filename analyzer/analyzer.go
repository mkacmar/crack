// Package analyzer provides binary security analysis.
// Use this package to programmatically analyze ELF binaries for security hardening.
package analyzer

import (
	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// ApplicabilityResult indicates whether a rule applies to a binary.
type ApplicabilityResult int

const (
	Applicable ApplicabilityResult = iota
	NotApplicableArchitecture
	NotApplicableCompiler
)

// Options configures analyzer behavior.
type Options struct {
	IncludePassed  bool
	IncludeSkipped bool
}

// Analyzer runs rules against parsed binaries.
// Safe for concurrent use after creation.
type Analyzer struct {
	rules          []rule.ELFRule
	includePassed  bool
	includeSkipped bool
}

// NewAnalyzer creates an analyzer with the given rules.
// Pass nil for opts to use defaults. Pass nil or empty rules for a no-op analyzer.
func NewAnalyzer(rules []rule.ELFRule, opts *Options) *Analyzer {
	a := &Analyzer{rules: rules}
	if opts != nil {
		a.includePassed = opts.IncludePassed
		a.includeSkipped = opts.IncludeSkipped
	}
	return a
}

// Analyze runs rules against binary and returns findings.
func (a *Analyzer) Analyze(bin *binary.ELFBinary) []rule.Finding {
	var results []rule.Finding

	for _, r := range a.rules {
		reason := CheckApplicability(r.Applicability(), bin)
		if reason != Applicable {
			if a.includeSkipped {
				results = append(results, rule.Finding{
					Result: rule.Result{
						Status:  rule.StatusSkipped,
						Message: skipMessage(reason, bin),
					},
					RuleID: r.ID(),
					Name:   r.Name(),
				})
			}
			continue
		}

		execResult := r.Execute(bin)
		if execResult.Status == rule.StatusPassed && !a.includePassed {
			continue
		}
		results = append(results, rule.Finding{
			Result: execResult,
			RuleID: r.ID(),
			Name:   r.Name(),
		})
	}

	return results
}

// CheckApplicability determines whether a rule applies to the binary.
func CheckApplicability(app rule.Applicability, bin *binary.ELFBinary) ApplicabilityResult {
	if !bin.Architecture.Matches(app.Platform.Architecture) {
		return NotApplicableArchitecture
	}

	if bin.Build.Compiler != toolchain.Unknown {
		hasCompiler := false
		for comp := range app.Compilers {
			if comp == bin.Build.Compiler {
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

func skipMessage(reason ApplicabilityResult, bin *binary.ELFBinary) string {
	switch reason {
	case NotApplicableArchitecture:
		return "rule not applicable to " + bin.Architecture.String() + " architecture"
	case NotApplicableCompiler:
		return "rule not applicable to " + bin.Build.Compiler.String() + " binaries"
	default:
		return ""
	}
}
