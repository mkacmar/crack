package scanner

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/internal/analyzer"
)

type Scanner struct {
	dispatcher *analyzer.Dispatcher
	logger     *slog.Logger
	workers    int
}

type Options struct {
	Logger  *slog.Logger
	Workers int
}

func NewScanner(dispatcher *analyzer.Dispatcher, opts Options) *Scanner {
	return &Scanner{
		dispatcher: dispatcher,
		logger:     opts.Logger.With(slog.String("component", "scanner")),
		workers:    opts.Workers,
	}
}

func (s *Scanner) ScanPaths(ctx context.Context, paths []string, recursive bool) <-chan analyzer.FileResult {
	var filesToScan []string
	var failures []analyzer.FileResult

	for _, path := range paths {
		files, failed := s.collectFiles(path, recursive)
		for _, f := range failed {
			s.logger.Warn("failed to collect files", slog.String("path", f.Path), slog.Any("error", f.Error))
		}
		filesToScan = append(filesToScan, files...)
		failures = append(failures, failed...)
	}

	s.logger.Debug("collected files to scan", slog.Int("count", len(filesToScan)), slog.Int("failures", len(failures)))

	return s.scanFiles(ctx, filesToScan, failures)
}

func (s *Scanner) scanFiles(ctx context.Context, files []string, failures []analyzer.FileResult) <-chan analyzer.FileResult {
	results := make(chan analyzer.FileResult)

	if len(files) == 0 && len(failures) == 0 {
		close(results)
		return results
	}

	s.logger.Debug("starting parallel scan", slog.Int("workers", s.workers), slog.Int("files", len(files)))

	go func() {
		defer close(results)

		for _, res := range failures {
			select {
			case results <- res:
			case <-ctx.Done():
				return
			}
		}

		g, ctx := errgroup.WithContext(ctx)
		g.SetLimit(s.workers)

		for _, path := range files {
			g.Go(func() error {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				fileResults := s.scanFile(ctx, path)
				for _, res := range fileResults {
					select {
					case results <- res:
					case <-ctx.Done():
						return ctx.Err()
					}
				}
				return nil
			})
		}

		_ = g.Wait()
	}()

	return results
}

// scanFile returns a slice of FileResult to support future fat/universal binaries.
func (s *Scanner) scanFile(ctx context.Context, path string) []analyzer.FileResult {
	s.logger.Debug("scanning file", slog.String("path", path))

	f, err := os.Open(path) // #nosec G304 -- user-provided paths are the tool's input
	if err != nil {
		s.logger.Warn("failed to open file", slog.String("path", path), slog.Any("error", err))
		return []analyzer.FileResult{{Path: path, Error: err}}
	}
	defer f.Close()

	results, err := s.dispatcher.Analyze(ctx, f)
	if err != nil {
		if errors.Is(err, analyzer.ErrUnrecognizedFormat) {
			s.logger.Debug("skipping unsupported format", slog.String("path", path))
			return []analyzer.FileResult{{Path: path, Skipped: true}}
		}
		s.logger.Warn("failed to analyze file", slog.String("path", path), slog.Any("error", err))
		return []analyzer.FileResult{{Path: path, Error: err}}
	}

	hash, err := hashFile(f)
	if err != nil {
		s.logger.Warn("failed to compute SHA256", slog.String("path", path), slog.Any("error", err))
	}

	// Assemble FileResults from AnalysisResults
	fileResults := make([]analyzer.FileResult, len(results))
	for i, r := range results {
		fileResults[i] = analyzer.FileResult{
			Path:     path,
			Format:   r.Format,
			Identity: binary.Identity{BuildID: r.Identity.BuildID, SHA256: hash},
			Profile:  r.Profile,
			Findings: r.Findings,
		}
	}
	return fileResults
}

// collectFiles returns the files found under path plus a result for every entry that could not be read.
func (s *Scanner) collectFiles(path string, recursive bool) ([]string, []analyzer.FileResult) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, []analyzer.FileResult{{Path: path, Error: fmt.Errorf("failed to stat path: %w", err)}}
	}

	if !info.IsDir() {
		return []string{path}, nil
	}

	var files []string
	var failures []analyzer.FileResult

	if recursive {
		// Reported and skipped: returning walkErr would abort the walk and discard the rest of the tree.
		_ = filepath.WalkDir(path, func(p string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				failures = append(failures, analyzer.FileResult{Path: p, Error: fmt.Errorf("failed to walk directory: %w", walkErr)})
				return nil
			}
			if !d.IsDir() {
				files = append(files, p)
			}
			return nil
		})
	} else {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, []analyzer.FileResult{{Path: path, Error: fmt.Errorf("failed to read directory: %w", err)}}
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				files = append(files, filepath.Join(path, entry.Name()))
			}
		}
	}

	return files, failures
}

func hashFile(f *os.File) (string, error) {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
