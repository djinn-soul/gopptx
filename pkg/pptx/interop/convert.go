package interop

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ConvertFromPpt uses LibreOffice (soffice) headlessly to convert a legacy .ppt file into a modern .pptx file.
// If outDir is empty, the resulting file is placed in the same directory as the input file.
// The returned string is the absolute path to the generated .pptx file.
func ConvertFromPpt(inputPath string, outDir string) (string, error) {
	inputPath = strings.TrimSpace(inputPath)
	if inputPath == "" {
		return "", errors.New("inputPath is empty")
	}

	absInput, err := filepath.Abs(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for input: %w", err)
	}

	if _, err := os.Stat(absInput); err != nil {
		return "", fmt.Errorf("input file not found: %w", err)
	}

	if outDir == "" {
		outDir = filepath.Dir(absInput)
	} else {
		var err error
		outDir, err = filepath.Abs(outDir)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path for outDir: %w", err)
		}
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return "", fmt.Errorf("failed to create outDir: %w", err)
		}
	}

	soffice, err := findSoffice()
	if err != nil {
		return "", fmt.Errorf("libreoffice required for conversion: %w", err)
	}

	// Calculate the expected output path
	baseName := filepath.Base(absInput)
	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)
	expectedOutPath := filepath.Join(outDir, nameWithoutExt+".pptx")

	// LibreOffice headless conversion command
	args := []string{
		"--headless",
		"--invisible",
		"--nologo",
		"--nodefault",
		"--norestore",
		"--convert-to", "pptx",
		"--outdir", outDir,
		absInput,
	}

	cmd := exec.Command(soffice, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("conversion failed: %w\nOutput: %s", err, string(out))
	}

	// Verify the file was actually created
	if _, err := os.Stat(expectedOutPath); err != nil {
		return "", fmt.Errorf(
			"conversion appeared successful but output file not found at %s: %w",
			expectedOutPath,
			err,
		)
	}

	return expectedOutPath, nil
}

// findSoffice attempts to locate the LibreOffice binary in the system path or common install locations.
func findSoffice() (string, error) {
	// 1. Try standard system PATH
	exeName := "soffice"
	if runtime.GOOS == "windows" {
		exeName = "soffice.exe"
	}

	foundPath, err := exec.LookPath(exeName)
	if err == nil {
		return foundPath, nil
	}

	// 2. Try common default installation directories
	var commonPaths []string
	if runtime.GOOS == "windows" {
		commonPaths = []string{
			`C:\Program Files\LibreOffice\program\soffice.exe`,
			`C:\Program Files (x86)\LibreOffice\program\soffice.exe`,
		}
	} else if runtime.GOOS == "darwin" {
		commonPaths = []string{
			`/Applications/LibreOffice.app/Contents/MacOS/soffice`,
		}
	} else {
		commonPaths = []string{
			`/usr/bin/soffice`,
			`/usr/local/bin/soffice`,
			`/opt/libreoffice/program/soffice`,
		}
	}

	for _, p := range commonPaths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", errors.New("soffice binary not found in PATH or standard installation directories")
}
