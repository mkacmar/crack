package elf

import (
	"debug/elf"
	"fmt"

	"github.com/mkacmar/crack/internal/model"
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

func (r FortifySourceRule) ID() string                     { return "fortify-source" }
func (r FortifySourceRule) Name() string                   { return "FORTIFY_SOURCE" }
func (r FortifySourceRule) Format() model.BinaryFormat     { return model.FormatELF }
func (r FortifySourceRule) FlagType() model.FlagType       { return model.FlagTypeCompile }
func (r FortifySourceRule) TargetArch() model.Architecture { return model.ArchAll }
func (r FortifySourceRule) HasPerfImpact() bool            { return false }

func (r FortifySourceRule) Feature() model.FeatureAvailability {
	return model.FeatureAvailability{
		Requirements: []model.CompilerRequirement{
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 12, Minor: 0}, Flag: "-D_FORTIFY_SOURCE=3 -O1"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 12, Minor: 0}, Flag: "-D_FORTIFY_SOURCE=3 -O1"},
		},
	}
}

func (r FortifySourceRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	if info != nil {
		// FORTIFY_SOURCE is a glibc feature - musl libc does not implement it.
		// https://wiki.musl-libc.org/future-ideas#fortify-source
		if info.LibC == model.LibCMusl {
			return model.RuleResult{
				State:   model.CheckStateSkipped,
				Message: "musl libc does not support FORTIFY_SOURCE",
			}
		}

		// Not applicable to standard Go/Rust toolchains
		switch info.Language {
		case model.LangGo:
			return model.RuleResult{
				State:   model.CheckStateSkipped,
				Message: "Standard Go toolchain (FORTIFY_SOURCE not applicable)",
			}
		case model.LangRust:
			return model.RuleResult{
				State:   model.CheckStateSkipped,
				Message: "Standard Rust toolchain (FORTIFY_SOURCE not applicable)",
			}
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
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: msg,
		}
	}

	// If we see unfortified functions but no _chk variants, we report a failure.
	// While the compiler could theoretically optimize away all _chk calls when it can
	// prove safety, this is rare in real-world binaries. Production code typically has
	// at least some calls where the compiler cannot prove buffer sizes at compile time.
	// Therefore, seeing only unfortified functions is strong evidence that FORTIFY_SOURCE
	// is not enabled.
	if len(unfortifiedFuncs) > 0 {
		return model.RuleResult{
			State:   model.CheckStateFailed,
			Message: fmt.Sprintf("FORTIFY_SOURCE is NOT enabled, unfortified: %v", unfortifiedFuncs),
		}
	}

	return model.RuleResult{
		State:   model.CheckStateSkipped,
		Message: "No fortifiable functions detected",
	}
}
