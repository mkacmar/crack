package elf

import (
	"bytes"
	"debug/elf"
	"fmt"
	"os"

	"github.com/mkacmar/crack/internal/model"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) CanParse(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	magic := make([]byte, 4)
	if _, err := f.Read(magic); err != nil {
		return false, nil
	}

	return magic[0] == 0x7f && magic[1] == 'E' && magic[2] == 'L' && magic[3] == 'F', nil
}

func (p *Parser) Parse(path string) (*model.ParsedBinary, error) {
	f, err := elf.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open ELF file: %w", err)
	}

	info := &model.ParsedBinary{
		Path:    path,
		Format:  model.FormatELF,
		ELFFile: f,
	}

	info.Architecture = parseArchitecture(f.Machine)
	if f.Class == elf.ELFCLASS64 {
		info.Bits = 64
	} else {
		info.Bits = 32
	}

	info.Build = model.CompilerInfo{
		BuildID:   p.extractBuildID(f),
		Toolchain: p.detectToolchain(f),
	}

	return info, nil
}

func (p *Parser) detectToolchain(f *elf.File) model.Toolchain {
	compilerInfo := p.extractCompilerInfo(f)
	return model.ParseToolchain(compilerInfo)
}

func (p *Parser) extractCompilerInfo(f *elf.File) string {
	section := f.Section(".comment")
	if section == nil {
		return ""
	}

	data, err := section.Data()
	if err != nil {
		return ""
	}

	return parseFirstComment(data)
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

func parseFirstComment(data []byte) string {
	for len(data) > 0 {
		idx := bytes.IndexByte(data, 0)
		if idx == -1 {
			return ""
		}
		if idx > 0 {
			return string(data[:idx])
		}
		data = data[idx+1:]
	}
	return ""
}

func parseArchitecture(machine elf.Machine) model.Architecture {
	switch machine {
	case elf.EM_386:
		return model.ArchX86
	case elf.EM_X86_64:
		return model.ArchX86_64
	case elf.EM_ARM:
		return model.ArchARM
	case elf.EM_AARCH64:
		return model.ArchARM64
	case elf.EM_RISCV:
		return model.ArchRISCV
	case elf.EM_PPC64:
		return model.ArchPPC64
	case elf.EM_MIPS:
		return model.ArchMIPS
	case elf.EM_S390:
		return model.ArchS390X
	default:
		return model.ArchUnknown
	}
}
