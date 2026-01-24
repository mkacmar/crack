package separate_code_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestSeparateCodeRule(t *testing.T) {
	e2e.RunRuleTests(t, "separate-code", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-separate-code", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-separate-code-static", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-separate-code-shared", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-object-file", Expect: e2e.Skip},

		// x86_64 Clang
		{Binary: "x86_64-clang-separate-code", Expect: e2e.Pass},
		{Binary: "x86_64-clang-separate-code-stripped", Expect: e2e.Pass},

		// aarch64 GCC
		{Binary: "aarch64-gcc-separate-code", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-separate-code-static", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-separate-code-shared", Expect: e2e.Pass},

		// aarch64 Clang
		{Binary: "aarch64-clang-separate-code", Expect: e2e.Pass},
		{Binary: "aarch64-clang-separate-code-stripped", Expect: e2e.Pass},

		// armv7 GCC
		{Binary: "armv7-gcc-separate-code", Expect: e2e.Pass},
		{Binary: "armv7-gcc-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "armv7-gcc-separate-code-static", Expect: e2e.Pass},
		{Binary: "armv7-gcc-separate-code-shared", Expect: e2e.Pass},

		// armv7 Clang
		{Binary: "armv7-clang-separate-code", Expect: e2e.Pass},
		{Binary: "armv7-clang-separate-code-stripped", Expect: e2e.Pass},

		// older toolchain - no separate-code by default
		{Binary: "x86_64-gcc7-no-separate-code", Expect: e2e.Fail},
		{Binary: "aarch64-gcc7-no-separate-code", Expect: e2e.Fail},
	})
}
