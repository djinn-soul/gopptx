package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

// runRepairCommand attempts to fix structural issues in a PPTX file.
func runRepairCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("repair", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var filePath string
	var outPath string
	var dryRun bool
	var format string
	jsonMode := false
	fs.StringVar(&filePath, "file", "", "Input PPTX file path")
	fs.StringVar(&outPath, "out", "", "Output (repaired) PPTX file path (optional, overwrites input if empty)")
	fs.BoolVar(&dryRun, "dry-run", false, "Simulate repair without writing to disk")
	fs.StringVar(&format, "format", "text", "Output format: text or json")

	if err := fs.Parse(args); err != nil {
		printErrorf(stderr, "invalid repair arguments: %v", err)
		printRepairUsage(stderr)
		return exitUsage
	}
	if len(fs.Args()) > 0 {
		printErrorf(stderr, "unexpected trailing arguments: %s", strings.Join(fs.Args(), " "))
		printRepairUsage(stderr)
		return exitUsage
	}

	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		printErrorf(stderr, "repair requires -file")
		printRepairUsage(stderr)
		return exitUsage
	}

	if outPath == "" {
		outPath = filePath
	}
	jsonMode = strings.EqualFold(strings.TrimSpace(format), "json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		if jsonMode {
			outputJSONError(stdout, fmt.Sprintf("failed to read file: %v", err))
			return exitIO
		}
		printErrorf(stderr, "failed to read file: %v", err)
		return exitIO
	}

	repairedData, result, err := pptx.Repair(data)
	if err != nil {
		if jsonMode {
			outputJSONError(stdout, fmt.Sprintf("repair failed: %v", err))
			return exitIO
		}
		printErrorf(stderr, "repair failed: %v", err)
		return exitIO
	}

	var jsonOut []byte
	if jsonMode {
		outObj := map[string]any{
			"dry_run":           dryRun,
			"issues_repaired":   result.IssuesRepaired,
			"issues_unrepaired": result.IssuesUnrepaired,
		}
		jsonOut, err = json.MarshalIndent(outObj, "", "  ")
		if err != nil {
			outputJSONError(stdout, fmt.Sprintf("failed to marshal JSON: %v", err))
			return exitIO
		}
	}

	if !jsonMode && len(result.IssuesRepaired) > 0 {
		_, _ = fmt.Fprintf(stdout, "Successfully repaired %d issues:\n", len(result.IssuesRepaired))
		for _, issue := range result.IssuesRepaired {
			_, _ = fmt.Fprintf(stdout, "  - %s\n", issue.Description)
		}
	}

	if !jsonMode && len(result.IssuesUnrepaired) > 0 {
		_, _ = fmt.Fprintf(stderr, "Could not repair %d issues:\n", len(result.IssuesUnrepaired))
		for _, issue := range result.IssuesUnrepaired {
			_, _ = fmt.Fprintf(stderr, "  - %s\n", issue.Description)
		}
	}

	if dryRun {
		if jsonMode {
			_, _ = fmt.Fprintln(stdout, string(jsonOut))
		} else {
			_, _ = fmt.Fprintln(stdout, "Dry run complete. No files written.")
		}
		if len(result.IssuesUnrepaired) > 0 {
			return exitValidate
		}
		return exitOK
	}

	// Write to a temporary file first, then rename for atomic overwrite
	// This prevents data loss if the write fails midway
	outDir := filepath.Dir(outPath)
	tmpFile, err := os.CreateTemp(outDir, ".repair-*.pptx")
	if err != nil {
		if jsonMode {
			outputJSONError(stdout, fmt.Sprintf("failed to create temp file: %v", err))
			return exitIO
		}
		printErrorf(stderr, "failed to create temp file: %v", err)
		return exitIO
	}
	tmpPath := tmpFile.Name()

	// Write repaired data to temp file
	_, writeErr := tmpFile.Write(repairedData)
	closeErr := tmpFile.Close()
	if writeErr != nil {
		os.Remove(tmpPath)
		if jsonMode {
			outputJSONError(stdout, fmt.Sprintf("failed to write repaired file: %v", writeErr))
			return exitIO
		}
		printErrorf(stderr, "failed to write repaired file: %v", writeErr)
		return exitIO
	}
	if closeErr != nil {
		os.Remove(tmpPath)
		if jsonMode {
			outputJSONError(stdout, fmt.Sprintf("failed to close temp file: %v", closeErr))
			return exitIO
		}
		printErrorf(stderr, "failed to close temp file: %v", closeErr)
		return exitIO
	}

	// Rename temp file to target (atomic on most filesystems)
	if err := os.Rename(tmpPath, outPath); err != nil {
		os.Remove(tmpPath)
		if jsonMode {
			outputJSONError(stdout, fmt.Sprintf("failed to save repaired file: %v", err))
			return exitIO
		}
		printErrorf(stderr, "failed to save repaired file: %v", err)
		return exitIO
	}

	if jsonMode {
		_, _ = fmt.Fprintln(stdout, string(jsonOut))
	} else {
		_, _ = fmt.Fprintf(stdout, "Repaired file saved to: %s\n", outPath)
	}
	if len(result.IssuesUnrepaired) > 0 {
		return exitValidate
	}
	return exitOK
}

func printRepairUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Usage: pptcli repair -file file.pptx [-out fixed.pptx] [-dry-run] [-format text|json]")
}
