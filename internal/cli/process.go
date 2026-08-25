package cli

import (
	"context"
	"log/slog"
	"os"
	"time"

	"go.kacmar.sk/crack/internal/analyzer"
	"go.kacmar.sk/crack/internal/output"
	"go.kacmar.sk/crack/internal/suggestions"
)

func (a *App) processResults(ctx context.Context, resultsChan <-chan analyzer.FileResult, opts *outputOptions, invocation *output.InvocationInfo) int {
	textFormatter := &output.TextFormatter{IncludePassed: opts.includePassed, IncludeSkipped: opts.includeSkipped}

	var sarifWriter *output.SARIFWriter
	if opts.sarifOutput != "" {
		f, err := os.Create(opts.sarifOutput)
		if err != nil {
			a.logger.Error("failed to create SARIF file", slog.String("path", opts.sarifOutput), slog.Any("error", err))
			return ExitError
		}
		defer f.Close()
		sarifWriter = output.NewSARIFWriter(f, opts.includePassed, opts.includeSkipped)
	}

	var hasFindings, hasErrors bool
	for res := range resultsChan {
		if res.Skipped {
			continue
		}
		if res.FailedRules() > 0 {
			hasFindings = true
		}
		if res.Error != nil {
			hasErrors = true
		}

		decorated := output.DecoratedFileResult{
			FileResult: res,
			Findings:   suggestions.Decorate(res.Findings, res.Profile),
		}

		if err := textFormatter.Format(&output.DecoratedReport{Results: []output.DecoratedFileResult{decorated}}, os.Stdout); err != nil {
			a.logger.Error("failed to format output", slog.Any("error", err))
		}
		if sarifWriter != nil {
			if err := sarifWriter.Write(decorated); err != nil {
				a.logger.Error("failed to write SARIF result", slog.Any("error", err))
			}
		}
	}

	interrupted := ctx.Err() != nil

	if sarifWriter != nil {
		invocation.EndTime = time.Now()
		invocation.Successful = !hasErrors && !interrupted
		if interrupted {
			sarifWriter.Notify("error", "Scan interrupted before completion, results are incomplete")
		}
		if err := sarifWriter.Close(invocation); err != nil {
			a.logger.Error("failed to write SARIF report", slog.String("path", opts.sarifOutput), slog.Any("error", err))
			return ExitError
		}
		a.logger.Info("SARIF report saved", slog.String("path", opts.sarifOutput))
	}

	return exitCode(hasFindings, hasErrors, interrupted, opts.exitZero)
}

// exitCode maps run outcomes to a process exit code.
// File errors and interruption take precedence over findings, and --exit-zero suppresses only findings.
func exitCode(hasFindings, hasErrors, interrupted, exitZero bool) int {
	switch {
	case hasErrors, interrupted:
		return ExitError
	case hasFindings && !exitZero:
		return ExitFindings
	default:
		return ExitSuccess
	}
}
