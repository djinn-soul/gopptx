package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type validationResult struct {
	Path   string
	Errors []string
}

func main() {
	dirFlag := flag.String("dir", "smoke_samples", "directory to scan for .pptx files")
	fileFlag := flag.String("file", "", "optional single .pptx file to validate")
	flag.Parse()

	files, err := collectPPTXFiles(*dirFlag, *fileFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "collect files: %v\n", err)
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "no .pptx files found (dir=%q, file=%q)\n", *dirFlag, *fileFlag)
		os.Exit(1)
	}

	var failed int
	for _, path := range files {
		result := validatePPTXPath(path)
		if len(result.Errors) == 0 {
			fmt.Printf("PASS %s\n", path)
			continue
		}
		failed++
		fmt.Printf("FAIL %s\n", path)
		for _, msg := range result.Errors {
			fmt.Printf("  - %s\n", msg)
		}
	}

	fmt.Printf("Validated %d file(s); failures=%d\n", len(files), failed)
	if failed > 0 {
		os.Exit(1)
	}
}

func collectPPTXFiles(dirPath, singleFile string) ([]string, error) {
	paths := make([]string, 0)
	seen := make(map[string]struct{})

	if strings.TrimSpace(singleFile) != "" {
		info, err := os.Stat(singleFile)
		if err != nil {
			return nil, fmt.Errorf("stat file %q: %w", singleFile, err)
		}
		if info.IsDir() {
			return nil, fmt.Errorf("file %q is a directory", singleFile)
		}
		if strings.ToLower(filepath.Ext(singleFile)) != ".pptx" {
			return nil, fmt.Errorf("file %q is not a .pptx", singleFile)
		}
		abs, err := filepath.Abs(singleFile)
		if err != nil {
			return nil, err
		}
		seen[abs] = struct{}{}
		paths = append(paths, abs)
	}

	if strings.TrimSpace(dirPath) != "" {
		info, err := os.Stat(dirPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Errorf("directory %q does not exist", dirPath)
			}
			return nil, err
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("dir %q is not a directory", dirPath)
		}

		err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				return nil
			}
			if strings.ToLower(filepath.Ext(path)) != ".pptx" {
				return nil
			}
			abs, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			if _, ok := seen[abs]; ok {
				return nil
			}
			seen[abs] = struct{}{}
			paths = append(paths, abs)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	sort.Strings(paths)
	return paths, nil
}

func validatePPTXPath(path string) validationResult {
	result := validationResult{Path: path}
	data, err := os.ReadFile(path)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("read file: %v", err))
		return result
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("open zip: %v", err))
		return result
	}

	names := make(map[string]struct{}, len(zr.File))
	slideCount := 0
	for _, f := range zr.File {
		names[f.Name] = struct{}{}
		if strings.HasPrefix(f.Name, "ppt/slides/slide") && strings.HasSuffix(f.Name, ".xml") {
			slideCount++
		}
	}

	required := []string{
		"[Content_Types].xml",
		"_rels/.rels",
		"ppt/presentation.xml",
		"ppt/_rels/presentation.xml.rels",
	}
	for _, req := range required {
		if _, ok := names[req]; !ok {
			result.Errors = append(result.Errors, fmt.Sprintf("missing required part: %s", req))
		}
	}
	if slideCount == 0 {
		result.Errors = append(result.Errors, "missing slide parts: ppt/slides/slide*.xml")
	}
	return result
}
