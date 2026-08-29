package output

import (
	"bytes"
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"go.kacmar.sk/crack/internal/suggestions"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/rule/elf"
)

// renderSARIF streams the given results through a SARIFWriter and parses the finished document.
func renderSARIF(t *testing.T, inv *InvocationInfo, includePassed, includeSkipped bool, results ...DecoratedFileResult) SARIFReport {
	t.Helper()

	var buf bytes.Buffer
	w := NewSARIFWriter(&buf, includePassed, includeSkipped)
	for _, res := range results {
		if err := w.Write(res); err != nil {
			t.Fatalf("Write() error = %v", err)
		}
	}
	if err := w.Close(inv); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	var report SARIFReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("failed to parse SARIF output: %v\n%s", err, buf.String())
	}
	return report
}

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
			res := DecoratedFileResult{
				Path: "/usr/bin/test",
				Findings: []suggestions.DecoratedFinding{{
					Status:  tt.status,
					Message: "test message",
					RuleID:  "test-rule",
					Name:    "Test Rule",
				}},
			}

			sarifReport := renderSARIF(t, nil, tt.includePassed, tt.includeSkipped, res)
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

	res := DecoratedFileResult{
		Path: "/usr/bin/test",
		Findings: []suggestions.DecoratedFinding{{
			Status:  rule.StatusPassed,
			Message: "test passed",
			RuleID:  "test-rule",
			Name:    "Test Rule",
		}},
	}

	t.Run("with invocation info", func(t *testing.T) {
		sarifReport := renderSARIF(t, &InvocationInfo{
			CommandLine: "crack analyze --preset=recommended /usr/bin",
			Arguments:   []string{"analyze", "--preset=recommended", "/usr/bin"},
			StartTime:   startTime,
			EndTime:     endTime,
			WorkingDir:  "/home/user/project",
			Successful:  true,
		}, true, false, res)

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
		sarifReport := renderSARIF(t, nil, true, false, res)

		if len(sarifReport.Runs[0].Invocations) != 0 {
			t.Errorf("expected 0 invocations, got %d", len(sarifReport.Runs[0].Invocations))
		}
	})
}

func TestToFileURI(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{path: "/usr/bin/ls", want: "file:///usr/bin/ls"},
		{path: "/tmp/with space/x", want: "file:///tmp/with%20space/x"},
		{path: "/tmp/hash#1/x", want: "file:///tmp/hash%231/x"},
		{path: "/tmp/query?y/x", want: "file:///tmp/query%3Fy/x"},
		{path: "/tmp/100%done/x", want: "file:///tmp/100%25done/x"},
		{path: "/tmp/caf\u00e9.bin", want: "file:///tmp/caf%C3%A9.bin"},
		{path: "relative/path", want: "relative/path"},
		{path: "relative/with space", want: "relative/with%20space"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := toFileURI(tt.path)
			if got != tt.want {
				t.Errorf("toFileURI(%q) = %q, want %q", tt.path, got, tt.want)
			}
			u, err := url.Parse(got)
			if err != nil {
				t.Fatalf("toFileURI(%q) = %q, which does not parse: %v", tt.path, got, err)
			}
			if u.Path != tt.path {
				t.Errorf("toFileURI(%q) round-trips to %q", tt.path, u.Path)
			}
		})
	}
}

func TestSARIFRuleFullDescription(t *testing.T) {
	res := DecoratedFileResult{
		Path: "/usr/bin/test",
		Findings: []suggestions.DecoratedFinding{{
			Status:  rule.StatusFailed,
			Message: "Not PIE",
			RuleID:  elf.PIERuleID,
			Name:    "Position Independent Executable",
		}},
	}

	rules := renderSARIF(t, nil, false, false, res).Runs[0].Tool.Driver.Rules
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}

	got, want := rules[0].FullDescription.Text, (elf.PIERule{}).Description()
	if got != want {
		t.Errorf("fullDescription = %q, want the rule's description %q", got, want)
	}
}
