package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

func TestCLI_ValidateSubcommand_JSON(t *testing.T) {
	badPath := filepath.Join(t.TempDir(), "bad.pptx")
	if err := os.WriteFile(badPath, []byte("not-a-zip"), 0o600); err != nil {
		t.Fatalf("write bad file: %v", err)
	}

	stdout, stderr, code := runCLI(t, "validate", "-file", badPath, "-format", "json")
	if code != exitValidate {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitValidate, code, stdout, stderr)
	}
	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("expected empty stderr on json mode, got %q", stderr)
	}
	if !strings.Contains(stdout, `"error":`) || !strings.Contains(stdout, "not a valid ZIP archive") {
		t.Fatalf("expected JSON output containing zip error, got %q", stdout)
	}
}

func TestCLI_RepairSubcommand_DryRun_JSON(t *testing.T) {
	badPath := filepath.Join(t.TempDir(), "bad.pptx")
	if err := os.WriteFile(badPath, []byte("not-a-zip"), 0o600); err != nil {
		t.Fatalf("write bad file: %v", err)
	}

	stdout, stderr, code := runCLI(t, "repair", "-file", badPath, "-dry-run", "-format", "json")
	// For "not a zip", repair fails cleanly instead of returning repaired issues.
	if code != exitIO {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitIO, code, stdout, stderr)
	}
	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("expected empty stderr on json mode, got %q", stderr)
	}
	if !strings.Contains(stdout, `"error":`) {
		t.Fatalf("expected JSON error output, got %q", stdout)
	}
}

func TestCLI_RepairSubcommand_JSON_WritesOutputFile(t *testing.T) {
	inPath := filepath.Join(t.TempDir(), "in.pptx")
	_, createErr, createCode := runCLI(t, "create", "-out", inPath, "-title", "Repair JSON", "-slides", "1")
	if createCode != exitOK {
		t.Fatalf("create failed: exit=%d stderr=%s", createCode, createErr)
	}
	outPath := filepath.Join(t.TempDir(), "out.pptx")
	stdout, stderr, code := runCLI(t, "repair", "-file", inPath, "-out", outPath, "-format", "json")
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitOK, code, stdout, stderr)
	}
	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("expected empty stderr on json mode, got %q", stderr)
	}
	if !strings.Contains(stdout, `"issues_repaired":`) || !strings.Contains(stdout, `"issues_unrepaired":`) {
		t.Fatalf("expected JSON output, got %q", stdout)
	}
	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("expected repaired output file: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("expected non-empty repaired output")
	}
}
