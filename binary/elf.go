package binary

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mkacmar/crack/toolchain"
)

// ErrUnsupportedFormat is returned when the file is not a supported binary format.
var ErrUnsupportedFormat = errors.New("unsupported binary format")

// GNU property constants for feature detection.
const (
	NT_GNU_PROPERTY_TYPE_0 = 5

	GNU_PROPERTY_X86_FEATURE_1_AND   = 0xc0000002
	GNU_PROPERTY_X86_FEATURE_1_IBT   = 0x1
	GNU_PROPERTY_X86_FEATURE_1_SHSTK = 0x2

	GNU_PROPERTY_AARCH64_FEATURE_1_AND = 0xc0000000
	GNU_PROPERTY_AARCH64_FEATURE_1_BTI = 0x1
	GNU_PROPERTY_AARCH64_FEATURE_1_PAC = 0x2
)

// DynEntry represents an entry in the .dynamic section.
type DynEntry struct {
	Tag uint64
	Val uint64
}

// ELFBinary represents a parsed ELF executable or shared library.
type ELFBinary struct {
	Binary

	file *elf.File

	Type     elf.Type
	Progs    []elf.Prog
	Sections []elf.SectionHeader
	// Symbols from .symtab section, or nil if missing/stripped.
	Symbols []elf.Symbol
	// DynSymbols from .dynsym section, or nil if missing.
	DynSymbols []elf.Symbol

	// DynEntries contains parsed entries from the .dynamic section.
	// Exposed for custom rules that need direct access to dynamic tags
	// beyond what HasDynFlag, HasDynTag, and DynString provide.
	DynEntries []DynEntry
}

// ParseELF parses an ELF binary.
// The caller owns the io.ReaderAt and must keep it open while ELFBinary is in use.
func ParseELF(r io.ReaderAt) (*ELFBinary, error) {
	return ParseELFWithDetector(r, toolchain.ELFCommentDetector{})
}

// ParseELFWithDetector parses an ELF binary using a custom compiler detector.
// The caller owns the io.ReaderAt and must keep it open while ELFBinary is in use.
func ParseELFWithDetector(r io.ReaderAt, detector toolchain.ELFDetector) (*ELFBinary, error) {
	f, err := elf.NewFile(r)
	if err != nil {
		if isNotELFError(err) {
			return nil, ErrUnsupportedFormat
		}
		return nil, fmt.Errorf("failed to open ELF file: %w", err)
	}

	bin := &ELFBinary{
		Binary: Binary{
			Format: FormatELF,
		},
		file: f,
		Type: f.Type,
	}

	bin.Progs = make([]elf.Prog, len(f.Progs))
	for i, p := range f.Progs {
		bin.Progs[i] = *p
	}

	bin.Sections = make([]elf.SectionHeader, len(f.Sections))
	for i, s := range f.Sections {
		bin.Sections[i] = s.SectionHeader
	}

	bin.DynEntries = parseDyn(f)

	bin.Architecture = parseArchitecture(f.Machine)
	if f.Class == elf.ELFCLASS64 {
		bin.Bits = Bits64
	} else {
		bin.Bits = Bits32
	}

	bin.Build = toolchain.BuildInfo{
		BuildID: extractBuildID(f),
	}
	bin.Build.Compiler, bin.Build.Version = detectToolchain(f, detector)

	bin.LibC = detectLibC(f)

	syms, err := f.Symbols()
	if err != nil && !errors.Is(err, elf.ErrNoSymbols) {
		return nil, fmt.Errorf("failed to read symbols: %w", err)
	}
	bin.Symbols = syms

	dynsyms, err := f.DynamicSymbols()
	if err != nil && !errors.Is(err, elf.ErrNoSymbols) {
		return nil, fmt.Errorf("failed to read dynamic symbols: %w", err)
	}
	bin.DynSymbols = dynsyms

	return bin, nil
}

// HasDynFlag reports whether a dynamic tag has the specified flag set.
func (b *ELFBinary) HasDynFlag(tag elf.DynTag, flag uint64) bool {
	for _, entry := range b.DynEntries {
		if entry.Tag == uint64(tag) && (entry.Val&flag) != 0 {
			return true
		}
	}
	return false
}

// HasDynTag reports whether a dynamic tag exists.
func (b *ELFBinary) HasDynTag(tag elf.DynTag) bool {
	for _, entry := range b.DynEntries {
		if entry.Tag == uint64(tag) {
			return true
		}
	}
	return false
}

// DynString returns the string value associated with a dynamic tag.
func (b *ELFBinary) DynString(tag elf.DynTag) string {
	dynstr := b.file.Section(".dynstr")
	if dynstr == nil {
		return ""
	}
	strtab, err := dynstr.Data()
	if err != nil {
		return ""
	}

	for _, entry := range b.DynEntries {
		if entry.Tag == uint64(tag) {
			if int(entry.Val) >= len(strtab) {
				return ""
			}
			end := int(entry.Val)
			for end < len(strtab) && strtab[end] != 0 {
				end++
			}
			return string(strtab[entry.Val:end])
		}
	}
	return ""
}

