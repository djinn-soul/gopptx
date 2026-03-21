# API Overview

`gopptx` exposes three API surfaces.

## 1) Python Presentation API

- Entry point: `gopptx.Presentation`
- Location: `python/gopptx/presentation/*`
- Best for: application developers who want typed, ergonomic APIs

Use this page next:

- [Python Presentation API](reference/python-presentation-api.md)

## 1b) Go API Reference

- Entry points: `pkg/pptx`, `pkg/pptx/templates`, `pkg/pptx/charts`
- Best for: direct Go integrations and fluent deck construction

Use this page next:

- [Go API Reference](reference/go-api.md)
- [Go Slides Reference](reference/go-slides.md)
- [Go Notes, Comments, and Sections Reference](reference/go-notes-comments-sections.md)
- [Go Tables Reference](reference/go-tables.md)
- [Go Metadata, Protection, and Export Reference](reference/go-metadata-export.md)
- [Go Templates Reference](reference/go-templates.md)
- [Go Charts Reference](reference/go-charts.md)
- [Go Shapes and Connectors Reference](reference/go-shapes.md)

## 2) JSON Bridge Operations

- Operation identifiers: `add_slide`, `set_notes`, `update_chart_data`, etc.
- Source of truth:
  - `pkg/pptx/editor/opspec.go`
  - `python/gopptx/ops.py`
- Best for: C/Python bridge clients and non-Python integrations

Use this page next:

- [JSON Bridge Operations](reference/bridge-operations.md)

## 3) Go API

Go APIs live primarily under `pkg/pptx` and `pkg/pptx/editor`.

This surface is best for direct engine integrations and core development.
