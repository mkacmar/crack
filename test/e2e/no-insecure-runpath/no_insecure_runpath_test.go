package no_insecure_runpath_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNoInsecureRUNPATHRule(t *testing.T) {
	e2e.RunRuleTests(t, "no-insecure-runpath", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-no-runpath", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-runpath-absolute", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-runpath-multiple-absolute", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-runpath-dot", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-runpath-dotdot", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-runpath-relative", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-runpath-parent-relative", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-runpath-tmp", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-runpath-var-tmp", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-runpath-tmp-subdir", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-runpath-empty-component", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-runpath-mixed", Expect: e2e.Fail},

		// x86_64 Clang
		{Binary: "x86_64-clang-no-runpath", Expect: e2e.Pass},
		{Binary: "x86_64-clang-runpath-absolute", Expect: e2e.Pass},
		{Binary: "x86_64-clang-runpath-dot", Expect: e2e.Fail},
		{Binary: "x86_64-clang-runpath-tmp", Expect: e2e.Fail},

		// aarch64 GCC
		{Binary: "aarch64-gcc-no-runpath", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-runpath-absolute", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-runpath-multiple-absolute", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-runpath-dot", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-runpath-dotdot", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-runpath-relative", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-runpath-parent-relative", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-runpath-tmp", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-runpath-var-tmp", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-runpath-tmp-subdir", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-runpath-empty-component", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-runpath-mixed", Expect: e2e.Fail},

		// aarch64 Clang
		{Binary: "aarch64-clang-no-runpath", Expect: e2e.Pass},
		{Binary: "aarch64-clang-runpath-absolute", Expect: e2e.Pass},
		{Binary: "aarch64-clang-runpath-dot", Expect: e2e.Fail},
		{Binary: "aarch64-clang-runpath-tmp", Expect: e2e.Fail},

		// armv7 GCC
		{Binary: "armv7-gcc-no-runpath", Expect: e2e.Pass},
		{Binary: "armv7-gcc-runpath-absolute", Expect: e2e.Pass},
		{Binary: "armv7-gcc-runpath-multiple-absolute", Expect: e2e.Pass},
		{Binary: "armv7-gcc-runpath-dot", Expect: e2e.Fail},
		{Binary: "armv7-gcc-runpath-dotdot", Expect: e2e.Fail},
		{Binary: "armv7-gcc-runpath-relative", Expect: e2e.Fail},
		{Binary: "armv7-gcc-runpath-parent-relative", Expect: e2e.Fail},
		{Binary: "armv7-gcc-runpath-tmp", Expect: e2e.Fail},
		{Binary: "armv7-gcc-runpath-var-tmp", Expect: e2e.Fail},
		{Binary: "armv7-gcc-runpath-tmp-subdir", Expect: e2e.Fail},
		{Binary: "armv7-gcc-runpath-empty-component", Expect: e2e.Fail},
		{Binary: "armv7-gcc-runpath-mixed", Expect: e2e.Fail},

		// armv7 Clang
		{Binary: "armv7-clang-no-runpath", Expect: e2e.Pass},
		{Binary: "armv7-clang-runpath-absolute", Expect: e2e.Pass},
		{Binary: "armv7-clang-runpath-dot", Expect: e2e.Fail},
		{Binary: "armv7-clang-runpath-tmp", Expect: e2e.Fail},
	})
}
