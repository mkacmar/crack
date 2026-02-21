package pie_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestPIERule(t *testing.T) {
	e2e.RunRuleTests(t, "pie", []e2e.TestCase{
		{Binary: "amd64-gcc-pie-explicit", Expect: e2e.Pass},
		{Binary: "amd64-gcc-no-pie", Expect: e2e.Fail},
		{Binary: "amd64-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "amd64-gcc-shared", Expect: e2e.Skip},
		{Binary: "amd64-gcc-pie-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-pie-strip-debug", Expect: e2e.Pass},
		{Binary: "amd64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "amd64-clang-pie-explicit", Expect: e2e.Pass},
		{Binary: "amd64-clang-no-pie", Expect: e2e.Fail},
		{Binary: "amd64-clang-pie-stripped", Expect: e2e.Pass},

		{Binary: "amd64-rustc-pie", Expect: e2e.Pass},
		{Binary: "amd64-rustc-no-pie", Expect: e2e.Fail},
		{Binary: "amd64-rustc-pie-stripped", Expect: e2e.Pass},

		{Binary: "amd64-gcc-old-pie", Expect: e2e.Pass},

		{Binary: "arm64-gcc-pie-explicit", Expect: e2e.Pass},
		{Binary: "arm64-gcc-no-pie", Expect: e2e.Fail},
		{Binary: "arm64-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "arm64-gcc-shared", Expect: e2e.Skip},
		{Binary: "arm64-gcc-pie-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-pie-strip-debug", Expect: e2e.Pass},
		{Binary: "arm64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-clang-pie-explicit", Expect: e2e.Pass},
		{Binary: "arm64-clang-no-pie", Expect: e2e.Fail},
		{Binary: "arm64-clang-pie-stripped", Expect: e2e.Pass},

		{Binary: "arm64-rustc-pie", Expect: e2e.Pass},
		{Binary: "arm64-rustc-no-pie", Expect: e2e.Fail},
		{Binary: "arm64-rustc-pie-stripped", Expect: e2e.Pass},

		{Binary: "arm-gcc-pie-explicit", Expect: e2e.Pass},
		{Binary: "arm-gcc-no-pie", Expect: e2e.Fail},
		{Binary: "arm-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "arm-gcc-shared", Expect: e2e.Skip},
		{Binary: "arm-gcc-pie-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-pie-strip-debug", Expect: e2e.Pass},
		{Binary: "arm-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-clang-pie-explicit", Expect: e2e.Pass},
		{Binary: "arm-clang-no-pie", Expect: e2e.Fail},
		{Binary: "arm-clang-pie-stripped", Expect: e2e.Pass},

		{Binary: "arm-rustc-pie", Expect: e2e.Pass},
		{Binary: "arm-rustc-no-pie", Expect: e2e.Fail},
		{Binary: "arm-rustc-pie-stripped", Expect: e2e.Pass},
	})
}
