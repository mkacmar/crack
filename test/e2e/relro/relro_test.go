package relro_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestRELRORule(t *testing.T) {
	e2e.RunRuleTests(t, "relro", []e2e.TestCase{
		{Binary: "x86_64-gcc-partial-relro", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-full-relro", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-no-relro", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-full-relro-static", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-full-relro-shared", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "x86_64-clang-partial-relro", Expect: e2e.Pass},
		{Binary: "x86_64-clang-full-relro", Expect: e2e.Pass},
		{Binary: "x86_64-clang-no-relro", Expect: e2e.Fail},
		{Binary: "x86_64-clang-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "aarch64-gcc-partial-relro", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-full-relro", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-no-relro", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-full-relro-static", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-full-relro-shared", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "aarch64-clang-partial-relro", Expect: e2e.Pass},
		{Binary: "aarch64-clang-full-relro", Expect: e2e.Pass},
		{Binary: "aarch64-clang-no-relro", Expect: e2e.Fail},
		{Binary: "aarch64-clang-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "armv7-gcc-partial-relro", Expect: e2e.Pass},
		{Binary: "armv7-gcc-full-relro", Expect: e2e.Pass},
		{Binary: "armv7-gcc-no-relro", Expect: e2e.Fail},
		{Binary: "armv7-gcc-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "armv7-gcc-full-relro-static", Expect: e2e.Pass},
		{Binary: "armv7-gcc-full-relro-shared", Expect: e2e.Pass},
		{Binary: "armv7-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "armv7-clang-partial-relro", Expect: e2e.Pass},
		{Binary: "armv7-clang-full-relro", Expect: e2e.Pass},
		{Binary: "armv7-clang-no-relro", Expect: e2e.Fail},
		{Binary: "armv7-clang-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "armv7-clang-relocatable.o", Expect: e2e.Skip},
	})
}

func TestFullRELRORule(t *testing.T) {
	e2e.RunRuleTests(t, "full-relro", []e2e.TestCase{
		{Binary: "x86_64-gcc-partial-relro", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-full-relro", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-no-relro", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-full-relro-static", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-full-relro-shared", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "x86_64-clang-partial-relro", Expect: e2e.Fail},
		{Binary: "x86_64-clang-full-relro", Expect: e2e.Pass},
		{Binary: "x86_64-clang-no-relro", Expect: e2e.Fail},
		{Binary: "x86_64-clang-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "x86_64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "aarch64-gcc-partial-relro", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-full-relro", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-no-relro", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-full-relro-static", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-full-relro-shared", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "aarch64-clang-partial-relro", Expect: e2e.Fail},
		{Binary: "aarch64-clang-full-relro", Expect: e2e.Pass},
		{Binary: "aarch64-clang-no-relro", Expect: e2e.Fail},
		{Binary: "aarch64-clang-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "aarch64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "armv7-gcc-partial-relro", Expect: e2e.Fail},
		{Binary: "armv7-gcc-full-relro", Expect: e2e.Pass},
		{Binary: "armv7-gcc-no-relro", Expect: e2e.Fail},
		{Binary: "armv7-gcc-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "armv7-gcc-full-relro-static", Expect: e2e.Pass},
		{Binary: "armv7-gcc-full-relro-shared", Expect: e2e.Pass},
		{Binary: "armv7-gcc-relocatable.o", Expect: e2e.Skip},

		{Binary: "armv7-clang-partial-relro", Expect: e2e.Fail},
		{Binary: "armv7-clang-full-relro", Expect: e2e.Pass},
		{Binary: "armv7-clang-no-relro", Expect: e2e.Fail},
		{Binary: "armv7-clang-full-relro-stripped", Expect: e2e.Pass},
		{Binary: "armv7-clang-relocatable.o", Expect: e2e.Skip},
	})
}
