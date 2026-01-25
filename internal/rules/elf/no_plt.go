package elf

import (
	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const NoPLTRuleID = "no-plt"

// NoPLTRule checks if binary was compiled with -fno-plt
// This avoids PLT indirection, reducing ROP gadgets and improving RELRO effectiveness
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Code-Gen-Options.html#index-fno-plt
// Gentoo: https://wiki.gentoo.org/wiki/Hardened/Toolchain
// Debian: https://wiki.debian.org/Hardening
type NoPLTRule struct{}

func (r NoPLTRule) ID() string   { return NoPLTRuleID }
func (r NoPLTRule) Name() string { return "No PLT" }

func (r NoPLTRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 6, Minor: 0}, Flag: "-fno-plt"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 9}, Flag: "-fno-plt"},
		},
	}
}

func (r NoPLTRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	// Skip for static binaries - PLT only applies to dynamically linked binaries
	if bin.File.Section(".dynamic") == nil {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Static binary (PLT not applicable)",
		}
	}

	pltSection := bin.File.Section(".plt")
	pltSecSection := bin.File.Section(".plt.sec")

	// No PLT section at all - definitely compiled with -fno-plt
	if pltSection == nil {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "No PLT section (direct GOT access)",
		}
	}

	// .plt.sec is used by Intel CET - if present, PLT is being used
	// but in a hardened way, so we consider this acceptable
	if pltSecSection != nil {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "Using secure PLT (.plt.sec for CET compatibility)",
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "PLT is used",
	}
}
