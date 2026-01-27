package noplt_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNoPLTRule(t *testing.T) {
	e2e.RunRuleTests(t, "no-plt", []e2e.TestCase{
		{Binary: "gcc-no-plt-cet", Expect: e2e.Skip},
		{Binary: "gcc-no-plt", Expect: e2e.Pass},
		{Binary: "gcc-no-plt-stripped", Expect: e2e.Pass},
		{Binary: "clang-no-plt", Expect: e2e.Pass},
		{Binary: "clang-no-plt-stripped", Expect: e2e.Pass},
		{Binary: "gcc-plt", Expect: e2e.Fail},
		{Binary: "clang-plt", Expect: e2e.Fail},

		{Binary: "i386-gcc-no-plt-cet", Expect: e2e.Skip},
		{Binary: "i386-gcc-no-plt", Expect: e2e.Pass},
		{Binary: "i386-gcc-plt", Expect: e2e.Fail},
		{Binary: "i386-clang-no-plt", Expect: e2e.Pass},
		{Binary: "i386-clang-plt", Expect: e2e.Fail},
	})
}
