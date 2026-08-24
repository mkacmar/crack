package scanner

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"go.kacmar.sk/crack/internal/analyzer"
)

func newTestScanner() *Scanner {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	dispatcher := analyzer.NewDispatcher(analyzer.DispatcherOptions{
		ELF:    analyzer.NewELFAnalyzer(analyzer.ELFAnalyzerOptions{Logger: logger}),
		Logger: logger,
	})
	return NewScanner(dispatcher, Options{Logger: logger, Workers: 1})
}

func TestScanPathsReportsUnreadablePaths(t *testing.T) {
	dir := t.TempDir()
	readable := filepath.Join(dir, "readable")
	if err := os.WriteFile(readable, []byte("contents"), 0o644); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(dir, "missing")

	results := make(map[string]analyzer.FileResult)
	for res := range newTestScanner().ScanPaths(context.Background(), []string{readable, missing}, false) {
		results[res.Path] = res
	}

	res, ok := results[missing]
	if !ok {
		t.Errorf("no result for missing path %q, want one carrying an error", missing)
	} else if res.Error == nil {
		t.Errorf("result for %q: error = nil, want non-nil", missing)
	}

	if _, ok := results[readable]; !ok {
		t.Errorf("no result for %q, want it scanned alongside the missing path", readable)
	}
}
