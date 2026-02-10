package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"reflect"
	"sort"
	"text/template"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/internal/rules"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

//go:embed rules.md.tpl
var templateContent string

var docTemplate = template.Must(template.New("doc").Parse(templateContent))

type docData struct {
	Rules []ruleData
}

type ruleData struct {
	ID          string
	Name        string
	StructName  string
	Description string
	Platform    string
	Compilers   []compilerData
}

type compilerData struct {
	Name           string
	MinVersion     string
	DefaultVersion string
	Flag           string
}

func main() {
	allRules := rules.All()
	sort.Slice(allRules, func(i, j int) bool {
		return allRules[i].ID() < allRules[j].ID()
	})

	doc, err := generateDoc(allRules)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error generating documentation: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(doc)
}

func generateDoc(rules []rule.Rule) (string, error) {
	data := docData{Rules: make([]ruleData, 0, len(rules))}

	for _, r := range rules {
		data.Rules = append(data.Rules, toRuleData(r))
	}

	var buf bytes.Buffer
	if err := docTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func toRuleData(r rule.Rule) ruleData {
	applicability := r.Applicability()

	compilers := make([]toolchain.Compiler, 0, len(applicability.Compilers))
	for c := range applicability.Compilers {
		compilers = append(compilers, c)
	}
	sort.Slice(compilers, func(i, j int) bool {
		return compilers[i].String() < compilers[j].String()
	})

	compilerList := make([]compilerData, 0, len(compilers))
	for _, c := range compilers {
		req := applicability.Compilers[c]
		defaultVer := "-"
		if req.DefaultVersion != (toolchain.Version{}) {
			defaultVer = req.DefaultVersion.String()
		}
		compilerList = append(compilerList, toCompilerData(c, req, defaultVer))
	}

	structName := reflect.TypeOf(r).Name()

	return ruleData{
		ID:          r.ID(),
		Name:        r.Name(),
		StructName:  structName,
		Description: r.Description(),
		Platform:    formatPlatform(applicability.Platform),
		Compilers:   compilerList,
	}
}

func toCompilerData(c toolchain.Compiler, req rule.CompilerRequirement, defaultVer string) compilerData {
	return compilerData{
		Name:           c.String(),
		MinVersion:     req.MinVersion.String(),
		DefaultVersion: defaultVer,
		Flag:           req.Flag,
	}
}

func formatPlatform(p binary.Platform) string {
	if p.MinISA == (binary.ISA{}) {
		return p.Architecture.String()
	}
	return fmt.Sprintf("%s (requires ISA %s+)", p.Architecture.String(), p.MinISA.String())
}
