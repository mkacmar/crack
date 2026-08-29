package no_insecure_runpath_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestNoInsecureRUNPATHRule(t *testing.T) {
	const (
		insecurePrefix = "Insecure RUNPATH: "
		pathUnset      = "No RUNPATH set"
		pathSecure     = "RUNPATH secure"
	)

	e2e.RunRuleTests(t, "no-insecure-runpath", []e2e.TestCase{
		{Binary: "amd64-gcc-no-runpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "amd64-gcc-runpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "amd64-gcc-runpath-multiple-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "amd64-gcc-runpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "amd64-gcc-runpath-dotdot", Expect: e2e.Fail, Message: insecurePrefix + ".."},
		{Binary: "amd64-gcc-runpath-relative", Expect: e2e.Fail, Message: insecurePrefix + "./lib"},
		{Binary: "amd64-gcc-runpath-parent-relative", Expect: e2e.Fail, Message: insecurePrefix + "../lib"},
		{Binary: "amd64-gcc-runpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "amd64-gcc-runpath-var-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/var/tmp"},
		{Binary: "amd64-gcc-runpath-tmp-subdir", Expect: e2e.Fail, Message: insecurePrefix + "/tmp/mylibs"},
		{Binary: "amd64-gcc-runpath-empty-component", Expect: e2e.Fail, Message: insecurePrefix + "(empty)"},
		{Binary: "amd64-gcc-runpath-mixed", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "amd64-gcc-runpath-bare-relative", Expect: e2e.Fail, Message: insecurePrefix + "lib"},
		{Binary: "amd64-gcc-runpath-subdir-relative", Expect: e2e.Fail, Message: insecurePrefix + "subdir/lib"},
		{Binary: "amd64-gcc-runpath-origin-relative", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "amd64-gcc-runpath-origin-parent", Expect: e2e.Fail, Message: insecurePrefix + "$ORIGIN/.."},
		{Binary: "amd64-gcc-runpath-dev-shm", Expect: e2e.Fail, Message: insecurePrefix + "/dev/shm"},

		{Binary: "amd64-clang-no-runpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "amd64-clang-runpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "amd64-clang-runpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "amd64-clang-runpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "amd64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-gcc-no-runpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "arm64-gcc-runpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm64-gcc-runpath-multiple-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm64-gcc-runpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "arm64-gcc-runpath-dotdot", Expect: e2e.Fail, Message: insecurePrefix + ".."},
		{Binary: "arm64-gcc-runpath-relative", Expect: e2e.Fail, Message: insecurePrefix + "./lib"},
		{Binary: "arm64-gcc-runpath-parent-relative", Expect: e2e.Fail, Message: insecurePrefix + "../lib"},
		{Binary: "arm64-gcc-runpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "arm64-gcc-runpath-var-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/var/tmp"},
		{Binary: "arm64-gcc-runpath-tmp-subdir", Expect: e2e.Fail, Message: insecurePrefix + "/tmp/mylibs"},
		{Binary: "arm64-gcc-runpath-empty-component", Expect: e2e.Fail, Message: insecurePrefix + "(empty)"},
		{Binary: "arm64-gcc-runpath-mixed", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "arm64-gcc-runpath-bare-relative", Expect: e2e.Fail, Message: insecurePrefix + "lib"},
		{Binary: "arm64-gcc-runpath-subdir-relative", Expect: e2e.Fail, Message: insecurePrefix + "subdir/lib"},
		{Binary: "arm64-gcc-runpath-origin-relative", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm64-gcc-runpath-origin-parent", Expect: e2e.Fail, Message: insecurePrefix + "$ORIGIN/.."},
		{Binary: "arm64-gcc-runpath-dev-shm", Expect: e2e.Fail, Message: insecurePrefix + "/dev/shm"},

		{Binary: "arm64-clang-no-runpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "arm64-clang-runpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm64-clang-runpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "arm64-clang-runpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "arm64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-gcc-no-runpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "arm-gcc-runpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm-gcc-runpath-multiple-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm-gcc-runpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "arm-gcc-runpath-dotdot", Expect: e2e.Fail, Message: insecurePrefix + ".."},
		{Binary: "arm-gcc-runpath-relative", Expect: e2e.Fail, Message: insecurePrefix + "./lib"},
		{Binary: "arm-gcc-runpath-parent-relative", Expect: e2e.Fail, Message: insecurePrefix + "../lib"},
		{Binary: "arm-gcc-runpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "arm-gcc-runpath-var-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/var/tmp"},
		{Binary: "arm-gcc-runpath-tmp-subdir", Expect: e2e.Fail, Message: insecurePrefix + "/tmp/mylibs"},
		{Binary: "arm-gcc-runpath-empty-component", Expect: e2e.Fail, Message: insecurePrefix + "(empty)"},
		{Binary: "arm-gcc-runpath-mixed", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "arm-gcc-runpath-bare-relative", Expect: e2e.Fail, Message: insecurePrefix + "lib"},
		{Binary: "arm-gcc-runpath-subdir-relative", Expect: e2e.Fail, Message: insecurePrefix + "subdir/lib"},
		{Binary: "arm-gcc-runpath-origin-relative", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm-gcc-runpath-origin-parent", Expect: e2e.Fail, Message: insecurePrefix + "$ORIGIN/.."},
		{Binary: "arm-gcc-runpath-dev-shm", Expect: e2e.Fail, Message: insecurePrefix + "/dev/shm"},

		{Binary: "arm-clang-no-runpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "arm-clang-runpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm-clang-runpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "arm-clang-runpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "arm-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "riscv64-gcc-no-runpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "riscv64-gcc-runpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "riscv64-gcc-runpath-multiple-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "riscv64-gcc-runpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "riscv64-gcc-runpath-dotdot", Expect: e2e.Fail, Message: insecurePrefix + ".."},
		{Binary: "riscv64-gcc-runpath-relative", Expect: e2e.Fail, Message: insecurePrefix + "./lib"},
		{Binary: "riscv64-gcc-runpath-parent-relative", Expect: e2e.Fail, Message: insecurePrefix + "../lib"},
		{Binary: "riscv64-gcc-runpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "riscv64-gcc-runpath-var-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/var/tmp"},
		{Binary: "riscv64-gcc-runpath-tmp-subdir", Expect: e2e.Fail, Message: insecurePrefix + "/tmp/mylibs"},
		{Binary: "riscv64-gcc-runpath-empty-component", Expect: e2e.Fail, Message: insecurePrefix + "(empty)"},
		{Binary: "riscv64-gcc-runpath-mixed", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "riscv64-gcc-runpath-bare-relative", Expect: e2e.Fail, Message: insecurePrefix + "lib"},
		{Binary: "riscv64-gcc-runpath-subdir-relative", Expect: e2e.Fail, Message: insecurePrefix + "subdir/lib"},
		{Binary: "riscv64-gcc-runpath-origin-relative", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "riscv64-gcc-runpath-origin-parent", Expect: e2e.Fail, Message: insecurePrefix + "$ORIGIN/.."},
		{Binary: "riscv64-gcc-runpath-dev-shm", Expect: e2e.Fail, Message: insecurePrefix + "/dev/shm"},

		{Binary: "riscv64-clang-no-runpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "riscv64-clang-runpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "riscv64-clang-runpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "riscv64-clang-runpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "riscv64-clang-relocatable.o", Expect: e2e.Skip},
	})
}
