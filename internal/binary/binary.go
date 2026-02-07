package binary

import (
	"errors"

	"github.com/mkacmar/crack/internal/toolchain"
)

var ErrUnsupportedFormat = errors.New("unsupported binary format")

type Binary struct {
	Path         string
	Format       Format
	Architecture Architecture
	Bits         BitWidth
	Build        toolchain.CompilerInfo
	LibC         toolchain.LibC
}
