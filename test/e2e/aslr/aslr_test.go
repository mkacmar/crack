package aslr_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestASLRRule(t *testing.T) {
	e2e.RunRuleTests(t, "aslr", []e2e.TestCase{
		// x86_64 GCC
		{"x86_64-gcc-aslr-full", "pass"},
		{"x86_64-gcc-no-pie", "fail"},
		{"x86_64-gcc-execstack", "fail"},
		{"x86_64-gcc-shared", "skip"},
		{"x86_64-gcc-static-pie", "pass"},
		{"x86_64-gcc-aslr-stripped", "pass"},

		// x86_64 Clang
		{"x86_64-clang-aslr-full", "pass"},
		{"x86_64-clang-no-pie", "fail"},
		{"x86_64-clang-execstack", "fail"},

		// aarch64 GCC
		{"aarch64-gcc-aslr-full", "pass"},
		{"aarch64-gcc-no-pie", "fail"},
		{"aarch64-gcc-execstack", "fail"},
		{"aarch64-gcc-shared", "skip"},
		{"aarch64-gcc-static-pie", "pass"},
		{"aarch64-gcc-aslr-stripped", "pass"},

		// aarch64 Clang
		{"aarch64-clang-aslr-full", "pass"},
		{"aarch64-clang-no-pie", "fail"},
		{"aarch64-clang-execstack", "fail"},

		// armv7 GCC
		{"armv7-gcc-aslr-full", "pass"},
		{"armv7-gcc-no-pie", "fail"},
		{"armv7-gcc-execstack", "fail"},
		{"armv7-gcc-shared", "skip"},
		{"armv7-gcc-static-pie", "pass"},
		{"armv7-gcc-aslr-stripped", "pass"},

		// armv7 Clang
		{"armv7-clang-aslr-full", "pass"},
		{"armv7-clang-no-pie", "fail"},
		{"armv7-clang-execstack", "fail"},
	})
}
