package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

// RELRORule checks for partial RELRO
// ld: https://sourceware.org/binutils/docs/ld/Options.html
type RELRORule struct{}

func (r RELRORule) ID() string   { return "relro" }
func (r RELRORule) Name() string { return "Partial RELRO" }

func (r RELRORule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 3, Minor: 0}, DefaultVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,relro"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 0}, DefaultVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-Wl,-z,relro"},
		},
	}
}

func (r RELRORule) Execute(f *elf.File, info *binary.Parsed) rule.ExecuteResult {
	hasRELRO := false
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_GNU_RELRO {
			hasRELRO = true
			break
		}
	}

	if hasRELRO {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "Partial RELRO is enabled (some ELF sections read-only after load)",
		}
	}
	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "Partial RELRO is NOT enabled (ELF sections remain writable)",
	}
}
