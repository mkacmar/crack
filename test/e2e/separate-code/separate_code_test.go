package separate_code_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestSeparateCodeRule(t *testing.T) {
	e2e.RunRuleTests(t, "separate-code", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-separate-code", Expect: "pass"},
		{Binary: "x86_64-gcc-separate-code-stripped", Expect: "pass"},
		{Binary: "x86_64-gcc-separate-code-static", Expect: "pass"},
		{Binary: "x86_64-gcc-separate-code-shared", Expect: "pass"},

		// x86_64 Clang
		{Binary: "x86_64-clang-separate-code", Expect: "pass"},
		{Binary: "x86_64-clang-separate-code-stripped", Expect: "pass"},

		// aarch64 GCC
		{Binary: "aarch64-gcc-separate-code", Expect: "pass"},
		{Binary: "aarch64-gcc-separate-code-stripped", Expect: "pass"},
		{Binary: "aarch64-gcc-separate-code-static", Expect: "pass"},
		{Binary: "aarch64-gcc-separate-code-shared", Expect: "pass"},

		// aarch64 Clang
		{Binary: "aarch64-clang-separate-code", Expect: "pass"},
		{Binary: "aarch64-clang-separate-code-stripped", Expect: "pass"},

		// armv7 GCC
		{Binary: "armv7-gcc-separate-code", Expect: "pass"},
		{Binary: "armv7-gcc-separate-code-stripped", Expect: "pass"},
		{Binary: "armv7-gcc-separate-code-static", Expect: "pass"},
		{Binary: "armv7-gcc-separate-code-shared", Expect: "pass"},

		// armv7 Clang
		{Binary: "armv7-clang-separate-code", Expect: "pass"},
		{Binary: "armv7-clang-separate-code-stripped", Expect: "pass"},

		// older toolchain - no separate-code by default
		{Binary: "x86_64-gcc7-no-separate-code", Expect: "fail"},
		{Binary: "aarch64-gcc7-no-separate-code", Expect: "fail"},
	})
}
