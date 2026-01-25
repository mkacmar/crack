package elf

import (
	"context"
	"errors"
	"log/slog"

	"github.com/mkacmar/crack/internal/analyzer"
	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/debuginfo"
	"github.com/mkacmar/crack/internal/parser/elf"
	"github.com/mkacmar/crack/internal/rules"
)

type Analyzer struct {
	engine           *rules.Engine
	parser           *elf.Parser
	debuginfodClient *debuginfo.Client
	logger           *slog.Logger
}

type Options struct {
	RuleIDs          []string
	DebuginfodClient *debuginfo.Client
	Logger           *slog.Logger
}

func NewAnalyzer(opts Options) *Analyzer {

	engine := rules.NewEngine(opts.Logger)
	engine.LoadRules(opts.RuleIDs)

	return &Analyzer{
		engine:           engine,
		parser:           elf.NewParser(),
		debuginfodClient: opts.DebuginfodClient,
		logger:           opts.Logger.With(slog.String("component", "elf-analyzer")),
	}
}

func (a *Analyzer) Analyze(ctx context.Context, path string) analyzer.Result {
	res := analyzer.Result{
		Path: path,
	}

	bin, err := a.parser.Parse(path)
	if err != nil {
		if errors.Is(err, binary.ErrUnsupportedFormat) {
			a.logger.Debug("skipping non-ELF file", slog.String("path", path))
			res.Skipped = true
			return res
		}
		a.logger.Warn("failed to parse binary", slog.String("path", path), slog.Any("error", err))
		res.Error = err
		return res
	}

	if bin.File != nil {
		defer bin.File.Close()
	}

	res.Format = bin.Format
	a.logger.Debug("parsed binary",
		slog.String("path", path),
		slog.String("format", bin.Format.String()),
		slog.String("arch", bin.Architecture.String()))

	if a.debuginfodClient != nil && bin.Build.BuildID != "" {
		a.logger.Debug("fetching debug symbols",
			slog.String("path", path),
			slog.String("build_id", bin.Build.BuildID))

		debugPath, err := a.debuginfodClient.FetchDebugInfo(ctx, bin.Build.BuildID)
		if err != nil {
			a.logger.Debug("debug symbols not available",
				slog.String("path", path),
				slog.Any("error", err))
		} else {
			if err := debuginfo.EnhanceWithDebugInfo(bin, debugPath, a.logger); err != nil {
				a.logger.Warn("failed to enhance with debug info",
					slog.String("path", path),
					slog.Any("error", err))
			}
		}
	}

	res.Toolchain = bin.Build.Toolchain
	res.Results = a.engine.ExecuteRules(bin)

	a.logger.Debug("analysis complete",
		slog.String("path", path),
		slog.Int("passed", res.PassedRules()),
		slog.Int("failed", res.FailedRules()))

	return res
}
