package analyzer

import (
	"context"
	"log/slog"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/internal/debuginfo"
	"github.com/mkacmar/crack/rule"
)

// ELFAnalyzer runs ELF-specific analysis and returns findings.
type ELFAnalyzer struct {
	rules            []rule.ELFRule
	debuginfodClient *debuginfo.Client
	logger           *slog.Logger
}

// ELFAnalyzerOptions configures ELFAnalyzer creation.
type ELFAnalyzerOptions struct {
	Rules            []rule.ELFRule
	DebuginfodClient *debuginfo.Client
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

		debugPath, err := a.debuginfodClient.FetchDebugInfo(ctx, bin.Build.BuildID)
		if err != nil {
			a.logger.Debug("debug symbols not available", slog.Any("error", err))
		} else {
			if err := debuginfo.EnhanceWithDebugInfo(bin, debugPath, a.logger); err != nil {
				a.logger.Warn("failed to enhance with debug info", slog.Any("error", err))
			}
		}
	}

	return rule.Check(a.rules, bin.Info, func(r rule.ELFRule) rule.Result {
		return r.Execute(bin)
	})
}
