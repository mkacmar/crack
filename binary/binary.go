package binary

import (
	"fmt"

	"github.com/mkacmar/crack/toolchain"
)

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

// BitWidth represents the binary's word size.
type BitWidth int

const (
	BitsUnknown BitWidth = 0
	Bits32      BitWidth = 32
	Bits64      BitWidth = 64
)

func (b BitWidth) String() string {
	if b == BitsUnknown {
		return "unknown"
	}
	return fmt.Sprintf("%d-bit", b)
}

// LibC identifies the C library the binary is linked against.
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

// Binary holds common metadata for any executable format.
type Binary struct {
	Format       Format
	Architecture Architecture
	Bits         BitWidth
	Build        toolchain.BuildInfo
	LibC         LibC
}
