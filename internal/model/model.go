package model

import (
	"debug/elf"
	"fmt"
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

type Compiler int

const (
	CompilerUnknown Compiler = iota
	CompilerGCC
	CompilerClang
)

func (c Compiler) String() string {
	switch c {
	case CompilerGCC:
		return "gcc"
	case CompilerClang:
		return "clang"
	default:
		return "unknown"
	}
}

type Toolchain struct {
	Compiler Compiler
	Version  Version
}

type CompilerRequirement struct {
	Compiler       Compiler
	MinVersion     Version
	DefaultVersion Version
	Flag           string
}

type FeatureAvailability struct {
	Requirements []CompilerRequirement
}

func (f FeatureAvailability) GetRequirement(compiler Compiler) *CompilerRequirement {
	for i := range f.Requirements {
		if f.Requirements[i].Compiler == compiler {
			return &f.Requirements[i]
		}
	}
	return nil
}

type BinaryFormat int

const (
	FormatUnknown BinaryFormat = iota
	FormatELF
)

func (f BinaryFormat) String() string {
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

type ParsedBinary struct {
	Path         string
	Format       BinaryFormat
	Architecture Architecture
	Bits         int
	ELFFile      *elf.File
	Build        CompilerInfo
	LibC         LibC
}

type BinaryParser interface {
	CanParse(path string) (bool, error)
	Parse(path string) (*ParsedBinary, error)
}

type CheckState int

const (
	CheckStatePassed CheckState = iota
	CheckStateFailed
	CheckStateSkipped
)

func (cs CheckState) String() string {
	switch cs {
	case CheckStatePassed:
		return "passed"
	case CheckStateFailed:
		return "failed"
	case CheckStateSkipped:
		return "skipped"
	default:
		return "unknown"
	}
}

type RuleResult struct {
	RuleID     string
	Name       string
	State      CheckState
	Message    string
	Suggestion string
}

type FlagType int

const (
	FlagTypeCompile FlagType = iota
	FlagTypeLink
	FlagTypeBoth
)

type Rule interface {
	ID() string
	Name() string
	Format() BinaryFormat
	FlagType() FlagType
	TargetArch() Architecture
	HasPerfImpact() bool
	Feature() FeatureAvailability
	Execute(f *elf.File, info *ParsedBinary) RuleResult
}

type FileScanResult struct {
	Path      string
	Format    BinaryFormat
	Toolchain Toolchain
	SHA256    string
	Results   []RuleResult
	Error     error
	Skipped   bool
}

func (sr *FileScanResult) PassedChecks() int {
	count := 0
	for _, result := range sr.Results {
		if result.State == CheckStatePassed {
			count++
		}
	}
	return count
}

func (sr *FileScanResult) FailedChecks() int {
	count := 0
	for _, result := range sr.Results {
		if result.State == CheckStateFailed {
			count++
		}
	}
	return count
}

type ScanResults struct {
	Results []FileScanResult
}
