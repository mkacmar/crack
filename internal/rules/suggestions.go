package rules

import (
	"fmt"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

func buildSuggestion(tc model.Toolchain, feature model.FeatureAvailability) string {
	if tc.Compiler == model.CompilerUnknown {
		return buildGenericSuggestion(feature)
	}
	return buildCompilerSuggestion(tc, feature)
}

func buildGenericSuggestion(feature model.FeatureAvailability) string {
	var parts []string
	parts = append(parts, "Toolchain not detected (binary likely stripped), use")

	var options []string
	if gccReq := feature.GetRequirement(model.CompilerGCC); gccReq != nil && gccReq.Flag != "" {
		options = append(options, fmt.Sprintf("GCC %s+ with \"%s\"", gccReq.MinVersion.String(), gccReq.Flag))
	}
	if clangReq := feature.GetRequirement(model.CompilerClang); clangReq != nil && clangReq.Flag != "" {
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

func buildCompilerSuggestion(tc model.Toolchain, feature model.FeatureAvailability) string {
	req := feature.GetRequirement(tc.Compiler)
	if req == nil {
		other := model.CompilerGCC
		if tc.Compiler == model.CompilerGCC {
			other = model.CompilerClang
		}
		otherReq := feature.GetRequirement(other)
		if otherReq != nil {
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

	if req.DefaultVersion != (model.Version{}) && !tc.Version.IsAtLeast(req.DefaultVersion) {
		return fmt.Sprintf("Use \"%s\" (default in %s %s+).",
			flag, compilerName, req.DefaultVersion.String())
	}

	if req.DefaultVersion == (model.Version{}) {
		return fmt.Sprintf("Use \"%s\".", flag)
	}

	return fmt.Sprintf("Should be enabled by default in %s %s+. Check build configuration or use \"%s\".",
		compilerName, req.DefaultVersion.String(), flag)
}
