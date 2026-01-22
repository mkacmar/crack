package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const (
	DF_BIND_NOW = 0x8
	DF_1_NOW    = 0x1
)

// FullRELRORule checks for full RELRO protection
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type FullRELRORule struct{}

func (r FullRELRORule) ID() string   { return "full-relro" }
func (r FullRELRORule) Name() string { return "Full RELRO" }

func (r FullRELRORule) Applicability() rule.Applicability {
	return rule.Applicability{
		Arch: binary.ArchAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,relro,-z,now"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,relro,-z,now"},
		},
	}
}

func (r FullRELRORule) Execute(f *elf.File, info *binary.Parsed) rule.Result {
	hasRELRO := false
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_GNU_RELRO {
			hasRELRO = true
			break
		}
	}

	if !hasRELRO {
		return rule.Result{
			State:   rule.CheckStateFailed,
			Message: "Full RELRO is NOT enabled (no PT_GNU_RELRO segment)",
		}
	}

	// Check for BIND_NOW flag which makes RELRO "full" by disabling lazy binding.
	// This can be indicated by DT_BIND_NOW, DT_FLAGS with DF_BIND_NOW, or DT_FLAGS_1 with DF_1_NOW.
	if HasDynTag(f, elf.DT_BIND_NOW) ||
		HasDynFlag(f, elf.DT_FLAGS, DF_BIND_NOW) ||
		HasDynFlag(f, elf.DT_FLAGS_1, DF_1_NOW) {
		return rule.Result{
			State:   rule.CheckStatePassed,
			Message: "Full RELRO is enabled (GOT read-only, lazy binding disabled)",
		}
	}

	return rule.Result{
		State:   rule.CheckStateFailed,
		Message: "Full RELRO is NOT enabled (GOT may be writable)",
	}
}
