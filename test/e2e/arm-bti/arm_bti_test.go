package arm_bti_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestARMBTIRule(t *testing.T) {
	e2e.RunRuleTests(t, "arm-bti", []e2e.TestCase{
		{Binary: "gcc-bti-enabled", Expect: "pass"},
		{Binary: "gcc-bti-disabled", Expect: "fail"},
		{Binary: "gcc-bti-stripped", Expect: "pass"},
		{Binary: "clang-bti-enabled", Expect: "pass"},
		{Binary: "clang-bti-disabled", Expect: "fail"},
		{Binary: "clang-bti-stripped", Expect: "pass"},
		{Binary: "musl-gcc-bti-enabled", Expect: "pass"},
		{Binary: "musl-gcc-bti-disabled", Expect: "fail"},
		{Binary: "musl-gcc-bti-stripped", Expect: "pass"},
		{Binary: "musl-clang-bti-enabled", Expect: "pass"},
		{Binary: "musl-clang-bti-disabled", Expect: "fail"},
		{Binary: "musl-clang-bti-stripped", Expect: "pass"},
	})
}
