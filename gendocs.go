//go:build ignore

package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	elfrules "github.com/mkacmar/crack/internal/rules/elf"
	"github.com/mkacmar/crack/internal/toolchain"
)

func main() {
	elfrules.RegisterRules()

	rules := rule.GetAll()
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].ID() < rules[j].ID()
	})

	fmt.Print(generateDoc(rules))
}

func generateDoc(rules []rule.Rule) string {
	var sb strings.Builder

	for i, r := range rules {
		sb.WriteString(generateRuleDoc(r))
		if i < len(rules)-1 {
			sb.WriteString("\n---\n\n")
		}
	}

	return sb.String()
}

func generateRuleDoc(r rule.Rule) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## %s (`%s`)\n\n", r.Name(), r.ID()))
	sb.WriteString(fmt.Sprintf("%s\n\n", r.Description()))

	applicability := r.Applicability()

	sb.WriteString("### Platform\n\n")
	sb.WriteString(formatPlatform(applicability.Platform))
	sb.WriteString("\n\n")

	sb.WriteString("### Toolchain\n\n")
	if len(applicability.Compilers) == 0 {
		sb.WriteString("No specific compiler requirements.\n")
	} else {
		sb.WriteString("| Compiler | Minimal Version | Default Version | Flag |\n")
		sb.WriteString("|:---------|:----------------|:----------------|:-----|\n")

		compilers := make([]toolchain.Compiler, 0, len(applicability.Compilers))
		for c := range applicability.Compilers {
			compilers = append(compilers, c)
		}
		sort.Slice(compilers, func(i, j int) bool {
			return compilers[i].String() < compilers[j].String()
		})

		for _, compiler := range compilers {
			req := applicability.Compilers[compiler]
			defaultVer := "-"
			if req.DefaultVersion != (toolchain.Version{}) {
				defaultVer = req.DefaultVersion.String()
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | `%s` |\n",
				compiler.String(),
				req.MinVersion.String(),
				defaultVer,
				req.Flag,
			))
		}
	}

	return sb.String()
}

func formatPlatform(p binary.Platform) string {
	arch := p.Architecture.String()
	switch p.Architecture {
	case binary.ArchAll:
		arch = "All architectures"
	case binary.ArchAllX86:
		arch = "x86, x86_64"
	case binary.ArchAllARM:
		arch = "ARM, ARM64"
	}
	if p.MinISA == (binary.ISA{}) {
		return arch
	}
	return fmt.Sprintf("%s (requires ISA %s+)", arch, p.MinISA.String())
}
