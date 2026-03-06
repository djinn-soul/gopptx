package pptx

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestPresentationParityFixture_PreservesUntouchedPackageParts(t *testing.T) {
	data, err := Create("Untouched", 1)
	if err != nil {
		t.Fatalf("create fixture deck: %v", err)
	}
	withOpaque, err := zipAddOrReplaceParts(data, map[string][]byte{
		"customXml/item99.xml": []byte(`<root><opaque attr="1">keep-me</opaque></root>`),
		"ppt/extra/opaque.bin": []byte{0x00, 0xFF, 0x01, 0x02, 0x03},
	})
	if err != nil {
		t.Fatalf("inject opaque parts: %v", err)
	}

	path := filepath.Join(t.TempDir(), "opaque-roundtrip.pptx")
	if err := os.WriteFile(path, withOpaque, 0o600); err != nil {
		t.Fatalf("write fixture deck: %v", err)
	}

	prs, err := Open(path)
	if err != nil {
		t.Fatalf("open fixture deck: %v", err)
	}
	prs.SetTitle("Updated Title")
	if err := prs.Save(); err != nil {
		_ = prs.Close()
		t.Fatalf("save fixture deck: %v", err)
	}
	_ = prs.Close()

	outData, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read saved deck: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(outData), int64(len(outData)))
	if err != nil {
		t.Fatalf("open saved zip: %v", err)
	}
	item := testutil.ReadZipFile(t, zr, "customXml/item99.xml")
	if item != `<root><opaque attr="1">keep-me</opaque></root>` {
		t.Fatalf("custom part changed unexpectedly: %s", item)
	}
	opaque := testutil.ReadZipFile(t, zr, "ppt/extra/opaque.bin")
	if len(opaque) == 0 {
		t.Fatal("opaque binary part missing after round-trip")
	}
}

func TestPresentationParityFixture_CorePropertiesReadWrite(t *testing.T) {
	data, err := Create("Seed", 1)
	if err != nil {
		t.Fatalf("create seed deck: %v", err)
	}
	coreFixture := []byte(
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" ` +
			`xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" ` +
			`xmlns:dcmitype="http://purl.org/dc/dcmitype/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">` +
			`<dc:title>Fixture Title</dc:title>` +
			`<dc:subject>Fixture Subject</dc:subject>` +
			`<dc:creator>Fixture Creator</dc:creator>` +
			`<cp:keywords>fixture,core,props</cp:keywords>` +
			`<dc:description>Fixture Description</dc:description>` +
			`<cp:lastModifiedBy>Fixture Modifier</cp:lastModifiedBy>` +
			`<cp:revision>7</cp:revision>` +
			`<dcterms:created xsi:type="dcterms:W3CDTF">2024-01-01T00:00:00Z</dcterms:created>` +
			`<dcterms:modified xsi:type="dcterms:W3CDTF">2024-01-02T00:00:00Z</dcterms:modified>` +
			`<cp:category>Fixture Category</cp:category>` +
			`<cp:contentStatus>Fixture Status</cp:contentStatus>` +
			`</cp:coreProperties>`,
	)
	withFixtureCore, err := zipAddOrReplaceParts(data, map[string][]byte{
		"docProps/core.xml": coreFixture,
	})
	if err != nil {
		t.Fatalf("inject core fixture: %v", err)
	}

	path := filepath.Join(t.TempDir(), "core-fixture.pptx")
	if err := os.WriteFile(path, withFixtureCore, 0o600); err != nil {
		t.Fatalf("write core fixture deck: %v", err)
	}

	prs, err := Open(path)
	if err != nil {
		t.Fatalf("open core fixture deck: %v", err)
	}
	if prs.Title() != "Fixture Title" || prs.Subject() != "Fixture Subject" || prs.Author() != "Fixture Creator" {
		_ = prs.Close()
		t.Fatalf("fixture core values not loaded correctly")
	}
	prs.SetKeywords("updated,keywords")
	prs.SetRevision("8")
	prs.SetDescription(`Updated & Reviewed <ok>`)
	if err := prs.Save(); err != nil {
		_ = prs.Close()
		t.Fatalf("save core fixture deck: %v", err)
	}
	_ = prs.Close()

	outData, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read saved core fixture deck: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(outData), int64(len(outData)))
	if err != nil {
		t.Fatalf("open saved core fixture zip: %v", err)
	}
	coreXML := testutil.ReadZipFile(t, zr, "docProps/core.xml")
	needs := []string{
		"<cp:keywords>updated,keywords</cp:keywords>",
		"<cp:revision>8</cp:revision>",
		"Updated &amp; Reviewed &lt;ok&gt;",
		"<dc:title>Fixture Title</dc:title>",
		"<dc:subject>Fixture Subject</dc:subject>",
	}
	for _, needle := range needs {
		if !strings.Contains(coreXML, needle) {
			t.Fatalf("expected %q in core.xml, got: %s", needle, coreXML)
		}
	}
}

func zipAddOrReplaceParts(src []byte, updates map[string][]byte) ([]byte, error) {
	srcReader, err := zip.NewReader(bytes.NewReader(src), int64(len(src)))
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	zw := zip.NewWriter(&out)

	for _, f := range srcReader.File {
		if _, replace := updates[f.Name]; replace {
			continue
		}
		w, err := zw.Create(f.Name)
		if err != nil {
			_ = zw.Close()
			return nil, err
		}
		r, err := f.Open()
		if err != nil {
			_ = zw.Close()
			return nil, err
		}
		_, err = io.Copy(w, r)
		_ = r.Close()
		if err != nil {
			_ = zw.Close()
			return nil, err
		}
	}

	for name, content := range updates {
		w, err := zw.Create(name)
		if err != nil {
			_ = zw.Close()
			return nil, err
		}
		if _, err := w.Write(content); err != nil {
			_ = zw.Close()
			return nil, err
		}
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
