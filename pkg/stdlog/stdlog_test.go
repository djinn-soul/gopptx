package stdlog

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const stdlogModeEnv = "STDLOG_TEST_MODE"

func TestPrintln(t *testing.T) {
	if runStdlogModeFromEnv() {
		return
	}
	output, err := runStdlogSubprocess(t, "println")
	if err != nil {
		t.Fatalf("println subprocess failed: %v, output=%q", err, output)
	}
	if !strings.Contains(output, "hello 7") {
		t.Fatalf("unexpected Println output: %q", output)
	}
}

func TestPrintf(t *testing.T) {
	if runStdlogModeFromEnv() {
		return
	}
	output, err := runStdlogSubprocess(t, "printf")
	if err != nil {
		t.Fatalf("printf subprocess failed: %v, output=%q", err, output)
	}
	if !strings.Contains(output, "value=42") {
		t.Fatalf("unexpected Printf output: %q", output)
	}
}

func TestFatal(t *testing.T) {
	if runStdlogModeFromEnv() {
		return
	}
	output, err := runStdlogSubprocess(t, "fatal")
	assertExitCodeOne(t, err, output)
	if !strings.Contains(output, "boom") {
		t.Fatalf("expected fatal output to contain boom, got %q", output)
	}
}

func TestFatalf(t *testing.T) {
	if runStdlogModeFromEnv() {
		return
	}
	output, err := runStdlogSubprocess(t, "fatalf")
	assertExitCodeOne(t, err, output)
	if !strings.Contains(output, "code=9") {
		t.Fatalf("expected fatalf output to contain code=9, got %q", output)
	}
}

func runStdlogModeFromEnv() bool {
	switch os.Getenv(stdlogModeEnv) {
	case "println":
		Println("hello", 7)
		return true
	case "printf":
		Printf("value=%d", 42)
		return true
	case "fatal":
		Fatal("boom")
		return true
	case "fatalf":
		Fatalf("code=%d", 9)
		return true
	default:
		return false
	}
}

func runStdlogSubprocess(t *testing.T, mode string) (string, error) {
	t.Helper()

	cmd := exec.Command(os.Args[0], "-test.run", "^"+t.Name()+"$")
	cmd.Env = append(os.Environ(), stdlogModeEnv+"="+mode)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func assertExitCodeOne(t *testing.T, err error, output string) {
	t.Helper()

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		t.Fatalf("expected exit code 1, err=%v, output=%q", err, output)
	}
}
