package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func writeOutputFile(path string, data []byte) error {
	cleanPath := strings.TrimSpace(path)
	if cleanPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}
	if err := ensureParentDir(cleanPath); err != nil {
		return err
	}
	return os.WriteFile(cleanPath, data, 0o600)
}

func ensureParentDir(path string) error {
	parent := filepath.Dir(path)
	if parent == "." || parent == "" {
		return nil
	}
	return os.MkdirAll(parent, 0o755)
}

func defaultOutputPathFromMarkdown(markdownPath string) string {
	cleanPath := strings.TrimSpace(markdownPath)
	if cleanPath == "" {
		return "output.pptx"
	}
	dir := filepath.Dir(cleanPath)
	base := filepath.Base(cleanPath)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	if stem == "" {
		stem = "output"
	}
	out := stem + ".pptx"
	if dir == "." || dir == "" {
		return out
	}
	return filepath.Join(dir, out)
}
