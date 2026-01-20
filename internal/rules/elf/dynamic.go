package elf

import "debug/elf"

type DynEntry struct {
	Tag uint64
	Val uint64
}

func ParseDynamic(f *elf.File) []DynEntry {
	dynSec := f.Section(".dynamic")
	if dynSec == nil {
		return nil
	}

	data, err := dynSec.Data()
	if err != nil {
		return nil
	}

	var entrySize int
	var readEntry func([]byte) (tag, val uint64)

	if f.Class == elf.ELFCLASS64 {
		entrySize = 16
		readEntry = func(b []byte) (uint64, uint64) {
			return f.ByteOrder.Uint64(b[:8]), f.ByteOrder.Uint64(b[8:16])
		}
	} else {
		entrySize = 8
		readEntry = func(b []byte) (uint64, uint64) {
			return uint64(f.ByteOrder.Uint32(b[:4])), uint64(f.ByteOrder.Uint32(b[4:8]))
		}
	}

	var entries []DynEntry
	for i := 0; i+entrySize <= len(data); i += entrySize {
		tag, val := readEntry(data[i:])
		if tag == uint64(elf.DT_NULL) {
			break
		}
		entries = append(entries, DynEntry{Tag: tag, Val: val})
	}

	return entries
}

func HasDynFlag(f *elf.File, tag elf.DynTag, flag uint64) bool {
	for _, entry := range ParseDynamic(f) {
		if entry.Tag == uint64(tag) && (entry.Val&flag) != 0 {
			return true
		}
	}
	return false
}

func HasDynTag(f *elf.File, tag elf.DynTag) bool {
	for _, entry := range ParseDynamic(f) {
		if entry.Tag == uint64(tag) {
			return true
		}
	}
	return false
}
