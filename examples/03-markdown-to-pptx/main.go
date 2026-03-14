package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
)

func main() {
	outDir := filepath.Join("examples", "output")
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	if err := generateFromAsset(
		"Task 03: Markdown + Mermaid (Complex)",
		filepath.Join("examples", "assets", "03", "markdown_mermaid_complex.md"),
		filepath.Join(outDir, "03_markdown_mermaid_complex"),
		true,
	); err != nil {
		log.Fatalf("failed complex markdown sample: %v", err)
	}

	if err := generateFromAsset(
		"Task 03: Markdown Links + Gallery",
		filepath.Join("examples", "assets", "03", "markdown_links_gallery.md"),
		filepath.Join(outDir, "03_markdown_links_gallery"),
		false,
	); err != nil {
		log.Fatalf("failed links/gallery sample: %v", err)
	}
}

func generateFromAsset(title, assetMarkdownPath, outputBase string, exportPDFEnabled bool) error {
	markdownPath := outputBase + ".md"
	pptxPath := outputBase + ".pptx"
	pdfPath := outputBase + ".pdf"

	markdownBytes, err := os.ReadFile(assetMarkdownPath)
	if err != nil {
		return err
	}

	// Keep a copy in examples/output for generated artifact inspection.
	if err := os.WriteFile(markdownPath, markdownBytes, 0o600); err != nil {
		return err
	}

	slides, err := pptx.SlidesFromMarkdownFile(assetMarkdownPath)
	if err != nil {
		return err
	}

	deck, err := pptx.CreateWithSlides(title, slides)
	if err != nil {
		return err
	}
	if err := os.WriteFile(pptxPath, deck, 0o600); err != nil {
		return err
	}

	if exportPDFEnabled {
		if err := export.PDFWithOptions(
			title,
			slides,
			pdfPath,
			export.PDFOptions{Driver: export.PDFDriverNative},
		); err != nil {
			return err
		}
		log.Printf("Generated pdf (native): %s", pdfPath)
	}

	log.Printf("Generated markdown: %s", markdownPath)
	log.Printf("Generated pptx: %s", pptxPath)
	return nil
}
