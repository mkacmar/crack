package x86_retpoline_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestX86RetpolineRule(t *testing.T) {
	e2e.RunRuleTests(t, "x86-retpoline", []e2e.TestCase{
		{Binary: "gcc-retpoline", Expect: e2e.Pass, Message: "Retpoline enabled (GCC)"},
		{Binary: "gcc-no-retpoline", Expect: e2e.Fail},
		{Binary: "gcc-retpoline-stripped", Expect: e2e.Skip, Message: "Stripped binary, retpoline detection limited"},
		{Binary: "gcc-cet-ibt", Expect: e2e.Skip, Message: "CET IBT enabled, retpoline not needed"},
		{Binary: "clang-retpoline", Expect: e2e.Pass, Message: "Retpoline enabled (LLVM)"},
		{Binary: "clang-no-retpoline", Expect: e2e.Fail},
	})
}
