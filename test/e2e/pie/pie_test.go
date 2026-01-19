package pie_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestPIERule(t *testing.T) {
	e2e.RunRuleTests(t, "pie", []e2e.TestCase{
		// x86_64 GCC
		{"x86_64-gcc-pie-explicit", "pass"},
		{"x86_64-gcc-no-pie", "fail"},
		{"x86_64-gcc-static-pie", "pass"},
		{"x86_64-gcc-shared", "skip"},
		{"x86_64-gcc-pie-stripped", "pass"},
		{"x86_64-gcc-pie-strip-debug", "pass"},

		// x86_64 Clang
		{"x86_64-clang-pie-explicit", "pass"},
		{"x86_64-clang-no-pie", "fail"},
		{"x86_64-clang-pie-stripped", "pass"},

		// aarch64 GCC
		{"aarch64-gcc-pie-explicit", "pass"},
		{"aarch64-gcc-no-pie", "fail"},
		{"aarch64-gcc-static-pie", "pass"},
		{"aarch64-gcc-shared", "skip"},
		{"aarch64-gcc-pie-stripped", "pass"},
		{"aarch64-gcc-pie-strip-debug", "pass"},

		// aarch64 Clang
		{"aarch64-clang-pie-explicit", "pass"},
		{"aarch64-clang-no-pie", "fail"},
		{"aarch64-clang-pie-stripped", "pass"},

		// armv7 GCC
		{"armv7-gcc-pie-explicit", "pass"},
		{"armv7-gcc-no-pie", "fail"},
		{"armv7-gcc-static-pie", "pass"},
		{"armv7-gcc-shared", "skip"},
		{"armv7-gcc-pie-stripped", "pass"},
		{"armv7-gcc-pie-strip-debug", "pass"},

		// armv7 Clang
		{"armv7-clang-pie-explicit", "pass"},
		{"armv7-clang-no-pie", "fail"},
		{"armv7-clang-pie-stripped", "pass"},
	})
}
