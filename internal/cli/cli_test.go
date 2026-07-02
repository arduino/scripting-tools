package cli

import (
	"bytes"
	"testing"
)

func TestRunVersion(t *testing.T) {
	var out bytes.Buffer
	var err bytes.Buffer

	exitCode := Run([]string{"version"}, &out, &err)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if out.Len() == 0 {
		t.Fatal("expected version output")
	}
	if err.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", err.String())
	}
}
