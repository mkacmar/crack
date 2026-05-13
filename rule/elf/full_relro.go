package elf

import (
	stdelf "debug/elf"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/binary/elf"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
)

// FullRELRORuleID is the rule ID for full RELRO.
const FullRELRORuleID = "full-relro"

// FullRELRORule checks for full RELRO protection.
//
// References:
//   - https://sourceware.org/binutils/docs/ld/Options.html
type FullRELRORule struct{}

func (r FullRELRORule) ID() string   { return FullRELRORuleID }
func (r FullRELRORule) Name() string { return "Full RELRO" }
func (r FullRELRORule) Description() string {
	return "Checks for full RELRO (Relocation Read-Only) protection. Full RELRO makes the Global Offset Table (GOT) read-only after initialization, preventing GOT overwrite attacks that redirect function calls to malicious code."
}

func (r FullRELRORule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM | binary.ArchRISCV},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, DefaultVersion: toolchain.Version{Major: 6, Minor: 1}, Flag: "-Wl,-z,relro,-z,now"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, DefaultVersion: toolchain.Version{Major: 4, Minor: 0}, Flag: "-Wl,-z,relro,-z,now"},
		},
		LibC: binary.LibCAll,
	}
}

func (r FullRELRORule) Execute(bin elf.Binary) rule.Result {
	if bin.Type() != stdelf.ET_EXEC && bin.Type() != stdelf.ET_DYN {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	hasRELRO := false
	for _, prog := range bin.Progs() {
		if prog.Type == stdelf.PT_GNU_RELRO {
			hasRELRO = true
			break
		}
	}

	if !hasRELRO {
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: "Full RELRO not enabled, no RELRO segment",
		}
	}

	bindNow, err := elf.HasDynTag(bin, stdelf.DT_BIND_NOW)
	if err != nil {
		return rule.Skip("failed to read dynamic section", err)
	}
	dfBindNow, err := elf.HasDynFlag(bin, stdelf.DT_FLAGS, uint64(stdelf.DF_BIND_NOW))
	if err != nil {
		return rule.Skip("failed to read dynamic section", err)
	}
	df1Now, err := elf.HasDynFlag(bin, stdelf.DT_FLAGS_1, uint64(stdelf.DF_1_NOW))
	if err != nil {
		return rule.Skip("failed to read dynamic section", err)
	}
	if bindNow || dfBindNow || df1Now {
		return rule.Result{
			Status:  rule.StatusPassed,
			Message: "Full RELRO enabled",
		}
	}

	return rule.Result{
		Status:  rule.StatusFailed,
		Message: "Full RELRO not enabled, partial RELRO only",
	}
}
