package rule

import (
	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/toolchain"
)

// Status indicates whether a rule passed, failed, or was skipped.
type Status int

const (
	StatusPassed Status = iota
	StatusFailed
	StatusSkipped
)

func (s Status) String() string {
	switch s {
	case StatusPassed:
		return "passed"
	case StatusFailed:
		return "failed"
	case StatusSkipped:
		return "skipped"
	default:
		return "unknown"
	}
}

// Result is the outcome of executing a rule.
type Result struct {
	Status  Status
	Message string
}

// Finding is a Result with rule metadata attached.
type Finding struct {
	Result
	RuleID string
	Name   string
}

// CompilerRequirement specifies version and flag requirements for a compiler.
type CompilerRequirement struct {
	MinVersion     toolchain.Version
	DefaultVersion toolchain.Version
	Flag           string
}

// Applicability defines which platforms and compilers a rule applies to.
type Applicability struct {
	Platform  binary.Platform
	Compilers map[toolchain.Compiler]CompilerRequirement
}

// Rule is a check that can be executed against a binary.
type Rule interface {
	ID() string
	Name() string
	Description() string
	Applicability() Applicability
}

// ELFRule is a Rule that operates on ELF binaries.
type ELFRule interface {
	Rule
	Execute(bin *binary.ELFBinary) Result
}
