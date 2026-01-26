package noplt_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNoPLTRule(t *testing.T) {
	e2e.RunRuleTests(t, "no-plt", []e2e.TestCase{
		{Binary: "x86_64-gcc-no-plt", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-no-plt-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-default", Expect: e2e.Fail},

		{Binary: "x86_64-clang-no-plt", Expect: e2e.Pass},
		{Binary: "x86_64-clang-no-plt-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-clang-default", Expect: e2e.Fail},

		{Binary: "aarch64-gcc-no-plt", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-no-plt-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-default", Expect: e2e.Fail},

		{Binary: "aarch64-clang-no-plt", Expect: e2e.Pass},
		{Binary: "aarch64-clang-no-plt-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-clang-default", Expect: e2e.Fail},

		{Binary: "armv7-gcc-no-plt", Expect: e2e.Pass},
		{Binary: "armv7-gcc-no-plt-stripped", Expect: e2e.Pass},
		{Binary: "armv7-gcc-default", Expect: e2e.Fail},

		{Binary: "armv7-clang-no-plt", Expect: e2e.Pass},
		{Binary: "armv7-clang-no-plt-stripped", Expect: e2e.Pass},
		{Binary: "armv7-clang-default", Expect: e2e.Fail},
	})
}
