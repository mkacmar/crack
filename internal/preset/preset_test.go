package preset

import (
	"testing"

	"github.com/mkacmar/crack/internal/rules/elf"
)

func TestAllRulesInPresets(t *testing.T) {
	rulesInPresets := make(map[string]bool)
	for _, preset := range All {
		for _, ruleID := range preset.Rules {
			rulesInPresets[ruleID] = true
		}
	}

	var missingRules []string
	for _, rule := range elf.AllRules {
		if !rulesInPresets[rule.ID()] {
			missingRules = append(missingRules, rule.ID())
		}
	}

	if len(missingRules) > 0 {
		t.Errorf("the following rules are not in any preset: %v", missingRules)
	}
}
