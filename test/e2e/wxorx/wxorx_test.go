package wxorx_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestWXorXRule(t *testing.T) {
	e2e.RunRuleTests(t, "wxorx", []e2e.TestCase{
		{Binary: "x86_64-gcc-wxorx", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-execstack", Expect: e2e.Fail},
		{Binary: "x86_64-gcc-shared-wxorx", Expect: e2e.Pass},
		{Binary: "x86_64-gcc-shared-execstack", Expect: e2e.Fail},
		{Binary: "x86_64-clang-wxorx", Expect: e2e.Pass},
		{Binary: "x86_64-clang-execstack", Expect: e2e.Fail},
		{Binary: "x86_64-clang-shared-wxorx", Expect: e2e.Pass},
		{Binary: "x86_64-clang-shared-execstack", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-wxorx", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-execstack", Expect: e2e.Fail},
		{Binary: "aarch64-gcc-shared-wxorx", Expect: e2e.Pass},
		{Binary: "aarch64-gcc-shared-execstack", Expect: e2e.Fail},
		{Binary: "aarch64-clang-wxorx", Expect: e2e.Pass},
		{Binary: "aarch64-clang-execstack", Expect: e2e.Fail},
		{Binary: "aarch64-clang-shared-wxorx", Expect: e2e.Pass},
		{Binary: "aarch64-clang-shared-execstack", Expect: e2e.Fail},
		{Binary: "armv7-gcc-wxorx", Expect: e2e.Pass},
		{Binary: "armv7-gcc-execstack", Expect: e2e.Fail},
		{Binary: "armv7-gcc-shared-wxorx", Expect: e2e.Pass},
		{Binary: "armv7-gcc-shared-execstack", Expect: e2e.Fail},
		{Binary: "armv7-clang-wxorx", Expect: e2e.Pass},
		{Binary: "armv7-clang-execstack", Expect: e2e.Fail},
		{Binary: "armv7-clang-shared-wxorx", Expect: e2e.Pass},
		{Binary: "armv7-clang-shared-execstack", Expect: e2e.Fail},
	})
}
