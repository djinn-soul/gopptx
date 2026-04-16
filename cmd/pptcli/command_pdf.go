package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
)

func runPDFCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("pdf", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var inPath string
	var outPath string
	var title string
	var driver string

	fs.StringVar(&inPath, "in", "", "input file (.pptx or .md)")
	fs.StringVar(&outPath, "out", "", "output PDF file (default: built from input name)")
	fs.StringVar(&title, "title", "Presentation", "presentation title")
	fs.StringVar(&driver, "driver", "auto", "PDF driver: auto|native|libreoffice|powerpoint")

	if err := fs.Parse(args); err != nil {
		printErrorf(stderr, "invalid pdf arguments: %v", err)
		printPDFUsage(stderr)
		return exitUsage
	}

	inPath = strings.TrimSpace(inPath)
	if inPath == "" {
		printErrorf(stderr, "pdf requires -in")
		printPDFUsage(stderr)
		return exitUsage
	}

	if strings.TrimSpace(outPath) == "" {
		outPath = defaultSiblingFilePath(inPath, "export", ".pdf")
	}

	pdfDriver, err := export.ParsePDFDriver(driver)
	if err != nil {
		printErrorf(stderr, "%v", err)
		printPDFUsage(stderr)
		return exitUsage
	}
	opts := export.PDFOptions{Driver: pdfDriver}

	inputKind, err := detectPresentationInputKind(inPath)
	if err != nil {
		printErrorf(stderr, "%v", err)
		printPDFUsage(stderr)
		return exitUsage
	}

	if inputKind == inputKindPPTX {
		return pdfFromPPTXFile(inPath, outPath, opts, stdout, stderr)
	}
	return pdfFromMarkdown(inPath, outPath, title, opts, stdout, stderr)
}

func pdfFromPPTXFile(inPath, outPath string, opts export.PDFOptions, stdout, stderr io.Writer) int {
	if err := export.PDFFromFileWithOptions(inPath, outPath, opts); err != nil {
		printErrorf(stderr, "PDF generation from PPTX failed: %v", err)
		return exitInternal
	}
	_, _ = fmt.Fprintf(stdout, "OK: Exported to %s\n", outPath)
	return exitOK
}

func pdfFromMarkdown(inPath, outPath, title string, opts export.PDFOptions, stdout, stderr io.Writer) int {
	markdown, err := os.ReadFile(inPath)
	if err != nil {
		printErrorf(stderr, "failed to read markdown %q: %v", inPath, err)
		return exitIO
	}

	slides, err := pptx.SlidesFromMarkdown(string(markdown))
	if err != nil {
		printErrorf(stderr, "markdown parse failed: %v", err)
		return exitParse
	}

	if err := export.PDFWithOptions(strings.TrimSpace(title), slides, outPath, opts); err != nil {
		printErrorf(stderr, "PDF generation failed: %v", err)
		return exitInternal
	}

	_, _ = fmt.Fprintf(stdout, "OK: Exported to %s\n", outPath)
	return exitOK
}

func printPDFUsage(w io.Writer) {
	_, _ = fmt.Fprintln(
		w,
		"Usage: pptcli pdf -in <file.pptx|file.pptm|file.md> [-out file.pdf] [-title TITLE] [-driver auto|native|libreoffice|powerpoint]",
	)
}
