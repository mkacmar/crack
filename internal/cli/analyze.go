package cli

import (
	"context"
	"flag"
	"fmt"
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
	fmt.Fprintf(os.Stderr, `Usage: %s analyze [options] <binary|directory>...

Analyze binaries for security hardening features.

Options:
  -P, --preset string         Security preset to use (default %q)
  -R, --rules string          Comma-separated list of rule IDs to run (mutually exclusive with --preset)
  -i, --input string          Read file paths from file (use "-" for stdin, mutually exclusive with positional args)
  -o, --sarif string          Save detailed SARIF report to file
  -a, --aggregate             Aggregate findings into actionable recommendations
  -r, --recursive             Recursively scan directories
      --log-level string      Log level: none, debug, info, warn, error (default "error")
      --show-passed           Show passing checks in output
      --show-skipped          Show skipped checks in output
  -p, --parallel int          Number of files to analyze in parallel (default %d)

Debuginfod options:
  -d, --debuginfod            Fetch debug symbols from debuginfod servers
      --debuginfod-urls       Debuginfod server URLs (default %q)
      --debuginfod-cache      Debuginfod cache directory (default "%s")
      --debuginfod-timeout    Debuginfod HTTP timeout (default %v)
      --debuginfod-retries    Debuginfod max retries per server (default %d)

Presets:
`, prog, preset.Default, runtime.NumCPU(), debuginfo.DefaultServerURL, debuginfo.DefaultCacheDir(), debuginfo.DefaultTimeout, debuginfo.DefaultRetries)

	for _, name := range preset.Names() {
		if name == preset.Default {
			fmt.Fprintf(os.Stderr, "  %s (default)\n", name)
		} else {
			fmt.Fprintf(os.Stderr, "  %s\n", name)
		}
	}
}

