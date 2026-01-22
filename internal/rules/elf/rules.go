package elf

import "github.com/mkacmar/crack/internal/model"

var allRules = []model.ELFRule{
	PIERule{},
	NXBitRule{},
	RELRORule{},
	FullRELRORule{},
	ASLRRule{},
	WXorXRule{},
	StackCanaryRule{},
	SafeStackRule{},
	StackLimitRule{},
	FortifySourceRule{},
	CFIRule{},
	KernelCFIRule{},
	ASANRule{},
	UBSanRule{},
	NoInsecureRPATHRule{},
	NoInsecureRUNPATHRule{},
	NoDLOpenRule{},
	NoPLTRule{},
	SeparateCodeRule{},
	NoDumpRule{},
	StrippedRule{},
	X86CETIBTRule{},
	X86CETShadowStackRule{},
	X86RetpolineRule{},
	ARMBranchProtectionRule{},
	ARMPACRule{},
	ARMBTIRule{},
	ARMMTERule{},
}

var AllRules = func() map[string]model.ELFRule {
	m := make(map[string]model.ELFRule, len(allRules))
	for _, r := range allRules {
		m[r.ID()] = r
	}
	return m
}()

func GetRuleByID(id string) model.ELFRule {
	return AllRules[id]
}
