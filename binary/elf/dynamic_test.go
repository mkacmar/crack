package elf

import (
	"debug/elf"
	"testing"
)

func TestHasDynFlag(t *testing.T) {
	tests := []struct {
		name    string
		entries []DynEntry
		tag     elf.DynTag
		flag    uint64
		want    bool
	}{
		{
			name:    "flag set",
			entries: []DynEntry{{Tag: elf.DT_FLAGS_1, Val: uint64(elf.DF_1_NOW) | uint64(elf.DF_1_PIE)}},
			tag:     elf.DT_FLAGS_1,
			flag:    uint64(elf.DF_1_PIE),
			want:    true,
		},
		{
			name:    "flag not set",
			entries: []DynEntry{{Tag: elf.DT_FLAGS_1, Val: uint64(elf.DF_1_NOW)}},
			tag:     elf.DT_FLAGS_1,
			flag:    uint64(elf.DF_1_PIE),
			want:    false,
		},
		{
			name:    "tag absent",
			entries: []DynEntry{{Tag: elf.DT_NEEDED, Val: 1}},
			tag:     elf.DT_FLAGS_1,
			flag:    uint64(elf.DF_1_PIE),
			want:    false,
		},
		{
			name: "flag set on a later entry with same tag",
			entries: []DynEntry{
				{Tag: elf.DT_FLAGS_1, Val: uint64(elf.DF_1_NOW)},
				{Tag: elf.DT_FLAGS_1, Val: uint64(elf.DF_1_PIE)},
			},
			tag:  elf.DT_FLAGS_1,
			flag: uint64(elf.DF_1_PIE),
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := HasDynFlag(&fakeBinary{dynEntry: tc.entries}, tc.tag, tc.flag)
			if err != nil {
				t.Fatalf("HasDynFlag: %v", err)
			}
			if got != tc.want {
				t.Errorf("HasDynFlag = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHasDynTag(t *testing.T) {
	tests := []struct {
		name    string
		entries []DynEntry
		tag     elf.DynTag
		want    bool
	}{
		{"present", []DynEntry{{Tag: elf.DT_BIND_NOW}}, elf.DT_BIND_NOW, true},
		{"absent", []DynEntry{{Tag: elf.DT_NEEDED, Val: 1}}, elf.DT_BIND_NOW, false},
		{"empty entries", nil, elf.DT_BIND_NOW, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := HasDynTag(&fakeBinary{dynEntry: tc.entries}, tc.tag)
			if err != nil {
				t.Fatalf("HasDynTag: %v", err)
			}
			if got != tc.want {
				t.Errorf("HasDynTag = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDynString(t *testing.T) {
	sec, entries := makeDynamic("libc.so.6", "libm.so.6")
	soname := DynEntry{Tag: elf.DT_SONAME, Val: entries[1].Val}
	fb := &fakeBinary{
		sections: []Section{sec},
		dynEntry: append(entries, soname),
	}

	tests := []struct {
		name string
		tag  elf.DynTag
		want string
	}{
		{"first DT_NEEDED resolves", elf.DT_NEEDED, "libc.so.6"},
		{"DT_SONAME resolves to second string", elf.DT_SONAME, "libm.so.6"},
		{"absent tag returns empty", elf.DT_RPATH, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DynString(fb, tc.tag)
			if err != nil {
				t.Fatalf("DynString: %v", err)
			}
			if got != tc.want {
				t.Errorf("DynString(%v) = %q, want %q", tc.tag, got, tc.want)
			}
		})
	}
}

func TestDynStringOutOfBoundsOffset(t *testing.T) {
	sec, _ := makeDynamic("libc.so.6")
	fb := &fakeBinary{
		sections: []Section{sec},
		dynEntry: []DynEntry{{Tag: elf.DT_NEEDED, Val: 9999}},
	}
	got, err := DynString(fb, elf.DT_NEEDED)
	if err != nil {
		t.Fatalf("DynString: %v", err)
	}
	if got != "" {
		t.Errorf("DynString on OOB offset = %q, want empty", got)
	}
}

func TestImportedLibraries(t *testing.T) {
	t.Run("multiple libs", func(t *testing.T) {
		sec, entries := makeDynamic("libc.so.6", "libm.so.6", "libdl.so.2")
		got, err := ImportedLibraries(&fakeBinary{sections: []Section{sec}, dynEntry: entries})
		if err != nil {
			t.Fatalf("ImportedLibraries: %v", err)
		}
		want := []string{"libc.so.6", "libm.so.6", "libdl.so.2"}
		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("[%d] got %q, want %q", i, got[i], want[i])
			}
		}
	})

	t.Run("no DT_NEEDED returns nil", func(t *testing.T) {
		got, err := ImportedLibraries(&fakeBinary{dynEntry: []DynEntry{{Tag: elf.DT_FLAGS_1, Val: 0}}})
		if err != nil {
			t.Fatalf("ImportedLibraries: %v", err)
		}
		if got != nil {
			t.Errorf("got %v, want nil", got)
		}
	})

	t.Run("out-of-bounds offsets are skipped", func(t *testing.T) {
		sec, entries := makeDynamic("libc.so.6")
		entries = append(entries, DynEntry{Tag: elf.DT_NEEDED, Val: 9999})
		got, err := ImportedLibraries(&fakeBinary{sections: []Section{sec}, dynEntry: entries})
		if err != nil {
			t.Fatalf("ImportedLibraries: %v", err)
		}
		if len(got) != 1 || got[0] != "libc.so.6" {
			t.Errorf("got %v, want [libc.so.6]", got)
		}
	})
}
