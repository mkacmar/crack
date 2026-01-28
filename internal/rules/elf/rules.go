package elf

import "github.com/mkacmar/crack/internal/rule"

var allRules = []rule.ELFRule{
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
	ASANRule{},
	UBSanRule{},
	NoInsecureRPATHRule{},
	NoInsecureRUNPATHRule{},
	NoDLOpenRule{},
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

func RegisterRules() {
	for _, r := range allRules {
		rule.Register(r)
	}
}
