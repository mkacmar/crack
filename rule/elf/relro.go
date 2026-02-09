package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// RELRORuleID is the rule ID for partial RELRO.
const RELRORuleID = "relro"

// RELRORule checks for partial RELRO.
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type RELRORule struct{}

func (r RELRORule) ID() string   { return RELRORuleID }
func (r RELRORule) Name() string { return "Partial RELRO" }
func (r RELRORule) Description() string {
	return "Checks for partial RELRO (Relocation Read-Only) protection. Partial RELRO reorders ELF sections to protect internal data structures and makes some segments read-only, but leaves the GOT writable for lazy binding."
}

func (r RELRORule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, DefaultVersion: toolchain.Version{Major: 6, Minor: 1}, Flag: "-Wl,-z,relro"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, DefaultVersion: toolchain.Version{Major: 3, Minor: 9}, Flag: "-Wl,-z,relro"},
			toolchain.Rustc: {MinVersion: toolchain.Version{Major: 1, Minor: 21}, DefaultVersion: toolchain.Version{Major: 1, Minor: 21}, Flag: "-C link-arg=-z -C link-arg=relro"},
		},
	}
}

func (r RELRORule) Execute(bin *binary.ELFBinary) rule.Result {
	if bin.Type != elf.ET_EXEC && bin.Type != elf.ET_DYN {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	hasRELRO := false
	for _, prog := range bin.Progs {
		if prog.Type == elf.PT_GNU_RELRO {
			hasRELRO = true
			break
		}
	}

	if hasRELRO {
		return rule.Result{
			Status:  rule.StatusPassed,
			Message: "RELRO enabled",
		}
	}
	return rule.Result{
		Status:  rule.StatusFailed,
		Message: "RELRO not enabled",
	}
}
