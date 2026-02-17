package output

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/mkacmar/crack/toolchain"
)

func sortedPaths(m PathSet) []string {
	return slices.Sorted(maps.Keys(m))
}

func (agg *AggregatedReport) Format() string {
	var sb strings.Builder

	if len(agg.Upgrades) > 0 {
		sb.WriteString("Ensure minimum toolchain version (prerequisite). The following features require at least these compiler versions:\n\n")

		for _, compiler := range slices.Sorted(maps.Keys(agg.Upgrades)) {
			upgrades := agg.Upgrades[compiler]
			ver, ok := getHighestVersion(upgrades)
			if ok {
				formatCompilerUpgrades(&sb, compiler, ver.String(), mergePaths(upgrades))
			}
		}
	}

	allBinaries := make(PathSet)
	for _, paths := range agg.Flags {
		for p := range paths {
			allBinaries[p] = true
		}
	}
	totalWithFindings := len(allBinaries)

	if len(agg.Flags) > 0 {
		sb.WriteString("Add following flags, even with the correct toolchain, these flags must be explicitly added:\n\n")
		formatFlagSection(&sb, agg.Flags, totalWithFindings, "  ")
		sb.WriteString("\n")
	}

	if len(agg.PassedAll) > 0 {
		sb.WriteString(fmt.Sprintf("No failed checks (%d binaries): %s\n",
			len(agg.PassedAll),
			strings.Join(agg.PassedAll, ", ")))
	}
	if len(agg.NoApplicable) > 0 {
		sb.WriteString(fmt.Sprintf("No applicable checks (%d binaries): %s\n",
			len(agg.NoApplicable),
			strings.Join(agg.NoApplicable, ", ")))
	}

	if sb.Len() == 0 {
		sb.WriteString("No binaries analyzed.\n")
	}

	return sb.String()
}

func getHighestVersion(upgrades CompilerUpgrades) (toolchain.Version, bool) {
	var highest toolchain.Version
	found := false
	for v := range upgrades {
		if !found || (v.IsAtLeast(highest) && !highest.IsAtLeast(v)) {
			highest = v
			found = true
		}
	}
	return highest, found
}

func mergePaths(upgrades CompilerUpgrades) []string {
	all := make(PathSet)
	for _, paths := range upgrades {
		for p := range paths {
			all[p] = true
		}
	}
	return sortedPaths(all)
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

func formatFlagSection(sb *strings.Builder, flags map[string]PathSet, totalWithFindings int, prefix string) {
	sortedFlags := slices.Sorted(maps.Keys(flags))

	var universalFlags []string
	partialFlags := make(map[string][]string)

	for _, flag := range sortedFlags {
		binaries := sortedPaths(flags[flag])
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
