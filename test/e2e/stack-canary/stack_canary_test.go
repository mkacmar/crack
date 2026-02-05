package stack_canary_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestStackCanaryRule(t *testing.T) {
	e2e.RunRuleTests(t, "stack-canary", []e2e.TestCase{
		{Binary: "x86_64-gcc-stack-protector-strong", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-stack-protector-all", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-stack-protector", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-no-stack-protector", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-stack-protector-simple", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-stack-protector-all-simple", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-stack-protector-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-stack-protector-static", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-stack-protector-static-stripped", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-stack-protector-lto", Expect: e2e.Pass},

		{Binary: "x86_64-clang-stack-protector-strong", Expect: e2e.Pass},
		{Binary: "x86_64-clang-stack-protector-all", Expect: e2e.Pass},
		{Binary: "x86_64-clang-no-stack-protector", Expect: e2e.Fail},
		{Binary: "x86_64-clang-stack-protector-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-clang-stack-protector-static", Expect: e2e.Pass},
		{Binary: "x86_64-clang-stack-protector-static-stripped", Expect: e2e.Fail},
		{Binary: "x86_64-clang-stack-protector-lto", Expect: e2e.Pass},

		{Binary: "aarch64-gcc-stack-protector-strong", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-stack-protector-all", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-stack-protector", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-no-stack-protector", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-stack-protector-simple", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-stack-protector-all-simple", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-stack-protector-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-stack-protector-static", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-stack-protector-static-stripped", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-stack-protector-lto", Expect: e2e.Pass},

		{Binary: "aarch64-clang-stack-protector-strong", Expect: e2e.Pass},
		{Binary: "aarch64-clang-stack-protector-all", Expect: e2e.Pass},
		{Binary: "aarch64-clang-no-stack-protector", Expect: e2e.Fail},
		{Binary: "aarch64-clang-stack-protector-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-clang-stack-protector-static", Expect: e2e.Pass},
		{Binary: "aarch64-clang-stack-protector-static-stripped", Expect: e2e.Fail},
		{Binary: "aarch64-clang-stack-protector-lto", Expect: e2e.Pass},

		{Binary: "armv7-gcc-stack-protector-strong", Expect: e2e.Pass},
		{Binary: "armv7-gcc-stack-protector-all", Expect: e2e.Pass},
		{Binary: "armv7-gcc-stack-protector", Expect: e2e.Pass},
		{Binary: "armv7-gcc-no-stack-protector", Expect: e2e.Fail},
		{Binary: "armv7-gcc-stack-protector-simple", Expect: e2e.Fail},
		{Binary: "armv7-gcc-stack-protector-all-simple", Expect: e2e.Pass},
		{Binary: "armv7-gcc-stack-protector-stripped", Expect: e2e.Pass},
		{Binary: "armv7-gcc-stack-protector-static", Expect: e2e.Pass},
		{Binary: "armv7-gcc-stack-protector-static-stripped", Expect: e2e.Fail},
		{Binary: "armv7-gcc-stack-protector-lto", Expect: e2e.Pass},

		{Binary: "armv7-clang-stack-protector-strong", Expect: e2e.Pass},
		{Binary: "armv7-clang-stack-protector-all", Expect: e2e.Pass},
		{Binary: "armv7-clang-no-stack-protector", Expect: e2e.Fail},
		{Binary: "armv7-clang-stack-protector-stripped", Expect: e2e.Pass},
		{Binary: "armv7-clang-stack-protector-static", Expect: e2e.Pass},
		{Binary: "armv7-clang-stack-protector-static-stripped", Expect: e2e.Fail},
		{Binary: "armv7-clang-stack-protector-lto", Expect: e2e.Pass},
	})
}
