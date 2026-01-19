package stack_canary_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestStackCanaryRule(t *testing.T) {
	e2e.RunRuleTests(t, "stack-canary", []e2e.TestCase{
		// x86_64 GCC
		{"x86_64-gcc-stack-protector-strong", "pass"},
		{"x86_64-gcc-stack-protector-all", "pass"},
		{"x86_64-gcc-stack-protector", "pass"},
		{"x86_64-gcc-no-stack-protector", "fail"},
		{"x86_64-gcc-stack-protector-simple", "fail"},
		{"x86_64-gcc-stack-protector-all-simple", "pass"},
		{"x86_64-gcc-stack-protector-stripped", "pass"},
		{"x86_64-gcc-stack-protector-static", "pass"},
		{"x86_64-gcc-stack-protector-static-stripped", "fail"},
		{"x86_64-gcc-stack-protector-lto", "pass"},

		// x86_64 Clang
		{"x86_64-clang-stack-protector-strong", "pass"},
		{"x86_64-clang-stack-protector-all", "pass"},
		{"x86_64-clang-no-stack-protector", "fail"},
		{"x86_64-clang-stack-protector-stripped", "pass"},
		{"x86_64-clang-stack-protector-static", "pass"},
		{"x86_64-clang-stack-protector-static-stripped", "fail"},
		{"x86_64-clang-stack-protector-lto", "pass"},

		// aarch64 GCC
		{"aarch64-gcc-stack-protector-strong", "pass"},
		{"aarch64-gcc-stack-protector-all", "pass"},
		{"aarch64-gcc-stack-protector", "pass"},
		{"aarch64-gcc-no-stack-protector", "fail"},
		{"aarch64-gcc-stack-protector-simple", "fail"},
		{"aarch64-gcc-stack-protector-all-simple", "pass"},
		{"aarch64-gcc-stack-protector-stripped", "pass"},
		{"aarch64-gcc-stack-protector-static", "pass"},
		{"aarch64-gcc-stack-protector-static-stripped", "fail"},
		{"aarch64-gcc-stack-protector-lto", "pass"},

		// aarch64 Clang
		{"aarch64-clang-stack-protector-strong", "pass"},
		{"aarch64-clang-stack-protector-all", "pass"},
		{"aarch64-clang-no-stack-protector", "fail"},
		{"aarch64-clang-stack-protector-stripped", "pass"},
		{"aarch64-clang-stack-protector-static", "pass"},
		{"aarch64-clang-stack-protector-static-stripped", "fail"},
		{"aarch64-clang-stack-protector-lto", "pass"},

		// armv7 GCC
		{"armv7-gcc-stack-protector-strong", "pass"},
		{"armv7-gcc-stack-protector-all", "pass"},
		{"armv7-gcc-stack-protector", "pass"},
		{"armv7-gcc-no-stack-protector", "fail"},
		{"armv7-gcc-stack-protector-simple", "fail"},
		{"armv7-gcc-stack-protector-all-simple", "pass"},
		{"armv7-gcc-stack-protector-stripped", "pass"},
		{"armv7-gcc-stack-protector-static", "pass"},
		{"armv7-gcc-stack-protector-static-stripped", "fail"},
		{"armv7-gcc-stack-protector-lto", "pass"},

		// armv7 Clang
		{"armv7-clang-stack-protector-strong", "pass"},
		{"armv7-clang-stack-protector-all", "pass"},
		{"armv7-clang-no-stack-protector", "fail"},
		{"armv7-clang-stack-protector-stripped", "pass"},
		{"armv7-clang-stack-protector-static", "pass"},
		{"armv7-clang-stack-protector-static-stripped", "fail"},
		{"armv7-clang-stack-protector-lto", "pass"},
	})
}
