package elf

import (
	"debug/elf"
	"fmt"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// StackLimitRuleID is the rule ID for stack limit.
const StackLimitRuleID = "stack-limit"

// StackLimitRule checks for explicit stack size limit.
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type StackLimitRule struct{}

func (r StackLimitRule) ID() string   { return StackLimitRuleID }
func (r StackLimitRule) Name() string { return "Explicit Stack Size Limit" }
func (r StackLimitRule) Description() string {
	return "Checks if an explicit stack size limit is set. Defining a maximum stack size helps prevent stack exhaustion attacks."
}

func (r StackLimitRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, Flag: "-Wl,-z,stack-size=<bytes>"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, Flag: "-Wl,-z,stack-size=<bytes>"},
		},
	}
}

func (r StackLimitRule) Execute(bin *binary.ELFBinary) rule.Result {
	if bin.Type != elf.ET_EXEC && bin.Type != elf.ET_DYN {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	var stackSize uint64

	for _, prog := range bin.Progs {
		if prog.Type == elf.PT_GNU_STACK {
			stackSize = prog.Memsz
			break
		}
	}

	if stackSize > 0 {
		return rule.Result{
			Status:  rule.StatusPassed,
			Message: fmt.Sprintf("Explicit stack limit set (%d bytes)", stackSize),
		}
	}

	return rule.Result{
		Status:  rule.StatusFailed,
		Message: "No explicit stack limit",
	}
}
