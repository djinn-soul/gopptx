package pptxxml_test

import (
	"archive/zip"
	"bytes"
	"io"
	"strconv"
	"testing"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

func TestPackageWriterWritesTextAndBinaryParts(t *testing.T) {
	pw := pptxxml.NewPackageWriter()
	pw.AddPart("doc.txt", "hello")
	pw.AddBinaryPart("bin.dat", []byte{0x01, 0x02, 0x03})

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	if err := pw.WriteTo(zw); err != nil {
		t.Fatalf("write package: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}

	got := map[string][]byte{}
	for _, f := range zr.File {
		rc, openErr := f.Open()
		if openErr != nil {
			t.Fatalf("open file %s: %v", f.Name, openErr)
		}
		data, readErr := io.ReadAll(rc)
		_ = rc.Close()
		if readErr != nil {
			t.Fatalf("read file %s: %v", f.Name, readErr)
		}
		got[f.Name] = data
	}

	if string(got["doc.txt"]) != "hello" {
		t.Fatalf("doc.txt mismatch: got %q", string(got["doc.txt"]))
	}
	if !bytes.Equal(got["bin.dat"], []byte{0x01, 0x02, 0x03}) {
		t.Fatalf("bin.dat mismatch: got %v", got["bin.dat"])
	}
}

func BenchmarkPackageWriter(b *testing.B) {
	paths := make([]string, 200)
	for i := range paths {
		paths[i] = strconv.Itoa(i)
	}
	payload := bytes.Repeat([]byte("x"), 4096)
	b.ReportAllocs()
	for b.Loop() {
		pw := pptxxml.NewPackageWriter()
		for _, s := range paths {
			pw.AddPart("ppt/slides/slide"+s+".xml", "<p:sld/>")
			pw.AddBinaryPart("ppt/media/media"+s+".bin", payload)
		}
		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)
		if err := pw.WriteTo(zw); err != nil {
			b.Fatalf("write package: %v", err)
		}
		if err := zw.Close(); err != nil {
			b.Fatalf("close zip: %v", err)
		}
	}
}

func BenchmarkPackageWriterAddOnly(b *testing.B) {
	paths := make([]string, 200)
	for i := range paths {
		paths[i] = strconv.Itoa(i)
	}
	payload := bytes.Repeat([]byte("x"), 4096)

	b.ReportAllocs()
	for b.Loop() {
		pw := pptxxml.NewPackageWriter()
		for _, s := range paths {
			pw.AddPart("ppt/slides/slide"+s+".xml", "<p:sld/>")
			pw.AddBinaryPart("ppt/media/media"+s+".bin", payload)
		}
	}
}

func BenchmarkPackageWriterWriteToOnlyDiscard(b *testing.B) {
	paths := make([]string, 200)
	for i := range paths {
		paths[i] = strconv.Itoa(i)
	}
	payload := bytes.Repeat([]byte("x"), 4096)
	pw := pptxxml.NewPackageWriter()
	for _, s := range paths {
		pw.AddPart("ppt/slides/slide"+s+".xml", "<p:sld/>")
		pw.AddBinaryPart("ppt/media/media"+s+".bin", payload)
	}

	b.ReportAllocs()
	for b.Loop() {
		zw := zip.NewWriter(io.Discard)
		if err := pw.WriteTo(zw); err != nil {
			b.Fatalf("write package: %v", err)
		}
		if err := zw.Close(); err != nil {
			b.Fatalf("close zip: %v", err)
		}
	}
}
