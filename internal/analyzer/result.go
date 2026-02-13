package analyzer

import (
	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
)

// AnalysisResult contains findings and binary metadata from analysis.
type AnalysisResult struct {
	Info     binary.Info
	Findings []rule.Finding
}
