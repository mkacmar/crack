package stripped_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestStrippedRule(t *testing.T) {
	e2e.RunRuleTests(t, "stripped", []e2e.TestCase{
		{Binary: "amd64-gcc-not-stripped", Expect: e2e.Fail},
		{Binary: "amd64-gcc-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-strip-debug", Expect: e2e.Fail},
		{Binary: "amd64-gcc-strip-symbols", Expect: e2e.Pass},
		{Binary: "amd64-gcc-link-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-partial-stripped", Expect: e2e.Fail},

		{Binary: "amd64-clang-not-stripped", Expect: e2e.Fail},
		{Binary: "amd64-clang-stripped", Expect: e2e.Pass},
		{Binary: "amd64-clang-strip-debug", Expect: e2e.Fail},

		{Binary: "amd64-rustc-not-stripped", Expect: e2e.Fail},
		{Binary: "amd64-rustc-stripped", Expect: e2e.Pass},
		{Binary: "amd64-rustc-strip-debuginfo", Expect: e2e.Fail},

		{Binary: "arm64-gcc-not-stripped", Expect: e2e.Fail},
		{Binary: "arm64-gcc-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-strip-debug", Expect: e2e.Fail},
		{Binary: "arm64-gcc-strip-symbols", Expect: e2e.Pass},
		{Binary: "arm64-gcc-link-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-partial-stripped", Expect: e2e.Fail},

		{Binary: "arm64-clang-not-stripped", Expect: e2e.Fail},
		{Binary: "arm64-clang-stripped", Expect: e2e.Pass},
		{Binary: "arm64-clang-strip-debug", Expect: e2e.Fail},

		{Binary: "arm64-rustc-not-stripped", Expect: e2e.Fail},
		{Binary: "arm64-rustc-stripped", Expect: e2e.Pass},
		{Binary: "arm64-rustc-strip-debuginfo", Expect: e2e.Fail},

		{Binary: "arm-gcc-not-stripped", Expect: e2e.Fail},
		{Binary: "arm-gcc-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-strip-debug", Expect: e2e.Fail},
		{Binary: "arm-gcc-strip-symbols", Expect: e2e.Pass},
		{Binary: "arm-gcc-link-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-partial-stripped", Expect: e2e.Fail},

		{Binary: "arm-clang-not-stripped", Expect: e2e.Fail},
		{Binary: "arm-clang-stripped", Expect: e2e.Pass},
		{Binary: "arm-clang-strip-debug", Expect: e2e.Fail},

		{Binary: "arm-rustc-not-stripped", Expect: e2e.Fail},
		{Binary: "arm-rustc-stripped", Expect: e2e.Pass},
		{Binary: "arm-rustc-strip-debuginfo", Expect: e2e.Fail},
	})
}
