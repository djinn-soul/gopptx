package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

const minMergeInputs = 2

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
	if len(inputs) < minMergeInputs {
		printErrorf(stderr, "at least two input PPTX files are required for merge")
		return exitUsage
	}

	editor, err := pptx.OpenPresentationEditor(inputs[0])
	if err != nil {
		printErrorf(stderr, "failed to open first input %q: %v", inputs[0], err)
		return exitInternal
	}

	for _, input := range inputs[1:] {
		if mergeErr := editor.MergeFromFile(input); mergeErr != nil {
			printErrorf(stderr, "failed to merge %q: %v", input, mergeErr)
			return exitInternal
		}
	}

	if saveErr := editor.Save(output); saveErr != nil {
		printErrorf(stderr, "failed to save merged presentation to %q: %v", output, saveErr)
		return exitInternal
	}

	_, _ = fmt.Fprintf(stdout, "Successfully merged %d files into %q\n", len(inputs), output)
	return exitOK
}
