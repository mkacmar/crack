package ubsan_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestUBSanRule(t *testing.T) {
	e2e.RunRuleTests(t, "ubsan", []e2e.TestCase{
		{Binary: "x86_64-gcc-ubsan", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-ubsan-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-no-ubsan", Expect: e2e.Fail},

		{Binary: "x86_64-clang-ubsan", Expect: e2e.Pass},
		{Binary: "x86_64-clang-ubsan-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-clang-no-ubsan", Expect: e2e.Fail},

		{Binary: "x86_64-rustc-no-ubsan", Expect: e2e.Skip},

		{Binary: "aarch64-gcc-ubsan", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-ubsan-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-no-ubsan", Expect: e2e.Fail},

		{Binary: "aarch64-clang-ubsan", Expect: e2e.Pass},
		{Binary: "aarch64-clang-ubsan-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-clang-no-ubsan", Expect: e2e.Fail},

		{Binary: "aarch64-rustc-no-ubsan", Expect: e2e.Skip},

		{Binary: "armv7-gcc-ubsan", Expect: e2e.Pass},
		{Binary: "armv7-gcc-ubsan-stripped", Expect: e2e.Pass},
		{Binary: "armv7-gcc-no-ubsan", Expect: e2e.Fail},

		{Binary: "armv7-clang-ubsan", Expect: e2e.Pass},
		{Binary: "armv7-clang-ubsan-stripped", Expect: e2e.Pass},
		{Binary: "armv7-clang-no-ubsan", Expect: e2e.Fail},

		{Binary: "armv7-rustc-no-ubsan", Expect: e2e.Skip},
	})
}
