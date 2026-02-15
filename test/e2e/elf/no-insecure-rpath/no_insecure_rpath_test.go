package no_insecure_rpath_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNoInsecureRPATHRule(t *testing.T) {
	e2e.RunRuleTests(t, "no-insecure-rpath", []e2e.TestCase{
		{Binary: "amd64-gcc-no-rpath", Expect: e2e.Pass},
		{Binary: "amd64-gcc-rpath-absolute", Expect: e2e.Pass},
		{Binary: "amd64-gcc-rpath-multiple-absolute", Expect: e2e.Pass},
		{Binary: "amd64-gcc-rpath-dot", Expect: e2e.Fail},
		{Binary: "amd64-gcc-rpath-dotdot", Expect: e2e.Fail},
		{Binary: "amd64-gcc-rpath-relative", Expect: e2e.Fail},
		{Binary: "amd64-gcc-rpath-parent-relative", Expect: e2e.Fail},
		{Binary: "amd64-gcc-rpath-tmp", Expect: e2e.Fail},
		{Binary: "amd64-gcc-rpath-var-tmp", Expect: e2e.Fail},
		{Binary: "amd64-gcc-rpath-tmp-subdir", Expect: e2e.Fail},
		{Binary: "amd64-gcc-rpath-empty-component", Expect: e2e.Fail},
		{Binary: "amd64-gcc-rpath-mixed", Expect: e2e.Fail},
		{Binary: "amd64-gcc-rpath-bare-relative", Expect: e2e.Fail},
		{Binary: "amd64-gcc-rpath-subdir-relative", Expect: e2e.Fail},
		{Binary: "amd64-gcc-rpath-origin-relative", Expect: e2e.Fail},
		{Binary: "amd64-gcc-rpath-dev-shm", Expect: e2e.Fail},

		{Binary: "amd64-clang-no-rpath", Expect: e2e.Pass},
		{Binary: "amd64-clang-rpath-absolute", Expect: e2e.Pass},
		{Binary: "amd64-clang-rpath-dot", Expect: e2e.Fail},
		{Binary: "amd64-clang-rpath-tmp", Expect: e2e.Fail},
		{Binary: "amd64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-gcc-no-rpath", Expect: e2e.Pass},
		{Binary: "arm64-gcc-rpath-absolute", Expect: e2e.Pass},
		{Binary: "arm64-gcc-rpath-multiple-absolute", Expect: e2e.Pass},
		{Binary: "arm64-gcc-rpath-dot", Expect: e2e.Fail},
		{Binary: "arm64-gcc-rpath-dotdot", Expect: e2e.Fail},
		{Binary: "arm64-gcc-rpath-relative", Expect: e2e.Fail},
		{Binary: "arm64-gcc-rpath-parent-relative", Expect: e2e.Fail},
		{Binary: "arm64-gcc-rpath-tmp", Expect: e2e.Fail},
		{Binary: "arm64-gcc-rpath-var-tmp", Expect: e2e.Fail},
		{Binary: "arm64-gcc-rpath-tmp-subdir", Expect: e2e.Fail},
		{Binary: "arm64-gcc-rpath-empty-component", Expect: e2e.Fail},
		{Binary: "arm64-gcc-rpath-mixed", Expect: e2e.Fail},
		{Binary: "arm64-gcc-rpath-bare-relative", Expect: e2e.Fail},
		{Binary: "arm64-gcc-rpath-subdir-relative", Expect: e2e.Fail},
		{Binary: "arm64-gcc-rpath-origin-relative", Expect: e2e.Fail},
		{Binary: "arm64-gcc-rpath-dev-shm", Expect: e2e.Fail},

		{Binary: "arm64-clang-no-rpath", Expect: e2e.Pass},
		{Binary: "arm64-clang-rpath-absolute", Expect: e2e.Pass},
		{Binary: "arm64-clang-rpath-dot", Expect: e2e.Fail},
		{Binary: "arm64-clang-rpath-tmp", Expect: e2e.Fail},
		{Binary: "arm64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-gcc-no-rpath", Expect: e2e.Pass},
		{Binary: "arm-gcc-rpath-absolute", Expect: e2e.Pass},
		{Binary: "arm-gcc-rpath-multiple-absolute", Expect: e2e.Pass},
		{Binary: "arm-gcc-rpath-dot", Expect: e2e.Fail},
		{Binary: "arm-gcc-rpath-dotdot", Expect: e2e.Fail},
		{Binary: "arm-gcc-rpath-relative", Expect: e2e.Fail},
		{Binary: "arm-gcc-rpath-parent-relative", Expect: e2e.Fail},
		{Binary: "arm-gcc-rpath-tmp", Expect: e2e.Fail},
		{Binary: "arm-gcc-rpath-var-tmp", Expect: e2e.Fail},
		{Binary: "arm-gcc-rpath-tmp-subdir", Expect: e2e.Fail},
		{Binary: "arm-gcc-rpath-empty-component", Expect: e2e.Fail},
		{Binary: "arm-gcc-rpath-mixed", Expect: e2e.Fail},
		{Binary: "arm-gcc-rpath-bare-relative", Expect: e2e.Fail},
		{Binary: "arm-gcc-rpath-subdir-relative", Expect: e2e.Fail},
		{Binary: "arm-gcc-rpath-origin-relative", Expect: e2e.Fail},
		{Binary: "arm-gcc-rpath-dev-shm", Expect: e2e.Fail},

		{Binary: "arm-clang-no-rpath", Expect: e2e.Pass},
		{Binary: "arm-clang-rpath-absolute", Expect: e2e.Pass},
		{Binary: "arm-clang-rpath-dot", Expect: e2e.Fail},
		{Binary: "arm-clang-rpath-tmp", Expect: e2e.Fail},
		{Binary: "arm-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "amd64-rustc-no-rpath", Expect: e2e.Pass},
		{Binary: "amd64-rustc-rpath-absolute", Expect: e2e.Pass},
		{Binary: "amd64-rustc-rpath-dot", Expect: e2e.Fail},
		{Binary: "amd64-rustc-rpath-tmp", Expect: e2e.Fail},

		{Binary: "arm64-rustc-no-rpath", Expect: e2e.Pass},
		{Binary: "arm64-rustc-rpath-absolute", Expect: e2e.Pass},
		{Binary: "arm64-rustc-rpath-dot", Expect: e2e.Fail},
		{Binary: "arm64-rustc-rpath-tmp", Expect: e2e.Fail},

		{Binary: "arm-rustc-no-rpath", Expect: e2e.Pass},
		{Binary: "arm-rustc-rpath-absolute", Expect: e2e.Pass},
		{Binary: "arm-rustc-rpath-dot", Expect: e2e.Fail},
		{Binary: "arm-rustc-rpath-tmp", Expect: e2e.Fail},
	})
}