// HasGNUProperty reports whether the binary has a GNU property with the specified flag.
func (b *ELFBinary) HasGNUProperty(propertyType, featureFlag uint32) bool {
	sec := b.file.Section(".note.gnu.property")
	if sec == nil {
		return false
	}

	data, err := sec.Data()
	if err != nil || len(data) < 16 {
		return false
	}

	align := 4
	if b.file.Class == elf.ELFCLASS64 {
		align = 8
	}

	offset := 0
	for offset+12 <= len(data) {
		namesz := b.file.ByteOrder.Uint32(data[offset : offset+4])
		descsz := b.file.ByteOrder.Uint32(data[offset+4 : offset+8])
		noteType := b.file.ByteOrder.Uint32(data[offset+8 : offset+12])

		alignedNamesz := (int(namesz) + 3) &^ 3
		alignedDescsz := (int(descsz) + align - 1) &^ (align - 1)

		nameStart := offset + 12
		descStart := nameStart + alignedNamesz

		if descStart+alignedDescsz > len(data) {
			break
		}

		if noteType == NT_GNU_PROPERTY_TYPE_0 && namesz >= 4 {
			name := string(data[nameStart : nameStart+4])
			if name == "GNU\x00" {
				propOffset := descStart
				propEnd := descStart + int(descsz)
				for propOffset+8 <= propEnd {
					propType := b.file.ByteOrder.Uint32(data[propOffset : propOffset+4])
					propSize := b.file.ByteOrder.Uint32(data[propOffset+4 : propOffset+8])

					if propType == propertyType && propSize >= 4 {
						features := b.file.ByteOrder.Uint32(data[propOffset+8 : propOffset+12])
						if features&featureFlag != 0 {
							return true
						}
					}

					alignedPropSize := (int(propSize) + align - 1) &^ (align - 1)
					propOffset += 8 + alignedPropSize
				}
			}
		}

		offset += 12 + alignedNamesz + alignedDescsz
	}

	return false
}

var compilerPriority = map[toolchain.Compiler]int{
	toolchain.GCC:   1,
	toolchain.Clang: 2,
	toolchain.Rustc: 3,
}

func isNotELFError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "bad magic number") ||
		strings.Contains(msg, "invalid argument")
}

func detectToolchain(f *elf.File, detector toolchain.ELFDetector) (toolchain.Compiler, toolchain.Version) {
	comments := extractCompilerComments(f)

	var bestComp toolchain.Compiler
	var bestVer toolchain.Version
	bestPriority := 0
	for _, comment := range comments {
		comp, ver := detector.Detect(comment)
		if comp == toolchain.Unknown {
			continue
		}
		priority := compilerPriority[comp]
		if bestComp == toolchain.Unknown || priority > bestPriority {
			bestComp = comp
			bestVer = ver
			bestPriority = priority
		}
	}
	return bestComp, bestVer
}

func extractCompilerComments(f *elf.File) []string {
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

func extractBuildID(f *elf.File) string {
	section := f.Section(".note.gnu.build-id")
	if section == nil {
		return ""
	}

	data, err := section.Data()
	if err != nil {
		return ""
	}

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

func parseArchitecture(machine elf.Machine) Architecture {
	switch machine {
	case elf.EM_386:
		return ArchX86
	case elf.EM_X86_64:
		return ArchAMD64
	case elf.EM_ARM:
		return ArchARM
	case elf.EM_AARCH64:
		return ArchARM64
	case elf.EM_RISCV:
		return ArchRISCV
	case elf.EM_PPC64:
		return ArchPPC64
	case elf.EM_MIPS:
		return ArchMIPS
	case elf.EM_S390:
		return ArchS390X
	default:
		return ArchUnknown
	}
}

func detectLibC(f *elf.File) LibC {
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_INTERP {
			data := make([]byte, prog.Filesz)
			if _, err := prog.ReadAt(data, 0); err != nil {
				continue
			}
			interpreter := string(bytes.TrimRight(data, "\x00"))

			if strings.Contains(interpreter, "ld-musl") {
				return LibCMusl
			}
			if strings.Contains(interpreter, "ld-linux") {
				return LibCGlibc
			}
		}
	}

	return LibCUnknown
}

func parseDyn(f *elf.File) []DynEntry {
	dynSec := f.Section(".dynamic")
	if dynSec == nil {
		return nil
	}

	data, err := dynSec.Data()
	if err != nil {
		return nil
	}

	var entrySize int
	var readEntry func([]byte) (tag, val uint64)

	if f.Class == elf.ELFCLASS64 {
		entrySize = 16
		readEntry = func(d []byte) (uint64, uint64) {
			return f.ByteOrder.Uint64(d[:8]), f.ByteOrder.Uint64(d[8:16])
		}
	} else {
		entrySize = 8
		readEntry = func(d []byte) (uint64, uint64) {
			return uint64(f.ByteOrder.Uint32(d[:4])), uint64(f.ByteOrder.Uint32(d[4:8]))
		}
	}

	var entries []DynEntry
	for i := 0; i+entrySize <= len(data); i += entrySize {
		tag, val := readEntry(data[i:])
		if tag == uint64(elf.DT_NULL) {
			break
		}
		entries = append(entries, DynEntry{Tag: tag, Val: val})
	}

	return entries
}
