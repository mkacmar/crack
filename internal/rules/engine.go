package rules

import (
	"log/slog"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/internal/analyzer"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

type Engine struct {
	rules  []rule.Rule
	logger *slog.Logger
}

func NewEngine(rules []rule.Rule, logger *slog.Logger) *Engine {
	return &Engine{
		rules:  rules,
		logger: logger.With(slog.String("component", "rules-engine")),
	}
}

func skipFinding(r rule.Rule, message string) analyzer.FindingWithSuggestion {
	return analyzer.FindingWithSuggestion{
		Finding: rule.Finding{
			Result: rule.Result{
				Status:  rule.StatusSkipped,
				Message: message,
			},
			RuleID: r.ID(),
			Name:   r.Name(),
		},
	}
}

func (e *Engine) ExecuteRules(bin *binary.ELFBinary) []analyzer.FindingWithSuggestion {
	findings := make([]analyzer.FindingWithSuggestion, 0, len(e.rules))

	for _, r := range e.rules {
		elfRule, ok := r.(rule.ELFRule)
		if !ok {
			continue
		}

		applicability := r.Applicability()
		if !bin.Architecture.Matches(applicability.Platform.Architecture) {
			findings = append(findings, skipFinding(r, "rule not applicable to "+bin.Architecture.String()+" architecture"))
			continue
		}

		if len(applicability.Compilers) > 0 && bin.Build.Compiler != toolchain.Unknown {
			hasCompiler := false
			for comp := range applicability.Compilers {
				if comp == bin.Build.Compiler {
					hasCompiler = true
					break
				}
			}
			if !hasCompiler {
				findings = append(findings, skipFinding(r, "rule not applicable to "+bin.Build.Compiler.String()+" binaries"))
				continue
			}
		}

		result := elfRule.Execute(bin)

		finding := analyzer.FindingWithSuggestion{
			Finding: rule.Finding{
				Result: result,
				RuleID: r.ID(),
				Name:   r.Name(),
			},
		}

		if result.Status == rule.StatusFailed {
			finding.Suggestion = buildSuggestion(bin.Build, applicability)
		}

		findings = append(findings, finding)
	}

	return findings
}
