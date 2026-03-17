package main

import (
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "20_validation_repair.pptx"
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

	// 1. Build a presentation to validate and repair.
	slides := []pptx.SlideContent{
		pptx.NewSlide("Validation & Repair").
			AddBullet("Create a presentation").
			AddBullet("Validate structural integrity").
			AddBullet("Repair any issues automatically"),
		pptx.NewSlide("Second Slide").
			AddBullet("Multiple slides supported").
			AddBullet("Round-trip safe"),
	}

	data, err := pptx.CreateWithSlides("Task 20: Validation & Repair", slides)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	// 2. Validate the generated file.
	issues, err := pptx.Validate(data)
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	log.Printf("Validation found %d issue(s)\n", len(issues))
	for _, issue := range issues {
		log.Printf("  - [%s] %s: %s\n", issue.Severity, issue.Code, issue.Description)
	}

	// 3. Repair any issues.
	repaired, result, err := pptx.Repair(data)
	if err != nil {
		return fmt.Errorf("repair: %w", err)
	}
	log.Printf("Repair applied %d fix(es), %d unrepaired\n", len(result.IssuesRepaired), len(result.IssuesUnrepaired))

	// 4. Write the repaired file.
	outputPath := outputDir + "/" + outputFile
	if err = os.WriteFile(outputPath, repaired, 0o600); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
