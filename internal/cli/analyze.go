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

	"go.kacmar.sk/crack/internal/analyzer"
	"go.kacmar.sk/crack/internal/debuginfo"
	"go.kacmar.sk/crack/internal/output"
	"go.kacmar.sk/crack/internal/preset"
	"go.kacmar.sk/crack/internal/scanner"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/rule/registry"
	"go.kacmar.sk/debuginfod/cache"
)

var errNoPathsSpecified = fmt.Errorf("no paths specified")

// defaultDebuginfodServer is the public elfutils debuginfod server.
const defaultDebuginfodServer = "https://debuginfod.elfutils.org"

type outputOptions struct {
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
	useLocalDebuginfo bool
	localDebuginfoDir string
	profile           profileConfig
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
      --exit-zero             Exit with 0 even when findings are detected
      --include-passed        Include passing checks in output
      --include-skipped       Include skipped checks in output
      --sarif string          Save detailed SARIF report to file

Logging options:
      --log string            Log output file (default stderr)
      --log-level string      Log level: none, debug, info, warn, error (default "error")

`)

	defaultCacheDir, err := debuginfo.DefaultCacheDir()
	if err != nil {
		defaultCacheDir = "(unavailable)"
	}
	fmt.Fprintf(os.Stderr, `Debuginfod options:
      --debuginfod                    Fetch debug symbols from debuginfod servers
      --debuginfod-cache-dir string   Debuginfod cache directory (default "%s")
      --debuginfod-retries int        Debuginfod max retries per server (default %d)
      --debuginfod-servers string     Comma-separated debuginfod server URLs (default %q)
      --debuginfod-timeout duration   Debuginfod HTTP timeout (default %v)
`, defaultCacheDir, debuginfo.DefaultMaxRetries, defaultDebuginfodServer, 30*time.Second)

	fmt.Fprintf(os.Stderr, `
Local debuginfo:
      --local-debuginfo                Resolve missing sections from a local build-id-indexed debug directory
      --local-debuginfo-dir string     Root directory of the local debuginfo store (default %q)
`, debuginfo.DefaultBuildIDDir)

	if usage := profileUsage(); usage != "" {
		fmt.Fprint(os.Stderr, usage)
	}
}

func parseRules(rulesFlag, targetPlatform, targetCompiler string) ([]rule.ELFRule, error) {
	var selectedRules []rule.ELFRule
	if rulesFlag != "" {
		ids := strings.Split(rulesFlag, ",")
		for _, id := range ids {
			id = strings.TrimSpace(id)
			r, ok := registry.Find[rule.ELFRule](registry.ByID(id))
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

	stopProfile, err := startProfiling(&cfg.profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitError
	}
	if stopProfile != nil {
		defer stopProfile()
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

	closeLog, err := a.setupLogging(cfg.logFile, cfg.logLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitError
	}
	defer closeLog()

	debuginfodCache, err := a.setupDebuginfod(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitError
	}

	elfAnalyzer := analyzer.NewELFAnalyzer(analyzer.ELFAnalyzerOptions{
		Rules:   selectedRules,
		Sources: a.buildDebuginfoSources(cfg, debuginfodCache),
		Logger:  a.logger,
	})

	dispatcher := analyzer.NewDispatcher(analyzer.DispatcherOptions{
		ELF:    elfAnalyzer,
		Logger: a.logger,
	})

	scan := scanner.NewScanner(dispatcher, scanner.Options{
		Logger:  a.logger,
		Workers: cfg.parallel,
	})

	ctx := cancelOnSignal(context.Background())
	a.logger.Info("starting scan", slog.Int("paths", len(paths)), slog.Bool("recursive", cfg.recursive))

	resultsChan := scan.ScanPaths(ctx, paths, cfg.recursive)

	invocation := &output.InvocationInfo{
		CommandLine: strings.Join(append([]string{prog}, args...), " "),
		Arguments:   args,
		StartTime:   startTime,
		WorkingDir:  workingDir,
	}

	if opts.sarifOutput != "" {
		return a.processFullReport(resultsChan, opts, invocation)
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
	fs.BoolVar(&cfg.recursive, "recursive", false, "")
	fs.StringVar(&cfg.logFile, "log", "", "")
	fs.StringVar(&cfg.logLevel, "log-level", "error", "")
	fs.BoolVar(&opts.includePassed, "include-passed", false, "")
	fs.BoolVar(&opts.includeSkipped, "include-skipped", false, "")
	fs.IntVar(&cfg.parallel, "parallel", runtime.NumCPU(), "")
	fs.BoolVar(&opts.exitZero, "exit-zero", false, "")
	fs.BoolVar(&cfg.useDebuginfod, "debuginfod", false, "")
	fs.StringVar(&cfg.debuginfodServers, "debuginfod-servers", defaultDebuginfodServer, "")
	fs.StringVar(&cfg.debuginfodCache, "debuginfod-cache-dir", "", "")
	fs.DurationVar(&cfg.debuginfodTimeout, "debuginfod-timeout", 30*time.Second, "")
	fs.IntVar(&cfg.debuginfodRetries, "debuginfod-retries", debuginfo.DefaultMaxRetries, "")
	fs.BoolVar(&cfg.useLocalDebuginfo, "local-debuginfo", false, "")
	fs.StringVar(&cfg.localDebuginfoDir, "local-debuginfo-dir", "", "")
	registerProfileFlags(fs, &cfg.profile)

	fs.Usage = func() { a.printAnalyzeUsage(prog) }

	return fs, opts, cfg
}

func (a *App) setupLogging(logFile, logLevel string) (func(), error) {
	var logOutput io.Writer = os.Stderr
	cleanup := func() {}
	if logFile != "" {
		f, err := os.Create(logFile) // #nosec G304 -- user-provided log file path
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		logOutput = f
		cleanup = func() { _ = f.Close() }
	}
	a.logger = setupLogger(logLevel, logOutput)
	return cleanup, nil
}

func (a *App) setupDebuginfod(cfg *analyzeConfig) (*cache.DiskCache, error) {
	if !cfg.useDebuginfod {
		return nil, nil
	}
	servers := strings.Split(cfg.debuginfodServers, ",")
	var filtered []string
	for _, s := range servers {
		if s = strings.TrimSpace(s); s != "" {
			filtered = append(filtered, s)
		}
	}

	return debuginfo.NewCache(debuginfo.Options{
		ServerURLs: filtered,
		CacheDir:   cfg.debuginfodCache,
		Timeout:    cfg.debuginfodTimeout,
		MaxRetries: cfg.debuginfodRetries,
		Logger:     a.logger,
	})
}

// buildDebuginfoSources assembles the configured debug-information sources in priority order.
func (a *App) buildDebuginfoSources(cfg *analyzeConfig, debuginfodCache *cache.DiskCache) []debuginfo.Source {
	var sources []debuginfo.Source
	if cfg.useLocalDebuginfo {
		root := cfg.localDebuginfoDir
		if root == "" {
			root = debuginfo.DefaultBuildIDDir
		}
		sources = append(sources, debuginfo.NewBuildIDDirSource(root, a.logger))
	}
	if debuginfodCache != nil {
		sources = append(sources, debuginfo.NewDebuginfodSource(debuginfodCache, a.logger))
	}
	return sources
}
