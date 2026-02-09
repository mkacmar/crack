package rules

import (
	"fmt"
	"strings"

	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

func buildSuggestion(build toolchain.BuildInfo, applicability rule.Applicability) string {
	if build.Compiler == toolchain.Unknown {
		return buildGenericSuggestion(applicability)
	}
	return buildCompilerSuggestion(build, applicability)
}

func buildGenericSuggestion(applicability rule.Applicability) string {
	var parts []string
	parts = append(parts, "Toolchain not detected (binary likely stripped), use")

	var options []string
	if gccReq, ok := getCompilerRequirement(applicability.Compilers, toolchain.GCC); ok && gccReq.Flag != "" {
		options = append(options, fmt.Sprintf("GCC %s+ with \"%s\"", gccReq.MinVersion.String(), gccReq.Flag))
	}
	if clangReq, ok := getCompilerRequirement(applicability.Compilers, toolchain.Clang); ok && clangReq.Flag != "" {
		options = append(options, fmt.Sprintf("Clang %s+ with \"%s\"", clangReq.MinVersion.String(), clangReq.Flag))
	}

	if len(options) > 0 {
		parts = append(parts, strings.Join(options, " or "))
	}

	result := strings.Join(parts, " ")
	if !strings.HasSuffix(result, ".") {
		result += "."
	}
	return result
}

func getCompilerRequirement(compilers map[toolchain.Compiler]rule.CompilerRequirement, target toolchain.Compiler) (rule.CompilerRequirement, bool) {
	req, ok := compilers[target]
	return req, ok
}

func buildCompilerSuggestion(build toolchain.BuildInfo, applicability rule.Applicability) string {
	req, ok := getCompilerRequirement(applicability.Compilers, build.Compiler)
	if !ok {
		other := toolchain.GCC
		if build.Compiler == toolchain.GCC {
			other = toolchain.Clang
		}
		if otherReq, ok := getCompilerRequirement(applicability.Compilers, other); ok {
			return fmt.Sprintf("Feature requires %s %s+. Consider switching or use alternatives.",
				other.String(), otherReq.MinVersion.String())
		}
		return "Feature not supported by detected compilers."
	}

	flag := req.Flag
	compilerName := build.Compiler.String()

	if !build.Version.IsAtLeast(req.MinVersion) {
		return fmt.Sprintf("Requires %s %s+ (you have %s %s), update and use \"%s\".",
			compilerName, req.MinVersion.String(), compilerName, build.Version.String(), flag)
	}

	if req.DefaultVersion != (toolchain.Version{}) && !build.Version.IsAtLeast(req.DefaultVersion) {
		return fmt.Sprintf("Use \"%s\" (default in %s %s+).",
			flag, compilerName, req.DefaultVersion.String())
	}

	if req.DefaultVersion == (toolchain.Version{}) {
		return fmt.Sprintf("Use \"%s\".", flag)
	}

	return fmt.Sprintf("Should be enabled by default in %s %s+. Check build configuration or use \"%s\".",
		compilerName, req.DefaultVersion.String(), flag)
}
