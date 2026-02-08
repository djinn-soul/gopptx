package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/djinn09/gopptx/pkg/pptx"
)

func runCreateCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("create", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var outPath string
	var title string
	var slideCount int
	fs.StringVar(&outPath, "out", "output.pptx", "output PPTX file")
	fs.StringVar(&title, "title", "Presentation", "presentation title")
	fs.IntVar(&slideCount, "slides", 1, "number of slides")

	if err := fs.Parse(args); err != nil {
		printErrorf(stderr, "invalid create arguments: %v", err)
		printCreateUsage(stderr)
		return exitUsage
	}
	if len(fs.Args()) > 0 {
		printErrorf(stderr, "unexpected trailing arguments: %s", strings.Join(fs.Args(), " "))
		printCreateUsage(stderr)
		return exitUsage
	}

	data, err := pptx.Create(title, slideCount)
	if err != nil {
		printErrorf(stderr, "create failed: %v", err)
		return exitUsage
	}
	if err := writeOutputFile(outPath, data); err != nil {
		printErrorf(stderr, "failed to write %q: %v", outPath, err)
		return exitIO
	}

	_, _ = fmt.Fprintf(stdout, "OK: wrote %s (%d slide(s))\n", strings.TrimSpace(outPath), slideCount)
	return exitOK
}

func printCreateUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Usage: pptcli create -out file.pptx [-title TITLE] [-slides N]")
}



