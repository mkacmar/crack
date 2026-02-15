package stack_canary_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestStackCanaryRule(t *testing.T) {
	e2e.RunRuleTests(t, "stack-canary", []e2e.TestCase{
		{Binary: "amd64-gcc-stack-protector-strong", Expect: e2e.Pass},
		{Binary: "amd64-gcc-stack-protector-all", Expect: e2e.Pass},
		{Binary: "amd64-gcc-stack-protector", Expect: e2e.Pass},
		{Binary: "amd64-gcc-no-stack-protector", Expect: e2e.Fail},
		{Binary: "amd64-gcc-stack-protector-simple", Expect: e2e.Fail},
		{Binary: "amd64-gcc-stack-protector-all-simple", Expect: e2e.Pass},
		{Binary: "amd64-gcc-stack-protector-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-stack-protector-static", Expect: e2e.Pass},
		{Binary: "amd64-gcc-stack-protector-static-stripped", Expect: e2e.Fail},
		{Binary: "amd64-gcc-stack-protector-lto", Expect: e2e.Pass},

		{Binary: "amd64-clang-stack-protector-strong", Expect: e2e.Pass},
		{Binary: "amd64-clang-stack-protector-all", Expect: e2e.Pass},
		{Binary: "amd64-clang-no-stack-protector", Expect: e2e.Fail},
		{Binary: "amd64-clang-stack-protector-stripped", Expect: e2e.Pass},
		{Binary: "amd64-clang-stack-protector-static", Expect: e2e.Pass},
		{Binary: "amd64-clang-stack-protector-static-stripped", Expect: e2e.Fail},
		{Binary: "amd64-clang-stack-protector-lto", Expect: e2e.Pass},

		{Binary: "arm64-gcc-stack-protector-strong", Expect: e2e.Pass},
		{Binary: "arm64-gcc-stack-protector-all", Expect: e2e.Pass},
		{Binary: "arm64-gcc-stack-protector", Expect: e2e.Pass},
		{Binary: "arm64-gcc-no-stack-protector", Expect: e2e.Fail},
		{Binary: "arm64-gcc-stack-protector-simple", Expect: e2e.Fail},
		{Binary: "arm64-gcc-stack-protector-all-simple", Expect: e2e.Pass},
		{Binary: "arm64-gcc-stack-protector-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-stack-protector-static", Expect: e2e.Pass},
		{Binary: "arm64-gcc-stack-protector-static-stripped", Expect: e2e.Fail},
		{Binary: "arm64-gcc-stack-protector-lto", Expect: e2e.Pass},

		{Binary: "arm64-clang-stack-protector-strong", Expect: e2e.Pass},
		{Binary: "arm64-clang-stack-protector-all", Expect: e2e.Pass},
		{Binary: "arm64-clang-no-stack-protector", Expect: e2e.Fail},
		{Binary: "arm64-clang-stack-protector-stripped", Expect: e2e.Pass},
		{Binary: "arm64-clang-stack-protector-static", Expect: e2e.Pass},
		{Binary: "arm64-clang-stack-protector-static-stripped", Expect: e2e.Fail},
		{Binary: "arm64-clang-stack-protector-lto", Expect: e2e.Pass},

		{Binary: "arm-gcc-stack-protector-strong", Expect: e2e.Pass},
		{Binary: "arm-gcc-stack-protector-all", Expect: e2e.Pass},
		{Binary: "arm-gcc-stack-protector", Expect: e2e.Pass},
		{Binary: "arm-gcc-no-stack-protector", Expect: e2e.Fail},
		{Binary: "arm-gcc-stack-protector-simple", Expect: e2e.Fail},
		{Binary: "arm-gcc-stack-protector-all-simple", Expect: e2e.Pass},
		{Binary: "arm-gcc-stack-protector-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-stack-protector-static", Expect: e2e.Pass},
		{Binary: "arm-gcc-stack-protector-static-stripped", Expect: e2e.Fail},
		{Binary: "arm-gcc-stack-protector-lto", Expect: e2e.Pass},

		{Binary: "arm-clang-stack-protector-strong", Expect: e2e.Pass},
		{Binary: "arm-clang-stack-protector-all", Expect: e2e.Pass},
		{Binary: "arm-clang-no-stack-protector", Expect: e2e.Fail},
		{Binary: "arm-clang-stack-protector-stripped", Expect: e2e.Pass},
		{Binary: "arm-clang-stack-protector-static", Expect: e2e.Pass},
		{Binary: "arm-clang-stack-protector-static-stripped", Expect: e2e.Fail},
		{Binary: "arm-clang-stack-protector-lto", Expect: e2e.Pass},

		{Binary: "amd64-rustc-stack-protector", Expect: e2e.Skip},
		{Binary: "amd64-rustc-stack-protector-stripped", Expect: e2e.Skip},
		{Binary: "arm64-rustc-stack-protector", Expect: e2e.Skip},
		{Binary: "arm64-rustc-stack-protector-stripped", Expect: e2e.Skip},
		{Binary: "arm-rustc-stack-protector", Expect: e2e.Skip},
		{Binary: "arm-rustc-stack-protector-stripped", Expect: e2e.Skip},
	})
}
