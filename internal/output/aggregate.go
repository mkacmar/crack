package output

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/mkacmar/crack/internal/suggestions"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

type AggregatedReport struct {
	Upgrades  map[string]map[string]map[string]bool // compiler name -> version -> paths
	Flags     map[string]map[string]bool            // flag -> paths
	PassedAll []string
}

func NewAggregatedReport() *AggregatedReport {
	return &AggregatedReport{
		Upgrades: make(map[string]map[string]map[string]bool),
		Flags:    make(map[string]map[string]bool),
	}
}

func AggregateFindings(report *DecoratedReport, rules []rule.ELFRule) *AggregatedReport {
	agg := NewAggregatedReport()

	rulesMap := make(map[string]rule.ELFRule, len(rules))
	for _, r := range rules {
		rulesMap[r.ID()] = r
	}

	for _, res := range report.Results {
		processResult(agg, res, rulesMap)
	}

	slices.Sort(agg.PassedAll)
	return agg
}

func processResult(agg *AggregatedReport, result DecoratedFileResult, rules map[string]rule.ELFRule) {
	if result.Error != nil {
		return
	}

	var failedFindings []suggestions.DecoratedFinding
	allPassed := true

	for _, f := range result.Findings {
		if f.Status == rule.StatusFailed {
			failedFindings = append(failedFindings, f)
			allPassed = false
		}
	}

	if allPassed {
		agg.PassedAll = append(agg.PassedAll, result.Path)
		return
	}

	for _, f := range failedFindings {
		processFailedFinding(agg, f, result.Path, result.Info.Build.Compiler, rules)
	}
}

func processFailedFinding(agg *AggregatedReport, f suggestions.DecoratedFinding, path string, detectedCompiler toolchain.Compiler, rules map[string]rule.ELFRule) {
	r := rules[f.RuleID]
	if r == nil {
		return
	}

	applicability := r.Applicability()

	for compiler, req := range applicability.Compilers {
		processRequirement(agg, compiler.String(), req, path, detectedCompiler)
	}
}

func processRequirement(agg *AggregatedReport, compilerName string, req rule.CompilerRequirement, path string, detectedCompiler toolchain.Compiler) {
	if detectedCompiler != toolchain.Unknown && compilerName != detectedCompiler.String() {
		return
	}

	ver := req.DefaultVersion
	if ver == (toolchain.Version{}) {
		ver = req.MinVersion
	}
	if ver != (toolchain.Version{}) {
		agg.addUpgrade(compilerName, ver.String(), path)
	}

	if req.Flag == "" {
		return
	}

	agg.addFlag(req.Flag, path)
}

func (agg *AggregatedReport) addUpgrade(compilerName, version, path string) {
	if agg.Upgrades[compilerName] == nil {
		agg.Upgrades[compilerName] = make(map[string]map[string]bool)
	}
	if agg.Upgrades[compilerName][version] == nil {
		agg.Upgrades[compilerName][version] = make(map[string]bool)
	}
	agg.Upgrades[compilerName][version][path] = true
}

func (agg *AggregatedReport) addFlag(flag, path string) {
	if agg.Flags[flag] == nil {
		agg.Flags[flag] = make(map[string]bool)
	}
	agg.Flags[flag][path] = true
}

func mapKeys(m map[string]bool) []string {
	return slices.Sorted(maps.Keys(m))
}

