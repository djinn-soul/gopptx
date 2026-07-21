// examples/58-bridge-performance/main.go demonstrates batch command execution
// through the JSON bridge for performance-sensitive workflows.
//
// A single batch_execute command applies multiple operations in one round-trip,
// avoiding the overhead of individual ExecuteCommand calls. This is the same
// pattern used by the Python binding's Presentation.batch() context manager.
//
// Run with: go run ./examples/58-bridge-performance/main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "58_bridge_performance.pptx"

	// JSON bridge envelope keys.
	keyAPIVersion = "api_version"
	keyPayload    = "payload"
	keyRequestID  = "request_id"
	keySlideIndex = "slide_index"
	keyTitle      = "title"
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

	// --- Build a base presentation with 3 slides ---
	baseSlides := []pptx.SlideContent{
		pptx.NewSlide("Original Title 1").AddBullet("Slide 1 content"),
		pptx.NewSlide("Original Title 2").AddBullet("Slide 2 content"),
		pptx.NewSlide("Original Title 3").AddBullet("Slide 3 content"),
	}

	baseData, err := pptx.CreateWithSlides("Bridge Performance Demo", baseSlides)
	if err != nil {
		return fmt.Errorf("create base: %w", err)
	}

	tmpPath := filepath.Join(os.TempDir(), "gopptx_58_perf_base.pptx")
	if err = os.WriteFile(tmpPath, baseData, 0o600); err != nil {
		return fmt.Errorf("write tmp: %w", err)
	}
	defer os.Remove(tmpPath)

	// --- Open editor ---
	e, err := editor.OpenPresentationEditor(tmpPath)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	defer func() { _ = e.Close() }()

	// -----------------------------------------------------------------------
	// Approach A: Individual commands (3 separate round-trips)
	// -----------------------------------------------------------------------
	startIndividual := time.Now()
	runIndividualCommands(e)
	durationIndividual := time.Since(startIndividual)
	log.Printf("Individual commands duration: %v", durationIndividual)

	// -----------------------------------------------------------------------
	// Approach B: Batch command (1 round-trip for all 3 updates + a read)
	// -----------------------------------------------------------------------
	startBatch := time.Now()

	batchResp := executeJSON(e, map[string]any{
		keyAPIVersion: 1,
		"op":          editor.OpBatchExecute,
		keyPayload: map[string]any{
			"commands": []map[string]any{
				{
					"op":         editor.OpSetSlideTitle,
					keyPayload:   map[string]any{keySlideIndex: 0, keyTitle: "Slide 1 - Batch Updated"},
					keyRequestID: "batch-1",
				},
				{
					"op":         editor.OpSetSlideTitle,
					keyPayload:   map[string]any{keySlideIndex: 1, keyTitle: "Slide 2 - Batch Updated"},
					keyRequestID: "batch-2",
				},
				{
					"op":         editor.OpSetSlideTitle,
					keyPayload:   map[string]any{keySlideIndex: 2, keyTitle: "Slide 3 - Batch Updated"},
					keyRequestID: "batch-3",
				},
				{
					"op":         editor.OpSlideCount,
					keyPayload:   map[string]any{},
					keyRequestID: "batch-4",
				},
			},
			"stop_on_error": false,
		},
		keyRequestID: "batch-req-001",
	})

	durationBatch := time.Since(startBatch)
	log.Printf("Batch command duration:      %v", durationBatch)
	log.Printf("Batch response: %s", batchResp)

	// -----------------------------------------------------------------------
	// Verify final state: list all slide titles
	// -----------------------------------------------------------------------
	listResp := executeJSON(e, map[string]any{
		keyAPIVersion: 1,
		"op":          editor.OpListSlides,
		keyPayload:    map[string]any{},
		keyRequestID:  "verify-001",
	})
	log.Printf("Final slide list: %s", listResp)

	// --- Save the modified presentation ---
	outputPath := filepath.Join(outputDir, outputFile)
	if err = e.Save(outputPath); err != nil {
		return fmt.Errorf("save: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

// runIndividualCommands sends three set_slide_title commands one at a time.
func runIndividualCommands(e *editor.PresentationEditor) {
	for i, title := range []string{
		"Slide 1 - Individual Update",
		"Slide 2 - Individual Update",
		"Slide 3 - Individual Update",
	} {
		resp := executeJSON(e, map[string]any{
			keyAPIVersion: 1,
			"op":          editor.OpSetSlideTitle,
			keyPayload: map[string]any{
				keySlideIndex: i,
				keyTitle:      title,
			},
			keyRequestID: fmt.Sprintf("individual-%d", i+1),
		})
		log.Printf("individual set_slide_title [%d]: %s", i, resp)
	}
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
