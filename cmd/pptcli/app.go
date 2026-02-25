package main

import (
	"fmt"
	"io"
	"strings"
)

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		return runLegacyFlags(nil, stdout, stderr)
	}

	switch strings.ToLower(strings.TrimSpace(args[0])) {
	case "help", "-h", "--help":
		printRootUsage(stdout)
		return exitOK
	case "create":
		return runCreateCommand(args[1:], stdout, stderr)
	case "md2ppt":
		return runMD2PPTCommand(args[1:], stdout, stderr)
	case "info":
		return runInfoCommand(args[1:], stdout, stderr)
	case "validate":
		return runValidateCommand(args[1:], stdout, stderr)
	case "repair":
		return runRepairCommand(args[1:], stdout, stderr)
	case "merge":
		return runMergeCommand(args[1:], stdout, stderr)
	case "completion":
		return runCompletionCommand(args[1:], stdout, stderr)
	case "version", "-version", "--version":
		return runVersionCommand(args[1:], stdout, stderr)
	default:
		if strings.HasPrefix(args[0], "-") {
			return runLegacyFlags(args, stdout, stderr)
		}
		printErrorf(stderr, "unknown command %q", args[0])
		printRootUsage(stderr)
		return exitUsage
	}
}

func printRootUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "gopptx CLI")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Usage:")
	_, _ = fmt.Fprintln(w, "  pptcli create   -out file.pptx [-title TITLE] [-slides N]")
	_, _ = fmt.Fprintln(w, "  pptcli md2ppt   -in deck.md [-out file.pptx] [-title TITLE]")
	_, _ = fmt.Fprintln(w, "  pptcli info     -file file.pptx")
	_, _ = fmt.Fprintln(w, "  pptcli validate -file file.pptx")
	_, _ = fmt.Fprintln(w, "  pptcli repair   -file file.pptx [-out fixed.pptx]")
	_, _ = fmt.Fprintln(w, "  pptcli merge    -out merged.pptx file1.pptx file2.pptx ...")
	_, _ = fmt.Fprintln(w, "  pptcli completion -shell bash|zsh")
	_, _ = fmt.Fprintln(w, "  pptcli version")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Legacy mode:")
	_, _ = fmt.Fprintln(w, "  pptcli [-out output.pptx] [-md input.md] [-title TITLE]")
}

func printErrorf(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, "ERROR: "+format+"\n", args...)
}
