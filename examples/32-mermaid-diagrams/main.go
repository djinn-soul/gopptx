// examples/32-mermaid-diagrams/main.go demonstrates Mermaid diagram rendering
// using the mermaid package to produce shape-based PPTX slides from diagram code.
//
// Three diagram types are shown: flowchart, sequence, and pie chart.
//
// Run with: go run ./examples/32-mermaid-diagrams/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/mermaid"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "32_mermaid_diagrams.pptx"
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

	// 1. Flowchart — decision and action nodes.
	flowchart := `flowchart LR
    A[Start] --> B{Validate Input}
    B -- Valid --> C[Process Data]
    B -- Invalid --> D[Return Error]
    C --> E{Save to DB}
    E -- Success --> F[Send Notification]
    E -- Failure --> G[Rollback]
    F --> H[End]
    G --> H`

	flowSlide, err := diagramSlide("Flowchart: Request Processing", flowchart)
	if err != nil {
		return fmt.Errorf("flowchart: %w", err)
	}
	slides = append(slides, flowSlide)

	// 2. Sequence diagram — actor message exchange.
	sequence := `sequenceDiagram
    participant C as Client
    participant A as API Gateway
    participant S as Service
    participant D as Database
    C->>A: POST /orders
    A->>S: Forward request
    S->>D: INSERT order
    D-->>S: Row ID
    S-->>A: 201 Created
    A-->>C: { "id": 42 }`

	seqSlide, err := diagramSlide("Sequence Diagram: Order Creation", sequence)
	if err != nil {
		return fmt.Errorf("sequence: %w", err)
	}
	slides = append(slides, seqSlide)

	// 3. Pie chart diagram — categorical distribution.
	pie := `pie title Programming Language Usage
    "Go" : 38
    "Python" : 29
    "TypeScript" : 18
    "Rust" : 10
    "Other" : 5`

	pieSlide, err := diagramSlide("Pie Chart: Language Distribution", pie)
	if err != nil {
		return fmt.Errorf("pie: %w", err)
	}
	slides = append(slides, pieSlide)

	// 4. Overview slide explaining the mermaid package.
	overview := pptx.NewSlide("Mermaid Diagram Support").
		AddBullet("mermaid.CreateDiagram(code) returns shapes and connectors").
		AddBullet("Shapes are added to a SlideContent via AddShape()").
		AddBullet("Connectors are added via AddConnector()").
		AddBullet("Supported: flowchart, sequenceDiagram, pie, classDiagram, and more")
	slides = append(slides, overview)

	data, err := pptx.CreateWithSlides("Task 32: Mermaid Diagrams", slides)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err = os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

// diagramSlide renders a Mermaid code string and returns a slide with the
// resulting shapes and connectors placed on it.
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
