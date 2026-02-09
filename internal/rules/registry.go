package rules

import (
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/rule/elf"
)

var registry = []rule.ELFRule{
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

func All() []rule.ELFRule {
	result := make([]rule.ELFRule, len(registry))
	copy(result, registry)
	return result
}

func Get(id string) rule.ELFRule {
	for _, r := range registry {
		if r.ID() == id {
			return r
		}
	}
	return nil
}
