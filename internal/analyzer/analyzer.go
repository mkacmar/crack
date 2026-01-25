package analyzer

import (
	"context"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

type Result struct {
	Path      string
	Format    binary.Format
	Toolchain toolchain.Toolchain
	SHA256    string
	Results   []rule.ProcessedResult
	Error     error
	Skipped   bool
}

func (r *Result) PassedRules() int {
	count := 0
	for _, res := range r.Results {
		if res.Status == rule.StatusPassed {
			count++
		}
	}
	return count
}

func (r *Result) FailedRules() int {
	count := 0
	for _, res := range r.Results {
		if res.Status == rule.StatusFailed {
			count++
		}
	}
	return count
}

type Results struct {
	Results []Result
}

type FileAnalyzer interface {
	Analyze(ctx context.Context, path string) Result
}
