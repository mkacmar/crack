package relro_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestRELRORule(t *testing.T) {
	e2e.RunRuleTests(t, "relro", relroCases(e2e.Pass, "", ""))
}

func TestFullRELRORule(t *testing.T) {
	e2e.RunRuleTests(t, "full-relro", relroCases(e2e.Fail,
		"Full RELRO not enabled, partial RELRO only",
		"Full RELRO not enabled, no RELRO segment"))
}

// relro has one message per status, so it passes empty strings and asserts the status alone.
func relroCases(partialExpect e2e.Expectation, partialMsg, noRELROMsg string) []e2e.TestCase {
	return []e2e.TestCase{
		{Binary: "amd64-gcc-partial-relro", Expect: partialExpect, Message: partialMsg},
		{Binary: "amd64-gcc-full-relro", Expect: e2e.Pass},
		{Binary: "amd64-gcc-no-relro", Expect: e2e.Fail, Message: noRELROMsg},
		{Binary: "amd64-gcc-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-full-relro-static", Expect: e2e.Pass},
		{Binary: "amd64-gcc-full-relro-shared", Expect: e2e.Pass},
		{Binary: "amd64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "amd64-clang-partial-relro", Expect: partialExpect, Message: partialMsg},
		{Binary: "amd64-clang-full-relro", Expect: e2e.Pass},
		{Binary: "amd64-clang-no-relro", Expect: e2e.Fail, Message: noRELROMsg},
		{Binary: "amd64-clang-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "amd64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-gcc-partial-relro", Expect: partialExpect, Message: partialMsg},
		{Binary: "arm64-gcc-full-relro", Expect: e2e.Pass},
		{Binary: "arm64-gcc-no-relro", Expect: e2e.Fail, Message: noRELROMsg},
		{Binary: "arm64-gcc-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-full-relro-static", Expect: e2e.Pass},
		{Binary: "arm64-gcc-full-relro-shared", Expect: e2e.Pass},
		{Binary: "arm64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-clang-partial-relro", Expect: partialExpect, Message: partialMsg},
		{Binary: "arm64-clang-full-relro", Expect: e2e.Pass},
		{Binary: "arm64-clang-no-relro", Expect: e2e.Fail, Message: noRELROMsg},
		{Binary: "arm64-clang-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "arm64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-gcc-partial-relro", Expect: partialExpect, Message: partialMsg},
		{Binary: "arm-gcc-full-relro", Expect: e2e.Pass},
		{Binary: "arm-gcc-no-relro", Expect: e2e.Fail, Message: noRELROMsg},
		{Binary: "arm-gcc-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-full-relro-static", Expect: e2e.Pass},
		{Binary: "arm-gcc-full-relro-shared", Expect: e2e.Pass},
		{Binary: "arm-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-clang-partial-relro", Expect: partialExpect, Message: partialMsg},
		{Binary: "arm-clang-full-relro", Expect: e2e.Pass},
		{Binary: "arm-clang-no-relro", Expect: e2e.Fail, Message: noRELROMsg},
		{Binary: "arm-clang-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "arm-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "riscv64-gcc-partial-relro", Expect: partialExpect, Message: partialMsg},
		{Binary: "riscv64-gcc-full-relro", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-no-relro", Expect: e2e.Fail, Message: noRELROMsg},
		{Binary: "riscv64-gcc-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-full-relro-static", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-full-relro-shared", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "riscv64-clang-partial-relro", Expect: partialExpect, Message: partialMsg},
		{Binary: "riscv64-clang-full-relro", Expect: e2e.Pass},
		{Binary: "riscv64-clang-no-relro", Expect: e2e.Fail, Message: noRELROMsg},
		{Binary: "riscv64-clang-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "riscv64-clang-relocatable.o", Expect: e2e.Skip},
	}
}
