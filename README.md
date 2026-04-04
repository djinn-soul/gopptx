# gopptx

High-performance PowerPoint (PPTX) engine powered by Go — with a Python library and a stable JSON command bridge for cross-language use.

## What is gopptx?

`gopptx` is a Go library for generating and manipulating PPTX files at scale.
It exposes a stable [JSON command bridge](docs/architecture/bridge-phase1-ops.md) so it can be driven from Python or any other language.

**Key capabilities**

| Feature | Description |
|---|---|
| Fast & Efficient | High-performance PPTX generation and manipulation |
| Concurrency-Oriented | Built for large workloads and batch processing |
| Lazy-Loading | Low memory footprint for heavy file editing |
| Rich Diagrams | Native Mermaid rendering (12+ types) and SmartArt |
| Deep Editing | Full control over shapes, images, charts, and tables |
| Accessibility | Alt-text, sections, comments, and speaker notes |
| C-Bridge | Stable JSON command API for Python and cross-language interop |

---

## Quick Start — Python

### Install

1. Build the Go shared library:
   ```powershell
   .\scripts\build_python.ps1
   ```
2. Install the Python package:
   ```bash
   pip install -e .
   ```

### Basic usage

```python
from gopptx import Presentation

with Presentation.new("Hello World") as pres:
    pres.add_slide("Hello from Python")
    pres.save("output.pptx")
```

### Batch operations (recommended for performance)

Reduce Python↔Go boundary crossings by grouping operations:

```python
from gopptx import Presentation

with Presentation.new("Batch Demo") as pres:
    with pres.batch(stop_on_error=True) as batch:
        for i in range(100):
            batch.add_slide(f"Slide {i}")
```

---

## Quick Start — Go

```bash
go get github.com/djinn-soul/gopptx
```

## PDF Export Warning

- PPTX to PDF export is currently **experimental** for visual fidelity.
- The Go-native PDF renderer (`driver="native"`) is experimental and can differ from PowerPoint output.
- For production output, prefer `auto` (default), which now tries `LibreOffice` / `PowerPoint` first and uses native only as fallback.
- Use native explicitly only when you accept rendering differences (notably in advanced SmartArt/layout-heavy decks).
- Windows prerequisites:
  - `driver="powerpoint"` requires Microsoft PowerPoint desktop installation (COM automation) and PowerShell (`powershell` or `pwsh`).
  - `driver="libreoffice"` requires LibreOffice and `soffice` available on `PATH`.
  - If LibreOffice is installed but `soffice` is not recognized, add `C:\Program Files\LibreOffice\program` to your `PATH`.

---

## JSON Command Bridge

All bridge operations use a JSON envelope.

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

Full operation reference: [bridge-phase1-ops.md](docs/architecture/bridge-phase1-ops.md) (145+ commands).

---

## Documentation & Examples

| Resource | Description |
|---|---|
| [Examples Index](examples/README.md) | Sequential walkthrough of 90+ features |
| [Python Guide](docs/guides/python-library.md) | Detailed Python usage and API reference |
| [Operation Reference](docs/architecture/bridge-phase1-ops.md) | All 145+ supported bridge commands |

---

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for setup, guidelines, and the PR checklist.

## License

Apache License 2.0 — see [`LICENSE`](LICENSE).
