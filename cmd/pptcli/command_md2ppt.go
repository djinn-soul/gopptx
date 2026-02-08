package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/djinn09/gopptx/pkg/pptx"
)

func runMD2PPTCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("md2ppt", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var inPath string
	var outPath string
	var title string
	fs.StringVar(&inPath, "in", "", "input markdown file")
	fs.StringVar(&outPath, "out", "", "output PPTX file (default: <input>.pptx)")
	fs.StringVar(&title, "title", "Presentation from Markdown", "presentation title")

	if err := fs.Parse(args); err != nil {
		printErrorf(stderr, "invalid md2ppt arguments: %v", err)
		printMD2PPTUsage(stderr)
		return exitUsage
	}
	if len(fs.Args()) > 0 {
		printErrorf(stderr, "unexpected trailing arguments: %s", strings.Join(fs.Args(), " "))
		printMD2PPTUsage(stderr)
		return exitUsage
	}

	inPath = strings.TrimSpace(inPath)
	if inPath == "" {
		printErrorf(stderr, "md2ppt requires -in")
		printMD2PPTUsage(stderr)
		return exitUsage
	}
	if strings.TrimSpace(outPath) == "" {
		outPath = defaultOutputPathFromMarkdown(inPath)
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

	data, err := pptx.CreateWithSlides(strings.TrimSpace(title), slides)
	if err != nil {
		printErrorf(stderr, "pptx generation failed: %v", err)
		return exitInternal
	}
	if err := writeOutputFile(outPath, data); err != nil {
		printErrorf(stderr, "failed to write %q: %v", outPath, err)
		return exitIO
	}

	_, _ = fmt.Fprintf(stdout, "OK: wrote %s from %s (%d slide(s))\n", strings.TrimSpace(outPath), inPath, len(slides))
	return exitOK
}

func printMD2PPTUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Usage: pptcli md2ppt -in deck.md [-out file.pptx] [-title TITLE]")
}



