package elf

import (
	stdbinary "encoding/binary"
	"testing"

	"debug/elf"

	"go.kacmar.sk/crack/binary"
)

// fakeBinary is a minimal Binary used to drive DetectLibC.
type fakeBinary struct {
	progs    []Prog
	dynEntry []DynEntry
	sections []Section
}

func (f *fakeBinary) Class() elf.Class                  { return elf.ELFCLASS64 }
func (f *fakeBinary) Type() elf.Type                    { return elf.ET_EXEC }
func (f *fakeBinary) Machine() elf.Machine              { return elf.EM_X86_64 }
func (f *fakeBinary) OSABI() elf.OSABI                  { return elf.ELFOSABI_NONE }
func (f *fakeBinary) ByteOrder() stdbinary.ByteOrder    { return stdbinary.LittleEndian }
func (f *fakeBinary) Entry() uint64                     { return 0 }
func (f *fakeBinary) BuildID() string                   { return "" }
func (f *fakeBinary) Progs() []Prog                     { return f.progs }
func (f *fakeBinary) Sections() []Section               { return f.sections }
func (f *fakeBinary) Symbols() ([]elf.Symbol, error)    { return nil, nil }
func (f *fakeBinary) DynSymbols() ([]elf.Symbol, error) { return nil, nil }
func (f *fakeBinary) DynEntries() ([]DynEntry, error)   { return f.dynEntry, nil }

// makeInterp builds a PT_INTERP segment carrying the given (NUL-terminated) interpreter path.
func makeInterp(path string) Prog {
	data := []byte(path + "\x00")
	return Prog{
		ProgHeader: elf.ProgHeader{Type: elf.PT_INTERP},
		data:       func() ([]byte, error) { return data, nil },
	}
}

// makeDynamic builds a .dynstr section + DT_NEEDED entries for the given library names.
func makeDynamic(libs ...string) (Section, []DynEntry) {
	var strtab []byte
	strtab = append(strtab, 0)

	entries := make([]DynEntry, 0, len(libs))
	for _, lib := range libs {
		off := uint64(len(strtab))
		strtab = append(strtab, []byte(lib)...)
		strtab = append(strtab, 0)
		entries = append(entries, DynEntry{Tag: elf.DT_NEEDED, Val: off})
	}

	sec := Section{
		SectionHeader: elf.SectionHeader{Name: ".dynstr"},
		data:          func() ([]byte, error) { return strtab, nil },
	}
	return sec, entries
}

func TestDetectLibC(t *testing.T) {
	tests := []struct {
		name  string
		progs []Prog
		libs  []string
		want  binary.LibC
	}{
		{
			name: "static binary: no interp, no needed",
			want: binary.LibCNone,
		},
		{
			name:  "glibc via interpreter (ld-linux)",
			progs: []Prog{makeInterp("/lib64/ld-linux-x86-64.so.2")},
			want:  binary.LibCGlibc,
		},
		{
			name:  "musl via interpreter (ld-musl)",
			progs: []Prog{makeInterp("/lib/ld-musl-x86_64.so.1")},
			want:  binary.LibCMusl,
		},
		{
			name:  "bionic-like interpreter, no libc DT_NEEDED",
			progs: []Prog{makeInterp("/system/bin/linker64")},
			want:  binary.LibCUnknown,
		},
		{
			name: "glibc via DT_NEEDED libc.so.6",
			libs: []string{"libc.so.6"},
			want: binary.LibCGlibc,
		},
		{
			name: "musl via DT_NEEDED libc.musl-*.so.1",
			libs: []string{"libc.musl-x86_64.so.1"},
			want: binary.LibCMusl,
		},
		{
			name: "unrecognized libc.so* DT_NEEDED (bionic)",
			libs: []string{"libc.so"},
			want: binary.LibCUnknown,
		},
		{
			name: "self-contained shared object (only libdl)",
			libs: []string{"libdl.so.2"},
			want: binary.LibCNone,
		},
		{
			name:  "interp unrecognized but DT_NEEDED resolves to glibc",
			progs: []Prog{makeInterp("/system/bin/linker64")},
			libs:  []string{"libc.so.6"},
			want:  binary.LibCGlibc,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fb := &fakeBinary{progs: tc.progs}
			if len(tc.libs) > 0 {
				sec, entries := makeDynamic(tc.libs...)
				fb.sections = []Section{sec}
				fb.dynEntry = entries
			}
			got := DetectLibC(fb)
			if got != tc.want {
				t.Errorf("DetectLibC() = %v, want %v", got, tc.want)
			}
		})
	}
}
