package pie_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestPIERule(t *testing.T) {
	e2e.RunRuleTests(t, "pie", []e2e.TestCase{
		{Binary: "x86_64-gcc-pie-explicit", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-no-pie", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-shared", Expect: e2e.Skip},
		{Binary: "x86_64-gcc-pie-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-pie-strip-debug", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-object-file", Expect: e2e.Skip},

		{Binary: "x86_64-clang-pie-explicit", Expect: e2e.Pass},
		{Binary: "x86_64-clang-no-pie", Expect: e2e.Fail},
		{Binary: "x86_64-clang-pie-stripped", Expect: e2e.Pass},

		{Binary: "x86_64-rustc-pie", Expect: e2e.Pass},
		{Binary: "x86_64-rustc-no-pie", Expect: e2e.Fail},
		{Binary: "x86_64-rustc-pie-stripped", Expect: e2e.Pass},

		{Binary: "x86_64-gcc-old-pie", Expect: e2e.Pass},

		{Binary: "aarch64-gcc-pie-explicit", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-no-pie", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-shared", Expect: e2e.Skip},
		{Binary: "aarch64-gcc-pie-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-pie-strip-debug", Expect: e2e.Pass},

		{Binary: "aarch64-clang-pie-explicit", Expect: e2e.Pass},
		{Binary: "aarch64-clang-no-pie", Expect: e2e.Fail},
		{Binary: "aarch64-clang-pie-stripped", Expect: e2e.Pass},

		{Binary: "aarch64-rustc-pie", Expect: e2e.Pass},
		{Binary: "aarch64-rustc-no-pie", Expect: e2e.Fail},
		{Binary: "aarch64-rustc-pie-stripped", Expect: e2e.Pass},

		{Binary: "armv7-gcc-pie-explicit", Expect: e2e.Pass},
		{Binary: "armv7-gcc-no-pie", Expect: e2e.Fail},
		{Binary: "armv7-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "armv7-gcc-shared", Expect: e2e.Skip},
		{Binary: "armv7-gcc-pie-stripped", Expect: e2e.Pass},
		{Binary: "armv7-gcc-pie-strip-debug", Expect: e2e.Pass},

		{Binary: "armv7-clang-pie-explicit", Expect: e2e.Pass},
		{Binary: "armv7-clang-no-pie", Expect: e2e.Fail},
		{Binary: "armv7-clang-pie-stripped", Expect: e2e.Pass},

		{Binary: "armv7-rustc-pie", Expect: e2e.Pass},
		{Binary: "armv7-rustc-no-pie", Expect: e2e.Fail},
		{Binary: "armv7-rustc-pie-stripped", Expect: e2e.Pass},
	})
}
