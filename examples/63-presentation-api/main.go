// examples/63-presentation-api demonstrates core Presentation API operations:
// creating new presentations, opening existing ones, saving (in-place and as copy),
// and working with core properties (title, author, subject, keywords, etc.).
//
// Run with: go run ./examples/63-presentation-api/main.go
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
	outputFile = "63_presentation_api.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err := buildAndSavePresentation(outputPath); err != nil {
		return err
	}

	if err := updateAndVerifyPresentation(outputPath); err != nil {
		return err
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

func buildAndSavePresentation(outputPath string) error {
	// --- Part 1: Create a presentation using PresentationBuilder ---
	builder := pptx.NewPresentationBuilder("Presentation API Demo").
		WithMetadata(pptx.Metadata{
			Metadata: pptx.MetadataFields{
				Title:       "Presentation API Demo",
				Description: "Demonstrates creation and metadata management",
			},
		}).
		AddBulletSlide("Overview", []string{
			"Create presentations with NewPresentationBuilder",
			"Open existing files with pptx.Open",
			"Save with Save() or SaveAs()",
			"Manage core properties: title, author, subject, keywords",
		}).
		AddBulletSlide("Core Properties", []string{
			"Title, Subject, Creator/Author",
			"Keywords, Description, Category",
			"ContentStatus, Revision",
			"Created and Modified timestamps",
		})

	if err := builder.WriteToFile(outputPath); err != nil {
		return fmt.Errorf("write initial file: %w", err)
	}
	log.Printf("Created initial presentation: %s\n", outputPath)

	// --- Part 2: Open the file and read core properties ---
	prs, err := pptx.Open(outputPath)
	if err != nil {
		return fmt.Errorf("open presentation: %w", err)
	}
	defer prs.Close()

	log.Printf("Slide count: %d\n", prs.SlideCount())
	log.Printf("Title: %q\n", prs.Title())
	log.Printf("Creator: %q\n", prs.Creator())

	return updateCoreProperties(prs)
}

func updateCoreProperties(prs *pptx.Presentation) error {
	// --- Part 3: Update core properties ---
	prs.SetTitle("Updated Presentation API Demo")
	prs.SetAuthor("Go Developer")
	prs.SetCreator("Go Developer")
	prs.SetSubject("gopptx Presentation API")
	prs.SetKeywords("gopptx, pptx, go, api")
	prs.SetDescription("A thorough demonstration of the presentation creation and metadata API.")
	prs.SetCategory("Developer Examples")
	prs.SetContentStatus("Final")
	prs.SetRevision("2")

	if err := prs.Save(); err != nil {
		return fmt.Errorf("save presentation: %w", err)
	}
	log.Printf("Saved with updated metadata\n")

	// --- Part 4: SaveAs to produce a copy ---
	copyPath := filepath.Join(outputDir, "63_presentation_api_copy.pptx")
	prs.SetTitle("Copy of Presentation API Demo")
	if err := prs.SaveAs(copyPath); err != nil {
		return fmt.Errorf("save as copy: %w", err)
	}
	log.Printf("Saved copy to: %s\n", copyPath)
	return nil
}

func updateAndVerifyPresentation(outputPath string) error {
	// --- Part 5: Verify core properties survive round-trip ---
	prs2, err := pptx.Open(outputPath)
	if err != nil {
		return fmt.Errorf("reopen for verification: %w", err)
	}
	defer prs2.Close()

	log.Printf("Verified title: %q\n", prs2.Title())
	log.Printf("Verified author: %q\n", prs2.Author())
	log.Printf("Verified keywords: %q\n", prs2.Keywords())
	log.Printf("Verified category: %q\n", prs2.Category())
	log.Printf("Verified revision: %q\n", prs2.Revision())

	// --- Part 6: CoreProperties bulk read/write ---
	props := prs2.CoreProperties()
	log.Printf("CoreProperties.Title: %q\n", props.Title)
	log.Printf("CoreProperties.Creator: %q\n", props.Creator)

	props.Title = "Final Title via CoreProperties"
	props.Revision = "3"
	prs2.SetCoreProperties(props)
	if err := prs2.Save(); err != nil {
		return fmt.Errorf("save after CoreProperties update: %w", err)
	}

	validateAndLogIssues(prs2)
	logSlideSizes()
	return nil
}

func validateAndLogIssues(prs *pptx.Presentation) {
	// --- Part 7: Validate ---
	issues := prs.Validate()
	if len(issues) == 0 {
		log.Printf("Validation passed with no issues.\n")
	} else {
		log.Printf("Validation issues: %d\n", len(issues))
		for _, issue := range issues {
			log.Printf("  - %s\n", issue.Description)
		}
	}
}

func logSlideSizes() {
	// --- Part 8: SlideSize helpers ---
	size4x3 := pptx.SlideSize4x3()
	size16x9 := pptx.SlideSize16x9()
	log.Printf("4:3 size: %dx%d EMU\n", size4x3.Width, size4x3.Height)
	log.Printf("16:9 size: %dx%d EMU\n", size16x9.Width, size16x9.Height)
}
