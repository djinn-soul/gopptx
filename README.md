# gopptx

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

PowerPoint (PPTX) engine written in Go — with a Python library and a stable JSON command bridge for cross-language use.

---

## Overview

`gopptx` is a Go library for generating and manipulating PPTX files. It is designed for correctness and ease of use across languages. A stable [JSON command bridge](docs/architecture/bridge-phase1-ops.md) lets you drive it from Python, or any other language without rewriting logic.

### Why gopptx?

gopptx combines the native performance of Go with cross-language accessibility, making it ideal for high-throughput presentation generation. Use it for automating business reports, creating dynamic slides from data sources, or integrating PPTX creation into your existing workflows without sacrificing speed or reliability.

### Key Capabilities

| Feature | Description |
|---|---|
| Rich Diagrams | Native Mermaid rendering (12+ types) and SmartArt |
| Deep Editing | Full control over shapes, images, charts, tables, and text |
| Themes & Layouts | Built-in themes, color schemes, font schemes, and layout helpers |
| Accessibility | Alt-text, sections, comments, and speaker notes |
| Export | PDF export (native, LibreOffice, or PowerPoint driver) |
| Cross-Language | Stable JSON command API for Python and other language interop |

---

## Prerequisites

- **Go**: 1.25.9 or later
- **Python**: 3.10 or later (required for Python bindings)
- **Platforms**: Windows, Linux, macOS

---

## Installation

### Go

```bash
go get github.com/djinn-soul/gopptx
```

### Python

The Python package uses a high-performance Go shared library for PPTX processing. You need to build this library before installing the package.

1. **Build the Go shared library:**
   - **Windows (PowerShell):**
     ```powershell
     .\scripts\build_python.ps1
     ```
   - **Linux/macOS (Bash):**
     ```bash
     ./scripts/build_python.sh
     ```

2. **Install the Python package:**
   ```bash
   pip install -e .
   ```

---

## Quick Start

### Python Setup

Build the native bridge and install the package locally:

```powershell
# Windows
.\scripts\build_python.ps1
pip install -e .
```

```bash
# Linux/macOS
./scripts/build_python.sh
pip install -e .
```

### Python

Create and save a new deck:

```python
from gopptx import Presentation

with Presentation.new("Quarterly Update") as pres:
    pres.add_slide("Overview")
    pres.add_bullet_slide("Highlights", ["Growth +12%", "Retention +4%"])
    pres.save("output.pptx")
```

Open and edit an existing deck:

```python
from gopptx import Presentation

with Presentation("input.pptx") as pres:
    pres.set_slide_title(0, "Updated Title")
    pres.add_slide("New Closing Slide")
    pres.save("edited.pptx")
```

Run many writes in one batch:

```python
from gopptx import Presentation

with Presentation.new("Batch Demo") as pres:
    with pres.batch(stop_on_error=True) as batch:
        for i in range(100):
            batch.add_slide(f"Slide {i}")
    pres.save("batch.pptx")
```

### Go

Install and run:

```bash
go get github.com/djinn-soul/gopptx
go run ./your-main.go
```

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

### Next Docs

- Python API reference: [docs/reference/python-presentation-api.md](docs/reference/python-presentation-api.md)
- Go API reference: [docs/reference/go-api.md](docs/reference/go-api.md)
- Full quickstart page: [docs/quickstart.md](docs/quickstart.md)
- Runnable examples: [examples/README.md](examples/README.md)

---

## JSON Command Bridge

All bridge operations use a JSON envelope. This is the primary interface for cross-language use.

**Request**
```json
{
  "api_version": 1,
  "request_id": "uuid-123",
  "op": "add_slide",
  "payload": { "title": "New Slide", "layout": "TITLE_AND_CONTENT" }
}
```

**Response**
```json
{
  "ok": true,
  "result": { "index": 1 },
  "request_id": "uuid-123"
}
```

Full operation reference: [bridge-phase1-ops.md](docs/architecture/bridge-phase1-ops.md) — 145+ supported commands.

---

## PDF Export

PDF export supports three drivers: `native`, `libreoffice`, and `powerpoint`. The default `auto` mode tries LibreOffice or PowerPoint first, and falls back to the native renderer.

> **Note:** The native Go PDF renderer is experimental. It may differ from PowerPoint output for advanced SmartArt or layout-heavy decks. Use `auto` for production unless you specifically need the native renderer.

**Example (Python):**

```python
from gopptx import Presentation

with Presentation("input.pptx") as pres:
    pres.export_pdf("output.pdf", driver="auto")
```

**Driver requirements (Windows):**
- `powerpoint` — requires Microsoft PowerPoint installed (COM automation) and PowerShell (`powershell` or `pwsh`).
- `libreoffice` — requires LibreOffice with `soffice` on `PATH`. If `soffice` is not recognized, add `C:\Program Files\LibreOffice\program` to your `PATH`.

---

## Documentation

| Resource | Description |
|---|---|
| [Quickstart](docs/quickstart.md) | Python and Go quickstart examples |
| [Examples Index](examples/README.md) | 90+ runnable examples covering all features |
| [Python Guide](docs/guides/python-library.md) | Full Python API reference and usage |
| [Go API Reference](docs/reference/go-api.md) | Go library reference |
| [Bridge Operations](docs/architecture/bridge-phase1-ops.md) | All 145+ JSON bridge commands |
| [Showcase](docs/showcase/usages/index.md) | Real-world usage patterns (simple → complex) |

---

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for setup instructions, code guidelines, and the PR checklist.

## License

Apache License 2.0 — see [`LICENSE`](LICENSE).
