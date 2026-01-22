package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

// NoPLTRule checks if binary was compiled with -fno-plt
// This avoids PLT indirection, reducing ROP gadgets and improving RELRO effectiveness
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Code-Gen-Options.html#index-fno-plt
// Gentoo: https://wiki.gentoo.org/wiki/Hardened/Toolchain
// Debian: https://wiki.debian.org/Hardening
type NoPLTRule struct{}

func (r NoPLTRule) ID() string   { return "no-plt" }
func (r NoPLTRule) Name() string { return "No PLT" }

func (r NoPLTRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Arch: binary.ArchAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 6, Minor: 0}, Flag: "-fno-plt"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 9}, Flag: "-fno-plt"},
		},
	}
}

func (r NoPLTRule) Execute(f *elf.File, info *binary.Parsed) rule.Result {
	// Skip for static binaries - PLT only applies to dynamically linked binaries
	if f.Section(".dynamic") == nil {
		return rule.Result{
			State:   rule.CheckStateSkipped,
			Message: "Static binary (PLT not applicable)",
		}
	}

	pltSection := f.Section(".plt")
	pltSecSection := f.Section(".plt.sec")

	// No PLT section at all - definitely compiled with -fno-plt
	if pltSection == nil {
		return rule.Result{
			State:   rule.CheckStatePassed,
			Message: "No PLT section (direct GOT access)",
		}
	}

	// .plt.sec is used by Intel CET - if present, PLT is being used
	// but in a hardened way, so we consider this acceptable
	if pltSecSection != nil {
		return rule.Result{
			State:   rule.CheckStatePassed,
			Message: "Using secure PLT (.plt.sec for CET compatibility)",
		}
	}

	return rule.Result{
		State:   rule.CheckStateFailed,
		Message: "PLT is used",
	}
}
