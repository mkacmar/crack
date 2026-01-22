package result

import (
	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

type FileScanResult struct {
	Path      string
	Format    binary.Format
	Toolchain toolchain.Toolchain
	SHA256    string
	Results   []rule.Result
	Error     error
	Skipped   bool
}

func (sr *FileScanResult) PassedChecks() int {
	count := 0
	for _, r := range sr.Results {
		if r.State == rule.CheckStatePassed {
			count++
		}
	}
	return count
}

func (sr *FileScanResult) FailedChecks() int {
	count := 0
	for _, r := range sr.Results {
		if r.State == rule.CheckStateFailed {
			count++
		}
	}
	return count
}

type ScanResults struct {
	Results []FileScanResult
}
