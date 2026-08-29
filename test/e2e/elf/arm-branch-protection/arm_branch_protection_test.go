package arm_branch_protection_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestARMBranchProtectionRule(t *testing.T) {
	const (
		btiMissing   = "ARM branch protection partial, BTI missing"
		pacMissing   = "ARM branch protection partial, PAC missing (libc may lack PAC support)"
		noProtection = "ARM branch protection not enabled"
	)

	e2e.RunRuleTests(t, "arm-branch-protection", []e2e.TestCase{
		{Binary: "gcc-branch-protection-standard", Expect: e2e.Pass},
		{Binary: "gcc-branch-protection-pac-ret", Expect: e2e.Fail, Message: btiMissing},
		{Binary: "gcc-branch-protection-bti", Expect: e2e.Fail, Message: pacMissing},
		{Binary: "gcc-no-branch-protection", Expect: e2e.Fail, Message: noProtection},
		{Binary: "gcc-branch-protection-stripped", Expect: e2e.Pass},
		{Binary: "clang-branch-protection-standard", Expect: e2e.Pass},
		{Binary: "clang-branch-protection-pac-ret", Expect: e2e.Fail, Message: btiMissing},
		{Binary: "clang-branch-protection-bti", Expect: e2e.Fail, Message: pacMissing},
		{Binary: "clang-no-branch-protection", Expect: e2e.Fail, Message: noProtection},
		{Binary: "clang-branch-protection-stripped", Expect: e2e.Pass},
	})
}
