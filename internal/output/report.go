package output

import (
	"github.com/mkacmar/crack/internal/analyzer"
	"github.com/mkacmar/crack/internal/suggestions"
)

// DecoratedFileResult extends FileResult with suggestion-enriched findings.
type DecoratedFileResult struct {
	analyzer.FileResult
	Findings []suggestions.DecoratedFinding
}

// DecoratedReport contains decorated analysis results for output formatting.
type DecoratedReport struct {
	Results []DecoratedFileResult
}
