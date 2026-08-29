package cfi_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestCFIRule(t *testing.T) {
	const crossDSO = "CFI enabled (cross-DSO mode)"

	e2e.RunRuleTests(t, "cfi", []e2e.TestCase{
		{Binary: "amd64-clang-cfi", Expect: e2e.Pass, Message: crossDSO},
		{Binary: "amd64-clang-cfi-stripped", Expect: e2e.Pass, Message: crossDSO},
		{Binary: "amd64-clang-no-cfi", Expect: e2e.Fail},
		{Binary: "amd64-gcc-no-cfi", Expect: e2e.Skip},

		{Binary: "arm64-clang-cfi", Expect: e2e.Pass, Message: crossDSO},
		{Binary: "arm64-clang-cfi-stripped", Expect: e2e.Pass, Message: crossDSO},
		{Binary: "arm64-clang-no-cfi", Expect: e2e.Fail},
		{Binary: "arm64-gcc-no-cfi", Expect: e2e.Skip},

		{Binary: "amd64-clang-no-cfi-stripped", Expect: e2e.Skip},
		{Binary: "amd64-gcc-no-cfi-stripped", Expect: e2e.Skip},
		{Binary: "arm64-clang-no-cfi-stripped", Expect: e2e.Skip},
		{Binary: "arm64-gcc-no-cfi-stripped", Expect: e2e.Skip},
	})
}
