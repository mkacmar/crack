package elf

import (
	"debug/elf"
)

// HasDynFlag reports whether a dynamic tag has the specified flag set.
func HasDynFlag(b Binary, tag elf.DynTag, flag uint64) (bool, error) {
	entries, err := b.DynEntries()
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if entry.Tag == tag && (entry.Val&flag) != 0 {
			return true, nil
		}
	}
	return false, nil
}

// HasDynTag reports whether a dynamic tag exists.
func HasDynTag(b Binary, tag elf.DynTag) (bool, error) {
	entries, err := b.DynEntries()
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if entry.Tag == tag {
			return true, nil
		}
	}
	return false, nil
}

// DynString returns the string value associated with the first occurrence of the given dynamic tag, or "" if the tag is absent.
func DynString(b Binary, tag elf.DynTag) (string, error) {
	entries, err := b.DynEntries()
	if err != nil {
		return "", err
	}

	var val uint64
	var found bool
	for _, entry := range entries {
		if entry.Tag == tag {
			val = entry.Val
			found = true
			break
		}
	}
	if !found {
		return "", nil
	}

	return readDynstrEntry(b, val)
}

// ImportedLibraries reports the dynamically-linked shared library dependencies via DT_NEEDED entries.
func ImportedLibraries(b Binary) ([]string, error) {
	entries, err := b.DynEntries()
	if err != nil {
		return nil, err
	}

	var needed []uint64
	for _, entry := range entries {
		if entry.Tag == elf.DT_NEEDED {
			needed = append(needed, entry.Val)
		}
	}
	if len(needed) == 0 {
		return nil, nil
	}

	strtab, err := findSectionData(b, ".dynstr")
	if err != nil || strtab == nil {
		return nil, err
	}

	size := uint64(len(strtab))
	libs := make([]string, 0, len(needed))
	for _, off := range needed {
		if off >= size {
			continue
		}
		end := off
		for end < size && strtab[end] != 0 {
			end++
		}
		libs = append(libs, string(strtab[off:end]))
	}
	return libs, nil
}

func readDynstrEntry(b Binary, off uint64) (string, error) {
	strtab, err := findSectionData(b, ".dynstr")
	if err != nil || strtab == nil {
		return "", err
	}
	if off > uint64(^uint32(0)) {
		return "", nil
	}
	return lookupStr(strtab, uint32(off)), nil
}
