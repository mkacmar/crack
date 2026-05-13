package elf

import (
	"debug/elf"
	"encoding/binary"
	"reflect"
	"testing"
)

func TestParseSymbols_EmptyData(t *testing.T) {
	syms, err := parseSymbols(nil, nil, elf.ELFCLASS64, binary.LittleEndian)
	if err != nil || syms != nil {
		t.Fatalf("nil data: got (%v, %v)", syms, err)
	}
	syms, err = parseSymbols([]byte{}, []byte{}, elf.ELFCLASS64, binary.LittleEndian)
	if err != nil || syms != nil {
		t.Fatalf("empty data: got (%v, %v)", syms, err)
	}
}

func TestParseSymbols_UnsupportedClass(t *testing.T) {
	_, err := parseSymbols(make([]byte, 16), nil, elf.ELFCLASSNONE, binary.LittleEndian)
	if err == nil {
		t.Fatal("expected error for unsupported class")
	}
}

func TestParseSymbols_TruncatedData(t *testing.T) {
	_, err := parseSymbols(make([]byte, 23), nil, elf.ELFCLASS64, binary.LittleEndian)
	if err == nil {
		t.Fatal("expected error for non-multiple data length")
	}
}

func TestParseSymbols_Matrix(t *testing.T) {
	strtab := append([]byte{0}, append([]byte("main"), 0)...)
	strtab = append(strtab, []byte("printf")...)
	strtab = append(strtab, 0)

	mainOff := uint32(1)
	printfOff := uint32(6)

	cases := []struct {
		name  string
		class elf.Class
		bo    binary.ByteOrder
	}{
		{"elf64-le", elf.ELFCLASS64, binary.LittleEndian},
		{"elf64-be", elf.ELFCLASS64, binary.BigEndian},
		{"elf32-le", elf.ELFCLASS32, binary.LittleEndian},
		{"elf32-be", elf.ELFCLASS32, binary.BigEndian},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := buildSymtab(t, tc.class, tc.bo, []elf.Symbol{
				{Name: "main", Info: byte(elf.STT_FUNC), Value: 0x1000, Size: 0x40, Section: 1},
				{Name: "printf", Info: byte(elf.STT_FUNC), Value: 0x2000, Size: 0, Section: elf.SectionIndex(elf.SHN_UNDEF)},
			}, []uint32{mainOff, printfOff})

			got, err := parseSymbols(data, strtab, tc.class, tc.bo)
			if err != nil {
				t.Fatalf("parseSymbols: %v", err)
			}
			want := []elf.Symbol{
				{Name: "main", Info: byte(elf.STT_FUNC), Value: 0x1000, Size: 0x40, Section: 1},
				{Name: "printf", Info: byte(elf.STT_FUNC), Value: 0x2000, Size: 0, Section: elf.SectionIndex(elf.SHN_UNDEF)},
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("symbols mismatch:\n got=%+v\nwant=%+v", got, want)
			}
		})
	}
}

func TestParseSymbols_MissingStrtab(t *testing.T) {
	data := buildSymtab(t, elf.ELFCLASS64, binary.LittleEndian, []elf.Symbol{
		{Info: byte(elf.STT_FUNC), Value: 0x1000},
	}, []uint32{1})

	got, err := parseSymbols(data, nil, elf.ELFCLASS64, binary.LittleEndian)
	if err != nil {
		t.Fatalf("parseSymbols: %v", err)
	}
	if len(got) != 1 || got[0].Name != "" {
		t.Fatalf("expected one symbol with empty name, got %+v", got)
	}
}

func TestParseSymbols_NameZero(t *testing.T) {
	strtab := []byte("\x00main\x00")
	data := buildSymtab(t, elf.ELFCLASS64, binary.LittleEndian, []elf.Symbol{
		{Info: byte(elf.STT_FUNC)},
	}, []uint32{0})

	got, err := parseSymbols(data, strtab, elf.ELFCLASS64, binary.LittleEndian)
	if err != nil {
		t.Fatalf("parseSymbols: %v", err)
	}
	if len(got) != 1 || got[0].Name != "" {
		t.Fatalf("expected one symbol with empty name, got %+v", got)
	}
}

func TestParseSymbols_NameOverflow(t *testing.T) {
	strtab := []byte("\x00main\x00")
	data := buildSymtab(t, elf.ELFCLASS64, binary.LittleEndian, []elf.Symbol{
		{Info: byte(elf.STT_FUNC)},
	}, []uint32{99})

	got, err := parseSymbols(data, strtab, elf.ELFCLASS64, binary.LittleEndian)
	if err != nil {
		t.Fatalf("parseSymbols: %v", err)
	}
	if len(got) != 1 || got[0].Name != "" {
		t.Fatalf("expected one symbol with empty name, got %+v", got)
	}
}

func TestParseSymbols_UnterminatedName(t *testing.T) {
	strtab := []byte("\x00main") // no trailing NUL
	data := buildSymtab(t, elf.ELFCLASS64, binary.LittleEndian, []elf.Symbol{
		{Info: byte(elf.STT_FUNC)},
	}, []uint32{1})

	got, err := parseSymbols(data, strtab, elf.ELFCLASS64, binary.LittleEndian)
	if err != nil {
		t.Fatalf("parseSymbols: %v", err)
	}
	if len(got) != 1 || got[0].Name != "main" {
		t.Fatalf("expected name 'main' from unterminated strtab, got %+v", got)
	}
}

// buildSymtab encodes syms (preceded by the conventional zero entry at index 0) using nameOffsets per symbol.
func buildSymtab(t *testing.T, class elf.Class, bo binary.ByteOrder, syms []elf.Symbol, nameOffsets []uint32) []byte {
	t.Helper()
	if len(syms) != len(nameOffsets) {
		t.Fatalf("syms / nameOffsets length mismatch")
	}

	var entrySize int
	if class == elf.ELFCLASS64 {
		entrySize = sym64Size
	} else {
		entrySize = sym32Size
	}

	out := make([]byte, entrySize*(len(syms)+1))
	for i, sym := range syms {
		off := (i + 1) * entrySize
		if class == elf.ELFCLASS64 {
			bo.PutUint32(out[off:off+4], nameOffsets[i])
			out[off+4] = sym.Info
			out[off+5] = sym.Other
			bo.PutUint16(out[off+6:off+8], uint16(sym.Section))
			bo.PutUint64(out[off+8:off+16], sym.Value)
			bo.PutUint64(out[off+16:off+24], sym.Size)
		} else {
			bo.PutUint32(out[off:off+4], nameOffsets[i])
			bo.PutUint32(out[off+4:off+8], uint32(sym.Value))
			bo.PutUint32(out[off+8:off+12], uint32(sym.Size))
			out[off+12] = sym.Info
			out[off+13] = sym.Other
			bo.PutUint16(out[off+14:off+16], uint16(sym.Section))
		}
	}
	return out
}
