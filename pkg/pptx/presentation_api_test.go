package pptx

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

// TestPresentation_OpenExisting verifies that a valid PPTX can be opened.
func TestPresentation_OpenExisting(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "test.pptx")

	// First create a simple presentation using the existing Create API
	data, err := Create("Test Presentation", 3)
	if err != nil {
		t.Fatalf("failed to create test presentation: %v", err)
	}
	if err := os.WriteFile(testFilePath, data, 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Now open it using the Presentation API
	prs, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to open presentation: %v", err)
	}
	defer prs.Close()

	// Verify basic properties
	if prs.SlideCount() != 3 {
		t.Errorf("expected 3 slides, got %d", prs.SlideCount())
	}
}

// TestPresentation_OpenNonexistent verifies an error is returned for missing files.
func TestPresentation_OpenNonexistent(t *testing.T) {
	prs, err := Open("nonexistent.pptx")
	if err == nil {
		prs.Close()
		t.Error("expected error for nonexistent file, got nil")
	}
}

// TestPresentation_MetadataGetters verifies all metadata getters work correctly.
// Note: This test uses the Open+Set+Save pattern because CreateWithMetadata
// doesn't serialize all CoreProperties fields to the initial file.
func TestPresentation_MetadataGetters(t *testing.T) {
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "test.pptx")

	// Create a simple presentation
	data, err := Create("Test Title", 1)
	if err != nil {
		t.Fatalf("failed to create test presentation: %v", err)
	}
	if err := os.WriteFile(testFilePath, data, 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Open and set all metadata
	prs, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to open presentation: %v", err)
	}

	prs.SetTitle("Test Title")
	prs.SetSubject("Test Subject")
	prs.SetCreator("Test Creator")
	prs.SetKeywords("test, keyword, metadata")
	prs.SetDescription("Test Description")
	prs.SetLastModifiedBy("Original Author")
	prs.SetRevision("1")
	prs.SetCategory("Test Category")
	prs.SetContentStatus("Draft")

	if err := prs.Save(); err != nil {
		prs.Close()
		t.Fatalf("failed to save presentation: %v", err)
	}
	prs.Close()

	// Reopen and verify metadata
	prs2, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to reopen presentation: %v", err)
	}
	defer prs2.Close()

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"Title", prs2.Title(), "Test Title"},
		{"Subject", prs2.Subject(), "Test Subject"},
		{"Creator", prs2.Creator(), "Test Creator"},
		{"Author", prs2.Author(), "Test Creator"},
		{"Keywords", prs2.Keywords(), "test, keyword, metadata"},
		{"Description", prs2.Description(), "Test Description"},
		{"LastModifiedBy", prs2.LastModifiedBy(), "Original Author"},
		{"Revision", prs2.Revision(), "1"},
		{"Category", prs2.Category(), "Test Category"},
		{"ContentStatus", prs2.ContentStatus(), "Draft"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.name, tt.expected, tt.got)
			}
		})
	}

	// Verify CoreProperties returns all at once
	props := prs2.CoreProperties()
	if props.Title != "Test Title" {
		t.Errorf("CoreProperties().Title: expected %q, got %q", "Test Title", props.Title)
	}
}

