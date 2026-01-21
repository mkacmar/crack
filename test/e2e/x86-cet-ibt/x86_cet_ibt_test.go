package x86_cet_ibt_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestX86CETIBTRule(t *testing.T) {
	e2e.RunRuleTests(t, "x86-cet-ibt", []e2e.TestCase{
		{Binary: "gcc-cet-full", Expect: "pass"},
		{Binary: "gcc-cet-branch", Expect: "pass"},
		{Binary: "gcc-cet-none", Expect: "fail"},
		{Binary: "gcc-cet-full-stripped", Expect: "pass"},
		{Binary: "clang-cet-full", Expect: "pass"},
		{Binary: "clang-cet-branch", Expect: "pass"},
		{Binary: "clang-cet-none", Expect: "fail"},
	})
}
