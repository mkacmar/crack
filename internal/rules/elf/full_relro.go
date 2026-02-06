package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const FullRELRORuleID = "full-relro"

// FullRELRORule checks for full RELRO protection
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type FullRELRORule struct{}

func (r FullRELRORule) ID() string   { return FullRELRORuleID }
func (r FullRELRORule) Name() string { return "Full RELRO" }
func (r FullRELRORule) Description() string {
	return "Checks for full RELRO (Relocation Read-Only) protection. Full RELRO makes the Global Offset Table (GOT) read-only after initialization, preventing GOT overwrite attacks that redirect function calls to malicious code."
}

func (r FullRELRORule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, DefaultVersion: toolchain.Version{Major: 6, Minor: 1}, Flag: "-Wl,-z,relro,-z,now"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, DefaultVersion: toolchain.Version{Major: 4, Minor: 0}, Flag: "-Wl,-z,relro,-z,now"},
			toolchain.CompilerRustc: {MinVersion: toolchain.Version{Major: 1, Minor: 21}, DefaultVersion: toolchain.Version{Major: 1, Minor: 21}, Flag: "-C link-arg=-z -C link-arg=relro -C link-arg=-z -C link-arg=now"},
		},
	}
}

func (r FullRELRORule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	if bin.File.Type != elf.ET_EXEC && bin.File.Type != elf.ET_DYN {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	hasRELRO := false
	for _, prog := range bin.File.Progs {
		if prog.Type == elf.PT_GNU_RELRO {
			hasRELRO = true
			break
		}
	}

	if !hasRELRO {
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: "Full RELRO not enabled, no RELRO segment",
		}
	}

	// Check for BIND_NOW flag which makes RELRO "full" by disabling lazy binding.
	// This can be indicated by DT_BIND_NOW, DT_FLAGS with DF_BIND_NOW, or DT_FLAGS_1 with DF_1_NOW.
	if HasDynTag(bin.File, elf.DT_BIND_NOW) ||
		HasDynFlag(bin.File, elf.DT_FLAGS, uint64(elf.DF_BIND_NOW)) ||
		HasDynFlag(bin.File, elf.DT_FLAGS_1, uint64(elf.DF_1_NOW)) {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "Full RELRO enabled",
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "Full RELRO not enabled, partial RELRO only",
	}
}
