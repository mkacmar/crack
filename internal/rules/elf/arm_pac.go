package elf

import (
	"debug/elf"

	"github.com/mkacmar/crack/internal/model"
)

const (
	GNU_PROPERTY_AARCH64_FEATURE_1_AND = 0xc0000000
	GNU_PROPERTY_AARCH64_FEATURE_1_BTI = 0x1
	GNU_PROPERTY_AARCH64_FEATURE_1_PAC = 0x2
)

// ARMPACRule checks for ARM Pointer Authentication Code
// ARM: https://developer.arm.com/documentation/ddi0487/latest
// GCC: https://gcc.gnu.org/onlinedocs/gcc/AArch64-Options.html#index-mbranch-protection
// Clang: https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-mbranch-protection
type ARMPACRule struct{}

func (r ARMPACRule) ID() string                     { return "arm-pac" }
func (r ARMPACRule) Name() string                   { return "ARM Pointer Authentication" }
func (r ARMPACRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r ARMPACRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r ARMPACRule) TargetArch() model.Architecture { return model.ArchARM64 }
func (r ARMPACRule) HasPerfImpact() bool            { return false }

func (r ARMPACRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 9, Minor: 1}, Flag: "-mbranch-protection=pac-ret"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 8, Minor: 0}, Flag: "-mbranch-protection=pac-ret"},
		},
	}
}

func (r ARMPACRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasPAC := parseGNUPropertyForARM64Feature(f, GNU_PROPERTY_AARCH64_FEATURE_1_PAC)

	if hasPAC {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "ARM PAC (Pointer Authentication Code) is enabled",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "ARM PAC is NOT enabled (requires ARMv8.3+ hardware)",
	}
}

func parseGNUPropertyForARM64Feature(f *elf.File, featureFlag uint32) bool {
	sec := f.Section(".note.gnu.property")
	if sec == nil {
		return false
	}

	data, err := sec.Data()
	if err != nil || len(data) < 16 {
		return false
	}

	// Note format:
	// - namesz (4 bytes)
	// - descsz (4 bytes)
	// - type (4 bytes)
	// - name (aligned to 4 bytes)
	// - desc (aligned to 4 bytes)

	offset := 0
	for offset+12 <= len(data) {
		namesz := f.ByteOrder.Uint32(data[offset : offset+4])
		descsz := f.ByteOrder.Uint32(data[offset+4 : offset+8])
		noteType := f.ByteOrder.Uint32(data[offset+8 : offset+12])

		// Align namesz to 4 or 8 bytes depending on ELF class
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

					if propType == GNU_PROPERTY_AARCH64_FEATURE_1_AND && propSize >= 4 {
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
