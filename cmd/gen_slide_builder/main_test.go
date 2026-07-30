package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func parseFixture(t *testing.T, src string) (*token.FileSet, *ast.File) {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "fixture.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	return fset, file
}

// Only exported SlideContent methods returning SlideContent get wrapped;
// anything else would not chain.
func TestIsChainableSelectsBuilderMethods(t *testing.T) {
	_, file := parseFixture(t, `package elements
type SlideContent struct{}
// AddBullet appends a bullet.
func (s SlideContent) AddBullet(text string) SlideContent { return s }
func (s SlideContent) Validate(i int) error { return nil }
func (s SlideContent) internal() SlideContent { return s }
func (s *SlideContent) Pointer() SlideContent { return *s }
func NotAMethod() SlideContent { return SlideContent{} }
func (s SlideContent) Two() (SlideContent, error) { return s, nil }
`)

	var got []string
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && isChainable(fn) {
			got = append(got, fn.Name.Name)
		}
	}

	if len(got) != 1 || got[0] != "AddBullet" {
		t.Errorf("chainable methods = %v, want [AddBullet]", got)
	}
}

func TestRenderParamsNamesAndForwardsArguments(t *testing.T) {
	fset, file := parseFixture(t, `package elements
type SlideContent struct{}
func (s SlideContent) A(text string, style Style) SlideContent { return s }
func (s SlideContent) B(runs ...Run) SlideContent { return s }
func (s SlideContent) C(int) SlideContent { return s }
func (s SlideContent) D() SlideContent { return s }
`)

	want := map[string][2]string{
		"A": {"text string, style Style", "text, style"},
		"B": {"runs ...Run", "runs..."},
		"C": {"arg0 int", "arg0"},
		"D": {"", ""},
	}

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || !isChainable(fn) {
			continue
		}
		params, args := renderParams(fset, fn.Type.Params, map[string]bool{})
		expected, known := want[fn.Name.Name]
		if !known {
			continue
		}
		if params != expected[0] {
			t.Errorf("%s params = %q, want %q", fn.Name.Name, params, expected[0])
		}
		if args != expected[1] {
			t.Errorf("%s args = %q, want %q", fn.Name.Name, args, expected[1])
		}
	}
}

// A chunk must import exactly the packages its own signatures mention, so that
// splitting the output does not produce unused or missing imports.
func TestImportsForChunkIsScopedToThatChunk(t *testing.T) {
	available := map[string]string{
		"charts": "example.com/charts",
		"shapes": "example.com/shapes",
	}
	chunk := []method{{Name: "WithChart", Quals: []string{"charts"}}}

	got := importsFor(chunk, available)
	if len(got) != 1 || got[0] != "example.com/charts" {
		t.Errorf("imports = %v, want [example.com/charts]", got)
	}
}

func TestChunkMethodsRespectsMaxPerFile(t *testing.T) {
	methods := make([]method, maxMethodsPerFile*2+3)
	chunks := chunkMethods(methods)

	if len(chunks) != 3 {
		t.Fatalf("got %d chunks, want 3", len(chunks))
	}
	for i, chunk := range chunks {
		if len(chunk) > maxMethodsPerFile {
			t.Errorf("chunk %d has %d methods, over the %d limit",
				i, len(chunk), maxMethodsPerFile)
		}
	}
}

func TestRenderChunkEmitsPointerReceiverDelegation(t *testing.T) {
	out, err := renderChunk([]method{{
		Name:   "AddBullet",
		Params: "text string",
		Args:   "text",
		Doc:    "appends a bullet.",
	}}, nil)
	if err != nil {
		t.Fatalf("renderChunk: %v", err)
	}

	for _, want := range []string{
		"DO NOT EDIT",
		"func (b *SlideBuilder) AddBullet(text string) *SlideBuilder {",
		"b.content = b.content.AddBullet(text)",
		"return b",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("generated chunk missing %q:\n%s", want, out)
		}
	}
}

// The committed builder must match what the current SlideContent methods
// produce, mirroring `task check:generated` at the unit level.
func TestGeneratedBuilderMatchesSource(t *testing.T) {
	root := repoRoot(t)
	dir := filepath.Join(root, "pkg", "pptx", "elements")

	methods, available, err := collect(dir)
	if err != nil {
		t.Fatalf("collect: %v", err)
	}
	if len(methods) == 0 {
		t.Fatal("no chainable SlideContent methods found")
	}

	for index, chunk := range chunkMethods(methods) {
		path := chunkPath(filepath.Join(dir, "slide_builder_gen.go"), index+1)
		committed, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Fatalf("read %s: %v", path, readErr)
		}
		want, renderErr := renderChunk(chunk, importsFor(chunk, available))
		if renderErr != nil {
			t.Fatalf("renderChunk: %v", renderErr)
		}
		if string(committed) != want {
			t.Errorf("%s is stale; run 'go generate ./...' and commit the result",
				filepath.Base(path))
		}
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}
