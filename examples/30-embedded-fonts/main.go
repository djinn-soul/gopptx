package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/fonts"
	logx "github.com/djinn-soul/gopptx/pkg/stdlog"
)

func main() {
	// 1. Prepare dummy font data (simulating a .ttf or .otf file)
	fontData := make([]byte, 1024)
	for i := range fontData {
		fontData[i] = byte(i % 256)
	}

	// 2. Define a GUID for font obfuscation.
	// In a real scenario, this would be the presentation GUID.
	guid := "7A7E4A1C-4B1C-4B1C-4B1C-4B1C4B1C4B1C"

	// 3. Obfuscate the font data
	obfuscated := fonts.ObfuscateFont(fontData, guid)

	// 4. Create the EmbeddedFont entry
	embeddedFont := fonts.New("CustomFont", fonts.StyleRegular, obfuscated).
		WithCharset(fonts.CharsetAnsi).
		WithPanose("020B0604020202020204").
		WithPitchFamily(0x22)

	// 5. Create Metadata with embedded fonts
	meta := pptx.Metadata{
		Metadata: common.Metadata{
			Title: "Embedded Fonts Example",
		},
		EmbeddedFonts: []fonts.EmbeddedFont{*embeddedFont},
	}

	// 6. Create slide content
	slide := pptx.NewSlide("Slide with Embedded Font").
		AddBullet("This presentation has an embedded font: 'CustomFont'")

	// 7. Generate PPTX
	data, err := pptx.CreateWithMetadata(meta, []pptx.SlideContent{slide})
	if err != nil {
		log.Fatalf("failed to create presentation: %v", err)
	}

	// 8. Save the presentation to disk
	outDir := "examples/output"
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}
	outputPath := filepath.Join(outDir, "30_embedded_fonts.pptx")
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		log.Fatalf("failed to save presentation: %v", err)
	}

	logx.Printf("Successfully created %s with embedded font.\n", outputPath)
}
