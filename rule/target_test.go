package rule

import (
	"testing"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/toolchain"
)

var (
	appARMPAC = Applicability{
		Platform: binary.PlatformARM64v83,
		Compilers: map[toolchain.Compiler]CompilerRequirement{
			toolchain.GCC:   {},
			toolchain.Clang: {},
		},
	}
	appX86CET = Applicability{
		Platform: binary.PlatformAllX86,
		Compilers: map[toolchain.Compiler]CompilerRequirement{
			toolchain.GCC:   {},
			toolchain.Clang: {},
		},
	}
	appClang18OrHigherOnly = Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]CompilerRequirement{
			toolchain.Clang: {MinVersion: toolchain.Version{Major: 18}},
		},
	}
	appGCCOnly = Applicability{
		Platform: binary.PlatformAll,
		Compilers: map[toolchain.Compiler]CompilerRequirement{
			toolchain.GCC: {},
		},
	}
	appARM64Clang = Applicability{
		Platform: binary.Platform{Architecture: binary.ArchARM64},
		Compilers: map[toolchain.Compiler]CompilerRequirement{
			toolchain.Clang: {},
		},
	}
	appARM64GCC = Applicability{
		Platform: binary.Platform{Architecture: binary.ArchARM64},
		Compilers: map[toolchain.Compiler]CompilerRequirement{
			toolchain.GCC: {},
		},
	}
)

func TestTargetFilterMatches(t *testing.T) {
	tests := []struct {
		name   string
		filter TargetFilter
		app    Applicability
		want   bool
	}{
		{
			name:   "empty filter matches everything",
			filter: TargetFilter{},
			app:    appARMPAC,
			want:   true,
		},
		{
			name:   "single platform match",
			filter: TargetFilter{Platforms: []PlatformTarget{{Architecture: binary.ArchARM64}}},
			app:    appARMPAC,
			want:   true,
		},
		{
			name:   "single platform mismatch",
			filter: TargetFilter{Platforms: []PlatformTarget{{Architecture: binary.ArchAMD64}}},
			app:    appARMPAC,
			want:   false,
		},
		{
			name: "multi platform union keeps arm64 rule",
			filter: TargetFilter{Platforms: []PlatformTarget{
				{Architecture: binary.ArchAMD64},
				{Architecture: binary.ArchARM64},
			}},
			app:  appARMPAC,
			want: true,
		},
		{
			name: "multi platform union keeps x86 rule",
			filter: TargetFilter{Platforms: []PlatformTarget{
				{Architecture: binary.ArchAMD64},
				{Architecture: binary.ArchARM64},
			}},
			app:  appX86CET,
			want: true,
		},
		{
			name: "platform ISA ceiling at or above rule minimum includes",
			filter: TargetFilter{Platforms: []PlatformTarget{
				{Architecture: binary.ArchARM64, MaxISA: &binary.ISA{Major: 8, Minor: 5}},
			}},
			app:  appARMPAC,
			want: true,
		},
		{
			name: "platform ISA ceiling below rule minimum excludes",
			filter: TargetFilter{Platforms: []PlatformTarget{
				{Architecture: binary.ArchARM64, MaxISA: &binary.ISA{Major: 8, Minor: 2}},
			}},
			app:  appARMPAC,
			want: false,
		},
		{
			name:   "single compiler match",
			filter: TargetFilter{Compilers: []CompilerTarget{{Compiler: toolchain.Clang}}},
			app:    appClang18OrHigherOnly,
			want:   true,
		},
		{
			name:   "single compiler mismatch",
			filter: TargetFilter{Compilers: []CompilerTarget{{Compiler: toolchain.GCC}}},
			app:    appClang18OrHigherOnly,
			want:   false,
		},
		{
			name: "multi compiler union keeps clang-only rule",
			filter: TargetFilter{Compilers: []CompilerTarget{
				{Compiler: toolchain.GCC},
				{Compiler: toolchain.Clang},
			}},
			app:  appClang18OrHigherOnly,
			want: true,
		},
		{
			name: "multi compiler union keeps gcc-only rule",
			filter: TargetFilter{Compilers: []CompilerTarget{
				{Compiler: toolchain.GCC},
				{Compiler: toolchain.Clang},
			}},
			app:  appGCCOnly,
			want: true,
		},
		{
			name: "compiler version ceiling below rule minimum excludes",
			filter: TargetFilter{Compilers: []CompilerTarget{
				{Compiler: toolchain.Clang, MaxVersion: &toolchain.Version{Major: 15}},
			}},
			app:  appClang18OrHigherOnly,
			want: false,
		},
		{
			name: "compiler version ceiling at or above rule minimum includes",
			filter: TargetFilter{Compilers: []CompilerTarget{
				{Compiler: toolchain.Clang, MaxVersion: &toolchain.Version{Major: 20}},
			}},
			app:  appClang18OrHigherOnly,
			want: true,
		},
		{
			name: "cross axis intersection keeps arm64 clang rule",
			filter: TargetFilter{
				Platforms: []PlatformTarget{{Architecture: binary.ArchARM64}},
				Compilers: []CompilerTarget{{Compiler: toolchain.Clang}},
			},
			app:  appARM64Clang,
			want: true,
		},
		{
			name: "cross axis intersection drops arm64 gcc-only rule when compiler is clang",
			filter: TargetFilter{
				Platforms: []PlatformTarget{{Architecture: binary.ArchARM64}},
				Compilers: []CompilerTarget{{Compiler: toolchain.Clang}},
			},
			app:  appARM64GCC,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.matches(tt.app); got != tt.want {
				t.Errorf("matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

type fakeRule struct {
	id  string
	app Applicability
}

func (r fakeRule) ID() string                   { return r.id }
func (r fakeRule) Name() string                 { return r.id }
func (r fakeRule) Description() string          { return r.id }
func (r fakeRule) Applicability() Applicability { return r.app }

func TestFilterRules(t *testing.T) {
	rules := []fakeRule{
		{id: "arm-pac", app: appARMPAC},
		{id: "x86-cet", app: appX86CET},
		{id: "clang-only", app: appClang18OrHigherOnly},
		{id: "gcc-only", app: appGCCOnly},
	}

	filter := &TargetFilter{Platforms: []PlatformTarget{
		{Architecture: binary.ArchAMD64},
		{Architecture: binary.ArchARM64},
	}}

	got := FilterRules(rules, filter)

	gotIDs := make(map[string]bool, len(got))
	for _, r := range got {
		gotIDs[r.id] = true
	}

	want := []string{"arm-pac", "x86-cet", "clang-only", "gcc-only"}
	if len(got) != len(want) {
		t.Fatalf("FilterRules() returned %d rules, want %d: %v", len(got), len(want), gotIDs)
	}
	for _, id := range want {
		if !gotIDs[id] {
			t.Errorf("FilterRules() missing rule %q", id)
		}
	}
}

func TestFilterRulesNilReturnsAll(t *testing.T) {
	rules := []fakeRule{
		{id: "arm-pac", app: appARMPAC},
		{id: "gcc-only", app: appGCCOnly},
	}

	if got := FilterRules(rules, nil); len(got) != len(rules) {
		t.Errorf("FilterRules(nil) returned %d rules, want %d", len(got), len(rules))
	}
}
