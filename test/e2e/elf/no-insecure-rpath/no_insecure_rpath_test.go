package no_insecure_rpath_test

import (
	"testing"

	"go.kacmar.sk/crack/test/e2e"
)

func TestNoInsecureRPATHRule(t *testing.T) {
	const (
		insecurePrefix = "Insecure RPATH: "
		pathUnset      = "No RPATH set"
		pathSecure     = "RPATH secure"
	)

	e2e.RunRuleTests(t, "no-insecure-rpath", []e2e.TestCase{
		{Binary: "amd64-gcc-no-rpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "amd64-gcc-rpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "amd64-gcc-rpath-multiple-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "amd64-gcc-rpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "amd64-gcc-rpath-dotdot", Expect: e2e.Fail, Message: insecurePrefix + ".."},
		{Binary: "amd64-gcc-rpath-relative", Expect: e2e.Fail, Message: insecurePrefix + "./lib"},
		{Binary: "amd64-gcc-rpath-parent-relative", Expect: e2e.Fail, Message: insecurePrefix + "../lib"},
		{Binary: "amd64-gcc-rpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "amd64-gcc-rpath-var-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/var/tmp"},
		{Binary: "amd64-gcc-rpath-tmp-subdir", Expect: e2e.Fail, Message: insecurePrefix + "/tmp/mylibs"},
		{Binary: "amd64-gcc-rpath-empty-component", Expect: e2e.Fail, Message: insecurePrefix + "(empty)"},
		{Binary: "amd64-gcc-rpath-mixed", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "amd64-gcc-rpath-bare-relative", Expect: e2e.Fail, Message: insecurePrefix + "lib"},
		{Binary: "amd64-gcc-rpath-subdir-relative", Expect: e2e.Fail, Message: insecurePrefix + "subdir/lib"},
		{Binary: "amd64-gcc-rpath-origin-relative", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "amd64-gcc-rpath-origin-parent", Expect: e2e.Fail, Message: insecurePrefix + "$ORIGIN/.."},
		{Binary: "amd64-gcc-rpath-dev-shm", Expect: e2e.Fail, Message: insecurePrefix + "/dev/shm"},

		{Binary: "amd64-clang-no-rpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "amd64-clang-rpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "amd64-clang-rpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "amd64-clang-rpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "amd64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm64-gcc-no-rpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "arm64-gcc-rpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm64-gcc-rpath-multiple-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm64-gcc-rpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "arm64-gcc-rpath-dotdot", Expect: e2e.Fail, Message: insecurePrefix + ".."},
		{Binary: "arm64-gcc-rpath-relative", Expect: e2e.Fail, Message: insecurePrefix + "./lib"},
		{Binary: "arm64-gcc-rpath-parent-relative", Expect: e2e.Fail, Message: insecurePrefix + "../lib"},
		{Binary: "arm64-gcc-rpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "arm64-gcc-rpath-var-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/var/tmp"},
		{Binary: "arm64-gcc-rpath-tmp-subdir", Expect: e2e.Fail, Message: insecurePrefix + "/tmp/mylibs"},
		{Binary: "arm64-gcc-rpath-empty-component", Expect: e2e.Fail, Message: insecurePrefix + "(empty)"},
		{Binary: "arm64-gcc-rpath-mixed", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "arm64-gcc-rpath-bare-relative", Expect: e2e.Fail, Message: insecurePrefix + "lib"},
		{Binary: "arm64-gcc-rpath-subdir-relative", Expect: e2e.Fail, Message: insecurePrefix + "subdir/lib"},
		{Binary: "arm64-gcc-rpath-origin-relative", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm64-gcc-rpath-origin-parent", Expect: e2e.Fail, Message: insecurePrefix + "$ORIGIN/.."},
		{Binary: "arm64-gcc-rpath-dev-shm", Expect: e2e.Fail, Message: insecurePrefix + "/dev/shm"},

		{Binary: "arm64-clang-no-rpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "arm64-clang-rpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm64-clang-rpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "arm64-clang-rpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "arm64-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "arm-gcc-no-rpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "arm-gcc-rpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm-gcc-rpath-multiple-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm-gcc-rpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "arm-gcc-rpath-dotdot", Expect: e2e.Fail, Message: insecurePrefix + ".."},
		{Binary: "arm-gcc-rpath-relative", Expect: e2e.Fail, Message: insecurePrefix + "./lib"},
		{Binary: "arm-gcc-rpath-parent-relative", Expect: e2e.Fail, Message: insecurePrefix + "../lib"},
		{Binary: "arm-gcc-rpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "arm-gcc-rpath-var-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/var/tmp"},
		{Binary: "arm-gcc-rpath-tmp-subdir", Expect: e2e.Fail, Message: insecurePrefix + "/tmp/mylibs"},
		{Binary: "arm-gcc-rpath-empty-component", Expect: e2e.Fail, Message: insecurePrefix + "(empty)"},
		{Binary: "arm-gcc-rpath-mixed", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "arm-gcc-rpath-bare-relative", Expect: e2e.Fail, Message: insecurePrefix + "lib"},
		{Binary: "arm-gcc-rpath-subdir-relative", Expect: e2e.Fail, Message: insecurePrefix + "subdir/lib"},
		{Binary: "arm-gcc-rpath-origin-relative", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm-gcc-rpath-origin-parent", Expect: e2e.Fail, Message: insecurePrefix + "$ORIGIN/.."},
		{Binary: "arm-gcc-rpath-dev-shm", Expect: e2e.Fail, Message: insecurePrefix + "/dev/shm"},

		{Binary: "arm-clang-no-rpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "arm-clang-rpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "arm-clang-rpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "arm-clang-rpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "arm-clang-relocatable.o", Expect: e2e.Skip},

		{Binary: "riscv64-gcc-no-rpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "riscv64-gcc-rpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "riscv64-gcc-rpath-multiple-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "riscv64-gcc-rpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "riscv64-gcc-rpath-dotdot", Expect: e2e.Fail, Message: insecurePrefix + ".."},
		{Binary: "riscv64-gcc-rpath-relative", Expect: e2e.Fail, Message: insecurePrefix + "./lib"},
		{Binary: "riscv64-gcc-rpath-parent-relative", Expect: e2e.Fail, Message: insecurePrefix + "../lib"},
		{Binary: "riscv64-gcc-rpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "riscv64-gcc-rpath-var-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/var/tmp"},
		{Binary: "riscv64-gcc-rpath-tmp-subdir", Expect: e2e.Fail, Message: insecurePrefix + "/tmp/mylibs"},
		{Binary: "riscv64-gcc-rpath-empty-component", Expect: e2e.Fail, Message: insecurePrefix + "(empty)"},
		{Binary: "riscv64-gcc-rpath-mixed", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "riscv64-gcc-rpath-bare-relative", Expect: e2e.Fail, Message: insecurePrefix + "lib"},
		{Binary: "riscv64-gcc-rpath-subdir-relative", Expect: e2e.Fail, Message: insecurePrefix + "subdir/lib"},
		{Binary: "riscv64-gcc-rpath-origin-relative", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "riscv64-gcc-rpath-origin-parent", Expect: e2e.Fail, Message: insecurePrefix + "$ORIGIN/.."},
		{Binary: "riscv64-gcc-rpath-dev-shm", Expect: e2e.Fail, Message: insecurePrefix + "/dev/shm"},

		{Binary: "riscv64-clang-no-rpath", Expect: e2e.Pass, Message: pathUnset},
		{Binary: "riscv64-clang-rpath-absolute", Expect: e2e.Pass, Message: pathSecure},
		{Binary: "riscv64-clang-rpath-dot", Expect: e2e.Fail, Message: insecurePrefix + "."},
		{Binary: "riscv64-clang-rpath-tmp", Expect: e2e.Fail, Message: insecurePrefix + "/tmp"},
		{Binary: "riscv64-clang-relocatable.o", Expect: e2e.Skip},
	})
}
