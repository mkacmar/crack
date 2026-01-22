package rules

import (
	"log/slog"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/preset"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/rules/elf"
)

type Engine struct {
	rules  []rule.Rule
	logger *slog.Logger
}

func NewEngine(logger *slog.Logger) *Engine {
	return &Engine{
		rules:  make([]rule.Rule, 0),
		logger: logger.With(slog.String("component", "rules")),
	}
}

func (e *Engine) LoadPreset(p preset.Preset) {
	e.rules = make([]rule.Rule, 0)

	for _, id := range p.Rules {
		if r, ok := elf.AllRules[id]; ok {
			e.rules = append(e.rules, r)
		} else {
			e.logger.Warn("unknown rule ID in preset, skipping", slog.String("rule_id", id))
		}
	}
}

func (e *Engine) ExecuteRules(info *binary.Parsed) []rule.ProcessedResult {
	if len(e.rules) == 0 {
		e.logger.Warn("no rules loaded, call LoadPreset() first")
		return nil
	}

	results := make([]rule.ProcessedResult, 0, len(e.rules))

	for _, r := range e.rules {
		applicability := r.Applicability()
		if !info.Architecture.Matches(applicability.Arch) {
			continue
		}

		var result rule.ExecuteResult
		switch elfRule := r.(type) {
		case rule.ELFRule:
			if info.Format != binary.FormatELF {
				continue
			}
			result = elfRule.Execute(info.ELF, info)
		default:
			continue
		}

		evaluated := rule.ProcessedResult{
			ExecuteResult: result,
			RuleID:        r.ID(),
			Name:          r.Name(),
		}

		if result.Status == rule.StatusFailed {
			evaluated.Suggestion = buildSuggestion(info.Build.Toolchain, applicability)
		}

		results = append(results, evaluated)
	}

	return results
}
