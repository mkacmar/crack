package asan_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestASANRule(t *testing.T) {
	e2e.RunRuleTests(t, "asan", []e2e.TestCase{
		{Binary: "amd64-gcc-asan", Expect: e2e.Pass},
		{Binary: "amd64-gcc-asan-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-no-asan", Expect: e2e.Fail},

		{Binary: "amd64-clang-asan", Expect: e2e.Pass},
		{Binary: "amd64-clang-asan-stripped", Expect: e2e.Pass},
		{Binary: "amd64-clang-no-asan", Expect: e2e.Fail},

		{Binary: "arm64-gcc-asan", Expect: e2e.Pass},
		{Binary: "arm64-gcc-asan-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-no-asan", Expect: e2e.Fail},

		{Binary: "arm64-clang-asan", Expect: e2e.Pass},
		{Binary: "arm64-clang-asan-stripped", Expect: e2e.Pass},
		{Binary: "arm64-clang-no-asan", Expect: e2e.Fail},

		{Binary: "arm-gcc-asan", Expect: e2e.Pass},
		{Binary: "arm-gcc-asan-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-no-asan", Expect: e2e.Fail},

		{Binary: "arm-clang-asan", Expect: e2e.Pass},
		{Binary: "arm-clang-asan-stripped", Expect: e2e.Pass},
		{Binary: "arm-clang-no-asan", Expect: e2e.Fail},

		{Binary: "amd64-rustc-no-asan", Expect: e2e.Skip},
		{Binary: "arm64-rustc-no-asan", Expect: e2e.Skip},
		{Binary: "arm-rustc-no-asan", Expect: e2e.Skip},
	})
}
