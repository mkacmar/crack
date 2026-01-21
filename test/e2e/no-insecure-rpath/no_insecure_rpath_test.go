package no_insecure_rpath_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNoInsecureRPATHRule(t *testing.T) {
	e2e.RunRuleTests(t, "no-insecure-rpath", []e2e.TestCase{
		// x86_64 GCC
		{"x86_64-gcc-no-rpath", "pass"},
		{"x86_64-gcc-rpath-absolute", "pass"},
		{"x86_64-gcc-rpath-multiple-absolute", "pass"},
		{"x86_64-gcc-rpath-dot", "fail"},
		{"x86_64-gcc-rpath-dotdot", "fail"},
		{"x86_64-gcc-rpath-relative", "fail"},
		{"x86_64-gcc-rpath-parent-relative", "fail"},
		{"x86_64-gcc-rpath-tmp", "fail"},
		{"x86_64-gcc-rpath-var-tmp", "fail"},
		{"x86_64-gcc-rpath-tmp-subdir", "fail"},
		{"x86_64-gcc-rpath-empty-component", "fail"},
		{"x86_64-gcc-rpath-mixed", "fail"},

		// x86_64 Clang
		{"x86_64-clang-no-rpath", "pass"},
		{"x86_64-clang-rpath-absolute", "pass"},
		{"x86_64-clang-rpath-dot", "fail"},
		{"x86_64-clang-rpath-tmp", "fail"},

		// aarch64 GCC
		{"aarch64-gcc-no-rpath", "pass"},
		{"aarch64-gcc-rpath-absolute", "pass"},
		{"aarch64-gcc-rpath-multiple-absolute", "pass"},
		{"aarch64-gcc-rpath-dot", "fail"},
		{"aarch64-gcc-rpath-dotdot", "fail"},
		{"aarch64-gcc-rpath-relative", "fail"},
		{"aarch64-gcc-rpath-parent-relative", "fail"},
		{"aarch64-gcc-rpath-tmp", "fail"},
		{"aarch64-gcc-rpath-var-tmp", "fail"},
		{"aarch64-gcc-rpath-tmp-subdir", "fail"},
		{"aarch64-gcc-rpath-empty-component", "fail"},
		{"aarch64-gcc-rpath-mixed", "fail"},

		// aarch64 Clang
		{"aarch64-clang-no-rpath", "pass"},
		{"aarch64-clang-rpath-absolute", "pass"},
		{"aarch64-clang-rpath-dot", "fail"},
		{"aarch64-clang-rpath-tmp", "fail"},

		// armv7 GCC
		{"armv7-gcc-no-rpath", "pass"},
		{"armv7-gcc-rpath-absolute", "pass"},
		{"armv7-gcc-rpath-multiple-absolute", "pass"},
		{"armv7-gcc-rpath-dot", "fail"},
		{"armv7-gcc-rpath-dotdot", "fail"},
		{"armv7-gcc-rpath-relative", "fail"},
		{"armv7-gcc-rpath-parent-relative", "fail"},
		{"armv7-gcc-rpath-tmp", "fail"},
		{"armv7-gcc-rpath-var-tmp", "fail"},
		{"armv7-gcc-rpath-tmp-subdir", "fail"},
		{"armv7-gcc-rpath-empty-component", "fail"},
		{"armv7-gcc-rpath-mixed", "fail"},

		// armv7 Clang
		{"armv7-clang-no-rpath", "pass"},
		{"armv7-clang-rpath-absolute", "pass"},
		{"armv7-clang-rpath-dot", "fail"},
		{"armv7-clang-rpath-tmp", "fail"},
	})
}
