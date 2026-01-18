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
	CFIICallRule{},
	CFIVCallRule{},
	KernelCFIRule{},
	ASANRule{},
	UBSanRule{},
	IntegerOverflowRule{},
	NoInsecureRPATHRule{},
	NoInsecureRUNPATHRule{},
	NoDLOpenRule{},
	NoPLTRule{},
	SeparateCodeRule{},
	NoDumpRule{},
	HiddenSymbolsRule{},
	StrippedRule{},
	GCSectionsRule{},
	IntelCETIBTRule{},
	IntelCETShadowStackRule{},
	RetpolineRule{},
	ARMBranchProtectionRule{},
	ARMPACRule{},
	ARMBTIRule{},
	ARMMTERule{},
	ARMShadowCallStackRule{},
}

func GetRuleByID(id string) model.Rule {
	for _, rule := range AllRules {
		if rule.ID() == id {
			return rule
		}
	}
	return nil
}
