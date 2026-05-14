package elf

import (
	"bytes"
	"debug/elf"
	"errors"
	"strings"

	"go.kacmar.sk/crack/binary"
)

// DetectArchitecture returns the architecture of the binary.
func DetectArchitecture(b Binary) binary.Architecture {
	switch b.Machine() {
	case elf.EM_386:
		return binary.ArchX86
	case elf.EM_X86_64:
		return binary.ArchAMD64
	case elf.EM_ARM:
		return binary.ArchARM
	case elf.EM_AARCH64:
		return binary.ArchARM64
	case elf.EM_RISCV:
		return binary.ArchRISCV
	case elf.EM_PPC64:
		return binary.ArchPPC64
	case elf.EM_MIPS:
		return binary.ArchMIPS
	case elf.EM_S390:
		return binary.ArchS390X
	default:
		return binary.ArchUnknown
	}
}

// DetectLibC identifies the C library that the binary links against, using PT_INTERP and DT_NEEDED as evidence.
// Returns LibCNone when the binary declares no libc dependency at all (static executables and self-contained shared objects).
// Returns LibCUnknown when the binary references a libc but the specific implementation can't be classified.
func DetectLibC(b Binary) binary.LibC {
	hasInterp := false
	for _, prog := range b.Progs() {
		if prog.Type != elf.PT_INTERP {
			continue
		}
		hasInterp = true
		data, err := prog.Data()
		if err != nil {
			continue
		}
		interpreter := string(bytes.TrimRight(data, "\x00"))
		if strings.Contains(interpreter, "ld-musl") {
			return binary.LibCMusl
		}
		if strings.Contains(interpreter, "ld-linux") {
			return binary.LibCGlibc
		}
	}

	libs, err := ImportedLibraries(b)
	if err != nil && !errors.Is(err, ErrSectionMissing) {
		return binary.LibCUnknown
	}
	hasLibcDep := false
	for _, lib := range libs {
		if strings.Contains(lib, "musl") {
			return binary.LibCMusl
		}
		if lib == "libc.so.6" {
			return binary.LibCGlibc
		}
		if strings.HasPrefix(lib, "libc.so") {
			hasLibcDep = true
		}
	}

	if hasInterp || hasLibcDep {
		return binary.LibCUnknown
	}
	return binary.LibCNone
}
