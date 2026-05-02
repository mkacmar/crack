package binary

import (
	"bytes"
	"debug/elf"
	"fmt"
	"strings"

	"go.kacmar.sk/crack/toolchain"
)

var compilerPriority = map[toolchain.Compiler]int{
	toolchain.GCC:   1,
	toolchain.Clang: 2,
	toolchain.Rustc: 3,
}

func isNotELFError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "bad magic number") ||
		strings.Contains(msg, "invalid argument")
}

func detectToolchain(f *elf.File, detector toolchain.ELFDetector) (toolchain.Compiler, toolchain.Version) {
	comments := extractCompilerComments(f)

	var bestComp toolchain.Compiler
	var bestVer toolchain.Version
	bestPriority := 0
	for _, comment := range comments {
		comp, ver := detector.Detect(comment)
		if comp == toolchain.Unknown {
			continue
		}
		priority := compilerPriority[comp]
		if bestComp == toolchain.Unknown || priority > bestPriority {
			bestComp = comp
			bestVer = ver
			bestPriority = priority
		}
	}
	return bestComp, bestVer
}

func extractCompilerComments(f *elf.File) []string {
	section := f.Section(".comment")
	if section == nil {
		return nil
	}

	data, err := section.Data()
	if err != nil {
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

func extractBuildID(f *elf.File) string {
	section := f.Section(".note.gnu.build-id")
	if section == nil {
		return ""
	}

	data, err := section.Data()
	if err != nil {
		return ""
	}

	const noteHeaderSize = 12
	if len(data) < noteHeaderSize {
		return ""
	}

	namesz := f.ByteOrder.Uint32(data[0:4])
	descsz := f.ByteOrder.Uint32(data[4:8])

	align := 4
	if f.Class == elf.ELFCLASS64 {
		align = 8
	}

	alignedNamesz := (int(namesz) + align - 1) &^ (align - 1)
	descOffset := noteHeaderSize + alignedNamesz

	if descOffset+int(descsz) > len(data) {
		return ""
	}

	return fmt.Sprintf("%x", data[descOffset:descOffset+int(descsz)])
}

func parseArchitecture(machine elf.Machine) Architecture {
	switch machine {
	case elf.EM_386:
		return ArchX86
	case elf.EM_X86_64:
		return ArchAMD64
	case elf.EM_ARM:
		return ArchARM
	case elf.EM_AARCH64:
		return ArchARM64
	case elf.EM_RISCV:
		return ArchRISCV
	case elf.EM_PPC64:
		return ArchPPC64
	case elf.EM_MIPS:
		return ArchMIPS
	case elf.EM_S390:
		return ArchS390X
	default:
		return ArchUnknown
	}
}

func detectLibC(f *elf.File) LibC {
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_INTERP {
			data := make([]byte, prog.Filesz)
			if _, err := prog.ReadAt(data, 0); err != nil {
				continue
			}
			interpreter := string(bytes.TrimRight(data, "\x00"))

			if strings.Contains(interpreter, "ld-musl") {
				return LibCMusl
			}
			if strings.Contains(interpreter, "ld-linux") {
				return LibCGlibc
			}
		}
	}

	// Fall back to DT_NEEDED entries for shared libraries without PT_INTERP.
	libs, err := f.ImportedLibraries()
	if err == nil {
		for _, lib := range libs {
			if strings.Contains(lib, "musl") {
				return LibCMusl
			}
			if lib == "libc.so.6" {
				return LibCGlibc
			}
		}
	}

	return LibCUnknown
}
