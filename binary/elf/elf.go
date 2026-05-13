package elf

import (
	"debug/elf"
	"encoding/binary"
	"errors"
)

// ErrSectionMissing is returned when a section cannot be obtained.
var ErrSectionMissing = errors.New("section missing")

// Alignment is a byte boundary used to pad ELF fields per the spec.
type Alignment int

// Pad returns n rounded up to the next multiple of a.
func (a Alignment) Pad(n int) int {
	return ((n + int(a) - 1) / int(a)) * int(a)
}

// Binary exposes ELF metadata and section/segment contents for analysis.
type Binary interface {
	// Class reports the ELF class (32-bit or 64-bit).
	Class() elf.Class
	// Type reports the ELF type (e.g. ET_EXEC, ET_DYN).
	Type() elf.Type
	// Machine reports the target machine.
	Machine() elf.Machine
	// OSABI reports the OS/ABI identification.
	OSABI() elf.OSABI
	// ByteOrder reports the byte order in which this binary's data is encoded.
	ByteOrder() binary.ByteOrder
	// Entry returns the virtual address of the binary's entry point.
	Entry() uint64
	// BuildID returns an opaque identifier used to look up debug artifacts for this binary.
	// The default implementation returns the GNU build ID hex string from .note.gnu.build-id, or "" if absent.
	// Wrappers may substitute any stable identifier (e.g. a manifest hash) that the configured Resolver understands.
	BuildID() string

	// Progs returns the program header table.
	Progs() []Prog
	// Sections returns the section header table.
	Sections() []Section
	// Symbols returns entries from .symtab.
	// Returns (nil, nil) when no .symtab is present, matching stdlib's elf.ErrNoSymbols treatment.
	Symbols() ([]elf.Symbol, error)
	// DynSymbols returns entries from .dynsym.
	// Returns (nil, nil) when the binary has no dynamic symbols.
	DynSymbols() ([]elf.Symbol, error)
	// DynEntries returns parsed entries from the .dynamic section.
	// Returns (nil, nil) when .dynamic is absent.
	DynEntries() ([]DynEntry, error)
}

// Resolver supplies the raw bytes of an ELF section that is not present in the local file.
// Implementations should return ErrSectionMissing when the section is unavailable on the remote side so callers can distinguish "section legitimately absent" from transport or parsing failures.
type Resolver interface {
	FetchSection(name string) ([]byte, error)
}

// Section is an ELF section header bundled with a lazy accessor for its content.
type Section struct {
	elf.SectionHeader
	data func() ([]byte, error)
}

// Data returns the section's raw bytes.
func (s Section) Data() ([]byte, error) {
	if s.data == nil {
		return nil, ErrSectionMissing
	}
	return s.data()
}

// Prog is an ELF program header bundled with a lazy accessor for its segment content.
type Prog struct {
	elf.ProgHeader
	data func() ([]byte, error)
}

// Data returns the segment's raw bytes, or (nil, nil) for segments with no file content.
func (p Prog) Data() ([]byte, error) {
	if p.data == nil {
		return nil, nil
	}
	return p.data()
}

// DynEntry is a parsed entry from the .dynamic section.
type DynEntry struct {
	Tag elf.DynTag
	Val uint64
}

// FindSection returns the named section from the binary, or ErrSectionMissing if it isn't present.
func FindSection(b Binary, name string) (Section, error) {
	for _, sec := range b.Sections() {
		if sec.Name == name {
			return sec, nil
		}
	}
	return Section{}, ErrSectionMissing
}

// findSectionData returns the raw bytes of the named section, or (nil, nil) when the section is absent.
// Any other error is returned unwrapped so callers can add their own context.
func findSectionData(b Binary, name string) ([]byte, error) {
	sec, err := FindSection(b, name)
	if err != nil {
		if errors.Is(err, ErrSectionMissing) {
			return nil, nil
		}
		return nil, err
	}
	data, err := sec.Data()
	if err != nil {
		if errors.Is(err, ErrSectionMissing) {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}
