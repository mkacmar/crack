package pie_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestPIERule(t *testing.T) {
	e2e.RunRuleTests(t, "pie", []e2e.TestCase{
		// x86_64 GCC
		{"x86_64-gcc-pie-explicit", "pass"},
		{"x86_64-gcc-pie-default", "pass"},
		{"x86_64-gcc-no-pie", "fail"},
		{"x86_64-gcc-static-pie", "pass"},
		{"x86_64-gcc-shared", "skip"},

		// x86_64 Clang
		{"x86_64-clang-pie-explicit", "pass"},
		{"x86_64-clang-pie-default", "pass"},
		{"x86_64-clang-no-pie", "fail"},

		// aarch64 GCC
		{"aarch64-gcc-pie-explicit", "pass"},
		{"aarch64-gcc-pie-default", "pass"},
		{"aarch64-gcc-no-pie", "fail"},
		{"aarch64-gcc-static-pie", "pass"},
		{"aarch64-gcc-shared", "skip"},

		// aarch64 Clang
		{"aarch64-clang-pie-explicit", "pass"},
		{"aarch64-clang-pie-default", "pass"},
		{"aarch64-clang-no-pie", "fail"},

		// armv7 (32-bit) GCC
		{"armv7-gcc-pie-explicit", "pass"},
		{"armv7-gcc-pie-default", "pass"},
		{"armv7-gcc-no-pie", "fail"},
		{"armv7-gcc-static-pie", "pass"},
		{"armv7-gcc-shared", "skip"},

		// armv7 (32-bit) Clang
		{"armv7-clang-pie-explicit", "pass"},
		{"armv7-clang-pie-default", "pass"},
		{"armv7-clang-no-pie", "fail"},
	})
}
