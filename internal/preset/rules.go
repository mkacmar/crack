package preset

import (
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/rule/elf"
)

func Default() []rule.ELFRule {
	return []rule.ELFRule{
		elf.ASLRRule{},
		elf.FortifySourceRule{},
		elf.FullRELRORule{},
		elf.NoInsecureRPATHRule{},
		elf.NoInsecureRUNPATHRule{},
		elf.NXBitRule{},
		elf.PIERule{},
		elf.RELRORule{},
		elf.SeparateCodeRule{},
		elf.StackCanaryRule{},
	}
}
