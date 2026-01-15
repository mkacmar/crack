package rules

import (
	"strings"
	"testing"

	"github.com/mkacmar/crack/internal/model"
)

func TestBuildSuggestion(t *testing.T) {
	tests := []struct {
		name        string
		toolchain   model.Toolchain
		feature     model.FeatureAvailability
		wantContain []string
		wantExact   string
	}{
		{
			name:      "unknown compiler shows both options",
			toolchain: model.Toolchain{Compiler: model.CompilerUnknown},
			feature: model.FeatureAvailability{
				Requirements: []model.CompilerRequirement{
					{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 4, Minor: 9}, Flag: "-fstack-protector-strong"},
					{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 5}, Flag: "-fstack-protector-strong"},
				},
			},
			wantContain: []string{"Toolchain not detected", "GCC 4.9+", "Clang 3.5+"},
		},
		{
			name: "compiler below minimum version",
			toolchain: model.Toolchain{
				Compiler: model.CompilerGCC,
				Version:  model.Version{Major: 4, Minor: 8},
			},
			feature: model.FeatureAvailability{
				Requirements: []model.CompilerRequirement{
					{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 4, Minor: 9}, Flag: "-fstack-protector-strong"},
				},
			},
			wantContain: []string{"Requires gcc 4.9+", "you have gcc 4.8"},
		},
		{
			name: "compiler above min but below default version",
			toolchain: model.Toolchain{
				Compiler: model.CompilerGCC,
				Version:  model.Version{Major: 10, Minor: 0},
			},
			feature: model.FeatureAvailability{
				Requirements: []model.CompilerRequirement{
					{
						Compiler:       model.CompilerGCC,
						MinVersion:     model.Version{Major: 8, Minor: 0},
						DefaultVersion: model.Version{Major: 12, Minor: 0},
						Flag:           "-fstack-clash-protection",
					},
				},
			},
			wantContain: []string{"Use \"-fstack-clash-protection\"", "default in gcc 12.0+"},
		},
		{
			name: "compiler has no default version",
			toolchain: model.Toolchain{
				Compiler: model.CompilerClang,
				Version:  model.Version{Major: 15, Minor: 0},
			},
			feature: model.FeatureAvailability{
				Requirements: []model.CompilerRequirement{
					{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 7, Minor: 0}, Flag: "-fsanitize=safe-stack"},
				},
			},
			wantExact: "Use \"-fsanitize=safe-stack\".",
		},
		{
			name: "feature not supported by detected compiler",
			toolchain: model.Toolchain{
				Compiler: model.CompilerGCC,
				Version:  model.Version{Major: 12, Minor: 0},
			},
			feature: model.FeatureAvailability{
				Requirements: []model.CompilerRequirement{
					{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 7}, Flag: "-fsanitize=cfi"},
				},
			},
			wantContain: []string{"requires clang"},
		},
		{
			name: "compiler above default version",
			toolchain: model.Toolchain{
				Compiler: model.CompilerGCC,
				Version:  model.Version{Major: 14, Minor: 0},
			},
			feature: model.FeatureAvailability{
				Requirements: []model.CompilerRequirement{
					{
						Compiler:       model.CompilerGCC,
						MinVersion:     model.Version{Major: 8, Minor: 0},
						DefaultVersion: model.Version{Major: 12, Minor: 0},
						Flag:           "-fstack-clash-protection",
					},
				},
			},
			wantContain: []string{"Should be enabled by default"},
		},
		{
			name:      "empty requirements",
			toolchain: model.Toolchain{Compiler: model.CompilerGCC, Version: model.Version{Major: 12, Minor: 0}},
			feature:   model.FeatureAvailability{Requirements: nil},
			wantExact: "Feature not supported by detected compilers.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildSuggestion(tt.toolchain, tt.feature)

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
		feature        model.FeatureAvailability
		wantContain    []string
		wantNotContain []string
	}{
		{
			name: "only GCC requirement",
			feature: model.FeatureAvailability{
				Requirements: []model.CompilerRequirement{
					{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 7, Minor: 0}, Flag: "-mindirect-branch=thunk"},
				},
			},
			wantContain:    []string{"GCC 7.0+"},
			wantNotContain: []string{"Clang"},
		},
		{
			name: "only Clang requirement",
			feature: model.FeatureAvailability{
				Requirements: []model.CompilerRequirement{
					{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 7}, Flag: "-fsanitize=cfi"},
				},
			},
			wantContain:    []string{"Clang 3.7+"},
			wantNotContain: []string{"GCC"},
		},
		{
			name: "both compilers",
			feature: model.FeatureAvailability{
				Requirements: []model.CompilerRequirement{
					{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 4, Minor: 9}, Flag: "-fPIE"},
					{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 3, Minor: 0}, Flag: "-fPIE"},
				},
			},
			wantContain: []string{"GCC 4.9+", "Clang 3.0+", " or "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildGenericSuggestion(tt.feature)

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
