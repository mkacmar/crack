package elf

import (
	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const ARMMTERuleID = "arm-mte"

// ARMMTERule checks for ARM Memory Tagging Extension
// ARM: https://developer.arm.com/documentation/ddi0487/latest
// Clang: https://clang.llvm.org/docs/MemTagSanitizer.html
type ARMMTERule struct{}

func (r ARMMTERule) ID() string   { return ARMMTERuleID }
func (r ARMMTERule) Name() string { return "ARM Memory Tagging Extension" }

func (r ARMMTERule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformARM64v8_5,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 11, Minor: 0}, Flag: "-march=armv8.5-a+memtag -fsanitize=memtag"},
		},
	}
}

func (r ARMMTERule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	hasMTE := false
	for _, sec := range bin.File.Sections {
		if sec.Name == ".note.memtag" {
			hasMTE = true
			break
		}
	}

	if hasMTE {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "ARM MTE (Memory Tagging Extension) is enabled",
		}
	}
	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "ARM MTE is NOT enabled (requires ARMv8.5+ hardware)",
	}
}
