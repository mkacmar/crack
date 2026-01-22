package rule

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/toolchain"
)

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

type ExecuteResult struct {
	Status  Status
	Message string
}

type ProcessedResult struct {
	ExecuteResult
	RuleID     string
	Name       string
	Suggestion string
}

type CompilerRequirement struct {
	MinVersion     toolchain.Version
	DefaultVersion toolchain.Version
	Flag           string
}

type Applicability struct {
	Arch      binary.Architecture
	Compilers map[toolchain.Compiler]CompilerRequirement
}

type Rule interface {
	ID() string
	Name() string
	Applicability() Applicability
}

type ELFRule interface {
	Rule
	Execute(f *elf.File, info *binary.Parsed) ExecuteResult
}
