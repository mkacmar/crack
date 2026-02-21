package stack_limit_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestStackLimitRule(t *testing.T) {
	e2e.RunRuleTests(t, "stack-limit", []e2e.TestCase{
		{Binary: "amd64-gcc-stack-limit", Expect: e2e.Pass},
		{Binary: "amd64-gcc-no-stack-limit", Expect: e2e.Fail},
		{Binary: "amd64-gcc-relocatable.o", Expect: e2e.Skip},
		{Binary: "amd64-clang-stack-limit", Expect: e2e.Pass},
		{Binary: "amd64-clang-no-stack-limit", Expect: e2e.Fail},
		{Binary: "amd64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-gcc-stack-limit", Expect: e2e.Pass},
		{Binary: "arm64-gcc-no-stack-limit", Expect: e2e.Fail},
		{Binary: "arm64-gcc-relocatable.o", Expect: e2e.Skip},
		{Binary: "arm64-clang-stack-limit", Expect: e2e.Pass},
		{Binary: "arm64-clang-no-stack-limit", Expect: e2e.Fail},
		{Binary: "arm64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-gcc-stack-limit", Expect: e2e.Pass},
		{Binary: "arm-gcc-no-stack-limit", Expect: e2e.Fail},
		{Binary: "arm-gcc-relocatable.o", Expect: e2e.Skip},
		{Binary: "arm-clang-stack-limit", Expect: e2e.Pass},
		{Binary: "arm-clang-no-stack-limit", Expect: e2e.Fail},
		{Binary: "arm-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "amd64-rustc-no-stack-limit", Expect: e2e.Skip},
		{Binary: "arm64-rustc-no-stack-limit", Expect: e2e.Skip},
		{Binary: "arm-rustc-no-stack-limit", Expect: e2e.Skip},
	})
}
