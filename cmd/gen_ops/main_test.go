package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUnquoteEscapes(t *testing.T) {
	got, err := unquote(`"line\nnext"`)
	if err != nil {
		t.Fatalf("unquote error: %v", err)
	}
	if got != "line\nnext" {
		t.Fatalf("unexpected unquote value: %q", got)
	}
}

func TestToSnakeCase(t *testing.T) {
	if got := toSnakeCase("BatchExecute"); got != "BATCH_EXECUTE" {
		t.Fatalf("unexpected snake case: %q", got)
	}
	if got := toSnakeCase("URLFetch"); got != "URL_FETCH" {
		t.Fatalf("unexpected snake case for acronym: %q", got)
	}
}

func TestParseOpsFromGo(t *testing.T) {
	tmp := t.TempDir()
	input := filepath.Join(tmp, "opspec.go")
	src := `package sample
const (
	OpSlideCount = "slide_count"
	NotOp = "ignore"
	OpBatchExecute = "batch_execute"
)`
	if err := os.WriteFile(input, []byte(src), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	ops, err := parseOpsFromGo(input)
	if err != nil {
		t.Fatalf("parseOpsFromGo: %v", err)
	}
	if len(ops) != 2 {
		t.Fatalf("expected 2 ops, got %d", len(ops))
	}

	got := map[string]string{}
	for _, op := range ops {
		got[op.PyName] = op.Value
	}
	if got["OP_SLIDE_COUNT"] != "slide_count" {
		t.Fatalf("missing OP_SLIDE_COUNT mapping: %+v", got)
	}
	if got["OP_BATCH_EXECUTE"] != "batch_execute" {
		t.Fatalf("missing OP_BATCH_EXECUTE mapping: %+v", got)
	}
}

func TestWriteOpsOutputs(t *testing.T) {
	tmp := t.TempDir()
	pyPath := filepath.Join(tmp, "ops.py")
	pyiPath := filepath.Join(tmp, "ops.pyi")

	py, err := os.Create(pyPath)
	if err != nil {
		t.Fatalf("create py: %v", err)
	}
	pyi, err := os.Create(pyiPath)
	if err != nil {
		_ = py.Close()
		t.Fatalf("create pyi: %v", err)
	}

	ops := []opSpec{
		{PyName: "OP_BATCH_EXECUTE", Value: "batch_execute"},
		{PyName: "OP_SLIDE_COUNT", Value: "slide_count"},
	}
	if err := writeOpsPy(py, ops); err != nil {
		t.Fatalf("writeOpsPy: %v", err)
	}
	if err := writeOpsPyi(pyi, ops); err != nil {
		t.Fatalf("writeOpsPyi: %v", err)
	}
	_ = py.Close()
	_ = pyi.Close()

	pyData, err := os.ReadFile(pyPath)
	if err != nil {
		t.Fatalf("read py: %v", err)
	}
	pyText := string(pyData)
	if !strings.Contains(pyText, `OP_BATCH_EXECUTE = "batch_execute"`) {
		t.Fatalf("missing constant in py output: %s", pyText)
	}
	if !strings.Contains(pyText, "SUPPORTED_OPS = (") {
		t.Fatalf("missing SUPPORTED_OPS tuple in py output")
	}

	pyiData, err := os.ReadFile(pyiPath)
	if err != nil {
		t.Fatalf("read pyi: %v", err)
	}
	pyiText := string(pyiData)
	if !strings.Contains(pyiText, "OP_BATCH_EXECUTE: str") {
		t.Fatalf("missing constant annotation in pyi output: %s", pyiText)
	}
	if !strings.Contains(pyiText, "SUPPORTED_OPS_SET: frozenset[str]") {
		t.Fatalf("missing supported ops set annotation in pyi output")
	}
}
