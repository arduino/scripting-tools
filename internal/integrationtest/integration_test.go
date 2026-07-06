package integrationtest

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type runResult struct {
	exitCode int
	stdout   string
	stderr   string
}

func runScriptingTools(t *testing.T, args ...string) runResult {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to determine current file path")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "../.."))

	cmdArgs := append([]string{"run", "."}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = repoRoot

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to execute scripting-tools: %v", err)
		}
	}

	return runResult{
		exitCode: exitCode,
		stdout:   stdout.String(),
		stderr:   stderr.String(),
	}
}

func TestNoArgumentsShowsUsage(t *testing.T) {
	t.Parallel()

	res := runScriptingTools(t)

	if res.exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", res.exitCode)
	}
	if !strings.Contains(res.stdout, "Usage:") {
		t.Fatalf("expected usage text in stdout, got %q", res.stdout)
	}
	if res.stderr != "" {
		t.Fatalf("expected empty stderr, got %q", res.stderr)
	}
}

func TestUnknownCommandReturnsError(t *testing.T) {
	t.Parallel()

	res := runScriptingTools(t, "does-not-exist")

	if res.exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", res.exitCode)
	}
	if !strings.Contains(res.stderr, `error: unknown command "does-not-exist"`) {
		t.Fatalf("expected unknown command error in stderr, got %q", res.stderr)
	}
	if !strings.Contains(res.stderr, "Commands:") {
		t.Fatalf("expected usage text in stderr, got %q", res.stderr)
	}
}

func TestConditionalScriptRunsCommandWhenConditionTrue(t *testing.T) {
	t.Parallel()

	res := runScriptingTools(t, "if", "eq", "a", "a", "echo", "OK")

	if res.exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstderr=%q", res.exitCode, res.stderr)
	}
	if !strings.Contains(res.stdout, "Running command: echo OK") {
		t.Fatalf("expected command execution trace in stdout, got %q", res.stdout)
	}
	if !strings.Contains(res.stdout, "OK") {
		t.Fatalf("expected command output in stdout, got %q", res.stdout)
	}
	if res.stderr != "" {
		t.Fatalf("expected empty stderr, got %q", res.stderr)
	}
}

func TestConditionalScriptSkipsCommandWhenConditionFalse(t *testing.T) {
	t.Parallel()

	res := runScriptingTools(t, "if", "eq", "a", "b", "echo", "SHOULD_NOT_RUN")

	if res.exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstderr=%q", res.exitCode, res.stderr)
	}
	if strings.Contains(res.stdout, "SHOULD_NOT_RUN") {
		t.Fatalf("expected command not to run, got stdout %q", res.stdout)
	}
	if res.stderr != "" {
		t.Fatalf("expected empty stderr, got %q", res.stderr)
	}
}

func TestTwoPhaseUploadSkipsLoaderWhenUpToDate(t *testing.T) {
	t.Parallel()

	res := runScriptingTools(
		t,
		"two-phase-upload",
		"::loader-installed", "1.2.3",
		"::loader-required", "1.2.3",
		"::upload-loader", "echo", "LOADER",
		"::upload", "echo", "FIRMWARE",
	)

	if res.exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstderr=%q", res.exitCode, res.stderr)
	}
	if !strings.Contains(res.stdout, "Installed loader is up-to-date.") {
		t.Fatalf("expected up-to-date message in stdout, got %q", res.stdout)
	}
	if !strings.Contains(res.stdout, "Running command: echo FIRMWARE") {
		t.Fatalf("expected firmware upload command to run, got %q", res.stdout)
	}
	if strings.Contains(res.stdout, "Running command: echo LOADER") {
		t.Fatalf("expected loader upload command not to run, got %q", res.stdout)
	}
	if res.stderr != "" {
		t.Fatalf("expected empty stderr, got %q", res.stderr)
	}
}

func TestMultipleScriptsWithSeparator(t *testing.T) {
	t.Parallel()

	res := runScriptingTools(
		t,
		"if", "eq", "1", "2", "echo", "YES",
		"::",
		"if", "eq", "2", "2", "echo", "NO",
	)

	if res.exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstderr=%q", res.exitCode, res.stderr)
	}
	if strings.Contains(res.stdout, "YES") {
		t.Fatalf("expected first command not to run, got stdout %q", res.stdout)
	}
	if !strings.Contains(res.stdout, "Running command: echo NO") {
		t.Fatalf("expected second command to run, got stdout %q", res.stdout)
	}
	if !strings.Contains(res.stdout, "NO") {
		t.Fatalf("expected second command output, got stdout %q", res.stdout)
	}
	if res.stderr != "" {
		t.Fatalf("expected empty stderr, got %q", res.stderr)
	}
}

func TestMultipleScriptsWithCustomSeparator(t *testing.T) {
	t.Parallel()

	res := runScriptingTools(
		t,
		"--sep", "AA",
		"if", "eq", "1", "2", "echo", "YES",
		"AA",
		"if", "eq", "2", "2", "echo", "NO",
	)

	if res.exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstderr=%q", res.exitCode, res.stderr)
	}
	if strings.Contains(res.stdout, "YES") {
		t.Fatalf("expected first command not to run, got stdout %q", res.stdout)
	}
	if !strings.Contains(res.stdout, "Running command: echo NO") {
		t.Fatalf("expected second command to run, got stdout %q", res.stdout)
	}
	if !strings.Contains(res.stdout, "NO") {
		t.Fatalf("expected second command output, got stdout %q", res.stdout)
	}
	if res.stderr != "" {
		t.Fatalf("expected empty stderr, got %q", res.stderr)
	}
}
