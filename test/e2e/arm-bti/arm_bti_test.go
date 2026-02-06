package arm_bti_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestARMBTIRule(t *testing.T) {
	e2e.RunRuleTests(t, "arm-bti", []e2e.TestCase{
		{Binary: "gcc-bti-enabled", Expect: e2e.Pass},
		{Binary: "gcc-bti-disabled", Expect: e2e.Fail},
		{Binary: "gcc-bti-stripped", Expect: e2e.Pass},
		{Binary: "clang-bti-enabled", Expect: e2e.Pass},
		{Binary: "clang-bti-disabled", Expect: e2e.Fail},
		{Binary: "clang-bti-stripped", Expect: e2e.Pass},
		{Binary: "musl-gcc-bti-enabled", Expect: e2e.Pass},
		{Binary: "musl-gcc-bti-disabled", Expect: e2e.Fail},
		{Binary: "musl-gcc-bti-stripped", Expect: e2e.Pass},
		{Binary: "musl-clang-bti-enabled", Expect: e2e.Pass},
		{Binary: "musl-clang-bti-disabled", Expect: e2e.Fail},
		{Binary: "musl-clang-bti-stripped", Expect: e2e.Pass},

		{Binary: "rustc-no-bti", Expect: e2e.Skip},
		{Binary: "musl-rustc-no-bti", Expect: e2e.Skip},
	})
}
