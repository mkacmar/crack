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
	HiddenSymbolsRule{},
	StrippedRule{},
	IntelCETIBTRule{},
	IntelCETShadowStackRule{},
	RetpolineRule{},
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
