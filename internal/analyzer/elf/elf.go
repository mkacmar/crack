package elf

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/internal/analyzer"
	"github.com/mkacmar/crack/internal/debuginfo"
	"github.com/mkacmar/crack/internal/rules"
	"github.com/mkacmar/crack/rule"
)

type Analyzer struct {
	engine           *rules.Engine
	debuginfodClient *debuginfo.Client
	logger           *slog.Logger
}

type Options struct {
	Rules            []rule.ELFRule
	DebuginfodClient *debuginfo.Client
	Logger           *slog.Logger
}

func NewAnalyzer(opts Options) *Analyzer {
	ruleList := make([]rule.Rule, len(opts.Rules))
	for i, r := range opts.Rules {
		ruleList[i] = r
	}

	return &Analyzer{
		engine:           rules.NewEngine(ruleList, opts.Logger),
		debuginfodClient: opts.DebuginfodClient,
		logger:           opts.Logger.With(slog.String("component", "elf-analyzer")),
	}
}

func (a *Analyzer) Analyze(ctx context.Context, path string) analyzer.FileResult {
	res := analyzer.FileResult{
		Path: path,
	}

	f, err := os.Open(path)
	if err != nil {
		a.logger.Warn("failed to open file", slog.String("path", path), slog.Any("error", err))
		res.Error = err
		return res
	}
	defer f.Close()

	bin, err := binary.ParseELF(f)
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

	res.Build = bin.Build
	res.Findings = a.engine.ExecuteRules(bin)

	a.logger.Debug("analysis complete",
		slog.String("path", path),
		slog.Int("passed", res.PassedRules()),
		slog.Int("failed", res.FailedRules()))

	return res
}
