package preset

import (
	"testing"

	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/rule/registry"
)

func TestDefaultRulesAreRegistered(t *testing.T) {
	for _, id := range Default() {
		if _, ok := registry.Find[rule.Rule](registry.ByID(id)); !ok {
			t.Errorf("preset rule %q is not in the registry", id)
		}
	}
}
