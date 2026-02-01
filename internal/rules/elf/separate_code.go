package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const SeparateCodeRuleID = "separate-code"

// SeparateCodeRule checks if code and data are in separate pages
// ld: https://sourceware.org/binutils/docs/ld/Options.html#index-z-keyword
type SeparateCodeRule struct{}

func (r SeparateCodeRule) ID() string   { return SeparateCodeRuleID }
func (r SeparateCodeRule) Name() string { return "Separate Code Segments" }
func (r SeparateCodeRule) Description() string {
	return "Checks if code and data are in separate memory pages. This prevents code pages from being writable and data pages from being executable, reducing the attack surface for memory corruption exploits."
}

func (r SeparateCodeRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 8, Minor: 1}, DefaultVersion: toolchain.Version{Major: 8, Minor: 1}, Flag: "-Wl,-z,separate-code"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 6, Minor: 0}, DefaultVersion: toolchain.Version{Major: 6, Minor: 0}, Flag: "-Wl,-z,separate-code"},
		},
	}
}

func (r SeparateCodeRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	if bin.File.Type != elf.ET_EXEC && bin.File.Type != elf.ET_DYN {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "Not an executable or shared library",
		}
	}

	// Check file offsets at 4KB page granularity.
	const pageSize uint64 = 4096

	var codePages, dataPages [][2]uint64 // [start, end) page ranges

	for _, prog := range bin.File.Progs {
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
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "No code segments found",
		}
	}

	// Check if any code page range overlaps with any data page range.
	for _, code := range codePages {
		for _, data := range dataPages {
			if code[0] < data[1] && code[1] > data[0] {
				return rule.ExecuteResult{
					Status:  rule.StatusFailed,
					Message: "Code and data share pages",
				}
			}
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusPassed,
		Message: "Code and data separated",
	}
}
