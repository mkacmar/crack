package cli

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/debuginfo"
	"github.com/mkacmar/crack/internal/output"
	"github.com/mkacmar/crack/internal/preset"
	"github.com/mkacmar/crack/internal/result"
	"github.com/mkacmar/crack/internal/rules"
	"github.com/mkacmar/crack/internal/rules/elf"
	"github.com/mkacmar/crack/internal/scanner"
	"github.com/mkacmar/crack/internal/version"
)

type App struct {
	logger *slog.Logger
}

func New() *App {
	return &App{}
}

func (a *App) Run(args []string) int {
	if len(args) < 2 {
		a.printUsage(args[0])
		return 1
	}

	cmd := args[1]
	switch cmd {
	case "analyze":
		return a.runAnalyze(args[0], args[2:])
	case "version", "-v", "--version":
		a.printVersion()
		return 0
	case "help", "-h", "--help":
		a.printUsage(args[0])
		return 0
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		a.printUsage(args[0])
		return 1
	}
}

func (a *App) printVersion() {
	fmt.Printf("crack %s\n", version.Version)
	if version.GitCommit != "unknown" {
		fmt.Printf("  commit: %s\n", version.GitCommit)
	}
	if version.BuildTime != "unknown" {
		fmt.Printf("  built:  %s\n", version.BuildTime)
	}
}

func (a *App) printUsage(prog string) {
	fmt.Fprintf(os.Stderr, "CRACK - Compiler Hardening Checker\n\n")
	fmt.Fprintf(os.Stderr, "Usage: %s <command> [options]\n\n", prog)
	fmt.Fprintf(os.Stderr, "Commands:\n")
	fmt.Fprintf(os.Stderr, "  analyze    Analyze binaries for security hardening features\n")
	fmt.Fprintf(os.Stderr, "  version    Show version information\n")
	fmt.Fprintf(os.Stderr, "  help       Show this help message\n")
	fmt.Fprintf(os.Stderr, "\nRun '%s <command> -h' for more information on a command.\n", prog)
}

func (a *App) printAnalyzeUsage(prog string) {
	fmt.Fprintf(os.Stderr, "Usage: %s analyze [options] <binary|directory>...\n\n", prog)
	fmt.Fprintf(os.Stderr, "Analyze binaries for security hardening features.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -P, --preset string         Security preset to use (default \"recommended\")\n")
	fmt.Fprintf(os.Stderr, "  -R, --rules string          Comma-separated list of rule IDs to run (mutually exclusive with --preset)\n")
	fmt.Fprintf(os.Stderr, "  -o, --sarif string          Save detailed SARIF report to file\n")
	fmt.Fprintf(os.Stderr, "  -a, --aggregate             Aggregate findings into actionable recommendations\n")
	fmt.Fprintf(os.Stderr, "  -r, --recursive             Recursively scan directories\n")
	fmt.Fprintf(os.Stderr, "      --log-level string      Log level: none, debug, info, warn, error (default \"error\")\n")
	fmt.Fprintf(os.Stderr, "      --show-passed           Show passing checks in output\n")
	fmt.Fprintf(os.Stderr, "      --show-skipped          Show skipped checks in output\n")
	fmt.Fprintf(os.Stderr, "  -p, --parallel int          Number of files to analyze in parallel (default %d)\n", runtime.NumCPU())
	fmt.Fprintf(os.Stderr, "\nDebuginfod options:\n")
	fmt.Fprintf(os.Stderr, "  -d, --debuginfod            Fetch debug symbols from debuginfod servers\n")
	fmt.Fprintf(os.Stderr, "      --debuginfod-urls       Debuginfod server URLs (default %q)\n", debuginfo.DefaultServerURL)
	fmt.Fprintf(os.Stderr, "      --debuginfod-cache      Debuginfod cache directory (default \"%s\")\n", getDefaultCacheDir())
	fmt.Fprintf(os.Stderr, "      --debuginfod-timeout    Debuginfod HTTP timeout (default %v)\n", debuginfo.DefaultTimeout)
	fmt.Fprintf(os.Stderr, "      --debuginfod-retries    Debuginfod max retries per server (default %d)\n", debuginfo.DefaultRetries)
	fmt.Fprintf(os.Stderr, "\nPresets:\n")
	for _, name := range preset.Names() {
		if name == "recommended" {
			fmt.Fprintf(os.Stderr, "  %s (default)\n", name)
		} else {
			fmt.Fprintf(os.Stderr, "  %s\n", name)
		}
	}
	fmt.Fprintf(os.Stderr, "\nUse 'crack analyze --preset=<name> --list-rules' to see rules in a preset.\n")
}

