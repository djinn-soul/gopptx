package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func writeOutputFile(path string, data []byte) error {
	cleanPath := strings.TrimSpace(path)
	if cleanPath == "" {
		return errors.New("output path cannot be empty")
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
	return os.MkdirAll(parent, 0o750)
}

func defaultOutputPathFromMarkdown(markdownPath string) string {
	return defaultSiblingFilePath(markdownPath, "output", ".pptx")
}