func (a *App) runAnalyze(prog string, args []string) int {
	startTime := time.Now()
	workingDir, _ := os.Getwd()

	fs := flag.NewFlagSet("analyze", flag.ExitOnError)

	var (
		presetName        string
		rulesFlag         string
		inputFile         string
		sarifOutput       string
		aggregate         bool
		recursive         bool
		logLevel          string
		showPassed        bool
		showSkipped       bool
		parallel          int
		useDebuginfod     bool
		debuginfodURLs    string
		debuginfodCache   string
		debuginfodTimeout time.Duration
		debuginfodRetries int
	)

	fs.StringVar(&presetName, "preset", "", "")
	fs.StringVar(&presetName, "P", "", "")
	fs.StringVar(&rulesFlag, "rules", "", "")
	fs.StringVar(&rulesFlag, "R", "", "")
	fs.StringVar(&inputFile, "input", "", "")
	fs.StringVar(&inputFile, "i", "", "")
	fs.StringVar(&sarifOutput, "sarif", "", "")
	fs.StringVar(&sarifOutput, "o", "", "")
	fs.BoolVar(&aggregate, "aggregate", false, "")
	fs.BoolVar(&aggregate, "a", false, "")
	fs.BoolVar(&recursive, "recursive", false, "")
	fs.BoolVar(&recursive, "r", false, "")
	fs.StringVar(&logLevel, "log-level", "error", "")
	fs.BoolVar(&showPassed, "show-passed", false, "")
	fs.BoolVar(&showSkipped, "show-skipped", false, "")
	fs.IntVar(&parallel, "parallel", runtime.NumCPU(), "")
	fs.IntVar(&parallel, "p", runtime.NumCPU(), "")
	fs.BoolVar(&useDebuginfod, "debuginfod", false, "")
	fs.BoolVar(&useDebuginfod, "d", false, "")
	fs.StringVar(&debuginfodURLs, "debuginfod-urls", debuginfo.DefaultServerURL, "")
	fs.StringVar(&debuginfodCache, "debuginfod-cache", "", "")
	fs.DurationVar(&debuginfodTimeout, "debuginfod-timeout", debuginfo.DefaultTimeout, "")
	fs.IntVar(&debuginfodRetries, "debuginfod-retries", debuginfo.DefaultRetries, "")

	fs.Usage = func() {
		a.printAnalyzeUsage(prog)
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if rulesFlag != "" && presetName != "" {
		fmt.Fprintf(os.Stderr, "Error: --rules and --preset are mutually exclusive\n")
		return 1
	}

	var p preset.Preset
	if rulesFlag != "" {
		ruleIDs := strings.Split(rulesFlag, ",")
		for i, id := range ruleIDs {
			ruleIDs[i] = strings.TrimSpace(id)
		}
		for _, id := range ruleIDs {
			if rule.Get(id) == nil {
				fmt.Fprintf(os.Stderr, "Error: unknown rule %q\n", id)
				fmt.Fprintf(os.Stderr, "Use 'crack list-rules' to see available rules.\n")
				return 1
			}
		}
		p = preset.Preset{Rules: ruleIDs}
	} else {
		if presetName == "" {
			presetName = preset.Default
		}
		var ok bool
		p, ok = preset.Get(presetName)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unknown preset %q\n", presetName)
			fmt.Fprintf(os.Stderr, "Available presets: %s\n", strings.Join(preset.Names(), ", "))
			return 1
		}
	}

	if fs.NArg() == 0 && inputFile == "" {
		fs.Usage()
		return 1
	}

	if fs.NArg() > 0 && inputFile != "" {
		fmt.Fprintf(os.Stderr, "Error: --input and positional arguments are mutually exclusive\n")
		return 1
	}

	var paths []string
	if inputFile != "" {
		var err error
		paths, err = readPathsFromInput(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return 1
		}
		if len(paths) == 0 {
			fmt.Fprintf(os.Stderr, "Error: no paths found in input\n")
			return 1
		}
	} else {
		paths = fs.Args()
	}

	if parallel < 1 {
		fmt.Fprintf(os.Stderr, "Error: --parallel must be at least 1\n")
		return 1
	}

	var logLevelValid bool
	a.logger, logLevelValid = setupLogger(logLevel)
	if !logLevelValid {
		fmt.Fprintf(os.Stderr, "Warning: invalid log level %q, using \"error\"\n", logLevel)
	}

	var debuginfodClient *debuginfo.Client
	if useDebuginfod {
		client, err := debuginfo.NewClient(debuginfo.Options{
			ServerURLs: parseURLList(debuginfodURLs),
			CacheDir:   debuginfodCache,
			Timeout:    debuginfodTimeout,
			MaxRetries: debuginfodRetries,
			Logger:     a.logger,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to initialize debuginfod client: %v\n", err)
			return 1
		}
		debuginfodClient = client
	}

	analyzer := elfanalyzer.NewAnalyzer(elfanalyzer.Options{
		RuleIDs:          p.Rules,
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
		return a.processFullReport(resultsChan, aggregate, showPassed, showSkipped, sarifOutput, invocation)
	}
	return a.processStreaming(resultsChan, showPassed, showSkipped)
}

func (a *App) processFullReport(resultsChan <-chan analyzer.Result, aggregate, showPassed, showSkipped bool, sarifOutput string, invocation *output.InvocationInfo) int {
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
		textFormatter, _ := output.GetFormatter("text", output.FormatterOptions{ShowPassed: showPassed, ShowSkipped: showSkipped})
		if err := textFormatter.Format(report, os.Stdout); err != nil {
			a.logger.Error("failed to format output", slog.Any("error", err))
			return 1
		}
	}

	if sarifOutput != "" {
		invocation.EndTime = time.Now()
		invocation.Successful = totalFailed == 0

		sarifFormatter, _ := output.GetFormatter("sarif", output.FormatterOptions{
			ShowPassed:  showPassed,
			ShowSkipped: showSkipped,
			Invocation:  invocation,
		})
		f, err := os.Create(sarifOutput)
		if err != nil {
			a.logger.Error("failed to create SARIF file", slog.String("path", sarifOutput), slog.Any("error", err))
			return 1
		}
		defer f.Close()
		if err := sarifFormatter.Format(report, f); err != nil {
			a.logger.Error("failed to write SARIF report", slog.Any("error", err))
			return 1
		}
		a.logger.Info("SARIF report saved", slog.String("path", sarifOutput))
	}

	if totalFailed > 0 {
		return 1
	}
	return 0
}

func (a *App) processStreaming(resultsChan <-chan analyzer.Result, showPassed, showSkipped bool) int {
	var totalFailed int
	textFormatter, _ := output.GetFormatter("text", output.FormatterOptions{ShowPassed: showPassed, ShowSkipped: showSkipped})

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

	if totalFailed > 0 {
		return 1
	}
	return 0
}
