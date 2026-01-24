package arm_pac_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestARMPACRule(t *testing.T) {
	e2e.RunRuleTests(t, "arm-pac", []e2e.TestCase{
		{Binary: "gcc-pac-enabled", Expect: "pass"},
		{Binary: "gcc-pac-disabled", Expect: "fail"},
		{Binary: "gcc-pac-stripped", Expect: "pass"},
		{Binary: "clang-pac-enabled", Expect: "pass"},
		{Binary: "clang-pac-disabled", Expect: "fail"},
		{Binary: "clang-pac-stripped", Expect: "pass"},
	})
}
