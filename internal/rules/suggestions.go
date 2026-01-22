package rules

import (
	"fmt"
	"strings"

	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

func buildSuggestion(tc toolchain.Toolchain, applicability rule.Applicability) string {
	if tc.Compiler == toolchain.CompilerUnknown {
		return buildGenericSuggestion(applicability)
	}
	return buildCompilerSuggestion(tc, applicability)
}

func buildGenericSuggestion(applicability rule.Applicability) string {
	var parts []string
	parts = append(parts, "Toolchain not detected (binary likely stripped), use")

	var options []string
	if gccReq, ok := applicability.Compilers[toolchain.CompilerGCC]; ok && gccReq.Flag != "" {
		options = append(options, fmt.Sprintf("GCC %s+ with \"%s\"", gccReq.MinVersion.String(), gccReq.Flag))
	}
	if clangReq, ok := applicability.Compilers[toolchain.CompilerClang]; ok && clangReq.Flag != "" {
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

func buildCompilerSuggestion(tc toolchain.Toolchain, applicability rule.Applicability) string {
	req, ok := applicability.Compilers[tc.Compiler]
	if !ok {
		other := toolchain.CompilerGCC
		if tc.Compiler == toolchain.CompilerGCC {
			other = toolchain.CompilerClang
		}
		if otherReq, ok := applicability.Compilers[other]; ok {
			return fmt.Sprintf("Feature requires %s %s+. Consider switching or use alternatives.",
				other.String(), otherReq.MinVersion.String())
		}
		return "Feature not supported by detected compilers."
	}

	flag := req.Flag
	compilerName := tc.Compiler.String()

	if !tc.Version.IsAtLeast(req.MinVersion) {
		return fmt.Sprintf("Requires %s %s+ (you have %s %s), update and use \"%s\".",
			compilerName, req.MinVersion.String(), compilerName, tc.Version.String(), flag)
	}

	if req.DefaultVersion != (toolchain.Version{}) && !tc.Version.IsAtLeast(req.DefaultVersion) {
		return fmt.Sprintf("Use \"%s\" (default in %s %s+).",
			flag, compilerName, req.DefaultVersion.String())
	}

	if req.DefaultVersion == (toolchain.Version{}) {
		return fmt.Sprintf("Use \"%s\".", flag)
	}

	return fmt.Sprintf("Should be enabled by default in %s %s+. Check build configuration or use \"%s\".",
		compilerName, req.DefaultVersion.String(), flag)
}
