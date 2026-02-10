package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func runMergeCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var output string
	fs.StringVar(&output, "out", "", "Output PPTX file path")

	if err := fs.Parse(args); err != nil {
		return exitUsage
	}

	if output == "" {
		printErrorf(stderr, "-out is required")
		return exitUsage
	}

	inputs := fs.Args()
	if len(inputs) < 2 {
		printErrorf(stderr, "at least two input PPTX files are required for merge")
		return exitUsage
	}

	editor, err := pptx.OpenPresentationEditor(inputs[0])
	if err != nil {
		printErrorf(stderr, "failed to open first input %q: %v", inputs[0], err)
		return exitInternal
	}

	for _, input := range inputs[1:] {
		if err := editor.MergeFromFile(input); err != nil {
			printErrorf(stderr, "failed to merge %q: %v", input, err)
			return exitInternal
		}
	}

	if err := editor.Save(output); err != nil {
		printErrorf(stderr, "failed to save merged presentation to %q: %v", output, err)
		return exitInternal
	}

	_, _ = fmt.Fprintf(stdout, "Successfully merged %d files into %q\n", len(inputs), output)
	return exitOK
}
