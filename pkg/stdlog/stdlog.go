package stdlog

import (
	"fmt"
	"os"
	"strings"
)

// Println writes a line to stderr.
func Println(args ...any) {
	_, _ = fmt.Fprintln(os.Stderr, args...)
}

// Printf writes formatted output to stderr.
func Printf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
}

// Fatal writes a line to stderr and exits with status 1.
func Fatal(args ...any) {
	Println(args...)
	os.Exit(1)
}

// Fatalf writes formatted output to stderr and exits with status 1.
func Fatalf(format string, args ...any) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	Printf(format, args...)
	os.Exit(1)
}
