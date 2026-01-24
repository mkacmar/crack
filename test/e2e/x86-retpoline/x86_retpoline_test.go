package x86_retpoline_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestX86RetpolineRule(t *testing.T) {
	e2e.RunRuleTests(t, "x86-retpoline", []e2e.TestCase{
		{Binary: "gcc-retpoline", Expect: e2e.Pass},
		{Binary: "gcc-no-retpoline", Expect: e2e.Fail},
		{Binary: "gcc-retpoline-stripped", Expect: e2e.Skip},
		{Binary: "gcc-cet-ibt", Expect: e2e.Skip},
		{Binary: "clang-retpoline", Expect: e2e.Pass},
		{Binary: "clang-no-retpoline", Expect: e2e.Fail},
	})
}
