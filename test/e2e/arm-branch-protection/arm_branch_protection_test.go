package arm_branch_protection_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestARMBranchProtectionRule(t *testing.T) {
	e2e.RunRuleTests(t, "arm-branch-protection", []e2e.TestCase{
		{Binary: "gcc-branch-protection-standard", Expect: "pass"},
		{Binary: "gcc-branch-protection-pac-ret", Expect: "fail"},
		{Binary: "gcc-branch-protection-bti", Expect: "fail"},
		{Binary: "gcc-no-branch-protection", Expect: "fail"},
		{Binary: "gcc-branch-protection-stripped", Expect: "pass"},
		{Binary: "clang-branch-protection-standard", Expect: "pass"},
		{Binary: "clang-branch-protection-pac-ret", Expect: "fail"},
		{Binary: "clang-branch-protection-bti", Expect: "fail"},
		{Binary: "clang-no-branch-protection", Expect: "fail"},
		{Binary: "clang-branch-protection-stripped", Expect: "pass"},
	})
}
