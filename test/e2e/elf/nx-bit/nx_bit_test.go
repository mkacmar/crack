package nx_bit_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestNXBitRule(t *testing.T) {
	const execStack = "NX not enabled, stack executable"

	e2e.RunRuleTests(t, "nx-bit", []e2e.TestCase{
		{Binary: "amd64-gcc-nx-explicit", Expect: e2e.Pass},
		{Binary: "amd64-gcc-no-nx", Expect: e2e.Fail, Message: execStack},
		{Binary: "amd64-gcc-nx-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-nx-static", Expect: e2e.Pass},
		{Binary: "amd64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "amd64-clang-nx-explicit", Expect: e2e.Pass},
		{Binary: "amd64-clang-no-nx", Expect: e2e.Fail, Message: execStack},
		{Binary: "amd64-clang-nx-stripped", Expect: e2e.Pass},
		{Binary: "amd64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-gcc-nx-explicit", Expect: e2e.Pass},
		{Binary: "arm64-gcc-no-nx", Expect: e2e.Fail, Message: execStack},
		{Binary: "arm64-gcc-nx-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-nx-static", Expect: e2e.Pass},
		{Binary: "arm64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-clang-nx-explicit", Expect: e2e.Pass},
		{Binary: "arm64-clang-no-nx", Expect: e2e.Fail, Message: execStack},
		{Binary: "arm64-clang-nx-stripped", Expect: e2e.Pass},
		{Binary: "arm64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-gcc-nx-explicit", Expect: e2e.Pass},
		{Binary: "arm-gcc-no-nx", Expect: e2e.Fail, Message: execStack},
		{Binary: "arm-gcc-nx-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-nx-static", Expect: e2e.Pass},
		{Binary: "arm-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-clang-nx-explicit", Expect: e2e.Pass},
		{Binary: "arm-clang-no-nx", Expect: e2e.Fail, Message: execStack},
		{Binary: "arm-clang-nx-stripped", Expect: e2e.Pass},
		{Binary: "arm-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "riscv64-gcc-nx-explicit", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-no-nx", Expect: e2e.Fail, Message: execStack},
		{Binary: "riscv64-gcc-nx-stripped", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-nx-static", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "riscv64-clang-nx-explicit", Expect: e2e.Pass},
		{Binary: "riscv64-clang-no-nx", Expect: e2e.Fail, Message: execStack},
		{Binary: "riscv64-clang-nx-stripped", Expect: e2e.Pass},
		{Binary: "riscv64-clang-relocatable.o", Expect: e2e.Skip},
	})
}
