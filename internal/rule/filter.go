package rule

import (
	"fmt"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/toolchain"
)

type PlatformTarget struct {
	Architecture binary.Architecture
	MaxISA       *binary.ISA
}

type CompilerTarget struct {
	Compiler   toolchain.Compiler
	MaxVersion *toolchain.Version
}

type TargetFilter struct {
	Platforms []PlatformTarget
	Compilers []CompilerTarget
}

func (f *TargetFilter) IsEmpty() bool {
	return len(f.Platforms) == 0 && len(f.Compilers) == 0
}

func (f *TargetFilter) Matches(app Applicability) bool {
	if f.IsEmpty() {
		return true
	}

	for _, pt := range f.Platforms {
		if !app.Platform.Architecture.Matches(pt.Architecture) {
			return false
		}
		if pt.MaxISA != nil && app.Platform.MinISA.Major > 0 {
			if !pt.MaxISA.IsAtLeast(app.Platform.MinISA) {
				return false
			}
		}
	}

	for _, ct := range f.Compilers {
		req, exists := app.Compilers[ct.Compiler]
		if !exists {
			return false
		}
		if ct.MaxVersion != nil && req.MinVersion.Major > 0 {
			if !ct.MaxVersion.IsAtLeast(req.MinVersion) {
				return false
			}
		}
	}

	return true
}

func ParsePlatformTarget(s string) (PlatformTarget, error) {
	name, version := splitTarget(s)

	arch, ok := binary.ParseArchitecture(name)
	if !ok {
		return PlatformTarget{}, fmt.Errorf("unknown architecture %q, valid values: %s",
			name, strings.Join(binary.ValidArchitectureNames(), ", "))
	}

	pt := PlatformTarget{Architecture: arch}
	if version != "" {
		isa, err := binary.ParseISA(version)
		if err != nil {
			return PlatformTarget{}, fmt.Errorf("invalid ISA version %q: %w", version, err)
		}
		pt.MaxISA = &isa
	}

	return pt, nil
}

func ParseCompilerTarget(s string) (CompilerTarget, error) {
	name, version := splitTarget(s)

	compiler, ok := toolchain.ParseCompiler(name)
	if !ok {
		return CompilerTarget{}, fmt.Errorf("unknown compiler %q, valid values: %s",
			name, strings.Join(toolchain.ValidCompilerNames(), ", "))
	}

	ct := CompilerTarget{Compiler: compiler}
	if version != "" {
		v, err := toolchain.ParseVersion(version)
		if err != nil {
			return CompilerTarget{}, fmt.Errorf("invalid compiler version %q: %w", version, err)
		}
		ct.MaxVersion = &v
	}

	return ct, nil
}

func splitTarget(s string) (name, version string) {
	if idx := strings.Index(s, ":"); idx != -1 {
		return s[:idx], s[idx+1:]
	}
	return s, ""
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

func ParseTargetFilter(platforms, compilers string) (*TargetFilter, error) {
	p, err := parseList(platforms, ParsePlatformTarget)
	if err != nil {
		return nil, err
	}

	c, err := parseList(compilers, ParseCompilerTarget)
	if err != nil {
		return nil, err
	}

	return &TargetFilter{Platforms: p, Compilers: c}, nil
}

func FilterRules(ruleIDs []string, filter *TargetFilter) []string {
	if filter == nil || filter.IsEmpty() {
		return ruleIDs
	}

	var filtered []string
	for _, id := range ruleIDs {
		r := Get(id)
		if r == nil {
			continue
		}
		if filter.Matches(r.Applicability()) {
			filtered = append(filtered, id)
		}
	}
	return filtered
}
