package preset

import (
	"maps"
	"slices"

	"github.com/mkacmar/crack/internal/rules/elf"
)

const (
	NameMinimal     = "minimal"
	NameRecommended = "recommended"
	NameHardened    = "hardened"
	NameBuildTest   = "build-test"
	NameKernel      = "kernel"

	Default = NameRecommended
)

type Preset struct {
	Rules []string
}

var Minimal = Preset{
	Rules: []string{
		elf.ASLRRuleID,
		elf.FortifySourceRuleID,
		elf.NXBitRuleID,
		elf.PIERuleID,
		elf.RELRORuleID,
		elf.StackCanaryRuleID,
	},
}

var Recommended = Preset{
	Rules: []string{
		elf.ASLRRuleID,
		elf.FortifySourceRuleID,
		elf.FullRELRORuleID,
		elf.NoInsecureRPATHRuleID,
		elf.NoInsecureRUNPATHRuleID,
		elf.NXBitRuleID,
		elf.PIERuleID,
		elf.RELRORuleID,
		elf.SeparateCodeRuleID,
		elf.StackCanaryRuleID,
		elf.WXorXRuleID,
	},
}

var Hardened = Preset{
	Rules: []string{
		elf.ARMBranchProtectionRuleID,
		elf.ARMBTIRuleID,
		elf.ARMMTERuleID,
		elf.ARMPACRuleID,
		elf.ASLRRuleID,
		elf.CFIRuleID,
		elf.FortifySourceRuleID,
		elf.FullRELRORuleID,
		elf.NoDLOpenRuleID,
		elf.NoDumpRuleID,
		elf.NoInsecureRPATHRuleID,
		elf.NoInsecureRUNPATHRuleID,
		elf.NoPLTRuleID,
		elf.NXBitRuleID,
		elf.PIERuleID,
		elf.RELRORuleID,
		elf.SafeStackRuleID,
		elf.SeparateCodeRuleID,
		elf.StackCanaryRuleID,
		elf.StackLimitRuleID,
		elf.StrippedRuleID,
		elf.UBSanRuleID,
		elf.WXorXRuleID,
		elf.X86CETIBTRuleID,
		elf.X86CETShadowStackRuleID,
		elf.X86RetpolineRuleID,
	},
}

var BuildTest = Preset{
	Rules: []string{
		elf.ASANRuleID,
		elf.CFIRuleID,
		elf.SafeStackRuleID,
		elf.UBSanRuleID,
	},
}

var Kernel = Preset{
	Rules: []string{
		elf.ARMBranchProtectionRuleID,
		elf.ARMBTIRuleID,
		elf.ARMPACRuleID,
		elf.KernelCFIRuleID,
		elf.StackCanaryRuleID,
		elf.X86RetpolineRuleID,
	},
}

var All = map[string]Preset{
	NameMinimal:     Minimal,
	NameRecommended: Recommended,
	NameHardened:    Hardened,
	NameBuildTest:   BuildTest,
	NameKernel:      Kernel,
}

func Get(name string) (Preset, bool) {
	p, ok := All[name]
	return p, ok
}

func Names() []string {
	return slices.Sorted(maps.Keys(All))
}
