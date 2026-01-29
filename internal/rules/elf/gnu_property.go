package elf

import "debug/elf"

const (
	GNU_PROPERTY_X86_FEATURE_1_AND   = 0xc0000002
	GNU_PROPERTY_X86_FEATURE_1_IBT   = 0x1
	GNU_PROPERTY_X86_FEATURE_1_SHSTK = 0x2
)

const (
	GNU_PROPERTY_AARCH64_FEATURE_1_AND = 0xc0000000
	GNU_PROPERTY_AARCH64_FEATURE_1_BTI = 0x1
	GNU_PROPERTY_AARCH64_FEATURE_1_PAC = 0x2
)

// parseGNUProperty parses .note.gnu.property section and checks for a feature flag.
// propertyType is the property type (e.g., GNU_PROPERTY_X86_FEATURE_1_AND)
// featureFlag is the specific feature bit to check (e.g., GNU_PROPERTY_X86_FEATURE_1_IBT)
func parseGNUProperty(f *elf.File, propertyType uint32, featureFlag uint32) bool {
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

		// Note name is always 4-byte aligned per ELF spec.
		// Descriptor alignment depends on ELF class for .note.gnu.property.
		align := 4
		if f.Class == elf.ELFCLASS64 {
			align = 8
		}
		alignedNamesz := (int(namesz) + 3) &^ 3
		alignedDescsz := (int(descsz) + align - 1) &^ (align - 1)

		nameStart := offset + 12
		descStart := nameStart + alignedNamesz

		if descStart+alignedDescsz > len(data) {
			break
		}

		// Check if this is a GNU note (NT_GNU_PROPERTY_TYPE_0 = 5).
		if noteType == 5 && namesz >= 4 {
			name := string(data[nameStart : nameStart+4])
			if name == "GNU\x00" {
				// Parse properties in the descriptor.
				propOffset := descStart
				propEnd := descStart + int(descsz)
				for propOffset+8 <= propEnd {
					propType := f.ByteOrder.Uint32(data[propOffset : propOffset+4])
					propSize := f.ByteOrder.Uint32(data[propOffset+4 : propOffset+8])

					if propType == propertyType && propSize >= 4 {
						features := f.ByteOrder.Uint32(data[propOffset+8 : propOffset+12])
						if features&featureFlag != 0 {
							return true
						}
					}

					// Move to next property (aligned).
					alignedPropSize := (int(propSize) + align - 1) &^ (align - 1)
					propOffset += 8 + alignedPropSize
				}
			}
		}

		offset += 12 + alignedNamesz + alignedDescsz
	}

	return false
}
