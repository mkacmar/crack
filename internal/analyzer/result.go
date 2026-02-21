package analyzer

import (
	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/rule"
)

// AnalysisResult contains findings and binary metadata from analysis.
type AnalysisResult struct {
	Info     binary.Info
	Findings []rule.Finding
}
