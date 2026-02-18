package editor

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"testing"
)

func pyOpLineRegex() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^OP_[A-Z_]+\s*=\s*"([^"]+)"\s*$`)
}

func pyiOpLineRegex() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^(OP_[A-Z_]+)\s*:\s*str\s*$`)
}

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
	for _, match := range pyOpLineRegex().FindAllStringSubmatch(string(content), -1) {
		pyOps = append(pyOps, match[1])
	}
	if len(pyOps) == 0 {
		t.Fatal("no OP_* constants found in python/gopptx/ops.py")
	}

	goOps := SupportedOps()
	slices.Sort(goOps)
	slices.Sort(pyOps)

	if !slices.Equal(goOps, pyOps) {
		t.Fatalf("go/python op mismatch\nGo: %v\nPython: %v", goOps, pyOps)
	}
}

func TestSupportedOpsMatchCommandHandlers(t *testing.T) {
	goOps := SupportedOps()
	handlerOps := make([]string, 0, len(goOps))
	for _, op := range goOps {
		_, ok := commandHandlerFor(op)
		if !ok {
			t.Fatalf("missing command handler for op %q", op)
		}
		handlerOps = append(handlerOps, op)
	}
	if len(goOps) == 0 {
		t.Fatal("no command handlers registered")
	}
	if _, ok := commandHandlerFor("___unknown___"); ok {
		t.Fatal("unexpected command handler for unknown op")
	}
	slices.Sort(handlerOps)
	slices.Sort(goOps)

	if !slices.Equal(goOps, handlerOps) {
		t.Fatalf("supported ops/handler mismatch\nSupportedOps: %v\nHandlers: %v", goOps, handlerOps)
	}
}

func TestPythonOpsStubDeclaresAllConstants(t *testing.T) {
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
	pyiPath := filepath.Join(root, "python", "gopptx", "ops.pyi")

	opsContent, err := os.ReadFile(opsPath)
	if err != nil {
		t.Fatalf("read %s: %v", opsPath, err)
	}
	pyiContent, err := os.ReadFile(pyiPath)
	if err != nil {
		t.Fatalf("read %s: %v", pyiPath, err)
	}

	pyNames := make(map[string]struct{})
	for _, line := range regexp.MustCompile(`(?m)^(OP_[A-Z_]+)\s*=\s*"[^"]+"\s*$`).FindAllStringSubmatch(string(opsContent), -1) {
		pyNames[line[1]] = struct{}{}
	}
	pyiNames := make(map[string]struct{})
	for _, line := range pyiOpLineRegex().FindAllStringSubmatch(string(pyiContent), -1) {
		pyiNames[line[1]] = struct{}{}
	}
	for name := range pyNames {
		if _, ok := pyiNames[name]; !ok {
			t.Fatalf("ops.pyi missing constant declaration for %s", name)
		}
	}
}
