package analyzer

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"go.kacmar.sk/crack/binary"
)

// Dispatcher routes binary analysis to format-specific analyzers.
type Dispatcher struct {
	elf    *ELFAnalyzer
	logger *slog.Logger
}

// DispatcherOptions configures Dispatcher creation.
type DispatcherOptions struct {
	ELF    *ELFAnalyzer
	Logger *slog.Logger
}

// NewDispatcher creates a dispatcher with the given analyzers.
func NewDispatcher(opts DispatcherOptions) *Dispatcher {
	return &Dispatcher{
		elf:    opts.ELF,
		logger: opts.Logger.With(slog.String("component", "dispatcher")),
	}
}

// Analyze parses the binary and returns analysis results.
// Returns slice to handle fat/universal binaries (one result per arch slice).
// Returns ErrUnsupportedFormat if no parser matches, other errors for parse failures.
func (d *Dispatcher) Analyze(ctx context.Context, r io.ReaderAt) ([]AnalysisResult, error) {
	if bin, err := binary.ParseELF(r); err == nil {
		d.logger.Debug("parsed ELF binary",
			slog.String("format", bin.Format.String()),
			slog.String("arch", bin.Architecture.String()))

		findings := d.elf.Analyze(ctx, bin)
		return []AnalysisResult{{
			Info:     bin.Info,
			Findings: findings,
		}}, nil
	} else if !errors.Is(err, binary.ErrUnsupportedFormat) {
		return nil, err
	}

	return nil, ErrUnrecognizedFormat
}
