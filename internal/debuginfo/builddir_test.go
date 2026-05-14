package debuginfo

import (
	"bytes"
	"context"
	stdelf "debug/elf"
	stdbinary "encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"go.kacmar.sk/crack/binary/elf"
)

const testBuildID = "abcdef0123456789abcdef0123456789abcdef01"

func newTestResolver(root string) elf.Resolver {
	return NewBuildIDDirSource(root, nil).ResolverFor(context.Background(), testBuildID)
}

func TestBuildIDDirSourceMissingFile(t *testing.T) {
	root := t.TempDir()
	_, err := newTestResolver(root).FetchSection(".symtab")
	if !errors.Is(err, elf.ErrSectionMissing) {
		t.Fatalf("got %v, want ErrSectionMissing", err)
	}
}

func TestBuildIDDirSourceFetchSection(t *testing.T) {
	root := t.TempDir()
	want := []byte("hello world")
	writeFakeDebugFile(t, root, testBuildID, ".symtab", want)

	got, err := newTestResolver(root).FetchSection(".symtab")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBuildIDDirSourceSectionMissingInFile(t *testing.T) {
	root := t.TempDir()
	writeFakeDebugFile(t, root, testBuildID, ".symtab", []byte("data"))

	_, err := newTestResolver(root).FetchSection(".strtab")
	if !errors.Is(err, elf.ErrSectionMissing) {
		t.Fatalf("got %v, want ErrSectionMissing", err)
	}
}

func TestBuildIDDirSourceInvalidELF(t *testing.T) {
	root := t.TempDir()
	path := debugFilePathFor(root, testBuildID)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("not an elf"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := newTestResolver(root).FetchSection(".symtab")
	if err == nil {
		t.Fatal("expected error for invalid ELF")
	}
	if errors.Is(err, elf.ErrSectionMissing) {
		t.Fatalf("invalid ELF should surface a real error, got ErrSectionMissing")
	}
}

func TestBuildIDDirSourceDefaultRoot(t *testing.T) {
	src := NewBuildIDDirSource("", nil)
	if src.root != DefaultBuildIDDir {
		t.Fatalf("empty root should fall back to default, got %q", src.root)
	}
}

func TestBuildIDDirSourceEmptyBuildIDPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for empty buildID")
		}
	}()
	_ = NewBuildIDDirSource("/tmp", nil).ResolverFor(context.Background(), "")
}

func TestBuildIDDirSourcePathLayout(t *testing.T) {
	root := t.TempDir()
	writeFakeDebugFile(t, root, testBuildID, ".symtab", []byte("payload"))

	expected := filepath.Join(root, ".build-id", testBuildID[:2], testBuildID[2:]+".debug")
	if _, err := os.Stat(expected); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("expected debug file at %s", expected)
		}
		t.Fatal(err)
	}
}

// debugFilePathFor mirrors buildIDDirResolver.debugFilePath for tests.
func debugFilePathFor(root, buildID string) string {
	return filepath.Join(root, ".build-id", buildID[:2], buildID[2:]+".debug")
}

// writeFakeDebugFile writes a minimal ELF64 file containing one named section with the given payload.
// The file is placed at the build-id-indexed path under root.
func writeFakeDebugFile(t *testing.T, root, buildID, sectionName string, payload []byte) {
	t.Helper()
	path := debugFilePathFor(root, buildID)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, makeMinimalELF(sectionName, payload), 0o644); err != nil {
		t.Fatal(err)
	}
}

// makeMinimalELF builds a minimal valid little-endian ELF64 file containing a single named section.
// Layout: [ELF header][section data][.shstrtab][section header table].
func makeMinimalELF(sectionName string, data []byte) []byte {
	names := []byte{0}
	sectionNameOff := uint32(len(names))
	names = append(append(names, sectionName...), 0)
	shstrtabNameOff := uint32(len(names))
	names = append(append(names, ".shstrtab"...), 0)

	const headerSize = 64
	const sectionHeaderSize = 64
	dataOff := uint64(headerSize)
	namesOff := dataOff + uint64(len(data))
	sectionHeadersOff := namesOff + uint64(len(names))

	header := stdelf.Header64{
		Ident: [16]byte{
			0x7f, 'E', 'L', 'F',
			byte(stdelf.ELFCLASS64),
			byte(stdelf.ELFDATA2LSB),
			byte(stdelf.EV_CURRENT),
		},
		Type:      uint16(stdelf.ET_EXEC),
		Machine:   uint16(stdelf.EM_X86_64),
		Version:   uint32(stdelf.EV_CURRENT),
		Shoff:     sectionHeadersOff,
		Ehsize:    headerSize,
		Shentsize: sectionHeaderSize,
		Shnum:     3,
		Shstrndx:  2,
	}

	sections := []stdelf.Section64{
		{},
		{
			Name:      sectionNameOff,
			Type:      uint32(stdelf.SHT_PROGBITS),
			Off:       dataOff,
			Size:      uint64(len(data)),
			Addralign: 1,
		},
		{
			Name:      shstrtabNameOff,
			Type:      uint32(stdelf.SHT_STRTAB),
			Off:       namesOff,
			Size:      uint64(len(names)),
			Addralign: 1,
		},
	}

	var buf bytes.Buffer
	mustWrite(&buf, &header)
	buf.Write(data)
	buf.Write(names)
	for i := range sections {
		mustWrite(&buf, &sections[i])
	}
	return buf.Bytes()
}

func mustWrite(buf *bytes.Buffer, v any) {
	if err := stdbinary.Write(buf, stdbinary.LittleEndian, v); err != nil {
		panic(fmt.Sprintf("test setup: binary.Write: %v", err))
	}
}
