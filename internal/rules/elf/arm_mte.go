package elf

import (
	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const ARMMTERuleID = "arm-mte"

// ARMMTERule checks for ARM Memory Tagging Extension
// ARM: https://developer.arm.com/documentation/ddi0487/latest
// LLVM: https://llvm.org/docs/MemTagSanitizer.html
type ARMMTERule struct{}

func (r ARMMTERule) ID() string   { return ARMMTERuleID }
func (r ARMMTERule) Name() string { return "ARM Memory Tagging Extension" }
func (r ARMMTERule) Description() string {
	return "Checks for ARM Memory Tagging Extension (MTE). MTE assigns 4-bit tags to memory regions and pointers, detecting use-after-free and buffer overflows when tags mismatch during memory access."
}

func (r ARMMTERule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformARM64v8_5,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			// gcc: MTE intrinsics supported but not easily detectable
			// rustc: nightly-only via -Z sanitizer=memtag
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 12, Minor: 0}, Flag: "-march=armv8.5-a+memtag -fsanitize=memtag"},
		},
	}
}

func (r ARMMTERule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	hasMTE := false
	for _, sec := range bin.File.Sections {
		// Currently Android NDK specific https://github.com/llvm/llvm-project/blob/main/lld/ELF/SyntheticSections.h
		if sec.Name == ".note.android.memtag" {
			hasMTE = true
			break
		}
	}

	if hasMTE {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "ARM MTE enabled",
		}
	}
	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "ARM MTE not enabled",
	}
}
