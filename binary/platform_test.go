package binary

import "testing"

func TestParseISA(t *testing.T) {
	tests := []struct {
		in      string
		want    ISA
		wantErr bool
	}{
		{"v8.3", ISA{Major: 8, Minor: 3}, false},
		{"8.3", ISA{Major: 8, Minor: 3}, false},
		{"v8", ISA{Major: 8}, false},
		{"1", ISA{Major: 1}, false},
		{"v0.0", ISA{}, false},
		{"v8.3.1", ISA{}, true},
		{"vabc", ISA{}, true},
		{"v8.x", ISA{}, true},
		{"", ISA{}, true},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, err := ParseISA(tc.in)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ParseISA(%q) err = %v, wantErr = %v", tc.in, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("ParseISA(%q) = %+v, want %+v", tc.in, got, tc.want)
			}
		})
	}
}

func TestISAString(t *testing.T) {
	tests := []struct {
		isa  ISA
		want string
	}{
		{ISA{Major: 8, Minor: 3}, "v8.3"},
		{ISA{Major: 8}, "v8"},
		{ISA{Major: 1, Minor: 0}, "v1"},
		{ISA{}, "v0"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			if got := tc.isa.String(); got != tc.want {
				t.Errorf("ISA%+v.String() = %q, want %q", tc.isa, got, tc.want)
			}
		})
	}
}

func TestISAIsAtLeast(t *testing.T) {
	tests := []struct {
		name     string
		have     ISA
		required ISA
		want     bool
	}{
		{"equal", ISA{Major: 8, Minor: 3}, ISA{Major: 8, Minor: 3}, true},
		{"higher minor", ISA{Major: 8, Minor: 5}, ISA{Major: 8, Minor: 3}, true},
		{"lower minor", ISA{Major: 8, Minor: 1}, ISA{Major: 8, Minor: 3}, false},
		{"higher major beats lower minor", ISA{Major: 9}, ISA{Major: 8, Minor: 5}, true},
		{"lower major beats higher minor", ISA{Major: 8, Minor: 9}, ISA{Major: 9}, false},
		{"zero required matches anything", ISA{Major: 1}, ISA{}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.have.IsAtLeast(tc.required); got != tc.want {
				t.Errorf("%+v.IsAtLeast(%+v) = %v, want %v", tc.have, tc.required, got, tc.want)
			}
		})
	}
}

func TestParseISARoundTrip(t *testing.T) {
	cases := []ISA{
		{Major: 8, Minor: 3},
		{Major: 8, Minor: 5},
		{Major: 1},
	}
	for _, want := range cases {
		t.Run(want.String(), func(t *testing.T) {
			got, err := ParseISA(want.String())
			if err != nil {
				t.Fatalf("ParseISA(%q): %v", want.String(), err)
			}
			if got != want {
				t.Errorf("round trip = %+v, want %+v", got, want)
			}
		})
	}
}
