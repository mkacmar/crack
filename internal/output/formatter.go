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
	Verbose bool
}

func (f *TextFormatter) Format(report *model.ScanResults, w io.Writer) error {
	for _, result := range report.Results {
		if result.Error != nil {
			fmt.Fprintf(w, "%s: ERROR: %v\n", result.Path, result.Error)
			continue
		}

		for _, check := range result.Results {
			switch check.State {
			case model.CheckStatePassed:
				if f.Verbose {
					fmt.Fprintf(w, "%s: %s [pass]: %s\n", result.Path, check.RuleID, check.Message)
				}
			case model.CheckStateFailed:
				if check.Suggestion != "" {
					fmt.Fprintf(w, "%s: %s [fail]: %s %s\n", result.Path, check.RuleID, check.Message, check.Suggestion)
				} else {
					fmt.Fprintf(w, "%s: %s [fail]: %s\n", result.Path, check.RuleID, check.Message)
				}
			}
		}
	}

	return nil
}

var formatters = map[string]func(bool) Formatter{
	"text":  func(verbose bool) Formatter { return &TextFormatter{Verbose: verbose} },
	"":      func(verbose bool) Formatter { return &TextFormatter{Verbose: verbose} },
	"sarif": func(verbose bool) Formatter { return &SARIFFormatter{IncludePassed: verbose} },
}

func GetFormatter(format string, verbose bool) (Formatter, error) {
	if constructor, ok := formatters[format]; ok {
		return constructor(verbose), nil
	}
	return nil, fmt.Errorf("unsupported format: %s (supported: text, sarif)", format)
}
