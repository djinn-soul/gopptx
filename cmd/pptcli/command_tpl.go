package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx/tplx"
)

func runTplCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("tpl", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		templatePath string
		dataPath     string
		outPath      string
		strict       bool
	)

	fs.StringVar(&templatePath, "template", "", "Path to the input .pptx template")
	fs.StringVar(&dataPath, "data", "", "Path to JSON data file")
	fs.StringVar(&outPath, "out", "", "Output .pptx file path")
	fs.BoolVar(&strict, "strict", false, "Fail if a template token is missing from data")

	fs.Usage = func() {
		fmt.Fprintln(stderr, "Usage: pptcli tpl -template tpl.pptx -data data.json -out out.pptx [-strict]")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return exitOK
		}
		return exitUsage
	}

	if templatePath == "" || dataPath == "" || outPath == "" {
		fmt.Fprintln(stderr, "Error: -template, -data, and -out are required")
		fs.Usage()
		return exitUsage
	}

	dataBytes, err := os.ReadFile(dataPath)
	if err != nil {
		printErrorf(stderr, "read data file: %v", err)
		return exitInternal
	}

	var ctx tplx.Context
	if err := json.Unmarshal(dataBytes, &ctx); err != nil {
		printErrorf(stderr, "parse json data: %v", err)
		return exitInternal
	}

	opts := tplx.Options{Strict: strict}
	result, err := tplx.RenderWithOptions(templatePath, ctx, opts)
	if err != nil {
		printErrorf(stderr, "render template: %v", err)
		return exitInternal
	}

	if err := result.Save(outPath); err != nil {
		printErrorf(stderr, "save output: %v", err)
		return exitInternal
	}

	fmt.Fprintf(stdout, "Rendered %s -> %s\n", templatePath, outPath)
	return exitOK
}