func FormatAggregated(agg *AggregatedReport) string {
	var sb strings.Builder

	gccUpgrades := agg.Upgrades[toolchain.GCC.String()]
	clangUpgrades := agg.Upgrades[toolchain.Clang.String()]

	gccVer := getHighestVersion(gccUpgrades)
	clangVer := getHighestVersion(clangUpgrades)

	if gccVer != "" || clangVer != "" {
		sb.WriteString("Ensure minimum toolchain version (prerequisite). The following features require at least these compiler versions:\n\n")

		if gccVer != "" && clangVer != "" {
			gccBinaries := mapKeys(gccUpgrades[gccVer])
			clangBinaries := mapKeys(clangUpgrades[clangVer])

			if slices.Equal(gccBinaries, clangBinaries) {
				sb.WriteString(fmt.Sprintf("  GCC %s+ or Clang %s+:\n", gccVer, clangVer))
				for _, b := range gccBinaries {
					sb.WriteString(fmt.Sprintf("    %s\n", b))
				}
				sb.WriteString("\n")
			} else {
				formatCompilerUpgrades(&sb, toolchain.GCC.String(), gccVer, gccBinaries)
				formatCompilerUpgrades(&sb, toolchain.Clang.String(), clangVer, clangBinaries)
			}
		} else {
			if gccVer != "" {
				formatCompilerUpgrades(&sb, toolchain.GCC.String(), gccVer, mapKeys(gccUpgrades[gccVer]))
			}
			if clangVer != "" {
				formatCompilerUpgrades(&sb, toolchain.Clang.String(), clangVer, mapKeys(clangUpgrades[clangVer]))
			}
		}
	}

	allBinaries := collectAllBinaries(agg.Flags)
	totalWithFindings := len(allBinaries)

	if len(agg.Flags) > 0 {
		sb.WriteString("Add following flags, even with the correct toolchain, these flags must be explicitly added:\n\n")
		formatFlagSection(&sb, agg.Flags, totalWithFindings, "  ")
		sb.WriteString("\n")
	}

	if len(agg.PassedAll) > 0 {
		sb.WriteString(fmt.Sprintf("Fully hardened (%d binaries): %s\n",
			len(agg.PassedAll),
			strings.Join(agg.PassedAll, ", ")))
	}

	if sb.Len() == 0 {
		sb.WriteString("No binaries analyzed.\n")
	}

	return sb.String()
}

func getHighestVersion(upgrades map[string]map[string]bool) string {
	if len(upgrades) == 0 {
		return ""
	}
	var highest string
	var highestVer toolchain.Version
	for v := range upgrades {
		parsed, err := toolchain.ParseVersion(v)
		if err != nil {
			continue
		}
		if highest == "" || parsed.IsAtLeast(highestVer) && parsed != highestVer {
			highest = v
			highestVer = parsed
		}
	}
	return highest
}

func collectAllBinaries(flags map[string]map[string]bool) map[string]bool {
	all := make(map[string]bool)
	for _, paths := range flags {
		for p := range paths {
			all[p] = true
		}
	}
	return all
}

func formatCompilerUpgrades(sb *strings.Builder, compilerName string, ver string, binaries []string) {
	if ver == "" || len(binaries) == 0 {
		return
	}
	sb.WriteString(fmt.Sprintf("  %s %s+:\n", compilerName, ver))
	for _, b := range binaries {
		sb.WriteString(fmt.Sprintf("    %s\n", b))
	}
	sb.WriteString("\n")
}

func formatFlagSection(sb *strings.Builder, flags map[string]map[string]bool, totalWithFindings int, prefix string) {
	sortedFlags := slices.Sorted(maps.Keys(flags))

	var universalFlags []string
	partialFlags := make(map[string][]string)

	for _, flag := range sortedFlags {
		binaries := mapKeys(flags[flag])
		if len(binaries) == totalWithFindings {
			universalFlags = append(universalFlags, flag)
		} else {
			partialFlags[flag] = binaries
		}
	}

	if len(universalFlags) > 0 {
		sb.WriteString(fmt.Sprintf("%s%s\n", prefix, strings.Join(universalFlags, " ")))
	}

	for _, flag := range sortedFlags {
		binaries, isPartial := partialFlags[flag]
		if !isPartial {
			continue
		}
		sb.WriteString(fmt.Sprintf("%s%s\n", prefix, flag))
		sb.WriteString(fmt.Sprintf("%s  Only: %s\n", prefix, strings.Join(binaries, ", ")))
	}
}
