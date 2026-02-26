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

func TestCLI_VersionSubcommand(t *testing.T) {
	stdout, stderr, code := runCLI(t, "version")
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitOK, code, stdout, stderr)
	}
	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}
	if !strings.Contains(stdout, "gopptx version") {
		t.Fatalf("expected version output, got %q", stdout)
	}
}

func TestCLI_CompletionSubcommand(t *testing.T) {
	stdout, stderr, code := runCLI(t, "completion", "-shell", "bash")
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitOK, code, stdout, stderr)
	}
	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}
	if !strings.Contains(stdout, "complete -F _pptcli_complete pptcli") {
		t.Fatalf("expected bash completion script output, got %q", stdout)
	}
}

func TestCLI_CompletionSubcommand_UnsupportedShell(t *testing.T) {
	stdout, stderr, code := runCLI(t, "completion", "-shell", "fish")
	if code != exitUsage {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitUsage, code, stdout, stderr)
	}
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("expected empty stdout on unsupported shell, got %q", stdout)
	}
	if !strings.Contains(stderr, "unsupported shell") {
		t.Fatalf("expected unsupported shell error, got %q", stderr)
	}
}

func TestCLI_CreateSubcommand(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "create.pptx")
	stdout, stderr, code := runCLI(t, "create", "-out", outPath, "-title", "CLI Deck", "-slides", "2")
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitOK, code, stdout, stderr)
	}
	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}
	if !strings.Contains(stdout, "OK: wrote") {
		t.Fatalf("expected success output, got %q", stdout)
	}
	if info, err := os.Stat(outPath); err != nil {
		t.Fatalf("expected output file: %v", err)
	} else if info.Size() == 0 {
		t.Fatalf("expected non-empty pptx file")
	}
}

func TestCLI_MD2PPTSubcommand_DefaultOutput(t *testing.T) {
	tmpDir := t.TempDir()
	inPath := filepath.Join(tmpDir, "deck.md")
	markdown := "# Intro\n- one\n- two\n"
	if err := os.WriteFile(inPath, []byte(markdown), 0o600); err != nil {
		t.Fatalf("write markdown: %v", err)
	}

	stdout, stderr, code := runCLI(t, "md2ppt", "-in", inPath, "-title", "From Markdown")
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitOK, code, stdout, stderr)
	}
	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}

	outPath := filepath.Join(tmpDir, "deck.pptx")
	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("expected derived output file %s: %v", outPath, err)
	}
}

func TestCLI_InfoAndValidateSubcommands(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "info.pptx")
	_, stderr, code := runCLI(t, "create", "-out", outPath, "-title", "Info Deck", "-slides", "1")
	if code != exitOK {
		t.Fatalf("create failed: exit=%d stderr=%s", code, stderr)
	}

	infoStdout, infoStderr, infoCode := runCLI(t, "info", "-file", outPath)
	if infoCode != exitOK {
		t.Fatalf("info failed: exit=%d stderr=%s", infoCode, infoStderr)
	}
	if strings.TrimSpace(infoStderr) != "" {
		t.Fatalf("expected empty info stderr, got %q", infoStderr)
	}
	if !strings.Contains(infoStdout, "Slide count: 1") {
		t.Fatalf("expected slide count in info output, got %q", infoStdout)
	}

	validateStdout, validateStderr, validateCode := runCLI(t, "validate", "-file", outPath)
	if validateCode != exitOK {
		t.Fatalf("validate failed: exit=%d stderr=%s", validateCode, validateStderr)
	}
	if strings.TrimSpace(validateStderr) != "" {
		t.Fatalf("expected empty validate stderr, got %q", validateStderr)
	}
	if !strings.Contains(strings.ToLower(validateStdout), "validation passed") {
		t.Fatalf("expected validation success output, got %q", validateStdout)
	}
}

func TestCLI_ValidateSubcommand_InvalidZip(t *testing.T) {
	badPath := filepath.Join(t.TempDir(), "bad.pptx")
	if err := os.WriteFile(badPath, []byte("not-a-zip"), 0o600); err != nil {
		t.Fatalf("write bad file: %v", err)
	}

	stdout, stderr, code := runCLI(t, "validate", "-file", badPath)
	if code != exitValidate {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitValidate, code, stdout, stderr)
	}
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("expected empty stdout on validation failure, got %q", stdout)
	}
	if !strings.Contains(stderr, "not a valid ZIP archive") {
		t.Fatalf("expected zip validation error, got %q", stderr)
	}
}

func TestCLI_PDFSubcommand_InvalidDriver(t *testing.T) {
	stdout, stderr, code := runCLI(t, "pdf", "-in", "deck.md", "-driver", "chromedp")
	if code != exitUsage {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitUsage, code, stdout, stderr)
	}
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "invalid PDF driver") {
		t.Fatalf("expected invalid driver error, got %q", stderr)
	}
}

func runCLI(t *testing.T, args ...string) (string, string, int) {
	t.Helper()

	cmd := exec.Command(cliBinary(t), args...)

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err == nil {
		return outBuf.String(), errBuf.String(), exitOK
	}

	var ee *exec.ExitError
	if !os.IsNotExist(err) && strings.Contains(err.Error(), "executable file not found") {
		t.Fatalf("failed to run CLI binary: %v", err)
	}
	if errors.As(err, &ee) {
		return outBuf.String(), errBuf.String(), ee.ExitCode()
	}
	t.Fatalf("unexpected run error: %v", err)
	return "", "", exitInternal
}

func TestMergeCommand(t *testing.T) {
	tmpDir := t.TempDir()

	a := filepath.Join(tmpDir, "a.pptx")
	runCLI(t, "create", "-out", a, "-title", "Presentation A")

	b := filepath.Join(tmpDir, "b.pptx")
	runCLI(t, "create", "-out", b, "-title", "Presentation B")

	merged := filepath.Join(tmpDir, "merged.pptx")
	out, stderr, code := runCLI(t, "merge", "-out", merged, a, b)
	if code != exitOK {
		t.Fatalf("merge failed (%d): %s\n%s", code, out, stderr)
	}

	if !strings.Contains(out, "Successfully merged 2 files") {
		t.Fatalf("unexpected success output: %s", out)
	}

	infoOut, _, _ := runCLI(t, "info", "-file", merged)
	if !strings.Contains(infoOut, "Slide count: 2") {
		t.Fatalf("expected 2 slides in merged file, got output: %s", infoOut)
	}
}

func cliBinary(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	name := "pptcli-test-bin"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	cliBinPath := filepath.Join(tmpDir, name)

	build := exec.Command("go", "build", "-o", cliBinPath, ".")
	build.Dir = "."
	var buildOut bytes.Buffer
	build.Stdout = &buildOut
	build.Stderr = &buildOut
	if err := build.Run(); err != nil {
		t.Fatalf("failed to build CLI binary: %v", fmt.Errorf("build failed: %w: %s", err, buildOut.String()))
	}

	return cliBinPath
}
