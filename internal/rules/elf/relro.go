package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const RELRORuleID = "relro"

// RELRORule checks for partial RELRO
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type RELRORule struct{}

func (r RELRORule) ID() string   { return RELRORuleID }
func (r RELRORule) Name() string { return "Partial RELRO" }
func (r RELRORule) Description() string {
	return "Checks for partial RELRO (Relocation Read-Only) protection. Partial RELRO reorders ELF sections to protect internal data structures and makes some segments read-only, but leaves the GOT writable for lazy binding."
}

func (r RELRORule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, DefaultVersion: toolchain.Version{Major: 6, Minor: 1}, Flag: "-Wl,-z,relro"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, DefaultVersion: toolchain.Version{Major: 3, Minor: 9}, Flag: "-Wl,-z,relro"},
		},
	}
}

func (r RELRORule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	hasRELRO := false
	for _, prog := range bin.File.Progs {
		if prog.Type == elf.PT_GNU_RELRO {
			hasRELRO = true
			break
		}
	}

	if hasRELRO {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "RELRO enabled",
		}
	}
	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "RELRO not enabled",
	}
}
