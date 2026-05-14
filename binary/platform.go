package binary

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Architecture identifies CPU architecture as a bitmask, allowing combinations.
type Architecture uint32

const (
	ArchUnknown Architecture = 0
	ArchX86     Architecture = 1 << 0
	ArchAMD64   Architecture = 1 << 1
	ArchARM     Architecture = 1 << 2
	ArchARM64   Architecture = 1 << 3
	ArchRISCV   Architecture = 1 << 4
	ArchPPC64   Architecture = 1 << 5
	ArchMIPS    Architecture = 1 << 6
	ArchS390X   Architecture = 1 << 7

	ArchAllX86 = ArchX86 | ArchAMD64
	ArchAllARM = ArchARM | ArchARM64
	ArchAll    = ArchX86 | ArchAMD64 | ArchARM | ArchARM64 | ArchRISCV | ArchPPC64 | ArchMIPS | ArchS390X
)

var architectureNames = map[Architecture]string{
	ArchX86:   "x86",
	ArchAMD64: "amd64",
	ArchARM:   "arm",
	ArchARM64: "arm64",
	ArchRISCV: "riscv",
	ArchPPC64: "ppc64",
	ArchMIPS:  "mips",
	ArchS390X: "s390x",
}

func (a Architecture) String() string {
	if name, ok := architectureNames[a]; ok {
		return name
	}
	var names []string
	for arch, name := range architectureNames {
		if a&arch != 0 {
			names = append(names, name)
		}
	}
	if len(names) > 0 {
		sort.Strings(names)
		return strings.Join(names, ", ")
	}
	return "unknown"
}

// Matches reports whether a has any overlap with target.
func (a Architecture) Matches(target Architecture) bool {
	return a&target != 0
}

// ParseArchitecture converts a string to Architecture.
// Returns false if unknown.
func ParseArchitecture(s string) (Architecture, bool) {
	for arch, name := range architectureNames {
		if name == s {
			return arch, true
		}
	}
	return ArchUnknown, false
}

// ISA represents an instruction set architecture version
type ISA struct {
	Major int
	Minor int
}

var (
	ARM64v83 = ISA{Major: 8, Minor: 3}
	ARM64v85 = ISA{Major: 8, Minor: 5}
)

func (i ISA) String() string {
	if i.Minor > 0 {
		return fmt.Sprintf("v%d.%d", i.Major, i.Minor)
	}
	return fmt.Sprintf("v%d", i.Major)
}

var (
	ErrInvalidISAFormat = errors.New("invalid ISA format")
	ErrInvalidISAMajor  = errors.New("invalid ISA major version component")
	ErrInvalidISAMinor  = errors.New("invalid ISA minor version component")
)

// ParseISA parses a version string into an ISA.
func ParseISA(s string) (ISA, error) {
	s = strings.TrimPrefix(s, "v")
	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return ISA{}, fmt.Errorf("%w: %s", ErrInvalidISAFormat, s)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return ISA{}, fmt.Errorf("%w: %w", ErrInvalidISAMajor, err)
	}

	var minor int
	if len(parts) == 2 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil {
			return ISA{}, fmt.Errorf("%w: %w", ErrInvalidISAMinor, err)
		}
	}

	return ISA{Major: major, Minor: minor}, nil
}

// IsAtLeast reports whether ISA is at least the required version.
func (i ISA) IsAtLeast(required ISA) bool {
	if i.Major != required.Major {
		return i.Major > required.Major
	}
	return i.Minor >= required.Minor
}

// Platform combines architecture with optional minimum ISA requirement.
type Platform struct {
	Architecture Architecture
	MinISA       ISA
}

var (
	PlatformAll    = Platform{Architecture: ArchAll}
	PlatformAllX86 = Platform{Architecture: ArchAllX86}

	PlatformARM64v83 = Platform{Architecture: ArchARM64, MinISA: ARM64v83}
	PlatformARM64v85 = Platform{Architecture: ArchARM64, MinISA: ARM64v85}
)

func (p Platform) String() string {
	if p.MinISA == (ISA{}) {
		return p.Architecture.String()
	}
	return fmt.Sprintf("%s %s", p.Architecture.String(), p.MinISA.String())
}
