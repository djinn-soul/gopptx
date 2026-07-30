package pptxxml

import (
	"archive/zip"
	"bytes"
	"testing"
)

// Issue #131: the package writer iterated its part maps, so [Content_Types].xml
// landed at a random position — usually last. OPC requires it first, and a
// consumer that reads the package as a stream cannot type any part until it has
// been seen. The random order also made the output non-reproducible.

func writeTestPackage(t *testing.T) []*zip.File {
	t.Helper()
	pw := NewPackageWriter()
	pw.AddPart("ppt/presentation.xml", "<p:presentation/>")
	pw.AddPart("ppt/slides/slide1.xml", "<p:sld/>")
	pw.AddPart(ContentTypesPartName, "<Types/>")
	pw.AddPart("_rels/.rels", "<Relationships/>")
	pw.AddBinaryPart("ppt/media/image1.png", []byte("not really a png"))

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	if err := pw.WriteTo(zw); err != nil {
		t.Fatalf("write package: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close package: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("read package: %v", err)
	}
	return zr.File
}

func TestPackageWriter_ContentTypesIsFirstEntry(t *testing.T) {
	files := writeTestPackage(t)

	if len(files) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(files))
	}
	if files[0].Name != ContentTypesPartName {
		t.Fatalf("expected %s first, got %q", ContentTypesPartName, files[0].Name)
	}
	if files[1].Name != "_rels/.rels" {
		t.Fatalf("expected _rels/.rels second, got %q", files[1].Name)
	}

	rest := []string{files[2].Name, files[3].Name, files[4].Name}
	want := []string{"ppt/media/image1.png", "ppt/presentation.xml", "ppt/slides/slide1.xml"}
	for i, name := range want {
		if rest[i] != name {
			t.Fatalf("remaining entries not sorted: got %v, want %v", rest, want)
		}
	}
}

// An entry written without a timestamp carries an all-zero MS-DOS date, which
// decodes to a month and day of zero.
func TestPackageWriter_EntriesCarryValidDOSTimestamp(t *testing.T) {
	for _, f := range writeTestPackage(t) {
		modified := f.Modified
		if modified.Month() < 1 || modified.Day() < 1 {
			t.Fatalf("entry %q has an invalid MS-DOS date: %s", f.Name, modified)
		}
		if modified.Year() < 1980 {
			t.Fatalf("entry %q predates the MS-DOS epoch: %s", f.Name, modified)
		}
	}
}

func TestPackageWriter_OutputIsReproducible(t *testing.T) {
	build := func() []byte {
		pw := NewPackageWriter()
		pw.AddPart(ContentTypesPartName, "<Types/>")
		pw.AddPart("_rels/.rels", "<Relationships/>")
		for _, name := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
			pw.AddPart("ppt/slides/"+name+".xml", "<p:sld/>")
		}
		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)
		if err := pw.WriteTo(zw); err != nil {
			t.Fatalf("write package: %v", err)
		}
		if err := zw.Close(); err != nil {
			t.Fatalf("close package: %v", err)
		}
		return buf.Bytes()
	}

	first := build()
	for i := range 5 {
		if !bytes.Equal(first, build()) {
			t.Fatalf("package bytes differ between builds (run %d)", i+2)
		}
	}
}
