package elf

import (
	"debug/elf"
	"encoding/binary"
	"fmt"
)

// ELF symbol table entry sizes per the ELF specification.
// See https://refspecs.linuxfoundation.org/elf/gabi4+/ch4.symtab.html (Elf32_Sym, Elf64_Sym).
const (
	sym32Size = 16
	sym64Size = 24
)

// parseSymbols decodes a raw .symtab byte stream into stdlib elf.Symbol values, resolving names against strtab.
//
// data and strtab may be nil. Nil or empty data returns (nil, nil), and names not resolvable in strtab become "".
// Returns an error only when an entry can't be decoded due to truncation or an unsupported ELF class.
func parseSymbols(data, strtab []byte, class elf.Class, bo binary.ByteOrder) ([]elf.Symbol, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var entrySize int
	switch class {
	case elf.ELFCLASS32:
		entrySize = sym32Size
	case elf.ELFCLASS64:
		entrySize = sym64Size
	default:
		return nil, fmt.Errorf("parseSymbols: unsupported ELF class %v", class)
	}

	if len(data)%entrySize != 0 {
		return nil, fmt.Errorf("parseSymbols: data length %d not a multiple of entry size %d", len(data), entrySize)
	}

	count := len(data) / entrySize
	if count == 0 {
		return nil, nil
	}

	syms := make([]elf.Symbol, 0, count-1)
	for i := 0; i < count; i++ {
		off := i * entrySize
		var sym elf.Symbol
		var nameOff uint32

		if class == elf.ELFCLASS64 {
			nameOff = bo.Uint32(data[off : off+4])
			sym.Info = data[off+4]
			sym.Other = data[off+5]
			sym.Section = elf.SectionIndex(bo.Uint16(data[off+6 : off+8]))
			sym.Value = bo.Uint64(data[off+8 : off+16])
			sym.Size = bo.Uint64(data[off+16 : off+24])
		} else {
			nameOff = bo.Uint32(data[off : off+4])
			sym.Value = uint64(bo.Uint32(data[off+4 : off+8]))
			sym.Size = uint64(bo.Uint32(data[off+8 : off+12]))
			sym.Info = data[off+12]
			sym.Other = data[off+13]
			sym.Section = elf.SectionIndex(bo.Uint16(data[off+14 : off+16]))
		}

		// Skip the conventional STN_UNDEF entry at index 0 to mirror stdlib elf.File.Symbols behavior.
		if i == 0 {
			continue
		}

		sym.Name = lookupStr(strtab, nameOff)
		syms = append(syms, sym)
	}

	return syms, nil
}
