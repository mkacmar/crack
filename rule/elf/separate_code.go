package elf

import (
	"debug/elf"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
)

// SeparateCodeRuleID is the rule ID for separate code.
const SeparateCodeRuleID = "separate-code"

// SeparateCodeRule checks if code and data are in separate pages.
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type SeparateCodeRule struct{}

func (r SeparateCodeRule) ID() string   { return SeparateCodeRuleID }
func (r SeparateCodeRule) Name() string { return "Separate Code Segments" }
func (r SeparateCodeRule) Description() string {
	return "Checks if code and data are in separate memory pages. This prevents code pages from being writable and data pages from being executable, reducing the attack surface for memory corruption exploits."
}

func (r SeparateCodeRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 8, Minor: 1}, DefaultVersion: toolchain.Version{Major: 8, Minor: 1}, Flag: "-Wl,-z,separate-code"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 6, Minor: 0}, DefaultVersion: toolchain.Version{Major: 6, Minor: 0}, Flag: "-Wl,-z,separate-code"},
			toolchain.Rustc: {MinVersion: toolchain.Version{Major: 1, Minor: 31}, Flag: "-C link-arg=-z -C link-arg=separate-code"},
		},
	}
}

func (r SeparateCodeRule) Execute(bin *binary.ELFBinary) rule.Result {
	if bin.Type != elf.ET_EXEC && bin.Type != elf.ET_DYN {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	const pageSize uint64 = 4096

	var codePages, dataPages [][2]uint64

	for _, prog := range bin.Progs {
		if prog.Type != elf.PT_LOAD {
			continue
		}

		startPage := prog.Off / pageSize
		endPage := (prog.Off + prog.Filesz + pageSize - 1) / pageSize

		if (prog.Flags & elf.PF_X) != 0 {
			codePages = append(codePages, [2]uint64{startPage, endPage})
		}
		if (prog.Flags & elf.PF_W) != 0 {
			dataPages = append(dataPages, [2]uint64{startPage, endPage})
		}
	}

	if len(codePages) == 0 {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "No code segments found",
		}
	}

	for _, code := range codePages {
		for _, data := range dataPages {
			if code[0] < data[1] && code[1] > data[0] {
				return rule.Result{
					Status:  rule.StatusFailed,
					Message: "Code and data share pages",
				}
			}
		}
	}

	return rule.Result{
		Status:  rule.StatusPassed,
		Message: "Code and data separated",
	}
}
