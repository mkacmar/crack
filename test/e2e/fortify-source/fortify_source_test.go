package fortify_source_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestFortifySourceRule(t *testing.T) {
	e2e.RunRuleTests(t, "fortify-source", []e2e.TestCase{
		// x86_64 GCC
		{"x86_64-gcc-fortify2-O2", "pass"},
		{"x86_64-gcc-fortify1-O1", "pass"},
		{"x86_64-gcc-fortify3-O2", "pass"},
		{"x86_64-gcc-no-fortify", "fail"},
		{"x86_64-gcc-fortify2-O0", "fail"},
		{"x86_64-gcc-fortify2-stripped", "pass"},
		{"x86_64-gcc-fortify2-static", "pass"},
		{"x86_64-gcc-fortify2-static-stripped", "skip"}, // no symbols = no fortifiable functions detected
		{"x86_64-gcc-fortify2-lto", "pass"},
		{"x86_64-gcc-fortify2-simple", "skip"},

		// x86_64 Clang
		{"x86_64-clang-fortify2-O2", "pass"},
		{"x86_64-clang-fortify1-O1", "pass"},
		{"x86_64-clang-no-fortify", "fail"},
		{"x86_64-clang-fortify2-O0", "fail"},
		{"x86_64-clang-fortify2-stripped", "pass"},
		{"x86_64-clang-fortify2-lto", "pass"},

		// aarch64 GCC
		{"aarch64-gcc-fortify2-O2", "pass"},
		{"aarch64-gcc-fortify1-O1", "pass"},
		{"aarch64-gcc-fortify3-O2", "pass"},
		{"aarch64-gcc-no-fortify", "fail"},
		{"aarch64-gcc-fortify2-O0", "fail"},
		{"aarch64-gcc-fortify2-stripped", "pass"},
		{"aarch64-gcc-fortify2-static", "pass"},
		{"aarch64-gcc-fortify2-static-stripped", "skip"},
		{"aarch64-gcc-fortify2-lto", "pass"},
		{"aarch64-gcc-fortify2-simple", "skip"},

		// aarch64 Clang
		{"aarch64-clang-fortify2-O2", "pass"},
		{"aarch64-clang-fortify1-O1", "pass"},
		{"aarch64-clang-no-fortify", "fail"},
		{"aarch64-clang-fortify2-O0", "fail"},
		{"aarch64-clang-fortify2-stripped", "pass"},
		{"aarch64-clang-fortify2-lto", "pass"},

		// armv7 GCC
		{"armv7-gcc-fortify2-O2", "pass"},
		{"armv7-gcc-fortify1-O1", "pass"},
		{"armv7-gcc-fortify3-O2", "pass"},
		{"armv7-gcc-no-fortify", "fail"},
		{"armv7-gcc-fortify2-O0", "fail"},
		{"armv7-gcc-fortify2-stripped", "pass"},
		{"armv7-gcc-fortify2-static", "pass"},
		{"armv7-gcc-fortify2-static-stripped", "skip"},
		{"armv7-gcc-fortify2-lto", "pass"},
		{"armv7-gcc-fortify2-simple", "skip"},

		// armv7 Clang
		{"armv7-clang-fortify2-O2", "pass"},
		{"armv7-clang-fortify1-O1", "pass"},
		{"armv7-clang-no-fortify", "fail"},
		{"armv7-clang-fortify2-O0", "fail"},
		{"armv7-clang-fortify2-stripped", "pass"},
		{"armv7-clang-fortify2-lto", "pass"},
	})
}
