package relro_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestRELRORule(t *testing.T) {
	e2e.RunRuleTests(t, "relro", []e2e.TestCase{
		// x86_64 GCC
		{"x86_64-gcc-partial-relro", "pass"},
		{"x86_64-gcc-full-relro", "pass"},
		{"x86_64-gcc-no-relro", "fail"},
		{"x86_64-gcc-full-relro-stripped", "pass"},
		{"x86_64-gcc-full-relro-static", "pass"},
		{"x86_64-gcc-full-relro-shared", "pass"},

		// x86_64 Clang
		{"x86_64-clang-partial-relro", "pass"},
		{"x86_64-clang-full-relro", "pass"},
		{"x86_64-clang-no-relro", "fail"},
		{"x86_64-clang-full-relro-stripped", "pass"},

		// aarch64 GCC
		{"aarch64-gcc-partial-relro", "pass"},
		{"aarch64-gcc-full-relro", "pass"},
		{"aarch64-gcc-no-relro", "fail"},
		{"aarch64-gcc-full-relro-stripped", "pass"},
		{"aarch64-gcc-full-relro-static", "pass"},
		{"aarch64-gcc-full-relro-shared", "pass"},

		// aarch64 Clang
		{"aarch64-clang-partial-relro", "pass"},
		{"aarch64-clang-full-relro", "pass"},
		{"aarch64-clang-no-relro", "fail"},
		{"aarch64-clang-full-relro-stripped", "pass"},

		// armv7 GCC
		{"armv7-gcc-partial-relro", "pass"},
		{"armv7-gcc-full-relro", "pass"},
		{"armv7-gcc-no-relro", "fail"},
		{"armv7-gcc-full-relro-stripped", "pass"},
		{"armv7-gcc-full-relro-static", "pass"},
		{"armv7-gcc-full-relro-shared", "pass"},

		// armv7 Clang
		{"armv7-clang-partial-relro", "pass"},
		{"armv7-clang-full-relro", "pass"},
		{"armv7-clang-no-relro", "fail"},
		{"armv7-clang-full-relro-stripped", "pass"},
	})
}

func TestFullRELRORule(t *testing.T) {
	e2e.RunRuleTests(t, "full-relro", []e2e.TestCase{
		// x86_64 GCC
		{"x86_64-gcc-partial-relro", "fail"},
		{"x86_64-gcc-full-relro", "pass"},
		{"x86_64-gcc-no-relro", "fail"},
		{"x86_64-gcc-full-relro-stripped", "pass"},
		{"x86_64-gcc-full-relro-static", "pass"},
		{"x86_64-gcc-full-relro-shared", "pass"},

		// x86_64 Clang
		{"x86_64-clang-partial-relro", "fail"},
		{"x86_64-clang-full-relro", "pass"},
		{"x86_64-clang-no-relro", "fail"},
		{"x86_64-clang-full-relro-stripped", "pass"},

		// aarch64 GCC
		{"aarch64-gcc-partial-relro", "fail"},
		{"aarch64-gcc-full-relro", "pass"},
		{"aarch64-gcc-no-relro", "fail"},
		{"aarch64-gcc-full-relro-stripped", "pass"},
		{"aarch64-gcc-full-relro-static", "pass"},
		{"aarch64-gcc-full-relro-shared", "pass"},

		// aarch64 Clang
		{"aarch64-clang-partial-relro", "fail"},
		{"aarch64-clang-full-relro", "pass"},
		{"aarch64-clang-no-relro", "fail"},
		{"aarch64-clang-full-relro-stripped", "pass"},

		// armv7 GCC
		{"armv7-gcc-partial-relro", "fail"},
		{"armv7-gcc-full-relro", "pass"},
		{"armv7-gcc-no-relro", "fail"},
		{"armv7-gcc-full-relro-stripped", "pass"},
		{"armv7-gcc-full-relro-static", "pass"},
		{"armv7-gcc-full-relro-shared", "pass"},

		// armv7 Clang
		{"armv7-clang-partial-relro", "fail"},
		{"armv7-clang-full-relro", "pass"},
		{"armv7-clang-no-relro", "fail"},
		{"armv7-clang-full-relro-stripped", "pass"},
	})
}
