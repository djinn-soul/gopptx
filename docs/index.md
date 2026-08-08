# gopptx

**A PowerPoint (`.pptx`) engine written in Go, with a first-class Python API and a stable JSON
command bridge.**

gopptx creates, reads, edits and exports OOXML presentations. The engine is pure Go — no
Microsoft Office, no LibreOffice and no headless browser is needed to produce a valid deck. The
Python package is a thin, typed layer over the same engine through a C shared library, so Python
code runs at Go speed without reimplementing any logic.

```python
from gopptx import Presentation

with Presentation.new("Quarterly Update") as pres:
    pres.add_bullet_slide("Highlights", ["Revenue +12%", "Retention +4%"])
    pres.save("deck.pptx")
```

## Pick your path

=== "Python application developer"

    1. [Installation](installation.md) — build the shared library, install the package
    2. [Quickstart](quickstart.md) — first deck
    3. [Core concepts](concepts.md) — **read the units section before writing any geometry**
    4. [Python library guide](guides/python-library.md) — the working reference

=== "Go developer"

    1. `go get github.com/djinn-soul/gopptx` — no native dependency
    2. [Quickstart](quickstart.md) — first deck
    3. [Go library guide](guides/go-library.md) — builders, editor, shapes, charts
    4. [Go API reference](reference/go-api.md)

=== "Integrating another language"

    1. [Core concepts](concepts.md) — the command envelope and handle model
    2. [Bridge operations](reference/bridge-operations.md) — all 179 operations
    3. [C bridge guidance](architecture/cbridge-guidance.md) — the ABI

## What you can build

- Automated reporting — QBRs, KPI decks, board updates generated from a data source
- Read-modify-save pipelines over decks a human authored
- Branded template systems with masters, layouts and placeholder overrides
- Markdown, HTML or URL content converted into slides
- Diagram-heavy decks — SmartArt and Mermaid rendered natively
- High-throughput generation, batching thousands of writes per boundary crossing

## Three things that trip people up

| | |
|---|---|
| **Geometry is EMU** | 914 400 EMU per inch. `bounds=(40, 120, 600, 220)` is a shape smaller than a pixel. Use `Inches()` / `Point()` / `Emu()`. See [Units](concepts.md#units-and-geometry). |
| **A new deck is not empty** | `Presentation.new(title)` already contains a title slide at index `0`. |
| **Reads are blocked inside `batch()`** | Buffered writes have not executed, so a read would see stale state. Use `execute_batch()` for interleaved reads and writes. See [Batch execution](guides/batch-execution.md). |

## Where the truth lives

Documentation drifts; source does not. When this site and the code disagree, the code wins:

| Question | Authoritative source |
|---|---|
| What operations does the bridge accept? | `pkg/pptx/editor/opspec.go` |
| What does a Python method take? | The `.pyi` stubs in `python/gopptx/` |
| What does the Go API expose? | `go doc github.com/djinn-soul/gopptx/pkg/pptx` |
| What is actually implemented vs. claimed? | [Feature matrix](reference/feature-matrix.md) |

The Python operation constants, chart-type enum, shape-type enum and slide builder are
**generated** from the Go declarations. `task check:generated` fails the build if they drift.
