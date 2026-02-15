package aslr_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestASLRRule(t *testing.T) {
	e2e.RunRuleTests(t, "aslr", []e2e.TestCase{
		{Binary: "amd64-gcc-aslr-full", Expect: e2e.Pass},
		{Binary: "amd64-gcc-no-pie", Expect: e2e.Fail},
		{Binary: "amd64-gcc-execstack", Expect: e2e.Fail},
		{Binary: "amd64-gcc-shared", Expect: e2e.Skip},
		{Binary: "amd64-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "amd64-gcc-aslr-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-static-no-pie", Expect: e2e.Fail},
		{Binary: "amd64-gcc-textrel-patched", Expect: e2e.Fail},
		{Binary: "amd64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "amd64-clang-aslr-full", Expect: e2e.Pass},
		{Binary: "amd64-clang-no-pie", Expect: e2e.Fail},
		{Binary: "amd64-clang-execstack", Expect: e2e.Fail},
		{Binary: "amd64-clang-static-no-pie", Expect: e2e.Fail},

		{Binary: "amd64-gcc-old-pie", Expect: e2e.Pass},

		{Binary: "arm64-gcc-aslr-full", Expect: e2e.Pass},
		{Binary: "arm64-gcc-no-pie", Expect: e2e.Fail},
		{Binary: "arm64-gcc-execstack", Expect: e2e.Fail},
		{Binary: "arm64-gcc-shared", Expect: e2e.Skip},
		{Binary: "arm64-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "arm64-gcc-aslr-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-static-no-pie", Expect: e2e.Fail},
		{Binary: "arm64-gcc-textrel-patched", Expect: e2e.Fail},
		{Binary: "arm64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-clang-aslr-full", Expect: e2e.Pass},
		{Binary: "arm64-clang-no-pie", Expect: e2e.Fail},
		{Binary: "arm64-clang-execstack", Expect: e2e.Fail},
		{Binary: "arm64-clang-static-no-pie", Expect: e2e.Fail},

		{Binary: "arm-gcc-aslr-full", Expect: e2e.Pass},
		{Binary: "arm-gcc-no-pie", Expect: e2e.Fail},
		{Binary: "arm-gcc-execstack", Expect: e2e.Fail},
		{Binary: "arm-gcc-shared", Expect: e2e.Skip},
		{Binary: "arm-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "arm-gcc-aslr-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-static-no-pie", Expect: e2e.Fail},
		{Binary: "arm-gcc-textrel-patched", Expect: e2e.Fail},
		{Binary: "arm-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-clang-aslr-full", Expect: e2e.Pass},
		{Binary: "arm-clang-no-pie", Expect: e2e.Fail},
		{Binary: "arm-clang-execstack", Expect: e2e.Fail},
		{Binary: "arm-clang-static-no-pie", Expect: e2e.Fail},

		{Binary: "amd64-rustc-pie", Expect: e2e.Pass},
		{Binary: "amd64-rustc-no-pie", Expect: e2e.Fail},
		{Binary: "amd64-rustc-stripped", Expect: e2e.Pass},
		{Binary: "arm64-rustc-pie", Expect: e2e.Pass},
		{Binary: "arm64-rustc-no-pie", Expect: e2e.Fail},
		{Binary: "arm64-rustc-stripped", Expect: e2e.Pass},
		{Binary: "arm-rustc-pie", Expect: e2e.Pass},
		{Binary: "arm-rustc-no-pie", Expect: e2e.Fail},
		{Binary: "arm-rustc-stripped", Expect: e2e.Pass},
	})
}
