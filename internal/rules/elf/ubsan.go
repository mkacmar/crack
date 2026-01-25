package elf

import (
	"fmt"
	"strings"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

const UBSanRuleID = "ubsan"

var ubsanHandlers = []string{
	"__ubsan_handle_add_overflow",
	"__ubsan_handle_alignment_assumption",
	"__ubsan_handle_builtin_unreachable",
	"__ubsan_handle_divrem_overflow",
	"__ubsan_handle_dynamic_type_cache_miss",
	"__ubsan_handle_float_cast_overflow",
	"__ubsan_handle_function_type_mismatch",
	"__ubsan_handle_implicit_conversion",
	"__ubsan_handle_invalid_builtin",
	"__ubsan_handle_load_invalid_value",
	"__ubsan_handle_missing_return",
	"__ubsan_handle_mul_overflow",
	"__ubsan_handle_negate_overflow",
	"__ubsan_handle_nonnull_arg",
	"__ubsan_handle_nonnull_return",
	"__ubsan_handle_nullability_arg",
	"__ubsan_handle_nullability_return",
	"__ubsan_handle_out_of_bounds",
	"__ubsan_handle_pointer_overflow",
	"__ubsan_handle_shift_out_of_bounds",
	"__ubsan_handle_sub_overflow",
	"__ubsan_handle_type_mismatch",
	"__ubsan_handle_vla_bound_not_positive",
}

// UBSanRule checks for Undefined Behavior Sanitizer
// Clang: https://clang.llvm.org/docs/UndefinedBehaviorSanitizer.html
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-fsanitize=undefined
type UBSanRule struct{}

func (r UBSanRule) ID() string   { return UBSanRuleID }
func (r UBSanRule) Name() string { return "Undefined Behavior Sanitizer" }

func (r UBSanRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 4, Minor: 9}, Flag: "-fsanitize=undefined"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 3, Minor: 3}, Flag: "-fsanitize=undefined"},
		},
	}
}

func (r UBSanRule) Execute(bin *binary.ELFBinary) rule.ExecuteResult {
	allSymbols := make(map[string]struct{})
	for _, sym := range bin.Symbols {
		allSymbols[sym.Name] = struct{}{}
	}
	for _, sym := range bin.DynSymbols {
		allSymbols[sym.Name] = struct{}{}
	}

	var foundHandlers []string
	for _, handler := range ubsanHandlers {
		for symName := range allSymbols {
			if strings.Contains(symName, handler) {
				foundHandlers = append(foundHandlers, handler)
				break
			}
		}
	}

	if len(foundHandlers) > 0 {
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: fmt.Sprintf("UBSan is enabled, found: %v", foundHandlers),
		}
	}
	return rule.ExecuteResult{
		Status:  rule.StatusFailed,
		Message: "UBSan is NOT enabled",
	}
}
