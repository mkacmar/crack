package output

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/mkacmar/crack/internal/analyzer"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

type AggregatedReport struct {
	Upgrades  map[toolchain.Compiler]map[string]map[string]bool // compiler -> version -> paths
	Flags     map[string]map[string]bool                        // flag -> paths
	PassedAll []string
}

func NewAggregatedReport() *AggregatedReport {
	return &AggregatedReport{
		Upgrades: make(map[toolchain.Compiler]map[string]map[string]bool),
		Flags:    make(map[string]map[string]bool),
	}
}

func AggregateFindings(report *analyzer.Results) *AggregatedReport {
	agg := NewAggregatedReport()

	for _, res := range report.Results {
		processResult(agg, res)
	}

	slices.Sort(agg.PassedAll)
	return agg
}

func processResult(agg *AggregatedReport, result analyzer.Result) {
	if result.Error != nil {
		return
	}

	var failedResults []rule.ProcessedResult
	allPassed := true

	for _, res := range result.Results {
		if res.Status == rule.StatusFailed {
			failedResults = append(failedResults, res)
			allPassed = false
		}
	}

	if allPassed {
		agg.PassedAll = append(agg.PassedAll, result.Path)
		return
	}

	for _, res := range failedResults {
		processFailedResult(agg, res, result.Path, result.Toolchain.Compiler)
	}
}

func processFailedResult(agg *AggregatedReport, res rule.ProcessedResult, path string, detectedCompiler toolchain.Compiler) {
	r := rule.Get(res.RuleID)
	if r == nil {
		return
	}

	applicability := r.Applicability()

	for compiler, req := range applicability.Compilers {
		processRequirement(agg, compiler, req, path, detectedCompiler)
	}
}

func processRequirement(agg *AggregatedReport, compiler toolchain.Compiler, req rule.CompilerRequirement, path string, detectedCompiler toolchain.Compiler) {
	if detectedCompiler != toolchain.CompilerUnknown && compiler != detectedCompiler {
		return
	}

	ver := req.DefaultVersion
	if ver == (toolchain.Version{}) {
		ver = req.MinVersion
	}
	if ver != (toolchain.Version{}) {
		addToUpgrades(agg, compiler, ver.String(), path)
	}

	if req.Flag == "" {
		return
	}

	addToFlags(agg.Flags, req.Flag, path)
}

func addToUpgrades(agg *AggregatedReport, compiler toolchain.Compiler, version, path string) {
	if agg.Upgrades[compiler] == nil {
		agg.Upgrades[compiler] = make(map[string]map[string]bool)
	}
	if agg.Upgrades[compiler][version] == nil {
		agg.Upgrades[compiler][version] = make(map[string]bool)
	}
	agg.Upgrades[compiler][version][path] = true
}

func addToFlags(flags map[string]map[string]bool, flag, path string) {
	if flags[flag] == nil {
		flags[flag] = make(map[string]bool)
	}
	flags[flag][path] = true
}

func mapKeys(m map[string]bool) []string {
	return slices.Sorted(maps.Keys(m))
}

func FormatAggregated(agg *AggregatedReport) string {
	var sb strings.Builder

	gccUpgrades := agg.Upgrades[toolchain.CompilerGCC]
	clangUpgrades := agg.Upgrades[toolchain.CompilerClang]

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
				formatCompilerUpgrades(&sb, toolchain.CompilerGCC, gccVer, gccBinaries)
				formatCompilerUpgrades(&sb, toolchain.CompilerClang, clangVer, clangBinaries)
			}
		} else {
			if gccVer != "" {
				formatCompilerUpgrades(&sb, toolchain.CompilerGCC, gccVer, mapKeys(gccUpgrades[gccVer]))
			}
			if clangVer != "" {
				formatCompilerUpgrades(&sb, toolchain.CompilerClang, clangVer, mapKeys(clangUpgrades[clangVer]))
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

func formatCompilerUpgrades(sb *strings.Builder, compiler toolchain.Compiler, ver string, binaries []string) {
	if ver == "" || len(binaries) == 0 {
		return
	}
	sb.WriteString(fmt.Sprintf("  %s %s+:\n", compiler.String(), ver))
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
