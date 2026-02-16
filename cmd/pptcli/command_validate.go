package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

var requiredPPTXParts = []string{
	"[Content_Types].xml",
	"_rels/.rels",
	"ppt/presentation.xml",
	"docProps/core.xml",
}

func runValidateCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var filePath string
	fs.StringVar(&filePath, "file", "", "PPTX file path")

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

	issues, err := validatePPTXFile(filePath)
	if err != nil {
		printErrorf(stderr, "validation failed to run: %v", err)
		return exitIO
	}
	if len(issues) > 0 {
		sort.Strings(issues)
		for _, issue := range issues {
			_, _ = fmt.Fprintf(stderr, "ERROR: %s\n", issue)
		}
		return exitValidate
	}

	_, _ = fmt.Fprintln(stdout, "OK: validation passed")
	return exitOK
}

func validatePPTXFile(path string) ([]string, error) {
	meta, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !meta.Mode().IsRegular() {
		return nil, errors.New("path is not a regular file")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	zr, err := zip.NewReader(file, meta.Size())
	if err != nil {
		return []string{"file is not a valid ZIP archive"}, nil
	}

	issues := make([]string, 0, 8)
	names := make(map[string]struct{}, len(zr.File))
	slideCount := 0
	for _, entry := range zr.File {
		name := strings.TrimSpace(entry.Name)
		if entry.FileInfo().IsDir() {
			continue
		}
		names[name] = struct{}{}
		lower := strings.ToLower(name)
		if strings.HasPrefix(lower, "ppt/slides/slide") && strings.HasSuffix(lower, ".xml") &&
			!strings.Contains(lower, "/_rels/") {
			slideCount++
		}
		if strings.HasSuffix(lower, ".xml") || strings.HasSuffix(lower, ".rels") {
			if err := validateEntryXML(entry); err != nil {
				issues = append(issues, fmt.Sprintf("%s: %v", name, err))
			}
		}
	}

	for _, required := range requiredPPTXParts {
		if _, ok := names[required]; !ok {
			issues = append(issues, fmt.Sprintf("missing required part %q", required))
		}
	}
	if slideCount == 0 {
		issues = append(issues, "no slide parts found under ppt/slides")
	}

	return issues, nil
}

func validateEntryXML(entry *zip.File) error {
	reader, err := entry.Open()
	if err != nil {
		return err
	}
	defer func() { _ = reader.Close() }()

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return errors.New("empty XML content")
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	for {
		if _, err := decoder.Token(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("invalid XML: %w", err)
		}
	}
	return nil
}

func printValidateUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Usage: pptcli validate -file file.pptx")
}
