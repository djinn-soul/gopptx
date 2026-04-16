package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func TestDetectPresentationInputKind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		want    presentationInputKind
		wantErr string
	}{
		{name: "markdown lower", path: "deck.md", want: inputKindMarkdown},
		{name: "markdown upper", path: "deck.MD", want: inputKindMarkdown},
		{name: "pptx upper", path: "deck.PPTX", want: inputKindPPTX},
		{name: "pptm lower", path: "deck.pptm", want: inputKindPPTX},
		{name: "unsupported", path: "deck.txt", wantErr: `unsupported input file extension ".txt"`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := detectPresentationInputKind(tc.path)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("detectPresentationInputKind(%q): %v", tc.path, err)
			}
			if got != tc.want {
				t.Fatalf("detectPresentationInputKind(%q) = %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}

func TestDefaultSiblingPathsHandleUppercaseExtensions(t *testing.T) {
	t.Parallel()

	if got := defaultSiblingFilePath(
		filepath.Join("tmp", "DECK.MD"),
		"export",
		".html",
	); got != filepath.Join("tmp", "DECK.html") {
		t.Fatalf("defaultSiblingFilePath html = %q", got)
	}
	if got := defaultSiblingFilePath(
		filepath.Join("tmp", "DECK.PPTM"),
		"export",
		".pdf",
	); got != filepath.Join("tmp", "DECK.pdf") {
		t.Fatalf("defaultSiblingFilePath pdf = %q", got)
	}
	if got := defaultSiblingDirPath(
		filepath.Join("tmp", "DECK.PPTX"),
		"export",
		"_png",
	); got != filepath.Join("tmp", "DECK_png") {
		t.Fatalf("defaultSiblingDirPath png = %q", got)
	}
}

func TestNormalizeInputForPNGRejectsUnsupportedExtension(t *testing.T) {
	t.Parallel()

	_, _, err := normalizeInputForPNG("deck.txt", "Ignored")
	if err == nil || !strings.Contains(err.Error(), `unsupported input file extension ".txt"`) {
		t.Fatalf("expected unsupported extension error, got %v", err)
	}
}

func TestRunHTMLCommandExportsPPTXInput(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	inPath := filepath.Join(tmpDir, "deck.PPTX")
	data, err := pptx.Create("Deck Title", 1)
	if err != nil {
		t.Fatalf("create pptx: %v", err)
	}
	if err := os.WriteFile(inPath, data, 0o600); err != nil {
		t.Fatalf("write pptx: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runHTMLCommand([]string{"-in", inPath}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d\nstdout=%s\nstderr=%s", exitOK, code, stdout.String(), stderr.String())
	}
	if strings.TrimSpace(stderr.String()) != "" {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}

	outPath := filepath.Join(tmpDir, "deck.html")
	htmlData, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read html output: %v", err)
	}
	htmlText := string(htmlData)
	if !strings.Contains(strings.ToLower(htmlText), "<!doctype html>") {
		t.Fatalf("expected HTML document, got %q", htmlText)
	}
	if !strings.Contains(htmlText, "<title>Deck Title</title>") {
		t.Fatalf("expected PPTX metadata title in html, got %q", htmlText)
	}
}

func TestRunPDFCommandRejectsUnsupportedExtension(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runPDFCommand([]string{"-in", "deck.txt"}, &stdout, &stderr)
	if code != exitUsage {
		t.Fatalf("expected exit %d, got %d", exitUsage, code)
	}
	if !strings.Contains(stderr.String(), `unsupported input file extension ".txt"`) {
		t.Fatalf("expected unsupported extension error, got %q", stderr.String())
	}
}
