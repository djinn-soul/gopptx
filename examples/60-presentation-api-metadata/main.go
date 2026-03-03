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
	// Example 1: Open an existing presentation
	examplePath := "example.pptx"

	// First create a sample presentation
	if err := createSamplePresentation(examplePath); err != nil {
		log.Fatalf("failed to create sample presentation: %v", err)
	}
	defer os.Remove(examplePath)

	printLine("=== Example 1: Opening a presentation ===")
	prs, err := pptx.Open(examplePath)
	if err != nil {
		log.Fatalf("failed to open presentation: %v", err)
	}
	defer prs.Close()

	printFmtf("Opened presentation with %d slides\n", prs.SlideCount())
	printFmtf("Current title: %s\n", prs.Title())
	printFmtf("Current author: %s\n", prs.Creator())

	// Example 2: Reading metadata
	printLine("\n=== Example 2: Reading all metadata ===")
	readAllMetadata(prs)

	// Example 3: Modifying metadata
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

	// Save the changes
	if err := prs.Save(); err != nil {
		log.Fatalf("failed to save presentation: %v", err)
	}
	printLine("Changes saved!")

	// Example 4: Reopen and verify persistence
	printLine("\n=== Example 4: Reopening to verify persistence ===")
	prs2, err := pptx.Open(examplePath)
	if err != nil {
		log.Fatalf("failed to reopen presentation: %v", err)
	}
	defer prs2.Close()

	printFmtf("Title: %s\n", prs2.Title())
	printFmtf("Author: %s\n", prs2.Author())
	printFmtf("Keywords: %s\n", prs2.Keywords())
	printFmtf("Category: %s\n", prs2.Category())

	// Example 5: SaveAs to create a copy
	printLine("\n=== Example 5: SaveAs to create a copy ===")
	copyPath := "example_copy.pptx"
	defer os.Remove(copyPath)

	prs2.SetTitle("Copy of the Presentation")
	if err := prs2.SaveAs(copyPath); err != nil {
		log.Fatalf("failed to save copy: %v", err)
	}
	printFmtf("Saved copy to %s\n", copyPath)

	// Verify original wasn't modified
	prs3, err := pptx.Open(examplePath)
	if err != nil {
		log.Fatalf("failed to verify original presentation: %v", err)
	}
	defer prs3.Close()

	if prs3.Title() != "Updated Presentation Title" {
		printFmtf("Warning: Original presentation was modified\n")
	}
	printFmtf("Original title unchanged: %s\n", prs3.Title())

	// Example 6: Using CoreProperties directly
	printLine("\n=== Example 6: Using CoreProperties directly ===")
	prs4, err := pptx.Open(examplePath)
	if err != nil {
		log.Fatalf("failed to open presentation: %v", err)
	}
	defer prs4.Close()

	props := prs4.CoreProperties()
	printFmtf("CoreProperties.Title: %s\n", props.Title)
	printFmtf("CoreProperties.Creator: %s\n", props.Creator)
	printFmtf("CoreProperties.Revision: %s\n", props.Revision)

	// Modify using SetCoreProperties
	props.Title = "Final Title"
	props.Revision = "3"
	prs4.SetCoreProperties(props)

	if err := prs4.Save(); err != nil {
		log.Fatalf("failed to save: %v", err)
	}

	// Verify
	prs5, err := pptx.Open(examplePath)
	if err != nil {
		log.Fatalf("failed to verify: %v", err)
	}
	defer prs5.Close()

	printFmtf("After core props update - Title: %s, Revision: %s\n",
		prs5.Title(), prs5.Revision())

	// Example 7: Validate presentation
	printLine("\n=== Example 7: Validate presentation ===")
	issues := prs5.Validate()
	if len(issues) == 0 {
		printLine("Presentation is valid!")
	} else {
		printFmtf("Found %d validation issues:\n", len(issues))
		for _, issue := range issues {
			printFmtf("  - %s\n", issue.Description)
		}
	}

	printLine("\nAll examples completed successfully!")
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
