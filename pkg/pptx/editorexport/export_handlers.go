// Package editorexport registers export operation handlers (PDF, HTML) into the
// editor bridge. It lives in a separate package to break the import cycle:
//
//	editor → export → pptx → editor
//
// Call Register() once at startup (e.g. from bindings/c/bridge.go init).
package editorexport

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
)

// Register wires the PDF and HTML export handlers into the editor command
// dispatcher. It must be called before any bridge requests are processed.
func Register() {
	editor.RegisterExportHandler(editor.OpExportPDF, handleExportPDF)
	editor.RegisterExportHandler(editor.OpExportHTML, handleExportHTML)
}

func handleExportPDF(e *editor.PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := editor.ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := editor.NewPayloadValidator()
	outputPath, _ := v.RequireString(p, "output_path")
	if v.HasErrors() {
		return nil, v.Error()
	}

	driver := export.PDFDriverAuto
	if d, ok := p["driver"].(string); ok && d != "" {
		parsed, parseErr := export.ParsePDFDriver(d)
		if parseErr != nil {
			return nil, parseErr
		}
		driver = parsed
	}

	tmp, err := os.CreateTemp("", "gopptx-export-*.pptx")
	if err != nil {
		return nil, err
	}
	tmpPath := tmp.Name()
	if closeErr := tmp.Close(); closeErr != nil {
		_ = os.Remove(tmpPath)
		return nil, closeErr
	}
	defer os.Remove(tmpPath)

	if err = e.Save(tmpPath); err != nil {
		return nil, err
	}

	outPath, err := filepath.Abs(outputPath)
	if err != nil {
		return nil, err
	}

	opts := export.PDFOptions{Driver: driver}
	if err = export.PDFFromFileWithOptions(tmpPath, outPath, opts); err != nil {
		return nil, err
	}

	return map[string]string{"output_path": outPath}, nil
}

func handleExportHTML(e *editor.PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := editor.ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	outputPath, _ := p["output_path"].(string)

	opts := export.DefaultHTMLOptions()
	if v, ok := p["embed_images"].(bool); ok {
		opts.EmbedImages = v
	}
	if v, ok := p["include_navigation"].(bool); ok {
		opts.IncludeNavigation = v
	}
	if v, ok := p["include_slide_numbers"].(bool); ok {
		opts.IncludeSlideNumbers = v
	}
	if v, ok := p["base_url"].(string); ok {
		opts.BaseURL = v
	}

	tmp, err := os.CreateTemp("", "gopptx-export-*.pptx")
	if err != nil {
		return nil, err
	}
	tmpPath := tmp.Name()
	if closeErr := tmp.Close(); closeErr != nil {
		_ = os.Remove(tmpPath)
		return nil, closeErr
	}
	defer os.Remove(tmpPath)

	if err = e.Save(tmpPath); err != nil {
		return nil, err
	}

	title, slides, err := export.SlidesFromPPTX(tmpPath)
	if err != nil {
		return nil, err
	}

	htmlStr := export.HTMLWithOptions(title, slides, opts)

	if outputPath != "" {
		outPath, absErr := filepath.Abs(outputPath)
		if absErr != nil {
			return nil, absErr
		}
		if writeErr := os.WriteFile(outPath, []byte(htmlStr), 0o600); writeErr != nil {
			return nil, writeErr
		}
		return map[string]string{"output_path": outPath}, nil
	}

	return map[string]string{"html": htmlStr}, nil
}
