package registry

import (
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/rule/elf"
)

var registry = []rule.Rule{
	elf.ARMBTIRule{},
	elf.ARMBranchProtectionRule{},
	elf.ARMMTERule{},
	elf.ARMPACRule{},
	elf.ASANRule{},
	elf.ASLRRule{},
	elf.CFIRule{},
	elf.FortifySourceRule{},
	elf.FullRELRORule{},
	elf.NXBitRule{},
	elf.NoDLOpenRule{},
	elf.NoDumpRule{},
	elf.NoInsecureRPATHRule{},
	elf.NoInsecureRUNPATHRule{},
	elf.PIERule{},
	elf.RELRORule{},
	elf.SafeStackRule{},
	elf.SeparateCodeRule{},
	elf.StackCanaryRule{},
	elf.StackLimitRule{},
	elf.StrippedRule{},
	elf.UBSanRule{},
	elf.X86CETIBTRule{},
	elf.X86CETShadowStackRule{},
	elf.X86RetpolineRule{},
}

// All returns all registered rules.
func All() []rule.Rule {
	result := make([]rule.Rule, len(registry))
	copy(result, registry)
	return result
}

// Where returns rules of type T matching the predicate.
// If predicate is nil, returns all rules of type T.
func Where[T rule.Rule](predicate func(rule.Rule) bool) []T {
	var result []T
	for _, r := range registry {
		if typed, ok := r.(T); ok {
			if predicate == nil || predicate(r) {
				result = append(result, typed)
			}
		}
	}
	return result
}

// Find returns the first rule of type T matching the predicate.
// If predicate is nil, returns the first rule of type T.
func Find[T rule.Rule](predicate func(rule.Rule) bool) (T, bool) {
	var zero T
	for _, r := range registry {
		if typed, ok := r.(T); ok {
			if predicate == nil || predicate(r) {
				return typed, true
			}
		}
	}
	return zero, false
}

// ByID returns a predicate matching rules by ID.
func ByID(id string) func(rule.Rule) bool {
	return func(r rule.Rule) bool { return r.ID() == id }
}
