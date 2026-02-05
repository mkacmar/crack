package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/mkacmar/crack/internal/analyzer"
	elfanalyzer "github.com/mkacmar/crack/internal/analyzer/elf"
	"github.com/mkacmar/crack/internal/debuginfo"
	"github.com/mkacmar/crack/internal/output"
	"github.com/mkacmar/crack/internal/preset"
	"github.com/mkacmar/crack/internal/rule"
	elfrules "github.com/mkacmar/crack/internal/rules/elf"
	"github.com/mkacmar/crack/internal/scanner"
)

func init() {
	elfrules.RegisterRules()
}

func (a *App) printAnalyzeUsage(prog string) {
	fmt.Fprintf(os.Stderr, `Usage: %s analyze [options] [<path>...]

Analyze binaries for security hardening features.

Options:
  -i, --input string          Read file paths from file, one path per line (use "-" for stdin, mutually exclusive with positional args)
  -p, --parallel int          Number of files to analyze in parallel (default %d)
  -r, --recursive             Recursively scan directories
      --rules string          Comma-separated list of rule IDs to run

Output options:
  -a, --aggregate             Aggregate findings into actionable recommendations
      --exit-zero             Exit with 0 even when findings are detected
      --include-passed        Include passing checks in output
      --include-skipped       Include skipped checks in output
      --sarif string          Save detailed SARIF report to file

Logging options:
      --log string            Log output file (default stderr)
      --log-level string      Log level: none, debug, info, warn, error (default "error")

Debuginfod options:
      --debuginfod                  Fetch debug symbols from debuginfod servers
      --debuginfod-cache string     Debuginfod cache directory (default "%s")
      --debuginfod-retries int      Debuginfod max retries per server (default %d)
      --debuginfod-servers string   Comma-separated debuginfod server URLs (default %q)
      --debuginfod-timeout duration Debuginfod HTTP timeout (default %v)
`, prog, runtime.NumCPU(), debuginfo.DefaultCacheDir(), debuginfo.DefaultRetries, debuginfo.DefaultServerURL, debuginfo.DefaultTimeout)
}

