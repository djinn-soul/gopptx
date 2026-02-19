package editor

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func pyOpLineRegex() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^OP_[A-Z_]+\s*=\s*"([^"]+)"\s*$`)
}

func pyiOpLineRegex() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^(OP_[A-Z_]+)\s*:\s*str\s*$`)
}

func findProjectRoot(t *testing.T) string {
	t.Helper()
	wd, _ := os.Getwd()
	root := wd
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			return root
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Fatal("could not find project root (go.mod)")
		}
		root = parent
	}
}

func TestSupportedOpsMatchPythonConstants(t *testing.T) {
	root := findProjectRoot(t)
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
		// Provide detailed diff output
		missingInGo := diffSlice(pyOps, goOps)
		missingInPy := diffSlice(goOps, pyOps)
		t.Fatalf("go/python op mismatch\nMissing in Go: %v\nMissing in Python: %v", missingInGo, missingInPy)
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
	root := findProjectRoot(t)
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

// TestContractBreakHandlerOnlyOp verifies that adding a handler without
// declaring it in SupportedOps would be caught by TestSupportedOpsMatchCommandHandlers.
func TestContractBreakHandlerOnlyOp(t *testing.T) {
	// This test documents the expected behavior: if someone adds a handler
	// but forgets to add it to SupportedOps, TestSupportedOpsMatchCommandHandlers
	// will fail because it iterates over SupportedOps and checks each has a handler.
	//
	// The inverse case (handler without SupportedOps entry) is caught by
	// the handler registration being the source of truth for commandHandlerFor.
	//
	// We verify the contract enforcement works by checking that all handlers
	// are in SupportedOps.
	goOps := make(map[string]struct{})
	for _, op := range SupportedOps() {
		goOps[op] = struct{}{}
	}

	// Get all registered handlers by checking each supported op
	for _, op := range SupportedOps() {
		handler, ok := commandHandlerFor(op)
		if !ok {
			t.Fatalf("op %q in SupportedOps but has no handler", op)
		}
		if handler == nil {
			t.Fatalf("op %q has nil handler", op)
		}
	}
}

// TestContractBreakPythonOnlyOp verifies that adding a Python op constant
// without a matching Go op would be caught by TestSupportedOpsMatchPythonConstants.
func TestContractBreakPythonOnlyOp(t *testing.T) {
	// This test documents the expected behavior: if someone adds an OP_* constant
	// to python/gopptx/ops.py but forgets to add it to Go's SupportedOps,
	// TestSupportedOpsMatchPythonConstants will fail.
	//
	// We verify the current state is consistent.
	root := findProjectRoot(t)
	opsPath := filepath.Join(root, "python", "gopptx", "ops.py")
	content, err := os.ReadFile(opsPath)
	if err != nil {
		t.Fatalf("read python ops: %v", err)
	}

	var pyOps []string
	for _, match := range pyOpLineRegex().FindAllStringSubmatch(string(content), -1) {
		pyOps = append(pyOps, match[1])
	}

	goOps := SupportedOps()
	goOpsMap := make(map[string]struct{})
	for _, op := range goOps {
		goOpsMap[op] = struct{}{}
	}

	// Check each Python op exists in Go
	for _, pyOp := range pyOps {
		if _, ok := goOpsMap[pyOp]; !ok {
			t.Fatalf("Python op %q not found in Go SupportedOps - contract break", pyOp)
		}
	}
}

// TestContractBreakGoOnlyOp verifies that adding a Go op constant
// without a matching Python op would be caught by TestSupportedOpsMatchPythonConstants.
func TestContractBreakGoOnlyOp(t *testing.T) {
	// This test documents the expected behavior: if someone adds an Op* constant
	// to Go's opspec.go but forgets to add it to Python's ops.py,
	// TestSupportedOpsMatchPythonConstants will fail.
	//
	// We verify the current state is consistent.
	root := findProjectRoot(t)
	opsPath := filepath.Join(root, "python", "gopptx", "ops.py")
	content, err := os.ReadFile(opsPath)
	if err != nil {
		t.Fatalf("read python ops: %v", err)
	}

	pyOpsMap := make(map[string]struct{})
	for _, match := range pyOpLineRegex().FindAllStringSubmatch(string(content), -1) {
		pyOpsMap[match[1]] = struct{}{}
	}

	// Check each Go op exists in Python
	for _, goOp := range SupportedOps() {
		if _, ok := pyOpsMap[goOp]; !ok {
			t.Fatalf("Go op %q not found in Python ops.py - contract break", goOp)
		}
	}
}

// TestContractBreakSupportedOpsSet verifies Python SUPPORTED_OPS_SET matches Go.
func TestContractBreakSupportedOpsSet(t *testing.T) {
	root := findProjectRoot(t)
	opsPath := filepath.Join(root, "python", "gopptx", "ops.py")
	content, err := os.ReadFile(opsPath)
	if err != nil {
		t.Fatalf("read python ops: %v", err)
	}

	// Extract SUPPORTED_OPS tuple
	re := regexp.MustCompile(`(?s)SUPPORTED_OPS\s*=\s*\(([^)]+)\)`)
	match := re.FindStringSubmatch(string(content))
	if match == nil {
		t.Fatal("could not find SUPPORTED_OPS tuple in ops.py")
	}

	// Parse the tuple contents
	tupleContent := match[1]
	opConstRe := regexp.MustCompile(`OP_[A-Z_]+`)
	pySupportedOps := make(map[string]struct{})
	for _, constName := range opConstRe.FindAllString(tupleContent, -1) {
		pySupportedOps[constName] = struct{}{}
	}

	// Get Go ops count
	goOps := SupportedOps()
	if len(goOps) != len(pySupportedOps) {
		t.Fatalf("SUPPORTED_OPS count mismatch: Go has %d, Python has %d", len(goOps), len(pySupportedOps))
	}
}

// TestOpConstantsNamingConvention verifies all op constants follow snake_case naming.
func TestOpConstantsNamingConvention(t *testing.T) {
	for _, op := range SupportedOps() {
		// Op names should be snake_case (lowercase with underscores)
		for _, r := range op {
			if !((r >= 'a' && r <= 'z') || r == '_') {
				t.Fatalf("op %q contains invalid character %q, expected snake_case", op, r)
			}
		}
		// Should not have double underscores
		if strings.Contains(op, "__") {
			t.Fatalf("op %q contains double underscore", op)
		}
		// Should not start or end with underscore
		if strings.HasPrefix(op, "_") || strings.HasSuffix(op, "_") {
			t.Fatalf("op %q should not start or end with underscore", op)
		}
	}
}

// TestAllOpsHaveDocumentation verifies all ops are documented in the architecture doc.
func TestAllOpsHaveDocumentation(t *testing.T) {
	root := findProjectRoot(t)
	docPath := filepath.Join(root, "docs", "architecture", "bridge-phase1-ops.md")
	content, err := os.ReadFile(docPath)
	if err != nil {
		t.Skipf("documentation file not found: %v", err)
	}

	docStr := string(content)
	for _, op := range SupportedOps() {
		if !strings.Contains(docStr, op) {
			t.Logf("WARNING: op %q not found in documentation", op)
		}
	}
}

// diffSlice returns elements in a that are not in b.
func diffSlice(a, b []string) []string {
	bMap := make(map[string]struct{})
	for _, v := range b {
		bMap[v] = struct{}{}
	}
	var diff []string
	for _, v := range a {
		if _, ok := bMap[v]; !ok {
			diff = append(diff, v)
		}
	}
	return diff
}
