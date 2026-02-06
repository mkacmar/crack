package binary

import (
	"debug/elf"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mkacmar/crack/internal/toolchain"
)

var ErrUnsupportedFormat = errors.New("unsupported binary format")

type Format int

const (
	FormatUnknown Format = iota
	FormatELF
)

func (f Format) String() string {
	switch f {
	case FormatELF:
		return "ELF"
	default:
		return "Unknown"
	}
}

type Architecture uint32

const (
	ArchUnknown Architecture = 0
	ArchX86     Architecture = 1 << 0
	ArchX86_64  Architecture = 1 << 1
	ArchARM     Architecture = 1 << 2
	ArchARM64   Architecture = 1 << 3
	ArchRISCV   Architecture = 1 << 4
	ArchPPC64   Architecture = 1 << 5
	ArchMIPS    Architecture = 1 << 6
	ArchS390X   Architecture = 1 << 7

	ArchAllX86 = ArchX86 | ArchX86_64
	ArchAllARM = ArchARM | ArchARM64
	ArchAll    = ArchX86 | ArchX86_64 | ArchARM | ArchARM64 | ArchRISCV | ArchPPC64 | ArchMIPS | ArchS390X
)

func (a Architecture) String() string {
	switch a {
	case ArchX86:
		return "x86"
	case ArchX86_64:
		return "amd64"
	case ArchARM:
		return "arm"
	case ArchARM64:
		return "arm64"
	case ArchRISCV:
		return "riscv"
	case ArchPPC64:
		return "ppc64"
	case ArchMIPS:
		return "mips"
	case ArchS390X:
		return "s390x"
	default:
		return "unknown"
	}
}

func ParseArchitecture(s string) (Architecture, bool) {
	switch s {
	case "x86":
		return ArchX86, true
	case "amd64":
		return ArchX86_64, true
	case "arm":
		return ArchARM, true
	case "arm64":
		return ArchARM64, true
	case "riscv":
		return ArchRISCV, true
	case "ppc64":
		return ArchPPC64, true
	case "mips":
		return ArchMIPS, true
	case "s390x":
		return ArchS390X, true
	default:
		return ArchUnknown, false
	}
}

func ValidArchitectureNames() []string {
	return []string{"x86", "amd64", "arm", "arm64", "riscv", "ppc64", "mips", "s390x"}
}

func (a Architecture) Matches(target Architecture) bool {
	return a&target != 0
}

type ISA struct {
	Major int
	Minor int
}

// https://developer.arm.com/documentation/ddi0487/latest
var (
	ARM64v8_3 = ISA{Major: 8, Minor: 3}
	ARM64v8_5 = ISA{Major: 8, Minor: 5}
)

// https://gitlab.com/x86-psABIs/x86-64-ABI
var (
	X86_64v1 = ISA{Major: 1}
	X86_64v2 = ISA{Major: 2}
	X86_64v3 = ISA{Major: 3}
	X86_64v4 = ISA{Major: 4}
)

func (i ISA) String() string {
	if i.Minor > 0 {
		return fmt.Sprintf("v%d.%d", i.Major, i.Minor)
	}
	return fmt.Sprintf("v%d", i.Major)
}

func ParseISA(s string) (ISA, error) {
	s = strings.TrimPrefix(s, "v")
	parts := strings.Split(s, ".")
	if len(parts) == 0 || len(parts) > 2 {
		return ISA{}, fmt.Errorf("invalid ISA format: %s", s)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return ISA{}, fmt.Errorf("invalid ISA major version: %s", parts[0])
	}

	var minor int
	if len(parts) == 2 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil {
			return ISA{}, fmt.Errorf("invalid ISA minor version: %s", parts[1])
		}
	}

	return ISA{Major: major, Minor: minor}, nil
}

func (i ISA) IsAtLeast(required ISA) bool {
	if i.Major != required.Major {
		return i.Major > required.Major
	}
	return i.Minor >= required.Minor
}

type Platform struct {
	Architecture Architecture
	MinISA       ISA
}

func (p Platform) String() string {
	if p.MinISA == (ISA{}) {
		return p.Architecture.String()
	}
	return fmt.Sprintf("%s %s", p.Architecture.String(), p.MinISA.String())
}

var (
	PlatformAll    = Platform{Architecture: ArchAll}
	PlatformAllX86 = Platform{Architecture: ArchAllX86}
	PlatformAllARM = Platform{Architecture: ArchAllARM}

	PlatformARM64v8_3 = Platform{Architecture: ArchARM64, MinISA: ARM64v8_3}
	PlatformARM64v8_5 = Platform{Architecture: ArchARM64, MinISA: ARM64v8_5}
)

type BitWidth uint8

const (
	Bits32 BitWidth = 32
	Bits64 BitWidth = 64
)

func (b BitWidth) String() string {
	return fmt.Sprintf("%d-bit", b)
}

type Binary struct {
	Path         string
	Format       Format
	Architecture Architecture
	Bits         BitWidth
	Build        toolchain.CompilerInfo
	LibC         toolchain.LibC
}

type ELFBinary struct {
	Binary
	File       *elf.File
	Symbols    []elf.Symbol
	DynSymbols []elf.Symbol
}
