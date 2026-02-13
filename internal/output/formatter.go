package output

import (
	"io"
)

type Formatter interface {
	Format(report *DecoratedReport, w io.Writer) error
}
