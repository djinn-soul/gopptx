package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type pptxFileInfo struct {
	path       string
	size       int64
	modifiedAt time.Time
	isPPTXZip  bool
	slideCount int
	chartCount int
	imageCount int
}

func runInfoCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("info", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var filePath string
	fs.StringVar(&filePath, "file", "", "PPTX file path")

	if err := fs.Parse(args); err != nil {
		printErrorf(stderr, "invalid info arguments: %v", err)
		printInfoUsage(stderr)
		return exitUsage
	}
	if len(fs.Args()) > 0 {
		printErrorf(stderr, "unexpected trailing arguments: %s", strings.Join(fs.Args(), " "))
		printInfoUsage(stderr)
		return exitUsage
	}
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		printErrorf(stderr, "info requires -file")
		printInfoUsage(stderr)
		return exitUsage
	}

	info, err := inspectPPTX(filePath)
	if err != nil {
		printErrorf(stderr, "failed to inspect %q: %v", filePath, err)
		return exitIO
	}

	_, _ = fmt.Fprintf(stdout, "Path: %s\n", info.path)
	_, _ = fmt.Fprintf(stdout, "Size: %d bytes\n", info.size)
	_, _ = fmt.Fprintf(stdout, "Modified: %s\n", info.modifiedAt.Format(time.RFC3339))
	if info.isPPTXZip {
		_, _ = fmt.Fprintln(stdout, "Format: valid ZIP package")
		_, _ = fmt.Fprintf(stdout, "Slide count: %d\n", info.slideCount)
		_, _ = fmt.Fprintf(stdout, "Chart parts: %d\n", info.chartCount)
		_, _ = fmt.Fprintf(stdout, "Media items: %d\n", info.imageCount)
	} else {
		_, _ = fmt.Fprintln(stdout, "Format: not a valid PPTX zip package")
	}

	return exitOK
}

func inspectPPTX(path string) (pptxFileInfo, error) {
	meta, err := os.Stat(path)
	if err != nil {
		return pptxFileInfo{}, err
	}
	out := pptxFileInfo{
		path:       path,
		size:       meta.Size(),
		modifiedAt: meta.ModTime(),
	}
	if !meta.Mode().IsRegular() {
		return out, errors.New("path is not a regular file")
	}

	file, err := os.Open(path)
	if err != nil {
		return out, err
	}
	defer func() { _ = file.Close() }()

	zr, err := zip.NewReader(file, meta.Size())
	if err != nil {
		return out, nil
	}

	out.isPPTXZip = true
	for _, entry := range zr.File {
		if entry.FileInfo().IsDir() {
			continue
		}
		name := strings.ToLower(strings.TrimSpace(entry.Name))
		switch {
		case strings.HasPrefix(name, "ppt/slides/slide") && strings.HasSuffix(name, ".xml") && !strings.Contains(name, "/_rels/"):
			out.slideCount++
		case strings.HasPrefix(name, "ppt/charts/chart") && strings.HasSuffix(name, ".xml"):
			out.chartCount++
		case strings.HasPrefix(name, "ppt/media/"):
			out.imageCount++
		}
	}
	return out, nil
}

func printInfoUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Usage: pptcli info -file file.pptx")
}
