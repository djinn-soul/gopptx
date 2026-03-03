// 60-presentation-api-metadata demonstrates the Presentation API
// for opening existing presentations and manipulating their metadata.
//
// This API is designed to be similar to python-pptx's Presentation(pptx_path)
// constructor.
package main

import (
	"fmt"
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

	fmt.Println("=== Example 1: Opening a presentation ===")
	prs, err := pptx.Open(examplePath)
	if err != nil {
		log.Fatalf("failed to open presentation: %v", err)
	}
	defer prs.Close()

	fmt.Printf("Opened presentation with %d slides\n", prs.SlideCount())
	fmt.Printf("Current title: %s\n", prs.Title())
	fmt.Printf("Current author: %s\n", prs.Creator())

	// Example 2: Reading metadata
	fmt.Println("\n=== Example 2: Reading all metadata ===")
	readAllMetadata(prs)

	// Example 3: Modifying metadata
	fmt.Println("\n=== Example 3: Modifying metadata ===")
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
	fmt.Println("Changes saved!")

	// Example 4: Reopen and verify persistence
	fmt.Println("\n=== Example 4: Reopening to verify persistence ===")
	prs2, err := pptx.Open(examplePath)
	if err != nil {
		log.Fatalf("failed to reopen presentation: %v", err)
	}
	defer prs2.Close()

	fmt.Printf("Title: %s\n", prs2.Title())
	fmt.Printf("Author: %s\n", prs2.Author())
	fmt.Printf("Keywords: %s\n", prs2.Keywords())
	fmt.Printf("Category: %s\n", prs2.Category())

	// Example 5: SaveAs to create a copy
	fmt.Println("\n=== Example 5: SaveAs to create a copy ===")
	copyPath := "example_copy.pptx"
	defer os.Remove(copyPath)

	prs2.SetTitle("Copy of the Presentation")
	if err := prs2.SaveAs(copyPath); err != nil {
		log.Fatalf("failed to save copy: %v", err)
	}
	fmt.Printf("Saved copy to %s\n", copyPath)

	// Verify original wasn't modified
	prs3, err := pptx.Open(examplePath)
	if err != nil {
		log.Fatalf("failed to verify original presentation: %v", err)
	}
	defer prs3.Close()

	if prs3.Title() != "Updated Presentation Title" {
		fmt.Printf("Warning: Original presentation was modified\n")
	}
	fmt.Printf("Original title unchanged: %s\n", prs3.Title())

	// Example 6: Using CoreProperties directly
	fmt.Println("\n=== Example 6: Using CoreProperties directly ===")
	prs4, err := pptx.Open(examplePath)
	if err != nil {
		log.Fatalf("failed to open presentation: %v", err)
	}
	defer prs4.Close()

	props := prs4.CoreProperties()
	fmt.Printf("CoreProperties.Title: %s\n", props.Title)
	fmt.Printf("CoreProperties.Creator: %s\n", props.Creator)
	fmt.Printf("CoreProperties.Revision: %s\n", props.Revision)

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

	fmt.Printf("After core props update - Title: %s, Revision: %s\n",
		prs5.Title(), prs5.Revision())

	// Example 7: Validate presentation
	fmt.Println("\n=== Example 7: Validate presentation ===")
	issues := prs5.Validate()
	if issues == nil || len(issues) == 0 {
		fmt.Println("Presentation is valid!")
	} else {
		fmt.Printf("Found %d validation issues:\n", len(issues))
		for _, issue := range issues {
			fmt.Printf("  - %s\n", issue.Description)
		}
	}

	fmt.Println("\nAll examples completed successfully!")
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
	fmt.Printf("Title: %s\n", prs.Title())
	fmt.Printf("Subject: %s\n", prs.Subject())
	fmt.Printf("Creator: %s\n", prs.Creator())
	fmt.Printf("Author: %s\n", prs.Author()) // Alias for Creator
	fmt.Printf("Keywords: %s\n", prs.Keywords())
	fmt.Printf("Description: %s\n", prs.Description())
	fmt.Printf("Last Modified By: %s\n", prs.LastModifiedBy())
	fmt.Printf("Revision: %s\n", prs.Revision())
	fmt.Printf("Created: %s\n", prs.Created())
	fmt.Printf("Modified: %s\n", prs.Modified())
	fmt.Printf("Category: %s\n", prs.Category())
	fmt.Printf("Content Status: %s\n", prs.ContentStatus())
}
