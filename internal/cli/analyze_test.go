package cli

import (
	"os"
	"testing"
)

func TestRunAnalyzeParseExitCodes(t *testing.T) {
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer devNull.Close()

	stderr := os.Stderr
	os.Stderr = devNull
	t.Cleanup(func() { os.Stderr = stderr })

	tests := []struct {
		name string
		args []string
		want int
	}{
		{name: "unknown flag is an error, not a finding", args: []string{"--nonsuch"}, want: ExitError},
		{name: "help is a success", args: []string{"-h"}, want: ExitSuccess},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New().runAnalyze("crack", tt.args); got != tt.want {
				t.Errorf("runAnalyze(%q) = %d, want %d", tt.args, got, tt.want)
			}
		})
	}
}
