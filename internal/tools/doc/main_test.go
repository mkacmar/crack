package main

import (
	"os"
	"testing"
)

const docPath = "../../../docs/rules.md"

func TestDocUpToDate(t *testing.T) {
	want, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("read %s: %v", docPath, err)
	}

	got, err := generateDoc(sortedRules())
	if err != nil {
		t.Fatalf("generate rule docs: %v", err)
	}

	if got != string(want) {
		t.Errorf("%s is out of date, run 'make doc' and commit the result", docPath)
	}
}
