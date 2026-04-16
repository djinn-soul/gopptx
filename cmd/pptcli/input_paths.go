package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type presentationInputKind uint8

const (
	inputKindUnknown presentationInputKind = iota
	inputKindMarkdown
	inputKindPPTX
)

func detectPresentationInputKind(path string) (presentationInputKind, error) {
	cleanPath := strings.TrimSpace(path)
	if cleanPath == "" {
		return inputKindUnknown, errors.New("input path cannot be empty")
	}

	switch strings.ToLower(filepath.Ext(cleanPath)) {
	case ".md":
		return inputKindMarkdown, nil
	case ".pptx", ".pptm":
		return inputKindPPTX, nil
	default:
		return inputKindUnknown, fmt.Errorf(
			"unsupported input file extension %q (allowed: .pptx, .pptm, .md)",
			filepath.Ext(cleanPath),
		)
	}
}

func defaultSiblingFilePath(path, fallbackStem, newExt string) string {
	dir, stem := siblingStem(path, fallbackStem)
	name := stem + newExt
	if dir == "." || dir == "" {
		return name
	}
	return filepath.Join(dir, name)
}

func defaultSiblingDirPath(path, fallbackStem, suffix string) string {
	dir, stem := siblingStem(path, fallbackStem)
	name := stem + suffix
	if dir == "." || dir == "" {
		return name
	}
	return filepath.Join(dir, name)
}

func siblingStem(path, fallbackStem string) (string, string) {
	cleanPath := strings.TrimSpace(path)
	if cleanPath == "" {
		return "", fallbackStem
	}

	dir := filepath.Dir(cleanPath)
	base := filepath.Base(cleanPath)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	if strings.TrimSpace(stem) == "" {
		stem = fallbackStem
	}
	return dir, stem
}
