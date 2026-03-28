// examples/52-legacy-ppt-interop/main.go demonstrates how gopptx handles legacy
// binary .ppt files and the recommended LibreOffice conversion workflow.
//
// gopptx supports the .pptx (OpenXML/ZIP) format only. Opening a binary .ppt
// (Compound File Binary / CFB) returns a clear error with conversion guidance.
//
// Run with: go run ./examples/52-legacy-ppt-interop/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "52_legacy_ppt_interop.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// --- Demonstrate legacy .ppt detection ---
	// Create a temporary file whose first bytes match the CFB magic signature
	// (D0 CF 11 E0 A1 B1 1A E1) used by binary .ppt files.
	tmpDir, err := os.MkdirTemp("", "gopptx-legacy-*")
	if err != nil {
		return fmt.Errorf("mktemp: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	legacyPath := filepath.Join(tmpDir, "legacy.ppt")
	cfbMagic := []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1, 0x00, 0x00}
	if writeErr := os.WriteFile(legacyPath, cfbMagic, 0o600); writeErr != nil {
		return fmt.Errorf("write legacy fixture: %w", writeErr)
	}

	_, openErr := pptx.OpenPresentationEditor(legacyPath)
	if openErr != nil {
		log.Printf("Expected error for legacy .ppt: %v", openErr)
		log.Printf("Workaround: convert to .pptx with LibreOffice:")
		log.Printf("  soffice --headless --convert-to pptx legacy.ppt")
	} else {
		log.Printf("NOTE: legacy .ppt was not detected as CFB on this build")
	}

	// --- Build a documentation presentation describing the interop strategy ---
	slides := []pptx.SlideContent{
		pptx.NewSlide("Legacy PPT Interop").
			AddBullet("gopptx supports .pptx (OpenXML) format only").
			AddBullet("Legacy .ppt (binary CFB) files must be converted first").
			AddBullet("gopptx detects CFB magic bytes and returns a descriptive error"),

		pptx.NewSlide("Conversion with LibreOffice").
			AddNumbered("Install LibreOffice (https://www.libreoffice.org)").
			AddNumbered("Run: soffice --headless --convert-to pptx file.ppt").
			AddNumbered("The resulting file.pptx can be opened with gopptx"),

		pptx.NewSlide("CFB Detection").
			AddBullet("Binary .ppt uses the Compound File Binary (CFB) container").
			AddBullet("CFB magic: D0 CF 11 E0 A1 B1 1A E1 (first 8 bytes)").
			AddBullet("gopptx checks the magic bytes before attempting ZIP/OpenXML parse").
			AddBullet("Returns a clear error message with conversion instructions"),

		pptx.NewSlide("Format Comparison").
			AddBullet(".ppt  — binary CFB, Office 97-2003, not supported by gopptx").
			AddBullet(".pptx — OpenXML ZIP, Office 2007+, fully supported").
			AddBullet(".pptm — macro-enabled PPTX, supported (VBA macros preserved)"),
	}

	data, err2 := pptx.CreateWithSlides("Legacy PPT Interop Guide", slides)
	if err2 != nil {
		return fmt.Errorf("create: %w", err2)
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err = os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
