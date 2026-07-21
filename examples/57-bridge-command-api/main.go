// examples/57-bridge-command-api/main.go demonstrates the JSON bridge command API
// exposed by editor.ExecuteCommand. This is the same protocol used by language
// bindings (Python, C, Wasm) to drive the presentation editor over JSON.
//
// Run with: go run ./examples/57-bridge-command-api/main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "57_bridge_command_api.pptx"

	// JSON bridge envelope keys.
	keyAPIVersion = "api_version"
	keyPayload    = "payload"
	keyRequestID  = "request_id"
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

	// --- Build a base presentation and write it to a temp file ---
	baseSlides := []pptx.SlideContent{
		pptx.NewSlide("Slide One").
			AddBullet("First bullet on slide one").
			AddBullet("Second bullet on slide one"),
		pptx.NewSlide("Slide Two").
			AddBullet("Content on slide two"),
		pptx.NewSlide("Slide Three").
			AddBullet("Content on slide three"),
	}

	baseData, err := pptx.CreateWithSlides("Bridge Command API Demo", baseSlides)
	if err != nil {
		return fmt.Errorf("create base: %w", err)
	}

	tmpPath := filepath.Join(os.TempDir(), "gopptx_57_bridge_base.pptx")
	if err = os.WriteFile(tmpPath, baseData, 0o600); err != nil {
		return fmt.Errorf("write tmp: %w", err)
	}
	defer os.Remove(tmpPath)

	// --- Open the presentation editor ---
	e, err := editor.OpenPresentationEditor(tmpPath)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	defer func() { _ = e.Close() }()

	// --- Command 1: slide_count ---
	// Returns the number of slides in the presentation.
	resp1 := executeJSON(e, map[string]any{
		keyAPIVersion: 1,
		"op":          editor.OpSlideCount,
		keyPayload:    map[string]any{},
		keyRequestID:  "req-001",
	})
	log.Printf("slide_count response: %s", resp1)

	// --- Command 2: set_slide_title ---
	// Updates the title placeholder on slide index 0.
	resp2 := executeJSON(e, map[string]any{
		keyAPIVersion: 1,
		"op":          editor.OpSetSlideTitle,
		keyPayload: map[string]any{
			"slide_index": 0,
			"title":       "Updated via Bridge API",
		},
		keyRequestID: "req-002",
	})
	log.Printf("set_slide_title response: %s", resp2)

	// --- Command 3: list_slides ---
	// Returns a summary of all slides (index, title, layout).
	resp3 := executeJSON(e, map[string]any{
		keyAPIVersion: 1,
		"op":          editor.OpListSlides,
		keyPayload:    map[string]any{},
		keyRequestID:  "req-003",
	})
	log.Printf("list_slides response: %s", resp3)

	// --- Command 4: get_core_properties ---
	// Reads document metadata (author, title, description, etc.).
	resp4 := executeJSON(e, map[string]any{
		keyAPIVersion: 1,
		"op":          editor.OpGetCoreProperties,
		keyPayload:    map[string]any{},
		keyRequestID:  "req-004",
	})
	log.Printf("get_core_properties response: %s", resp4)

	// --- Save the modified presentation ---
	outputPath := filepath.Join(outputDir, outputFile)
	if err = e.Save(outputPath); err != nil {
		return fmt.Errorf("save: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

// executeJSON marshals a command map, sends it through editor.ExecuteCommand,
// and returns the raw JSON response string.
func executeJSON(e *editor.PresentationEditor, cmd map[string]any) string {
	b, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Sprintf(`{"ok":false,"error":{"code":"MARSHAL_ERROR","message":%q}}`, err.Error())
	}
	return editor.ExecuteCommand(e, string(b))
}
