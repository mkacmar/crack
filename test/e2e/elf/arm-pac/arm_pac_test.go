package arm_pac_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestARMPACRule(t *testing.T) {
	e2e.RunRuleTests(t, "arm-pac", []e2e.TestCase{
		{Binary: "gcc-pac-enabled", Expect: e2e.Pass},
		{Binary: "gcc-pac-disabled", Expect: e2e.Fail},
		{Binary: "gcc-pac-stripped", Expect: e2e.Pass},
		{Binary: "clang-pac-enabled", Expect: e2e.Pass},
		{Binary: "clang-pac-disabled", Expect: e2e.Fail},
		{Binary: "clang-pac-stripped", Expect: e2e.Pass},

		{Binary: "rustc-no-pac", Expect: e2e.Skip},
	})
}
