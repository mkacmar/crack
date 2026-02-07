package arm_mte_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestARMMTERule(t *testing.T) {
	e2e.RunRuleTests(t, "arm-mte", []e2e.TestCase{
		{Binary: "clang-mte", Expect: e2e.Pass},
		{Binary: "clang-mte-stripped", Expect: e2e.Pass},
		{Binary: "clang-no-mte", Expect: e2e.Fail},
	})
}
