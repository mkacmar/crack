package separate_code_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestSeparateCodeRule(t *testing.T) {
	e2e.RunRuleTests(t, "separate-code", []e2e.TestCase{
		{Binary: "x86_64-gcc-separate-code", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-separate-code-static", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-separate-code-shared", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "x86_64-clang-separate-code", Expect: e2e.Pass},
		{Binary: "x86_64-clang-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "aarch64-gcc-separate-code", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-separate-code-static", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-separate-code-shared", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "aarch64-clang-separate-code", Expect: e2e.Pass},
		{Binary: "aarch64-clang-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "armv7-gcc-separate-code", Expect: e2e.Pass},
		{Binary: "armv7-gcc-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "armv7-gcc-separate-code-static", Expect: e2e.Pass},
		{Binary: "armv7-gcc-separate-code-shared", Expect: e2e.Pass},
		{Binary: "armv7-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "armv7-clang-separate-code", Expect: e2e.Pass},
		{Binary: "armv7-clang-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "armv7-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "x86_64-gcc7-no-separate-code", Expect: e2e.Fail},
		{Binary: "aarch64-gcc7-no-separate-code", Expect: e2e.Fail},

		{Binary: "x86_64-rustc-separate-code", Expect: e2e.Pass},
		{Binary: "aarch64-rustc-separate-code", Expect: e2e.Pass},
		{Binary: "armv7-rustc-separate-code", Expect: e2e.Pass},
	})
}
