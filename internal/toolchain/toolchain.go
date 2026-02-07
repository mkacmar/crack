package toolchain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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

type Compiler string

const (
	CompilerUnknown Compiler = ""
	CompilerGCC     Compiler = "gcc"
	CompilerClang   Compiler = "clang"
	CompilerRustc   Compiler = "rustc"
)

func (c Compiler) String() string {
	if c == "" {
		return "unknown"
	}
	return string(c)
}

func ValidCompilerNames() []string {
	return []string{string(CompilerGCC), string(CompilerClang), string(CompilerRustc)}
}

func ParseCompiler(s string) (Compiler, bool) {
	switch Compiler(s) {
	case CompilerGCC, CompilerClang, CompilerRustc:
		return Compiler(s), true
	default:
		return CompilerUnknown, false
	}
}

type Toolchain struct {
	Compiler Compiler
	Version  Version
}

func ParseToolchain(info string) Toolchain {
	if info == "" {
		return Toolchain{Compiler: CompilerUnknown}
	}

	lower := strings.ToLower(info)

	var compiler Compiler
	if strings.Contains(lower, "rustc") {
		compiler = CompilerRustc
	} else if strings.Contains(lower, "gcc") || strings.Contains(lower, "gnu c") || strings.Contains(lower, "gnu gimple") {
		compiler = CompilerGCC
	} else if strings.Contains(lower, "clang") {
		compiler = CompilerClang
	} else {
		return Toolchain{Compiler: CompilerUnknown}
	}

	parts := strings.Fields(info)

	for i, part := range parts {
		if strings.ToLower(part) == "version" && i+1 < len(parts) {
			version := strings.TrimRight(parts[i+1], "(),;")
			if v, err := ParseVersion(version); err == nil {
				return Toolchain{Compiler: compiler, Version: v}
			}
		}
	}

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

type CompilerInfo struct {
	BuildID   string
	Toolchain Toolchain
}

type LibC int

const (
	LibCUnknown LibC = iota
	LibCGlibc
	LibCMusl
)

func (l LibC) String() string {
	switch l {
	case LibCGlibc:
		return "glibc"
	case LibCMusl:
		return "musl"
	default:
		return "unknown"
	}
}
