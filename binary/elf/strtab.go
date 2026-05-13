package elf

import "bytes"

// lookupStr returns the NUL-terminated string at off in strtab, or "" if strtab is missing or off is out of range.
func lookupStr(strtab []byte, off uint32) string {
	if off == 0 || int(off) >= len(strtab) {
		return ""
	}
	end := bytes.IndexByte(strtab[off:], 0)
	if end < 0 {
		return string(strtab[off:])
	}
	return string(strtab[off : int(off)+end])
}
