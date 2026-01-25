package rules

import (
	"log/slog"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
)

type Engine struct {
	rules  []rule.Rule
	logger *slog.Logger
}

func NewEngine(logger *slog.Logger) *Engine {
	return &Engine{
		rules:  make([]rule.Rule, 0),
		logger: logger.With(slog.String("component", "rules-engine")),
	}
}

func (e *Engine) LoadRules(ruleIDs []string) {
	e.rules = make([]rule.Rule, 0)

	for _, id := range ruleIDs {
		r := rule.Get(id)
		if r == nil {
			e.logger.Warn("unknown rule ID, skipping", slog.String("rule_id", id))
			continue
		}
		e.rules = append(e.rules, r)
	}
}

func (e *Engine) ExecuteRules(bin *binary.ELFBinary) []rule.ProcessedResult {
	if len(e.rules) == 0 {
		e.logger.Warn("no rules loaded, call LoadRules() first")
		return nil
	}

	results := make([]rule.ProcessedResult, 0, len(e.rules))

	for _, r := range e.rules {
		applicability := r.Applicability()
		if !bin.Architecture.Matches(applicability.Platform.Architecture) {
			continue
		}

		var result rule.ExecuteResult

		switch typedRule := r.(type) {
		case rule.ELFRule:
			if bin.Format != binary.FormatELF {
				continue
			}
			result = typedRule.Execute(bin)
		default:
			continue
		}

		evaluated := rule.ProcessedResult{
			ExecuteResult: result,
			RuleID:        r.ID(),
			Name:          r.Name(),
		}

		if result.Status == rule.StatusFailed {
			evaluated.Suggestion = buildSuggestion(bin.Build.Toolchain, applicability)
		}

		results = append(results, evaluated)
	}

	return results
}
