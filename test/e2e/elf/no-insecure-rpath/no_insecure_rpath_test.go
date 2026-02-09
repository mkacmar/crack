package no_insecure_rpath_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNoInsecureRPATHRule(t *testing.T) {
	e2e.RunRuleTests(t, "no-insecure-rpath", []e2e.TestCase{
		{Binary: "x86_64-gcc-no-rpath", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-rpath-absolute", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-rpath-multiple-absolute", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-rpath-dot", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-rpath-dotdot", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-rpath-relative", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-rpath-parent-relative", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-rpath-tmp", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-rpath-var-tmp", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-rpath-tmp-subdir", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-rpath-empty-component", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-rpath-mixed", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-rpath-bare-relative", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-rpath-subdir-relative", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-rpath-origin-relative", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-rpath-dev-shm", Expect: e2e.Fail},

		{Binary: "x86_64-clang-no-rpath", Expect: e2e.Pass},
		{Binary: "x86_64-clang-rpath-absolute", Expect: e2e.Pass},
		{Binary: "x86_64-clang-rpath-dot", Expect: e2e.Fail},
		{Binary: "x86_64-clang-rpath-tmp", Expect: e2e.Fail},
		{Binary: "x86_64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "aarch64-gcc-no-rpath", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-rpath-absolute", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-rpath-multiple-absolute", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-rpath-dot", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-rpath-dotdot", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-rpath-relative", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-rpath-parent-relative", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-rpath-tmp", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-rpath-var-tmp", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-rpath-tmp-subdir", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-rpath-empty-component", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-rpath-mixed", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-rpath-bare-relative", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-rpath-subdir-relative", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-rpath-origin-relative", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-rpath-dev-shm", Expect: e2e.Fail},

		{Binary: "aarch64-clang-no-rpath", Expect: e2e.Pass},
		{Binary: "aarch64-clang-rpath-absolute", Expect: e2e.Pass},
		{Binary: "aarch64-clang-rpath-dot", Expect: e2e.Fail},
		{Binary: "aarch64-clang-rpath-tmp", Expect: e2e.Fail},
		{Binary: "aarch64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "armv7-gcc-no-rpath", Expect: e2e.Pass},
		{Binary: "armv7-gcc-rpath-absolute", Expect: e2e.Pass},
		{Binary: "armv7-gcc-rpath-multiple-absolute", Expect: e2e.Pass},
		{Binary: "armv7-gcc-rpath-dot", Expect: e2e.Fail},
		{Binary: "armv7-gcc-rpath-dotdot", Expect: e2e.Fail},
		{Binary: "armv7-gcc-rpath-relative", Expect: e2e.Fail},
		{Binary: "armv7-gcc-rpath-parent-relative", Expect: e2e.Fail},
		{Binary: "armv7-gcc-rpath-tmp", Expect: e2e.Fail},
		{Binary: "armv7-gcc-rpath-var-tmp", Expect: e2e.Fail},
		{Binary: "armv7-gcc-rpath-tmp-subdir", Expect: e2e.Fail},
		{Binary: "armv7-gcc-rpath-empty-component", Expect: e2e.Fail},
		{Binary: "armv7-gcc-rpath-mixed", Expect: e2e.Fail},
		{Binary: "armv7-gcc-rpath-bare-relative", Expect: e2e.Fail},
		{Binary: "armv7-gcc-rpath-subdir-relative", Expect: e2e.Fail},
		{Binary: "armv7-gcc-rpath-origin-relative", Expect: e2e.Fail},
		{Binary: "armv7-gcc-rpath-dev-shm", Expect: e2e.Fail},

		{Binary: "armv7-clang-no-rpath", Expect: e2e.Pass},
		{Binary: "armv7-clang-rpath-absolute", Expect: e2e.Pass},
		{Binary: "armv7-clang-rpath-dot", Expect: e2e.Fail},
		{Binary: "armv7-clang-rpath-tmp", Expect: e2e.Fail},
		{Binary: "armv7-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "x86_64-rustc-no-rpath", Expect: e2e.Pass},
		{Binary: "x86_64-rustc-rpath-absolute", Expect: e2e.Pass},
		{Binary: "x86_64-rustc-rpath-dot", Expect: e2e.Fail},
		{Binary: "x86_64-rustc-rpath-tmp", Expect: e2e.Fail},

		{Binary: "aarch64-rustc-no-rpath", Expect: e2e.Pass},
		{Binary: "aarch64-rustc-rpath-absolute", Expect: e2e.Pass},
		{Binary: "aarch64-rustc-rpath-dot", Expect: e2e.Fail},
		{Binary: "aarch64-rustc-rpath-tmp", Expect: e2e.Fail},

		{Binary: "armv7-rustc-no-rpath", Expect: e2e.Pass},
		{Binary: "armv7-rustc-rpath-absolute", Expect: e2e.Pass},
		{Binary: "armv7-rustc-rpath-dot", Expect: e2e.Fail},
		{Binary: "armv7-rustc-rpath-tmp", Expect: e2e.Fail},
	})
}
