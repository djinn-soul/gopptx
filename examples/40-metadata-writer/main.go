package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
)

func main() {
	const outputDir = "examples/output"
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("Failed to create output dir: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "gopptx-metadata-writer-*")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Printf("warning: failed to remove temp dir %s: %v", tmpDir, err)
		}
	}()

	inputFile := filepath.Join(tmpDir, "metadata_input.pptx")
	outputFile := filepath.Join(outputDir, "40_metadata_output.pptx")

	// 1. Create a minimal valid PPTX file
	fmt.Printf("Generating minimal PPTX: %s...\n", inputFile)
	if err := createMinimalPPTX(inputFile); err != nil {
		log.Fatalf("Failed to create minimal PPTX: %v", err)
	}
	defer func() {
		// optional cleanup
		if err := os.Remove(inputFile); err != nil && !os.IsNotExist(err) {
			log.Printf("warning: failed to remove input file %s: %v", inputFile, err)
		}
	}()

	// 2. Open it
	fmt.Printf("Opening %s...\n", inputFile)
	ppt, err := editor.OpenPresentationEditor(inputFile)
	if err != nil {
		log.Fatalf("Failed to open presentation: %v", err)
	}
	defer func() { _ = ppt.Close() }()

	// 3. Check initial metadata
	props := ppt.GetCoreProperties()
	fmt.Printf("Initial Title: %s\n", props.Title)
	if props.Title != "Initial Title" {
		log.Fatalf("Expected 'Initial Title', got '%s'", props.Title)
	}

	// 4. Update metadata
	fmt.Println("Updating metadata...")
	newProps := common.CoreProperties{
		Title:       "Updated Title",
		Subject:     "Updated Subject",
		Creator:     "Updated Creator",
		Description: "Updated Description",
		Keywords:    "test, metadata",
	}
	ppt.SetCoreProperties(newProps)

	// 5. Save
	fmt.Printf("Saving to %s...\n", outputFile)
	if err := ppt.Save(outputFile); err != nil {
		log.Fatalf("Failed to save: %v", err)
	}

	// 6. Verify output
	verifyOutput(outputFile)

	fmt.Println("Done! Smoke test passed.")
}

func verifyOutput(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open output file: %v", err)
	}
	defer func() { _ = f.Close() }()

	fi, _ := f.Stat()
	z, err := zip.NewReader(f, fi.Size())
	if err != nil {
		log.Fatalf("Failed to open zip: %v", err)
	}

	foundCore := false
	for _, f := range z.File {
		if f.Name == "docProps/core.xml" {
			foundCore = true
			rc, _ := f.Open()
			content := make([]byte, f.UncompressedSize64)
			_, _ = rc.Read(content)
			_ = rc.Close()
			s := string(content)
			if !contains(s, "Updated Title") {
				log.Fatalf("Output missing updated title")
			}
			if !contains(s, "Updated Description") {
				log.Fatalf("Output missing updated description")
			}
		}
	}
	if !foundCore {
		log.Fatalf("Output missing docProps/core.xml")
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func createMinimalPPTX(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	z := zip.NewWriter(f)
	defer func() { _ = z.Close() }()

	// [Content_Types].xml
	ct := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>
  <Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>
</Types>`
	if err := writeZipFile(z, "[Content_Types].xml", ct); err != nil {
		return err
	}

	// _rels/.rels
	pkgRels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
</Relationships>`
	if err := writeZipFile(z, "_rels/.rels", pkgRels); err != nil {
		return err
	}

	// ppt/presentation.xml
	pres := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:sldIdLst/></p:presentation>`
	if err := writeZipFile(z, "ppt/presentation.xml", pres); err != nil {
		return err
	}

	// ppt/_rels/presentation.xml.rels (required)
	presRels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
	<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`
	if err := writeZipFile(z, "ppt/_rels/presentation.xml.rels", presRels); err != nil {
		return err
	}

	// docProps/core.xml
	core := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:dcmitype="http://purl.org/dc/dcmitype/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<dc:title>Initial Title</dc:title>
<dc:creator>Initial Creator</dc:creator>
<cp:keywords>initial</cp:keywords>
</cp:coreProperties>`
	if err := writeZipFile(z, "docProps/core.xml", core); err != nil {
		return err
	}

	return nil
}

func writeZipFile(z *zip.Writer, name, content string) error {
	w, err := z.Create(name)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(content))
	return err
}
