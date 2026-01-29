package preset

import (
	"github.com/mkacmar/crack/internal/rules/elf"
)

var DefaultRules = []string{
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
}
