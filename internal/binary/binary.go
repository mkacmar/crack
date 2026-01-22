package binary

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/toolchain"
)

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
	if target == ArchAll {
		return true
	}
	return a&target != 0
}

func (a Architecture) IsX86() bool {
	return a&ArchAllX86 != 0
}

func (a Architecture) IsARM() bool {
	return a&ArchAllARM != 0
}

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
	CanParse(path string) (bool, error)
	Parse(path string) (*Parsed, error)
}
