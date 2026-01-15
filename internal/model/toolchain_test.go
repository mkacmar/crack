package model

import (
	"testing"
)

func TestParseToolchain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		compiler Compiler
		major    int
		minor    int
	}{
		{
			name:     "empty string",
			input:    "",
			compiler: CompilerUnknown,
		},
		{
			name:     "GCC simple",
			input:    "GCC: (GNU) 11.2.0",
			compiler: CompilerGCC,
			major:    11,
			minor:    2,
		},
		{
			name:     "GCC with GNU C",
			input:    "GNU C17 11.4.0 -mtune=generic -march=x86-64",
			compiler: CompilerGCC,
			major:    11,
			minor:    4,
		},
		{
			name:     "GCC GIMPLE",
			input:    "GNU GIMPLE 12.1.0",
			compiler: CompilerGCC,
			major:    12,
			minor:    1,
		},
		{
			name:     "Clang simple",
			input:    "clang version 14.0.0",
			compiler: CompilerClang,
			major:    14,
			minor:    0,
		},
		{
			name:     "Clang with target",
			input:    "Ubuntu clang version 15.0.7 (target: x86_64-pc-linux-gnu)",
			compiler: CompilerClang,
			major:    15,
			minor:    0,
		},
		{
			name:     "Unknown compiler",
			input:    "some random string",
			compiler: CompilerUnknown,
		},
		{
			name:     "GCC from DWARF producer",
			input:    "GNU C23 15.2.1 20250813 -D_FORTIFY_SOURCE=3 -march=x86-64",
			compiler: CompilerGCC,
			major:    15,
			minor:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := ParseToolchain(tt.input)

			if tc.Compiler != tt.compiler {
				t.Errorf("ParseToolchain(%q).Compiler = %v, want %v", tt.input, tc.Compiler, tt.compiler)
			}

			if tt.compiler != CompilerUnknown {
				if tc.Version.Major != tt.major {
					t.Errorf("ParseToolchain(%q).Version.Major = %d, want %d", tt.input, tc.Version.Major, tt.major)
				}
				if tc.Version.Minor != tt.minor {
					t.Errorf("ParseToolchain(%q).Version.Minor = %d, want %d", tt.input, tc.Version.Minor, tt.minor)
				}
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		major int
		minor int
		patch int
	}{
		{
			name:  "version keyword",
			input: "clang version 14.0.6",
			major: 14,
			minor: 0,
			patch: 6,
		},
		{
			name:  "version without keyword",
			input: "GCC: (GNU) 11.2.0",
			major: 11,
			minor: 2,
			patch: 0,
		},
		{
			name:  "version with parentheses",
			input: "Ubuntu clang version 15.0.7 (something)",
			major: 15,
			minor: 0,
			patch: 7,
		},
		{
			name:  "no version",
			input: "unknown toolchain",
			major: 0,
			minor: 0,
			patch: 0,
		},
		{
			name:  "two digit minor",
			input: "GCC 4.12.0",
			major: 4,
			minor: 12,
			patch: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := parseVersion(tt.input)

			if v.Major != tt.major {
				t.Errorf("parseVersion(%q).Major = %d, want %d", tt.input, v.Major, tt.major)
			}
			if v.Minor != tt.minor {
				t.Errorf("parseVersion(%q).Minor = %d, want %d", tt.input, v.Minor, tt.minor)
			}
			if v.Patch != tt.patch {
				t.Errorf("parseVersion(%q).Patch = %d, want %d", tt.input, v.Patch, tt.patch)
			}
		})
	}
}

func TestTryParseVersion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		ok    bool
		major int
		minor int
		patch int
	}{
		{
			name:  "major.minor",
			input: "11.2",
			ok:    true,
			major: 11,
			minor: 2,
			patch: 0,
		},
		{
			name:  "major.minor.patch",
			input: "14.0.6",
			ok:    true,
			major: 14,
			minor: 0,
			patch: 6,
		},
		{
			name:  "single number",
			input: "14",
			ok:    false,
		},
		{
			name:  "invalid major",
			input: "abc.2.0",
			ok:    false,
		},
		{
			name:  "invalid minor",
			input: "14.abc.0",
			ok:    false,
		},
		{
			name:  "empty string",
			input: "",
			ok:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, ok := tryParseVersion(tt.input)

			if ok != tt.ok {
				t.Errorf("tryParseVersion(%q) ok = %v, want %v", tt.input, ok, tt.ok)
			}

			if tt.ok {
				if v.Major != tt.major {
					t.Errorf("tryParseVersion(%q).Major = %d, want %d", tt.input, v.Major, tt.major)
				}
				if v.Minor != tt.minor {
					t.Errorf("tryParseVersion(%q).Minor = %d, want %d", tt.input, v.Minor, tt.minor)
				}
				if v.Patch != tt.patch {
					t.Errorf("tryParseVersion(%q).Patch = %d, want %d", tt.input, v.Patch, tt.patch)
				}
			}
		})
	}
}
