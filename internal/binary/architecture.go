package binary

import (
	"fmt"
	"strconv"
	"strings"
)

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
	return "unknown"
}

func ParseArchitecture(s string) (Architecture, bool) {
	for arch, name := range architectureNames {
		if name == s {
			return arch, true
		}
	}
	return ArchUnknown, false
}

func ValidArchitectureNames() []string {
	names := make([]string, 0, len(architectureNames))
	for _, name := range architectureNames {
		names = append(names, name)
	}
	return names
}

func (a Architecture) Matches(target Architecture) bool {
	return a&target != 0
}

type ISA struct {
	Major int
	Minor int
}

var (
	ARM64v8_3 = ISA{Major: 8, Minor: 3}
	ARM64v8_5 = ISA{Major: 8, Minor: 5}
)

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
