package aslr_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestASLRRule(t *testing.T) {
	e2e.RunRuleTests(t, "aslr", []e2e.TestCase{
		{Binary: "x86_64-gcc-aslr-full", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-no-pie", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-execstack", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-shared", Expect: e2e.Skip},
		{Binary: "x86_64-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-aslr-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-static-no-pie", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-textrel-patched", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "x86_64-clang-aslr-full", Expect: e2e.Pass},
		{Binary: "x86_64-clang-no-pie", Expect: e2e.Fail},
		{Binary: "x86_64-clang-execstack", Expect: e2e.Fail},
		{Binary: "x86_64-clang-static-no-pie", Expect: e2e.Fail},

		{Binary: "x86_64-gcc-old-pie", Expect: e2e.Pass},

		{Binary: "aarch64-gcc-aslr-full", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-no-pie", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-execstack", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-shared", Expect: e2e.Skip},
		{Binary: "aarch64-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-aslr-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-static-no-pie", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-textrel-patched", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "aarch64-clang-aslr-full", Expect: e2e.Pass},
		{Binary: "aarch64-clang-no-pie", Expect: e2e.Fail},
		{Binary: "aarch64-clang-execstack", Expect: e2e.Fail},
		{Binary: "aarch64-clang-static-no-pie", Expect: e2e.Fail},

		{Binary: "armv7-gcc-aslr-full", Expect: e2e.Pass},
		{Binary: "armv7-gcc-no-pie", Expect: e2e.Fail},
		{Binary: "armv7-gcc-execstack", Expect: e2e.Fail},
		{Binary: "armv7-gcc-shared", Expect: e2e.Skip},
		{Binary: "armv7-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "armv7-gcc-aslr-stripped", Expect: e2e.Pass},
		{Binary: "armv7-gcc-static-no-pie", Expect: e2e.Fail},
		{Binary: "armv7-gcc-textrel-patched", Expect: e2e.Fail},
		{Binary: "armv7-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "armv7-clang-aslr-full", Expect: e2e.Pass},
		{Binary: "armv7-clang-no-pie", Expect: e2e.Fail},
		{Binary: "armv7-clang-execstack", Expect: e2e.Fail},
		{Binary: "armv7-clang-static-no-pie", Expect: e2e.Fail},

		{Binary: "x86_64-rustc-pie", Expect: e2e.Pass},
		{Binary: "x86_64-rustc-no-pie", Expect: e2e.Fail},
		{Binary: "x86_64-rustc-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-rustc-pie", Expect: e2e.Pass},
		{Binary: "aarch64-rustc-no-pie", Expect: e2e.Fail},
		{Binary: "aarch64-rustc-stripped", Expect: e2e.Pass},
		{Binary: "armv7-rustc-pie", Expect: e2e.Pass},
		{Binary: "armv7-rustc-no-pie", Expect: e2e.Fail},
		{Binary: "armv7-rustc-stripped", Expect: e2e.Pass},
	})
}
