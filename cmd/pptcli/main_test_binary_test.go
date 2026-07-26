package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const cliBinPathEnv = "GOPPTX_TEST_CLI_BIN"

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "pptcli-test-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create CLI test directory: %v\n", err)
		os.Exit(1)
	}

	name := "pptcli-test-bin"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	cliBinPath := filepath.Join(tmpDir, name)

	build := exec.Command("go", "build", "-o", cliBinPath, ".")
	build.Dir = "."
	if output, buildErr := build.CombinedOutput(); buildErr != nil {
		fmt.Fprintf(os.Stderr, "build CLI test binary: %v\n%s", buildErr, output)
		_ = os.RemoveAll(tmpDir)
		os.Exit(1)
	}
	if err := os.Setenv(cliBinPathEnv, cliBinPath); err != nil {
		fmt.Fprintf(os.Stderr, "configure CLI test binary: %v\n", err)
		_ = os.RemoveAll(tmpDir)
		os.Exit(1)
	}

	code := m.Run()
	if removeErr := os.RemoveAll(tmpDir); removeErr != nil && code == 0 {
		fmt.Fprintf(os.Stderr, "remove CLI test directory: %v\n", removeErr)
		code = 1
	}
	os.Exit(code)
}

func runCLI(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	return runCLIWithEnv(t, nil, args...)
}

func runCLIWithEnv(t *testing.T, env []string, args ...string) (string, string, int) {
	t.Helper()

	cliBinPath := os.Getenv(cliBinPathEnv)
	if cliBinPath == "" {
		t.Fatal("CLI test binary path is not configured")
	}
	cmd := exec.Command(cliBinPath, args...)
	cmd.Env = append(os.Environ(), env...)

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err == nil {
		return outBuf.String(), errBuf.String(), exitOK
	}

	var exitErr *exec.ExitError
	if !os.IsNotExist(err) && strings.Contains(err.Error(), "executable file not found") {
		t.Fatalf("failed to run CLI binary: %v", err)
	}
	if errors.As(err, &exitErr) {
		return outBuf.String(), errBuf.String(), exitErr.ExitCode()
	}
	t.Fatalf("unexpected run error: %v", err)
	return "", "", exitInternal
}
