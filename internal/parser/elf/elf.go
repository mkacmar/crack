package elf

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"strings"

	elfbinary "github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/toolchain"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(path string) (*elfbinary.ELFBinary, error) {
	f, err := elf.Open(path)
	if err != nil {
		if isNotELFError(err) {
			return nil, elfbinary.ErrUnsupportedFormat
		}
		return nil, fmt.Errorf("failed to open ELF file: %w", err)
	}

	bin := &elfbinary.ELFBinary{
		Binary: elfbinary.Binary{
			Path:   path,
			Format: elfbinary.FormatELF,
		},
		File: f,
	}

	bin.Architecture = parseArchitecture(f.Machine)
	if f.Class == elf.ELFCLASS64 {
		bin.Bits = elfbinary.Bits64
	} else {
		bin.Bits = elfbinary.Bits32
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

// extractBuildID parses the .note.gnu.build-id section.
// https://man7.org/linux/man-pages/man5/elf.5.html
func (p *Parser) extractBuildID(f *elf.File) string {
	section := f.Section(".note.gnu.build-id")
	if section == nil {
		return ""
	}

	data, err := section.Data()
	if err != nil {
		return ""
	}

	// namesz (4) + descsz (4) + type (4)
	const noteHeaderSize = 12
	if len(data) < noteHeaderSize {
		return ""
	}

	namesz := binary.LittleEndian.Uint32(data[0:4])
	descsz := binary.LittleEndian.Uint32(data[4:8])

	const align = 4
	nameAligned := (namesz + align - 1) &^ (align - 1)

	descOffset := noteHeaderSize + int(nameAligned)
	if descOffset+int(descsz) > len(data) {
		return ""
	}

	return fmt.Sprintf("%x", data[descOffset:descOffset+int(descsz)])
}

func parseArchitecture(machine elf.Machine) elfbinary.Architecture {
	switch machine {
	case elf.EM_386:
		return elfbinary.ArchX86
	case elf.EM_X86_64:
		return elfbinary.ArchAMD64
	case elf.EM_ARM:
		return elfbinary.ArchARM
	case elf.EM_AARCH64:
		return elfbinary.ArchARM64
	case elf.EM_RISCV:
		return elfbinary.ArchRISCV
	case elf.EM_PPC64:
		return elfbinary.ArchPPC64
	case elf.EM_MIPS:
		return elfbinary.ArchMIPS
	case elf.EM_S390:
		return elfbinary.ArchS390X
	default:
		return elfbinary.ArchUnknown
	}
}

func (p *Parser) detectLibC(f *elf.File) toolchain.LibC {
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_INTERP {
			data := make([]byte, prog.Filesz)
			if _, err := prog.ReadAt(data, 0); err != nil {
				continue
			}
			interpreter := string(bytes.TrimRight(data, "\x00"))

			if strings.Contains(interpreter, "ld-musl") {
				return toolchain.LibCMusl
			}
			if strings.Contains(interpreter, "ld-linux") {
				return toolchain.LibCGlibc
			}
		}
	}

	return toolchain.LibCUnknown
}
