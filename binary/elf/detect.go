package elf

import (
	"bytes"
	"debug/dwarf"
	"debug/elf"
	"errors"
	"strings"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/toolchain"
)

var compilerPriority = map[toolchain.Compiler]int{
	toolchain.GCC:   1,
	toolchain.Clang: 2,
	toolchain.Rustc: 3,
	toolchain.Go:    4,
}

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

// DetectToolchain identifies the compiler and version that produced the binary.
func DetectToolchain(b Binary, detector toolchain.ELFDetector) toolchain.Toolchain {
	if _, err := FindSection(b, ".go.buildinfo"); err == nil {
		return toolchain.Toolchain{Compiler: toolchain.Go}
	}

	comments := extractCompilerComments(b)

	var best toolchain.Toolchain
	bestPriority := 0
	for _, comment := range comments {
		comp, ver := detector.Detect(comment)
		if comp == toolchain.Unknown {
			continue
		}
		priority := compilerPriority[comp]
		if best.Compiler == toolchain.Unknown || priority > bestPriority {
			best = toolchain.Toolchain{Compiler: comp, Version: ver}
			bestPriority = priority
		}
	}
	return best
}

// DetectToolchainFromDWARF identifies the compiler and version from DW_AT_producer in .debug_info.
// Returns a zero-value Toolchain when DWARF is unavailable or no producer attribute is found.
func DetectToolchainFromDWARF(b Binary, detector toolchain.ELFDetector) toolchain.Toolchain {
	d, err := loadDWARF(b)
	if err != nil || d == nil {
		return toolchain.Toolchain{}
	}

	reader := d.Reader()
	for {
		entry, err := reader.Next()
		if err != nil || entry == nil {
			break
		}
		if entry.Tag != dwarf.TagCompileUnit {
			continue
		}
		producer, ok := entry.Val(dwarf.AttrProducer).(string)
		if !ok || producer == "" {
			continue
		}
		comp, ver := detector.Detect(producer)
		if comp != toolchain.Unknown {
			return toolchain.Toolchain{Compiler: comp, Version: ver}
		}
	}
	return toolchain.Toolchain{}
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

func extractCompilerComments(b Binary) []string {
	data, err := findSectionData(b, ".comment")
	if err != nil || data == nil {
		return nil
	}

	var comments []string
	for len(data) > 0 {
		idx := bytes.IndexByte(data, 0)
		if idx == -1 {
			break
		}
		if idx > 0 {
			comments = append(comments, string(data[:idx]))
		}
		data = data[idx+1:]
	}
	return comments
}

// loadDWARF assembles a *dwarf.Data sufficient for reading DW_AT_producer.
// Fetches only the sections needed for compile-unit walks and string attributes. Line, ranges, and loc are skipped to avoid pulling large debug sections via the resolver.
// Returns (nil, nil) when the mandatory sections (.debug_info, .debug_abbrev) aren't available.
func loadDWARF(b Binary) (*dwarf.Data, error) {
	abbrev, err := findSectionData(b, ".debug_abbrev")
	if err != nil || len(abbrev) == 0 {
		return nil, err
	}
	info, err := findSectionData(b, ".debug_info")
	if err != nil || len(info) == 0 {
		return nil, err
	}
	str, err := findSectionData(b, ".debug_str")
	if err != nil {
		return nil, err
	}

	// dwarf.New positional params: abbrev, aranges, frame, info, line, pubnames, ranges, str.
	return dwarf.New(abbrev, nil, nil, info, nil, nil, nil, str)
}
