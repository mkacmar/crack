package elf

import (
	"debug/elf"
	"fmt"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

// StackLimitRule checks for explicit stack size limit
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type StackLimitRule struct{}

func (r StackLimitRule) ID() string   { return "stack-limit" }
func (r StackLimitRule) Name() string { return "Explicit Stack Size Limit" }

func (r StackLimitRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,stack-size=8388608"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,stack-size=8388608"},
		},
	}
}

func (r StackLimitRule) Execute(f *elf.File, info *binary.Parsed) rule.ExecuteResult {
	hasExplicitStackLimit := false
	stackSize := uint64(0)
	foundGNUStack := false

	for _, prog := range f.Progs {
		if prog.Type == elf.PT_GNU_STACK {
			foundGNUStack = true
			stackSize = prog.Memsz
			hasExplicitStackLimit = stackSize > 0
			break
		}
	}

	if hasExplicitStackLimit {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: fmt.Sprintf("Stack has explicit size limit: %d bytes", stackSize),
		}
	}

	msg := "No PT_GNU_STACK segment found"
	if foundGNUStack {
		msg = "Stack uses system default size (no explicit limit set)"
	}
	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: msg,
	}
}
