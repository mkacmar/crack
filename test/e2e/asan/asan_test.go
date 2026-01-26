package asan_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestASANRule(t *testing.T) {
	e2e.RunRuleTests(t, "asan", []e2e.TestCase{
		{Binary: "x86_64-gcc-asan", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-asan-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-no-asan", Expect: e2e.Fail},

		{Binary: "x86_64-clang-asan", Expect: e2e.Pass},
		{Binary: "x86_64-clang-asan-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-clang-no-asan", Expect: e2e.Fail},

		{Binary: "aarch64-gcc-asan", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-asan-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-no-asan", Expect: e2e.Fail},

		{Binary: "aarch64-clang-asan", Expect: e2e.Pass},
		{Binary: "aarch64-clang-asan-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-clang-no-asan", Expect: e2e.Fail},

		{Binary: "armv7-gcc-asan", Expect: e2e.Pass},
		{Binary: "armv7-gcc-asan-stripped", Expect: e2e.Pass},
		{Binary: "armv7-gcc-no-asan", Expect: e2e.Fail},

		{Binary: "armv7-clang-asan", Expect: e2e.Pass},
		{Binary: "armv7-clang-asan-stripped", Expect: e2e.Pass},
		{Binary: "armv7-clang-no-asan", Expect: e2e.Fail},
	})
}
