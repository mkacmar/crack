package output

import (
	"fmt"
	"io"

	"go.kacmar.sk/crack/rule"
)

type TextFormatter struct {
	IncludePassed  bool
	IncludeSkipped bool
}

func (f *TextFormatter) Format(report *DecoratedReport, w io.Writer) error {
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
