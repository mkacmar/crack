package ubsan_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestUBSanRule(t *testing.T) {
	e2e.RunRuleTests(t, "ubsan", []e2e.TestCase{
		{Binary: "amd64-gcc-ubsan", Expect: e2e.Pass},
		{Binary: "amd64-gcc-ubsan-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-no-ubsan", Expect: e2e.Fail},

		{Binary: "amd64-clang-ubsan", Expect: e2e.Pass},
		{Binary: "amd64-clang-ubsan-stripped", Expect: e2e.Pass},
		{Binary: "amd64-clang-no-ubsan", Expect: e2e.Fail},

		{Binary: "amd64-rustc-no-ubsan", Expect: e2e.Skip},

		{Binary: "arm64-gcc-ubsan", Expect: e2e.Pass},
		{Binary: "arm64-gcc-ubsan-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-no-ubsan", Expect: e2e.Fail},

		{Binary: "arm64-clang-ubsan", Expect: e2e.Pass},
		{Binary: "arm64-clang-ubsan-stripped", Expect: e2e.Pass},
		{Binary: "arm64-clang-no-ubsan", Expect: e2e.Fail},

		{Binary: "arm64-rustc-no-ubsan", Expect: e2e.Skip},

		{Binary: "arm-gcc-ubsan", Expect: e2e.Pass},
		{Binary: "arm-gcc-ubsan-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-no-ubsan", Expect: e2e.Fail},

		{Binary: "arm-clang-ubsan", Expect: e2e.Pass},
		{Binary: "arm-clang-ubsan-stripped", Expect: e2e.Pass},
		{Binary: "arm-clang-no-ubsan", Expect: e2e.Fail},

		{Binary: "arm-rustc-no-ubsan", Expect: e2e.Skip},
	})
}
