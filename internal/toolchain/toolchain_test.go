package toolchain

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input   string
		want    Version
		wantErr bool
	}{
		{"14.2", Version{Major: 14, Minor: 2}, false},
		{"14.2.1", Version{Major: 14, Minor: 2, Patch: 1}, false},
		{"3.5.0", Version{Major: 3, Minor: 5, Patch: 0}, false},
		{"14", Version{}, true},
		{"", Version{}, true},
		{"abc.def", Version{}, true},
		{"14.abc", Version{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVersion(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseVersion(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestVersionIsAtLeast(t *testing.T) {
	tests := []struct {
		name     string
		v        Version
		required Version
		want     bool
	}{
		{"equal versions", Version{14, 2, 1}, Version{14, 2, 1}, true},
		{"higher major", Version{15, 0, 0}, Version{14, 2, 1}, true},
		{"lower major", Version{13, 9, 9}, Version{14, 0, 0}, false},
		{"higher minor", Version{14, 3, 0}, Version{14, 2, 1}, true},
		{"lower minor", Version{14, 1, 9}, Version{14, 2, 0}, false},
		{"higher patch", Version{14, 2, 2}, Version{14, 2, 1}, true},
		{"lower patch", Version{14, 2, 0}, Version{14, 2, 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.IsAtLeast(tt.required); got != tt.want {
				t.Errorf("%v.IsAtLeast(%v) = %v, want %v", tt.v, tt.required, got, tt.want)
			}
		})
	}
}

// TestParseToolchain tests parsing of compiler info from .comment section.
// To extract .comment from a binary: readelf -p .comment <binary>
func TestParseToolchain(t *testing.T) {
	tests := []struct {
		input        string
		wantCompiler Compiler
		wantVersion  Version
	}{
		{"", CompilerUnknown, Version{}},
		{"GCC: (Ubuntu 14.2.0-4ubuntu2~24.04) 14.2.0", CompilerGCC, Version{14, 2, 0}},
		{"Ubuntu clang version 18.1.3 (1ubuntu1)", CompilerClang, Version{18, 1, 3}},
		{"GCC: (Alpine 15.2.0) 15.2.0", CompilerGCC, Version{15, 2, 0}},
		{"GCC: (GNU) 15.2.1 20251211 (Red Hat 15.2.1-5)", CompilerGCC, Version{15, 2, 1}},
		{"Alpine clang version 21.1.2", CompilerClang, Version{21, 1, 2}},
		{"GNU GIMPLE 14.2.0", CompilerGCC, Version{14, 2, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseToolchain(tt.input)
			if got.Compiler != tt.wantCompiler {
				t.Errorf("ParseToolchain(%q).Compiler = %v, want %v", tt.input, got.Compiler, tt.wantCompiler)
			}
			if got.Version != tt.wantVersion {
				t.Errorf("ParseToolchain(%q).Version = %v, want %v", tt.input, got.Version, tt.wantVersion)
			}
		})
	}
}
