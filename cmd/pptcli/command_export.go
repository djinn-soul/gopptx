package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func runExportCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var inPath string
	var outPath string
	var format string
	var title string
	var driver string
	var embedImages bool
	var nav bool

	fs.StringVar(&inPath, "in", "", "input file (.pptx/.pptm/.md)")
	fs.StringVar(&outPath, "out", "", "output path (file for pdf/html, directory for png)")
	fs.StringVar(&format, "format", "pdf", "export format: pdf|html|png")
	fs.StringVar(&title, "title", "Presentation", "presentation title (used for markdown input)")
	fs.StringVar(&driver, "driver", "auto", "PDF driver: auto|native|libreoffice|powerpoint")
	fs.BoolVar(&embedImages, "embed", true, "embed images as base64 in HTML")
	fs.BoolVar(&nav, "nav", true, "include navigation shell script in HTML")

	if err := fs.Parse(args); err != nil {
		printErrorf(stderr, "invalid export arguments: %v", err)
		printExportUsage(stderr)
		return exitUsage
	}
	if len(fs.Args()) > 0 {
		printErrorf(stderr, "unexpected trailing arguments: %s", strings.Join(fs.Args(), " "))
		printExportUsage(stderr)
		return exitUsage
	}

	inPath = strings.TrimSpace(inPath)
	if inPath == "" {
		printErrorf(stderr, "export requires -in")
		printExportUsage(stderr)
		return exitUsage
	}

	switch strings.ToLower(strings.TrimSpace(format)) {
	case "pdf":
		return runPDFCommand(buildPDFArgs(inPath, outPath, title, driver), stdout, stderr)
	case "html":
		return runHTMLCommand(buildHTMLArgs(inPath, outPath, title, embedImages, nav), stdout, stderr)
	case "png":
		return runExportPNG(inPath, outPath, title, stdout, stderr)
	default:
		printErrorf(stderr, "unsupported export format %q (allowed: pdf|html|png)", format)
		printExportUsage(stderr)
		return exitUsage
	}
}

func runExportPNG(inPath, outPath, title string, stdout io.Writer, stderr io.Writer) int {
	pptxPath, cleanup, err := normalizeInputForPNG(inPath, title)
	if err != nil {
		printErrorf(stderr, "png export input error: %v", err)
		return exitIO
	}
	defer cleanup()

	outDir := strings.TrimSpace(outPath)
	if outDir == "" {
		outDir = defaultPNGOutputDir(inPath)
	}

	if err := ensureParentDir(filepath.Clean(outDir)); err != nil {
		printErrorf(stderr, "failed to prepare output directory: %v", err)
		return exitIO
	}
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		printErrorf(stderr, "failed to create output directory %q: %v", outDir, err)
		return exitIO
	}

	if err := exportPPTXToPNG(pptxPath, outDir); err != nil {
		printErrorf(stderr, "PNG export failed: %v", err)
		return exitInternal
	}

	_, _ = fmt.Fprintf(stdout, "OK: Exported PNG slides to %s\n", outDir)
	return exitOK
}

func normalizeInputForPNG(inPath, title string) (string, func(), error) {
	inputKind, err := detectPresentationInputKind(inPath)
	if err != nil {
		return "", nil, err
	}

	if inputKind == inputKindPPTX {
		return inPath, func() {}, nil
	}

	markdown, err := os.ReadFile(inPath)
	if err != nil {
		return "", nil, fmt.Errorf("read markdown %q: %w", inPath, err)
	}
	slides, err := pptx.SlidesFromMarkdown(string(markdown))
	if err != nil {
		return "", nil, fmt.Errorf("markdown parse failed: %w", err)
	}
	data, err := pptx.CreateWithSlides(strings.TrimSpace(title), slides)
	if err != nil {
		return "", nil, fmt.Errorf("pptx generation failed: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "pptcli-export-png-*")
	if err != nil {
		return "", nil, fmt.Errorf("create temp dir: %w", err)
	}
	tmpPath := filepath.Join(tmpDir, "input.pptx")
	if err := os.WriteFile(tmpPath, data, 0o600); err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("write temp pptx: %w", err)
	}
	return tmpPath, func() { _ = os.RemoveAll(tmpDir) }, nil
}

func defaultPNGOutputDir(inPath string) string {
	return defaultSiblingDirPath(inPath, "export", "_png")
}

func buildPDFArgs(inPath, outPath, title, driver string) []string {
	args := []string{"-in", strings.TrimSpace(inPath)}
	if strings.TrimSpace(outPath) != "" {
		args = append(args, "-out", strings.TrimSpace(outPath))
	}
	if strings.TrimSpace(title) != "" {
		args = append(args, "-title", strings.TrimSpace(title))
	}
	if strings.TrimSpace(driver) != "" {
		args = append(args, "-driver", strings.TrimSpace(driver))
	}
	return args
}

func buildHTMLArgs(inPath, outPath, title string, embedImages, nav bool) []string {
	args := []string{"-in", strings.TrimSpace(inPath)}
	if strings.TrimSpace(outPath) != "" {
		args = append(args, "-out", strings.TrimSpace(outPath))
	}
	if strings.TrimSpace(title) != "" {
		args = append(args, "-title", strings.TrimSpace(title))
	}
	args = append(args, "-embed", boolToFlag(embedImages), "-nav", boolToFlag(nav))
	return args
}

func boolToFlag(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func printExportUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Usage: pptcli export -in <file.pptx|file.pptm|file.md> [-out path] [-format pdf|html|png]")
}
