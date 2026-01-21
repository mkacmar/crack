package stack_canary_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestStackCanaryRule(t *testing.T) {
	e2e.RunRuleTests(t, "stack-canary", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-stack-protector-strong", Expect: "pass"},
		{Binary: "x86_64-gcc-stack-protector-all", Expect: "pass"},
		{Binary: "x86_64-gcc-stack-protector", Expect: "pass"},
		{Binary: "x86_64-gcc-no-stack-protector", Expect: "fail"},
		{Binary: "x86_64-gcc-stack-protector-simple", Expect: "fail"},
		{Binary: "x86_64-gcc-stack-protector-all-simple", Expect: "pass"},
		{Binary: "x86_64-gcc-stack-protector-stripped", Expect: "pass"},
		{Binary: "x86_64-gcc-stack-protector-static", Expect: "pass"},
		{Binary: "x86_64-gcc-stack-protector-static-stripped", Expect: "fail"},
		{Binary: "x86_64-gcc-stack-protector-lto", Expect: "pass"},

		// x86_64 Clang
		{Binary: "x86_64-clang-stack-protector-strong", Expect: "pass"},
		{Binary: "x86_64-clang-stack-protector-all", Expect: "pass"},
		{Binary: "x86_64-clang-no-stack-protector", Expect: "fail"},
		{Binary: "x86_64-clang-stack-protector-stripped", Expect: "pass"},
		{Binary: "x86_64-clang-stack-protector-static", Expect: "pass"},
		{Binary: "x86_64-clang-stack-protector-static-stripped", Expect: "fail"},
		{Binary: "x86_64-clang-stack-protector-lto", Expect: "pass"},

		// aarch64 GCC
		{Binary: "aarch64-gcc-stack-protector-strong", Expect: "pass"},
		{Binary: "aarch64-gcc-stack-protector-all", Expect: "pass"},
		{Binary: "aarch64-gcc-stack-protector", Expect: "pass"},
		{Binary: "aarch64-gcc-no-stack-protector", Expect: "fail"},
		{Binary: "aarch64-gcc-stack-protector-simple", Expect: "fail"},
		{Binary: "aarch64-gcc-stack-protector-all-simple", Expect: "pass"},
		{Binary: "aarch64-gcc-stack-protector-stripped", Expect: "pass"},
		{Binary: "aarch64-gcc-stack-protector-static", Expect: "pass"},
		{Binary: "aarch64-gcc-stack-protector-static-stripped", Expect: "fail"},
		{Binary: "aarch64-gcc-stack-protector-lto", Expect: "pass"},

		// aarch64 Clang
		{Binary: "aarch64-clang-stack-protector-strong", Expect: "pass"},
		{Binary: "aarch64-clang-stack-protector-all", Expect: "pass"},
		{Binary: "aarch64-clang-no-stack-protector", Expect: "fail"},
		{Binary: "aarch64-clang-stack-protector-stripped", Expect: "pass"},
		{Binary: "aarch64-clang-stack-protector-static", Expect: "pass"},
		{Binary: "aarch64-clang-stack-protector-static-stripped", Expect: "fail"},
		{Binary: "aarch64-clang-stack-protector-lto", Expect: "pass"},

		// armv7 GCC
		{Binary: "armv7-gcc-stack-protector-strong", Expect: "pass"},
		{Binary: "armv7-gcc-stack-protector-all", Expect: "pass"},
		{Binary: "armv7-gcc-stack-protector", Expect: "pass"},
		{Binary: "armv7-gcc-no-stack-protector", Expect: "fail"},
		{Binary: "armv7-gcc-stack-protector-simple", Expect: "fail"},
		{Binary: "armv7-gcc-stack-protector-all-simple", Expect: "pass"},
		{Binary: "armv7-gcc-stack-protector-stripped", Expect: "pass"},
		{Binary: "armv7-gcc-stack-protector-static", Expect: "pass"},
		{Binary: "armv7-gcc-stack-protector-static-stripped", Expect: "fail"},
		{Binary: "armv7-gcc-stack-protector-lto", Expect: "pass"},

		// armv7 Clang
		{Binary: "armv7-clang-stack-protector-strong", Expect: "pass"},
		{Binary: "armv7-clang-stack-protector-all", Expect: "pass"},
		{Binary: "armv7-clang-no-stack-protector", Expect: "fail"},
		{Binary: "armv7-clang-stack-protector-stripped", Expect: "pass"},
		{Binary: "armv7-clang-stack-protector-static", Expect: "pass"},
		{Binary: "armv7-clang-stack-protector-static-stripped", Expect: "fail"},
		{Binary: "armv7-clang-stack-protector-lto", Expect: "pass"},
	})
}
