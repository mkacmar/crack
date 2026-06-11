package cli

import "testing"

func TestExitCode(t *testing.T) {
	tests := []struct {
		name        string
		hasFindings bool
		hasErrors   bool
		exitZero    bool
		want        int
	}{
		{name: "clean run", want: ExitSuccess},
		{name: "findings only", hasFindings: true, want: ExitFindings},
		{name: "findings with exit-zero", hasFindings: true, exitZero: true, want: ExitSuccess},
		{name: "errors only", hasErrors: true, want: ExitError},
		{name: "errors and findings", hasFindings: true, hasErrors: true, want: ExitError},
		{name: "errors with exit-zero", hasErrors: true, exitZero: true, want: ExitError},
		{name: "errors and findings with exit-zero", hasFindings: true, hasErrors: true, exitZero: true, want: ExitError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exitCode(tt.hasFindings, tt.hasErrors, tt.exitZero); got != tt.want {
				t.Errorf("exitCode(%t, %t, %t) = %d, want %d", tt.hasFindings, tt.hasErrors, tt.exitZero, got, tt.want)
			}
		})
	}
}
