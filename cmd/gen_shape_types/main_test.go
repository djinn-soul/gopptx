package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestCategoryFromFilename(t *testing.T) {
	cases := map[string]string{
		"shape_types_basic.go":         "basic",
		"shape_types_rect_variants.go": "rect_variants",
		"unrelated.go":                 "",
		"shape_types_basic.txt":        "",
	}
	for in, want := range cases {
		if got := categoryFromFilename(in); got != want {
			t.Errorf("categoryFromFilename(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseShapeIdents(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "sample.go")
	src := `package shapes
const (
	ShapeTypeFoo = "foo"
	OtherConst = "bar"
	ShapeTypeBar = "baz"
)`
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	idents, err := parseShapeIdents(path)
	if err != nil {
		t.Fatalf("parseShapeIdents: %v", err)
	}
	if _, ok := idents["ShapeTypeFoo"]; !ok {
		t.Errorf("missing ShapeTypeFoo: %+v", idents)
	}
	if _, ok := idents["ShapeTypeBar"]; !ok {
		t.Errorf("missing ShapeTypeBar: %+v", idents)
	}
	if _, ok := idents["OtherConst"]; ok {
		t.Errorf("unexpected non-ShapeType ident kept: %+v", idents)
	}
}

// TestGeneratorDrift regenerates shape_types.go into a temp file and ensures it
// matches the checked-in copy. Fails if someone edited the generated file
// by hand or forgot to re-run `go generate ./pkg/pptx/`.
func TestGeneratorDrift(t *testing.T) {
	root := repoRoot(t)
	shapesDir := filepath.Join(root, "pkg", "pptx", "shapes")
	aliases := filepath.Join(root, "pkg", "pptx", "shape_aliases.go")
	committed := filepath.Join(root, "pkg", "pptx", "shape_types.go")

	tmpOut := filepath.Join(t.TempDir(), "shape_types.go")

	skip, err := parseShapeIdents(aliases)
	if err != nil {
		t.Fatalf("parse aliases: %v", err)
	}
	groups, err := loadGroups(shapesDir, skip)
	if err != nil {
		t.Fatalf("load groups: %v", err)
	}
	if err := writeOutput(tmpOut, groups); err != nil {
		t.Fatalf("write output: %v", err)
	}

	gotBytes, err := os.ReadFile(tmpOut)
	if err != nil {
		t.Fatalf("read generated: %v", err)
	}
	wantBytes, err := os.ReadFile(committed)
	if err != nil {
		t.Fatalf("read committed: %v", err)
	}
	gotNorm := strings.ReplaceAll(string(gotBytes), "\r\n", "\n")
	wantNorm := strings.ReplaceAll(string(wantBytes), "\r\n", "\n")
	if gotNorm != wantNorm {
		t.Fatalf("shape_types.go drift: generator output differs from checked-in file. "+
			"Re-run `go generate ./pkg/pptx/`.\nlen(got)=%d len(want)=%d",
			len(gotNorm), len(wantNorm))
	}
}
