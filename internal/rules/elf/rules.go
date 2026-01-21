package elf

import "github.com/mkacmar/crack/internal/model"

var AllRules = []model.Rule{
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

func GetRuleByID(id string) model.Rule {
	for _, rule := range AllRules {
		if rule.ID() == id {
			return rule
		}
	}
	return nil
}
