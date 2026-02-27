package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

// runValidateCommand validates a PPTX file structure.
func runValidateCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var filePath string
	var format string
	fs.StringVar(&filePath, "file", "", "PPTX file path")
	fs.StringVar(&format, "format", "text", "Output format: text or json")

	if err := fs.Parse(args); err != nil {
		printErrorf(stderr, "invalid validate arguments: %v", err)
		printValidateUsage(stderr)
		return exitUsage
	}
	if len(fs.Args()) > 0 {
		printErrorf(stderr, "unexpected trailing arguments: %s", strings.Join(fs.Args(), " "))
		printValidateUsage(stderr)
		return exitUsage
	}
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		printErrorf(stderr, "validate requires -file")
		printValidateUsage(stderr)
		return exitUsage
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if format == "json" {
			outputJSONError(stdout, fmt.Sprintf("failed to read file: %v", err))
			return exitIO
		}
		printErrorf(stderr, "failed to read file: %v", err)
		return exitIO
	}

	issues, err := pptx.Validate(data)
	if err != nil {
		if format == "json" {
			outputJSONError(stdout, fmt.Sprintf("validation failed: %v", err))
			return exitValidate
		}
		printErrorf(stderr, "validation failed: %v", err)
		return exitValidate
	}

	if format == "json" {
		out, err := json.MarshalIndent(issues, "", "  ")
		if err != nil {
			outputJSONError(stdout, fmt.Sprintf("failed to marshal JSON: %v", err))
			return exitValidate
		}
		_, _ = fmt.Fprintln(stdout, string(out))
		if len(issues) > 0 {
			return exitValidate
		}
		return exitOK
	}

	if len(issues) > 0 {
		for _, issue := range issues {
			_, _ = fmt.Fprintln(stderr, issue.String())
		}
		return exitValidate
	}

	_, _ = fmt.Fprintln(stdout, "OK: validation passed")
	return exitOK
}

func outputJSONError(w io.Writer, msg string) {
	errObj := map[string]string{"error": msg}
	out, _ := json.MarshalIndent(errObj, "", "  ")
	_, _ = fmt.Fprintln(w, string(out))
}

func printValidateUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Usage: pptcli validate -file file.pptx [-format text|json]")
}
