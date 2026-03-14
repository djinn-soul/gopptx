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

	assetMarkdownPath := filepath.Join("examples", "assets", "03", "markdown_mermaid_complex.md")
	markdownPath := filepath.Join(outDir, "03_markdown_mermaid_complex.md")
	pptxPath := filepath.Join(outDir, "03_markdown_mermaid_complex.pptx")
	pdfPath := filepath.Join(outDir, "03_markdown_mermaid_complex.pdf")

	markdownBytes, err := os.ReadFile(assetMarkdownPath)
	if err != nil {
		log.Fatalf("failed to read markdown asset %s: %v", assetMarkdownPath, err)
	}
	// Keep a copy in examples/output for generated artifact inspection.
	if err := os.WriteFile(markdownPath, markdownBytes, 0o600); err != nil {
		log.Fatalf("failed to write markdown output copy: %v", err)
	}

	slides, err := pptx.SlidesFromMarkdownFile(assetMarkdownPath)
	if err != nil {
		log.Fatalf("failed to parse markdown: %v", err)
	}

	deck, err := pptx.CreateWithSlides("Task 03: Markdown + Mermaid (Complex)", slides)
	if err != nil {
		log.Fatalf("failed to create pptx: %v", err)
	}
	if err := os.WriteFile(pptxPath, deck, 0o600); err != nil {
		log.Fatalf("failed to write pptx: %v", err)
	}

	if err := export.PDFWithOptions(
		"Task 03: Markdown + Mermaid (Complex)",
		slides,
		pdfPath,
		export.PDFOptions{Driver: export.PDFDriverNative},
	); err != nil {
		log.Fatalf("failed to export native pdf: %v", err)
	}

	log.Printf("Generated markdown: %s", markdownPath)
	log.Printf("Generated pptx: %s", pptxPath)
	log.Printf("Generated pdf (native): %s", pdfPath)
}