func (a *App) runAnalyze(prog string, args []string) int {
	fs := flag.NewFlagSet("analyze", flag.ExitOnError)

	var (
		presetName        string
		presetSet         bool
		rulesFlag         string
		listRules         bool
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

	fs.StringVar(&presetName, "preset", "recommended", "")
	fs.StringVar(&presetName, "P", "recommended", "")
	fs.StringVar(&rulesFlag, "rules", "", "")
	fs.StringVar(&rulesFlag, "R", "", "")
	fs.BoolVar(&listRules, "list-rules", false, "")
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

	// Check if --preset was explicitly set
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "preset" || f.Name == "P" {
			presetSet = true
		}
	})

	if rulesFlag != "" && presetSet {
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
			if elf.GetRuleByID(id) == nil {
				fmt.Fprintf(os.Stderr, "Error: unknown rule %q\n", id)
				fmt.Fprintf(os.Stderr, "Use 'crack analyze --preset=hardened --list-rules' to see available rules.\n")
				return 1
			}
		}
		p = preset.Preset{Rules: ruleIDs}
	} else {
		var ok bool
		p, ok = preset.Get(presetName)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unknown preset %q\n", presetName)
			fmt.Fprintf(os.Stderr, "Available presets: %s\n", strings.Join(preset.Names(), ", "))
			return 1
		}
	}

	if listRules {
		var general, x86, arm []string
		for _, ruleID := range p.Rules {
			rule := elf.GetRuleByID(ruleID)
			if rule == nil {
				general = append(general, ruleID)
				continue
			}
			arch := rule.Applicability().Arch
			if arch.Matches(binary.ArchAllX86) && !arch.Matches(binary.ArchAllARM) {
				x86 = append(x86, ruleID)
			} else if arch.Matches(binary.ArchAllARM) && !arch.Matches(binary.ArchAllX86) {
				arm = append(arm, ruleID)
			} else {
				general = append(general, ruleID)
			}
		}

		sort.Strings(general)
		sort.Strings(x86)
		sort.Strings(arm)

		if len(general) > 0 {
			fmt.Println("General:")
			for _, ruleID := range general {
				rule := elf.GetRuleByID(ruleID)
				if rule != nil {
					fmt.Printf("  %-24s %s\n", ruleID, rule.Name())
				} else {
					fmt.Printf("  %-24s (unknown)\n", ruleID)
				}
			}
		}

		if len(x86) > 0 {
			if len(general) > 0 {
				fmt.Println()
			}
			fmt.Println("x86:")
			for _, ruleID := range x86 {
				rule := elf.GetRuleByID(ruleID)
				fmt.Printf("  %-24s %s\n", ruleID, rule.Name())
			}
		}

		if len(arm) > 0 {
			if len(general) > 0 || len(x86) > 0 {
				fmt.Println()
			}
			fmt.Println("ARM:")
			for _, ruleID := range arm {
				rule := elf.GetRuleByID(ruleID)
				fmt.Printf("  %-24s %s\n", ruleID, rule.Name())
			}
		}

		return 0
	}

	if fs.NArg() == 0 {
		fs.Usage()
		return 1
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

	ruleEngine := rules.NewEngine(a.logger)
	a.logger.Debug("loading preset", slog.String("preset", presetName))
	ruleEngine.LoadPreset(p)

	cache := debuginfodCache
	if cache == "" {
		cache = getDefaultCacheDir()
	}

	scannerOpts := scanner.Options{
		Logger:            a.logger,
		Workers:           parallel,
		UseDebuginfod:     useDebuginfod,
		DebuginfodURLs:    parseURLList(debuginfodURLs),
		DebuginfodCache:   cache,
		DebuginfodTimeout: debuginfodTimeout,
		DebuginfodRetries: debuginfodRetries,
	}
	scan := scanner.NewScanner(ruleEngine, scannerOpts)

	ctx := context.Background()
	paths := fs.Args()

	a.logger.Info("starting scan", slog.Int("paths", len(paths)), slog.Bool("recursive", recursive))

	resultsChan := scan.ScanPaths(ctx, paths, recursive)

	needsFullReport := aggregate || sarifOutput != ""

	if needsFullReport {
		return a.processFullReport(resultsChan, aggregate, showPassed, showSkipped, sarifOutput)
	}
	return a.processStreaming(resultsChan, showPassed, showSkipped)
}

func (a *App) processFullReport(resultsChan <-chan result.FileScanResult, aggregate, showPassed, showSkipped bool, sarifOutput string) int {
	var results []result.FileScanResult
	var totalFailed int

	for res := range resultsChan {
		if res.Skipped {
			continue
		}
		results = append(results, res)
		totalFailed += res.FailedChecks()
	}

	report := &result.ScanResults{Results: results}

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
		sarifFormatter, _ := output.GetFormatter("sarif", output.FormatterOptions{ShowPassed: true})
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

func (a *App) processStreaming(resultsChan <-chan result.FileScanResult, showPassed, showSkipped bool) int {
	var totalFailed int
	textFormatter, _ := output.GetFormatter("text", output.FormatterOptions{ShowPassed: showPassed, ShowSkipped: showSkipped})

	for res := range resultsChan {
		if res.Skipped {
			continue
		}
		totalFailed += res.FailedChecks()
		singleReport := &result.ScanResults{Results: []result.FileScanResult{res}}
		if err := textFormatter.Format(singleReport, os.Stdout); err != nil {
			a.logger.Error("failed to format output", slog.Any("error", err))
		}
	}

	if totalFailed > 0 {
		return 1
	}
	return 0
}

func setupLogger(level string) (*slog.Logger, bool) {
	var logLevel slog.Level
	valid := true

	switch level {
	case "none":
		logLevel = slog.Level(99)
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelError
		valid = false
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewTextHandler(os.Stderr, opts)
	return slog.New(handler), valid
}

func getDefaultCacheDir() string {
	cacheDir, _ := os.UserCacheDir()
	return filepath.Join(cacheDir, "crack", "debuginfo")
}

func parseURLList(s string) []string {
	var urls []string
	for _, url := range strings.Split(s, ",") {
		if url = strings.TrimSpace(url); url != "" {
			urls = append(urls, url)
		}
	}
	return urls
}
