package cli

import (
	"log/slog"
	"os"
	"time"

	"go.kacmar.sk/crack/internal/analyzer"
	"go.kacmar.sk/crack/internal/output"
	"go.kacmar.sk/crack/internal/suggestions"
)

func (a *App) processFullReport(resultsChan <-chan analyzer.FileResult, opts *outputOptions, invocation *output.InvocationInfo) int {
	var results []analyzer.FileResult
	var hasFindings, hasErrors bool

	for res := range resultsChan {
		if res.Skipped {
			continue
		}
		results = append(results, res)
		if res.FailedRules() > 0 {
			hasFindings = true
		}
		if res.Error != nil {
			hasErrors = true
		}
	}

	report := decorateReport(results)

	textFormatter := &output.TextFormatter{IncludePassed: opts.includePassed, IncludeSkipped: opts.includeSkipped}
	if err := textFormatter.Format(report, os.Stdout); err != nil {
		a.logger.Error("failed to format output", slog.Any("error", err))
		return ExitError
	}

	if opts.sarifOutput != "" {
		invocation.EndTime = time.Now()
		invocation.Successful = !hasErrors

		sarifFormatter := &output.SARIFFormatter{
			IncludePassed:  opts.includePassed,
			IncludeSkipped: opts.includeSkipped,
			Invocation:     invocation,
		}
		f, err := os.Create(opts.sarifOutput)
		if err != nil {
			a.logger.Error("failed to create SARIF file", slog.String("path", opts.sarifOutput), slog.Any("error", err))
			return ExitError
		}
		defer f.Close()
		if err := sarifFormatter.Format(report, f); err != nil {
			a.logger.Error("failed to write SARIF report", slog.Any("error", err))
			return ExitError
		}
		a.logger.Info("SARIF report saved", slog.String("path", opts.sarifOutput))
	}

	return exitCode(hasFindings, hasErrors, opts.exitZero)
}

func (a *App) processStreaming(resultsChan <-chan analyzer.FileResult, opts *outputOptions) int {
	var hasFindings, hasErrors bool
	textFormatter := &output.TextFormatter{IncludePassed: opts.includePassed, IncludeSkipped: opts.includeSkipped}

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
		report := decorateReport([]analyzer.FileResult{res})
		if err := textFormatter.Format(report, os.Stdout); err != nil {
			a.logger.Error("failed to format output", slog.Any("error", err))
		}
	}

	return exitCode(hasFindings, hasErrors, opts.exitZero)
}

// exitCode maps run outcomes to a process exit code, with file errors taking precedence over findings.
func exitCode(hasFindings, hasErrors, exitZero bool) int {
	switch {
	case hasErrors:
		return ExitError
	case hasFindings && !exitZero:
		return ExitFindings
	default:
		return ExitSuccess
	}
}

func decorateReport(results []analyzer.FileResult) *output.DecoratedReport {
	decorated := make([]output.DecoratedFileResult, len(results))
	for i, res := range results {
		decorated[i] = output.DecoratedFileResult{
			FileResult: res,
			Findings:   suggestions.Decorate(res.Findings, res.Profile),
		}
	}
	return &output.DecoratedReport{Results: decorated}
}
