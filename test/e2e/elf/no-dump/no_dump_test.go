package no_dump_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestNoDumpRule(t *testing.T) {
	e2e.RunRuleTests(t, "no-dump", []e2e.TestCase{
		{Binary: "amd64-gcc-nodump", Expect: e2e.Pass},
		{Binary: "amd64-gcc-nodump-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-default", Expect: e2e.Fail},
		{Binary: "amd64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "amd64-clang-nodump", Expect: e2e.Pass},
		{Binary: "amd64-clang-nodump-stripped", Expect: e2e.Pass},
		{Binary: "amd64-clang-default", Expect: e2e.Fail},
		{Binary: "amd64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-gcc-nodump", Expect: e2e.Pass},
		{Binary: "arm64-gcc-nodump-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-default", Expect: e2e.Fail},
		{Binary: "arm64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-clang-nodump", Expect: e2e.Pass},
		{Binary: "arm64-clang-nodump-stripped", Expect: e2e.Pass},
		{Binary: "arm64-clang-default", Expect: e2e.Fail},
		{Binary: "arm64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-gcc-nodump", Expect: e2e.Pass},
		{Binary: "arm-gcc-nodump-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-default", Expect: e2e.Fail},
		{Binary: "arm-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-clang-nodump", Expect: e2e.Pass},
		{Binary: "arm-clang-nodump-stripped", Expect: e2e.Pass},
		{Binary: "arm-clang-default", Expect: e2e.Fail},
		{Binary: "arm-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "amd64-rustc-default", Expect: e2e.Skip},
		{Binary: "arm64-rustc-default", Expect: e2e.Skip},
		{Binary: "arm-rustc-default", Expect: e2e.Skip},
	})
}
