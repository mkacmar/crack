package arm_branch_protection_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestARMBranchProtectionRule(t *testing.T) {
	e2e.RunRuleTests(t, "arm-branch-protection", []e2e.TestCase{
		{Binary: "gcc-branch-protection-standard", Expect: e2e.Pass},
		{Binary: "gcc-branch-protection-pac-ret", Expect: e2e.Fail},
		{Binary: "gcc-branch-protection-bti", Expect: e2e.Fail},
		{Binary: "gcc-no-branch-protection", Expect: e2e.Fail},
		{Binary: "gcc-branch-protection-stripped", Expect: e2e.Pass},
		{Binary: "clang-branch-protection-standard", Expect: e2e.Pass},
		{Binary: "clang-branch-protection-pac-ret", Expect: e2e.Fail},
		{Binary: "clang-branch-protection-bti", Expect: e2e.Fail},
		{Binary: "clang-no-branch-protection", Expect: e2e.Fail},
		{Binary: "clang-branch-protection-stripped", Expect: e2e.Pass},
	})
}
