# API overview

gopptx exposes three surfaces. They are layers of the same engine, not alternatives to each
other.

```
Python API      →  JSON bridge  →  Go API
gopptx.Presentation                pkg/pptx, pkg/pptx/editor
```

Anything the Python API can do, the bridge can do. Anything the bridge can do, the Go API can
do. The reverse is not true — the Go API is the widest surface.

## 1. Python API

**Use it when** you are writing an application, a report generator or a pipeline in Python.

```python
from gopptx import Presentation

with Presentation.new("Deck") as pres:
    pres.add_bullet_slide("Highlights", ["One", "Two"])
    pres.save("deck.pptx")
```

- Entry point: `gopptx.Presentation` — 208 methods and properties
- Object layer: `pres.slides[i]` → `Slide` (90 members), `Chart` (47), `Table` (23)
- Builders: `PresentationBuilder`, `ShapeBuilder`, `RunBuilder`
- Typed throughout — the package ships `py.typed` and `.pyi` stubs

| | |
|---|---|
| Guide | [Python library](guides/python-library.md) |
| Reference | [Python `Presentation` API](reference/python-presentation-api.md) |

## 2. Go API

**Use it when** you are writing Go, embedding the engine, or working on gopptx itself.

Two layers:

| Layer | Package | For |
|---|---|---|
| High-level | `pkg/pptx` (plus `shapes`, `charts`, `tables`, `text`, `styling`, `templates`, `elements`) | Fluent deck construction and editing |
| Engine | `pkg/pptx/editor` | The command dispatcher — the same object the bridge drives |

```go
err := pptx.NewPresentationBuilder("Deck").
	AddTitleSlide("Hello").
	WriteToFile("deck.pptx")
```

`pkg/pptx` re-exports the sub-package types, so one import usually suffices.

| | |
|---|---|
| Guide | [Go library](guides/go-library.md) |
| Reference | [Go API](reference/go-api.md), or `go doc github.com/djinn-soul/gopptx/pkg/pptx` |

## 3. JSON bridge

**Use it when** you are integrating a language that is neither Go nor Python, or you need an
operation the typed API does not wrap yet.

```json
{"api_version": 1, "op": "add_slide", "payload": {"title": "Agenda"}}
```

- 179 operations
- Handle-based C ABI in `bindings/c`
- Stable contract: `pkg/pptx/editor/opspec.go` is the source of truth, and the Python
  constants are generated from it

| | |
|---|---|
| Reference | [Bridge operations](reference/bridge-operations.md) |
| ABI | [C bridge guidance](architecture/cbridge-guidance.md) |
| Batching | [Batch execute envelope](architecture/batch_execute_envelope.md) |

## Choosing between them

| Situation | Surface |
|---|---|
| Python application | Python API |
| Go service or CLI | Go high-level API |
| Very high volume from Python | Python API with `batch()` |
| An operation with no typed wrapper | `pres.execute(op, payload)` |
| Another language entirely | JSON bridge over the C ABI |
| Working on gopptx | Go API, `pkg/pptx/editor` |

## Cross-cutting rules

These hold on every surface:

- **Geometry is EMU** — 914 400 per inch. See [Units](concepts.md#units-and-geometry).
- **Indices are zero-based.**
- **Shape identity is `shape_id`**, the value returned at creation — not a position.
- **Colours are RGB hex without `#`.**
- **Nothing is written implicitly** — you call `save`, `to_bytes` or `WriteToFile`.

## Generated code

The Python operation constants, the `ChartType` and `ShapeType` enums, the slide builder surface
and the cgo header are all generated from Go declarations. Edit the Go source, run
`task generate`; `task check:generated` fails the build on drift. See
[Core concepts](concepts.md#generated-code).
