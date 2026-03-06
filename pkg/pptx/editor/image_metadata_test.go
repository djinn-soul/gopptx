package editor

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestGetImageMetadataIncludesContentTypeAndHash(t *testing.T) {
	pptxPath := createPictureFixturePPTX(t, testutil.TinyPNG())

	ed, err := OpenPresentationEditor(pptxPath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	meta, err := ed.GetImageMetadata(0, 2)
	if err != nil {
		t.Fatalf("get image metadata: %v", err)
	}
	if meta.Width <= 0 || meta.Height <= 0 {
		t.Fatalf("expected non-zero dimensions, got %dx%d", meta.Width, meta.Height)
	}
	if meta.Format != "png" {
		t.Fatalf("expected png format, got %q", meta.Format)
	}
	if meta.ContentType != "image/png" {
		t.Fatalf("expected image/png content-type, got %q", meta.ContentType)
	}
	expectedHash := sha256.Sum256(testutil.TinyPNG())
	if meta.Hash != hex.EncodeToString(expectedHash[:]) {
		t.Fatalf("unexpected hash: got %q", meta.Hash)
	}
}

func createPictureFixturePPTX(t *testing.T, imageData []byte) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "image-fixture.pptx")
	files := map[string][]byte{
		"[Content_Types].xml": []byte(
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
				`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
				`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
				`<Default Extension="xml" ContentType="application/xml"/>` +
				`<Default Extension="png" ContentType="image/png"/>` +
				`<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>` +
				`<Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>` +
				`</Types>`,
		),
		"_rels/.rels": []byte(
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
				`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
				`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>` +
				`</Relationships>`,
		),
		"ppt/presentation.xml": []byte(
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
				`<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
				`<p:sldIdLst><p:sldId id="256" r:id="rId1"/></p:sldIdLst></p:presentation>`,
		),
		"ppt/_rels/presentation.xml.rels": []byte(
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
				`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
				`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/>` +
				`</Relationships>`,
		),
		"ppt/slides/slide1.xml": []byte(
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
				`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
				`<p:cSld><p:spTree>` +
				`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
				`<p:pic><p:nvPicPr><p:cNvPr id="2" name="Picture 2"/><p:cNvPicPr/><p:nvPr/></p:nvPicPr>` +
				`<p:blipFill><a:blip r:embed="rId2"/><a:stretch><a:fillRect/></a:stretch></p:blipFill>` +
				`<p:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="1000" cy="1000"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr>` +
				`</p:pic></p:spTree></p:cSld></p:sld>`,
		),
		"ppt/slides/_rels/slide1.xml.rels": []byte(
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
				`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
				`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="../media/image1.png"/>` +
				`</Relationships>`,
		),
		"ppt/media/image1.png": imageData,
	}
	if err := writeZipFixtureBytes(path, files); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

func writeZipFixtureBytes(path string, files map[string][]byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	zw := zip.NewWriter(f)
	for name, content := range files {
		w, createErr := zw.Create(name)
		if createErr != nil {
			return createErr
		}
		if _, writeErr := bytes.NewBuffer(content).WriteTo(w); writeErr != nil {
			return writeErr
		}
	}
	if err := zw.Close(); err != nil {
		return err
	}
	return nil
}
