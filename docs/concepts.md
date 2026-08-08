# Core concepts

Seven ideas explain almost every surprise people hit with gopptx. Read the units section even if
you skip the rest.

## Units and geometry

**All geometry crossing the API is in EMU** — English Metric Units, the native OOXML unit.

| Unit | EMU |
|---|---|
| 1 inch | 914 400 |
| 1 centimetre | 360 000 |
| 1 point | 12 700 |
| 1 pixel at 96 DPI | 9 525 |

| Slide size | EMU | Inches |
|---|---|---|
| 4:3 (the default for a new deck) | `9 144 000 × 6 858 000` | 10 × 7.5 |
| 16:9 widescreen | `12 192 000 × 6 858 000` | 13.33 × 7.5 |

Set it with `pres.set_slide_size(width, height)` in Python, or
`NewPresentationBuilder(...).WithSlideSize(pptx.SlideSize16x9())` in Go.

This is the single most common mistake:

```python
# WRONG — a table 600 EMU wide is 0.00066 inch. It renders as nothing.
pres.add_table_from_rows(0, rows, bounds=(40, 120, 600, 220))

# RIGHT
from gopptx import Inches
pres.add_table_from_rows(0, rows, bounds=(Inches(0.5), Inches(1.5), Inches(4.0), Inches(2.0)))
```

The file still saves, PowerPoint still opens it, and nothing tells you the shape is there but
invisible — which is why it costs people an afternoon.

### Helpers

=== "Python"

    ```python
    from gopptx import Emu, Inches, Point

    Inches(1.5)   # 1371600
    Point(12)     # 152400
    Emu(914400)   # 914400  (identity — for when you already have EMU)
    ```

=== "Go"

    ```go
    pptx.Inches(1.5)        // styling.Length
    pptx.Centimeters(2.0)
    pptx.Points(12)
    pptx.Emu(914400)
    ```

    Relative dimensions are also available where a length may be expressed against the slide:

    ```go
    pptx.Absolute(pptx.Inches(2))   // fixed
    pptx.PercentOf(50)              // half the slide
    pptx.Ratio(0.25)                // a quarter of the slide
    ```

### The one exception

Go's convenience shape constructors take **inches as plain `float64`**, not `Length`:

```go
shapes.NewRectangle(1.0, 1.8, 2.2, 1.3)          // inches
pptx.NewShape("rect", pptx.Inches(1.0), pptx.Inches(1.8),
    pptx.Inches(2.2), pptx.Inches(1.3))          // EMU Length
```

Font sizes are a separate scale again: points, as integers, everywhere.

## Architecture layers

```
python/gopptx      typed Python API — Presentation, Slide, Shape, Table, Chart
      │            encodes a JSON command envelope
bindings/c         handle-based C ABI (cgo, c-shared)
      │
pkg/pptx/editor    command dispatch against a loaded presentation
pkg/pptx           high-level Go API — builders, shapes, charts, tables, text
internal/pptxxml   OOXML serialisation
```

Python holds no presentation state. Every call is forwarded to the Go engine, which owns the
document. That is why Python performance is Go performance, and why boundary crossings — not
document operations — are the thing worth optimising.

## Handles and sessions

A presentation is a **handle** into Go-owned memory, not a Python object graph.

```python
with Presentation.new("Deck") as pres:   # allocates a handle
    pres.add_slide("One")                # mutates in Go memory
    pres.save("out.pptx")                # serialises to disk
# handle released on exit
```

| Entry point | Meaning |
|---|---|
| `Presentation.new(title)` | New deck — **already contains a title slide at index 0** |
| `Presentation(path)` | Open an existing file |
| `Presentation.open_deck(path)` | Same, as a named constructor |
| `Presentation.from_template(...)` | Start from a template deck |
| `pres.save(path)` | Write bytes; may be called more than once |
| `pres.to_bytes()` | Serialise without touching disk |
| `pres.close()` | Release the handle; the context manager does this for you |

Nothing is written implicitly. If you never call `save()` or `to_bytes()`, the work is discarded
when the handle closes.

Use the context manager. A leaked handle keeps Go memory alive for the life of the process.

## The command envelope

Every operation — from Python, from C, from any other client — is one JSON envelope in and one
out.

**Request**

```json
{
  "api_version": 1,
  "request_id": "b1c2…",
  "op": "add_slide",
  "payload": {"title": "Agenda"}
}
```

