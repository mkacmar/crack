package output

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/mkacmar/crack/internal/model"
	"github.com/mkacmar/crack/internal/rules/elf"
)

type AggregatedReport struct {
	Upgrades        map[model.Compiler]map[string]map[string]bool // compiler -> version -> paths
	CompileFlags    map[string]map[string]bool                    // flag -> paths
	LinkFlags       map[string]map[string]bool                    // flag -> paths
	PerfImpactFlags map[string]bool
	PassedAll       []string
}

func AggregateFindings(report *model.ScanResults) *AggregatedReport {
	agg := &AggregatedReport{
		Upgrades:        make(map[model.Compiler]map[string]map[string]bool),
		CompileFlags:    make(map[string]map[string]bool),
		LinkFlags:       make(map[string]map[string]bool),
		PerfImpactFlags: make(map[string]bool),
	}

	for _, result := range report.Results {
		if result.Error != nil {
			continue
		}

		var failedChecks []model.RuleResult
		allPassed := true

		for _, check := range result.Results {
			if check.State == model.CheckStateFailed {
				failedChecks = append(failedChecks, check)
				allPassed = false
			}
		}

		if allPassed {
			agg.PassedAll = append(agg.PassedAll, result.Path)
			continue
		}

		detectedCompiler := result.Toolchain.Compiler

		for _, check := range failedChecks {
			rule := elf.GetRuleByID(check.RuleID)
			if rule == nil {
				continue
			}

			feature := rule.Feature()
			hasPerfImpact := rule.HasPerfImpact()

			for _, req := range feature.Requirements {
				if detectedCompiler != model.CompilerUnknown && req.Compiler != detectedCompiler {
					continue
				}

				ver := req.DefaultVersion
				if ver.IsZero() {
					ver = req.MinVersion
				}
				if !ver.IsZero() {
					addToUpgrades(agg, req.Compiler, ver.String(), result.Path)
				}

				if req.Flag == "" {
					continue
				}

				if rule.FlagType() == model.FlagTypeCompile || rule.FlagType() == model.FlagTypeBoth {
					addToFlags(agg.CompileFlags, req.Flag, result.Path)
					if hasPerfImpact {
						agg.PerfImpactFlags[req.Flag] = true
					}
				}

				if rule.FlagType() == model.FlagTypeLink || rule.FlagType() == model.FlagTypeBoth {
					addToFlags(agg.LinkFlags, req.Flag, result.Path)
					if hasPerfImpact {
						agg.PerfImpactFlags[req.Flag] = true
					}
				}
			}
		}
	}

	sort.Strings(agg.PassedAll)
	return agg
}

func addToUpgrades(agg *AggregatedReport, compiler model.Compiler, version, path string) {
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
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func FormatAggregated(agg *AggregatedReport) string {
	var sb strings.Builder

	gccUpgrades := agg.Upgrades[model.CompilerGCC]
	clangUpgrades := agg.Upgrades[model.CompilerClang]

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
				formatCompilerUpgrades(&sb, model.CompilerGCC, gccVer, gccBinaries)
				formatCompilerUpgrades(&sb, model.CompilerClang, clangVer, clangBinaries)
			}
		} else {
			if gccVer != "" {
				formatCompilerUpgrades(&sb, model.CompilerGCC, gccVer, mapKeys(gccUpgrades[gccVer]))
			}
			if clangVer != "" {
				formatCompilerUpgrades(&sb, model.CompilerClang, clangVer, mapKeys(clangUpgrades[clangVer]))
			}
		}
	}

	allBinaries := collectAllBinaries(agg.CompileFlags, agg.LinkFlags)
	totalWithFindings := len(allBinaries)

	if len(agg.CompileFlags) > 0 || len(agg.LinkFlags) > 0 {
		sb.WriteString("Add following flags, even with the correct toolchain, these flags must be explicitly added:\n\n")
	}

	if len(agg.CompileFlags) > 0 {
		sb.WriteString("  Compiler flags (CFLAGS/CXXFLAGS):\n")
		formatFlagSection(&sb, agg.CompileFlags, totalWithFindings, "    ")
		sb.WriteString("\n")
	}

	if len(agg.LinkFlags) > 0 {
		sb.WriteString("  Linker flags (LDFLAGS):\n")
		formatFlagSection(&sb, agg.LinkFlags, totalWithFindings, "    ")
		sb.WriteString("\n")
	}

	if len(agg.PerfImpactFlags) > 0 {
		perfFlags := mapKeys(agg.PerfImpactFlags)
		sb.WriteString(fmt.Sprintf("Note: Some flags have performance impact: %s\n\n", strings.Join(perfFlags, ", ")))
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
	versions := make([]string, 0, len(upgrades))
	for v := range upgrades {
		versions = append(versions, v)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))
	return versions[0]
}

func collectAllBinaries(compileFlags, linkFlags map[string]map[string]bool) map[string]bool {
	all := make(map[string]bool)
	for _, paths := range compileFlags {
		for p := range paths {
			all[p] = true
		}
	}
	for _, paths := range linkFlags {
		for p := range paths {
			all[p] = true
		}
	}
	return all
}

func formatCompilerUpgrades(sb *strings.Builder, compiler model.Compiler, ver string, binaries []string) {
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
	sortedFlags := make([]string, 0, len(flags))
	for f := range flags {
		sortedFlags = append(sortedFlags, f)
	}
	sort.Strings(sortedFlags)

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
