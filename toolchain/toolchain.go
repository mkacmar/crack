// Package toolchain provides compiler and version detection for binaries.
package toolchain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Compiler identifies a compiler toolchain.
type Compiler string

const (
	Unknown Compiler = ""
	GCC     Compiler = "gcc"
	Clang   Compiler = "clang"
	Rustc   Compiler = "rustc"
)

func (c Compiler) String() string {
	if c == "" {
		return "unknown"
	}
	return string(c)
}

// Version represents a semantic version (major.minor.patch).
type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	if v.Patch > 0 {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

// IsAtLeast reports whether v is at least the required version.
func (v Version) IsAtLeast(required Version) bool {
	if v.Major != required.Major {
		return v.Major > required.Major
	}
	if v.Minor != required.Minor {
		return v.Minor > required.Minor
	}
	return v.Patch >= required.Patch
}

var (
	ErrInvalidVersionFormat = errors.New("invalid version format")
	ErrInvalidVersionMajor  = errors.New("invalid major version component")
	ErrInvalidVersionMinor  = errors.New("invalid minor version component")
)

// ParseVersion parses a version string like "1.2.3" or "1.2".
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

// BuildInfo contains compiler metadata extracted from a binary.
type BuildInfo struct {
	BuildID  string
	Compiler Compiler
	Version  Version
}

// ELFDetector detects compiler info from ELF binaries.
type ELFDetector interface {
	Detect(comment string) (Compiler, Version)
}

// ELFCommentDetector detects compiler from ELF .comment section.
type ELFCommentDetector struct{}

func (ELFCommentDetector) Detect(comment string) (Compiler, Version) {
	if comment == "" {
		return Unknown, Version{}
	}

	lower := strings.ToLower(comment)

	var comp Compiler
	if strings.Contains(lower, "rustc") {
		comp = Rustc
	} else if strings.Contains(lower, "gcc") || strings.Contains(lower, "gnu c") || strings.Contains(lower, "gnu gimple") {
		comp = GCC
	} else if strings.Contains(lower, "clang") {
		comp = Clang
	} else {
		return Unknown, Version{}
	}

	parts := strings.Fields(comment)

	for i, part := range parts {
		if strings.ToLower(part) == "version" && i+1 < len(parts) {
			version := strings.TrimRight(parts[i+1], "(),;")
			if v, err := ParseVersion(version); err == nil {
				return comp, v
			}
		}
	}

	for _, part := range parts {
		if strings.Count(part, ".") >= 1 {
			version := strings.TrimRight(part, "(),;")
			if v, err := ParseVersion(version); err == nil {
				return comp, v
			}
		}
	}

	return comp, Version{}
}
