package binary

import (
	"debug/elf"
	"errors"
	"fmt"

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
		return "x86_64"
	case ArchARM:
		return "ARM"
	case ArchARM64:
		return "ARM64"
	case ArchRISCV:
		return "RISC-V"
	case ArchPPC64:
		return "PPC64"
	case ArchMIPS:
		return "MIPS"
	case ArchS390X:
		return "s390x"
	default:
		return "Unknown"
	}
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
	ARM64v8_0 = ISA{Major: 8, Minor: 0}
	ARM64v8_3 = ISA{Major: 8, Minor: 3} // PAC (Pointer Authentication)
	ARM64v8_5 = ISA{Major: 8, Minor: 5} // BTI (Branch Target Identification), MTE
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
	PlatformX86    = Platform{Architecture: ArchX86}
	PlatformX86_64 = Platform{Architecture: ArchX86_64}
	PlatformAllX86 = Platform{Architecture: ArchAllX86}
	PlatformARM    = Platform{Architecture: ArchARM}
	PlatformARM64  = Platform{Architecture: ArchARM64}
	PlatformAllARM = Platform{Architecture: ArchAllARM}

	PlatformARM64v8_3 = Platform{Architecture: ArchARM64, MinISA: ARM64v8_3} // PAC
	PlatformARM64v8_5 = Platform{Architecture: ArchARM64, MinISA: ARM64v8_5} // BTI, MTE
)

type Parsed struct {
	Path         string
	Format       Format
	Architecture Architecture
	Bits         int
	ELF          *elf.File
	Build        toolchain.CompilerInfo
	LibC         toolchain.LibC
}

type Parser interface {
	Parse(path string) (*Parsed, error)
}
