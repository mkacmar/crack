package x86_cet_shstk_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestX86CETShadowStackRule(t *testing.T) {
	e2e.RunRuleTests(t, "x86-cet-shstk", []e2e.TestCase{
		{Binary: "gcc-cet-full", Expect: e2e.Pass},
		{Binary: "gcc-cet-return", Expect: e2e.Pass},
		{Binary: "gcc-cet-none", Expect: e2e.Fail},
		{Binary: "gcc-cet-full-stripped", Expect: e2e.Pass},
		{Binary: "clang-cet-full", Expect: e2e.Pass},
		{Binary: "clang-cet-return", Expect: e2e.Pass},
		{Binary: "clang-cet-none", Expect: e2e.Fail},
		{Binary: "rustc-default", Expect: e2e.Skip},
	})
}
