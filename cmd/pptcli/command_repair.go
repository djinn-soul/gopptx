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
	"github.com/djinn-soul/gopptx/pkg/pptx/validation/structural"
)

type repairCommandOptions struct {
	filePath string
	outPath  string
	dryRun   bool
	jsonMode bool
}

// runRepairCommand attempts to fix structural issues in a PPTX file.
func runRepairCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	opts, exitCode, ok := parseRepairCommandOptions(args, stderr)
	if !ok {
		return exitCode
	}

	data, err := os.ReadFile(opts.filePath)
	if err != nil {
		return handleRepairError(stdout, stderr, opts.jsonMode, "failed to read file: %v", err)
	}

	repairedData, result, err := pptx.Repair(data)
	if err != nil {
		return handleRepairError(stdout, stderr, opts.jsonMode, "repair failed: %v", err)
	}

	jsonOut, err := buildRepairOutputJSON(result, opts)
	if err != nil {
		return handleRepairError(stdout, stderr, opts.jsonMode, "failed to marshal JSON: %v", err)
	}

	printRepairSummary(stdout, stderr, result, opts)

	if opts.dryRun {
		return completeRepairDryRun(stdout, result, jsonOut, opts)
	}

	if err := writeRepairedFileAtomically(opts.outPath, repairedData, stderr); err != nil {
		return handleRepairError(stdout, stderr, opts.jsonMode, "failed to save repaired file: %v", err)
	}

	return finalizeRepair(stdout, result, jsonOut, opts)
}

func parseRepairCommandOptions(args []string, stderr io.Writer) (repairCommandOptions, int, bool) {
	opts := repairCommandOptions{}

	fs := flag.NewFlagSet("repair", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var format string
	fs.StringVar(&opts.filePath, "file", "", "Input PPTX file path")
	fs.StringVar(&opts.outPath, "out", "", "Output (repaired) PPTX file path (optional, overwrites input if empty)")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "Simulate repair without writing to disk")
	fs.StringVar(&format, "format", "text", "Output format: text or json")

	if err := fs.Parse(args); err != nil {
		printErrorf(stderr, "invalid repair arguments: %v", err)
		printRepairUsage(stderr)
		return repairCommandOptions{}, exitUsage, false
	}
	if len(fs.Args()) > 0 {
		printErrorf(stderr, "unexpected trailing arguments: %s", strings.Join(fs.Args(), " "))
		printRepairUsage(stderr)
		return repairCommandOptions{}, exitUsage, false
	}

	opts.filePath = strings.TrimSpace(opts.filePath)
	if opts.filePath == "" {
		printErrorf(stderr, "repair requires -file")
		printRepairUsage(stderr)
		return repairCommandOptions{}, exitUsage, false
	}

	if opts.outPath == "" {
		opts.outPath = opts.filePath
	}
	opts.jsonMode = strings.EqualFold(strings.TrimSpace(format), "json")
	return opts, exitOK, true
}

func buildRepairOutputJSON(result structural.RepairResult, opts repairCommandOptions) ([]byte, error) {
	if !opts.jsonMode {
		return nil, nil
	}
	outObj := map[string]any{
		"dry_run":           opts.dryRun,
		"issues_repaired":   result.IssuesRepaired,
		"issues_unrepaired": result.IssuesUnrepaired,
	}
	return json.MarshalIndent(outObj, "", "  ")
}

func printRepairSummary(stdout io.Writer, stderr io.Writer, result structural.RepairResult, opts repairCommandOptions) {
	if !opts.jsonMode && len(result.IssuesRepaired) > 0 {
		_, _ = fmt.Fprintf(stdout, "Successfully repaired %d issues:\n", len(result.IssuesRepaired))
		for _, issue := range result.IssuesRepaired {
			_, _ = fmt.Fprintf(stdout, "  - %s\n", issue.Description)
		}
	}

	if !opts.jsonMode && len(result.IssuesUnrepaired) > 0 {
		_, _ = fmt.Fprintf(stderr, "Could not repair %d issues:\n", len(result.IssuesUnrepaired))
		for _, issue := range result.IssuesUnrepaired {
			_, _ = fmt.Fprintf(stderr, "  - %s\n", issue.Description)
		}
	}
}

func completeRepairDryRun(
	stdout io.Writer,
	result structural.RepairResult,
	jsonOut []byte,
	opts repairCommandOptions,
) int {
	if opts.jsonMode {
		_, _ = fmt.Fprintln(stdout, string(jsonOut))
	} else {
		_, _ = fmt.Fprintln(stdout, "Dry run complete. No files written.")
	}
	if len(result.IssuesUnrepaired) > 0 {
		return exitValidate
	}
	return exitOK
}

func writeRepairedFileAtomically(outPath string, repairedData []byte, stderr io.Writer) error {
	outDir := filepath.Dir(outPath)
	tmpFile, err := os.CreateTemp(outDir, ".repair-*.pptx")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	_, writeErr := tmpFile.Write(repairedData)
	closeErr := tmpFile.Close()
	if writeErr != nil {
		cleanupRepairTempFile(tmpPath, stderr)
		return fmt.Errorf("failed to write repaired file: %w", writeErr)
	}
	if closeErr != nil {
		cleanupRepairTempFile(tmpPath, stderr)
		return fmt.Errorf("failed to close temp file: %w", closeErr)
	}

	if renameErr := os.Rename(tmpPath, outPath); renameErr != nil {
		cleanupRepairTempFile(tmpPath, stderr)
		return renameErr
	}
	return nil
}

func cleanupRepairTempFile(tmpPath string, stderr io.Writer) {
	if removeErr := os.Remove(tmpPath); removeErr != nil && !os.IsNotExist(removeErr) {
		printErrorf(stderr, "failed to cleanup temp file: %v", removeErr)
	}
}

func finalizeRepair(stdout io.Writer, result structural.RepairResult, jsonOut []byte, opts repairCommandOptions) int {
	if opts.jsonMode {
		_, _ = fmt.Fprintln(stdout, string(jsonOut))
	} else {
		_, _ = fmt.Fprintf(stdout, "Repaired file saved to: %s\n", opts.outPath)
	}
	if len(result.IssuesUnrepaired) > 0 {
		return exitValidate
	}
	return exitOK
}

func printRepairUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Usage: pptcli repair -file file.pptx [-out fixed.pptx] [-dry-run] [-format text|json]")
}

func handleRepairError(stdout io.Writer, stderr io.Writer, jsonMode bool, format string, a ...any) int {
	msg := fmt.Sprintf(format, a...)
	if jsonMode {
		outputJSONError(stdout, msg)
	} else {
		printErrorf(stderr, msg)
	}
	return exitIO
}
