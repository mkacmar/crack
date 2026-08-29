package preset

import "go.kacmar.sk/crack/rule/elf"

func Default() []string {
	return []string{
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
	}
}
