// Package binary provides types for parsing and representing executable binaries.
package binary

import (
	"errors"
	"sort"
	"strings"

	"go.kacmar.sk/crack/toolchain"
)

// ErrUnsupportedFormat is returned when the file is not a supported binary format.
var ErrUnsupportedFormat = errors.New("unsupported binary format")

// Format identifies the executable format.
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

// LibC identifies the C library the binary is linked against.
type LibC uint32

const (
	LibCUnknown LibC = 0
	LibCNone    LibC = 1 << 0
	LibCGlibc   LibC = 1 << 1
	LibCMusl    LibC = 1 << 2

	LibCAll = LibCNone | LibCGlibc | LibCMusl
)

var libcNames = map[LibC]string{
	LibCNone:  "none",
	LibCGlibc: "glibc",
	LibCMusl:  "musl",
}

func (l LibC) String() string {
	if name, ok := libcNames[l]; ok {
		return name
	}
	var names []string
	for libc, name := range libcNames {
		if l&libc != 0 {
			names = append(names, name)
		}
	}
	if len(names) > 0 {
		sort.Strings(names)
		return strings.Join(names, ", ")
	}
	return "unknown"
}

// Matches reports whether l has any overlap with target.
func (l LibC) Matches(target LibC) bool {
	return l&target != 0
}

// Profile holds the detected attributes of a binary.
type Profile struct {
	Architecture Architecture
	Toolchain    toolchain.Toolchain
	LibC         LibC
}

// Identity contains the unique fingerprints of a binary artifact.
type Identity struct {
	BuildID string
	SHA256  string
}
