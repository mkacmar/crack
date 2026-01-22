package elf

import (
	"debug/elf"
	"fmt"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

// WXorXRule checks for W^X (Write XOR Execute) policy
// GNU ld: https://sourceware.org/binutils/docs/ld/Options.html (-z noexecstack)
type WXorXRule struct{}

func (r WXorXRule) ID() string   { return "wxorx" }
func (r WXorXRule) Name() string { return "W^X (Write XOR Execute)" }

func (r WXorXRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Arch: binary.ArchAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 3, Minor: 0}, DefaultVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-z noexecstack"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 0}, DefaultVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-z noexecstack"},
		},
	}
}

func (r WXorXRule) Execute(f *elf.File, info *binary.Parsed) rule.ExecuteResult {
	for _, prog := range f.Progs {
		// Check PT_LOAD segments for W+X
		if prog.Type == elf.PT_LOAD {
			if (prog.Flags&elf.PF_W) != 0 && (prog.Flags&elf.PF_X) != 0 {
				return rule.ExecuteResult{
					Status: rule.StatusFailed,
					Message: fmt.Sprintf("W^X violation: segment at offset 0x%x is both writable and executable", prog.Off),
				}
			}
		}
		// Check PT_GNU_STACK for executable stack
		if prog.Type == elf.PT_GNU_STACK && (prog.Flags&elf.PF_X) != 0 {
			return rule.ExecuteResult{
				Status: rule.StatusFailed,
				Message: "W^X violation: executable stack",
			}
		}
	}

	return rule.ExecuteResult{
		Status: rule.StatusPassed,
		Message: "All memory segments follow W^X policy (no segment is both writable and executable)",
	}
}
