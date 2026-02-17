package no_insecure_runpath_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNoInsecureRUNPATHRule(t *testing.T) {
	e2e.RunRuleTests(t, "no-insecure-runpath", []e2e.TestCase{
		{Binary: "amd64-gcc-no-runpath", Expect: e2e.Pass},
		{Binary: "amd64-gcc-runpath-absolute", Expect: e2e.Pass},
		{Binary: "amd64-gcc-runpath-multiple-absolute", Expect: e2e.Pass},
		{Binary: "amd64-gcc-runpath-dot", Expect: e2e.Fail},
		{Binary: "amd64-gcc-runpath-dotdot", Expect: e2e.Fail},
		{Binary: "amd64-gcc-runpath-relative", Expect: e2e.Fail},
		{Binary: "amd64-gcc-runpath-parent-relative", Expect: e2e.Fail},
		{Binary: "amd64-gcc-runpath-tmp", Expect: e2e.Fail},
		{Binary: "amd64-gcc-runpath-var-tmp", Expect: e2e.Fail},
		{Binary: "amd64-gcc-runpath-tmp-subdir", Expect: e2e.Fail},
		{Binary: "amd64-gcc-runpath-empty-component", Expect: e2e.Fail},
		{Binary: "amd64-gcc-runpath-mixed", Expect: e2e.Fail},
		{Binary: "amd64-gcc-runpath-bare-relative", Expect: e2e.Fail},
		{Binary: "amd64-gcc-runpath-subdir-relative", Expect: e2e.Fail},
		{Binary: "amd64-gcc-runpath-origin-relative", Expect: e2e.Pass},
		{Binary: "amd64-gcc-runpath-origin-parent", Expect: e2e.Fail},
		{Binary: "amd64-gcc-runpath-dev-shm", Expect: e2e.Fail},

		{Binary: "amd64-clang-no-runpath", Expect: e2e.Pass},
		{Binary: "amd64-clang-runpath-absolute", Expect: e2e.Pass},
		{Binary: "amd64-clang-runpath-dot", Expect: e2e.Fail},
		{Binary: "amd64-clang-runpath-tmp", Expect: e2e.Fail},
		{Binary: "amd64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-gcc-no-runpath", Expect: e2e.Pass},
		{Binary: "arm64-gcc-runpath-absolute", Expect: e2e.Pass},
		{Binary: "arm64-gcc-runpath-multiple-absolute", Expect: e2e.Pass},
		{Binary: "arm64-gcc-runpath-dot", Expect: e2e.Fail},
		{Binary: "arm64-gcc-runpath-dotdot", Expect: e2e.Fail},
		{Binary: "arm64-gcc-runpath-relative", Expect: e2e.Fail},
		{Binary: "arm64-gcc-runpath-parent-relative", Expect: e2e.Fail},
		{Binary: "arm64-gcc-runpath-tmp", Expect: e2e.Fail},
		{Binary: "arm64-gcc-runpath-var-tmp", Expect: e2e.Fail},
		{Binary: "arm64-gcc-runpath-tmp-subdir", Expect: e2e.Fail},
		{Binary: "arm64-gcc-runpath-empty-component", Expect: e2e.Fail},
		{Binary: "arm64-gcc-runpath-mixed", Expect: e2e.Fail},
		{Binary: "arm64-gcc-runpath-bare-relative", Expect: e2e.Fail},
		{Binary: "arm64-gcc-runpath-subdir-relative", Expect: e2e.Fail},
		{Binary: "arm64-gcc-runpath-origin-relative", Expect: e2e.Pass},
		{Binary: "arm64-gcc-runpath-origin-parent", Expect: e2e.Fail},
		{Binary: "arm64-gcc-runpath-dev-shm", Expect: e2e.Fail},

		{Binary: "arm64-clang-no-runpath", Expect: e2e.Pass},
		{Binary: "arm64-clang-runpath-absolute", Expect: e2e.Pass},
		{Binary: "arm64-clang-runpath-dot", Expect: e2e.Fail},
		{Binary: "arm64-clang-runpath-tmp", Expect: e2e.Fail},
		{Binary: "arm64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-gcc-no-runpath", Expect: e2e.Pass},
		{Binary: "arm-gcc-runpath-absolute", Expect: e2e.Pass},
		{Binary: "arm-gcc-runpath-multiple-absolute", Expect: e2e.Pass},
		{Binary: "arm-gcc-runpath-dot", Expect: e2e.Fail},
		{Binary: "arm-gcc-runpath-dotdot", Expect: e2e.Fail},
		{Binary: "arm-gcc-runpath-relative", Expect: e2e.Fail},
		{Binary: "arm-gcc-runpath-parent-relative", Expect: e2e.Fail},
		{Binary: "arm-gcc-runpath-tmp", Expect: e2e.Fail},
		{Binary: "arm-gcc-runpath-var-tmp", Expect: e2e.Fail},
		{Binary: "arm-gcc-runpath-tmp-subdir", Expect: e2e.Fail},
		{Binary: "arm-gcc-runpath-empty-component", Expect: e2e.Fail},
		{Binary: "arm-gcc-runpath-mixed", Expect: e2e.Fail},
		{Binary: "arm-gcc-runpath-bare-relative", Expect: e2e.Fail},
		{Binary: "arm-gcc-runpath-subdir-relative", Expect: e2e.Fail},
		{Binary: "arm-gcc-runpath-origin-relative", Expect: e2e.Pass},
		{Binary: "arm-gcc-runpath-origin-parent", Expect: e2e.Fail},
		{Binary: "arm-gcc-runpath-dev-shm", Expect: e2e.Fail},

		{Binary: "arm-clang-no-runpath", Expect: e2e.Pass},
		{Binary: "arm-clang-runpath-absolute", Expect: e2e.Pass},
		{Binary: "arm-clang-runpath-dot", Expect: e2e.Fail},
		{Binary: "arm-clang-runpath-tmp", Expect: e2e.Fail},
		{Binary: "arm-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "amd64-rustc-no-runpath", Expect: e2e.Pass},
		{Binary: "amd64-rustc-runpath-absolute", Expect: e2e.Pass},
		{Binary: "amd64-rustc-runpath-dot", Expect: e2e.Fail},
		{Binary: "amd64-rustc-runpath-tmp", Expect: e2e.Fail},

		{Binary: "arm64-rustc-no-runpath", Expect: e2e.Pass},
		{Binary: "arm64-rustc-runpath-absolute", Expect: e2e.Pass},
		{Binary: "arm64-rustc-runpath-dot", Expect: e2e.Fail},
		{Binary: "arm64-rustc-runpath-tmp", Expect: e2e.Fail},

		{Binary: "arm-rustc-no-runpath", Expect: e2e.Pass},
		{Binary: "arm-rustc-runpath-absolute", Expect: e2e.Pass},
		{Binary: "arm-rustc-runpath-dot", Expect: e2e.Fail},
		{Binary: "arm-rustc-runpath-tmp", Expect: e2e.Fail},
	})
}
