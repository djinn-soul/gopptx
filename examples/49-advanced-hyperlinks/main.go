package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/action"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Smoke test failed: %v", err)
	}
	fmt.Println("Advanced Hyperlinks Smoke Test: SUCCESS")
	fmt.Println("To verify visually: Open examples/output/49_advanced_hyperlinks.pptx and try clicking the shapes.")
}

func run() error {
	const outputDir = "examples/output"
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	outFile := filepath.Join(outputDir, "49_advanced_hyperlinks.pptx")
	deck := pptx.NewPresentationBuilder("Advanced Hyperlinks")

	// Slide 1: Local File Link
	slide1 := pptx.NewSlide("Local File Link")
	rect1 := pptx.NewRectangle(1, 2, 4, 1).
		WithText("Open 'test.txt'").
		WithClickAction(action.NewHyperlink(action.HyperlinkFile("C:\\Temp\\test.txt")))
	slide1 = slide1.AddShape(rect1)
	deck.AddSlide(slide1)

	// Slide 2: Program Link
	slide2 := pptx.NewSlide("Program Link")
	rect2 := pptx.NewRectangle(1, 2, 4, 1).
		WithText("Launch Calculator").
		WithClickAction(action.NewHyperlink(action.HyperlinkProgram("C:\\Windows\\System32\\calc.exe")))
	slide2 = slide2.AddShape(rect2)
	deck.AddSlide(slide2)

	if err := deck.WriteToFile(outFile); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	// Verify the relationships in the output file
	return verifyRelationships(outFile)
}

func verifyRelationships(pptxPath string) (retErr error) {
	r, err := zip.OpenReader(pptxPath)
	if err != nil {
		return fmt.Errorf("failed to open pptx for verification: %w", err)
	}
	defer func() {
		if closeErr := r.Close(); closeErr != nil && retErr == nil {
			retErr = fmt.Errorf("close pptx reader: %w", closeErr)
		}
	}()

	// Map of slide index to expected partial target
	expected := map[int]string{
		0: "file:///C:/Temp/test.txt",
		1: "file:///C:/Windows/System32/calc.exe",
	}

	for i, target := range expected {
		relPath := fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", i+1)
		content, err := readZipFile(r, relPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", relPath, err)
		}

		// Check for TargetMode="External"
		if !strings.Contains(content, `TargetMode="External"`) {
			return fmt.Errorf("slide %d rels missing TargetMode=\"External\"", i+1)
		}
		// Check for Target
		if !strings.Contains(content, target) {
			return fmt.Errorf("slide %d rels missing target %q\nContent: %s", i+1, target, content)
		}
	}
	return nil
}

func readZipFile(r *zip.ReadCloser, name string) (string, error) {
	for _, f := range r.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			b, readErr := io.ReadAll(rc)
			closeErr := rc.Close()
			if readErr != nil {
				return "", readErr
			}
			if closeErr != nil {
				return "", closeErr
			}
			return string(b), nil
		}
	}
	return "", fmt.Errorf("file %s not found in zip", name)
}
