package elf

import (
	stdelf "debug/elf"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/binary/elf"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
)

// NoDumpRuleID is the rule ID for no dump.
const NoDumpRuleID = "no-dump"

// NoDumpRule checks if core dumps are disabled.
//
// References:
//   - https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type NoDumpRule struct{}

func (r NoDumpRule) ID() string   { return NoDumpRuleID }
func (r NoDumpRule) Name() string { return "Core Dump Protection" }
func (r NoDumpRule) Description() string {
	return "Checks if the binary has the DF_1_NODUMP flag set to prevent dldump(3) from copying the shared object. This restricts an attacker's ability to extract the mapped object from a running process."
}

func (r NoDumpRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM | binary.ArchRISCV},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 1}, Flag: "-Wl,-z,nodump"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 4}, Flag: "-Wl,-z,nodump"},
		},
		LibC: binary.LibCAll,
	}
}

func (r NoDumpRule) Execute(bin elf.Binary) rule.Result {
	if bin.Type() != stdelf.ET_EXEC && bin.Type() != stdelf.ET_DYN {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	nodump, err := elf.HasDynFlag(bin, stdelf.DT_FLAGS_1, uint64(stdelf.DF_1_NODUMP))
	if err != nil {
		return rule.Skip("failed to read dynamic section", err)
	}
	if nodump {
		return rule.Result{
			Status:  rule.StatusPassed,
			Message: "Core dumps disabled",
		}
	}

	return rule.Result{
		Status:  rule.StatusFailed,
		Message: "Core dumps not explicitly disabled",
	}
}
