package cli

import (
	"fmt"
	"strings"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

func ParseTargetFilter(platforms, compilers string) (*rule.TargetFilter, error) {
	p, err := parseList(platforms, parsePlatformTarget)
	if err != nil {
		return nil, err
	}

	c, err := parseList(compilers, parseCompilerTarget)
	if err != nil {
		return nil, err
	}

	return &rule.TargetFilter{Platforms: p, Compilers: c}, nil
}

func parsePlatformTarget(s string) (rule.PlatformTarget, error) {
	name, version := splitTarget(s)

	arch, ok := binary.ParseArchitecture(name)
	if !ok {
		return rule.PlatformTarget{}, fmt.Errorf("unknown architecture %q, valid values: %s",
			name, strings.Join(validArchitectureNames(), ", "))
	}

	pt := rule.PlatformTarget{Architecture: arch}
	if version != "" {
		isa, err := binary.ParseISA(version)
		if err != nil {
			return rule.PlatformTarget{}, fmt.Errorf("invalid ISA version %q: %w", version, err)
		}
		pt.MaxISA = &isa
	}

	return pt, nil
}

func parseCompilerTarget(s string) (rule.CompilerTarget, error) {
	name, version := splitTarget(s)

	compiler, ok := parseCompiler(name)
	if !ok {
		return rule.CompilerTarget{}, fmt.Errorf("unknown compiler %q, valid values: %s",
			name, strings.Join(validCompilerNames(), ", "))
	}

	ct := rule.CompilerTarget{Compiler: compiler}
	if version != "" {
		v, err := toolchain.ParseVersion(version)
		if err != nil {
			return rule.CompilerTarget{}, fmt.Errorf("invalid compiler version %q: %w", version, err)
		}
		ct.MaxVersion = &v
	}

	return ct, nil
}

func parseList[T any](input string, parse func(string) (T, error)) ([]T, error) {
	if input == "" {
		return nil, nil
	}

	var results []T
	for _, item := range strings.Split(input, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		parsed, err := parse(item)
		if err != nil {
			return nil, err
		}
		results = append(results, parsed)
	}
	return results, nil
}

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
