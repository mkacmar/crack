package safe_stack_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestSafeStackRule(t *testing.T) {
	e2e.RunRuleTests(t, "safe-stack", []e2e.TestCase{
		{Binary: "amd64-clang-safestack", Expect: e2e.Pass},
		{Binary: "amd64-clang-safestack-stripped", Expect: e2e.Pass},
		{Binary: "amd64-clang-no-safestack", Expect: e2e.Fail},
		{Binary: "amd64-gcc-no-safestack", Expect: e2e.Skip},
		{Binary: "amd64-rustc-no-safestack", Expect: e2e.Skip},

		{Binary: "arm64-clang-safestack", Expect: e2e.Pass},
		{Binary: "arm64-clang-safestack-stripped", Expect: e2e.Pass},
		{Binary: "arm64-clang-no-safestack", Expect: e2e.Fail},
		{Binary: "arm64-gcc-no-safestack", Expect: e2e.Skip},
		{Binary: "arm64-rustc-no-safestack", Expect: e2e.Skip},
	})
}
