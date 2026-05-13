package elf

import (
	"debug/elf"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	bin "go.kacmar.sk/crack/binary"
)

// File is the canonical Binary implementation.
// It wraps stdlib's *elf.File and, when configured with a Resolver, transparently fetches sections that have been stripped from the local file.
type File struct {
	file     *elf.File
	resolver Resolver
	buildID  string

	progs    []Prog
	sections []Section

	// Memoized parsed views. Each closure runs at most once thanks to sync.OnceValues.
	symbols    func() ([]elf.Symbol, error)
	dynSymbols func() ([]elf.Symbol, error)
	dynEntries func() ([]DynEntry, error)

	// Cache of raw section bytes keyed by section name, populated lazily on first fetch.
	// Guarded by sectionMu because individual sections may be requested concurrently and the resolver call is the slow path we want to deduplicate.
	sectionMu    sync.Mutex
	sectionBytes map[string]sectionResult
}

// sectionResult memoizes a single section lookup so repeated fetches return the same outcome without re-invoking the resolver.
type sectionResult struct {
	data []byte
	err  error
}

// Option configures a File at construction time.
type Option func(*openConfig)

type openConfig struct {
	resolverFactory func(buildID string) Resolver
}

// WithResolverFactory supplies a factory that produces a Resolver scoped to the binary's build ID.
// The factory is invoked once during Open after the build ID is known. Returning nil disables remote fetching.
// A Resolver is consulted when an accessor needs data that isn't present locally (typically debug symbols on stripped binaries).
func WithResolverFactory(factory func(buildID string) Resolver) Option {
	return func(c *openConfig) { c.resolverFactory = factory }
}

// Open parses the ELF header and section/program header tables from r.
// Section and segment contents are read lazily on demand.
// The caller owns r and must keep it open while the returned binary is in use.
func Open(r io.ReaderAt, opts ...Option) (*File, error) {
	var cfg openConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	f, err := elf.NewFile(r)
	if err != nil {
		if isFormatMismatch(err) {
			return nil, bin.ErrUnsupportedFormat
		}
		return nil, fmt.Errorf("failed to open ELF file: %w", err)
	}

	b := &File{
		file:         f,
		sectionBytes: make(map[string]sectionResult),
	}

	b.progs = make([]Prog, len(f.Progs))
	for i, p := range f.Progs {
		prog := p
		b.progs[i] = Prog{
			ProgHeader: prog.ProgHeader,
			data: func() ([]byte, error) {
				buf := make([]byte, prog.Filesz)
				if _, err := prog.ReadAt(buf, 0); err != nil {
					return nil, fmt.Errorf("failed to read segment: %w", err)
				}
				return buf, nil
			},
		}
	}

	b.sections = make([]Section, len(f.Sections))
	for i, s := range f.Sections {
		name := s.SectionHeader.Name
		b.sections[i] = Section{
			SectionHeader: s.SectionHeader,
			data:          func() ([]byte, error) { return b.sectionDataByName(name) },
		}
	}

	b.buildID = extractBuildID(b)

	if cfg.resolverFactory != nil {
		b.resolver = cfg.resolverFactory(b.buildID)
	}

	b.symbols = sync.OnceValues(b.loadLocalSymbols)
	b.dynSymbols = sync.OnceValues(b.loadLocalDynSymbols)
	b.dynEntries = sync.OnceValues(b.loadDynEntries)

	return b, nil
}

func (b *File) Class() elf.Class            { return b.file.Class }
func (b *File) Type() elf.Type              { return b.file.Type }
func (b *File) Machine() elf.Machine        { return b.file.Machine }
func (b *File) OSABI() elf.OSABI            { return b.file.OSABI }
func (b *File) ByteOrder() binary.ByteOrder { return b.file.ByteOrder }
func (b *File) Entry() uint64               { return b.file.Entry }
func (b *File) BuildID() string             { return b.buildID }
func (b *File) Progs() []Prog               { return b.progs }
func (b *File) Sections() []Section         { return b.sections }

func (b *File) Symbols() ([]elf.Symbol, error)    { return b.symbols() }
func (b *File) DynSymbols() ([]elf.Symbol, error) { return b.dynSymbols() }
func (b *File) DynEntries() ([]DynEntry, error)   { return b.dynEntries() }

