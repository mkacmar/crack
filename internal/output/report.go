package output

import (
	"go.kacmar.sk/crack/internal/analyzer"
	"go.kacmar.sk/crack/internal/suggestions"
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
