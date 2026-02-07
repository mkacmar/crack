package elf

import (
	"bytes"
	"debug/elf"
	"fmt"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/toolchain"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(path string) (*binary.ELFBinary, error) {
	f, err := elf.Open(path)
	if err != nil {
		if isNotELFError(err) {
			return nil, binary.ErrUnsupportedFormat
		}
		return nil, fmt.Errorf("failed to open ELF file: %w", err)
	}

	bin := &binary.ELFBinary{
		Binary: binary.Binary{
			Path:   path,
			Format: binary.FormatELF,
		},
		File: f,
	}

	bin.Architecture = parseArchitecture(f.Machine)
	if f.Class == elf.ELFCLASS64 {
		bin.Bits = binary.Bits64
	} else {
		bin.Bits = binary.Bits32
	}

	bin.Build = toolchain.CompilerInfo{
		BuildID:   p.extractBuildID(f),
		Toolchain: p.detectToolchain(f),
	}

	bin.LibC = p.detectLibC(f)

	if syms, err := f.Symbols(); err == nil {
		bin.Symbols = syms
	}
	if dynsyms, err := f.DynamicSymbols(); err == nil {
		bin.DynSymbols = dynsyms
	}

	return bin, nil
}

func isNotELFError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "bad magic number") ||
		strings.Contains(msg, "invalid argument")
}

// compilerPriority resolves conflicts when multiple .comment entries exist (e.g., glibc embeds GCC).
var compilerPriority = map[toolchain.Compiler]int{
	toolchain.CompilerGCC:   1,
	toolchain.CompilerClang: 2,
	toolchain.CompilerRustc: 3,
}

func (p *Parser) detectToolchain(f *elf.File) toolchain.Toolchain {
	comments := p.extractCompilerComments(f)

	var best toolchain.Toolchain
	for _, comment := range comments {
		tc := toolchain.ParseToolchain(comment)
		if tc.Compiler == toolchain.CompilerUnknown {
			continue
		}
		if best.Compiler == toolchain.CompilerUnknown || compilerPriority[tc.Compiler] > compilerPriority[best.Compiler] {
			best = tc
		}
	}
	return best
}

func (p *Parser) extractCompilerComments(f *elf.File) []string {
	section := f.Section(".comment")
	if section == nil {
		return nil
	}

	data, err := section.Data()
	if err != nil {
		return nil
	}

	return parseComments(data)
}

func (p *Parser) extractBuildID(f *elf.File) string {
	section := f.Section(".note.gnu.build-id")
	if section == nil {
		return ""
	}

	data, err := section.Data()
	if err != nil || len(data) < 16 {
		return ""
	}

	// ELF note structure:
	// namesz (4), descsz (4), type (4), name (aligned), desc (build-id)
	namesz := uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
	descsz := uint32(data[4]) | uint32(data[5])<<8 | uint32(data[6])<<16 | uint32(data[7])<<24

	descOffset := 12 + int((namesz+3)&^3)
	if descOffset+int(descsz) > len(data) {
		return ""
	}

	return fmt.Sprintf("%x", data[descOffset:descOffset+int(descsz)])
}

// parseComments extracts all compiler info strings from .comment section.
func parseComments(data []byte) []string {
	var comments []string
	for len(data) > 0 {
		idx := bytes.IndexByte(data, 0)
		if idx == -1 {
			break
		}
		if idx > 0 {
			comments = append(comments, string(data[:idx]))
		}
		data = data[idx+1:]
	}
	return comments
}

func parseArchitecture(machine elf.Machine) binary.Architecture {
	switch machine {
	case elf.EM_386:
		return binary.ArchX86
	case elf.EM_X86_64:
		return binary.ArchX86_64
	case elf.EM_ARM:
		return binary.ArchARM
	case elf.EM_AARCH64:
		return binary.ArchARM64
	case elf.EM_RISCV:
		return binary.ArchRISCV
	case elf.EM_PPC64:
		return binary.ArchPPC64
	case elf.EM_MIPS:
		return binary.ArchMIPS
	case elf.EM_S390:
		return binary.ArchS390X
	default:
		return binary.ArchUnknown
	}
}

func (p *Parser) detectLibC(f *elf.File) toolchain.LibC {
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_INTERP {
			data := make([]byte, prog.Filesz)
			if _, err := prog.ReadAt(data, 0); err != nil {
				continue
			}
			interp := string(bytes.TrimRight(data, "\x00"))

			if bytes.Contains([]byte(interp), []byte("ld-musl")) {
				return toolchain.LibCMusl
			}
			if bytes.Contains([]byte(interp), []byte("ld-linux")) {
				return toolchain.LibCGlibc
			}
		}
	}

	return toolchain.LibCUnknown
}
