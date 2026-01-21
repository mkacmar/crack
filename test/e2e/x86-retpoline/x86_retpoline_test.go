package x86_retpoline_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestX86RetpolineRule(t *testing.T) {
	e2e.RunRuleTests(t, "x86-retpoline", []e2e.TestCase{
		{Binary: "gcc-retpoline", Expect: "pass"},
		{Binary: "gcc-no-retpoline", Expect: "fail"},
		{Binary: "gcc-retpoline-stripped", Expect: "skip"},
		{Binary: "clang-retpoline", Expect: "pass"},
		{Binary: "clang-no-retpoline", Expect: "fail"},
	})
}
