package output

import (
	"fmt"
	"io"

	"github.com/mkacmar/crack/internal/analyzer"
	"github.com/mkacmar/crack/rule"
)

type Formatter interface {
	Format(report *analyzer.Report, w io.Writer) error
}

type TextFormatter struct {
	IncludePassed  bool
	IncludeSkipped bool
}

func (f *TextFormatter) Format(report *analyzer.Report, w io.Writer) error {
	for _, result := range report.Results {
		if result.Error != nil {
			fmt.Fprintf(w, "ERROR = %s: %v\n", result.Path, result.Error)
			continue
		}

		for _, finding := range result.Findings {
			switch finding.Status {
			case rule.StatusPassed:
				if f.IncludePassed {
					fmt.Fprintf(w, "PASS = %s @ %s: %s\n", finding.RuleID, result.Path, finding.Message)
				}
			case rule.StatusFailed:
				if finding.Suggestion != "" {
					fmt.Fprintf(w, "FAIL = %s @ %s: %s %s\n", finding.RuleID, result.Path, finding.Message, finding.Suggestion)
				} else {
					fmt.Fprintf(w, "FAIL = %s @ %s: %s\n", finding.RuleID, result.Path, finding.Message)
				}
			case rule.StatusSkipped:
				if f.IncludeSkipped {
					fmt.Fprintf(w, "SKIP = %s @ %s: %s\n", finding.RuleID, result.Path, finding.Message)
				}
			}
		}
	}

	return nil
}

type FormatterOptions struct {
	IncludePassed  bool
	IncludeSkipped bool
	Invocation     *InvocationInfo
}

var formatters = map[string]func(FormatterOptions) Formatter{
	"text": func(opts FormatterOptions) Formatter {
		return &TextFormatter{IncludePassed: opts.IncludePassed, IncludeSkipped: opts.IncludeSkipped}
	},
	"": func(opts FormatterOptions) Formatter {
		return &TextFormatter{IncludePassed: opts.IncludePassed, IncludeSkipped: opts.IncludeSkipped}
	},
	"sarif": func(opts FormatterOptions) Formatter {
		return &SARIFFormatter{
			IncludePassed:  opts.IncludePassed,
			IncludeSkipped: opts.IncludeSkipped,
			Invocation:     opts.Invocation,
		}
	},
}

func GetFormatter(format string, opts FormatterOptions) (Formatter, error) {
	if constructor, ok := formatters[format]; ok {
		return constructor(opts), nil
	}
	return nil, fmt.Errorf("unsupported format: %s (supported: text, sarif)", format)
}
