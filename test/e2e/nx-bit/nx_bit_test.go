package nx_bit_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNXBitRule(t *testing.T) {
	e2e.RunRuleTests(t, "nx-bit", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-nx-explicit", Expect: "pass"},
		{Binary: "x86_64-gcc-no-nx", Expect: "fail"},
		{Binary: "x86_64-gcc-nx-stripped", Expect: "pass"},
		{Binary: "x86_64-gcc-nx-static", Expect: "pass"},

		// x86_64 Clang
		{Binary: "x86_64-clang-nx-explicit", Expect: "pass"},
		{Binary: "x86_64-clang-no-nx", Expect: "fail"},
		{Binary: "x86_64-clang-nx-stripped", Expect: "pass"},

		// aarch64 GCC
		{Binary: "aarch64-gcc-nx-explicit", Expect: "pass"},
		{Binary: "aarch64-gcc-no-nx", Expect: "fail"},
		{Binary: "aarch64-gcc-nx-stripped", Expect: "pass"},
		{Binary: "aarch64-gcc-nx-static", Expect: "pass"},

		// aarch64 Clang
		{Binary: "aarch64-clang-nx-explicit", Expect: "pass"},
		{Binary: "aarch64-clang-no-nx", Expect: "fail"},
		{Binary: "aarch64-clang-nx-stripped", Expect: "pass"},

		// armv7 GCC
		{Binary: "armv7-gcc-nx-explicit", Expect: "pass"},
		{Binary: "armv7-gcc-no-nx", Expect: "fail"},
		{Binary: "armv7-gcc-nx-stripped", Expect: "pass"},
		{Binary: "armv7-gcc-nx-static", Expect: "pass"},

		// armv7 Clang
		{Binary: "armv7-clang-nx-explicit", Expect: "pass"},
		{Binary: "armv7-clang-no-nx", Expect: "fail"},
		{Binary: "armv7-clang-nx-stripped", Expect: "pass"},
	})
}
