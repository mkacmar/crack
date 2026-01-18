package scanner

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/mkacmar/crack/internal/debuginfo"
	"github.com/mkacmar/crack/internal/model"
	elfparser "github.com/mkacmar/crack/internal/parser/elf"
	"github.com/mkacmar/crack/internal/rules"
)

type Scanner struct {
	ruleEngine       *rules.Engine
	parsers          []model.BinaryParser
	logger           *slog.Logger
	workers          int
	debuginfodClient *debuginfo.Client
}

type Options struct {
	Logger            *slog.Logger
	Workers           int
	UseDebuginfod     bool
	DebuginfodURLs    []string
	DebuginfodCache   string
	DebuginfodTimeout time.Duration
	DebuginfodRetries int
}

func NewScanner(ruleEngine *rules.Engine, opts Options) *Scanner {
	var debuginfodClient *debuginfo.Client
	if opts.UseDebuginfod {
		client, err := debuginfo.NewClient(debuginfo.Options{
			ServerURLs: opts.DebuginfodURLs,
			CacheDir:   opts.DebuginfodCache,
			Timeout:    opts.DebuginfodTimeout,
			MaxRetries: opts.DebuginfodRetries,
			Logger:     opts.Logger,
		})
		if err != nil {
			opts.Logger.Warn("failed to initialize debuginfod client", slog.Any("error", err))
		} else {
			debuginfodClient = client
		}
	}

	return &Scanner{
		parsers: []model.BinaryParser{
			elfparser.NewParser(),
		},
		ruleEngine:       ruleEngine,
		logger:           opts.Logger,
		workers:          opts.Workers,
		debuginfodClient: debuginfodClient,
	}
}

func (s *Scanner) ScanPaths(ctx context.Context, paths []string, recursive bool) <-chan model.FileScanResult {
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

func (s *Scanner) scanFilesParallel(ctx context.Context, files []string) <-chan model.FileScanResult {
	results := make(chan model.FileScanResult)

	if len(files) == 0 {
		close(results)
		return results
	}

	s.logger.Debug("starting parallel scan", slog.Int("workers", s.workers), slog.Int("files", len(files)))

	go func() {
		g, ctx := errgroup.WithContext(ctx)
		g.SetLimit(s.workers)

		for _, file := range files {
			path := file
			g.Go(func() error {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				result := s.ScanFile(ctx, path)
				select {
				case results <- result:
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

func (s *Scanner) ScanFile(ctx context.Context, path string) model.FileScanResult {
	result := model.FileScanResult{
		Path: path,
	}

	s.logger.Debug("scanning file", slog.String("path", path))

	var parser model.BinaryParser
	for _, p := range s.parsers {
		canParse, err := p.CanParse(path)
		if err != nil {
			s.logger.Debug("failed to check file", slog.String("path", path), slog.Any("error", err))
			result.Error = err
			return result
		}
		if canParse {
			parser = p
			break
		}
	}

	if parser == nil {
		s.logger.Debug("skipping non-binary file", slog.String("path", path))
		result.Skipped = true
		return result
	}

	info, err := parser.Parse(path)
	if err != nil {
		s.logger.Debug("failed to parse binary", slog.String("path", path), slog.Any("error", err))
		result.Error = err
		return result
	}

	if info.ELFFile != nil {
		defer info.ELFFile.Close()
	}

	result.Format = info.Format
	s.logger.Debug("parsed binary", slog.String("path", path), slog.String("format", info.Format.String()), slog.String("arch", info.Architecture.String()))

	if s.debuginfodClient != nil && info.Build.BuildID != "" {
		s.logger.Debug("fetching debug symbols", slog.String("path", path), slog.String("build_id", info.Build.BuildID))

		debugPath, err := s.debuginfodClient.FetchDebugInfo(ctx, info.Build.BuildID)
		if err != nil {
			s.logger.Debug("failed to fetch debug symbols", slog.String("path", path), slog.Any("error", err))
		} else {
			if err := debuginfo.EnhanceWithDebugInfo(info, debugPath, s.logger); err != nil {
				s.logger.Debug("failed to parse debug symbols", slog.String("path", path), slog.Any("error", err))
			}
		}
	}

	result.Toolchain = info.Build.Toolchain

	checkResults := s.ruleEngine.ExecuteRules(info)
	result.Results = checkResults

	s.logger.Debug("scan complete",
		slog.String("path", path),
		slog.Int("passed", result.PassedChecks()),
		slog.Int("failed", result.FailedChecks()))

	return result
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
