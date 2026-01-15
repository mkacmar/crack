package model

import (
	"strconv"
	"strings"
)

// ParseToolchain extracts toolchain info from a build string
// (e.g., from .comment section or DWARF DW_AT_producer)
func ParseToolchain(info string) Toolchain {
	if info == "" {
		return Toolchain{Compiler: CompilerUnknown}
	}

	lower := strings.ToLower(info)

	if strings.Contains(lower, "gcc") || strings.Contains(lower, "gnu c") || strings.Contains(lower, "gnu gimple") {
		return Toolchain{
			Compiler: CompilerGCC,
			Version:  parseVersion(info),
		}
	}

	if strings.Contains(lower, "clang") {
		return Toolchain{
			Compiler: CompilerClang,
			Version:  parseVersion(info),
		}
	}

	return Toolchain{Compiler: CompilerUnknown}
}

func parseVersion(info string) Version {
	parts := strings.Fields(info)

	// First try: look for "version X.Y.Z" pattern (more reliable)
	for i, part := range parts {
		if strings.ToLower(part) == "version" && i+1 < len(parts) {
			version := strings.TrimRight(parts[i+1], "(),;")
			if v, ok := tryParseVersion(version); ok {
				return v
			}
		}
	}

	// Fallback: find any version-like string
	for _, part := range parts {
		if strings.Count(part, ".") >= 1 {
			version := strings.TrimRight(part, "(),;")
			if v, ok := tryParseVersion(version); ok {
				return v
			}
		}
	}

	return Version{}
}

func tryParseVersion(s string) (Version, bool) {
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return Version{}, false
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, false
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, false
	}

	var patch int
	if len(parts) >= 3 {
		patch, _ = strconv.Atoi(parts[2])
	}

	return Version{Major: major, Minor: minor, Patch: patch}, true
}
