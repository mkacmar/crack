package analyzer

import (
	"context"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// FindingWithSuggestion extends rule.Finding with a fix suggestion.
type FindingWithSuggestion struct {
	rule.Finding
	Suggestion string
}

type FileResult struct {
	Path     string
	Format   binary.Format
	Build    toolchain.BuildInfo
	SHA256   string
	Findings []FindingWithSuggestion
	Error    error
	Skipped  bool
}

func (r *FileResult) PassedRules() int {
	count := 0
	for _, res := range r.Findings {
		if res.Status == rule.StatusPassed {
			count++
		}
	}
	return count
}

func (r *FileResult) FailedRules() int {
	count := 0
	for _, res := range r.Findings {
		if res.Status == rule.StatusFailed {
			count++
		}
	}
	return count
}

type Report struct {
	Results []FileResult
}

type FileAnalyzer interface {
	Analyze(ctx context.Context, path string) FileResult
}
