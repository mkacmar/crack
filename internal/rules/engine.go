package rules

import (
	"log/slog"

	"github.com/mkacmar/crack/internal/model"
	"github.com/mkacmar/crack/internal/preset"
	"github.com/mkacmar/crack/internal/rules/elf"
)

type Engine struct {
	rules  []model.Rule
	logger *slog.Logger
}

func NewEngine(logger *slog.Logger) *Engine {
	return &Engine{
		rules:  make([]model.Rule, 0),
		logger: logger.With(slog.String("component", "rules")),
	}
}

func (e *Engine) LoadPreset(p preset.Preset) {
	e.rules = make([]model.Rule, 0)

	for _, id := range p.Rules {
		if rule, ok := elf.AllRules[id]; ok {
			e.rules = append(e.rules, rule)
		} else {
			e.logger.Warn("unknown rule ID in preset, skipping", slog.String("rule_id", id))
		}
	}
}

func (e *Engine) ExecuteRules(info *model.ParsedBinary) []model.RuleResult {
	if len(e.rules) == 0 {
		e.logger.Warn("no rules loaded, call LoadPreset() first")
		return nil
	}

	results := make([]model.RuleResult, 0, len(e.rules))

	for _, rule := range e.rules {
		applicability := rule.Applicability()
		if !info.Architecture.Matches(applicability.Arch) {
			continue
		}

		var result model.RuleResult
		switch r := rule.(type) {
		case model.ELFRule:
			if info.Format != model.FormatELF {
				continue
			}
			result = r.Execute(info.ELFFile, info)
		default:
			continue
		}

		result.RuleID = rule.ID()
		result.Name = rule.Name()

		if result.State == model.CheckStateFailed {
			result.Suggestion = buildSuggestion(info.Build.Toolchain, applicability)
		}

		results = append(results, result)
	}

	return results
}
