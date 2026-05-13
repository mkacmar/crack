package analyzer

import (
	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/rule"
)

// AnalysisResult contains findings and binary metadata from analysis.
// Identity.SHA256 is left empty at this layer. The scanner fills it after hashing the file.
type AnalysisResult struct {
	Format   binary.Format
	Identity binary.Identity
	Profile  binary.Profile
	Findings []rule.Finding
}
