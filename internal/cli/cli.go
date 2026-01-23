package cli

import (
	"fmt"
	"log/slog"
	"os"

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
	case "list-rules":
		return a.runListRules(args[0], args[2:])
	case "version", "-v", "--version":
		a.printVersion()
		return 0
	case "help", "-h", "--help":
		a.printUsage(args[0])
		return 0
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command: %s\n\n", cmd)
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
	fmt.Fprintf(os.Stderr, `CRACK - Compiler Hardening Checker

Usage: %s <command> [options]

Commands:
  analyze      Analyze binaries for security hardening features
  list-rules   List available security rules
  version      Show version information
  help         Show this help message

Run '%s <command> -h' for more information on a command.
`, prog, prog)
}