**Response**

```json
{
  "ok": true,
  "request_id": "b1c2…",
  "result": {"index": 2}
}
```

**Error response**

```json
{
  "ok": false,
  "request_id": "b1c2…",
  "error": {"code": "INVALID_INDEX", "message": "slide index 9 out of range", "details": null}
}
```

Python raises `GopptxError` carrying the same `code` on its `.code` attribute.

### Error codes

| Code | Meaning |
|---|---|
| `INVALID_JSON` | The envelope itself did not parse |
| `UNSUPPORTED_VERSION` | `api_version` is not supported |
| `UNKNOWN_OP` | No such operation |
| `INVALID_PAYLOAD` | Payload failed structural validation |
| `MISSING_FIELD` | A required payload field was absent |
| `INVALID_TYPE` | A field had the wrong JSON type |
| `INVALID_VALUE` | A field parsed but the value is not acceptable |
| `INVALID_INDEX` | Slide, shape, row or column index out of range |
| `INVALID_HANDLE` | The handle is closed or was never valid |
| `INVALID_BATCH_ITEM` | A command inside `batch_execute` is not batchable (e.g. a nested batch) |
| `OP_FAILED` | The operation ran and failed |
| `MARSHAL_ERROR` | The result could not be encoded |
| `INTERNAL_ERROR` | A bug — please report it |

You can call any operation directly, bypassing the typed wrapper:

```python
result = pres.execute("add_slide", {"title": "Raw"})
```

The 179 operation names live in `pkg/pptx/editor/opspec.go` and are mirrored into
`python/gopptx/_ops_constants.py` by a generator. See
[Bridge operations](reference/bridge-operations.md).

### Result keys

Results are JSON produced from Go structs. Handlers whose structs carry JSON tags return
snake_case (`code`, `path`); those without tags return the Go field names (`ID`, `RelID`).
The Python layer adds snake_case aliases to every response, keeping the original keys, so
`shape["id"]` and `shape["ID"]` both work.

## Batching

Each Python → C → Go crossing costs more than the operation it carries. For write-heavy loops,
send many operations per crossing.

```python
with pres.batch(stop_on_error=True) as batch:
    for i in range(500):
        batch.add_slide(f"Slide {i}")
```

or explicitly:

```python
from gopptx import ops

commands = [{"op": ops.OP_ADD_SLIDE, "payload": {"title": f"Slide {i}"}} for i in range(500)]
results = pres.execute_batch(commands, stop_on_error=False)
```

**Reads are rejected inside a `batch()` block.** The buffered writes have not run yet, so any
read would answer from stale state. Move reads outside the block, or use `execute_batch()` with
an explicit command list where you control the ordering.

`stop_on_error=True` aborts at the first failure. With `False`, every command is attempted and
`execute_batch` returns one result mapping per command — `{"ok": true, "op": "add_slide",
"result": {"index": 2}}`. Check `ok` per item rather than assuming success.

See [Batch execution](guides/batch-execution.md).

## Generated code

Parts of the Python surface are **generated from the Go declarations** and must never be edited
by hand:

| Generated file | Generated from |
|---|---|
| `python/gopptx/_ops_constants.py`, `ops.pyi` | `pkg/pptx/editor/opspec.go` |
| The `ChartType` enum | `XLChartType` in `pkg/pptx/enums` |
| The `ShapeType` enum | The shape preset declarations |
| The slide builder surface | `elements.SlideContent` |
| `python/gopptx/gopptx.h` | `bindings/c/bridge.go` (cgo) |

Change the Go source and run `task generate`. `task check:generated` fails the build when they
drift — that check exists because hand-edits used to survive for weeks.

## Validation

gopptx can check its own output before you ship it:

```python
issues = pres.validate()          # structural findings; empty list when clean
report = pres.repair()            # fixes what it safely can
```

Each issue is a dict with `code`, `severity`, `path`, `description` and `repairable`.
`repair()` returns `repaired`, `repaired_count`, `unrepaired` and `unrepaired_count`.

Validation covers OPC package rules, required parts, relationship targets and content types —
the classes of defect that make PowerPoint show its "repair" dialog. It validates the package as
it *would be saved*, so slides added since opening are not reported as orphans.

Validation is not a substitute for opening the file. A deck can be structurally valid and still
visually wrong — which is exactly what the EMU trap produces.
