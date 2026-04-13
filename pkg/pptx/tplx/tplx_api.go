// Package tplx provides a Jinja-style file-based template engine for PPTX files.
// Author a .pptx in PowerPoint with {{variable}} tokens in text boxes, table cells,
// or notes — then render it with a Go data map.
package tplx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
)

// Context is a map of template variable names to values.
// Values may be strings, bools, ints, float64, or []Row for loop data.
type Context map[string]any

// Row represents a single data record for table-row or slide loops.
type Row map[string]string

// Result holds the rendered PPTX bytes and exposes a Save helper.
type Result struct {
	data []byte
}

// Bytes returns the raw rendered PPTX bytes.
func (r *Result) Bytes() []byte { return r.data }

// Save writes the rendered PPTX to path.
func (r *Result) Save(path string) error {
	return os.WriteFile(path, r.data, 0o600)
}

// Options controls rendering behaviour.
type Options struct {
	Strict bool
}

// Render opens a PPTX template file and renders it with the given context.
func Render(path string, ctx Context) (*Result, error) {
	return RenderWithOptions(path, ctx, Options{})
}

// RenderWithOptions is like Render but accepts rendering options.
func RenderWithOptions(path string, ctx Context, opts Options) (*Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("tplx: read template %q: %w", path, err)
	}
	return RenderBytesWithOptions(data, ctx, opts)
}

// RenderBytes renders a PPTX template from in-memory bytes.
func RenderBytes(pptxBytes []byte, ctx Context) (*Result, error) {
	return RenderBytesWithOptions(pptxBytes, ctx, Options{})
}

// RenderBytesWithOptions is like RenderBytes but accepts rendering options.
func RenderBytesWithOptions(pptxBytes []byte, ctx Context, opts Options) (*Result, error) {
	zr, err := zip.NewReader(bytes.NewReader(pptxBytes), int64(len(pptxBytes)))
	if err != nil {
		return nil, fmt.Errorf("tplx: not a valid PPTX/ZIP: %w", err)
	}

	parts, err := readZipParts(zr)
	if err != nil {
		return nil, err
	}
	slideParts := collectSlideParts(parts)
	protectEscapedTemplateDelimiters(parts, slideParts)
	slideParts, parts = expandSlideLoops(slideParts, parts, ctx)
	if err = renderSlideParts(parts, slideParts, ctx, opts.Strict); err != nil {
		return nil, err
	}
	if err = renderNonSlideParts(parts, ctx, opts.Strict); err != nil {
		return nil, err
	}
	restoreEscapedTemplateDelimiters(parts, slideParts)

	out, err := repackZip(zr, parts)
	if err != nil {
		return nil, fmt.Errorf("tplx: repack: %w", err)
	}
	return &Result{data: out}, nil
}
