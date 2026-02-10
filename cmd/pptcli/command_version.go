package main

import (
	"fmt"
	"io"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func runVersionCommand(_ []string, stdout io.Writer, _ io.Writer) int {
	_, _ = fmt.Fprintf(stdout, "gopptx version %s\n", pptx.Version)
	return exitOK
}
