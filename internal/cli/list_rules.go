package cli

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/preset"
	"github.com/mkacmar/crack/internal/rule"
)

func (a *App) printListRulesUsage(prog string) {
	fmt.Fprintf(os.Stderr, `Usage: %s list-rules [options]

List available security rules.

Options:
  -P, --preset string    Show rules for a specific preset (default %q)

Presets:
`, prog, preset.Default)

	for _, name := range preset.Names() {
		if name == preset.Default {
			fmt.Fprintf(os.Stderr, "  %s (default)\n", name)
		} else {
			fmt.Fprintf(os.Stderr, "  %s\n", name)
		}
	}
}

func (a *App) runListRules(prog string, args []string) int {
	fs := flag.NewFlagSet("list-rules", flag.ExitOnError)

	var presetName string
	fs.StringVar(&presetName, "preset", preset.Default, "")
	fs.StringVar(&presetName, "P", preset.Default, "")

	fs.Usage = func() {
		a.printListRulesUsage(prog)
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	p, ok := preset.Get(presetName)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: unknown preset %q\n", presetName)
		fmt.Fprintf(os.Stderr, "Available presets: %s\n", strings.Join(preset.Names(), ", "))
		return 1
	}

	var general, x86, arm []string
	for _, ruleID := range p.Rules {
		r := rule.Get(ruleID)
		if r == nil {
			general = append(general, ruleID)
			continue
		}
		arch := r.Applicability().Platform.Architecture
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
			r := rule.Get(ruleID)
			if r != nil {
				fmt.Printf("  %-24s %s\n", ruleID, r.Name())
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
			r := rule.Get(ruleID)
			fmt.Printf("  %-24s %s\n", ruleID, r.Name())
		}
	}

	if len(arm) > 0 {
		if len(general) > 0 || len(x86) > 0 {
			fmt.Println()
		}
		fmt.Println("ARM:")
		for _, ruleID := range arm {
			r := rule.Get(ruleID)
			fmt.Printf("  %-24s %s\n", ruleID, r.Name())
		}
	}

	return 0
}
