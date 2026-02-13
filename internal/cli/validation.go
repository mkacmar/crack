package cli

import (
	"strings"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/toolchain"
)

func splitTarget(s string) (name, version string) {
	if idx := strings.Index(s, ":"); idx != -1 {
		return s[:idx], s[idx+1:]
	}
	return s, ""
}

func parseCompiler(s string) (toolchain.Compiler, bool) {
	switch s {
	case toolchain.GCC.String():
		return toolchain.GCC, true
	case toolchain.Clang.String():
		return toolchain.Clang, true
	case toolchain.Rustc.String():
		return toolchain.Rustc, true
	default:
		return toolchain.Unknown, false
	}
}

func validCompilerNames() []string {
	return []string{toolchain.GCC.String(), toolchain.Clang.String(), toolchain.Rustc.String()}
}

func validArchitectureNames() []string {
	return []string{
		binary.ArchX86.String(),
		binary.ArchAMD64.String(),
		binary.ArchARM.String(),
		binary.ArchARM64.String(),
		binary.ArchRISCV.String(),
		binary.ArchPPC64.String(),
		binary.ArchMIPS.String(),
		binary.ArchS390X.String(),
	}
}
