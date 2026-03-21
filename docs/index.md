# gopptx Documentation

`gopptx` is a PowerPoint generation and editing engine with:

- a Go-native API for direct engine usage
- a Python API (`Presentation`) for application workflows
- a JSON bridge contract for runtime integrations

## How To Use This Website

Use the docs in this order for fastest onboarding:

1. [Installation](installation.md): setup Go or Python runtime.
2. [Quickstart](quickstart.md): generate your first deck in minutes.
3. [Examples](showcase/index.md): copy runnable examples with screenshots + downloadable `.pptx` outputs.
4. [Guides](guides/python-library.md): deeper patterns for Python/Go and batch execution.
5. [Reference](api-reference.md): full API surfaces and bridge operation names.

## What You Can Build

- Automated reports (QBR, KPI, board updates)
- Slide editing pipelines (read-modify-save)
- Presentation templates, masters, and placeholder overrides
- Tables, charts, shapes, media, notes, comments, and metadata
- High-throughput generation flows through batch operations

## Choose Your Path

- Python app developers: start with [Quickstart](quickstart.md), then [Python Library](guides/python-library.md).
- Go backend/tooling developers: start with [Quickstart](quickstart.md), then [Go Library](guides/go-library.md).
- Integration/runtime developers: start with [Core Concepts](concepts.md), then [JSON Bridge Operations](reference/bridge-operations.md).

## Source Of Truth

For operation IDs and bridge payload contracts, use:

- `pkg/pptx/editor/opspec.go`
- `python/gopptx/ops.py`
