package fortify_source_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestFortifySourceRule(t *testing.T) {
	e2e.RunRuleTests(t, "fortify-source", []e2e.TestCase{
		// x86_64 GCC
		{Binary: "x86_64-gcc-fortify2-O2", Expect: "pass"},
		{Binary: "x86_64-gcc-fortify1-O1", Expect: "pass"},
		{Binary: "x86_64-gcc-fortify3-O2", Expect: "pass"},
		{Binary: "x86_64-gcc-no-fortify", Expect: "fail"},
		{Binary: "x86_64-gcc-fortify2-O0", Expect: "fail"}, // -O0 disables fortify optimization
		{Binary: "x86_64-gcc-fortify2-stripped", Expect: "pass"},
		{Binary: "x86_64-gcc-fortify2-static", Expect: "pass"},
		{Binary: "x86_64-gcc-fortify2-static-stripped", Expect: "skip"}, // no symbols to analyze
		{Binary: "x86_64-gcc-fortify2-lto", Expect: "pass"},
		{Binary: "x86_64-gcc-fortify2-simple", Expect: "skip"}, // no fortifiable functions

		// x86_64 Clang
		{Binary: "x86_64-clang-fortify2-O2", Expect: "pass"},
		{Binary: "x86_64-clang-fortify1-O1", Expect: "pass"},
		{Binary: "x86_64-clang-no-fortify", Expect: "fail"},  // no fortify enabled
		{Binary: "x86_64-clang-fortify2-O0", Expect: "fail"}, // -O0 disables fortify optimization
		{Binary: "x86_64-clang-fortify2-stripped", Expect: "pass"},
		{Binary: "x86_64-clang-fortify2-lto", Expect: "pass"},

		// aarch64 GCC
		{Binary: "aarch64-gcc-fortify2-O2", Expect: "pass"},
		{Binary: "aarch64-gcc-fortify1-O1", Expect: "pass"},
		{Binary: "aarch64-gcc-fortify3-O2", Expect: "pass"},
		{Binary: "aarch64-gcc-no-fortify", Expect: "fail"},
		{Binary: "aarch64-gcc-fortify2-O0", Expect: "fail"}, // -O0 disables fortify optimization
		{Binary: "aarch64-gcc-fortify2-stripped", Expect: "pass"},
		{Binary: "aarch64-gcc-fortify2-static", Expect: "pass"},
		{Binary: "aarch64-gcc-fortify2-static-stripped", Expect: "skip"}, // no symbols to analyze
		{Binary: "aarch64-gcc-fortify2-lto", Expect: "pass"},
		{Binary: "aarch64-gcc-fortify2-simple", Expect: "skip"}, // no fortifiable functions

		// aarch64 Clang
		{Binary: "aarch64-clang-fortify2-O2", Expect: "pass"},
		{Binary: "aarch64-clang-fortify1-O1", Expect: "pass"},
		{Binary: "aarch64-clang-no-fortify", Expect: "fail"},  // no fortify enabled
		{Binary: "aarch64-clang-fortify2-O0", Expect: "fail"}, // -O0 disables fortify optimization
		{Binary: "aarch64-clang-fortify2-stripped", Expect: "pass"},
		{Binary: "aarch64-clang-fortify2-lto", Expect: "pass"},

		// armv7 GCC - musl libc, fortify-source rule skips
		{Binary: "armv7-gcc-fortify2-O2", Expect: "skip"},
		{Binary: "armv7-gcc-fortify1-O1", Expect: "skip"},
		{Binary: "armv7-gcc-fortify3-O2", Expect: "skip"},
		{Binary: "armv7-gcc-no-fortify", Expect: "skip"},
		{Binary: "armv7-gcc-fortify2-O0", Expect: "skip"},
		{Binary: "armv7-gcc-fortify2-stripped", Expect: "skip"},
		{Binary: "armv7-gcc-fortify2-static", Expect: "fail"},          // static - can't detect musl, has unfortified funcs
		{Binary: "armv7-gcc-fortify2-static-stripped", Expect: "skip"}, // no symbols to analyze
		{Binary: "armv7-gcc-fortify2-lto", Expect: "skip"},
		{Binary: "armv7-gcc-fortify2-simple", Expect: "skip"}, // no fortifiable functions

		// armv7 Clang - musl libc, fortify-source rule skips
		{Binary: "armv7-clang-fortify2-O2", Expect: "skip"},
		{Binary: "armv7-clang-fortify1-O1", Expect: "skip"},
		{Binary: "armv7-clang-no-fortify", Expect: "skip"},
		{Binary: "armv7-clang-fortify2-O0", Expect: "skip"},
		{Binary: "armv7-clang-fortify2-stripped", Expect: "skip"},
		{Binary: "armv7-clang-fortify2-lto", Expect: "skip"},
	})
}
