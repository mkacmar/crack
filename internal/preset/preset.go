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
		elf.ASLRRule{}.ID(),
		elf.FortifySourceRule{}.ID(),
		elf.NXBitRule{}.ID(),
		elf.PIERule{}.ID(),
		elf.RELRORule{}.ID(),
		elf.StackCanaryRule{}.ID(),
	},
}

var Recommended = Preset{
	Rules: []string{
		elf.ASLRRule{}.ID(),
		elf.FortifySourceRule{}.ID(),
		elf.FullRELRORule{}.ID(),
		elf.NoInsecureRPATHRule{}.ID(),
		elf.NoInsecureRUNPATHRule{}.ID(),
		elf.NXBitRule{}.ID(),
		elf.PIERule{}.ID(),
		elf.RELRORule{}.ID(),
		elf.SeparateCodeRule{}.ID(),
		elf.StackCanaryRule{}.ID(),
		elf.WXorXRule{}.ID(),
	},
}

var Hardened = Preset{
	Rules: []string{
		elf.ARMBranchProtectionRule{}.ID(),
		elf.ARMBTIRule{}.ID(),
		elf.ARMMTERule{}.ID(),
		elf.ARMPACRule{}.ID(),
		elf.ASLRRule{}.ID(),
		elf.CFIRule{}.ID(),
		elf.FortifySourceRule{}.ID(),
		elf.FullRELRORule{}.ID(),
		elf.NoDLOpenRule{}.ID(),
		elf.NoDumpRule{}.ID(),
		elf.NoInsecureRPATHRule{}.ID(),
		elf.NoInsecureRUNPATHRule{}.ID(),
		elf.NoPLTRule{}.ID(),
		elf.NXBitRule{}.ID(),
		elf.PIERule{}.ID(),
		elf.RELRORule{}.ID(),
		elf.RetpolineRule{}.ID(),
		elf.SafeStackRule{}.ID(),
		elf.SeparateCodeRule{}.ID(),
		elf.StackCanaryRule{}.ID(),
		elf.StackLimitRule{}.ID(),
		elf.StrippedRule{}.ID(),
		elf.UBSanRule{}.ID(),
		elf.WXorXRule{}.ID(),
		elf.X86CETIBTRule{}.ID(),
		elf.X86CETShadowStackRule{}.ID(),
	},
}

var BuildTest = Preset{
	Rules: []string{
		elf.ASANRule{}.ID(),
		elf.CFIRule{}.ID(),
		elf.SafeStackRule{}.ID(),
		elf.UBSanRule{}.ID(),
	},
}

var Kernel = Preset{
	Rules: []string{
		elf.ARMBranchProtectionRule{}.ID(),
		elf.ARMBTIRule{}.ID(),
		elf.ARMPACRule{}.ID(),
		elf.KernelCFIRule{}.ID(),
		elf.RetpolineRule{}.ID(),
		elf.StackCanaryRule{}.ID(),
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
