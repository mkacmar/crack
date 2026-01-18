package output

import (
	"fmt"
	"io"

	"github.com/mkacmar/crack/internal/model"
)

type Formatter interface {
	Format(report *model.ScanResults, w io.Writer) error
}

type TextFormatter struct {
	ShowPassed  bool
	ShowSkipped bool
}

func (f *TextFormatter) Format(report *model.ScanResults, w io.Writer) error {
	for _, result := range report.Results {
		if result.Error != nil {
			fmt.Fprintf(w, "ERROR = %s: %v\n", result.Path, result.Error)
			continue
		}

		for _, check := range result.Results {
			switch check.State {
			case model.CheckStatePassed:
				if f.ShowPassed {
					fmt.Fprintf(w, "PASS = %s @ %s: %s\n", check.RuleID, result.Path, check.Message)
				}
			case model.CheckStateFailed:
				if check.Suggestion != "" {
					fmt.Fprintf(w, "FAIL = %s @ %s: %s %s\n", check.RuleID, result.Path, check.Message, check.Suggestion)
				} else {
					fmt.Fprintf(w, "FAIL = %s @ %s: %s\n", check.RuleID, result.Path, check.Message)
				}
			case model.CheckStateSkipped:
				if f.ShowSkipped {
					fmt.Fprintf(w, "SKIP = %s @ %s: %s\n", check.RuleID, result.Path, check.Message)
				}
			}
		}
	}

	return nil
}

type FormatterOptions struct {
	ShowPassed  bool
	ShowSkipped bool
}

var formatters = map[string]func(FormatterOptions) Formatter{
	"text": func(opts FormatterOptions) Formatter {
		return &TextFormatter{ShowPassed: opts.ShowPassed, ShowSkipped: opts.ShowSkipped}
	},
	"": func(opts FormatterOptions) Formatter {
		return &TextFormatter{ShowPassed: opts.ShowPassed, ShowSkipped: opts.ShowSkipped}
	},
	"sarif": func(opts FormatterOptions) Formatter { return &SARIFFormatter{IncludePassed: opts.ShowPassed} },
}

func GetFormatter(format string, opts FormatterOptions) (Formatter, error) {
	if constructor, ok := formatters[format]; ok {
		return constructor(opts), nil
	}
	return nil, fmt.Errorf("unsupported format: %s (supported: text, sarif)", format)
}
