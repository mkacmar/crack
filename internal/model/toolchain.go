package model

import (
	"errors"
	"fmt"
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

	var compiler Compiler
	if strings.Contains(lower, "gcc") || strings.Contains(lower, "gnu c") || strings.Contains(lower, "gnu gimple") {
		compiler = CompilerGCC
	} else if strings.Contains(lower, "clang") {
		compiler = CompilerClang
	} else {
		return Toolchain{Compiler: CompilerUnknown}
	}

	parts := strings.Fields(info)

	// First try: look for "version X.Y.Z" pattern (more reliable)
	for i, part := range parts {
		if strings.ToLower(part) == "version" && i+1 < len(parts) {
			version := strings.TrimRight(parts[i+1], "(),;")
			if v, err := ParseVersion(version); err == nil {
				return Toolchain{Compiler: compiler, Version: v}
			}
		}
	}

	// Fallback: find any version-like string
	for _, part := range parts {
		if strings.Count(part, ".") >= 1 {
			version := strings.TrimRight(part, "(),;")
			if v, err := ParseVersion(version); err == nil {
				return Toolchain{Compiler: compiler, Version: v}
			}
		}
	}

	return Toolchain{Compiler: compiler}
}

var (
	ErrInvalidVersionFormat = errors.New("invalid version format")
	ErrInvalidVersionMajor  = errors.New("invalid major version component")
	ErrInvalidVersionMinor  = errors.New("invalid minor version component")
)

func ParseVersion(s string) (Version, error) {
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return Version{}, ErrInvalidVersionFormat
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, fmt.Errorf("%w: %w", ErrInvalidVersionMajor, err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, fmt.Errorf("%w: %w", ErrInvalidVersionMinor, err)
	}

	var patch int
	if len(parts) >= 3 {
		patch, _ = strconv.Atoi(parts[2])
	}

	return Version{Major: major, Minor: minor, Patch: patch}, nil
}
