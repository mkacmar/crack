package stripped_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestStrippedRule(t *testing.T) {
	const (
		symsAndDebug = "Not stripped, has symbols and debug info"
		symsOnly     = "Not stripped, has symbols"
		debugOnly    = "Partially stripped, has debug info"
	)

	e2e.RunRuleTests(t, "stripped", []e2e.TestCase{
		{Binary: "amd64-gcc-not-stripped", Expect: e2e.Fail, Message: symsAndDebug},
		{Binary: "amd64-gcc-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-strip-debug", Expect: e2e.Fail, Message: symsOnly},
		{Binary: "amd64-gcc-strip-symbols", Expect: e2e.Pass},
		{Binary: "amd64-gcc-link-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-partial-stripped", Expect: e2e.Fail, Message: debugOnly},

		{Binary: "amd64-clang-not-stripped", Expect: e2e.Fail, Message: symsAndDebug},
		{Binary: "amd64-clang-stripped", Expect: e2e.Pass},
		{Binary: "amd64-clang-strip-debug", Expect: e2e.Fail, Message: symsOnly},

		{Binary: "arm64-gcc-not-stripped", Expect: e2e.Fail, Message: symsAndDebug},
		{Binary: "arm64-gcc-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-strip-debug", Expect: e2e.Fail, Message: symsOnly},
		{Binary: "arm64-gcc-strip-symbols", Expect: e2e.Pass},
		{Binary: "arm64-gcc-link-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-partial-stripped", Expect: e2e.Fail, Message: debugOnly},

		{Binary: "arm64-clang-not-stripped", Expect: e2e.Fail, Message: symsAndDebug},
		{Binary: "arm64-clang-stripped", Expect: e2e.Pass},
		{Binary: "arm64-clang-strip-debug", Expect: e2e.Fail, Message: symsOnly},

		{Binary: "arm-gcc-not-stripped", Expect: e2e.Fail, Message: symsAndDebug},
		{Binary: "arm-gcc-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-strip-debug", Expect: e2e.Fail, Message: symsOnly},
		{Binary: "arm-gcc-strip-symbols", Expect: e2e.Pass},
		{Binary: "arm-gcc-link-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-partial-stripped", Expect: e2e.Fail, Message: debugOnly},

		{Binary: "arm-clang-not-stripped", Expect: e2e.Fail, Message: symsAndDebug},
		{Binary: "arm-clang-stripped", Expect: e2e.Pass},
		{Binary: "arm-clang-strip-debug", Expect: e2e.Fail, Message: symsOnly},

		{Binary: "riscv64-gcc-not-stripped", Expect: e2e.Fail, Message: symsAndDebug},
		{Binary: "riscv64-gcc-stripped", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-strip-debug", Expect: e2e.Fail, Message: symsOnly},
		{Binary: "riscv64-gcc-strip-symbols", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-link-stripped", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-partial-stripped", Expect: e2e.Fail, Message: debugOnly},

		{Binary: "riscv64-clang-not-stripped", Expect: e2e.Fail, Message: symsAndDebug},
		{Binary: "riscv64-clang-stripped", Expect: e2e.Pass},
		{Binary: "riscv64-clang-strip-debug", Expect: e2e.Fail, Message: symsOnly},
	})
}
