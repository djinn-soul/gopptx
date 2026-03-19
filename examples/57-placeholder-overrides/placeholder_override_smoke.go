package main

import (
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	outputDir  = "examples/output"
	outputFile = "57_placeholder_override_smoke.pptx"
)

func main() {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	slide := pptx.NewSlide("Placeholder Override Smoke").
		WithPlaceholderText(1, "Original body placeholder text").
		WithPlaceholderOverride(
			shapes.PlaceholderTarget{Type: "body", Index: 1},
			shapes.PlaceholderOverrideOptions{
				X:  ptrLength(styling.Inches(1.0)),
				Y:  ptrLength(styling.Inches(2.0)),
				CX: ptrLength(styling.Inches(4.5)),
				CY: ptrLength(styling.Inches(2.2)),
				TextStyle: &shapes.PlaceholderTextStyle{
					Font:   strPtr("Calibri"),
					SizePt: intPtr(30),
					Bold:   boolPtr(true),
					Color:  strPtr("0B5FA5"),
				},
			},
		)

	data, err := pptx.CreateWithSlides("Placeholder Override Smoke", []pptx.SlideContent{slide})
	if err != nil {
		log.Fatalf("failed to create presentation: %v", err)
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		log.Fatalf("failed to write presentation: %v", err)
	}

	log.Printf("generated %s", outputPath)
}

func ptrLength(v styling.Length) *styling.Length { return &v }
func intPtr(v int) *int                          { return &v }
func boolPtr(v bool) *bool                       { return &v }
func strPtr(v string) *string                    { return &v }
