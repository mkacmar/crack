package profile

import (
	"testing"

	"github.com/mkacmar/crack/internal/rules/elf"
)

func TestAllRulesInProfiles(t *testing.T) {
	rulesInProfiles := make(map[string]bool)
	for _, profile := range All {
		for _, ruleID := range profile.Rules {
			rulesInProfiles[ruleID] = true
		}
	}

	var missingRules []string
	for _, rule := range elf.AllRules {
		if !rulesInProfiles[rule.ID()] {
			missingRules = append(missingRules, rule.ID())
		}
	}

	if len(missingRules) > 0 {
		t.Errorf("the following rules are not in any profile: %v", missingRules)
	}
}
