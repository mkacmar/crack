package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mkacmar/crack/internal/output"
)

const (
	Pass Expectation = "pass"
	Fail Expectation = "fail"
	Skip Expectation = "notApplicable"
)

type Expectation string

type TestCase struct {
	Binary string
	Expect Expectation
}

func RunRuleTests(t *testing.T, rule string, cases []TestCase) {
	t.Helper()

	_, thisFile, _, _ := runtime.Caller(0)
	e2eDir := filepath.Dir(thisFile)
	rootDir := filepath.Join(e2eDir, "..", "..")
	binariesDir := filepath.Join(e2eDir, rule, "binaries")
	crackBin := filepath.Join(rootDir, "crack")

	if _, err := os.Stat(binariesDir); os.IsNotExist(err) {
		t.Skipf("binaries directory %q not found", binariesDir)
	}

	if _, err := os.Stat(crackBin); os.IsNotExist(err) {
		t.Skipf("crack binary not found, run 'make build' first")
	}

	for _, tc := range cases {
		t.Run(tc.Binary, func(t *testing.T) {
			binaryPath := filepath.Join(binariesDir, tc.Binary)
			if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
				t.Skipf("binary %q not found", tc.Binary)
			}

			sarifPath := filepath.Join(t.TempDir(), "result.sarif")
			cmd := exec.Command(
				crackBin,
				"analyze",
				"--rules="+rule,
				"--show-passed",
				"--show-skipped",
				"--sarif="+sarifPath,
				binaryPath,
			)
			cmd.Run()

			state := getRuleState(t, sarifPath, rule)
			if state != tc.Expect {
				t.Errorf("expected %s, got %s", tc.Expect, state)
			}
		})
	}
}

func getRuleState(t *testing.T, sarifPath, rule string) Expectation {
	t.Helper()

	data, err := os.ReadFile(sarifPath)
	if err != nil {
		t.Fatalf("failed to read SARIF: %v", err)
	}

	var report output.SARIFReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("failed to parse SARIF: %v", err)
	}

	if len(report.Runs) == 0 {
		t.Fatal("no runs in SARIF output")
	}

	run := report.Runs[0]
	rules := run.Tool.Driver.Rules

	for _, r := range run.Results {
		if r.RuleIndex >= 0 && r.RuleIndex < len(rules) && rules[r.RuleIndex].ID == rule {
			return Expectation(r.Kind)
		}
	}

	t.Fatalf("no result found for rule %q in SARIF output", rule)
	return ""
}
