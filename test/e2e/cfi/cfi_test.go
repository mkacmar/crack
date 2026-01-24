package cfi_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestCFIRule(t *testing.T) {
	e2e.RunRuleTests(t, "cfi", []e2e.TestCase{
		// x86_64 Clang
		{Binary: "x86_64-clang-cfi", Expect: e2e.Pass},
		{Binary: "x86_64-clang-cfi-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-clang-no-cfi", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-no-cfi", Expect: e2e.Fail},

		// aarch64 Clang
		{Binary: "aarch64-clang-cfi", Expect: e2e.Pass},
		{Binary: "aarch64-clang-cfi-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-clang-no-cfi", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-no-cfi", Expect: e2e.Fail},
	})
}
