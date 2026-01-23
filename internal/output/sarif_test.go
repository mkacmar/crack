package output

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/mkacmar/crack/internal/result"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

func TestSARIFResultKind(t *testing.T) {
	tests := []struct {
		name           string
		status         rule.Status
		includePassed  bool
		includeSkipped bool
		wantKind       string
		wantLevel      string
		wantIncluded   bool
	}{
		{
			name:         "failed result",
			status:       rule.StatusFailed,
			wantKind:     "fail",
			wantLevel:    "warning",
			wantIncluded: true,
		},
		{
			name:          "passed result included",
			status:        rule.StatusPassed,
			includePassed: true,
			wantKind:      "pass",
			wantLevel:     "",
			wantIncluded:  true,
		},
		{
			name:          "passed result excluded",
			status:        rule.StatusPassed,
			includePassed: false,
			wantIncluded:  false,
		},
		{
			name:           "skipped result included",
			status:         rule.StatusSkipped,
			includeSkipped: true,
			wantKind:       "notApplicable",
			wantLevel:      "",
			wantIncluded:   true,
		},
		{
			name:           "skipped result excluded",
			status:         rule.StatusSkipped,
			includeSkipped: false,
			wantIncluded:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &result.ScanResults{
				Results: []result.FileScanResult{
					{
						Path:      "/usr/bin/test",
						Toolchain: toolchain.Toolchain{},
						Results: []rule.ProcessedResult{
							{
								ExecuteResult: rule.ExecuteResult{
									Status:  tt.status,
									Message: "test message",
								},
								RuleID: "test-rule",
								Name:   "Test Rule",
							},
						},
					},
				},
			}

			formatter := &SARIFFormatter{
				IncludePassed:  tt.includePassed,
				IncludeSkipped: tt.includeSkipped,
			}

			var buf bytes.Buffer
			if err := formatter.Format(report, &buf); err != nil {
				t.Fatalf("Format() error = %v", err)
			}

			var sarifReport SARIFReport
			if err := json.Unmarshal(buf.Bytes(), &sarifReport); err != nil {
				t.Fatalf("failed to parse SARIF output: %v", err)
			}

			if len(sarifReport.Runs) != 1 {
				t.Fatalf("expected 1 run, got %d", len(sarifReport.Runs))
			}

			results := sarifReport.Runs[0].Results

			if !tt.wantIncluded {
				if len(results) != 0 {
					t.Errorf("expected 0 results, got %d", len(results))
				}
				return
			}

			if len(results) != 1 {
				t.Fatalf("expected 1 result, got %d", len(results))
			}

			if results[0].Kind != tt.wantKind {
				t.Errorf("Kind = %q, want %q", results[0].Kind, tt.wantKind)
			}

			if results[0].Level != tt.wantLevel {
				t.Errorf("Level = %q, want %q", results[0].Level, tt.wantLevel)
			}
		})
	}
}
