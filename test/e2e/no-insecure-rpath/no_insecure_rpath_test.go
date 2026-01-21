package no_insecure_rpath_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNoInsecureRPATHRule(t *testing.T) {
	e2e.RunRuleTests(t, "no-insecure-rpath", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-no-rpath", Expect: "pass"},
		{Binary: "x86_64-gcc-rpath-absolute", Expect: "pass"},
		{Binary: "x86_64-gcc-rpath-multiple-absolute", Expect: "pass"},
		{Binary: "x86_64-gcc-rpath-dot", Expect: "fail"},
		{Binary: "x86_64-gcc-rpath-dotdot", Expect: "fail"},
		{Binary: "x86_64-gcc-rpath-relative", Expect: "fail"},
		{Binary: "x86_64-gcc-rpath-parent-relative", Expect: "fail"},
		{Binary: "x86_64-gcc-rpath-tmp", Expect: "fail"},
		{Binary: "x86_64-gcc-rpath-var-tmp", Expect: "fail"},
		{Binary: "x86_64-gcc-rpath-tmp-subdir", Expect: "fail"},
		{Binary: "x86_64-gcc-rpath-empty-component", Expect: "fail"},
		{Binary: "x86_64-gcc-rpath-mixed", Expect: "fail"},

		// x86_64 Clang
		{Binary: "x86_64-clang-no-rpath", Expect: "pass"},
		{Binary: "x86_64-clang-rpath-absolute", Expect: "pass"},
		{Binary: "x86_64-clang-rpath-dot", Expect: "fail"},
		{Binary: "x86_64-clang-rpath-tmp", Expect: "fail"},

		// aarch64 GCC
		{Binary: "aarch64-gcc-no-rpath", Expect: "pass"},
		{Binary: "aarch64-gcc-rpath-absolute", Expect: "pass"},
		{Binary: "aarch64-gcc-rpath-multiple-absolute", Expect: "pass"},
		{Binary: "aarch64-gcc-rpath-dot", Expect: "fail"},
		{Binary: "aarch64-gcc-rpath-dotdot", Expect: "fail"},
		{Binary: "aarch64-gcc-rpath-relative", Expect: "fail"},
		{Binary: "aarch64-gcc-rpath-parent-relative", Expect: "fail"},
		{Binary: "aarch64-gcc-rpath-tmp", Expect: "fail"},
		{Binary: "aarch64-gcc-rpath-var-tmp", Expect: "fail"},
		{Binary: "aarch64-gcc-rpath-tmp-subdir", Expect: "fail"},
		{Binary: "aarch64-gcc-rpath-empty-component", Expect: "fail"},
		{Binary: "aarch64-gcc-rpath-mixed", Expect: "fail"},

		// aarch64 Clang
		{Binary: "aarch64-clang-no-rpath", Expect: "pass"},
		{Binary: "aarch64-clang-rpath-absolute", Expect: "pass"},
		{Binary: "aarch64-clang-rpath-dot", Expect: "fail"},
		{Binary: "aarch64-clang-rpath-tmp", Expect: "fail"},

		// armv7 GCC
		{Binary: "armv7-gcc-no-rpath", Expect: "pass"},
		{Binary: "armv7-gcc-rpath-absolute", Expect: "pass"},
		{Binary: "armv7-gcc-rpath-multiple-absolute", Expect: "pass"},
		{Binary: "armv7-gcc-rpath-dot", Expect: "fail"},
		{Binary: "armv7-gcc-rpath-dotdot", Expect: "fail"},
		{Binary: "armv7-gcc-rpath-relative", Expect: "fail"},
		{Binary: "armv7-gcc-rpath-parent-relative", Expect: "fail"},
		{Binary: "armv7-gcc-rpath-tmp", Expect: "fail"},
		{Binary: "armv7-gcc-rpath-var-tmp", Expect: "fail"},
		{Binary: "armv7-gcc-rpath-tmp-subdir", Expect: "fail"},
		{Binary: "armv7-gcc-rpath-empty-component", Expect: "fail"},
		{Binary: "armv7-gcc-rpath-mixed", Expect: "fail"},

		// armv7 Clang
		{Binary: "armv7-clang-no-rpath", Expect: "pass"},
		{Binary: "armv7-clang-rpath-absolute", Expect: "pass"},
		{Binary: "armv7-clang-rpath-dot", Expect: "fail"},
		{Binary: "armv7-clang-rpath-tmp", Expect: "fail"},
	})
}