// TestPresentation_MetadataSetters verifies all metadata setters work correctly.
func TestPresentation_MetadataSetters(t *testing.T) {
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "test.pptx")

	// Create a simple presentation
	data, err := Create("Original Title", 1)
	if err != nil {
		t.Fatalf("failed to create test presentation: %v", err)
	}
	if err := os.WriteFile(testFilePath, data, 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Open and modify all metadata fields
	prs, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to open presentation: %v", err)
	}

	prs.SetTitle("Updated Title")
	prs.SetSubject("Updated Subject")
	prs.SetCreator("Updated Creator")
	prs.SetKeywords("updated, keywords")
	prs.SetDescription("Updated Description")
	prs.SetLastModifiedBy("Last Modifier")
	prs.SetRevision("2")
	prs.SetCategory("Updated Category")
	prs.SetContentStatus("Final")
	prs.SetCreated("2024-01-01T00:00:00Z")
	prs.SetModified("2024-12-31T23:59:59Z")

	// Verify changes before saving
	if prs.Title() != "Updated Title" {
		t.Errorf("SetTitle: expected %q, got %q", "Updated Title", prs.Title())
	}
	if prs.Creator() != "Updated Creator" {
		t.Errorf("SetCreator: expected %q, got %q", "Updated Creator", prs.Creator())
	}

	// Verify SetAuthor is an alias for SetCreator
	prs.SetAuthor("Updated Author")
	if prs.Creator() != "Updated Author" {
		t.Errorf("SetAuthor: expected %q, got %q", "Updated Author", prs.Creator())
	}
	if prs.Author() != "Updated Author" {
		t.Errorf("Author getter: expected %q, got %q", "Updated Author", prs.Author())
	}

	// Save and reopen to verify persistence
	if err := prs.Save(); err != nil {
		prs.Close()
		t.Fatalf("failed to save presentation: %v", err)
	}
	prs.Close()

	// Reopen and verify all values persist
	prs2, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to reopen presentation: %v", err)
	}
	defer prs2.Close()

	expectedValues := map[string]string{
		"Title":          "Updated Title",
		"Subject":        "Updated Subject",
		"Creator":        "Updated Author",
		"Author":         "Updated Author",
		"Keywords":       "updated, keywords",
		"Description":    "Updated Description",
		"LastModifiedBy": "Last Modifier",
		"Revision":       "2",
		"Category":       "Updated Category",
		"ContentStatus":  "Final",
		"Created":        "2024-01-01T00:00:00Z",
		"Modified":       "2024-12-31T23:59:59Z",
	}

	for field, expected := range expectedValues {
		var got string
		switch field {
		case "Title":
			got = prs2.Title()
		case "Subject":
			got = prs2.Subject()
		case "Creator":
			got = prs2.Creator()
		case "Author":
			got = prs2.Author()
		case "Keywords":
			got = prs2.Keywords()
		case "Description":
			got = prs2.Description()
		case "LastModifiedBy":
			got = prs2.LastModifiedBy()
		case "Revision":
			got = prs2.Revision()
		case "Category":
			got = prs2.Category()
		case "ContentStatus":
			got = prs2.ContentStatus()
		case "Created":
			got = prs2.Created()
		case "Modified":
			got = prs2.Modified()
		}

		if got != expected {
			t.Errorf("%s: expected %q, got %q", field, expected, got)
		}
	}
}

// TestPresentation_SetCoreProperties verifies setting all core properties at once.
func TestPresentation_SetCoreProperties(t *testing.T) {
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "test.pptx")

	data, err := Create("Original", 1)
	if err != nil {
		t.Fatalf("failed to create test presentation: %v", err)
	}
	if err := os.WriteFile(testFilePath, data, 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	prs, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to open presentation: %v", err)
	}

	newProps := common.CoreProperties{
		Title:          "New Title",
		Subject:        "New Subject",
		Creator:        "New Creator",
		Keywords:       "new, keywords",
		Description:    "New Description",
		LastModifiedBy: "New Modifier",
		Revision:       "3",
		Created:        "2024-02-01T00:00:00Z",
		Modified:       "2024-02-15T00:00:00Z",
		Category:       "New Category",
		ContentStatus:  "Published",
	}
	prs.SetCoreProperties(newProps)

	// Verify before saving
	if prs.Title() != "New Title" {
		t.Errorf("SetCoreProperties: expected Title %q, got %q", "New Title", prs.Title())
	}
	if prs.Creator() != "New Creator" {
		t.Errorf("SetCoreProperties: expected Creator %q, got %q", "New Creator", prs.Creator())
	}

	// Save and reopen to verify
	if err := prs.Save(); err != nil {
		prs.Close()
		t.Fatalf("failed to save presentation: %v", err)
	}
	prs.Close()

	prs2, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to reopen presentation: %v", err)
	}
	defer prs2.Close()

	props := prs2.CoreProperties()
	if props.Title != "New Title" {
		t.Errorf("persistent Title: expected %q, got %q", "New Title", props.Title)
	}
	if props.Creator != "New Creator" {
		t.Errorf("persistent Creator: expected %q, got %q", "New Creator", props.Creator)
	}
	if props.Category != "New Category" {
		t.Errorf("persistent Category: expected %q, got %q", "New Category", props.Category)
	}
}

