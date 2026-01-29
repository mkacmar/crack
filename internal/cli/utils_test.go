package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseURLList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single URL",
			input:    "https://example.com",
			expected: []string{"https://example.com"},
		},
		{
			name:     "multiple URLs",
			input:    "https://a.com,https://b.com,https://c.com",
			expected: []string{"https://a.com", "https://b.com", "https://c.com"},
		},
		{
			name:     "URLs with spaces",
			input:    "https://a.com , https://b.com , https://c.com",
			expected: []string{"https://a.com", "https://b.com", "https://c.com"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "only commas",
			input:    ",,,",
			expected: nil,
		},
		{
			name:     "trailing comma",
			input:    "https://a.com,",
			expected: []string{"https://a.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseURLList(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("parseURLList(%q) = %v, want %v", tt.input, result, tt.expected)
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("parseURLList(%q)[%d] = %q, want %q", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestReadPathsFromInput(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		expected  []string
		wantError bool
	}{
		{
			name:     "single path",
			content:  "/path/to/binary\n",
			expected: []string{"/path/to/binary"},
		},
		{
			name:     "multiple paths",
			content:  "/path/to/binary1\n/path/to/binary2\n/path/to/binary3\n",
			expected: []string{"/path/to/binary1", "/path/to/binary2", "/path/to/binary3"},
		},
		{
			name:     "skip empty lines",
			content:  "/path/to/binary1\n\n\n/path/to/binary2\n\n",
			expected: []string{"/path/to/binary1", "/path/to/binary2"},
		},
		{
			name:     "trim whitespace",
			content:  "  /path/to/binary1  \n\t/path/to/binary2\t\n",
			expected: []string{"/path/to/binary1", "/path/to/binary2"},
		},
		{
			name:     "empty file",
			content:  "",
			expected: nil,
		},
		{
			name:     "only whitespace",
			content:  "   \n\t\n  \n",
			expected: nil,
		},
		{
			name:     "no trailing newline",
			content:  "/path/to/binary1\n/path/to/binary2",
			expected: []string{"/path/to/binary1", "/path/to/binary2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			inputFile := filepath.Join(tmpDir, "paths.txt")
			if err := os.WriteFile(inputFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			paths, err := readPathsFromInput(inputFile)
			if err != nil {
				t.Fatalf("readPathsFromInput() error = %v", err)
			}

			if len(paths) != len(tt.expected) {
				t.Errorf("got %d paths, want %d", len(paths), len(tt.expected))
				return
			}
			for i := range paths {
				if paths[i] != tt.expected[i] {
					t.Errorf("paths[%d] = %q, want %q", i, paths[i], tt.expected[i])
				}
			}
		})
	}

	t.Run("file not found", func(t *testing.T) {
		_, err := readPathsFromInput("/nonexistent/file.txt")
		if err == nil {
			t.Error("expected error for nonexistent file")
		}
	})
}
