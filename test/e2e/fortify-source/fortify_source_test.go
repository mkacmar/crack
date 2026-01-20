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
		{"x86_64-gcc-no-fortify", "pass"},  // distro enables fortify by default - needs rebuild
		{"x86_64-gcc-fortify2-O0", "fail"}, // -O0 disables fortify optimization
		{"x86_64-gcc-fortify2-stripped", "pass"},
		{"x86_64-gcc-fortify2-static", "pass"},
		{"x86_64-gcc-fortify2-static-stripped", "skip"}, // no symbols to analyze
		{"x86_64-gcc-fortify2-lto", "pass"},
		{"x86_64-gcc-fortify2-simple", "skip"}, // no fortifiable functions

		// x86_64 Clang
		{"x86_64-clang-fortify2-O2", "pass"},
		{"x86_64-clang-fortify1-O1", "pass"},
		{"x86_64-clang-no-fortify", "fail"},  // no fortify enabled
		{"x86_64-clang-fortify2-O0", "fail"}, // -O0 disables fortify optimization
		{"x86_64-clang-fortify2-stripped", "pass"},
		{"x86_64-clang-fortify2-lto", "pass"},

		// aarch64 GCC
		{"aarch64-gcc-fortify2-O2", "pass"},
		{"aarch64-gcc-fortify1-O1", "pass"},
		{"aarch64-gcc-fortify3-O2", "pass"},
		{"aarch64-gcc-no-fortify", "pass"},  // distro enables fortify by default - needs rebuild
		{"aarch64-gcc-fortify2-O0", "fail"}, // -O0 disables fortify optimization
		{"aarch64-gcc-fortify2-stripped", "pass"},
		{"aarch64-gcc-fortify2-static", "pass"},
		{"aarch64-gcc-fortify2-static-stripped", "skip"}, // no symbols to analyze
		{"aarch64-gcc-fortify2-lto", "pass"},
		{"aarch64-gcc-fortify2-simple", "skip"}, // no fortifiable functions

		// aarch64 Clang
		{"aarch64-clang-fortify2-O2", "pass"},
		{"aarch64-clang-fortify1-O1", "pass"},
		{"aarch64-clang-no-fortify", "fail"},  // no fortify enabled
		{"aarch64-clang-fortify2-O0", "fail"}, // -O0 disables fortify optimization
		{"aarch64-clang-fortify2-stripped", "pass"},
		{"aarch64-clang-fortify2-lto", "pass"},

		// armv7 GCC - musl libc, fortify-source rule skips
		{"armv7-gcc-fortify2-O2", "skip"},
		{"armv7-gcc-fortify1-O1", "skip"},
		{"armv7-gcc-fortify3-O2", "skip"},
		{"armv7-gcc-no-fortify", "skip"},
		{"armv7-gcc-fortify2-O0", "skip"},
		{"armv7-gcc-fortify2-stripped", "skip"},
		{"armv7-gcc-fortify2-static", "fail"},          // static - can't detect musl, has unfortified funcs
		{"armv7-gcc-fortify2-static-stripped", "skip"}, // no symbols to analyze
		{"armv7-gcc-fortify2-lto", "skip"},
		{"armv7-gcc-fortify2-simple", "skip"}, // no fortifiable functions

		// armv7 Clang - musl libc, fortify-source rule skips
		{"armv7-clang-fortify2-O2", "skip"},
		{"armv7-clang-fortify1-O1", "skip"},
		{"armv7-clang-no-fortify", "skip"},
		{"armv7-clang-fortify2-O0", "skip"},
		{"armv7-clang-fortify2-stripped", "skip"},
		{"armv7-clang-fortify2-lto", "skip"},
	})
}