// TestPresentation_SaveModifiesFile verifies that Save modifies the file in place.
func TestPresentation_SaveModifiesFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "test.pptx")

	// Create initial file
	data, err := Create("Initial Title", 1)
	if err != nil {
		t.Fatalf("failed to create test presentation: %v", err)
	}
	if err := os.WriteFile(testFilePath, data, 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Get original file size and content
	originalInfo, err := os.Stat(testFilePath)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}
	originalSize := originalInfo.Size()

	// Open, modify, and save
	prs, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to open presentation: %v", err)
	}

	prs.SetTitle("Modified Title")
	prs.SetCreator("New Author")

	if err := prs.Save(); err != nil {
		prs.Close()
		t.Fatalf("failed to save presentation: %v", err)
	}
	prs.Close()

	// Verify file was modified (should be different size due to different content)
	newInfo, err := os.Stat(testFilePath)
	if err != nil {
		t.Fatalf("failed to stat modified file: %v", err)
	}
	newSize := newInfo.Size()

	// File size may be different due to content changes
	// But let's verify the content actually changed
	prs2, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to reopen presentation: %v", err)
	}
	defer prs2.Close()

	if prs2.Title() != "Modified Title" {
		t.Errorf("failed to persist Title: expected %q, got %q", "Modified Title", prs2.Title())
	}
	if prs2.Creator() != "New Author" {
		t.Errorf("failed to persist Creator: expected %q, got %q", "New Author", prs2.Creator())
	}

	_ = originalSize
	_ = newSize
}

// TestPresentation_SaveAsCreatesNewFile verifies SaveAs creates a new file.
func TestPresentation_SaveAsCreatesNewFile(t *testing.T) {
	tmpDir := t.TempDir()
	originalPath := filepath.Join(tmpDir, "original.pptx")
	newPath := filepath.Join(tmpDir, "copy.pptx")

	// Create original file
	data, err := Create("Original Title", 1)
	if err != nil {
		t.Fatalf("failed to create test presentation: %v", err)
	}
	if err := os.WriteFile(originalPath, data, 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Get original file info
	originalInfo, err := os.Stat(originalPath)
	if err != nil {
		t.Fatalf("failed to stat original file: %v", err)
	}
	originalModTime := originalInfo.ModTime()

	// Open and save to new path
	prs, err := Open(originalPath)
	if err != nil {
		t.Fatalf("failed to open presentation: %v", err)
	}

	prs.SetTitle("Copied Presentation")

	if err := prs.SaveAs(newPath); err != nil {
		prs.Close()
		t.Fatalf("failed to save as: %v", err)
	}
	prs.Close()

	// Verify new file exists
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Fatal("SaveAs did not create new file")
	}

	// Verify original file still has original title (was not modified)
	prs2, err := Open(originalPath)
	if err != nil {
		t.Fatalf("failed to reopen original: %v", err)
	}
	defer prs2.Close()

	if prs2.Title() != "Original Title" {
		t.Errorf("original file was modified: expected Title %q, got %q", "Original Title", prs2.Title())
	}

	// Verify new file has new title
	prs3, err := Open(newPath)
	if err != nil {
		t.Fatalf("failed to open new file: %v", err)
	}
	defer prs3.Close()

	if prs3.Title() != "Copied Presentation" {
		t.Errorf("new file Title incorrect: expected %q, got %q", "Copied Presentation", prs3.Title())
	}

	// Verify new file is a separate file (different mtime)
	newInfo, err := os.Stat(newPath)
	if err != nil {
		t.Fatalf("failed to stat new file: %v", err)
	}

	// Check that original file wasn't modified (mtimes should be the same)
	afterOriginalInfo, err := os.Stat(originalPath)
	if err != nil {
		t.Fatalf("failed to stat original file after SaveAs: %v", err)
	}

	if afterOriginalInfo.ModTime() != originalModTime {
		t.Error("original file was modified by SaveAs, it should be unchanged")
	}

	_ = newInfo // Used to verify file exists
	_ = originalInfo
	_ = afterOriginalInfo
}

