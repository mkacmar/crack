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
	elfanalyzer "github.com/mkacmar/crack/internal/analyzer/elf"
	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/debuginfo"
	"github.com/mkacmar/crack/internal/output"
	"github.com/mkacmar/crack/internal/preset"
	"github.com/mkacmar/crack/internal/rule"
	elfrules "github.com/mkacmar/crack/internal/rules/elf"
	"github.com/mkacmar/crack/internal/scanner"
	"github.com/mkacmar/crack/internal/toolchain"
)

func init() {
	elfrules.RegisterRules()
}

var errNoPathsSpecified = fmt.Errorf("no paths specified")

type outputOptions struct {
	aggregate      bool
	includePassed  bool
	includeSkipped bool
	exitZero       bool
	sarifOutput    string
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

`, strings.Join(toolchain.ValidCompilerNames(), ", "), strings.Join(binary.ValidArchitectureNames(), ", "))

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

	fmt.Fprintf(os.Stderr, `Debuginfod options:
      --debuginfod                  Fetch debug symbols from debuginfod servers
      --debuginfod-cache string     Debuginfod cache directory (default "%s")
      --debuginfod-retries int      Debuginfod max retries per server (default %d)
      --debuginfod-servers string   Comma-separated debuginfod server URLs (default %q)
      --debuginfod-timeout duration Debuginfod HTTP timeout (default %v)
`, debuginfo.DefaultCacheDir(), debuginfo.DefaultRetries, debuginfo.DefaultServerURL, debuginfo.DefaultTimeout)
}

func parseRules(rulesFlag, targetPlatform, targetCompiler string) ([]string, error) {
	var ruleIDs []string
	if rulesFlag != "" {
		ruleIDs = strings.Split(rulesFlag, ",")
		for i, id := range ruleIDs {
			ruleIDs[i] = strings.TrimSpace(id)
		}
		for _, id := range ruleIDs {
			if rule.Get(id) == nil {
				return nil, fmt.Errorf("unknown rule %q", id)
			}
		}
	} else {
		ruleIDs = preset.DefaultRules
	}

	if targetPlatform != "" || targetCompiler != "" {
		filter, err := rule.ParseTargetFilter(targetPlatform, targetCompiler)
		if err != nil {
			return nil, err
		}
		ruleIDs = rule.FilterRules(ruleIDs, filter)
		if len(ruleIDs) == 0 {
			return nil, fmt.Errorf("no rules match the specified target filter")
		}
	}

	return ruleIDs, nil
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

	fs := flag.NewFlagSet("analyze", flag.ExitOnError)

	var (
		rulesFlag         string
		targetPlatform    string
		targetCompiler    string
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
	fs.StringVar(&targetPlatform, "target-platform", "", "")
	fs.StringVar(&targetCompiler, "target-compiler", "", "")
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

	ruleIDs, err := parseRules(rulesFlag, targetPlatform, targetCompiler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitError
	}

	paths, err := parsePaths(fs, inputFile)
	if err != nil {
		if errors.Is(err, errNoPathsSpecified) {
			fs.Usage()
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		return ExitError
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

	opts := outputOptions{
		aggregate:      aggregate,
		includePassed:  includePassed,
		includeSkipped: includeSkipped,
		exitZero:       exitZero,
		sarifOutput:    sarifOutput,
	}

	if needsFullReport {
		return a.processFullReport(resultsChan, opts, invocation)
	}
	return a.processStreaming(resultsChan, opts)
}

func (a *App) processFullReport(resultsChan <-chan analyzer.Result, opts outputOptions, invocation *output.InvocationInfo) int {
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

	if opts.aggregate {
		agg := output.AggregateFindings(report)
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

func (a *App) processStreaming(resultsChan <-chan analyzer.Result, opts outputOptions) int {
	var totalFailed int
	textFormatter, _ := output.GetFormatter("text", output.FormatterOptions{IncludePassed: opts.includePassed, IncludeSkipped: opts.includeSkipped})

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

	if totalFailed > 0 && !opts.exitZero {
		return ExitFindings
	}
	return ExitSuccess
}
