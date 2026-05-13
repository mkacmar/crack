package rule

import (
	"testing"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/toolchain"
)

func TestCheckApplicability(t *testing.T) {
	gccOnly := map[toolchain.Compiler]CompilerRequirement{
		toolchain.GCC: {},
	}

	tests := []struct {
		name    string
		app     Applicability
		profile binary.Profile
		want    ApplicabilityResult
	}{
		{
			name: "matching arch, no compiler/libc constraints",
			app:  Applicability{Platform: binary.PlatformAll, LibC: binary.LibCAll},
			profile: binary.Profile{
				Architecture: binary.ArchAMD64,
				Toolchain:    toolchain.Toolchain{Compiler: toolchain.GCC},
				LibC:         binary.LibCGlibc,
			},
			want: Applicable,
		},
		{
			name: "architecture mismatch",
			app:  Applicability{Platform: binary.Platform{Architecture: binary.ArchARM64}},
			profile: binary.Profile{
				Architecture: binary.ArchAMD64,
			},
			want: NotApplicableArchitecture,
		},
		{
			name: "PlatformAll matches any arch",
			app:  Applicability{Platform: binary.PlatformAll},
			profile: binary.Profile{
				Architecture: binary.ArchRISCV,
			},
			want: Applicable,
		},
		{
			name: "compiler in allowed set",
			app: Applicability{
				Platform:  binary.PlatformAll,
				Compilers: gccOnly,
			},
			profile: binary.Profile{
				Architecture: binary.ArchAMD64,
				Toolchain:    toolchain.Toolchain{Compiler: toolchain.GCC},
			},
			want: Applicable,
		},
		{
			name: "compiler not in allowed set",
			app: Applicability{
				Platform:  binary.PlatformAll,
				Compilers: gccOnly,
			},
			profile: binary.Profile{
				Architecture: binary.ArchAMD64,
				Toolchain:    toolchain.Toolchain{Compiler: toolchain.Clang},
			},
			want: NotApplicableCompiler,
		},
		{
			name: "empty compiler map allows any compiler",
			app: Applicability{
				Platform:  binary.PlatformAll,
				Compilers: map[toolchain.Compiler]CompilerRequirement{},
			},
			profile: binary.Profile{
				Architecture: binary.ArchAMD64,
				Toolchain:    toolchain.Toolchain{Compiler: toolchain.Clang},
			},
			want: Applicable,
		},
		{
			name: "unknown compiler bypasses filter (best-effort)",
			app: Applicability{
				Platform:  binary.PlatformAll,
				Compilers: gccOnly,
			},
			profile: binary.Profile{
				Architecture: binary.ArchAMD64,
				Toolchain:    toolchain.Toolchain{Compiler: toolchain.Unknown},
			},
			want: Applicable,
		},
		{
			name: "libc bitmask match",
			app: Applicability{
				Platform: binary.PlatformAll,
				LibC:     binary.LibCAll,
			},
			profile: binary.Profile{
				Architecture: binary.ArchAMD64,
				LibC:         binary.LibCGlibc,
			},
			want: Applicable,
		},
		{
			name: "libc specific match (glibc to glibc)",
			app: Applicability{
				Platform: binary.PlatformAll,
				LibC:     binary.LibCGlibc,
			},
			profile: binary.Profile{
				Architecture: binary.ArchAMD64,
				LibC:         binary.LibCGlibc,
			},
			want: Applicable,
		},
		{
			name: "libc mismatch (glibc rule, musl binary)",
			app: Applicability{
				Platform: binary.PlatformAll,
				LibC:     binary.LibCGlibc,
			},
			profile: binary.Profile{
				Architecture: binary.ArchAMD64,
				LibC:         binary.LibCMusl,
			},
			want: NotApplicableLibC,
		},
		{
			name: "libc-specific rule skips static binary (LibCNone)",
			app: Applicability{
				Platform: binary.PlatformAll,
				LibC:     binary.LibCGlibc,
			},
			profile: binary.Profile{
				Architecture: binary.ArchAMD64,
				LibC:         binary.LibCNone,
			},
			want: NotApplicableLibC,
		},
		{
			name: "unknown libc bypasses filter (best-effort)",
			app: Applicability{
				Platform: binary.PlatformAll,
				LibC:     binary.LibCGlibc,
			},
			profile: binary.Profile{
				Architecture: binary.ArchAMD64,
				LibC:         binary.LibCUnknown,
			},
			want: Applicable,
		},
		{
			name: "architecture failure shadows compiler failure",
			app: Applicability{
				Platform:  binary.Platform{Architecture: binary.ArchARM64},
				Compilers: gccOnly,
			},
			profile: binary.Profile{
				Architecture: binary.ArchAMD64,
				Toolchain:    toolchain.Toolchain{Compiler: toolchain.Clang},
			},
			want: NotApplicableArchitecture,
		},
		{
			name: "compiler failure shadows libc failure",
			app: Applicability{
				Platform:  binary.PlatformAll,
				Compilers: gccOnly,
				LibC:      binary.LibCGlibc,
			},
			profile: binary.Profile{
				Architecture: binary.ArchAMD64,
				Toolchain:    toolchain.Toolchain{Compiler: toolchain.Clang},
				LibC:         binary.LibCMusl,
			},
			want: NotApplicableCompiler,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CheckApplicability(tc.app, tc.profile)
			if got != tc.want {
				t.Errorf("CheckApplicability() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestApplicabilityResultString(t *testing.T) {
	tests := []struct {
		r    ApplicabilityResult
		want string
	}{
		{Applicable, ""},
		{NotApplicableArchitecture, "architecture not applicable"},
		{NotApplicableCompiler, "compiler not applicable"},
		{NotApplicableLibC, "libc not applicable"},
	}
	for _, tc := range tests {
		if got := tc.r.String(); got != tc.want {
			t.Errorf("%v.String() = %q, want %q", tc.r, got, tc.want)
		}
	}
}

func TestApplicabilityResultReason(t *testing.T) {
	profile := binary.Profile{
		Architecture: binary.ArchAMD64,
		Toolchain:    toolchain.Toolchain{Compiler: toolchain.Clang},
		LibC:         binary.LibCMusl,
	}
	tests := []struct {
		r    ApplicabilityResult
		want string
	}{
		{Applicable, ""},
		{NotApplicableArchitecture, "rule not applicable to amd64 architecture"},
		{NotApplicableCompiler, "rule not applicable to clang binaries"},
		{NotApplicableLibC, "rule not applicable to musl binaries"},
	}
	for _, tc := range tests {
		if got := tc.r.Reason(profile); got != tc.want {
			t.Errorf("%v.Reason() = %q, want %q", tc.r, got, tc.want)
		}
	}
}
