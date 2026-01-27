package elf

import (
	"debug/elf"

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
		Platform: binary.PlatformAllX86,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 6, Minor: 0}, Flag: "-fno-plt"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 9}, Flag: "-fno-plt"},
		},
	}
}

func (r NoPLTRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	if bin.File.Section(".dynamic") == nil {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Static binary, PLT not applicable",
		}
	}

	hasIBT := parseGNUProperty(bin.File, GNU_PROPERTY_X86_FEATURE_1_AND, GNU_PROPERTY_X86_FEATURE_1_IBT)
	if hasIBT {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "CET-IBT enabled, PLT gadgets protected by hardware",
		}
	}

	var hasPLTRelocations bool
	if bin.File.Class == elf.ELFCLASS64 {
		hasPLTRelocations = bin.File.Section(".rela.plt") != nil
	} else {
		hasPLTRelocations = bin.File.Section(".rel.plt") != nil
	}

	if !hasPLTRelocations {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: "PLT not used, direct GOT access",
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "PLT in use, ROP gadgets present",
	}
}
