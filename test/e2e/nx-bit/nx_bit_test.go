package nx_bit_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNXBitRule(t *testing.T) {
	e2e.RunRuleTests(t, "nx-bit", []e2e.TestCase{
		// x86_64 GCC
		{"x86_64-gcc-nx-explicit", "pass"},
		{"x86_64-gcc-no-nx", "fail"},
		{"x86_64-gcc-nx-stripped", "pass"},
		{"x86_64-gcc-nx-static", "pass"},

		// x86_64 Clang
		{"x86_64-clang-nx-explicit", "pass"},
		{"x86_64-clang-no-nx", "fail"},
		{"x86_64-clang-nx-stripped", "pass"},

		// aarch64 GCC
		{"aarch64-gcc-nx-explicit", "pass"},
		{"aarch64-gcc-no-nx", "fail"},
		{"aarch64-gcc-nx-stripped", "pass"},
		{"aarch64-gcc-nx-static", "pass"},

		// aarch64 Clang
		{"aarch64-clang-nx-explicit", "pass"},
		{"aarch64-clang-no-nx", "fail"},
		{"aarch64-clang-nx-stripped", "pass"},

		// armv7 GCC
		{"armv7-gcc-nx-explicit", "pass"},
		{"armv7-gcc-no-nx", "fail"},
		{"armv7-gcc-nx-stripped", "pass"},
		{"armv7-gcc-nx-static", "pass"},

		// armv7 Clang
		{"armv7-clang-nx-explicit", "pass"},
		{"armv7-clang-no-nx", "fail"},
		{"armv7-clang-nx-stripped", "pass"},
	})
}