// TestPresentation_CloseReleasesResources verifies Close works correctly.
func TestPresentation_CloseReleasesResources(t *testing.T) {
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "test.pptx")

	data, err := Create("Test", 1)
	if err != nil {
		t.Fatalf("failed to create test presentation: %v", err)
	}
	if err := os.WriteFile(testFilePath, data, 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	prs, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to open presentation: %v", err)
	}

	// Verify we can access data before close
	if prs.SlideCount() != 1 {
		t.Errorf("expected 1 slide, got %d", prs.SlideCount())
	}

	// Close the presentation
	if err := prs.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Calling Close again should be safe (idempotent)
	if err := prs.Close(); err != nil {
		t.Errorf("second Close failed: %v", err)
	}

	// Accessing methods after close should be safe (may return zero values)
	// This depends on the implementation - we just verify it doesn't panic
	_ = prs.SlideCount()
	_ = prs.Title()
}

// TestPresentation_CloseOnNil verifies Close works on nil presentation.
func TestPresentation_CloseOnNil(t *testing.T) {
	var prs *Presentation
	if err := prs.Close(); err != nil {
		t.Errorf("Close on nil should not error: %v", err)
	}
}

// TestPresentation_SlideCount returns correct slide count.
func TestPresentation_SlideCount(t *testing.T) {
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "test.pptx")

	// Test with different slide counts
	for _, expectedCount := range []int{1, 2, 5, 10} {
		t.Run("", func(t *testing.T) {
			data, err := Create("Multi-Slide", expectedCount)
			if err != nil {
				t.Fatalf("failed to create test presentation: %v", err)
			}
			if err := os.WriteFile(testFilePath, data, 0o600); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			prs, err := Open(testFilePath)
			if err != nil {
				t.Fatalf("failed to open presentation: %v", err)
			}
			defer prs.Close()

			if count := prs.SlideCount(); count != expectedCount {
				t.Errorf("expected %d slides, got %d", expectedCount, count)
			}
		})
	}
}

// TestPresentation_Validate validates presentation structure.
func TestPresentation_Validate(t *testing.T) {
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "test.pptx")

	data, err := Create("Valid Presentation", 2)
	if err != nil {
		t.Fatalf("failed to create test presentation: %v", err)
	}
	if err := os.WriteFile(testFilePath, data, 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	prs, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to open presentation: %v", err)
	}
	defer prs.Close()

	issues := prs.Validate()
	// A valid presentation should have no critical issues
	// nil means no issues (editor returned no validation problems)
	// empty slice also means no issues
	if len(issues) > 0 {
		t.Errorf("expected no validation issues, got %d", len(issues))
	}

	// Test nil presentation returns nil
	var nilPrs *Presentation
	nilIssues := nilPrs.Validate()
	if nilIssues != nil {
		t.Errorf("expected nil validation result for nil presentation")
	}
}

