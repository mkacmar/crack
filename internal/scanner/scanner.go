package scanner

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	"github.com/mkacmar/crack/internal/analyzer"
)

type Scanner struct {
	analyzer analyzer.FileAnalyzer
	logger   *slog.Logger
	workers  int
}

type Options struct {
	Logger  *slog.Logger
	Workers int
}

func NewScanner(fileAnalyzer analyzer.FileAnalyzer, opts Options) *Scanner {
	return &Scanner{
		analyzer: fileAnalyzer,
		logger:   opts.Logger.With(slog.String("component", "scanner")),
		workers:  opts.Workers,
	}
}

func (s *Scanner) ScanPaths(ctx context.Context, paths []string, recursive bool) <-chan analyzer.FileResult {
	var filesToScan []string

	for _, path := range paths {
		files, err := s.collectFiles(path, recursive)
		if err != nil {
			s.logger.Error("failed to collect files", slog.String("path", path), slog.Any("error", err))
			continue
		}
		filesToScan = append(filesToScan, files...)
	}

	s.logger.Debug("collected files to scan", slog.Int("count", len(filesToScan)))

	return s.scanFilesParallel(ctx, filesToScan)
}

func (s *Scanner) scanFilesParallel(ctx context.Context, files []string) <-chan analyzer.FileResult {
	results := make(chan analyzer.FileResult)

	if len(files) == 0 {
		close(results)
		return results
	}

	s.logger.Debug("starting parallel scan", slog.Int("workers", s.workers), slog.Int("files", len(files)))

	go func() {
		g, ctx := errgroup.WithContext(ctx)
		g.SetLimit(s.workers)

		for _, path := range files {
			g.Go(func() error {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				res := s.scanFile(ctx, path)
				select {
				case results <- res:
				case <-ctx.Done():
					return ctx.Err()
				}
				return nil
			})
		}

		_ = g.Wait()
		close(results)
	}()

	return results
}

func (s *Scanner) scanFile(ctx context.Context, path string) analyzer.FileResult {
	s.logger.Debug("scanning file", slog.String("path", path))

	res := s.analyzer.Analyze(ctx, path)

	hash, err := computeSHA256(path)
	if err != nil {
		s.logger.Warn("failed to compute SHA256", slog.String("path", path), slog.Any("error", err))
	}
	res.SHA256 = hash

	return res
}

func (s *Scanner) collectFiles(path string, recursive bool) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}

	if !info.IsDir() {
		return []string{path}, nil
	}

	var files []string

	if recursive {
		err := filepath.WalkDir(path, func(p string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if !d.IsDir() {
				files = append(files, p)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to walk directory: %w", err)
		}
	} else {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory: %w", err)
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				files = append(files, filepath.Join(path, entry.Name()))
			}
		}
	}

	return files, nil
}

func computeSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
