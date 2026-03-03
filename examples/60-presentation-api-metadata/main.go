// 60-presentation-api-metadata demonstrates the Presentation API
// for opening existing presentations and manipulating their metadata.
//
// This API is designed to be similar to python-pptx's Presentation(pptx_path)
// constructor.
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("presentation metadata example failed: %v", err)
	}
}

func run() error {
	examplePath := "example.pptx"

	if err := createSamplePresentation(examplePath); err != nil {
		return fmt.Errorf("failed to create sample presentation: %w", err)
	}
	defer os.Remove(examplePath)

	printLine("=== Example 1: Opening a presentation ===")
	prs, err := pptx.Open(examplePath)
	if err != nil {
		return fmt.Errorf("failed to open presentation: %w", err)
	}
	defer prs.Close()

	printFmtf("Opened presentation with %d slides\n", prs.SlideCount())
	printFmtf("Current title: %s\n", prs.Title())
	printFmtf("Current author: %s\n", prs.Creator())

	// Example 2: Reading metadata
	printLine("\n=== Example 2: Reading all metadata ===")
	readAllMetadata(prs)

	if err := updateMetadata(prs); err != nil {
		return err
	}
	if err := reopenAndVerify(examplePath); err != nil {
		return err
	}
	if err := saveCopyExample(examplePath); err != nil {
		return err
	}
	if err := updateCoreProperties(examplePath); err != nil {
		return err
	}
	return validatePresentation(examplePath)
}

func updateMetadata(prs *pptx.Presentation) error {
	printLine("\n=== Example 3: Modifying metadata ===")
	prs.SetTitle("Updated Presentation Title")
	prs.SetAuthor("Jane Doe")
	prs.SetCreator("Jane Doe") // Author is an alias for Creator
	prs.SetSubject("Updated Subject")
	prs.SetKeywords("presentation, go, pptx, metadata")
	prs.SetDescription("This presentation was updated using the Presentation API")
	prs.SetCategory("Technical Documentation")
	prs.SetContentStatus("Draft")
	prs.SetRevision("2")
	if err := prs.Save(); err != nil {
		return fmt.Errorf("failed to save presentation: %w", err)
	}
	printLine("Changes saved!")
	return nil
}

func reopenAndVerify(examplePath string) error {
	printLine("\n=== Example 4: Reopening to verify persistence ===")
	prs, err := pptx.Open(examplePath)
	if err != nil {
		return fmt.Errorf("failed to reopen presentation: %w", err)
	}
	defer prs.Close()

	printFmtf("Title: %s\n", prs.Title())
	printFmtf("Author: %s\n", prs.Author())
	printFmtf("Keywords: %s\n", prs.Keywords())
	printFmtf("Category: %s\n", prs.Category())
	return nil
}

func saveCopyExample(examplePath string) error {
	printLine("\n=== Example 5: SaveAs to create a copy ===")
	copyPath := "example_copy.pptx"
	defer os.Remove(copyPath)

	prs, err := pptx.Open(examplePath)
	if err != nil {
		return fmt.Errorf("failed to open presentation for copy: %w", err)
	}
	defer prs.Close()

	prs.SetTitle("Copy of the Presentation")
	if err := prs.SaveAs(copyPath); err != nil {
		return fmt.Errorf("failed to save copy: %w", err)
	}
	printFmtf("Saved copy to %s\n", copyPath)

	original, err := pptx.Open(examplePath)
	if err != nil {
		return fmt.Errorf("failed to verify original presentation: %w", err)
	}
	defer original.Close()

	if original.Title() != "Updated Presentation Title" {
		printFmtf("Warning: Original presentation was modified\n")
	}
	printFmtf("Original title unchanged: %s\n", original.Title())
	return nil
}

func updateCoreProperties(examplePath string) error {
	printLine("\n=== Example 6: Using CoreProperties directly ===")
	prs, err := pptx.Open(examplePath)
	if err != nil {
		return fmt.Errorf("failed to open presentation: %w", err)
	}
	defer prs.Close()

	props := prs.CoreProperties()
	printFmtf("CoreProperties.Title: %s\n", props.Title)
	printFmtf("CoreProperties.Creator: %s\n", props.Creator)
	printFmtf("CoreProperties.Revision: %s\n", props.Revision)

	props.Title = "Final Title"
	props.Revision = "3"
	prs.SetCoreProperties(props)

	if err := prs.Save(); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	verified, err := pptx.Open(examplePath)
	if err != nil {
		return fmt.Errorf("failed to verify: %w", err)
	}
	defer verified.Close()

	printFmtf("After core props update - Title: %s, Revision: %s\n",
		verified.Title(), verified.Revision())
	return nil
}

func validatePresentation(examplePath string) error {
	printLine("\n=== Example 7: Validate presentation ===")
	prs, err := pptx.Open(examplePath)
	if err != nil {
		return fmt.Errorf("failed to open presentation for validation: %w", err)
	}
	defer prs.Close()

	issues := prs.Validate()
	if len(issues) == 0 {
		printLine("Presentation is valid!")
	} else {
		printFmtf("Found %d validation issues:\n", len(issues))
		for _, issue := range issues {
			printFmtf("  - %s\n", issue.Description)
		}
	}

	printLine("\nAll examples completed successfully!")
	return nil
}

func createSamplePresentation(path string) error {
	data, err := pptx.Create("Sample Presentation", 3)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func readAllMetadata(prs *pptx.Presentation) {
	printFmtf("Title: %s\n", prs.Title())
	printFmtf("Subject: %s\n", prs.Subject())
	printFmtf("Creator: %s\n", prs.Creator())
	printFmtf("Author: %s\n", prs.Author()) // Alias for Creator
	printFmtf("Keywords: %s\n", prs.Keywords())
	printFmtf("Description: %s\n", prs.Description())
	printFmtf("Last Modified By: %s\n", prs.LastModifiedBy())
	printFmtf("Revision: %s\n", prs.Revision())
	printFmtf("Created: %s\n", prs.Created())
	printFmtf("Modified: %s\n", prs.Modified())
	printFmtf("Category: %s\n", prs.Category())
	printFmtf("Content Status: %s\n", prs.ContentStatus())
}

func printLine(args ...any) {
	_, _ = io.WriteString(os.Stdout, fmt.Sprintln(args...))
}

func printFmtf(format string, args ...any) {
	_, _ = io.WriteString(os.Stdout, fmt.Sprintf(format, args...))
}
