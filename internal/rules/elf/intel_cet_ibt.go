package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

const (
	GNU_PROPERTY_X86_FEATURE_1_AND   = 0xc0000002
	GNU_PROPERTY_X86_FEATURE_1_IBT   = 0x1
	GNU_PROPERTY_X86_FEATURE_1_SHSTK = 0x2
)

// IntelCETIBTRule checks for Intel CET Indirect Branch Tracking
// Intel: https://www.intel.com/content/www/us/en/developer/articles/technical/technical-look-control-flow-enforcement-technology.html
// GCC: https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html#index-fcf-protection
type IntelCETIBTRule struct{}

func (r IntelCETIBTRule) ID() string                     { return "intel-cet-ibt" }
func (r IntelCETIBTRule) Name() string                   { return "Intel CET - Indirect Branch Tracking" }
func (r IntelCETIBTRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r IntelCETIBTRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r IntelCETIBTRule) TargetArch() model.Architecture { return model.ArchAllX86 }
func (r IntelCETIBTRule) HasPerfImpact() bool            { return false }

func (r IntelCETIBTRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 8, Minor: 1}, Flag: "-fcf-protection=full"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 7, Minor: 0}, Flag: "-fcf-protection=full"},
		},
	}
}

func (r IntelCETIBTRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasIBT := parseGNUPropertyForX86Feature(f, GNU_PROPERTY_X86_FEATURE_1_IBT)

	if hasIBT {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "Intel CET IBT (Indirect Branch Tracking) is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "Intel CET IBT is NOT enabled (requires Intel CET-capable CPU)",
	}
}

func parseGNUPropertyForX86Feature(f *elf.File, featureFlag uint32) bool {
	sec := f.Section(".note.gnu.property")
	if sec == nil {
		return false
	}

	data, err := sec.Data()
	if err != nil || len(data) < 16 {
		return false
	}

	offset := 0
	for offset+12 <= len(data) {
		namesz := f.ByteOrder.Uint32(data[offset : offset+4])
		descsz := f.ByteOrder.Uint32(data[offset+4 : offset+8])
		noteType := f.ByteOrder.Uint32(data[offset+8 : offset+12])

		// Align to 4 or 8 bytes depending on ELF class
		align := 4
		if f.Class == elf.ELFCLASS64 {
			align = 8
		}
		alignedNamesz := (int(namesz) + align - 1) &^ (align - 1)
		alignedDescsz := (int(descsz) + align - 1) &^ (align - 1)

		nameStart := offset + 12
		descStart := nameStart + alignedNamesz

		if descStart+alignedDescsz > len(data) {
			break
		}

		// Check if this is a GNU note (NT_GNU_PROPERTY_TYPE_0 = 5)
		if noteType == 5 && namesz >= 4 {
			name := string(data[nameStart : nameStart+4])
			if name == "GNU\x00" {
				// Parse properties in the descriptor
				propOffset := descStart
				propEnd := descStart + int(descsz)
				for propOffset+8 <= propEnd {
					propType := f.ByteOrder.Uint32(data[propOffset : propOffset+4])
					propSize := f.ByteOrder.Uint32(data[propOffset+4 : propOffset+8])

					if propType == GNU_PROPERTY_X86_FEATURE_1_AND && propSize >= 4 {
						features := f.ByteOrder.Uint32(data[propOffset+8 : propOffset+12])
						if features&featureFlag != 0 {
							return true
						}
					}

					// Move to next property (aligned)
					alignedPropSize := (int(propSize) + align - 1) &^ (align - 1)
					propOffset += 8 + alignedPropSize
				}
			}
		}

		offset += 12 + alignedNamesz + alignedDescsz
	}

	return false
}
