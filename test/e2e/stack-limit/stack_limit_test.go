package stacklimit_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestStackLimitRule(t *testing.T) {
	e2e.RunRuleTests(t, "stack-limit", []e2e.TestCase{
		{Binary: "x86_64-gcc-stack-limit", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-no-stack-limit", Expect: e2e.Fail},
		{Binary: "x86_64-clang-stack-limit", Expect: e2e.Pass},
		{Binary: "x86_64-clang-no-stack-limit", Expect: e2e.Fail},

		{Binary: "aarch64-gcc-stack-limit", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-no-stack-limit", Expect: e2e.Fail},
		{Binary: "aarch64-clang-stack-limit", Expect: e2e.Pass},
		{Binary: "aarch64-clang-no-stack-limit", Expect: e2e.Fail},

		{Binary: "armv7-gcc-stack-limit", Expect: e2e.Pass},
		{Binary: "armv7-gcc-no-stack-limit", Expect: e2e.Fail},
		{Binary: "armv7-clang-stack-limit", Expect: e2e.Pass},
		{Binary: "armv7-clang-no-stack-limit", Expect: e2e.Fail},
	})
}
