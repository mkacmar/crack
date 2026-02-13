package rule

import "github.com/mkacmar/crack/binary"

// Check runs all rules against binary metadata, using execFn to execute each rule.
func Check[R Rule](rules []R, info binary.Info, execFn func(R) Result) []Finding {
	findings := make([]Finding, 0, len(rules))

	for _, r := range rules {
		if reason := CheckApplicability(r.Applicability(), info); reason != Applicable {
			findings = append(findings, Finding{
				Result: Result{
					Status:  StatusSkipped,
					Message: reason.SkipMessage(info),
				},
				RuleID: r.ID(),
				Name:   r.Name(),
			})
			continue
		}

		result := execFn(r)
		findings = append(findings, Finding{
			Result: result,
			RuleID: r.ID(),
			Name:   r.Name(),
		})
	}

	return findings
}
