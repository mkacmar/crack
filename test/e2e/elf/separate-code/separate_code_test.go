package separate_code_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestSeparateCodeRule(t *testing.T) {
	const notExec = "Not an executable or shared library"

	e2e.RunRuleTests(t, "separate-code", []e2e.TestCase{
		{Binary: "amd64-gcc-separate-code", Expect: e2e.Pass},
		{Binary: "amd64-gcc-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-separate-code-static", Expect: e2e.Pass},
		{Binary: "amd64-gcc-separate-code-shared", Expect: e2e.Pass},
		{Binary: "amd64-gcc-relocatable.o", Expect: e2e.Skip, Message: notExec},

		{Binary: "amd64-clang-separate-code", Expect: e2e.Pass},
		{Binary: "amd64-clang-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "amd64-clang-relocatable.o", Expect: e2e.Skip, Message: notExec},

		{Binary: "arm64-gcc-separate-code", Expect: e2e.Pass},
		{Binary: "arm64-gcc-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-separate-code-static", Expect: e2e.Pass},
		{Binary: "arm64-gcc-separate-code-shared", Expect: e2e.Pass},
		{Binary: "arm64-gcc-relocatable.o", Expect: e2e.Skip, Message: notExec},

		{Binary: "arm64-clang-separate-code", Expect: e2e.Pass},
		{Binary: "arm64-clang-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "arm64-clang-relocatable.o", Expect: e2e.Skip, Message: notExec},

		{Binary: "arm-gcc-separate-code", Expect: e2e.Pass},
		{Binary: "arm-gcc-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-separate-code-static", Expect: e2e.Pass},
		{Binary: "arm-gcc-separate-code-shared", Expect: e2e.Pass},
		{Binary: "arm-gcc-relocatable.o", Expect: e2e.Skip, Message: notExec},

		{Binary: "arm-clang-separate-code", Expect: e2e.Pass},
		{Binary: "arm-clang-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "arm-clang-relocatable.o", Expect: e2e.Skip, Message: notExec},

		{Binary: "riscv64-gcc-separate-code", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-separate-code-static", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-separate-code-shared", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-relocatable.o", Expect: e2e.Skip, Message: notExec},

		{Binary: "riscv64-clang-separate-code", Expect: e2e.Pass},
		{Binary: "riscv64-clang-separate-code-stripped", Expect: e2e.Pass},
		{Binary: "riscv64-clang-relocatable.o", Expect: e2e.Skip, Message: notExec},

		{Binary: "amd64-gcc4-no-separate-code", Expect: e2e.Fail},
		{Binary: "arm64-gcc7-no-separate-code", Expect: e2e.Fail},
	})
}
