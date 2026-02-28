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

func runHTMLCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("html", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var inPath string
	var outPath string
	var title string
	var embedImages bool
	var nav bool

	fs.StringVar(&inPath, "in", "", "input markdown file")
	fs.StringVar(&outPath, "out", "", "output HTML file (default: built from input name)")
	fs.StringVar(&title, "title", "Presentation", "presentation title")
	fs.BoolVar(&embedImages, "embed", true, "embed images as base64 in HTML")
	fs.BoolVar(&nav, "nav", true, "include navigation shell script in HTML")

	if err := fs.Parse(args); err != nil {
		printErrorf(stderr, "invalid html arguments: %v", err)
		printHTMLUsage(stderr)
		return exitUsage
	}

	inPath = strings.TrimSpace(inPath)
	if inPath == "" {
		printErrorf(stderr, "html requires -in")
		printHTMLUsage(stderr)
		return exitUsage
	}

	if strings.TrimSpace(outPath) == "" {
		outPath = strings.TrimSuffix(inPath, ".md") + ".html"
	}

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

	opts := export.DefaultHTMLOptions()
	opts.EmbedImages = embedImages
	opts.IncludeNavigation = nav

	htmlContent := export.HTMLWithOptions(strings.TrimSpace(title), slides, opts)
	if err := os.WriteFile(outPath, []byte(htmlContent), 0o644); err != nil {
		printErrorf(stderr, "failed to write HTML to %q: %v", outPath, err)
		return exitIO
	}

	_, _ = fmt.Fprintf(stdout, "OK: Exported to %s\n", outPath)
	return exitOK
}

func printHTMLUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Usage: pptcli html -in deck.md [-out deck.html] [-title TITLE] [-embed=true] [-nav=true]")
}
