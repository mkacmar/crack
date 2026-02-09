package elf

import (
	"fmt"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/rule"
	"github.com/mkacmar/crack/toolchain"
)

// FortifySourceRuleID is the rule ID for FORTIFY_SOURCE.
const FortifySourceRuleID = "fortify-source"

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

// FortifySourceRule checks for FORTIFY_SOURCE protection.
// glibc: https://sourceware.org/glibc/wiki/FortifySourceLevel3
// GCC: https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-D_FORTIFY_SOURCE
type FortifySourceRule struct{}

func (r FortifySourceRule) ID() string   { return FortifySourceRuleID }
func (r FortifySourceRule) Name() string { return "FORTIFY_SOURCE" }
func (r FortifySourceRule) Description() string {
	return "Checks for FORTIFY_SOURCE buffer overflow protection. This glibc feature replaces unsafe C library functions (strcpy, memcpy, sprintf, etc.) with bounds-checked variants at compile time."
}

func (r FortifySourceRule) Applicability() rule.Applicability {
	return rule.Applicability{
		Platform: binary.Platform{Architecture: binary.ArchAllX86 | binary.ArchAllARM},
		Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
			toolchain.GCC:   {MinVersion: toolchain.Version{Major: 12, Minor: 1}, Flag: "-D_FORTIFY_SOURCE=3 -O1"},
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 12, Minor: 0}, Flag: "-D_FORTIFY_SOURCE=3 -O1"},
		},
	}
}

func (r FortifySourceRule) Execute(bin *binary.ELFBinary) rule.Result {
	if bin.LibC == binary.LibCMusl {
		return rule.Result{
			Status:  rule.StatusSkipped,
			Message: "musl libc, FORTIFY_SOURCE not supported",
		}
	}

	allSymbols := make(map[string]struct{})
	for _, sym := range bin.Symbols {
		allSymbols[sym.Name] = struct{}{}
	}
	for _, sym := range bin.DynSymbols {
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
		msg := fmt.Sprintf("FORTIFY_SOURCE enabled (%d fortified)", len(fortifiedFuncs))
		if len(unfortifiedFuncs) > 0 {
			msg = fmt.Sprintf("FORTIFY_SOURCE enabled (%d fortified, %d unfortified)", len(fortifiedFuncs), len(unfortifiedFuncs))
		}
		return rule.Result{
			Status:  rule.StatusPassed,
			Message: msg,
		}
	}

	if len(unfortifiedFuncs) > 0 {
		return rule.Result{
			Status:  rule.StatusFailed,
			Message: "FORTIFY_SOURCE not enabled",
		}
	}

	return rule.Result{
		Status:  rule.StatusSkipped,
		Message: "No fortifiable functions detected",
	}
}
