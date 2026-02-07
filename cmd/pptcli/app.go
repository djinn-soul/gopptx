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
	fmt.Fprintln(w, "goppt CLI")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  pptcli create   -out file.pptx [-title TITLE] [-slides N]")
	fmt.Fprintln(w, "  pptcli md2ppt   -in deck.md [-out file.pptx] [-title TITLE]")
	fmt.Fprintln(w, "  pptcli info     -file file.pptx")
	fmt.Fprintln(w, "  pptcli validate -file file.pptx")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Legacy mode:")
	fmt.Fprintln(w, "  pptcli [-out output.pptx] [-md input.md] [-title TITLE]")
}

func printErrorf(w io.Writer, format string, args ...any) {
	fmt.Fprintf(w, "ERROR: "+format+"\n", args...)
}