func (a *App) runAnalyze(prog string, args []string) int {
	startTime := time.Now()
	workingDir, _ := os.Getwd()

	fs := flag.NewFlagSet("analyze", flag.ExitOnError)

	var (
		rulesFlag         string
		inputFile         string
		sarifOutput       string
		aggregate         bool
		recursive         bool
		logFile           string
		logLevel          string
		includePassed     bool
		includeSkipped    bool
		parallel          int
		exitZero          bool
		useDebuginfod     bool
		debuginfodServers string
		debuginfodCache   string
		debuginfodTimeout time.Duration
		debuginfodRetries int
	)

	fs.StringVar(&rulesFlag, "rules", "", "")
	fs.StringVar(&inputFile, "input", "", "")
	fs.StringVar(&inputFile, "i", "", "")
	fs.StringVar(&sarifOutput, "sarif", "", "")
	fs.BoolVar(&aggregate, "aggregate", false, "")
	fs.BoolVar(&aggregate, "a", false, "")
	fs.BoolVar(&recursive, "recursive", false, "")
	fs.BoolVar(&recursive, "r", false, "")
	fs.StringVar(&logFile, "log", "", "")
	fs.StringVar(&logLevel, "log-level", "error", "")
	fs.BoolVar(&includePassed, "include-passed", false, "")
	fs.BoolVar(&includeSkipped, "include-skipped", false, "")
	fs.IntVar(&parallel, "parallel", runtime.NumCPU(), "")
	fs.IntVar(&parallel, "p", runtime.NumCPU(), "")
	fs.BoolVar(&exitZero, "exit-zero", false, "")
	fs.BoolVar(&useDebuginfod, "debuginfod", false, "")
	fs.StringVar(&debuginfodServers, "debuginfod-servers", debuginfo.DefaultServerURL, "")
	fs.StringVar(&debuginfodCache, "debuginfod-cache", "", "")
	fs.DurationVar(&debuginfodTimeout, "debuginfod-timeout", debuginfo.DefaultTimeout, "")
	fs.IntVar(&debuginfodRetries, "debuginfod-retries", debuginfo.DefaultRetries, "")

	fs.Usage = func() {
		a.printAnalyzeUsage(prog)
	}

	if err := fs.Parse(args); err != nil {
		return ExitError
	}

	var ruleIDs []string
	if rulesFlag != "" {
		ruleIDs = strings.Split(rulesFlag, ",")
		for i, id := range ruleIDs {
			ruleIDs[i] = strings.TrimSpace(id)
		}
		for _, id := range ruleIDs {
			if rule.Get(id) == nil {
				fmt.Fprintf(os.Stderr, "Error: unknown rule %q\n", id)
				return ExitError
			}
		}
	} else {
		ruleIDs = preset.DefaultRules
	}

	if fs.NArg() == 0 && inputFile == "" {
		fs.Usage()
		return ExitError
	}

	if fs.NArg() > 0 && inputFile != "" {
		fmt.Fprintf(os.Stderr, "Error: --input and positional arguments are mutually exclusive\n")
		return ExitError
	}

	var paths []string
	if inputFile != "" {
		var err error
		paths, err = readPathsFromInput(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return ExitError
		}
		if len(paths) == 0 {
			fmt.Fprintf(os.Stderr, "Error: no paths found in input\n")
			return ExitError
		}
	} else {
		paths = fs.Args()
	}

	if parallel < 1 {
		fmt.Fprintf(os.Stderr, "Error: --parallel must be at least 1\n")
		return ExitError
	}

	var logOutput io.Writer = os.Stderr
	if logFile != "" {
		f, err := os.Create(logFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to open log file: %v\n", err)
			return ExitError
		}
		defer f.Close()
		logOutput = f
	}
	a.logger = setupLogger(logLevel, logOutput)

	var debuginfodClient *debuginfo.Client
	if useDebuginfod {
		client, err := debuginfo.NewClient(debuginfo.Options{
			ServerURLs: parseURLList(debuginfodServers),
			CacheDir:   debuginfodCache,
			Timeout:    debuginfodTimeout,
			MaxRetries: debuginfodRetries,
			Logger:     a.logger,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to initialize debuginfod client: %v\n", err)
			return ExitError
		}
		debuginfodClient = client
	}

	analyzer := elfanalyzer.NewAnalyzer(elfanalyzer.Options{
		RuleIDs:          ruleIDs,
		DebuginfodClient: debuginfodClient,
		Logger:           a.logger,
	})

	scannerOpts := scanner.Options{
		Logger:  a.logger,
		Workers: parallel,
	}
	scan := scanner.NewScanner(analyzer, scannerOpts)

	ctx := context.Background()

	a.logger.Info("starting scan", slog.Int("paths", len(paths)), slog.Bool("recursive", recursive))

	resultsChan := scan.ScanPaths(ctx, paths, recursive)

	needsFullReport := aggregate || sarifOutput != ""

	invocation := &output.InvocationInfo{
		CommandLine: strings.Join(append([]string{prog}, args...), " "),
		Arguments:   args,
		StartTime:   startTime,
		WorkingDir:  workingDir,
	}

	if needsFullReport {
		return a.processFullReport(resultsChan, aggregate, includePassed, includeSkipped, sarifOutput, invocation, exitZero)
	}
	return a.processStreaming(resultsChan, includePassed, includeSkipped, exitZero)
}

func (a *App) processFullReport(resultsChan <-chan analyzer.Result, aggregate, includePassed, includeSkipped bool, sarifOutput string, invocation *output.InvocationInfo, exitZero bool) int {
	var results []analyzer.Result
	var totalFailed int

	for res := range resultsChan {
		if res.Skipped {
			continue
		}
		results = append(results, res)
		totalFailed += res.FailedRules()
	}

	report := &analyzer.Results{Results: results}

	if aggregate {
		agg := output.AggregateFindings(report)
		fmt.Print(output.FormatAggregated(agg))
	} else {
		textFormatter, _ := output.GetFormatter("text", output.FormatterOptions{IncludePassed: includePassed, IncludeSkipped: includeSkipped})
		if err := textFormatter.Format(report, os.Stdout); err != nil {
			a.logger.Error("failed to format output", slog.Any("error", err))
			return ExitError
		}
	}

	if sarifOutput != "" {
		invocation.EndTime = time.Now()
		invocation.Successful = totalFailed == 0

		sarifFormatter, _ := output.GetFormatter("sarif", output.FormatterOptions{
			IncludePassed:  includePassed,
			IncludeSkipped: includeSkipped,
			Invocation:     invocation,
		})
		f, err := os.Create(sarifOutput)
		if err != nil {
			a.logger.Error("failed to create SARIF file", slog.String("path", sarifOutput), slog.Any("error", err))
			return ExitError
		}
		defer f.Close()
		if err := sarifFormatter.Format(report, f); err != nil {
			a.logger.Error("failed to write SARIF report", slog.Any("error", err))
			return ExitError
		}
		a.logger.Info("SARIF report saved", slog.String("path", sarifOutput))
	}

	if totalFailed > 0 && !exitZero {
		return ExitFindings
	}
	return ExitSuccess
}

func (a *App) processStreaming(resultsChan <-chan analyzer.Result, includePassed, includeSkipped bool, exitZero bool) int {
	var totalFailed int
	textFormatter, _ := output.GetFormatter("text", output.FormatterOptions{IncludePassed: includePassed, IncludeSkipped: includeSkipped})

	for res := range resultsChan {
		if res.Skipped {
			continue
		}
		totalFailed += res.FailedRules()
		singleReport := &analyzer.Results{Results: []analyzer.Result{res}}
		if err := textFormatter.Format(singleReport, os.Stdout); err != nil {
			a.logger.Error("failed to format output", slog.Any("error", err))
		}
	}

	if totalFailed > 0 && !exitZero {
		return ExitFindings
	}
	return ExitSuccess
}