// TestPresentation_MetadataPersistenceRoundTrip verifies metadata persists correctly
// through a full create-open-modify-save-open cycle.
func TestPresentation_MetadataPersistenceRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "roundtrip.pptx")

	// Create with initial metadata
	meta := Metadata{}
	meta.Title = "Round Trip Test"
	meta.Subject = "Testing persistence"
	meta.Creator = "Initial Author"
	meta.Description = "Initial description"
	meta.CoreProperties = common.CoreProperties{
		Title:         "Round Trip Test",
		Subject:       "Testing persistence",
		Creator:       "Initial Author",
		Keywords:      "initial, test",
		Description:   "Initial description",
		Revision:      "1",
		Category:      "Testing",
		ContentStatus: "Draft",
	}

	data, err := CreateWithMetadata(meta, []SlideContent{
		NewSlide("Slide 1"),
		NewSlide("Slide 2"),
	})
	if err != nil {
		t.Fatalf("failed to create test presentation: %v", err)
	}
	if err := os.WriteFile(testFilePath, data, 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Open, modify, and save
	prs1, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("first open failed: %v", err)
	}

	// Modify all fields
	prs1.SetTitle("Updated Title")
	prs1.SetSubject("Updated Subject")
	prs1.SetCreator("Updated Creator")
	prs1.SetKeywords("updated, keywords, roundtrip")
	prs1.SetDescription("Updated description with more details")
	prs1.SetLastModifiedBy("Modified By User")
	prs1.SetRevision("2")
	prs1.SetCategory("Gopptx Testing")
	prs1.SetContentStatus("Final")

	if err := prs1.Save(); err != nil {
		prs1.Close()
		t.Fatalf("save failed: %v", err)
	}
	prs1.Close()

	// Reopen and verify all persisted values
	prs2, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("second open failed: %v", err)
	}
	defer prs2.Close()

	expected := map[string]string{
		"Title":          "Updated Title",
		"Subject":        "Updated Subject",
		"Creator":        "Updated Creator",
		"Keywords":       "updated, keywords, roundtrip",
		"Description":    "Updated description with more details",
		"LastModifiedBy": "Modified By User",
		"Revision":       "2",
		"Category":       "Gopptx Testing",
		"ContentStatus":  "Final",
	}

	for field, expectedValue := range expected {
		var got string
		switch field {
		case "Title":
			got = prs2.Title()
		case "Subject":
			got = prs2.Subject()
		case "Creator":
			got = prs2.Creator()
		case "Keywords":
			got = prs2.Keywords()
		case "Description":
			got = prs2.Description()
		case "LastModifiedBy":
			got = prs2.LastModifiedBy()
		case "Revision":
			got = prs2.Revision()
		case "Category":
			got = prs2.Category()
		case "ContentStatus":
			got = prs2.ContentStatus()
		}

		if got != expectedValue {
			t.Errorf("%s: expected %q, got %q", field, expectedValue, got)
		}
	}
}

// TestPresentation_PreservesZipStructure verifies that Save preserves
// the ZIP structure of the PPTX file.
func TestPresentation_PreservesZipStructure(t *testing.T) {
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "structure.pptx")

	// Create a multi-slide presentation
	data, err := Create("Structure Test", 3)
	if err != nil {
		t.Fatalf("failed to create test presentation: %v", err)
	}
	if err := os.WriteFile(testFilePath, data, 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Read original ZIP structure
	originalFile, err := os.Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to open original file: %v", err)
	}
	defer originalFile.Close()

	originalInfo, err := originalFile.Stat()
	if err != nil {
		t.Fatalf("failed to stat original file: %v", err)
	}
	originalZr, err := zip.NewReader(originalFile, originalInfo.Size())
	if err != nil {
		t.Fatalf("failed to read original as ZIP: %v", err)
	}

	var originalFiles []string
	for _, f := range originalZr.File {
		originalFiles = append(originalFiles, f.Name)
	}

	// Open, modify, and save
	prs, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to open presentation: %v", err)
	}

	prs.SetTitle("Structure Modified")
	prs.SetCategory("Structure Test")

	if err := prs.Save(); err != nil {
		prs.Close()
		t.Fatalf("failed to save presentation: %v", err)
	}
	prs.Close()

	// Read saved ZIP structure
	savedFile, err := os.Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to open saved file: %v", err)
	}
	defer savedFile.Close()

	savedInfo, err := savedFile.Stat()
	if err != nil {
		t.Fatalf("failed to stat saved file: %v", err)
	}
	savedZr, err := zip.NewReader(savedFile, savedInfo.Size())
	if err != nil {
		t.Fatalf("failed to read saved as ZIP: %v", err)
	}

	var savedFiles []string
	for _, f := range savedZr.File {
		savedFiles = append(savedFiles, f.Name)
	}

	// Verify core properties file exists
	corePropsPath := "docProps/core.xml"
	if !testutil.ZipHasFile(savedZr, corePropsPath) {
		t.Errorf("saved file missing %s", corePropsPath)
	}

	// Verify presentation file exists
	presentationPath := "ppt/presentation.xml"
	if !testutil.ZipHasFile(savedZr, presentationPath) {
		t.Errorf("saved file missing %s", presentationPath)
	}

	// Verify slides exist (at least 3 slide files)
	slideCount := 0
	for _, f := range savedFiles {
		if strings.HasPrefix(f, "ppt/slides/") && strings.HasSuffix(f, ".xml") {
			slideCount++
		}
	}
	if slideCount < 3 {
		t.Errorf("expected at least 3 slide files, got %d", slideCount)
	}

	// Verify saved metadata is correct
	corePropsContent := testutil.ReadZipFile(t, savedZr, corePropsPath)
	if !strings.Contains(corePropsContent, "Structure Modified") {
		t.Error("modified title not found in saved core properties")
	}

	_ = originalFiles // compare if needed
	_ = savedFiles
}

