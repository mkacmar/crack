package elf

import (
	"debug/elf"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

var fortifiedFunctions = []string{
	"__memcpy_chk",
	"__memmove_chk",
	"__memset_chk",
	"__strcpy_chk",
	"__strncpy_chk",
	"__strcat_chk",
	"__strncat_chk",
	"__sprintf_chk",
	"__snprintf_chk",
	"__vsprintf_chk",
	"__vsnprintf_chk",
	"__gets_chk",
	"__read_chk",
	"__pread_chk",
	"__fgets_chk",
	"__fread_chk",
	"__recv_chk",
	"__recvfrom_chk",
	"__realpath_chk",
	"__wcsncpy_chk",
	"__wcscpy_chk",
	"__wcscat_chk",
	"__stpcpy_chk",
	"__stpncpy_chk",
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
			{Compiler: model.CompilerGCC, MinVersion: model.Version{Major: 4, Minor: 0}, Flag: "-D_FORTIFY_SOURCE=3 -O1"},
			{Compiler: model.CompilerClang, MinVersion: model.Version{Major: 5, Minor: 0}, Flag: "-D_FORTIFY_SOURCE=3 -O1"},
		},
	}
}

func (r FortifySourceRule) Execute(f *elf.File, info *model.ParsedBinary) model.RuleResult {
	hasFortifySource := false

	symbols, err := f.Symbols()
	if err != nil {
		symbols = nil
	}

	dynsyms, err := f.DynamicSymbols()
	if err != nil {
		dynsyms = nil
	}

	for _, sym := range symbols {
		for _, fortified := range fortifiedFunctions {
			if strings.Contains(sym.Name, fortified) {
				hasFortifySource = true
				break
			}
		}
		if hasFortifySource {
			break
		}
	}

	if !hasFortifySource {
		for _, sym := range dynsyms {
			for _, fortified := range fortifiedFunctions {
				if strings.Contains(sym.Name, fortified) {
					hasFortifySource = true
					break
				}
			}
			if hasFortifySource {
				break
			}
		}
	}

	if hasFortifySource {
		return model.RuleResult{
			State:   model.CheckStatePassed,
			Message: "FORTIFY_SOURCE is enabled (fortified functions detected)",
		}
	}
	return model.RuleResult{
		State:   model.CheckStateFailed,
		Message: "FORTIFY_SOURCE is NOT enabled",
	}
}
