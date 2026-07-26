package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

func TestCLI_ExportSubcommand_HTMLFromMarkdown(t *testing.T) {
	tmpDir := t.TempDir()
	inPath := filepath.Join(tmpDir, "deck.md")
	if err := os.WriteFile(inPath, []byte("# Intro\n- one\n- two\n"), 0o600); err != nil {
		t.Fatalf("write markdown: %v", err)
	}
	outPath := filepath.Join(tmpDir, "deck.html")

	stdout, stderr, code := runCLI(
		t,
		"export",
		"-in",
		inPath,
		"-format",
		"html",
		"-out",
		outPath,
		"-title",
		"Export HTML",
	)
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitOK, code, stdout, stderr)
	}
	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("expected html output file: %v", err)
	}
	if !strings.Contains(strings.ToLower(string(data)), "<!doctype html>") {
		t.Fatalf("expected html output, got %q", string(data))
	}
}

func TestCLI_ExportSubcommand_RejectsUnknownFormat(t *testing.T) {
	stdout, stderr, code := runCLI(t, "export", "-in", "deck.md", "-format", "svg")
	if code != exitUsage {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitUsage, code, stdout, stderr)
	}
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "unsupported export format") {
		t.Fatalf("expected unsupported format error, got %q", stderr)
	}
}

func TestCLI_URLFetchSubcommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>URLFetch</title></head>
<body>
  <main>
    <h1>Hello URL</h1>
    <p>This paragraph has enough content to be parsed into deck bullets, and it intentionally contains more than one hundred characters so the urlfetch parser accepts this document as primary page content.</p>
    <p>Second paragraph to ensure there is meaningful body text for section extraction and slide generation.</p>
  </main>
</body>
</html>`))
	}))
	t.Cleanup(server.Close)

	outPath := filepath.Join(t.TempDir(), "urlfetch.pptx")
	stdout, stderr, code := runCLIWithEnv(
		t,
		[]string{urlfetchAllowPrivateHostsEnv + "=1"},
		"urlfetch",
		"-url",
		server.URL,
		"-out",
		outPath,
		"-title",
		"URL Deck",
	)
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

func TestCLI_URLFetchSubcommand_RejectsPrivateHostsByDefault(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("<html><body>blocked</body></html>"))
	}))
	t.Cleanup(server.Close)

	stdout, stderr, code := runCLI(t, "urlfetch", "-url", server.URL)
	if code != exitInternal {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitInternal, code, stdout, stderr)
	}
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "SSRF guard") {
		t.Fatalf("expected SSRF guard error, got %q", stderr)
	}
}

func TestCLI_URLFetchSubcommand_RequiresURL(t *testing.T) {
	stdout, stderr, code := runCLI(t, "urlfetch")
	if code != exitUsage {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitUsage, code, stdout, stderr)
	}
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "urlfetch requires -url") {
		t.Fatalf("expected missing-url error, got %q", stderr)
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
