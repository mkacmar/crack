package cli

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
)

func TestParsePlatformTarget(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    rule.PlatformTarget
		wantErr string
	}{
		{
			name: "bare architecture",
			in:   "amd64",
			want: rule.PlatformTarget{Architecture: binary.ArchAMD64},
		},
		{
			name: "architecture with ISA",
			in:   "arm64:8.5",
			want: rule.PlatformTarget{Architecture: binary.ArchARM64, MaxISA: &binary.ISA{Major: 8, Minor: 5}},
		},
		{
			name: "architecture with v-prefixed ISA",
			in:   "arm64:v9",
			want: rule.PlatformTarget{Architecture: binary.ArchARM64, MaxISA: &binary.ISA{Major: 9}},
		},
		{
			name:    "unknown architecture lists the valid ones",
			in:      "sparc",
			wantErr: `unknown architecture "sparc"`,
		},
		{
			name:    "malformed ISA",
			in:      "arm64:8.5.1",
			wantErr: `invalid ISA version "8.5.1"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePlatformTarget(tt.in)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("parsePlatformTarget(%q): %v", tt.in, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePlatformTarget(%q) = %+v, want %+v", tt.in, got, tt.want)
			}
		})
	}
}

func TestParsePlatformTargetErrorListsValidNames(t *testing.T) {
	_, err := parsePlatformTarget("sparc")
	if err == nil {
		t.Fatal("parsePlatformTarget(\"sparc\"): expected an error")
	}
	for _, name := range validArchitectureNames() {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("error %q does not offer %q as a valid value", err, name)
		}
	}
}

func TestParseCompilerTarget(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    rule.CompilerTarget
		wantErr string
	}{
		{
			name: "bare compiler",
			in:   "gcc",
			want: rule.CompilerTarget{Compiler: toolchain.GCC},
		},
		{
			name: "compiler with version",
			in:   "clang:18.1",
			want: rule.CompilerTarget{Compiler: toolchain.Clang, MaxVersion: &toolchain.Version{Major: 18, Minor: 1}},
		},
		{
			name: "rustc",
			in:   "rustc",
			want: rule.CompilerTarget{Compiler: toolchain.Rustc},
		},
		{
			name:    "unknown compiler",
			in:      "icc",
			wantErr: `unknown compiler "icc"`,
		},
		{
			name:    "version without a minor is rejected",
			in:      "gcc:13",
			wantErr: `invalid compiler version "13"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCompilerTarget(tt.in)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseCompilerTarget(%q): %v", tt.in, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCompilerTarget(%q) = %+v, want %+v", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseCompilerTargetWrapsVersionError(t *testing.T) {
	_, err := parseCompilerTarget("gcc:13")
	if !errors.Is(err, toolchain.ErrInvalidVersionFormat) {
		t.Errorf("error = %v, want it to wrap %v", err, toolchain.ErrInvalidVersionFormat)
	}
}

func TestParseList(t *testing.T) {
	identity := func(s string) (string, error) { return s, nil }

	tests := []struct {
		name string
		in   string
		want []string
	}{
		{name: "empty input yields no items", in: "", want: nil},
		{name: "single item", in: "a", want: []string{"a"}},
		{name: "surrounding whitespace is trimmed", in: " a , b ", want: []string{"a", "b"}},
		{name: "empty items are skipped", in: "a,,b", want: []string{"a", "b"}},
		{name: "only separators yields no items", in: ",,,", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseList(tt.in, identity)
			if err != nil {
				t.Fatalf("parseList(%q): %v", tt.in, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseList(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseListStopsAtFirstError(t *testing.T) {
	sentinel := errors.New("bad item")
	var seen []string
	parse := func(s string) (string, error) {
		seen = append(seen, s)
		if s == "b" {
			return "", sentinel
		}
		return s, nil
	}

	got, err := parseList("a,b,c", parse)
	if !errors.Is(err, sentinel) {
		t.Errorf("error = %v, want %v", err, sentinel)
	}
	if got != nil {
		t.Errorf("result = %v, want nil on error", got)
	}
	if !reflect.DeepEqual(seen, []string{"a", "b"}) {
		t.Errorf("parsed %v, want it to stop after \"b\"", seen)
	}
}

func TestParseTargetFilter(t *testing.T) {
	tests := []struct {
		name          string
		platforms     string
		compilers     string
		wantPlatforms int
		wantCompilers int
		wantErr       string
	}{
		{name: "both empty"},
		{name: "platforms only", platforms: "amd64,arm64", wantPlatforms: 2},
		{name: "compilers only", compilers: "gcc,clang", wantCompilers: 2},
		{name: "both", platforms: "amd64", compilers: "gcc:13.2", wantPlatforms: 1, wantCompilers: 1},
		{name: "platform error", platforms: "sparc", wantErr: "unknown architecture"},
		{name: "compiler error", compilers: "icc", wantErr: "unknown compiler"},
		{name: "platform error wins over compiler error", platforms: "sparc", compilers: "icc", wantErr: "unknown architecture"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTargetFilter(tt.platforms, tt.compilers)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %v, want it to contain %q", err, tt.wantErr)
				}
				if got != nil {
					t.Errorf("filter = %+v, want nil on error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseTargetFilter: %v", err)
			}
			if len(got.Platforms) != tt.wantPlatforms {
				t.Errorf("platforms = %d, want %d", len(got.Platforms), tt.wantPlatforms)
			}
			if len(got.Compilers) != tt.wantCompilers {
				t.Errorf("compilers = %d, want %d", len(got.Compilers), tt.wantCompilers)
			}
		})
	}
}
