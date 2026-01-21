package wxorx_test

import (
	"testing"

	"github.com/mkacmar/crack/test/e2e"
)

func TestWXorXRule(t *testing.T) {
	e2e.RunRuleTests(t, "wxorx", []e2e.TestCase{
		{Binary: "x86_64-gcc-wxorx", Expect: "pass"},
		{Binary: "x86_64-gcc-execstack", Expect: "fail"},
		{Binary: "x86_64-gcc-shared-wxorx", Expect: "pass"},
		{Binary: "x86_64-gcc-shared-execstack", Expect: "fail"},
		{Binary: "x86_64-clang-wxorx", Expect: "pass"},
		{Binary: "x86_64-clang-execstack", Expect: "fail"},
		{Binary: "x86_64-clang-shared-wxorx", Expect: "pass"},
		{Binary: "x86_64-clang-shared-execstack", Expect: "fail"},
		{Binary: "aarch64-gcc-wxorx", Expect: "pass"},
		{Binary: "aarch64-gcc-execstack", Expect: "fail"},
		{Binary: "aarch64-gcc-shared-wxorx", Expect: "pass"},
		{Binary: "aarch64-gcc-shared-execstack", Expect: "fail"},
		{Binary: "aarch64-clang-wxorx", Expect: "pass"},
		{Binary: "aarch64-clang-execstack", Expect: "fail"},
		{Binary: "aarch64-clang-shared-wxorx", Expect: "pass"},
		{Binary: "aarch64-clang-shared-execstack", Expect: "fail"},
		{Binary: "armv7-gcc-wxorx", Expect: "pass"},
		{Binary: "armv7-gcc-execstack", Expect: "fail"},
		{Binary: "armv7-gcc-shared-wxorx", Expect: "pass"},
		{Binary: "armv7-gcc-shared-execstack", Expect: "fail"},
		{Binary: "armv7-clang-wxorx", Expect: "pass"},
		{Binary: "armv7-clang-execstack", Expect: "fail"},
		{Binary: "armv7-clang-shared-wxorx", Expect: "pass"},
		{Binary: "armv7-clang-shared-execstack", Expect: "fail"},
	})
}
