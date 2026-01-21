package relro_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestRELRORule(t *testing.T) {
	e2e.RunRuleTests(t, "relro", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-partial-relro", Expect: "pass"},
		{Binary: "x86_64-gcc-full-relro", Expect: "pass"},
		{Binary: "x86_64-gcc-no-relro", Expect: "fail"},
		{Binary: "x86_64-gcc-full-relro-stripped", Expect: "pass"},
		{Binary: "x86_64-gcc-full-relro-static", Expect: "pass"},
		{Binary: "x86_64-gcc-full-relro-shared", Expect: "pass"},

		// x86_64 Clang
		{Binary: "x86_64-clang-partial-relro", Expect: "pass"},
		{Binary: "x86_64-clang-full-relro", Expect: "pass"},
		{Binary: "x86_64-clang-no-relro", Expect: "fail"},
		{Binary: "x86_64-clang-full-relro-stripped", Expect: "pass"},

		// aarch64 GCC
		{Binary: "aarch64-gcc-partial-relro", Expect: "pass"},
		{Binary: "aarch64-gcc-full-relro", Expect: "pass"},
		{Binary: "aarch64-gcc-no-relro", Expect: "fail"},
		{Binary: "aarch64-gcc-full-relro-stripped", Expect: "pass"},
		{Binary: "aarch64-gcc-full-relro-static", Expect: "pass"},
		{Binary: "aarch64-gcc-full-relro-shared", Expect: "pass"},

		// aarch64 Clang
		{Binary: "aarch64-clang-partial-relro", Expect: "pass"},
		{Binary: "aarch64-clang-full-relro", Expect: "pass"},
		{Binary: "aarch64-clang-no-relro", Expect: "fail"},
		{Binary: "aarch64-clang-full-relro-stripped", Expect: "pass"},

		// armv7 GCC
		{Binary: "armv7-gcc-partial-relro", Expect: "pass"},
		{Binary: "armv7-gcc-full-relro", Expect: "pass"},
		{Binary: "armv7-gcc-no-relro", Expect: "fail"},
		{Binary: "armv7-gcc-full-relro-stripped", Expect: "pass"},
		{Binary: "armv7-gcc-full-relro-static", Expect: "pass"},
		{Binary: "armv7-gcc-full-relro-shared", Expect: "pass"},

		// armv7 Clang
		{Binary: "armv7-clang-partial-relro", Expect: "pass"},
		{Binary: "armv7-clang-full-relro", Expect: "pass"},
		{Binary: "armv7-clang-no-relro", Expect: "fail"},
		{Binary: "armv7-clang-full-relro-stripped", Expect: "pass"},
	})
}

func TestFullRELRORule(t *testing.T) {
	e2e.RunRuleTests(t, "full-relro", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-partial-relro", Expect: "fail"},
		{Binary: "x86_64-gcc-full-relro", Expect: "pass"},
		{Binary: "x86_64-gcc-no-relro", Expect: "fail"},
		{Binary: "x86_64-gcc-full-relro-stripped", Expect: "pass"},
		{Binary: "x86_64-gcc-full-relro-static", Expect: "pass"},
		{Binary: "x86_64-gcc-full-relro-shared", Expect: "pass"},

		// x86_64 Clang
		{Binary: "x86_64-clang-partial-relro", Expect: "fail"},
		{Binary: "x86_64-clang-full-relro", Expect: "pass"},
		{Binary: "x86_64-clang-no-relro", Expect: "fail"},
		{Binary: "x86_64-clang-full-relro-stripped", Expect: "pass"},

		// aarch64 GCC
		{Binary: "aarch64-gcc-partial-relro", Expect: "fail"},
		{Binary: "aarch64-gcc-full-relro", Expect: "pass"},
		{Binary: "aarch64-gcc-no-relro", Expect: "fail"},
		{Binary: "aarch64-gcc-full-relro-stripped", Expect: "pass"},
		{Binary: "aarch64-gcc-full-relro-static", Expect: "pass"},
		{Binary: "aarch64-gcc-full-relro-shared", Expect: "pass"},

		// aarch64 Clang
		{Binary: "aarch64-clang-partial-relro", Expect: "fail"},
		{Binary: "aarch64-clang-full-relro", Expect: "pass"},
		{Binary: "aarch64-clang-no-relro", Expect: "fail"},
		{Binary: "aarch64-clang-full-relro-stripped", Expect: "pass"},

		// armv7 GCC
		{Binary: "armv7-gcc-partial-relro", Expect: "fail"},
		{Binary: "armv7-gcc-full-relro", Expect: "pass"},
		{Binary: "armv7-gcc-no-relro", Expect: "fail"},
		{Binary: "armv7-gcc-full-relro-stripped", Expect: "pass"},
		{Binary: "armv7-gcc-full-relro-static", Expect: "pass"},
		{Binary: "armv7-gcc-full-relro-shared", Expect: "pass"},

		// armv7 Clang
		{Binary: "armv7-clang-partial-relro", Expect: "fail"},
		{Binary: "armv7-clang-full-relro", Expect: "pass"},
		{Binary: "armv7-clang-no-relro", Expect: "fail"},
		{Binary: "armv7-clang-full-relro-stripped", Expect: "pass"},
	})
}
