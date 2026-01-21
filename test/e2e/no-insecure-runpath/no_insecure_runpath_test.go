package no_insecure_runpath_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNoInsecureRUNPATHRule(t *testing.T) {
	e2e.RunRuleTests(t, "no-insecure-runpath", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-no-runpath", Expect: "pass"},
		{Binary: "x86_64-gcc-runpath-absolute", Expect: "pass"},
		{Binary: "x86_64-gcc-runpath-multiple-absolute", Expect: "pass"},
		{Binary: "x86_64-gcc-runpath-dot", Expect: "fail"},
		{Binary: "x86_64-gcc-runpath-dotdot", Expect: "fail"},
		{Binary: "x86_64-gcc-runpath-relative", Expect: "fail"},
		{Binary: "x86_64-gcc-runpath-parent-relative", Expect: "fail"},
		{Binary: "x86_64-gcc-runpath-tmp", Expect: "fail"},
		{Binary: "x86_64-gcc-runpath-var-tmp", Expect: "fail"},
		{Binary: "x86_64-gcc-runpath-tmp-subdir", Expect: "fail"},
		{Binary: "x86_64-gcc-runpath-empty-component", Expect: "fail"},
		{Binary: "x86_64-gcc-runpath-mixed", Expect: "fail"},

		// x86_64 Clang
		{Binary: "x86_64-clang-no-runpath", Expect: "pass"},
		{Binary: "x86_64-clang-runpath-absolute", Expect: "pass"},
		{Binary: "x86_64-clang-runpath-dot", Expect: "fail"},
		{Binary: "x86_64-clang-runpath-tmp", Expect: "fail"},

		// aarch64 GCC
		{Binary: "aarch64-gcc-no-runpath", Expect: "pass"},
		{Binary: "aarch64-gcc-runpath-absolute", Expect: "pass"},
		{Binary: "aarch64-gcc-runpath-multiple-absolute", Expect: "pass"},
		{Binary: "aarch64-gcc-runpath-dot", Expect: "fail"},
		{Binary: "aarch64-gcc-runpath-dotdot", Expect: "fail"},
		{Binary: "aarch64-gcc-runpath-relative", Expect: "fail"},
		{Binary: "aarch64-gcc-runpath-parent-relative", Expect: "fail"},
		{Binary: "aarch64-gcc-runpath-tmp", Expect: "fail"},
		{Binary: "aarch64-gcc-runpath-var-tmp", Expect: "fail"},
		{Binary: "aarch64-gcc-runpath-tmp-subdir", Expect: "fail"},
		{Binary: "aarch64-gcc-runpath-empty-component", Expect: "fail"},
		{Binary: "aarch64-gcc-runpath-mixed", Expect: "fail"},

		// aarch64 Clang
		{Binary: "aarch64-clang-no-runpath", Expect: "pass"},
		{Binary: "aarch64-clang-runpath-absolute", Expect: "pass"},
		{Binary: "aarch64-clang-runpath-dot", Expect: "fail"},
		{Binary: "aarch64-clang-runpath-tmp", Expect: "fail"},

		// armv7 GCC
		{Binary: "armv7-gcc-no-runpath", Expect: "pass"},
		{Binary: "armv7-gcc-runpath-absolute", Expect: "pass"},
		{Binary: "armv7-gcc-runpath-multiple-absolute", Expect: "pass"},
		{Binary: "armv7-gcc-runpath-dot", Expect: "fail"},
		{Binary: "armv7-gcc-runpath-dotdot", Expect: "fail"},
		{Binary: "armv7-gcc-runpath-relative", Expect: "fail"},
		{Binary: "armv7-gcc-runpath-parent-relative", Expect: "fail"},
		{Binary: "armv7-gcc-runpath-tmp", Expect: "fail"},
		{Binary: "armv7-gcc-runpath-var-tmp", Expect: "fail"},
		{Binary: "armv7-gcc-runpath-tmp-subdir", Expect: "fail"},
		{Binary: "armv7-gcc-runpath-empty-component", Expect: "fail"},
		{Binary: "armv7-gcc-runpath-mixed", Expect: "fail"},

		// armv7 Clang
		{Binary: "armv7-clang-no-runpath", Expect: "pass"},
		{Binary: "armv7-clang-runpath-absolute", Expect: "pass"},
		{Binary: "armv7-clang-runpath-dot", Expect: "fail"},
		{Binary: "armv7-clang-runpath-tmp", Expect: "fail"},
	})
}
