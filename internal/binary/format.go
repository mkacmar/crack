package binary

import "fmt"

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

type BitWidth uint8

const (
	Bits32 BitWidth = 32
	Bits64 BitWidth = 64
)

func (b BitWidth) String() string {
	return fmt.Sprintf("%d-bit", b)
}
