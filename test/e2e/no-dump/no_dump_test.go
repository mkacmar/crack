package nodump_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNoDumpRule(t *testing.T) {
	e2e.RunRuleTests(t, "no-dump", []e2e.TestCase{
		{Binary: "x86_64-gcc-nodump", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-nodump-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-default", Expect: e2e.Fail},

		{Binary: "x86_64-clang-nodump", Expect: e2e.Pass},
		{Binary: "x86_64-clang-nodump-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-clang-default", Expect: e2e.Fail},

		{Binary: "aarch64-gcc-nodump", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-nodump-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-default", Expect: e2e.Fail},

		{Binary: "aarch64-clang-nodump", Expect: e2e.Pass},
		{Binary: "aarch64-clang-nodump-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-clang-default", Expect: e2e.Fail},

		{Binary: "armv7-gcc-nodump", Expect: e2e.Pass},
		{Binary: "armv7-gcc-nodump-stripped", Expect: e2e.Pass},
		{Binary: "armv7-gcc-default", Expect: e2e.Fail},

		{Binary: "armv7-clang-nodump", Expect: e2e.Pass},
		{Binary: "armv7-clang-nodump-stripped", Expect: e2e.Pass},
		{Binary: "armv7-clang-default", Expect: e2e.Fail},
	})
}
