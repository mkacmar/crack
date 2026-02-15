package no_dlopen_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNoDLOpenRule(t *testing.T) {
	e2e.RunRuleTests(t, "no-dlopen", []e2e.TestCase{
		{Binary: "amd64-gcc-nodlopen.so", Expect: e2e.Pass},
		{Binary: "amd64-gcc-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "amd64-gcc-default.so", Expect: e2e.Fail},
		{Binary: "amd64-gcc-pie-executable", Expect: e2e.Skip},

		{Binary: "amd64-clang-nodlopen.so", Expect: e2e.Pass},
		{Binary: "amd64-clang-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "amd64-clang-default.so", Expect: e2e.Fail},
		{Binary: "amd64-clang-pie-executable", Expect: e2e.Skip},

		{Binary: "arm64-gcc-nodlopen.so", Expect: e2e.Pass},
		{Binary: "arm64-gcc-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "arm64-gcc-default.so", Expect: e2e.Fail},
		{Binary: "arm64-gcc-pie-executable", Expect: e2e.Skip},

		{Binary: "arm64-clang-nodlopen.so", Expect: e2e.Pass},
		{Binary: "arm64-clang-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "arm64-clang-default.so", Expect: e2e.Fail},
		{Binary: "arm64-clang-pie-executable", Expect: e2e.Skip},

		{Binary: "arm-gcc-nodlopen.so", Expect: e2e.Pass},
		{Binary: "arm-gcc-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "arm-gcc-default.so", Expect: e2e.Fail},
		{Binary: "arm-gcc-pie-executable", Expect: e2e.Skip},

		{Binary: "arm-clang-nodlopen.so", Expect: e2e.Pass},
		{Binary: "arm-clang-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "arm-clang-default.so", Expect: e2e.Fail},
		{Binary: "arm-clang-pie-executable", Expect: e2e.Skip},

		{Binary: "amd64-rustc-executable", Expect: e2e.Skip},
		{Binary: "amd64-rustc-nodlopen.so", Expect: e2e.Pass},
		{Binary: "amd64-rustc-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "amd64-rustc-default.so", Expect: e2e.Fail},

		{Binary: "arm64-rustc-executable", Expect: e2e.Skip},
		{Binary: "arm64-rustc-nodlopen.so", Expect: e2e.Pass},
		{Binary: "arm64-rustc-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "arm64-rustc-default.so", Expect: e2e.Fail},

		{Binary: "arm-rustc-executable", Expect: e2e.Skip},
		{Binary: "arm-rustc-nodlopen.so", Expect: e2e.Pass},
		{Binary: "arm-rustc-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "arm-rustc-default.so", Expect: e2e.Fail},
	})
}
