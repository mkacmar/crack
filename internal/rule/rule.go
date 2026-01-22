package rule

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/toolchain"
)

type CheckState int

const (
	CheckStatePassed CheckState = iota
	CheckStateFailed
	CheckStateSkipped
)

func (cs CheckState) String() string {
	switch cs {
	case CheckStatePassed:
		return "passed"
	case CheckStateFailed:
		return "failed"
	case CheckStateSkipped:
		return "skipped"
	default:
		return "unknown"
	}
}

type Result struct {
	RuleID     string
	Name       string
	State      CheckState
	Message    string
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
	Execute(f *elf.File, info *binary.Parsed) Result
}
