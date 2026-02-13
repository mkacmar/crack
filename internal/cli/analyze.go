package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/mkacmar/crack/internal/analyzer"
	"github.com/mkacmar/crack/internal/debuginfo"
	"github.com/mkacmar/crack/internal/output"
	"github.com/mkacmar/crack/internal/preset"
	"github.com/mkacmar/crack/internal/rules"
	"github.com/mkacmar/crack/internal/scanner"
	"github.com/mkacmar/crack/internal/suggestions"
	"github.com/mkacmar/crack/rule"
)

var errNoPathsSpecified = fmt.Errorf("no paths specified")

type outputOptions struct {
	aggregate      bool
	includePassed  bool
	includeSkipped bool
	exitZero       bool
	sarifOutput    string
}

type analyzeConfig struct {
	rulesFlag         string
	targetPlatform    string
	targetCompiler    string
	inputFile         string
	recursive         bool
	logFile           string
	logLevel          string
	parallel          int
	useDebuginfod     bool
	debuginfodServers string
	debuginfodCache   string
	debuginfodTimeout time.Duration
	debuginfodRetries int
	debuginfodMaxSize int64
}

func (a *App) printAnalyzeUsage(prog string) {
	fmt.Fprintf(os.Stderr, `Usage: %s analyze [options] [<path>...]

Analyze binaries for security hardening features.

Options:
  -i, --input string          Read file paths from file, one path per line (use "-" for stdin, mutually exclusive with positional args)
  -p, --parallel int          Number of files to analyze in parallel (default %d)
  -r, --recursive             Recursively scan directories

`, prog, runtime.NumCPU())

	fmt.Fprintf(os.Stderr, `Rule selection:
      --rules string              Comma-separated list of rule IDs to run
      --target-compiler string    Only run rules available for these compilers: %s
      --target-platform string    Only run rules available for these platforms: %s

`, strings.Join(validCompilerNames(), ", "), strings.Join(validArchitectureNames(), ", "))

	fmt.Fprint(os.Stderr, `Output options:
      --aggregate             Aggregate findings into actionable recommendations
      --exit-zero             Exit with 0 even when findings are detected
      --include-passed        Include passing checks in output
      --include-skipped       Include skipped checks in output
      --sarif string          Save detailed SARIF report to file

Logging options:
      --log string            Log output file (default stderr)
      --log-level string      Log level: none, debug, info, warn, error (default "error")

`)

	defaultCacheDir, _ := debuginfo.DefaultCacheDir()
	fmt.Fprintf(os.Stderr, `Debuginfod options:
      --debuginfod                  Fetch debug symbols from debuginfod servers
      --debuginfod-cache string     Debuginfod cache directory (default "%s")
      --debuginfod-max-size bytes   Max debug file size per download (default %d)
      --debuginfod-retries int      Debuginfod max retries per server (default %d)
      --debuginfod-servers string   Comma-separated debuginfod server URLs (default %q)
      --debuginfod-timeout duration Debuginfod HTTP timeout (default %v)
`, defaultCacheDir, debuginfo.DefaultMaxFileSize, debuginfo.DefaultRetries, debuginfo.DefaultServerURL, debuginfo.DefaultTimeout)
}

func parseRules(rulesFlag, targetPlatform, targetCompiler string) ([]rule.ELFRule, error) {
	var selectedRules []rule.ELFRule
	if rulesFlag != "" {
		ids := strings.Split(rulesFlag, ",")
		for _, id := range ids {
			id = strings.TrimSpace(id)
			r, ok := rules.Find[rule.ELFRule](rules.ByID(id))
			if !ok {
				return nil, fmt.Errorf("unknown rule %q", id)
			}
			selectedRules = append(selectedRules, r)
		}
	} else {
		selectedRules = preset.Default()
	}

	if targetPlatform != "" || targetCompiler != "" {
		filter, err := ParseTargetFilter(targetPlatform, targetCompiler)
		if err != nil {
			return nil, err
		}
		selectedRules = rule.FilterRules(selectedRules, filter)
		if len(selectedRules) == 0 {
			return nil, fmt.Errorf("no rules match the specified target filter")
		}
	}

	return selectedRules, nil
}

func parsePaths(fs *flag.FlagSet, inputFile string) ([]string, error) {
	if fs.NArg() == 0 && inputFile == "" {
		return nil, errNoPathsSpecified
	}

	if fs.NArg() > 0 && inputFile != "" {
		return nil, fmt.Errorf("--input and positional arguments are mutually exclusive")
	}

	if inputFile != "" {
		paths, err := readPathsFromInput(inputFile)
		if err != nil {
			return nil, err
		}
		if len(paths) == 0 {
			return nil, fmt.Errorf("no paths found in input")
		}
		return paths, nil
	}

	return fs.Args(), nil
}

