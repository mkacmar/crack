package elf

import (
	"debug/elf"
	"encoding/binary"
	"testing"

	"go.kacmar.sk/crack/toolchain"
)

func makeSection(name string, data []byte) Section {
	return Section{
		SectionHeader: elf.SectionHeader{Name: name},
		data:          func() ([]byte, error) { return data, nil },
	}
}

func makeComment(comments ...string) Section {
	var data []byte
	for _, c := range comments {
		data = append(data, c+"\x00"...)
	}
	return makeSection(".comment", data)
}

// makeGoBuildInfo builds a .go.buildinfo section in the inline (Go 1.18+) format.
func makeGoBuildInfo(version string) Section {
	data := append([]byte("\xff Go buildinf:\x08\x02"), make([]byte, 16)...)
	var varintBuf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(varintBuf[:], uint64(len(version)))
	data = append(data, varintBuf[:n]...)
	data = append(data, version+"\x00"...)
	return makeSection(".go.buildinfo", data)
}

func makeGoBuildInfoPointerFormat() Section {
	data := append([]byte("\xff Go buildinf:\x08\x00"), make([]byte, 16)...)
	return makeSection(".go.buildinfo", data)
}

func TestParseGoVersion(t *testing.T) {
	tests := []struct {
		in     string
		want   toolchain.Version
		wantOK bool
	}{
		{"go1.21.3", toolchain.Version{Major: 1, Minor: 21, Patch: 3}, true},
		{"go1.21", toolchain.Version{Major: 1, Minor: 21}, true},
		{"go1.22-rc1", toolchain.Version{Major: 1, Minor: 22}, true},
		{"go1.22+something", toolchain.Version{Major: 1, Minor: 22}, true},
		{"go1.22 someinfo", toolchain.Version{Major: 1, Minor: 22}, true},
		{"devel go1.22", toolchain.Version{}, false},
		{"go1", toolchain.Version{}, false},
		{"", toolchain.Version{}, false},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, ok := parseGoVersion(tc.in)
			if ok != tc.wantOK {
				t.Fatalf("parseGoVersion(%q) ok = %v, want %v", tc.in, ok, tc.wantOK)
			}
			if got != tc.want {
				t.Errorf("parseGoVersion(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestDefaultToolchainDetector(t *testing.T) {
	tests := []struct {
		name     string
		sections []Section
		want     toolchain.Toolchain
	}{
		{
			name: "empty binary returns unknown",
			want: toolchain.Toolchain{},
		},
		{
			name:     "go inline buildinfo yields version",
			sections: []Section{makeGoBuildInfo("go1.21.3")},
			want:     toolchain.Toolchain{Compiler: toolchain.Go, Version: toolchain.Version{Major: 1, Minor: 21, Patch: 3}},
		},
		{
			name:     "go pointer-format buildinfo yields compiler without version",
			sections: []Section{makeGoBuildInfoPointerFormat()},
			want:     toolchain.Toolchain{Compiler: toolchain.Go},
		},
		{
			name:     "go buildinfo with unparseable version still reports go",
			sections: []Section{makeGoBuildInfo("devel go1.22-abc")},
			want:     toolchain.Toolchain{Compiler: toolchain.Go},
		},
		{
			name:     "go buildinfo precedes .comment",
			sections: []Section{makeGoBuildInfo("go1.21"), makeComment("GCC: (Debian 12.2.0-14) 12.2.0")},
			want:     toolchain.Toolchain{Compiler: toolchain.Go, Version: toolchain.Version{Major: 1, Minor: 21}},
		},
		{
			name:     "gcc detected from .comment",
			sections: []Section{makeComment("GCC: (Debian 12.2.0-14) 12.2.0")},
			want:     toolchain.Toolchain{Compiler: toolchain.GCC, Version: toolchain.Version{Major: 12, Minor: 2}},
		},
		{
			name:     "clang detected from .comment",
			sections: []Section{makeComment("Debian clang version 14.0.6")},
			want:     toolchain.Toolchain{Compiler: toolchain.Clang, Version: toolchain.Version{Major: 14, Minor: 0, Patch: 6}},
		},
		{
			name:     "rustc detected from .comment",
			sections: []Section{makeComment("rustc version 1.70.0")},
			want:     toolchain.Toolchain{Compiler: toolchain.Rustc, Version: toolchain.Version{Major: 1, Minor: 70}},
		},
		{
			name:     "precedence: rustc wins over gcc in mixed .comment",
			sections: []Section{makeComment("GCC: (Debian 12.2.0-14) 12.2.0", "rustc version 1.70.0")},
			want:     toolchain.Toolchain{Compiler: toolchain.Rustc, Version: toolchain.Version{Major: 1, Minor: 70}},
		},
		{
			name:     "precedence: clang wins over gcc in mixed .comment",
			sections: []Section{makeComment("GCC: (Debian 12.2.0-14) 12.2.0", "Debian clang version 14.0.6")},
			want:     toolchain.Toolchain{Compiler: toolchain.Clang, Version: toolchain.Version{Major: 14, Minor: 0, Patch: 6}},
		},
		{
			name:     "unrecognized .comment yields unknown",
			sections: []Section{makeComment("some random toolchain")},
			want:     toolchain.Toolchain{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fb := &fakeBinary{sections: tc.sections}
			got := DefaultToolchainDetector{}.Detect(fb)
			if got != tc.want {
				t.Errorf("Detect() = %+v, want %+v", got, tc.want)
			}
		})
	}
}

type stubStringDetector struct {
	called bool
	comp   toolchain.Compiler
	ver    toolchain.Version
}

func (s *stubStringDetector) Detect(string) (toolchain.Compiler, toolchain.Version) {
	s.called = true
	return s.comp, s.ver
}

func TestDefaultToolchainDetectorStringDetectorOverride(t *testing.T) {
	stub := &stubStringDetector{comp: toolchain.GCC, ver: toolchain.Version{Major: 99, Minor: 0}}
	fb := &fakeBinary{sections: []Section{makeComment("anything")}}

	got := DefaultToolchainDetector{StringDetector: stub}.Detect(fb)

	if !stub.called {
		t.Fatal("override StringDetector was not consulted")
	}
	want := toolchain.Toolchain{Compiler: toolchain.GCC, Version: toolchain.Version{Major: 99}}
	if got != want {
		t.Errorf("Detect() = %+v, want %+v", got, want)
	}
}
