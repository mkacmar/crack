package scanner

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.kacmar.sk/crack/internal/analyzer"
)

func newTestScanner(workers int) *Scanner {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	dispatcher := analyzer.NewDispatcher(analyzer.DispatcherOptions{
		ELF:    analyzer.NewELFAnalyzer(analyzer.ELFAnalyzerOptions{Logger: logger}),
		Logger: logger,
	})
	return NewScanner(dispatcher, Options{Logger: logger, Workers: workers})
}

// writeFiles creates n files under dir and returns their paths.
func writeFiles(t *testing.T, dir string, n int) []string {
	t.Helper()
	paths := make([]string, 0, n)
	for i := range n {
		p := filepath.Join(dir, fmt.Sprintf("file%02d", i))
		if err := os.WriteFile(p, []byte(strings.Repeat("not an ELF binary\n", 8)), 0o600); err != nil {
			t.Fatal(err)
		}
		paths = append(paths, p)
	}
	return paths
}

const scanFileCount = 64

// 1 is serial, 4 makes workers queue for a slot, scanFileCount never blocks.
var workerCounts = []int{1, 4, scanFileCount}

func collectPaths(results <-chan analyzer.FileResult) map[string]int {
	seen := make(map[string]int)
	for res := range results {
		seen[res.Path]++
	}
	return seen
}

func TestScanPathsReportsUnreadablePaths(t *testing.T) {
	dir := t.TempDir()
	readable := filepath.Join(dir, "readable")
	if err := os.WriteFile(readable, []byte("contents"), 0o644); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(dir, "missing")

	results := make(map[string]analyzer.FileResult)
	for res := range newTestScanner(1).ScanPaths(context.Background(), []string{readable, missing}, false) {
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

func TestScanPathsReturnsEveryFileExactlyOnce(t *testing.T) {
	for _, workers := range workerCounts {
		t.Run(fmt.Sprintf("workers=%d", workers), func(t *testing.T) {
			paths := writeFiles(t, t.TempDir(), scanFileCount)

			seen := collectPaths(newTestScanner(workers).ScanPaths(context.Background(), paths, false))

			if len(seen) != len(paths) {
				t.Errorf("got %d distinct results, want %d", len(seen), len(paths))
			}
			for _, p := range paths {
				switch seen[p] {
				case 1:
				case 0:
					t.Errorf("no result for %q", p)
				default:
					t.Errorf("%q reported %d times, want 1", p, seen[p])
				}
			}
		})
	}
}

func TestScanPathsStopsWhenContextIsCancelled(t *testing.T) {
	for _, workers := range workerCounts {
		t.Run(fmt.Sprintf("workers=%d", workers), func(t *testing.T) {
			paths := writeFiles(t, t.TempDir(), 256)

			ctx, cancel := context.WithCancel(context.Background())
			results := newTestScanner(workers).ScanPaths(ctx, paths, false)

			<-results
			cancel()

			// The channel is unbuffered, so the remaining workers sit parked mid-send until cancellation releases them.
			total := 1
			drained := make(chan struct{})
			go func() {
				defer close(drained)
				for range results {
					total++
				}
			}()

			select {
			case <-drained:
			case <-time.After(10 * time.Second):
				t.Fatal("ScanPaths did not close its result channel after cancellation")
			}

			if total >= len(paths) {
				t.Errorf("scanned %d of %d files, want cancellation to cut the scan short", total, len(paths))
			}
		})
	}
}

func TestScanPathsDirectoryWalk(t *testing.T) {
	tests := []struct {
		name      string
		recursive bool
		wantDeep  bool
	}{
		{name: "recursive descends into subdirectories", recursive: true, wantDeep: true},
		{name: "non-recursive stays at the top level", recursive: false, wantDeep: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			nested := filepath.Join(root, "a", "b")
			if err := os.MkdirAll(nested, 0o750); err != nil {
				t.Fatal(err)
			}
			top := writeFiles(t, root, 2)
			deep := writeFiles(t, nested, 3)

			// Collection finishes before any worker starts, so the fan-out cannot change what is found.
			seen := collectPaths(newTestScanner(1).ScanPaths(context.Background(), []string{root}, tt.recursive))

			for _, p := range top {
				if seen[p] != 1 {
					t.Errorf("%q seen %d times, want 1", p, seen[p])
				}
			}

			wantDeepCount := 0
			if tt.wantDeep {
				wantDeepCount = 1
			}
			for _, p := range deep {
				if seen[p] != wantDeepCount {
					t.Errorf("nested %q seen %d times, want %d", p, seen[p], wantDeepCount)
				}
			}
		})
	}
}

func TestNewScannerRejectsNonPositiveWorkers(t *testing.T) {
	for _, workers := range []int{0, -1} {
		t.Run(fmt.Sprintf("workers=%d", workers), func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf("NewScanner with Workers=%d: expected a panic", workers)
				}
			}()
			newTestScanner(workers)
		})
	}
}