func (a *App) runAnalyze(prog string, args []string) int {
	startTime := time.Now()
	workingDir, _ := os.Getwd()

	fs, opts, cfg := a.setupAnalyzeFlags(prog)
	if err := fs.Parse(args); err != nil {
		return ExitError
	}

	selectedRules, err := parseRules(cfg.rulesFlag, cfg.targetPlatform, cfg.targetCompiler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitError
	}

	paths, err := parsePaths(fs, cfg.inputFile)
	if err != nil {
		if errors.Is(err, errNoPathsSpecified) {
			fs.Usage()
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		return ExitError
	}

	if cfg.parallel < 1 {
		fmt.Fprintf(os.Stderr, "Error: --parallel must be at least 1\n")
		return ExitError
	}

	if err := a.setupLogging(cfg.logFile, cfg.logLevel); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitError
	}

	debuginfodClient, err := a.setupDebuginfod(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitError
	}

	elfAnalyzer := analyzer.NewELFAnalyzer(analyzer.ELFAnalyzerOptions{
		Rules:            selectedRules,
		DebuginfodClient: debuginfodClient,
		Logger:           a.logger,
	})

	dispatcher := analyzer.NewDispatcher(analyzer.DispatcherOptions{
		ELF:    elfAnalyzer,
		Logger: a.logger,
	})

	scan := scanner.NewScanner(dispatcher, scanner.Options{
		Logger:  a.logger,
		Workers: cfg.parallel,
	})

	ctx := context.Background()
	a.logger.Info("starting scan", slog.Int("paths", len(paths)), slog.Bool("recursive", cfg.recursive))

	resultsChan := scan.ScanPaths(ctx, paths, cfg.recursive)

	invocation := &output.InvocationInfo{
		CommandLine: strings.Join(append([]string{prog}, args...), " "),
		Arguments:   args,
		StartTime:   startTime,
		WorkingDir:  workingDir,
	}

	if opts.aggregate || opts.sarifOutput != "" {
		return a.processFullReport(resultsChan, opts, invocation, selectedRules)
	}
	return a.processStreaming(resultsChan, opts)
}

func (a *App) setupAnalyzeFlags(prog string) (*flag.FlagSet, *outputOptions, *analyzeConfig) {
	fs := flag.NewFlagSet("analyze", flag.ExitOnError)
	opts := &outputOptions{}
	cfg := &analyzeConfig{}

	fs.StringVar(&cfg.rulesFlag, "rules", "", "")
	fs.StringVar(&cfg.targetPlatform, "target-platform", "", "")
	fs.StringVar(&cfg.targetCompiler, "target-compiler", "", "")
	fs.StringVar(&cfg.inputFile, "input", "", "")
	fs.StringVar(&opts.sarifOutput, "sarif", "", "")
	fs.BoolVar(&opts.aggregate, "aggregate", false, "")
	fs.BoolVar(&cfg.recursive, "recursive", false, "")
	fs.StringVar(&cfg.logFile, "log", "", "")
	fs.StringVar(&cfg.logLevel, "log-level", "error", "")
	fs.BoolVar(&opts.includePassed, "include-passed", false, "")
	fs.BoolVar(&opts.includeSkipped, "include-skipped", false, "")
	fs.IntVar(&cfg.parallel, "parallel", runtime.NumCPU(), "")
	fs.BoolVar(&opts.exitZero, "exit-zero", false, "")
	fs.BoolVar(&cfg.useDebuginfod, "debuginfod", false, "")
	fs.StringVar(&cfg.debuginfodServers, "debuginfod-servers", debuginfo.DefaultServerURL, "")
	fs.StringVar(&cfg.debuginfodCache, "debuginfod-cache", "", "")
	fs.DurationVar(&cfg.debuginfodTimeout, "debuginfod-timeout", debuginfo.DefaultTimeout, "")
	fs.IntVar(&cfg.debuginfodRetries, "debuginfod-retries", debuginfo.DefaultRetries, "")
	fs.Int64Var(&cfg.debuginfodMaxSize, "debuginfod-max-size", debuginfo.DefaultMaxFileSize, "")

	fs.Usage = func() { a.printAnalyzeUsage(prog) }

	return fs, opts, cfg
}

func (a *App) setupLogging(logFile, logLevel string) error {
	var logOutput io.Writer = os.Stderr
	if logFile != "" {
		f, err := os.Create(logFile)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		logOutput = f
	}
	a.logger = setupLogger(logLevel, logOutput)
	return nil
}

func (a *App) setupDebuginfod(cfg *analyzeConfig) (*debuginfo.Client, error) {
	if !cfg.useDebuginfod {
		return nil, nil
	}
	return debuginfo.NewClient(debuginfo.Options{
		ServerURLs:  strings.Split(cfg.debuginfodServers, ","),
		CacheDir:    cfg.debuginfodCache,
		Timeout:     cfg.debuginfodTimeout,
		MaxRetries:  cfg.debuginfodRetries,
		MaxFileSize: cfg.debuginfodMaxSize,
		Logger:      a.logger,
	})
}

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
		fmt.Print(output.FormatAggregated(agg))
	} else {
		textFormatter, _ := output.GetFormatter("text", output.FormatterOptions{IncludePassed: opts.includePassed, IncludeSkipped: opts.includeSkipped})
		if err := textFormatter.Format(report, os.Stdout); err != nil {
			a.logger.Error("failed to format output", slog.Any("error", err))
			return ExitError
		}
	}

	if opts.sarifOutput != "" {
		invocation.EndTime = time.Now()
		invocation.Successful = totalFailed == 0

		sarifFormatter, _ := output.GetFormatter("sarif", output.FormatterOptions{
			IncludePassed:  opts.includePassed,
			IncludeSkipped: opts.includeSkipped,
			Invocation:     invocation,
		})
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
	textFormatter, _ := output.GetFormatter("text", output.FormatterOptions{IncludePassed: opts.includePassed, IncludeSkipped: opts.includeSkipped})

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
