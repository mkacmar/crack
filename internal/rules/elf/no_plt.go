package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

// NoPLTRule checks if binary was compiled with -fno-plt
// This avoids PLT indirection, reducing ROP gadgets and improving RELRO effectiveness
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Code-Gen-Options.html#index-fno-plt
// Gentoo: https://wiki.gentoo.org/wiki/Hardened/Toolchain
// Debian: https://wiki.debian.org/Hardening
type NoPLTRule struct{}

func (r NoPLTRule) ID() string                 { return "no-plt" }
func (r NoPLTRule) Name() string               { return "No PLT" }

func (r NoPLTRule) Applicability() model.Applicability {
	return model.Applicability{
		Arch: model.ArchAll,
		Compilers: map[model.Compiler]model.CompilerRequirement{
			model.CompilerGCC:   {MinVersion: model.Version{Major: 6, Minor: 0}, Flag: "-fno-plt"},
			model.CompilerClang: {MinVersion: model.Version{Major: 3, Minor: 9}, Flag: "-fno-plt"},
		},
	}
}

func (r NoPLTRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	// Skip for static binaries - PLT only applies to dynamically linked binaries
	if f.Section(".dynamic") == nil {
		return model.RuleResult{
			State:   model.CheckStateSkipped,
			Message: "Static binary (PLT not applicable)",
		}
	}

	pltSection := f.Section(".plt")
	pltSecSection := f.Section(".plt.sec")

	// No PLT section at all - definitely compiled with -fno-plt
	if pltSection == nil {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "No PLT section (direct GOT access)",
		}
	}

	// .plt.sec is used by Intel CET - if present, PLT is being used
	// but in a hardened way, so we consider this acceptable
	if pltSecSection != nil {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Using secure PLT (.plt.sec for CET compatibility)",
		}
	}

	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "PLT is used",
	}
}
