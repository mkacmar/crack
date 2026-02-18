package registry

import (
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/rule/elf"
)

var registry = []rule.Rule{
	elf.PIERule{},
	elf.NXBitRule{},
	elf.RELRORule{},
	elf.FullRELRORule{},
	elf.ASLRRule{},
	elf.StackCanaryRule{},
	elf.SafeStackRule{},
	elf.StackLimitRule{},
	elf.FortifySourceRule{},
	elf.CFIRule{},
	elf.ASANRule{},
	elf.UBSanRule{},
	elf.NoInsecureRPATHRule{},
	elf.NoInsecureRUNPATHRule{},
	elf.NoDLOpenRule{},
	elf.SeparateCodeRule{},
	elf.NoDumpRule{},
	elf.StrippedRule{},
	elf.X86CETIBTRule{},
	elf.X86CETShadowStackRule{},
	elf.X86RetpolineRule{},
	elf.ARMBranchProtectionRule{},
	elf.ARMPACRule{},
	elf.ARMBTIRule{},
	elf.ARMMTERule{},
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
