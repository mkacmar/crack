package aslr_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestASLRRule(t *testing.T) {
	const (
		notPIE    = "Not ASLR compatible, not PIE"
		execStack = "Not ASLR compatible, executable stack"
		textRel   = "Not ASLR compatible, text relocations present"
		notExec   = "Not an executable or shared library"
		sharedLib = "Shared library, ASLR not applicable"
	)

	e2e.RunRuleTests(t, "aslr", []e2e.TestCase{
		{Binary: "amd64-gcc-aslr-full", Expect: e2e.Pass},
		{Binary: "amd64-gcc-no-pie", Expect: e2e.Fail, Message: notPIE},
		{Binary: "amd64-gcc-execstack", Expect: e2e.Fail, Message: execStack},
		{Binary: "amd64-gcc-shared", Expect: e2e.Skip, Message: sharedLib},
		{Binary: "amd64-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "amd64-gcc-aslr-stripped", Expect: e2e.Pass},
		{Binary: "amd64-gcc-static-no-pie", Expect: e2e.Fail, Message: notPIE},
		{Binary: "amd64-gcc-textrel-patched", Expect: e2e.Fail, Message: textRel},
		{Binary: "amd64-gcc-relocatable.o", Expect: e2e.Skip, Message: notExec},

		{Binary: "amd64-clang-aslr-full", Expect: e2e.Pass},
		{Binary: "amd64-clang-no-pie", Expect: e2e.Fail, Message: notPIE},
		{Binary: "amd64-clang-execstack", Expect: e2e.Fail, Message: execStack},
		{Binary: "amd64-clang-static-no-pie", Expect: e2e.Fail, Message: notPIE},

		{Binary: "amd64-gcc-old-pie", Expect: e2e.Pass},

		{Binary: "arm64-gcc-aslr-full", Expect: e2e.Pass},
		{Binary: "arm64-gcc-no-pie", Expect: e2e.Fail, Message: notPIE},
		{Binary: "arm64-gcc-execstack", Expect: e2e.Fail, Message: execStack},
		{Binary: "arm64-gcc-shared", Expect: e2e.Skip, Message: sharedLib},
		{Binary: "arm64-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "arm64-gcc-aslr-stripped", Expect: e2e.Pass},
		{Binary: "arm64-gcc-static-no-pie", Expect: e2e.Fail, Message: notPIE},
		{Binary: "arm64-gcc-textrel-patched", Expect: e2e.Fail, Message: textRel},
		{Binary: "arm64-gcc-relocatable.o", Expect: e2e.Skip, Message: notExec},

		{Binary: "arm64-clang-aslr-full", Expect: e2e.Pass},
		{Binary: "arm64-clang-no-pie", Expect: e2e.Fail, Message: notPIE},
		{Binary: "arm64-clang-execstack", Expect: e2e.Fail, Message: execStack},
		{Binary: "arm64-clang-static-no-pie", Expect: e2e.Fail, Message: notPIE},

		{Binary: "arm-gcc-aslr-full", Expect: e2e.Pass},
		{Binary: "arm-gcc-no-pie", Expect: e2e.Fail, Message: notPIE},
		{Binary: "arm-gcc-execstack", Expect: e2e.Fail, Message: execStack},
		{Binary: "arm-gcc-shared", Expect: e2e.Skip, Message: sharedLib},
		{Binary: "arm-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "arm-gcc-aslr-stripped", Expect: e2e.Pass},
		{Binary: "arm-gcc-static-no-pie", Expect: e2e.Fail, Message: notPIE},
		{Binary: "arm-gcc-textrel-patched", Expect: e2e.Fail, Message: textRel},
		{Binary: "arm-gcc-relocatable.o", Expect: e2e.Skip, Message: notExec},

		{Binary: "arm-clang-aslr-full", Expect: e2e.Pass},
		{Binary: "arm-clang-no-pie", Expect: e2e.Fail, Message: notPIE},
		{Binary: "arm-clang-execstack", Expect: e2e.Fail, Message: execStack},
		{Binary: "arm-clang-static-no-pie", Expect: e2e.Fail, Message: notPIE},

		{Binary: "riscv64-gcc-aslr-full", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-no-pie", Expect: e2e.Fail, Message: notPIE},
		{Binary: "riscv64-gcc-execstack", Expect: e2e.Fail, Message: execStack},
		{Binary: "riscv64-gcc-shared", Expect: e2e.Skip, Message: sharedLib},
		{Binary: "riscv64-gcc-static-pie", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-aslr-stripped", Expect: e2e.Pass},
		{Binary: "riscv64-gcc-static-no-pie", Expect: e2e.Fail, Message: notPIE},
		{Binary: "riscv64-gcc-textrel-patched", Expect: e2e.Fail, Message: textRel},
		{Binary: "riscv64-gcc-relocatable.o", Expect: e2e.Skip, Message: notExec},

		{Binary: "riscv64-clang-aslr-full", Expect: e2e.Pass},
		{Binary: "riscv64-clang-no-pie", Expect: e2e.Fail, Message: notPIE},
		{Binary: "riscv64-clang-execstack", Expect: e2e.Fail, Message: execStack},
		{Binary: "riscv64-clang-static-no-pie", Expect: e2e.Fail, Message: notPIE},
	})
}
