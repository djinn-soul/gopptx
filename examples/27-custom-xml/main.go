package main

import (
	"log"
	"os"
	"path/filepath"

	logx "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/customxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
)

func main() {
	outDir := "examples/output"
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		log.Fatalf("failed to create output dir: %v", err)
	}
	outputPath := filepath.Join(outDir, "27_custom_xml_smoke.pptx")
	roundTripPath := filepath.Join(outDir, "27_custom_xml_roundtrip.pptx")

	// 1. Create a blank baseline presentation
	if err := pptx.NewPresentationBuilder("Presentation").
		AddSlide(pptx.NewSlide("Slide 1")).
		WriteToFile(outputPath); err != nil {
		log.Fatalf("failed to save baseline presentation: %v", err)
	}
	logx.Printf("Created baseline presentation at %s\n", outputPath)

	// 2. Open with Editor to manipulate Custom XML
	e, err := editor.OpenPresentationEditor(outputPath)
	if err != nil {
		log.Fatalf("failed to open presentation in editor: %v", err)
	}

	// 3. Create Custom XML Parts using the Builder
	store := customxml.NewStore()

	// Structured part
	store.Add("CompanyData").
		Namespace("http://schemas.example.com/company").
		Property("Name", "Acme Corp").
		Property("ID", "12345").
		Property("Active", "true")

	// Raw part
	store.Add("").Content(`<RawSettings><Theme>Dark</Theme><Version>2</Version></RawSettings>`)

	// Inject parts into the active presentation metadata
	m := e.Metadata()
	m.CustomXML = append(m.CustomXML, store.ToCommonParts()...)

	logx.Printf("Injected %d Custom XML parts into presentation\n", len(store.ToCommonParts()))

	// 4. Save Editor changes
	if err := e.Save(roundTripPath); err != nil {
		log.Fatalf("failed to save round-trip presentation: %v", err)
	}
	logx.Printf("Saved round-trip presentation to %s\n", roundTripPath)

	// 5. Verify the parts round-tripped correctly
	e2, err := editor.OpenPresentationEditor(roundTripPath)
	if err != nil {
		log.Fatalf("failed to re-open round-trip presentation: %v", err)
	}

	parts := e2.Metadata().CustomXML
	if len(parts) != 2 {
		log.Fatalf("expected 2 Custom XML parts, got %d", len(parts))
	}
	logx.Println("Successfully verified Custom XML preservation on round-trip!")
}
