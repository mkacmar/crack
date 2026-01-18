package rules

import (
	"log/slog"
	"slices"

	"github.com/mkacmar/crack/internal/model"
	"github.com/mkacmar/crack/internal/profile"
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

func (e *Engine) LoadProfile(p profile.Profile) {
	e.rules = make([]model.Rule, 0)
	loaded := make(map[string]bool)

	for _, rule := range elf.AllRules {
		if slices.Contains(p.Rules, rule.ID()) {
			e.rules = append(e.rules, rule)
			loaded[rule.ID()] = true
		}
	}

	for _, id := range p.Rules {
		if !loaded[id] {
			e.logger.Warn("unknown rule ID in profile, skipping", slog.String("rule_id", id))
		}
	}
}

func (e *Engine) ExecuteRules(info *model.ParsedBinary) []model.RuleResult {
	if len(e.rules) == 0 {
		e.logger.Warn("no rules loaded, call LoadProfile() first")
		return nil
	}

	results := make([]model.RuleResult, 0, len(e.rules))

	for _, rule := range e.rules {
		if rule.Format() != info.Format {
			continue
		}

		targetArch := rule.TargetArch()
		if targetArch != 0 && !info.Architecture.Matches(targetArch) {
			continue
		}

		result := rule.Execute(info.ELFFile, info)
		result.RuleID = rule.ID()
		result.Name = rule.Name()

		if result.State == model.CheckStateFailed {
			result.Suggestion = buildSuggestion(info.Build.Toolchain, rule.Feature())
		}

		results = append(results, result)
	}

	return results
}
