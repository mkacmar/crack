package preset

import (
	"maps"
	"slices"

	"github.com/mkacmar/crack/internal/rules/elf"
)

type Preset struct {
	Rules []string
}

var Minimal = Preset{
	Rules: []string{
		elf.PIERule{}.ID(),
		elf.NXBitRule{}.ID(),
		elf.RELRORule{}.ID(),
		elf.StackCanaryRule{}.ID(),
		elf.FortifySourceRule{}.ID(),
		elf.ASLRRule{}.ID(),
	},
}

var Recommended = Preset{
	Rules: []string{
		elf.PIERule{}.ID(),
		elf.NXBitRule{}.ID(),
		elf.RELRORule{}.ID(),
		elf.FullRELRORule{}.ID(),
		elf.StackCanaryRule{}.ID(),
		elf.FortifySourceRule{}.ID(),
		elf.ASLRRule{}.ID(),
		elf.WXorXRule{}.ID(),
		elf.NoInsecureRPATHRule{}.ID(),
		elf.NoInsecureRUNPATHRule{}.ID(),
		elf.SeparateCodeRule{}.ID(),
	},
}

var Hardened = Preset{
	Rules: []string{
		elf.PIERule{}.ID(),
		elf.NXBitRule{}.ID(),
		elf.RELRORule{}.ID(),
		elf.FullRELRORule{}.ID(),
		elf.StackCanaryRule{}.ID(),
		elf.FortifySourceRule{}.ID(),
		elf.ASLRRule{}.ID(),
		elf.WXorXRule{}.ID(),
		elf.NoInsecureRPATHRule{}.ID(),
		elf.NoInsecureRUNPATHRule{}.ID(),
		elf.SeparateCodeRule{}.ID(),
		elf.NoDLOpenRule{}.ID(),
		elf.NoPLTRule{}.ID(),
		elf.NoDumpRule{}.ID(),
		elf.StrippedRule{}.ID(),
		elf.StackLimitRule{}.ID(),
		elf.CFIRule{}.ID(),
		elf.SafeStackRule{}.ID(),
		elf.UBSanRule{}.ID(),
		elf.IntelCETIBTRule{}.ID(),
		elf.IntelCETShadowStackRule{}.ID(),
		elf.RetpolineRule{}.ID(),
		elf.ARMBranchProtectionRule{}.ID(),
		elf.ARMPACRule{}.ID(),
		elf.ARMBTIRule{}.ID(),
		elf.ARMMTERule{}.ID(),
	},
}

var BuildTest = Preset{
	Rules: []string{
		elf.ASANRule{}.ID(),
		elf.UBSanRule{}.ID(),
		elf.SafeStackRule{}.ID(),
		elf.CFIRule{}.ID(),
	},
}

var Kernel = Preset{
	Rules: []string{
		elf.KernelCFIRule{}.ID(),
		elf.StackCanaryRule{}.ID(),
		elf.RetpolineRule{}.ID(),
		elf.ARMBranchProtectionRule{}.ID(),
		elf.ARMPACRule{}.ID(),
		elf.ARMBTIRule{}.ID(),
	},
}

var All = map[string]Preset{
	"minimal":     Minimal,
	"recommended": Recommended,
	"hardened":    Hardened,
	"build-test":  BuildTest,
	"kernel":      Kernel,
}

func Get(name string) (Preset, bool) {
	p, ok := All[name]
	return p, ok
}

func Names() []string {
	return slices.Sorted(maps.Keys(All))
}
