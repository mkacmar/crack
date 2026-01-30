package elf

import (
	"debug/elf"
	"fmt"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const StackLimitRuleID = "stack-limit"

// StackLimitRule checks for explicit stack size limit
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type StackLimitRule struct{}

func (r StackLimitRule) ID() string   { return StackLimitRuleID }
func (r StackLimitRule) Name() string { return "Explicit Stack Size Limit" }

func (r StackLimitRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, Flag: "-Wl,-z,stack-size=<bytes>"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, Flag: "-Wl,-z,stack-size=<bytes>"},
		},
	}
}

func (r StackLimitRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	var stackSize uint64

	for _, prog := range bin.File.Progs {
		if prog.Type == elf.PT_GNU_STACK {
			stackSize = prog.Memsz
			break
		}
	}

	if stackSize > 0 {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: fmt.Sprintf("Explicit stack limit set (%d bytes)", stackSize),
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "No explicit stack limit",
	}
}
