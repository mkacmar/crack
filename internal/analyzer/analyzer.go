package analyzer

import (
	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/rule"
)

// FileResult contains analysis results for a single file (or arch slice for fat binaries).
type FileResult struct {
	Path     string
	Info     binary.Info
	SHA256   string
	Findings []rule.Finding
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
