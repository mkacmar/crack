package pie_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestPIERule(t *testing.T) {
	e2e.RunRuleTests(t, "pie", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-pie-explicit", Expect: "pass"},
		{Binary: "x86_64-gcc-no-pie", Expect: "fail"},
		{Binary: "x86_64-gcc-static-pie", Expect: "pass"},
		{Binary: "x86_64-gcc-shared", Expect: "skip"},
		{Binary: "x86_64-gcc-pie-stripped", Expect: "pass"},
		{Binary: "x86_64-gcc-pie-strip-debug", Expect: "pass"},

		// x86_64 Clang
		{Binary: "x86_64-clang-pie-explicit", Expect: "pass"},
		{Binary: "x86_64-clang-no-pie", Expect: "fail"},
		{Binary: "x86_64-clang-pie-stripped", Expect: "pass"},

		// aarch64 GCC
		{Binary: "aarch64-gcc-pie-explicit", Expect: "pass"},
		{Binary: "aarch64-gcc-no-pie", Expect: "fail"},
		{Binary: "aarch64-gcc-static-pie", Expect: "pass"},
		{Binary: "aarch64-gcc-shared", Expect: "skip"},
		{Binary: "aarch64-gcc-pie-stripped", Expect: "pass"},
		{Binary: "aarch64-gcc-pie-strip-debug", Expect: "pass"},

		// aarch64 Clang
		{Binary: "aarch64-clang-pie-explicit", Expect: "pass"},
		{Binary: "aarch64-clang-no-pie", Expect: "fail"},
		{Binary: "aarch64-clang-pie-stripped", Expect: "pass"},

		// armv7 GCC
		{Binary: "armv7-gcc-pie-explicit", Expect: "pass"},
		{Binary: "armv7-gcc-no-pie", Expect: "fail"},
		{Binary: "armv7-gcc-static-pie", Expect: "pass"},
		{Binary: "armv7-gcc-shared", Expect: "skip"},
		{Binary: "armv7-gcc-pie-stripped", Expect: "pass"},
		{Binary: "armv7-gcc-pie-strip-debug", Expect: "pass"},

		// armv7 Clang
		{Binary: "armv7-clang-pie-explicit", Expect: "pass"},
		{Binary: "armv7-clang-no-pie", Expect: "fail"},
		{Binary: "armv7-clang-pie-stripped", Expect: "pass"},
	})
}
