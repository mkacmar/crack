package elf

import (
	"bytes"
	"debug/dwarf"
	"encoding/binary"
	"strings"

	"go.kacmar.sk/crack/toolchain"
)

var defaultStringDetector toolchain.StringDetector = toolchain.DefaultStringDetector{}

// ToolchainDetector identifies the compiler and version that produced an ELF binary.
type ToolchainDetector interface {
	Detect(b Binary) toolchain.Toolchain
}

// DefaultToolchainDetector is the bundled toolchain-detection chain.
type DefaultToolchainDetector struct {
	// StringDetector overrides the free-form string classifier.
	// A nil value uses toolchain.DefaultStringDetector.
	StringDetector toolchain.StringDetector
}

// compilerPrecedence orders compilers from most specific to least.
// Multi-language binaries can carry markers for more than one toolchain.
// The earliest match in this list wins.
// New toolchain.Compiler values must be added here to participate in detection ranking.
var compilerPrecedence = []toolchain.Compiler{
	toolchain.Go,
	toolchain.Rustc,
	toolchain.Clang,
	toolchain.GCC,
}

func (d DefaultToolchainDetector) Detect(b Binary) toolchain.Toolchain {
	sd := d.StringDetector
	if sd == nil {
		sd = defaultStringDetector
	}
	if tc, ok := detectGoBuildInfo(b); ok {
		return tc
	}
	if tc := detectFromComment(b, sd); tc.Compiler != toolchain.Unknown {
		return tc
	}
	if tc := detectFromDWARF(b, sd); tc.Compiler != toolchain.Unknown {
		return tc
	}
	return toolchain.Toolchain{}
}

// detectGoBuildInfo reads the version embedded in .go.buildinfo by the Go linker.
// The section layout is documented at https://pkg.go.dev/debug/buildinfo.
// Binaries built with Go 1.18+ store the version inline.
// Older binaries reference it through a pointer and are reported without a version.
func detectGoBuildInfo(b Binary) (toolchain.Toolchain, bool) {
	data, err := findSectionData(b, ".go.buildinfo")
	if err != nil || len(data) < 32 {
		return toolchain.Toolchain{}, false
	}
	if !bytes.HasPrefix(data, []byte("\xff Go buildinf:")) {
		return toolchain.Toolchain{}, false
	}
	tc := toolchain.Toolchain{Compiler: toolchain.Go}
	const flagsVersionInl = 0x2
	if data[15]&flagsVersionInl == 0 {
		return tc, true
	}
	length, n := binary.Uvarint(data[32:])
	if n <= 0 {
		return tc, true
	}
	start := 32 + n
	// #nosec G115 -- start <= len(data) since binary.Uvarint cannot consume more bytes than its input.
	remaining := uint64(len(data) - start)
	if length > remaining {
		return tc, true
	}
	// #nosec G115 -- length <= remaining (checked above), which fits in int.
	vers := string(data[start : start+int(length)])
	if v, ok := parseGoVersion(vers); ok {
		tc.Version = v
	}
	return tc, true
}

// parseGoVersion extracts a semantic version from a Go toolchain version string.
// Accepts forms like "go1.21.3" or "go1.21" and ignores pre-release suffixes.
func parseGoVersion(s string) (toolchain.Version, bool) {
	s = strings.TrimPrefix(s, "go")
	if i := strings.IndexAny(s, "-+ "); i >= 0 {
		s = s[:i]
	}
	v, err := toolchain.ParseVersion(s)
	if err != nil {
		return toolchain.Version{}, false
	}
	return v, true
}

// detectFromComment scans the .comment section and returns the most specific recognized toolchain.
func detectFromComment(b Binary, sd toolchain.StringDetector) toolchain.Toolchain {
	detected := make(map[toolchain.Compiler]toolchain.Version)
	for _, comment := range extractCompilerComments(b) {
		comp, ver := sd.Detect(comment)
		if comp == toolchain.Unknown {
			continue
		}
		if _, seen := detected[comp]; !seen {
			detected[comp] = ver
		}
	}
	for _, comp := range compilerPrecedence {
		if ver, ok := detected[comp]; ok {
			return toolchain.Toolchain{Compiler: comp, Version: ver}
		}
	}
	return toolchain.Toolchain{}
}

// detectFromDWARF walks compile units in .debug_info and returns the first recognized DW_AT_producer.
// Returns a zero-value Toolchain when DWARF is unavailable or no producer attribute is found.
func detectFromDWARF(b Binary, sd toolchain.StringDetector) toolchain.Toolchain {
	d, err := loadDWARF(b)
	if err != nil || d == nil {
		return toolchain.Toolchain{}
	}

	reader := d.Reader()
	for {
		entry, err := reader.Next()
		if err != nil || entry == nil {
			break
		}
		if entry.Tag != dwarf.TagCompileUnit {
			continue
		}
		producer, ok := entry.Val(dwarf.AttrProducer).(string)
		if !ok || producer == "" {
			continue
		}
		comp, ver := sd.Detect(producer)
		if comp != toolchain.Unknown {
			return toolchain.Toolchain{Compiler: comp, Version: ver}
		}
	}
	return toolchain.Toolchain{}
}

func extractCompilerComments(b Binary) []string {
	data, err := findSectionData(b, ".comment")
	if err != nil || data == nil {
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

// loadDWARF assembles a *dwarf.Data sufficient for reading DW_AT_producer.
// Fetches only the sections needed for compile-unit walks and string attributes.
// Line, ranges, and loc are skipped to avoid pulling large debug sections via the resolver.
// Returns (nil, nil) when the mandatory sections (.debug_info, .debug_abbrev) aren't available.
func loadDWARF(b Binary) (*dwarf.Data, error) {
	abbrev, err := findSectionData(b, ".debug_abbrev")
	if err != nil || len(abbrev) == 0 {
		return nil, err
	}
	info, err := findSectionData(b, ".debug_info")
	if err != nil || len(info) == 0 {
		return nil, err
	}
	str, err := findSectionData(b, ".debug_str")
	if err != nil {
		return nil, err
	}

	return dwarf.New(abbrev, nil, nil, info, nil, nil, nil, str)
}