// TestPresentation_NilHandling verifies methods handle nil presentation gracefully.
func TestPresentation_NilHandling(t *testing.T) {
	var prs *Presentation

	// All getters should return zero values safely
	if prs.SlideCount() != 0 {
		t.Error("SlideCount() on nil should return 0")
	}

	if prs.Title() != "" {
		t.Error("Title() on nil should return empty string")
	}

	if prs.Creator() != "" {
		t.Error("Creator() on nil should return empty string")
	}

	if prs.Author() != "" {
		t.Error("Author() on nil should return empty string")
	}

	if prs.Subject() != "" {
		t.Error("Subject() on nil should return empty string")
	}

	if prs.Keywords() != "" {
		t.Error("Keywords() on nil should return empty string")
	}

	if prs.Description() != "" {
		t.Error("Description() on nil should return empty string")
	}

	if prs.LastModifiedBy() != "" {
		t.Error("LastModifiedBy() on nil should return empty string")
	}

	if prs.Revision() != "" {
		t.Error("Revision() on nil should return empty string")
	}

	if prs.Created() != "" {
		t.Error("Created() on nil should return empty string")
	}

	if prs.Modified() != "" {
		t.Error("Modified() on nil should return empty string")
	}

	if prs.Category() != "" {
		t.Error("Category() on nil should return empty string")
	}

	if prs.ContentStatus() != "" {
		t.Error("ContentStatus() on nil should return empty string")
	}

	if prs.CoreProperties().Title != "" {
		t.Error("CoreProperties() on nil should return empty CoreProperties")
	}

	if prs.Validate() != nil {
		t.Error("Validate() on nil should return nil")
	}

	// Setters should be safe to call on nil (no-op)
	prs.SetTitle("test")
	prs.SetSubject("test")
	prs.SetCreator("test")
	prs.SetAuthor("test")
	prs.SetKeywords("test")
	prs.SetDescription("test")
	prs.SetLastModifiedBy("test")
	prs.SetRevision("test")
	prs.SetCreated("2024-01-01")
	prs.SetModified("2024-01-01")
	prs.SetCategory("test")
	prs.SetContentStatus("test")
	prs.SetCoreProperties(common.CoreProperties{Title: "test"})
}

// TestPresentation_ModifiedTimestamp verifies modification timestamp handling.
func TestPresentation_ModifiedTimestamp(t *testing.T) {
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "timestamp.pptx")

	data, err := Create("Timestamp Test", 1)
	if err != nil {
		t.Fatalf("failed to create test presentation: %v", err)
	}
	if err := os.WriteFile(testFilePath, data, 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	prs, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to open presentation: %v", err)
	}

	prs.SetModified("2024-06-15T10:30:45Z")

	if err := prs.Save(); err != nil {
		prs.Close()
		t.Fatalf("failed to save presentation: %v", err)
	}
	prs.Close()

	// Reopen and verify
	prs2, err := Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to reopen presentation: %v", err)
	}
	defer prs2.Close()

	modified := prs2.Modified()
	if modified != "2024-06-15T10:30:45Z" {
		t.Errorf("Modified timestamp: expected %q, got %q", "2024-06-15T10:30:45Z", modified)
	}
}
