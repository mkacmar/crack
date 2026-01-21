package separate_code_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestSeparateCodeRule(t *testing.T) {
	e2e.RunRuleTests(t, "separate-code", []e2e.TestCase{
		// x86_64 GCC
		{"x86_64-gcc-separate-code", "pass"},
		{"x86_64-gcc-separate-code-stripped", "pass"},
		{"x86_64-gcc-separate-code-static", "pass"},
		{"x86_64-gcc-separate-code-shared", "pass"},

		// x86_64 Clang
		{"x86_64-clang-separate-code", "pass"},
		{"x86_64-clang-separate-code-stripped", "pass"},

		// aarch64 GCC
		{"aarch64-gcc-separate-code", "pass"},
		{"aarch64-gcc-separate-code-stripped", "pass"},
		{"aarch64-gcc-separate-code-static", "pass"},
		{"aarch64-gcc-separate-code-shared", "pass"},

		// aarch64 Clang
		{"aarch64-clang-separate-code", "pass"},
		{"aarch64-clang-separate-code-stripped", "pass"},

		// armv7 GCC
		{"armv7-gcc-separate-code", "pass"},
		{"armv7-gcc-separate-code-stripped", "pass"},
		{"armv7-gcc-separate-code-static", "pass"},
		{"armv7-gcc-separate-code-shared", "pass"},

		// armv7 Clang
		{"armv7-clang-separate-code", "pass"},
		{"armv7-clang-separate-code-stripped", "pass"},

		// older toolchain - no separate-code by default
		{"x86_64-gcc7-no-separate-code", "fail"},
		{"aarch64-gcc7-no-separate-code", "fail"},
	})
}