func (b *File) sectionDataByName(name string) ([]byte, error) {
	b.sectionMu.Lock()
	defer b.sectionMu.Unlock()
	if r, ok := b.sectionBytes[name]; ok {
		return r.data, r.err
	}
	data, err := b.fetchSectionData(name)
	b.sectionBytes[name] = sectionResult{data: data, err: err}
	return data, err
}

func (b *File) loadLocalSymbols() ([]elf.Symbol, error) {
	syms, err := b.file.Symbols()
	if err != nil && !errors.Is(err, elf.ErrNoSymbols) {
		return nil, fmt.Errorf("failed to read symbols: %w", err)
	}
	if len(syms) > 0 || b.resolver == nil {
		return syms, nil
	}
	return b.loadResolverSymbols()
}

// loadResolverSymbols fetches .symtab and .strtab via the configured resolver and parses them manually.
// Returns (nil, nil) when either section can't be obtained, matching the "no symbols available" semantics of stdlib.
func (b *File) loadResolverSymbols() ([]elf.Symbol, error) {
	symData, err := b.sectionDataByName(".symtab")
	if err != nil {
		if errors.Is(err, ErrSectionMissing) {
			return nil, nil
		}
		return nil, fmt.Errorf("fetch .symtab: %w", err)
	}

	strData, err := b.sectionDataByName(".strtab")
	if err != nil && !errors.Is(err, ErrSectionMissing) {
		return nil, fmt.Errorf("fetch .strtab: %w", err)
	}

	syms, err := parseSymbols(symData, strData, b.file.Class, b.file.ByteOrder)
	if err != nil {
		return nil, fmt.Errorf("parse remote symbols: %w", err)
	}
	return syms, nil
}

func (b *File) loadLocalDynSymbols() ([]elf.Symbol, error) {
	syms, err := b.file.DynamicSymbols()
	if err != nil && !errors.Is(err, elf.ErrNoSymbols) {
		return nil, fmt.Errorf("failed to read dynamic symbols: %w", err)
	}
	return syms, nil
}

func (b *File) loadDynEntries() ([]DynEntry, error) {
	dynSec := b.file.Section(".dynamic")
	if dynSec == nil {
		return nil, nil
	}

	data, err := dynSec.Data()
	if err != nil {
		return nil, fmt.Errorf("failed to read .dynamic: %w", err)
	}

	var entrySize int
	var readEntry func([]byte) (elf.DynTag, uint64)

	if b.file.Class == elf.ELFCLASS64 {
		entrySize = 16
		readEntry = func(d []byte) (elf.DynTag, uint64) {
			// #nosec G115 -- ELF spec tag values fit in DynTag's underlying int.
			tag := elf.DynTag(b.file.ByteOrder.Uint64(d[:8]))
			return tag, b.file.ByteOrder.Uint64(d[8:16])
		}
	} else {
		entrySize = 8
		readEntry = func(d []byte) (elf.DynTag, uint64) {
			return elf.DynTag(b.file.ByteOrder.Uint32(d[:4])), uint64(b.file.ByteOrder.Uint32(d[4:8]))
		}
	}

	var entries []DynEntry
	for i := 0; i+entrySize <= len(data); i += entrySize {
		tag, val := readEntry(data[i:])
		if tag == elf.DT_NULL {
			break
		}
		entries = append(entries, DynEntry{Tag: tag, Val: val})
	}

	return entries, nil
}

func (b *File) fetchSectionData(name string) ([]byte, error) {
	if sec := b.file.Section(name); sec != nil {
		data, err := sec.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to read section %s: %w", name, err)
		}
		return data, nil
	}
	if b.resolver == nil {
		return nil, ErrSectionMissing
	}
	data, err := b.resolver.FetchSection(name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch section %s: %w", name, err)
	}
	return data, nil
}

// isFormatMismatch reports whether err from elf.NewFile signals that the input is not an ELF file (rather than a malformed one).
// The stdlib returns these specific messages for a bad ELF magic or an invalid class byte.
func isFormatMismatch(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "bad magic number") ||
		strings.Contains(msg, "invalid argument")
}
