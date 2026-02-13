package editor

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"testing"
)

var pyOpLine = regexp.MustCompile(`(?m)^OP_[A-Z_]+\s*=\s*"([^"]+)"\s*$`)

func TestSupportedOpsMatchPythonConstants(t *testing.T) {
	// Robustly find the project root by walking up looking for go.mod
	wd, _ := os.Getwd()
	root := wd
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Fatal("could not find project root (go.mod)")
		}
		root = parent
	}
	opsPath := filepath.Join(root, "python", "gopptx", "ops.py")
	content, err := os.ReadFile(opsPath)
	if err != nil {
		t.Fatalf("read python ops map from %s: %v", opsPath, err)
	}

	var pyOps []string
	for _, match := range pyOpLine.FindAllStringSubmatch(string(content), -1) {
		pyOps = append(pyOps, match[1])
	}
	if len(pyOps) == 0 {
		t.Fatal("no OP_* constants found in python/gopptx/ops.py")
	}

	goOps := append([]string(nil), SupportedOps...)
	slices.Sort(goOps)
	slices.Sort(pyOps)

	if !slices.Equal(goOps, pyOps) {
		t.Fatalf("go/python op mismatch\nGo: %v\nPython: %v", goOps, pyOps)
	}
}

func TestSupportedOpsMatchCommandHandlers(t *testing.T) {
	handlerOps := make([]string, 0, len(commandHandlers))
	for op := range commandHandlers {
		handlerOps = append(handlerOps, op)
	}
	if len(handlerOps) == 0 {
		t.Fatal("no command handlers registered")
	}
	slices.Sort(handlerOps)

	goOps := append([]string(nil), SupportedOps...)
	slices.Sort(goOps)

	if !slices.Equal(goOps, handlerOps) {
		t.Fatalf("supported ops/handler mismatch\nSupportedOps: %v\nHandlers: %v", goOps, handlerOps)
	}
}
