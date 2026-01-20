package stripped_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestStrippedRule(t *testing.T) {
	e2e.RunRuleTests(t, "stripped", []e2e.TestCase{
		// x86_64 GCC
		{"x86_64-gcc-not-stripped", "fail"},
		{"x86_64-gcc-stripped", "pass"},
		{"x86_64-gcc-strip-debug", "fail"},
		{"x86_64-gcc-strip-symbols", "pass"},
		{"x86_64-gcc-link-stripped", "pass"},

		// x86_64 Clang
		{"x86_64-clang-not-stripped", "fail"},
		{"x86_64-clang-stripped", "pass"},
		{"x86_64-clang-strip-debug", "fail"},

		// aarch64 GCC
		{"aarch64-gcc-not-stripped", "fail"},
		{"aarch64-gcc-stripped", "pass"},
		{"aarch64-gcc-strip-debug", "fail"},
		{"aarch64-gcc-strip-symbols", "pass"},
		{"aarch64-gcc-link-stripped", "pass"},

		// aarch64 Clang
		{"aarch64-clang-not-stripped", "fail"},
		{"aarch64-clang-stripped", "pass"},
		{"aarch64-clang-strip-debug", "fail"},

		// armv7 GCC
		{"armv7-gcc-not-stripped", "fail"},
		{"armv7-gcc-stripped", "pass"},
		{"armv7-gcc-strip-debug", "fail"},
		{"armv7-gcc-strip-symbols", "pass"},
		{"armv7-gcc-link-stripped", "pass"},

		// armv7 Clang
		{"armv7-clang-not-stripped", "fail"},
		{"armv7-clang-stripped", "pass"},
		{"armv7-clang-strip-debug", "fail"},
	})
}
