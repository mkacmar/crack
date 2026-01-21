package stripped_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestStrippedRule(t *testing.T) {
	e2e.RunRuleTests(t, "stripped", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-not-stripped", Expect: "fail"},
		{Binary: "x86_64-gcc-stripped", Expect: "pass"},
		{Binary: "x86_64-gcc-strip-debug", Expect: "fail"},
		{Binary: "x86_64-gcc-strip-symbols", Expect: "pass"},
		{Binary: "x86_64-gcc-link-stripped", Expect: "pass"},

		// x86_64 Clang
		{Binary: "x86_64-clang-not-stripped", Expect: "fail"},
		{Binary: "x86_64-clang-stripped", Expect: "pass"},
		{Binary: "x86_64-clang-strip-debug", Expect: "fail"},

		// aarch64 GCC
		{Binary: "aarch64-gcc-not-stripped", Expect: "fail"},
		{Binary: "aarch64-gcc-stripped", Expect: "pass"},
		{Binary: "aarch64-gcc-strip-debug", Expect: "fail"},
		{Binary: "aarch64-gcc-strip-symbols", Expect: "pass"},
		{Binary: "aarch64-gcc-link-stripped", Expect: "pass"},

		// aarch64 Clang
		{Binary: "aarch64-clang-not-stripped", Expect: "fail"},
		{Binary: "aarch64-clang-stripped", Expect: "pass"},
		{Binary: "aarch64-clang-strip-debug", Expect: "fail"},

		// armv7 GCC
		{Binary: "armv7-gcc-not-stripped", Expect: "fail"},
		{Binary: "armv7-gcc-stripped", Expect: "pass"},
		{Binary: "armv7-gcc-strip-debug", Expect: "fail"},
		{Binary: "armv7-gcc-strip-symbols", Expect: "pass"},
		{Binary: "armv7-gcc-link-stripped", Expect: "pass"},

		// armv7 Clang
		{Binary: "armv7-clang-not-stripped", Expect: "fail"},
		{Binary: "armv7-clang-stripped", Expect: "pass"},
		{Binary: "armv7-clang-strip-debug", Expect: "fail"},
	})
}
