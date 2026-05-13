package elf

import (
	"debug/elf"
	"encoding/binary"
	"errors"
	"fmt"
)

// GNU note types and property constants for feature detection.
const (
	NT_GNU_BUILD_ID        = 3
	NT_GNU_PROPERTY_TYPE_0 = 5

	GNU_PROPERTY_X86_FEATURE_1_AND   = 0xc0000002
	GNU_PROPERTY_X86_FEATURE_1_IBT   = 0x1
	GNU_PROPERTY_X86_FEATURE_1_SHSTK = 0x2

	GNU_PROPERTY_AARCH64_FEATURE_1_AND = 0xc0000000
	GNU_PROPERTY_AARCH64_FEATURE_1_BTI = 0x1
	GNU_PROPERTY_AARCH64_FEATURE_1_PAC = 0x2
)

// gnuNoteName is the vendor name string stored in every GNU-defined note (NUL-terminated).
const gnuNoteName = "GNU\x00"

// noteHeaderSize is the size of an ELF note header: three uint32s (namesz, descsz, type).
const noteHeaderSize = 12

// noteNameAlign is the padding alignment applied to a note's name field per the ELF spec.
const noteNameAlign Alignment = 4

// HasGNUProperty reports whether the binary has a GNU property with the specified feature flag set under the given property type.
func HasGNUProperty(b Binary, propertyType, featureFlag uint32) (bool, error) {
	sec, err := FindSection(b, ".note.gnu.property")
	if err != nil {
		if errors.Is(err, ErrSectionMissing) {
			return false, nil
		}
		return false, err
	}
	data, err := sec.Data()
	if err != nil {
		if errors.Is(err, ErrSectionMissing) {
			return false, nil
		}
		return false, err
	}

	bo := b.ByteOrder()
	// Note descriptor padding: 4 by default, 8 on some 64-bit ABIs.
	descAlign := Alignment(4)
	if sec.Addralign == 8 {
		descAlign = 8
	}
	// GNU property records are padded to 4 bytes on 32-bit and 8 bytes on 64-bit ELF.
	propAlign := Alignment(4)
	if b.Class() == elf.ELFCLASS64 {
		propAlign = 8
	}

	var found bool
	walkNotes(data, bo, descAlign, func(noteType uint32, name, desc []byte) bool {
		if noteType != NT_GNU_PROPERTY_TYPE_0 || string(name) != gnuNoteName {
			return false
		}
		walkGNUProperties(desc, bo, propAlign, func(propType uint32, propData []byte) bool {
			if propType == propertyType && len(propData) >= 4 && bo.Uint32(propData[:4])&featureFlag != 0 {
				found = true
				return true
			}
			return false
		})
		return found
	})

	return found, nil
}

func extractBuildID(b Binary) string {
	sec, err := FindSection(b, ".note.gnu.build-id")
	if err != nil {
		return ""
	}
	data, err := sec.Data()
	if err != nil {
		return ""
	}

	// Note descriptor padding: 4 by default, 8 on some 64-bit ABIs.
	descAlign := Alignment(4)
	if sec.Addralign == 8 {
		descAlign = 8
	}

	var buildID string
	walkNotes(data, b.ByteOrder(), descAlign, func(noteType uint32, name, desc []byte) bool {
		if noteType != NT_GNU_BUILD_ID || string(name) != gnuNoteName {
			return false
		}
		buildID = fmt.Sprintf("%x", desc)
		return true
	})
	return buildID
}

// walkNotes iterates note records in data, invoking fn for each.
// Iteration stops if fn returns true or the remaining data can't hold another note.
func walkNotes(data []byte, bo binary.ByteOrder, descAlign Alignment, fn func(noteType uint32, name, desc []byte) bool) {
	offset := 0
	for offset+noteHeaderSize <= len(data) {
		namesz := bo.Uint32(data[offset : offset+4])
		descsz := bo.Uint32(data[offset+4 : offset+8])
		noteType := bo.Uint32(data[offset+8 : offset+12])

		paddedNameLen := noteNameAlign.Pad(int(namesz))
		paddedDescLen := descAlign.Pad(int(descsz))

		nameStart := offset + noteHeaderSize
		descStart := nameStart + paddedNameLen
		if descStart+paddedDescLen > len(data) {
			return
		}

		name := data[nameStart : nameStart+int(namesz)]
		desc := data[descStart : descStart+int(descsz)]
		if fn(noteType, name, desc) {
			return
		}

		offset = descStart + paddedDescLen
	}
}

// walkGNUProperties iterates property records inside the descriptor of a GNU property note, invoking fn for each.
// Iteration stops if fn returns true.
func walkGNUProperties(desc []byte, bo binary.ByteOrder, align Alignment, fn func(propType uint32, propData []byte) bool) {
	offset := 0
	for offset+8 <= len(desc) {
		propType := bo.Uint32(desc[offset : offset+4])
		propSize := bo.Uint32(desc[offset+4 : offset+8])

		dataStart := offset + 8
		dataEnd := dataStart + int(propSize)
		if dataEnd > len(desc) {
			return
		}

		if fn(propType, desc[dataStart:dataEnd]) {
			return
		}

		offset = dataStart + align.Pad(int(propSize))
	}
}
