# gopptx

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go 1.25+](https://img.shields.io/badge/go-1.25%2B-00ADD8?logo=go)](https://go.dev)
[![Python 3.10+](https://img.shields.io/badge/python-3.10%2B-3776AB?logo=python&logoColor=white)](https://www.python.org)

**A PowerPoint (`.pptx`) engine written in Go, with a first-class Python API and a stable JSON command bridge.**

gopptx creates, reads, edits and exports OOXML presentations. The engine is pure Go — no
Microsoft Office, no LibreOffice and no headless browser is required to produce a valid deck.
The Python package is a thin, typed layer over the same engine through a C shared library, so
Python code runs at Go speed without reimplementing any logic.

```python
from gopptx import Presentation

with Presentation.new("Quarterly Update") as pres:
    pres.add_bullet_slide("Highlights", ["Revenue +12%", "Retention +4%"])
    pres.save("deck.pptx")
```

---

## Table of contents

- [Why gopptx](#why-gopptx)
- [What it can do](#what-it-can-do)
- [Installation](#installation)
- [Quickstart — Python](#quickstart--python)
- [Quickstart — Go](#quickstart--go)
- [Editing an existing deck](#editing-an-existing-deck)
- [Exporting to PDF and HTML](#exporting-to-pdf-and-html)
- [Throughput and batching](#throughput-and-batching)
- [Architecture](#architecture)
- [Documentation map](#documentation-map)
- [Development](#development)
- [License](#license)

---

## Why gopptx

| | |
|---|---|
| **One engine, two languages** | The Go package and the Python package are the same code. Behaviour, validation and output bytes are identical. |
| **No Office dependency for authoring** | Decks are generated from OOXML primitives directly. Office is only needed if you choose the PowerPoint PDF driver. |
| **Reads *and* writes** | Most libraries do one or the other. gopptx opens an existing deck, walks its shape tree, edits in place and writes it back. |
| **Batching for throughput** | Write-heavy Python workloads cross the FFI boundary once per batch, not once per operation. |
| **Validated output** | `validate()` checks OPC package rules, part relationships and content types before you ship a file. `repair()` fixes what it can. |

## What it can do

| Area | Capability |
|---|---|
| **Slides** | Add, insert, remove, move, duplicate, hide, copy between decks, merge decks, rebind layouts |
| **Text** | Runs with bold/italic/underline/strike/caps/sub/superscript, fonts, sizes, colours, highlight, paragraph alignment/indent/spacing, autofit, anchors, RTL and complex-script handling |
| **Shapes** | ~200 preset geometries, freeforms, groups, connectors (straight/elbow/curved, auto-routing), solid/gradient/pattern/picture fills, lines, shadows, glow, soft edges, rotation, flips, z-order |
| **Tables** | Row/column insert and remove, cell merge and split, widths and heights, borders, banding flags, built-in and custom table styles, load from rows or dicts |
| **Charts** | 24 chart types, combo charts, multiple charts per slide, an embedded Excel workbook so the data stays editable in PowerPoint, axes, legends, data labels, trendlines, error bars, per-point formatting |
| **Images & media** | Local files, bytes, base64, URLs, natural-size placement, crops, effects, video, audio, online video, OLE objects, media extraction |
| **Diagrams** | SmartArt (every layout PowerPoint offers, node-level edits, quick styles, colour styles) and native Mermaid rendering |
| **Layouts & themes** | Slide masters, layouts, placeholders and overrides, notes master, handout master, colour schemes, font schemes, built-in theme presets, grid/stack/distribute layout helpers |
| **Deck metadata** | Core and app document properties, sections, comments and authors, speaker notes, headers and footers, custom XML parts, VBA projects, digital-signature detection, mark-as-final, modify password |
| **Motion** | Slide transitions (including morph), entrance/emphasis/exit animations, slide-show settings and custom shows |
| **Import** | Markdown → slides, HTML → slides, URL fetch → slides, Jinja2-style templating |
| **Export** | PDF (native Go renderer, LibreOffice or PowerPoint COM), HTML, flat XML, grayscale conversion, package compression and size analysis |

Every area above has a runnable example under [`examples/`](examples/) — 96 of them.

---

## Installation

### Prerequisites

| | Requirement |
|---|---|
| Go | 1.25 or later |
| Python | 3.10 or later (only for the Python package) |
| C toolchain | Required to build the shared library for Python (cgo) |
| Platforms | Windows, Linux, macOS |

### Go

```bash
go get github.com/djinn-soul/gopptx
```

That is all — the Go package has no cgo requirement and no native dependency.

### Python

The Python package calls into a Go shared library, which must be built first.

```powershell
# Windows (PowerShell)
.\scripts\build_python.ps1
pip install -e .
```

```bash
# Linux / macOS
./scripts/build_python.sh
pip install -e .
```

Optionally install `orjson` for faster JSON encoding on the bridge:

```bash
pip install orjson
```

The library filename is platform-specific — `gopptx.dll`, `libgopptx.so` or `libgopptx.dylib`.
Set `GOPPTX_LIB_PATH` if you keep it outside the package directory. See
[Installation](docs/installation.md) for details and
[Troubleshooting](docs/troubleshooting.md) if the library is not found.

---

## Quickstart — Python

### Create a deck

```python
from gopptx import ChartType, Inches, Presentation

with Presentation.new("Quarterly Update") as pres:
    pres.add_bullet_slide("Highlights", ["Revenue +12%", "Retention +4%"])

    slide = pres.add_slide("Numbers")
    pres.add_table_from_rows(
        slide.index,
        [["Region", "Revenue"], ["EMEA", "4.1M"], ["APAC", "2.8M"]],
        bounds=(Inches(0.5), Inches(1.5), Inches(4.0), Inches(2.0)),
    )
    pres.add_chart(
        slide.index,
        ChartType.COLUMN,
        ["Q1", "Q2", "Q3", "Q4"],
        [12.0, 15.5, 18.0, 21.0],
        bounds=(Inches(5.0), Inches(1.5), Inches(4.0), Inches(3.0)),
        title="Quarterly trend",
    )

    pres.save("quarterly.pptx")
```

Four things worth knowing straight away:

- **All geometry is in EMU** (English Metric Units, 914 400 per inch). Raw numbers like
  `(40, 120, 600, 220)` are 600 EMU wide — smaller than a pixel, and the shape will appear to
  be missing. Always wrap coordinates in `Inches()`, `Point()` or `Emu()`.
- `Presentation` is a context manager. Leaving the `with` block releases the native handle;
  `save()` writes the file. Both are explicit — nothing is written implicitly.
- `Presentation.new(title)` already contains a title slide at index `0`. Your first
  `add_slide()` becomes slide `1`.
- Pass a `ChartType` member, not a bare string. Strings still work but are deprecated and will
  be rejected in a future release.

### The object layer

Alongside the flat `pres.*` methods there is a navigable object layer:

```python
with Presentation("deck.pptx") as pres:
    for slide in pres.slides:
        print(slide.index, slide.title, len(slide.list_shapes()))

    first = pres.slides[0]
    first.set_transition("fade")
    first.notes = "Open with the revenue number."
```

## Quickstart — Go

### Generate a deck in one call

```go
package main

import (
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Hello from gopptx").AddBullet("Created with gopptx"),
	}

	data, err := pptx.CreateWithSlides("My Deck", slides)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("output.pptx", data, 0o600); err != nil {
		panic(err)
	}
}
```

### Build a deck fluently

`SlideBuilder` mutates in place, so a dropped return value cannot silently discard content —
unlike `SlideContent`'s value-receiver methods.

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	intro := pptx.NewSlideBuilder("Agenda")
	intro.AddBullet("Results")
	intro.AddBullet("Outlook")
	intro.WithNotes("Keep this to two minutes.")

	err := pptx.NewPresentationBuilder("Quarterly Update").
		WithTheme(pptx.ThemeCorporate).
		WithSlideSize(pptx.SlideSize16x9()).
		AddTitleSlide("FY26 Q3").
		AddSlide(intro.Build()).
		WriteToFile("quarterly.pptx")
	if err != nil {
		panic(err)
	}
}
```

See the [Go library guide](docs/guides/go-library.md) for shapes, charts, tables and layout helpers.

---

## Editing an existing deck

### Python

```python
from gopptx import Presentation

with Presentation("input.pptx") as pres:
    pres.set_slide_title(0, "Updated Title")
    pres.find_and_replace("Draft", "Final")
    pres.update_chart_data_by_index(
        2,
        0,
        {"categories": ["Q1", "Q2"], "series": [{"name": "Revenue", "values": [10.0, 14.0]}]},
    )
    pres.save("edited.pptx")
```

### Go

```go
p, err := pptx.Open("input.pptx")
if err != nil {
	panic(err)
}
defer p.Close()

p.SetTitle("Updated Title")
p.SetAuthor("Reporting Bot")

if err := p.SaveAs("edited.pptx"); err != nil {
	panic(err)
}
```

`pptx.Open` returns a metadata-and-charts facade. For full shape-level editing use
`pptx.OpenEditor(path)`, which returns the `PresentationEditor` that also backs the JSON bridge.

---

## Exporting to PDF and HTML

```python
from gopptx import HTMLOptions, PDFOptions, Presentation

with Presentation("deck.pptx") as pres:
    pres.export_pdf("deck.pdf", PDFOptions(driver="auto"))
    pres.export_html("deck.html", HTMLOptions(embed_images=True))
```

`export_pdf` is the current name; `save_as_pdf` remains as a deprecated alias. There are
three real drivers, plus `auto` which
picks between them:

| Driver | Requires | Notes |
|---|---|---|
| `auto` (default) | — | Tries LibreOffice, then PowerPoint, then falls back to `native`. Use this in production. |
| `native` | nothing | Pure Go renderer. No external process, works in containers. Still maturing on layout-heavy decks — it emits a warning. |
| `libreoffice` | `soffice` on `PATH` | Highest general fidelity without Office. On Windows add `C:\Program Files\LibreOffice\program` to `PATH`. |
| `powerpoint` | Microsoft PowerPoint + PowerShell (Windows) | Ground truth for fidelity; used to grade the native renderer. |

Full detail in the [export guide](docs/guides/export.md).

---

## Throughput and batching

Every Python call crosses a Python → C → Go boundary. For write-heavy loops, batch them so the
crossing happens once:

```python
from gopptx import Presentation

with Presentation.new("Batch Demo") as pres:
    with pres.batch(stop_on_error=True) as batch:
        for i in range(500):
            batch.add_slide(f"Slide {i}")
    pres.save("batch.pptx")
```

Read operations are rejected inside a `batch()` block by design — a buffered write has not
executed yet, so a read would see stale state. Use `execute_batch()` with an explicit command
list when you need reads and writes interleaved. See
[Batch execution](docs/guides/batch-execution.md).

---

## Architecture

```
┌──────────────────────────────┐
│  python/gopptx               │  typed Python API — Presentation, Slide, Shape, Table, Chart
│  (ctypes → JSON envelopes)   │
└───────────────┬──────────────┘
                │  JSON command envelope, 179 operations
┌───────────────▼──────────────┐
│  bindings/c                  │  handle-based C ABI (cgo, c-shared)
└───────────────┬──────────────┘
                │
┌───────────────▼──────────────┐
│  pkg/pptx/editor             │  command dispatch over a loaded presentation
│  pkg/pptx                    │  high-level Go API — builders, shapes, charts, tables, text
│  internal/pptxxml            │  OOXML serialisation
└──────────────────────────────┘
```

The JSON envelope is a stable contract, so any language that can call a C function or shell out
to the CLI can drive the engine. Operation identifiers and payload shapes are defined in
`pkg/pptx/editor/opspec.go`; the Python constants in `python/gopptx/_ops_constants.py` are
generated from it. See [Bridge operations](docs/reference/bridge-operations.md).

---

## Documentation map

**Start here**

| Page | What it covers |
|---|---|
| [Installation](docs/installation.md) | Go module, Python shared library, environment variables |
| [Quickstart](docs/quickstart.md) | First deck in both languages |
| [Core concepts](docs/concepts.md) | Handles, sessions, the command envelope, batching, units |

**Guides**

| Page | What it covers |
|---|---|
| [Python library](docs/guides/python-library.md) | The full Python workflow, method by area |
| [Go library](docs/guides/go-library.md) | Builders, the editor, shapes, charts, layout helpers |
| [Batch execution](docs/guides/batch-execution.md) | Throughput patterns and their limits |
| [Export](docs/guides/export.md) | PDF drivers, HTML options, fidelity notes |

**Reference**

| Page | What it covers |
|---|---|
| [API overview](docs/api-reference.md) | The three surfaces and when to use each |
| [Python `Presentation` API](docs/reference/python-presentation-api.md) | All 208 methods, grouped |
| [Go API](docs/reference/go-api.md) | Packages, types and constructors |
| [Bridge operations](docs/reference/bridge-operations.md) | All 179 JSON operations |
| [Feature matrix](docs/reference/feature-matrix.md) | Honest comparison against python-pptx |

**Examples**

- [`examples/`](examples/) — 96 runnable Go and Python programs
- [Showcase](docs/showcase/usages/index.md) — 30 annotated recipes, simple → complex, with screenshots

---

## Development

```bash
task build:go        # build the C shared library for Python
task test            # Go + Python test suites
task lint            # golangci-lint, ruff, basedpyright, generated-code drift
task generate        # regenerate everything derived from Go declarations
task docs:serve      # MkDocs Material on http://localhost:8000
```

Parts of the Python surface — operation constants, chart-type and shape-type enums, the slide
builder — are **generated from the Go declarations**. Edit the Go source and run `task generate`;
never hand-edit the generated files. `task check:generated` fails the build if they drift.

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the full setup, code standards and PR checklist.

## License

Apache License 2.0 — see [`LICENSE`](LICENSE).
