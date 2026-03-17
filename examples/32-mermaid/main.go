package main

import (
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/mermaid"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "32_mermaid.pptx"
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

	var slides []pptx.SlideContent

	// 1. Flowchart diagram.
	flowchart := `flowchart LR
    A[Start] --> B{Decision}
    B -- Yes --> C[Action A]
    B -- No --> D[Action B]
    C --> E[End]
    D --> E`

	flowSlide, err := diagramSlide("Mermaid Flowchart", flowchart)
	if err != nil {
		return fmt.Errorf("flowchart: %w", err)
	}
	slides = append(slides, flowSlide)

	// 2. Sequence diagram.
	sequence := `sequenceDiagram
    Alice->>Bob: Hello Bob, how are you?
    Bob-->>Alice: I am good thanks!
    Alice->>Bob: Great to hear`

	seqSlide, err := diagramSlide("Sequence Diagram", sequence)
	if err != nil {
		return fmt.Errorf("sequence: %w", err)
	}
	slides = append(slides, seqSlide)

	// 3. Pie chart diagram.
	pie := `pie title Browser Market Share
    "Chrome" : 65
    "Firefox" : 15
    "Safari" : 12
    "Edge" : 8`

	pieSlide, err := diagramSlide("Pie Chart Diagram", pie)
	if err != nil {
		return fmt.Errorf("pie: %w", err)
	}
	slides = append(slides, pieSlide)

	data, err := pptx.CreateWithSlides("Task 32: Mermaid Diagrams", slides)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	outputPath := outputDir + "/" + outputFile
	if err = os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

func diagramSlide(title, code string) (pptx.SlideContent, error) {
	diag, err := mermaid.CreateDiagram(code)
	if err != nil {
		return pptx.SlideContent{}, err
	}
	slide := pptx.NewSlide(title)
	for _, s := range diag.Shapes {
		slide = slide.AddShape(s)
	}
	for _, c := range diag.Connectors {
		slide = slide.AddConnector(c)
	}
	return slide, nil
}
