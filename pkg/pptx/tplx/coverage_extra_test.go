package tplx

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestRenderFunctions(t *testing.T) {
	// Create a dummy PPTX (ZIP)
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, _ := zw.Create("ppt/slides/slide1.xml")
	_, _ = f.Write([]byte(`<p:sp><a:t>{{name}}</a:t></p:sp>`))
	_ = zw.Close()

	tmpDir := t.TempDir()
	pptxPath := filepath.Join(tmpDir, "test.pptx")
	_ = os.WriteFile(pptxPath, buf.Bytes(), 0644)

	ctx := Context{"name": "World"}

	t.Run("Render", func(t *testing.T) {
		res, err := Render(pptxPath, ctx)
		if err != nil {
			t.Fatalf("Render failed: %v", err)
		}
		checkContent(t, res.Bytes(), "ppt/slides/slide1.xml", "World")
	})

	t.Run("RenderWithOptions", func(t *testing.T) {
		res, err := RenderWithOptions(pptxPath, ctx, Options{Strict: true})
		if err != nil {
			t.Fatalf("RenderWithOptions failed: %v", err)
		}
		checkContent(t, res.Bytes(), "ppt/slides/slide1.xml", "World")
	})

	t.Run("RenderMissingFile", func(t *testing.T) {
		_, err := Render("non-existent.pptx", ctx)
		if err == nil {
			t.Error("Expected error for missing file")
		}
	})
}

func checkContent(t *testing.T, data []byte, name, expected string) {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("Invalid ZIP: %v", err)
	}
	for _, f := range zr.File {
		if f.Name == name {
			rc, _ := f.Open()
			content, _ := io.ReadAll(rc)
			_ = rc.Close()
			if !bytes.Contains(content, []byte(expected)) {
				t.Errorf("Expected %q in %s, got %s", expected, name, string(content))
			}
			return
		}
	}
	t.Errorf("File %s not found in ZIP", name)
}

func TestSlideLoopsAndRels(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	// Slide 1 with each loop
	f1, _ := zw.Create("ppt/slides/slide1.xml")
	_, _ = f1.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
<p:sp><p:txBody><a:p><a:r><a:t>{{#each items}}{{val}}{{/each}}</a:t></a:r></a:p></p:txBody></p:sp>
</p:sld>`))

	// Rels for slide 1
	r1, _ := zw.Create("ppt/slides/_rels/slide1.xml.rels")
	_, _ = r1.Write([]byte(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`))

	_ = zw.Close()

	ctx := Context{
		"items": []Row{
			{"val": "A"},
			{"val": "B"},
		},
	}

	res, err := RenderBytes(buf.Bytes(), ctx)
	if err != nil {
		t.Fatalf("RenderBytes failed: %v", err)
	}

	// Should have slide 2 and slide 3 now, slide 1 removed
	outZr, _ := zip.NewReader(bytes.NewReader(res.Bytes()), int64(len(res.Bytes())))
	foundS2 := false
	foundS3 := false
	foundR2 := false
	foundR3 := false
	for _, f := range outZr.File {
		if f.Name == "ppt/slides/slide2.xml" {
			foundS2 = true
		}
		if f.Name == "ppt/slides/slide3.xml" {
			foundS3 = true
		}
		if f.Name == "ppt/slides/_rels/slide2.xml.rels" {
			foundR2 = true
		}
		if f.Name == "ppt/slides/_rels/slide3.xml.rels" {
			foundR3 = true
		}
	}

	if !foundS2 || !foundS3 {
		t.Errorf("Expected slide2 and slide3, got foundS2=%v, foundS3=%v", foundS2, foundS3)
	}
	if !foundR2 || !foundR3 {
		t.Errorf("Expected rels for slide2 and slide3, got foundR2=%v, foundR3=%v", foundR2, foundR3)
	}
}

func TestSlideLoopEmptyData(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f1, _ := zw.Create("ppt/slides/slide1.xml")
	_, _ = f1.Write([]byte(`{{#each items}}{{val}}{{/each}}`))
	_ = zw.Close()

	// Empty data should remove the slide
	ctx := Context{"items": []Row{}}
	res, _ := RenderBytes(buf.Bytes(), ctx)
	outZr, _ := zip.NewReader(bytes.NewReader(res.Bytes()), int64(len(res.Bytes())))
	for _, f := range outZr.File {
		if f.Name == "ppt/slides/slide1.xml" {
			t.Error("Slide1 should have been removed")
		}
	}

	// Missing key should also remove the slide
	res, _ = RenderBytes(buf.Bytes(), Context{})
	outZr, _ = zip.NewReader(bytes.NewReader(res.Bytes()), int64(len(res.Bytes())))
	for _, f := range outZr.File {
		if f.Name == "ppt/slides/slide1.xml" {
			t.Error("Slide1 should have been removed for missing key")
		}
	}
}

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		val    any
		expect bool
	}{
		{nil, false},
		{true, true},
		{false, false},
		{"", false},
		{"  ", false},
		{"hi", true},
		{0, false},
		{1, true},
		{int64(0), false},
		{int64(1), true},
		{0.0, false},
		{3.14, true},
		{[]Row{}, false},
		{[]Row{{"a": "b"}}, true},
		{[]map[string]string{}, false},
		{[]map[string]string{{"a": "b"}}, true},
		{struct{}{}, true}, // unknown type is truthy
	}

	for _, tc := range tests {
		if got := isTruthy(tc.val); got != tc.expect {
			t.Errorf("isTruthy(%v) = %v, expect %v", tc.val, got, tc.expect)
		}
	}
}

func TestToRows(t *testing.T) {
	t.Run("map[string]string", func(t *testing.T) {
		input := []map[string]string{{"a": "b"}}
		rows, ok := toRows(input)
		if !ok || len(rows) != 1 || rows[0]["a"] != "b" {
			t.Errorf("toRows failed for map[string]string")
		}
	})
	t.Run("invalid", func(t *testing.T) {
		_, ok := toRows("not a slice")
		if ok {
			t.Error("toRows should fail for string")
		}
	})
}

func TestHelpers(t *testing.T) {
	if lastSegment("a/b/c") != "c" {
		t.Errorf("lastSegment failed")
	}
	if lastSegment("abc") != "abc" {
		t.Errorf("lastSegment failed for no slash")
	}

	ss := []string{"c", "a", "b"}
	sortStrings(ss)
	if ss[0] != "a" || ss[1] != "b" || ss[2] != "c" {
		t.Errorf("sortStrings failed: %v", ss)
	}

	if trimPrefix("foobar", "foo") != "bar" {
		t.Errorf("trimPrefix failed")
	}
	if trimPrefix("foobar", "baz") != "foobar" {
		t.Errorf("trimPrefix failed for non-prefix")
	}
	if trimSuffix("foobar", "bar") != "foo" {
		t.Errorf("trimSuffix failed")
	}
	if trimSuffix("foobar", "baz") != "foobar" {
		t.Errorf("trimSuffix failed for non-suffix")
	}
}

func TestParseSlideNumber(t *testing.T) {
	if parseSlideNumber("ppt/slides/slide123.xml") != 123 {
		t.Errorf("parseSlideNumber failed")
	}
	if parseSlideNumber("ppt/slides/slideABC.xml") != 0 {
		t.Errorf("parseSlideNumber should return 0 for invalid")
	}
}
