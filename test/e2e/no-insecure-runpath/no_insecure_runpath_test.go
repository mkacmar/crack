package no_insecure_runpath_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNoInsecureRUNPATHRule(t *testing.T) {
	e2e.RunRuleTests(t, "no-insecure-runpath", []e2e.TestCase{
		// x86_64 GCC
		{"x86_64-gcc-no-runpath", "pass"},
		{"x86_64-gcc-runpath-absolute", "pass"},
		{"x86_64-gcc-runpath-multiple-absolute", "pass"},
		{"x86_64-gcc-runpath-dot", "fail"},
		{"x86_64-gcc-runpath-dotdot", "fail"},
		{"x86_64-gcc-runpath-relative", "fail"},
		{"x86_64-gcc-runpath-parent-relative", "fail"},
		{"x86_64-gcc-runpath-tmp", "fail"},
		{"x86_64-gcc-runpath-var-tmp", "fail"},
		{"x86_64-gcc-runpath-tmp-subdir", "fail"},
		{"x86_64-gcc-runpath-empty-component", "fail"},
		{"x86_64-gcc-runpath-mixed", "fail"},

		// x86_64 Clang
		{"x86_64-clang-no-runpath", "pass"},
		{"x86_64-clang-runpath-absolute", "pass"},
		{"x86_64-clang-runpath-dot", "fail"},
		{"x86_64-clang-runpath-tmp", "fail"},

		// aarch64 GCC
		{"aarch64-gcc-no-runpath", "pass"},
		{"aarch64-gcc-runpath-absolute", "pass"},
		{"aarch64-gcc-runpath-multiple-absolute", "pass"},
		{"aarch64-gcc-runpath-dot", "fail"},
		{"aarch64-gcc-runpath-dotdot", "fail"},
		{"aarch64-gcc-runpath-relative", "fail"},
		{"aarch64-gcc-runpath-parent-relative", "fail"},
		{"aarch64-gcc-runpath-tmp", "fail"},
		{"aarch64-gcc-runpath-var-tmp", "fail"},
		{"aarch64-gcc-runpath-tmp-subdir", "fail"},
		{"aarch64-gcc-runpath-empty-component", "fail"},
		{"aarch64-gcc-runpath-mixed", "fail"},

		// aarch64 Clang
		{"aarch64-clang-no-runpath", "pass"},
		{"aarch64-clang-runpath-absolute", "pass"},
		{"aarch64-clang-runpath-dot", "fail"},
		{"aarch64-clang-runpath-tmp", "fail"},

		// armv7 GCC
		{"armv7-gcc-no-runpath", "pass"},
		{"armv7-gcc-runpath-absolute", "pass"},
		{"armv7-gcc-runpath-multiple-absolute", "pass"},
		{"armv7-gcc-runpath-dot", "fail"},
		{"armv7-gcc-runpath-dotdot", "fail"},
		{"armv7-gcc-runpath-relative", "fail"},
		{"armv7-gcc-runpath-parent-relative", "fail"},
		{"armv7-gcc-runpath-tmp", "fail"},
		{"armv7-gcc-runpath-var-tmp", "fail"},
		{"armv7-gcc-runpath-tmp-subdir", "fail"},
		{"armv7-gcc-runpath-empty-component", "fail"},
		{"armv7-gcc-runpath-mixed", "fail"},

		// armv7 Clang
		{"armv7-clang-no-runpath", "pass"},
		{"armv7-clang-runpath-absolute", "pass"},
		{"armv7-clang-runpath-dot", "fail"},
		{"armv7-clang-runpath-tmp", "fail"},
	})
}
