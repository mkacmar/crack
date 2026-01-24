package stripped_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestStrippedRule(t *testing.T) {
	e2e.RunRuleTests(t, "stripped", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-not-stripped", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-strip-debug", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-strip-symbols", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-link-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-partial-stripped", Expect: e2e.Fail},

		// x86_64 Clang
		{Binary: "x86_64-clang-not-stripped", Expect: e2e.Fail},
		{Binary: "x86_64-clang-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-clang-strip-debug", Expect: e2e.Fail},

		// aarch64 GCC
		{Binary: "aarch64-gcc-not-stripped", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-strip-debug", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-strip-symbols", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-link-stripped", Expect: e2e.Pass},

		// aarch64 Clang
		{Binary: "aarch64-clang-not-stripped", Expect: e2e.Fail},
		{Binary: "aarch64-clang-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-clang-strip-debug", Expect: e2e.Fail},

		// armv7 GCC
		{Binary: "armv7-gcc-not-stripped", Expect: e2e.Fail},
		{Binary: "armv7-gcc-stripped", Expect: e2e.Pass},
		{Binary: "armv7-gcc-strip-debug", Expect: e2e.Fail},
		{Binary: "armv7-gcc-strip-symbols", Expect: e2e.Pass},
		{Binary: "armv7-gcc-link-stripped", Expect: e2e.Pass},

		// armv7 Clang
		{Binary: "armv7-clang-not-stripped", Expect: e2e.Fail},
		{Binary: "armv7-clang-stripped", Expect: e2e.Pass},
		{Binary: "armv7-clang-strip-debug", Expect: e2e.Fail},
	})
}
