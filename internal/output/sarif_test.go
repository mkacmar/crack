package output

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/mkacmar/crack/internal/analyzer"
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
		{name: "failed result", status: rule.StatusFailed, wantKind: "fail", wantLevel: "warning", wantIncluded: true},
		{name: "passed result included", status: rule.StatusPassed, includePassed: true, wantKind: "pass", wantLevel: "", wantIncluded: true},
		{name: "passed result excluded", status: rule.StatusPassed, includePassed: false, wantIncluded: false},
		{name: "skipped result included", status: rule.StatusSkipped, includeSkipped: true, wantKind: "notApplicable", wantLevel: "", wantIncluded: true},
		{name: "skipped result excluded", status: rule.StatusSkipped, includeSkipped: false, wantIncluded: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &analyzer.Results{
				Results: []analyzer.Result{{
					Path:      "/usr/bin/test",
					Toolchain: toolchain.Toolchain{},
					Results: []rule.ProcessedResult{{
						ExecuteResult: rule.ExecuteResult{Status: tt.status, Message: "test message"},
						RuleID:        "test-rule",
						Name:          "Test Rule",
					}},
				}},
			}

			formatter := &SARIFFormatter{IncludePassed: tt.includePassed, IncludeSkipped: tt.includeSkipped}

			var buf bytes.Buffer
			if err := formatter.Format(report, &buf); err != nil {
				t.Fatalf("Format() error = %v", err)
			}

			var sarifReport SARIFReport
			if err := json.Unmarshal(buf.Bytes(), &sarifReport); err != nil {
				t.Fatalf("failed to parse SARIF output: %v", err)
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
			if results[0].RuleIndex != 0 {
				t.Errorf("RuleIndex = %d, want 0", results[0].RuleIndex)
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

func TestSARIFInvocation(t *testing.T) {
	startTime := time.Date(2026, 1, 23, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2026, 1, 23, 10, 5, 0, 0, time.UTC)

	report := &analyzer.Results{
		Results: []analyzer.Result{{
			Path:      "/usr/bin/test",
			Toolchain: toolchain.Toolchain{},
			Results: []rule.ProcessedResult{{
				ExecuteResult: rule.ExecuteResult{Status: rule.StatusPassed, Message: "test passed"},
				RuleID:        "test-rule",
				Name:          "Test Rule",
			}},
		}},
	}

	t.Run("with invocation info", func(t *testing.T) {
		formatter := &SARIFFormatter{
			IncludePassed: true,
			Invocation: &InvocationInfo{
				CommandLine: "crack analyze --preset=recommended /usr/bin",
				Arguments:   []string{"analyze", "--preset=recommended", "/usr/bin"},
				StartTime:   startTime,
				EndTime:     endTime,
				WorkingDir:  "/home/user/project",
				Successful:  true,
			},
		}

		var buf bytes.Buffer
		if err := formatter.Format(report, &buf); err != nil {
			t.Fatalf("Format() error = %v", err)
		}

		var sarifReport SARIFReport
		if err := json.Unmarshal(buf.Bytes(), &sarifReport); err != nil {
			t.Fatalf("failed to parse SARIF output: %v", err)
		}

		invocations := sarifReport.Runs[0].Invocations
		if len(invocations) != 1 {
			t.Fatalf("expected 1 invocation, got %d", len(invocations))
		}

		inv := invocations[0]
		if inv.CommandLine != "crack analyze --preset=recommended /usr/bin" {
			t.Errorf("CommandLine = %q, want %q", inv.CommandLine, "crack analyze --preset=recommended /usr/bin")
		}
		if !inv.ExecutionSuccessful {
			t.Error("ExecutionSuccessful = false, want true")
		}
		if inv.StartTimeUtc != "2026-01-23T10:00:00Z" {
			t.Errorf("StartTimeUtc = %q, want %q", inv.StartTimeUtc, "2026-01-23T10:00:00Z")
		}
		if inv.EndTimeUtc != "2026-01-23T10:05:00Z" {
			t.Errorf("EndTimeUtc = %q, want %q", inv.EndTimeUtc, "2026-01-23T10:05:00Z")
		}
		if inv.WorkingDirectory == nil || inv.WorkingDirectory.URI != "file:///home/user/project" {
			t.Errorf("WorkingDirectory.URI = %q, want %q", inv.WorkingDirectory.URI, "file:///home/user/project")
		}
	})

	t.Run("without invocation info", func(t *testing.T) {
		formatter := &SARIFFormatter{IncludePassed: true, Invocation: nil}

		var buf bytes.Buffer
		if err := formatter.Format(report, &buf); err != nil {
			t.Fatalf("Format() error = %v", err)
		}

		var sarifReport SARIFReport
		if err := json.Unmarshal(buf.Bytes(), &sarifReport); err != nil {
			t.Fatalf("failed to parse SARIF output: %v", err)
		}

		if len(sarifReport.Runs[0].Invocations) != 0 {
			t.Errorf("expected 0 invocations, got %d", len(sarifReport.Runs[0].Invocations))
		}
	})
}
