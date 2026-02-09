package x86_cet_ibt_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestX86CETIBTRule(t *testing.T) {
	e2e.RunRuleTests(t, "x86-cet-ibt", []e2e.TestCase{
		{Binary: "gcc-cet-full", Expect: e2e.Pass},
		{Binary: "gcc-cet-branch", Expect: e2e.Pass},
		{Binary: "gcc-cet-none", Expect: e2e.Fail},
		{Binary: "gcc-cet-full-stripped", Expect: e2e.Pass},
		{Binary: "clang-cet-full", Expect: e2e.Pass},
		{Binary: "clang-cet-branch", Expect: e2e.Pass},
		{Binary: "clang-cet-none", Expect: e2e.Fail},
		{Binary: "rustc-default", Expect: e2e.Skip},
	})
}
