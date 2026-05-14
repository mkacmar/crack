package toolchain

import (
	"errors"
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		in      string
		want    Version
		wantErr error
	}{
		{"1.2.3", Version{Major: 1, Minor: 2, Patch: 3}, nil},
		{"1.21", Version{Major: 1, Minor: 21}, nil},
		{"12.2.0", Version{Major: 12, Minor: 2}, nil},
		{"1.2.3.4", Version{Major: 1, Minor: 2, Patch: 3}, nil},
		{"1.2.x", Version{Major: 1, Minor: 2}, nil},
		{"1", Version{}, ErrInvalidVersionFormat},
		{"", Version{}, ErrInvalidVersionFormat},
		{"x.2.3", Version{}, ErrInvalidVersionMajor},
		{"1.x.3", Version{}, ErrInvalidVersionMinor},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, err := ParseVersion(tc.in)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("ParseVersion(%q) err = %v, want %v", tc.in, err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseVersion(%q): %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("ParseVersion(%q) = %+v, want %+v", tc.in, got, tc.want)
			}
		})
	}
}

func TestVersionString(t *testing.T) {
	tests := []struct {
		v    Version
		want string
	}{
		{Version{Major: 1, Minor: 2, Patch: 3}, "1.2.3"},
		{Version{Major: 1, Minor: 21}, "1.21"},
		{Version{Major: 0, Minor: 0}, "0.0"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			if got := tc.v.String(); got != tc.want {
				t.Errorf("%+v.String() = %q, want %q", tc.v, got, tc.want)
			}
		})
	}
}

func TestVersionIsAtLeast(t *testing.T) {
	tests := []struct {
		name     string
		have     Version
		required Version
		want     bool
	}{
		{"equal", Version{1, 2, 3}, Version{1, 2, 3}, true},
		{"higher patch", Version{1, 2, 4}, Version{1, 2, 3}, true},
		{"lower patch", Version{1, 2, 2}, Version{1, 2, 3}, false},
		{"higher minor beats lower patch", Version{1, 3, 0}, Version{1, 2, 9}, true},
		{"higher major beats lower minor/patch", Version{2, 0, 0}, Version{1, 9, 9}, true},
		{"lower major loses regardless of minor", Version{1, 9, 9}, Version{2, 0, 0}, false},
		{"zero satisfies zero", Version{}, Version{}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.have.IsAtLeast(tc.required); got != tc.want {
				t.Errorf("%+v.IsAtLeast(%+v) = %v, want %v", tc.have, tc.required, got, tc.want)
			}
		})
	}
}

func TestDefaultStringDetector(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		wantComp Compiler
		wantVer  Version
	}{
		{"gcc with version keyword", "GCC: (Debian 12.2.0-14) 12.2.0", GCC, Version{Major: 12, Minor: 2}},
		{"clang with version keyword", "Debian clang version 14.0.6", Clang, Version{Major: 14, Minor: 0, Patch: 6}},
		{"rustc", "rustc version 1.70.0 (90c541806 2023-05-31)", Rustc, Version{Major: 1, Minor: 70}},
		{"gnu c", "GNU C17 12.2.0 -mtune=generic", GCC, Version{Major: 12, Minor: 2}},
		{"gnu gimple", "GNU GIMPLE 12.2.0", GCC, Version{Major: 12, Minor: 2}},
		{"no version still classifies", "GCC: (Debian)", GCC, Version{}},
		{"empty string is unknown", "", Unknown, Version{}},
		{"unrelated string is unknown", "some random producer", Unknown, Version{}},
		{"rustc takes precedence over gcc keyword", "rustc 1.70.0 built with gcc 12", Rustc, Version{Major: 1, Minor: 70}},
		{"clang version without keyword falls back to dotted token", "Apple clang 15.0.0", Clang, Version{Major: 15, Minor: 0}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			comp, ver := DefaultStringDetector{}.Detect(tc.in)
			if comp != tc.wantComp {
				t.Errorf("compiler = %v, want %v", comp, tc.wantComp)
			}
			if ver != tc.wantVer {
				t.Errorf("version = %+v, want %+v", ver, tc.wantVer)
			}
		})
	}
}
