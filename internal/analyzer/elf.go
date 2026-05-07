package analyzer

import (
	"context"
	"io"
	"log/slog"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/internal/debuginfo"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/debuginfod"
)

// ELFAnalyzer runs ELF-specific analysis and returns findings.
type ELFAnalyzer struct {
	rules            []rule.ELFRule
	debuginfodClient *debuginfod.Client
	logger           *slog.Logger
}

// ELFAnalyzerOptions configures ELFAnalyzer creation.
type ELFAnalyzerOptions struct {
	Rules            []rule.ELFRule
	DebuginfodClient *debuginfod.Client
	Logger           *slog.Logger
}

// NewELFAnalyzer creates an ELF analyzer with the given options.
func NewELFAnalyzer(opts ELFAnalyzerOptions) *ELFAnalyzer {
	return &ELFAnalyzer{
		rules:            opts.Rules,
		debuginfodClient: opts.DebuginfodClient,
		logger:           opts.Logger.With(slog.String("component", "elf-analyzer")),
	}
}

// Analyze runs ELF-specific rules against the binary and returns findings.
func (a *ELFAnalyzer) Analyze(ctx context.Context, bin *binary.ELFBinary) []rule.Finding {
	if a.debuginfodClient != nil && bin.Build.BuildID != "" {
		a.logger.Debug("fetching debug symbols", slog.String("build_id", bin.Build.BuildID))

		rc, err := a.debuginfodClient.FetchDebugInfo(ctx, bin.Build.BuildID)
		if err != nil {
			a.logger.Debug("debug symbols not available", slog.Any("error", err))
		} else {
			defer rc.Close()
			if ra, ok := rc.(io.ReaderAt); ok {
				if err := debuginfo.ApplyDebugInfo(bin, ra, a.logger); err != nil {
					a.logger.Warn("failed to apply debug info", slog.Any("error", err))
				}
			} else {
				a.logger.Warn("debug info source does not support random access, skipping")
			}
		}
	}

	return rule.Check(a.rules, bin.Info, func(r rule.ELFRule) rule.Result {
		return r.Execute(bin)
	})
}
