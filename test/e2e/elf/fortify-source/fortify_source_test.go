package fortify_source_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestFortifySourceRule(t *testing.T) {
	e2e.RunRuleTests(t, "fortify-source", []e2e.TestCase{
		{Binary: "amd64-gcc-fortify2-O2", Expect: e2e.Pass},
		{Binary: "amd64-gcc-fortify1-O1", Expect: e2e.Pass},
		{Binary: "amd64-gcc-fortify3-O2", Expect: e2e.Pass},
		{Binary: "amd64-gcc-no-fortify", Expect: e2e.Fail},
		{Binary: "amd64-gcc-fortify2-O0", Expect: e2e.Fail},
		{Binary: "amd64-gcc-fortify2-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-fortify2-static", Expect: e2e.Pass},
		{Binary: "amd64-gcc-fortify2-static-stripped", Expect: e2e.Skip},
		{Binary: "amd64-gcc-fortify2-lto", Expect: e2e.Pass},
		{Binary: "amd64-gcc-fortify2-simple", Expect: e2e.Skip},

		{Binary: "amd64-clang-fortify2-O2", Expect: e2e.Pass},
		{Binary: "amd64-clang-fortify1-O1", Expect: e2e.Pass},
		{Binary: "amd64-clang-no-fortify", Expect: e2e.Fail},
		{Binary: "amd64-clang-fortify2-O0", Expect: e2e.Fail},
		{Binary: "amd64-clang-fortify2-stripped", Expect: e2e.Pass},
		{Binary: "amd64-clang-fortify2-lto", Expect: e2e.Pass},

		{Binary: "amd64-gcc-fortify2-shared", Expect: e2e.Pass},
		{Binary: "amd64-gcc-no-fortify-shared", Expect: e2e.Fail},

		{Binary: "arm64-gcc-fortify2-O2", Expect: e2e.Pass},
		{Binary: "arm64-gcc-fortify1-O1", Expect: e2e.Pass},
		{Binary: "arm64-gcc-fortify3-O2", Expect: e2e.Pass},
		{Binary: "arm64-gcc-no-fortify", Expect: e2e.Fail},
		{Binary: "arm64-gcc-fortify2-O0", Expect: e2e.Fail},
		{Binary: "arm64-gcc-fortify2-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-fortify2-static", Expect: e2e.Pass},
		{Binary: "arm64-gcc-fortify2-static-stripped", Expect: e2e.Skip},
		{Binary: "arm64-gcc-fortify2-lto", Expect: e2e.Pass},
		{Binary: "arm64-gcc-fortify2-simple", Expect: e2e.Skip},

		{Binary: "arm64-clang-fortify2-O2", Expect: e2e.Pass},
		{Binary: "arm64-clang-fortify1-O1", Expect: e2e.Pass},
		{Binary: "arm64-clang-no-fortify", Expect: e2e.Fail},
		{Binary: "arm64-clang-fortify2-O0", Expect: e2e.Fail},
		{Binary: "arm64-clang-fortify2-stripped", Expect: e2e.Pass},
		{Binary: "arm64-clang-fortify2-lto", Expect: e2e.Pass},

		{Binary: "arm64-gcc-fortify2-shared", Expect: e2e.Pass},
		{Binary: "arm64-gcc-no-fortify-shared", Expect: e2e.Fail},

		// arm uses musl libc, fortify-source rule skips
		{Binary: "arm-gcc-fortify2-O2", Expect: e2e.Skip},
		{Binary: "arm-gcc-fortify1-O1", Expect: e2e.Skip},
		{Binary: "arm-gcc-fortify3-O2", Expect: e2e.Skip},
		{Binary: "arm-gcc-no-fortify", Expect: e2e.Skip},
		{Binary: "arm-gcc-fortify2-O0", Expect: e2e.Skip},
		{Binary: "arm-gcc-fortify2-stripped", Expect: e2e.Skip},
		{Binary: "arm-gcc-fortify2-static", Expect: e2e.Fail},
		{Binary: "arm-gcc-fortify2-static-stripped", Expect: e2e.Skip},
		{Binary: "arm-gcc-fortify2-lto", Expect: e2e.Skip},
		{Binary: "arm-gcc-fortify2-simple", Expect: e2e.Skip},

		{Binary: "arm-clang-fortify2-O2", Expect: e2e.Skip},
		{Binary: "arm-clang-fortify1-O1", Expect: e2e.Skip},
		{Binary: "arm-clang-no-fortify", Expect: e2e.Skip},
		{Binary: "arm-clang-fortify2-O0", Expect: e2e.Skip},
		{Binary: "arm-clang-fortify2-stripped", Expect: e2e.Skip},
		{Binary: "arm-clang-fortify2-lto", Expect: e2e.Skip},

		// musl shared libs: no PT_INTERP, musl detected via DT_NEEDED
		{Binary: "arm-gcc-fortify2-shared", Expect: e2e.Skip},
		{Binary: "arm-gcc-no-fortify-shared", Expect: e2e.Skip},

		{Binary: "amd64-rustc-no-fortify", Expect: e2e.Skip},
		{Binary: "arm64-rustc-no-fortify", Expect: e2e.Skip},
		{Binary: "arm-rustc-no-fortify", Expect: e2e.Skip},
	})
}
