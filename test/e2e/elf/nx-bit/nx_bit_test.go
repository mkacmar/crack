package nx_bit_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNXBitRule(t *testing.T) {
	e2e.RunRuleTests(t, "nx-bit", []e2e.TestCase{
		{Binary: "amd64-gcc-nx-explicit", Expect: e2e.Pass},
		{Binary: "amd64-gcc-no-nx", Expect: e2e.Fail},
		{Binary: "amd64-gcc-nx-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-nx-static", Expect: e2e.Pass},
		{Binary: "amd64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "amd64-clang-nx-explicit", Expect: e2e.Pass},
		{Binary: "amd64-clang-no-nx", Expect: e2e.Fail},
		{Binary: "amd64-clang-nx-stripped", Expect: e2e.Pass},
		{Binary: "amd64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-gcc-nx-explicit", Expect: e2e.Pass},
		{Binary: "arm64-gcc-no-nx", Expect: e2e.Fail},
		{Binary: "arm64-gcc-nx-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-nx-static", Expect: e2e.Pass},
		{Binary: "arm64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-clang-nx-explicit", Expect: e2e.Pass},
		{Binary: "arm64-clang-no-nx", Expect: e2e.Fail},
		{Binary: "arm64-clang-nx-stripped", Expect: e2e.Pass},
		{Binary: "arm64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-gcc-nx-explicit", Expect: e2e.Pass},
		{Binary: "arm-gcc-no-nx", Expect: e2e.Fail},
		{Binary: "arm-gcc-nx-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-nx-static", Expect: e2e.Pass},
		{Binary: "arm-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-clang-nx-explicit", Expect: e2e.Pass},
		{Binary: "arm-clang-no-nx", Expect: e2e.Fail},
		{Binary: "arm-clang-nx-stripped", Expect: e2e.Pass},
		{Binary: "arm-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "amd64-rustc-nx", Expect: e2e.Pass},
		{Binary: "amd64-rustc-nx-stripped", Expect: e2e.Pass},
		{Binary: "arm64-rustc-nx", Expect: e2e.Pass},
		{Binary: "arm64-rustc-nx-stripped", Expect: e2e.Pass},
		{Binary: "arm-rustc-nx", Expect: e2e.Pass},
		{Binary: "arm-rustc-nx-stripped", Expect: e2e.Pass},
	})
}
