package suggestions

import (
	"strings"
	"testing"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
)

func TestBuildSuggestion(t *testing.T) {
	tests := []struct {
		name          string
		profile       binary.Profile
		applicability rule.Applicability
		wantContain   []string
		wantExact     string
	}{
		{
			name:    "unknown compiler shows both options",
			profile: binary.Profile{Toolchain: toolchain.Toolchain{Compiler: toolchain.Unknown}},
			applicability: rule.Applicability{
				Platform: binary.PlatformAll,
				Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
					toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 9}, Flag: "-fstack-protector-strong"},
					toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 5}, Flag: "-fstack-protector-strong"},
				},
			},
			wantContain: []string{"Toolchain not detected", "GCC 4.9+", "Clang 3.5+"},
		},
		{
			name: "compiler below minimum version",
			profile: binary.Profile{
				Toolchain: toolchain.Toolchain{
					Compiler: toolchain.GCC,
					Version:  toolchain.Version{Major: 4, Minor: 8},
				},
			},
			applicability: rule.Applicability{
				Platform: binary.PlatformAll,
				Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
					toolchain.GCC: {MinVersion: toolchain.Version{Major: 4, Minor: 9}, Flag: "-fstack-protector-strong"},
				},
			},
			wantContain: []string{"Requires gcc 4.9+", "you have gcc 4.8"},
		},
		{
			name: "compiler above min but below default version",
			profile: binary.Profile{
				Toolchain: toolchain.Toolchain{
					Compiler: toolchain.GCC,
					Version:  toolchain.Version{Major: 10, Minor: 0},
				},
			},
			applicability: rule.Applicability{
				Platform: binary.PlatformAll,
				Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
					toolchain.GCC: {
						MinVersion:     toolchain.Version{Major: 8, Minor: 0},
						DefaultVersion: toolchain.Version{Major: 12, Minor: 0},
						Flag:           "-fstack-clash-protection",
					},
				},
			},
			wantContain: []string{"Use \"-fstack-clash-protection\"", "default in gcc 12.0+"},
		},
		{
			name: "compiler has no default version",
			profile: binary.Profile{
				Toolchain: toolchain.Toolchain{
					Compiler: toolchain.Clang,
					Version:  toolchain.Version{Major: 15, Minor: 0},
				},
			},
			applicability: rule.Applicability{
				Platform: binary.PlatformAll,
				Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
					toolchain.Clang: {MinVersion: toolchain.Version{Major: 7, Minor: 0}, Flag: "-fsanitize=safe-stack"},
				},
			},
			wantExact: "Use \"-fsanitize=safe-stack\".",
		},
		{
			name: "feature not supported by detected compiler",
			profile: binary.Profile{
				Toolchain: toolchain.Toolchain{
					Compiler: toolchain.GCC,
					Version:  toolchain.Version{Major: 12, Minor: 0},
				},
			},
			applicability: rule.Applicability{
				Platform: binary.PlatformAll,
				Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
					toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 7}, Flag: "-fsanitize=cfi"},
				},
			},
			wantContain: []string{"requires clang"},
		},
		{
			name: "compiler above default version",
			profile: binary.Profile{
				Toolchain: toolchain.Toolchain{
					Compiler: toolchain.GCC,
					Version:  toolchain.Version{Major: 14, Minor: 0},
				},
			},
			applicability: rule.Applicability{
				Platform: binary.PlatformAll,
				Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
					toolchain.GCC: {
						MinVersion:     toolchain.Version{Major: 8, Minor: 0},
						DefaultVersion: toolchain.Version{Major: 12, Minor: 0},
						Flag:           "-fstack-clash-protection",
					},
				},
			},
			wantContain: []string{"Should be enabled by default"},
		},
		{
			name:          "empty requirements",
			profile:       binary.Profile{Toolchain: toolchain.Toolchain{Compiler: toolchain.GCC, Version: toolchain.Version{Major: 12, Minor: 0}}},
			applicability: rule.Applicability{Platform: binary.PlatformAll, Compilers: nil},
			wantExact:     "Feature not supported by detected compilers.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildSuggestion(tt.profile, tt.applicability)

			if tt.wantExact != "" {
				if result != tt.wantExact {
					t.Errorf("got %q, want %q", result, tt.wantExact)
				}
				return
			}

			for _, want := range tt.wantContain {
				if !strings.Contains(result, want) {
					t.Errorf("result %q does not contain %q", result, want)
				}
			}
		})
	}
}

func TestBuildGenericSuggestion(t *testing.T) {
	tests := []struct {
		name           string
		applicability  rule.Applicability
		wantContain    []string
		wantNotContain []string
	}{
		{
			name: "only GCC requirement",
			applicability: rule.Applicability{
				Platform: binary.PlatformAll,
				Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
					toolchain.GCC: {MinVersion: toolchain.Version{Major: 7, Minor: 0}, Flag: "-mindirect-branch=thunk"},
				},
			},
			wantContain:    []string{"GCC 7.0+"},
			wantNotContain: []string{"Clang"},
		},
		{
			name: "only Clang requirement",
			applicability: rule.Applicability{
				Platform: binary.PlatformAll,
				Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
					toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 7}, Flag: "-fsanitize=cfi"},
				},
			},
			wantContain:    []string{"Clang 3.7+"},
			wantNotContain: []string{"GCC"},
		},
		{
			name: "both compilers",
			applicability: rule.Applicability{
				Platform: binary.PlatformAll,
				Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
					toolchain.GCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 9}, Flag: "-fPIE"},
					toolchain.Clang: {MinVersion: toolchain.Version{Major: 3, Minor: 0}, Flag: "-fPIE"},
				},
			},
			wantContain: []string{"GCC 4.9+", "Clang 3.0+", " or "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildGenericSuggestion(tt.applicability)

			for _, want := range tt.wantContain {
				if !strings.Contains(result, want) {
					t.Errorf("result %q does not contain %q", result, want)
				}
			}
			for _, notWant := range tt.wantNotContain {
				if strings.Contains(result, notWant) {
					t.Errorf("result %q should not contain %q", result, notWant)
				}
			}
		})
	}
}
