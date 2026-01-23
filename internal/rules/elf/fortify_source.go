package elf

import (
	"debug/elf"
	"fmt"

	"github.com/mkacmar/crack/internal/binary"
	"github.com/mkacmar/crack/internal/rule"
	"github.com/mkacmar/crack/internal/toolchain"
)

var fortifiableFunctions = map[string]string{
	"fgets":     "__fgets_chk",
	"fread":     "__fread_chk",
	"gets":      "__gets_chk",
	"memcpy":    "__memcpy_chk",
	"memmove":   "__memmove_chk",
	"memset":    "__memset_chk",
	"pread":     "__pread_chk",
	"read":      "__read_chk",
	"realpath":  "__realpath_chk",
	"recv":      "__recv_chk",
	"recvfrom":  "__recvfrom_chk",
	"snprintf":  "__snprintf_chk",
	"sprintf":   "__sprintf_chk",
	"stpcpy":    "__stpcpy_chk",
	"stpncpy":   "__stpncpy_chk",
	"strcat":    "__strcat_chk",
	"strcpy":    "__strcpy_chk",
	"strncat":   "__strncat_chk",
	"strncpy":   "__strncpy_chk",
	"vsnprintf": "__vsnprintf_chk",
	"vsprintf":  "__vsprintf_chk",
	"wcscat":    "__wcscat_chk",
	"wcscpy":    "__wcscpy_chk",
	"wcsncpy":   "__wcsncpy_chk",
}

// FortifySourceRule checks for FORTIFY_SOURCE protection
// glibc: https://sourceware.org/glibc/wiki/FortifySourceLevel3
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-D_FORTIFY_SOURCE
type FortifySourceRule struct{}

func (r FortifySourceRule) ID() string   { return "fortify-source" }
func (r FortifySourceRule) Name() string { return "FORTIFY_SOURCE" }

func (r FortifySourceRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.CompilerGCC:   {MinVersion: toolchain.Version{Major: 12, Minor: 0}, Flag: "-D_FORTIFY_SOURCE=3 -O1"},
			toolchain.CompilerClang: {MinVersion: toolchain.Version{Major: 12, Minor: 0}, Flag: "-D_FORTIFY_SOURCE=3 -O1"},
		},
	}
}

func (r FortifySourceRule) Execute(f *elf.File, info *binary.Parsed) rule.ExecuteResult {
	// FORTIFY_SOURCE is a glibc feature - musl libc does not implement it.
	// https://wiki.musl-libc.org/future-ideas#fortify-source
	if info != nil && info.LibC == toolchain.LibCMusl {
		return rule.ExecuteResult{
			Status:  rule.StatusSkipped,
			Message: "musl libc does not support FORTIFY_SOURCE",
		}
	}

	symbols, _ := f.Symbols()
	dynsyms, _ := f.DynamicSymbols()

	allSymbols := make(map[string]struct{})
	for _, sym := range symbols {
		allSymbols[sym.Name] = struct{}{}
	}
	for _, sym := range dynsyms {
		allSymbols[sym.Name] = struct{}{}
	}

	var fortifiedFuncs []string
	var unfortifiedFuncs []string

	for unfortified, fortified := range fortifiableFunctions {
		_, hasFortified := allSymbols[fortified]
		_, hasUnfortified := allSymbols[unfortified]

		if hasFortified {
			fortifiedFuncs = append(fortifiedFuncs, fortified)
		}
		if hasUnfortified {
			unfortifiedFuncs = append(unfortifiedFuncs, unfortified)
		}
	}

	if len(fortifiedFuncs) > 0 {
		msg := fmt.Sprintf("FORTIFY_SOURCE is enabled. %d fortified %v", len(fortifiedFuncs), fortifiedFuncs)
		if len(unfortifiedFuncs) > 0 {
			msg += fmt.Sprintf(", %d left unfortified %v", len(unfortifiedFuncs), unfortifiedFuncs)
		}
		return rule.ExecuteResult{
			Status:  rule.StatusPassed,
			Message: msg,
		}
	}

	// If we see unfortified functions but no _chk variants, we report a failure.
	// While the compiler might optimize them away if it can prove safety, real-world binaries typically have some unprovable buffer sizes.
	if len(unfortifiedFuncs) > 0 {
		return rule.ExecuteResult{
			Status:  rule.StatusFailed,
			Message: fmt.Sprintf("FORTIFY_SOURCE is NOT enabled, unfortified: %v", unfortifiedFuncs),
		}
	}

	return rule.ExecuteResult{
		Status:  rule.StatusSkipped,
		Message: "No fortifiable functions detected",
	}
}
