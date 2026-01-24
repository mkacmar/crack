package aslr_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestASLRRule(t *testing.T) {
	e2e.RunRuleTests(t, "aslr", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-aslr-full", Expect: "pass"},
		{Binary: "x86_64-gcc-no-pie", Expect: "fail"},
		{Binary: "x86_64-gcc-execstack", Expect: "fail"},
		{Binary: "x86_64-gcc-shared", Expect: "skip"},
		{Binary: "x86_64-gcc-static-pie", Expect: "pass"},
		{Binary: "x86_64-gcc-aslr-stripped", Expect: "pass"},
		{Binary: "x86_64-gcc-static-no-pie", Expect: "fail"},
		{Binary: "x86_64-gcc-textrel-patched", Expect: "fail"},

		// x86_64 Clang
		{Binary: "x86_64-clang-aslr-full", Expect: "pass"},
		{Binary: "x86_64-clang-no-pie", Expect: "fail"},
		{Binary: "x86_64-clang-execstack", Expect: "fail"},
		{Binary: "x86_64-clang-static-no-pie", Expect: "fail"},

		// x86_64 GCC old toolchain
		{Binary: "x86_64-gcc-old-pie", Expect: "pass"},

		// aarch64 GCC
		{Binary: "aarch64-gcc-aslr-full", Expect: "pass"},
		{Binary: "aarch64-gcc-no-pie", Expect: "fail"},
		{Binary: "aarch64-gcc-execstack", Expect: "fail"},
		{Binary: "aarch64-gcc-shared", Expect: "skip"},
		{Binary: "aarch64-gcc-static-pie", Expect: "pass"},
		{Binary: "aarch64-gcc-aslr-stripped", Expect: "pass"},
		{Binary: "aarch64-gcc-static-no-pie", Expect: "fail"},
		{Binary: "aarch64-gcc-textrel-patched", Expect: "fail"},

		// aarch64 Clang
		{Binary: "aarch64-clang-aslr-full", Expect: "pass"},
		{Binary: "aarch64-clang-no-pie", Expect: "fail"},
		{Binary: "aarch64-clang-execstack", Expect: "fail"},
		{Binary: "aarch64-clang-static-no-pie", Expect: "fail"},

		// armv7 GCC
		{Binary: "armv7-gcc-aslr-full", Expect: "pass"},
		{Binary: "armv7-gcc-no-pie", Expect: "fail"},
		{Binary: "armv7-gcc-execstack", Expect: "fail"},
		{Binary: "armv7-gcc-shared", Expect: "skip"},
		{Binary: "armv7-gcc-static-pie", Expect: "pass"},
		{Binary: "armv7-gcc-aslr-stripped", Expect: "pass"},
		{Binary: "armv7-gcc-static-no-pie", Expect: "fail"},
		{Binary: "armv7-gcc-textrel-patched", Expect: "fail"},

		// armv7 Clang
		{Binary: "armv7-clang-aslr-full", Expect: "pass"},
		{Binary: "armv7-clang-no-pie", Expect: "fail"},
		{Binary: "armv7-clang-execstack", Expect: "fail"},
		{Binary: "armv7-clang-static-no-pie", Expect: "fail"},
	})
}
