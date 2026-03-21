# API Overview

`gopptx` exposes three API surfaces.

## 1. Python API

- Entry point: `gopptx.Presentation`
- Location: `python/gopptx/presentation/*`
- Best for: application developers who want typed, ergonomic APIs

See: [Python Presentation API](reference/python-presentation-api.md)

## 2. Go API

- Best for: direct Go integrations, fluent deck construction, and core engine development.
- The Go API has two layers:
  - **High-level API**: `pkg/pptx`, `pkg/pptx/templates`, `pkg/pptx/charts` for fluent deck construction
  - **Engine API**: `pkg/pptx/editor` for direct engine integrations and core development

See:

- [Go API Reference](reference/go-api.md)
- [Go Slides Reference](reference/go-slides.md)
- [Go Notes, Comments, and Sections Reference](reference/go-notes-comments-sections.md)
- [Go Tables Reference](reference/go-tables.md)
- [Go Metadata, Protection, and Export Reference](reference/go-metadata-export.md)
- [Go Templates Reference](reference/go-templates.md)
- [Go Charts Reference](reference/go-charts.md)
- [Go Shapes and Connectors Reference](reference/go-shapes.md)

## 3. JSON Bridge API

- Operation identifiers: `add_slide`, `set_notes`, `update_chart_data`, etc.
- Source of truth:
  - `pkg/pptx/editor/opspec.go`
  - `python/gopptx/ops.py`
- Best for: C/Python bridge clients and non-Python/Go integrations

See: [JSON Bridge Operations](reference/bridge-operations.md)
