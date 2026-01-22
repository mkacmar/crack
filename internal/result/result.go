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
	Results   []rule.ProcessedResult
	Error     error
	Skipped   bool
}

func (sr *FileScanResult) PassedRules() int {
	count := 0
	for _, r := range sr.Results {
		if r.Status == rule.StatusPassed {
			count++
		}
	}
	return count
}

func (sr *FileScanResult) FailedRules() int {
	count := 0
	for _, r := range sr.Results {
		if r.Status == rule.StatusFailed {
			count++
		}
	}
	return count
}

type ScanResults struct {
	Results []FileScanResult
}
