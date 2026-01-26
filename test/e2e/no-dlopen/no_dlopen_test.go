package nodlopen_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNoDLOpenRule(t *testing.T) {
	e2e.RunRuleTests(t, "no-dlopen", []e2e.TestCase{
		{Binary: "x86_64-gcc-nodlopen.so", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-default.so", Expect: e2e.Fail},

		{Binary: "x86_64-clang-nodlopen.so", Expect: e2e.Pass},
		{Binary: "x86_64-clang-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "x86_64-clang-default.so", Expect: e2e.Fail},

		{Binary: "aarch64-gcc-nodlopen.so", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-default.so", Expect: e2e.Fail},

		{Binary: "aarch64-clang-nodlopen.so", Expect: e2e.Pass},
		{Binary: "aarch64-clang-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "aarch64-clang-default.so", Expect: e2e.Fail},

		{Binary: "armv7-gcc-nodlopen.so", Expect: e2e.Pass},
		{Binary: "armv7-gcc-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "armv7-gcc-default.so", Expect: e2e.Fail},

		{Binary: "armv7-clang-nodlopen.so", Expect: e2e.Pass},
		{Binary: "armv7-clang-nodlopen-stripped.so", Expect: e2e.Pass},
		{Binary: "armv7-clang-default.so", Expect: e2e.Fail},
	})
}
