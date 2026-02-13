package cli

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/mkacmar/crack/internal/analyzer"
	"github.com/mkacmar/crack/internal/output"
	"github.com/mkacmar/crack/internal/suggestions"
	"github.com/mkacmar/crack/rule"
)

func (a *App) processFullReport(resultsChan <-chan analyzer.FileResult, opts *outputOptions, invocation *output.InvocationInfo, rules []rule.ELFRule) int {
	var results []analyzer.FileResult
	var totalFailed int

	for res := range resultsChan {
		if res.Skipped {
			continue
		}
		results = append(results, res)
		totalFailed += res.FailedRules()
	}

	// Decorate findings with suggestions
	report := decorateReport(results)

	if opts.aggregate {
		agg := output.AggregateFindings(report, rules)
		fmt.Print(agg.Format())
	} else {
		textFormatter := &output.TextFormatter{IncludePassed: opts.includePassed, IncludeSkipped: opts.includeSkipped}
		if err := textFormatter.Format(report, os.Stdout); err != nil {
			a.logger.Error("failed to format output", slog.Any("error", err))
			return ExitError
		}
	}

	if opts.sarifOutput != "" {
		invocation.EndTime = time.Now()
		invocation.Successful = totalFailed == 0

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

	if totalFailed > 0 && !opts.exitZero {
		return ExitFindings
	}
	return ExitSuccess
}

func (a *App) processStreaming(resultsChan <-chan analyzer.FileResult, opts *outputOptions) int {
	var totalFailed int
	textFormatter := &output.TextFormatter{IncludePassed: opts.includePassed, IncludeSkipped: opts.includeSkipped}

	for res := range resultsChan {
		if res.Skipped {
			continue
		}
		totalFailed += res.FailedRules()
		report := decorateReport([]analyzer.FileResult{res})
		if err := textFormatter.Format(report, os.Stdout); err != nil {
			a.logger.Error("failed to format output", slog.Any("error", err))
		}
	}

	if totalFailed > 0 && !opts.exitZero {
		return ExitFindings
	}
	return ExitSuccess
}

func decorateReport(results []analyzer.FileResult) *output.DecoratedReport {
	decorated := make([]output.DecoratedFileResult, len(results))
	for i, res := range results {
		decorated[i] = output.DecoratedFileResult{
			FileResult: res,
			Findings:   suggestions.Decorate(res.Findings, res.Info),
		}
	}
	return &output.DecoratedReport{Results: decorated}
}
