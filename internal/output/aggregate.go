package output

import (
	"slices"

	"github.com/mkacmar/crack/internal/suggestions"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

type PathSet map[string]bool
type CompilerUpgrades map[toolchain.Version]PathSet

type AggregatedReport struct {
	Upgrades     map[string]CompilerUpgrades
	Flags        map[string]PathSet
	PassedAll    []string
	NoApplicable []string
}

func newAggregatedReport() *AggregatedReport {
	return &AggregatedReport{
		Upgrades: make(map[string]CompilerUpgrades),
		Flags:    make(map[string]PathSet),
	}
}

func AggregateFindings(report *DecoratedReport, rules []rule.ELFRule) *AggregatedReport {
	agg := newAggregatedReport()

	rulesMap := make(map[string]rule.ELFRule, len(rules))
	for _, r := range rules {
		rulesMap[r.ID()] = r
	}

	for _, res := range report.Results {
		agg.addResult(res, rulesMap)
	}

	slices.Sort(agg.PassedAll)
	slices.Sort(agg.NoApplicable)
	return agg
}

func (agg *AggregatedReport) addResult(result DecoratedFileResult, rules map[string]rule.ELFRule) {
	if result.Error != nil {
		return
	}

	var failedFindings []suggestions.DecoratedFinding
	hasApplicable := false
	hasFailed := false

	for _, f := range result.Findings {
		switch f.Status {
		case rule.StatusFailed:
			failedFindings = append(failedFindings, f)
			hasApplicable = true
			hasFailed = true
		case rule.StatusPassed:
			hasApplicable = true
		}
	}

	if !hasApplicable {
		agg.NoApplicable = append(agg.NoApplicable, result.Path)
		return
	}

	if !hasFailed {
		agg.PassedAll = append(agg.PassedAll, result.Path)
		return
	}

	for _, f := range failedFindings {
		agg.addFinding(f, result.Path, result.Info.Build.Compiler, rules)
	}
}

func (agg *AggregatedReport) addFinding(f suggestions.DecoratedFinding, path string, detectedCompiler toolchain.Compiler, rules map[string]rule.ELFRule) {
	r := rules[f.RuleID]
	if r == nil {
		return
	}

	applicability := r.Applicability()

	for compiler, req := range applicability.Compilers {
		agg.addRequirement(compiler.String(), req, path, detectedCompiler)
	}
}

func (agg *AggregatedReport) addRequirement(compilerName string, req rule.CompilerRequirement, path string, detectedCompiler toolchain.Compiler) {
	if detectedCompiler != toolchain.Unknown && compilerName != detectedCompiler.String() {
		return
	}

	ver := req.DefaultVersion
	if ver == (toolchain.Version{}) {
		ver = req.MinVersion
	}
	if ver != (toolchain.Version{}) {
		agg.addUpgrade(compilerName, ver, path)
	}

	if req.Flag == "" {
		return
	}

	agg.addFlag(req.Flag, path)
}

func (agg *AggregatedReport) addUpgrade(compilerName string, version toolchain.Version, path string) {
	if agg.Upgrades[compilerName] == nil {
		agg.Upgrades[compilerName] = make(CompilerUpgrades)
	}
	if agg.Upgrades[compilerName][version] == nil {
		agg.Upgrades[compilerName][version] = make(PathSet)
	}
	agg.Upgrades[compilerName][version][path] = true
}

func (agg *AggregatedReport) addFlag(flag, path string) {
	if agg.Flags[flag] == nil {
		agg.Flags[flag] = make(PathSet)
	}
	agg.Flags[flag][path] = true
}
