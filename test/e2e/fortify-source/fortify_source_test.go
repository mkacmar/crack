package fortify_source_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestFortifySourceRule(t *testing.T) {
	e2e.RunRuleTests(t, "fortify-source", []e2e.TestCase{
		{Binary: "x86_64-gcc-fortify2-O2", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-fortify1-O1", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-fortify3-O2", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-no-fortify", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-fortify2-O0", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-fortify2-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-fortify2-static", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-fortify2-static-stripped", Expect: e2e.Skip},
		{Binary: "x86_64-gcc-fortify2-lto", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-fortify2-simple", Expect: e2e.Skip},

		{Binary: "x86_64-clang-fortify2-O2", Expect: e2e.Pass},
		{Binary: "x86_64-clang-fortify1-O1", Expect: e2e.Pass},
		{Binary: "x86_64-clang-no-fortify", Expect: e2e.Fail},
		{Binary: "x86_64-clang-fortify2-O0", Expect: e2e.Fail},
		{Binary: "x86_64-clang-fortify2-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-clang-fortify2-lto", Expect: e2e.Pass},

		{Binary: "aarch64-gcc-fortify2-O2", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-fortify1-O1", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-fortify3-O2", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-no-fortify", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-fortify2-O0", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-fortify2-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-fortify2-static", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-fortify2-static-stripped", Expect: e2e.Skip},
		{Binary: "aarch64-gcc-fortify2-lto", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-fortify2-simple", Expect: e2e.Skip},

		{Binary: "aarch64-clang-fortify2-O2", Expect: e2e.Pass},
		{Binary: "aarch64-clang-fortify1-O1", Expect: e2e.Pass},
		{Binary: "aarch64-clang-no-fortify", Expect: e2e.Fail},
		{Binary: "aarch64-clang-fortify2-O0", Expect: e2e.Fail},
		{Binary: "aarch64-clang-fortify2-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-clang-fortify2-lto", Expect: e2e.Pass},

		// armv7 uses musl libc, fortify-source rule skips
		{Binary: "armv7-gcc-fortify2-O2", Expect: e2e.Skip},
		{Binary: "armv7-gcc-fortify1-O1", Expect: e2e.Skip},
		{Binary: "armv7-gcc-fortify3-O2", Expect: e2e.Skip},
		{Binary: "armv7-gcc-no-fortify", Expect: e2e.Skip},
		{Binary: "armv7-gcc-fortify2-O0", Expect: e2e.Skip},
		{Binary: "armv7-gcc-fortify2-stripped", Expect: e2e.Skip},
		{Binary: "armv7-gcc-fortify2-static", Expect: e2e.Fail},
		{Binary: "armv7-gcc-fortify2-static-stripped", Expect: e2e.Skip},
		{Binary: "armv7-gcc-fortify2-lto", Expect: e2e.Skip},
		{Binary: "armv7-gcc-fortify2-simple", Expect: e2e.Skip},

		{Binary: "armv7-clang-fortify2-O2", Expect: e2e.Skip},
		{Binary: "armv7-clang-fortify1-O1", Expect: e2e.Skip},
		{Binary: "armv7-clang-no-fortify", Expect: e2e.Skip},
		{Binary: "armv7-clang-fortify2-O0", Expect: e2e.Skip},
		{Binary: "armv7-clang-fortify2-stripped", Expect: e2e.Skip},
		{Binary: "armv7-clang-fortify2-lto", Expect: e2e.Skip},

		{Binary: "x86_64-rustc-no-fortify", Expect: e2e.Skip},
		{Binary: "aarch64-rustc-no-fortify", Expect: e2e.Skip},
		{Binary: "armv7-rustc-no-fortify", Expect: e2e.Skip},
	})
}
