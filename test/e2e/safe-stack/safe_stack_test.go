package safe_stack_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestSafeStackRule(t *testing.T) {
	e2e.RunRuleTests(t, "safe-stack", []e2e.TestCase{
		{Binary: "x86_64-clang-safestack", Expect: e2e.Pass},
		{Binary: "x86_64-clang-safestack-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-clang-no-safestack", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-no-safestack", Expect: e2e.Fail},

		{Binary: "aarch64-clang-safestack", Expect: e2e.Pass},
		{Binary: "aarch64-clang-safestack-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-clang-no-safestack", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-no-safestack", Expect: e2e.Fail},
	})
}
