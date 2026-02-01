package nx_bit_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNXBitRule(t *testing.T) {
	e2e.RunRuleTests(t, "nx-bit", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-nx-explicit", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-no-nx", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-nx-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-nx-static", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-relocatable.o", Expect: e2e.Skip},

		// x86_64 Clang
		{Binary: "x86_64-clang-nx-explicit", Expect: e2e.Pass},
		{Binary: "x86_64-clang-no-nx", Expect: e2e.Fail},
		{Binary: "x86_64-clang-nx-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-clang-relocatable.o", Expect: e2e.Skip},

		// aarch64 GCC
		{Binary: "aarch64-gcc-nx-explicit", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-no-nx", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-nx-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-nx-static", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-relocatable.o", Expect: e2e.Skip},

		// aarch64 Clang
		{Binary: "aarch64-clang-nx-explicit", Expect: e2e.Pass},
		{Binary: "aarch64-clang-no-nx", Expect: e2e.Fail},
		{Binary: "aarch64-clang-nx-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-clang-relocatable.o", Expect: e2e.Skip},

		// armv7 GCC
		{Binary: "armv7-gcc-nx-explicit", Expect: e2e.Pass},
		{Binary: "armv7-gcc-no-nx", Expect: e2e.Fail},
		{Binary: "armv7-gcc-nx-stripped", Expect: e2e.Pass},
		{Binary: "armv7-gcc-nx-static", Expect: e2e.Pass},
		{Binary: "armv7-gcc-relocatable.o", Expect: e2e.Skip},

		// armv7 Clang
		{Binary: "armv7-clang-nx-explicit", Expect: e2e.Pass},
		{Binary: "armv7-clang-no-nx", Expect: e2e.Fail},
		{Binary: "armv7-clang-nx-stripped", Expect: e2e.Pass},
		{Binary: "armv7-clang-relocatable.o", Expect: e2e.Skip},
	})
}
