package profile

import "github.com/mkacmar/crack/internal/rules/elf"

type Profile struct {
	Rules []string
}

var Minimal = Profile{
	Rules: []string{
		elf.PIERule{}.ID(),
		elf.NXBitRule{}.ID(),
		elf.RELRORule{}.ID(),
		elf.StackCanaryRule{}.ID(),
		elf.FortifySourceRule{}.ID(),
		elf.ASLRRule{}.ID(),
	},
}

var Recommended = Profile{
	Rules: []string{
		elf.PIERule{}.ID(),
		elf.NXBitRule{}.ID(),
		elf.RELRORule{}.ID(),
		elf.FullRELRORule{}.ID(),
		elf.StackCanaryRule{}.ID(),
		elf.FortifySourceRule{}.ID(),
		elf.FortifyLevel2Rule{}.ID(),
		elf.ASLRRule{}.ID(),
		elf.WXorXRule{}.ID(),
		elf.NoInsecureRPATHRule{}.ID(),
		elf.NoInsecureRUNPATHRule{}.ID(),
		elf.FormatStringRule{}.ID(),
		elf.HiddenSymbolsRule{}.ID(),
		elf.SeparateCodeRule{}.ID(),
		elf.NoDLOpenRule{}.ID(),
		elf.GCSectionsRule{}.ID(),
	},
}

var Hardened = Profile{
	Rules: []string{
		elf.PIERule{}.ID(),
		elf.NXBitRule{}.ID(),
		elf.RELRORule{}.ID(),
		elf.FullRELRORule{}.ID(),
		elf.StackCanaryRule{}.ID(),
		elf.FortifySourceRule{}.ID(),
		elf.FortifyLevel2Rule{}.ID(),
		elf.ASLRRule{}.ID(),
		elf.WXorXRule{}.ID(),
		elf.NoInsecureRPATHRule{}.ID(),
		elf.NoInsecureRUNPATHRule{}.ID(),
		elf.FormatStringRule{}.ID(),
		elf.HiddenSymbolsRule{}.ID(),
		elf.SeparateCodeRule{}.ID(),
		elf.NoDLOpenRule{}.ID(),
		elf.NoPLTRule{}.ID(),
		elf.NoDumpRule{}.ID(),
		elf.GCSectionsRule{}.ID(),
		elf.StrippedRule{}.ID(),
		elf.StackLimitRule{}.ID(),
		elf.CFIRule{}.ID(),
		elf.CFIICallRule{}.ID(),
		elf.CFIVCallRule{}.ID(),
		elf.SafeStackRule{}.ID(),
		elf.IntegerOverflowRule{}.ID(),
		elf.IntelCETIBTRule{}.ID(),
		elf.IntelCETShadowStackRule{}.ID(),
		elf.RetpolineRule{}.ID(),
		elf.ARMBranchProtectionRule{}.ID(),
		elf.ARMPACRule{}.ID(),
		elf.ARMBTIRule{}.ID(),
		elf.ARMMTERule{}.ID(),
		elf.ARMShadowCallStackRule{}.ID(),
	},
}

var BuildTest = Profile{
	Rules: []string{
		elf.ASANRule{}.ID(),
		elf.UBSanRule{}.ID(),
		elf.SafeStackRule{}.ID(),
		elf.IntegerOverflowRule{}.ID(),
		elf.CFIRule{}.ID(),
		elf.FortifyLevel2Rule{}.ID(),
		elf.KernelCFIRule{}.ID(),
	},
}

var All = map[string]Profile{
	"minimal":     Minimal,
	"recommended": Recommended,
	"hardened":    Hardened,
	"build-test":  BuildTest,
}

func Get(name string) (Profile, bool) {
	p, ok := All[name]
	return p, ok
}

func Names() []string {
	return []string{"minimal", "recommended", "hardened", "build-test"}
}
