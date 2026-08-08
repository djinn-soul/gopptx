# Troubleshooting

Ordered roughly by how often people hit them.

## A shape I added is not there

**Symptom.** The file saves, PowerPoint opens it without complaint, and the textbox, table,
image or chart you added is nowhere on the slide.

**Cause.** Geometry is in **EMU** — 914 400 per inch. A shape at `(40, 120, 600, 220)` is
0.00066 inch wide. It is on the slide; it is just far too small to see.

**Fix.**

```python
from gopptx import Inches

# not this
pres.add_table_from_rows(0, rows, bounds=(40, 120, 600, 220))

# this
pres.add_table_from_rows(0, rows, bounds=(Inches(0.5), Inches(1.5), Inches(4.0), Inches(2.0)))
```

In Go use `pptx.Inches()`, `pptx.Points()`, `pptx.Centimeters()` — except for the convenience
shape constructors (`shapes.NewRectangle` and friends), which take **inches as plain floats**.

See [Units](concepts.md#units-and-geometry).

---

## My first slide is not slide 0

**Symptom.** You call `Presentation.new("Deck")` then `add_slide("Intro")`, and `Intro` lands at
index 1.

**Cause.** `Presentation.new(title)` creates a deck that already contains a title slide.

**Fix.** Either use the returned `Slide` object rather than assuming an index:

```python
slide = pres.add_slide("Intro")
pres.add_textbox(slide.index, ...)
```

or retitle slide 0 instead of adding another: `pres.set_slide_title(0, "Intro")`.

---

## Bridge library not found

```
Could not find shared library ... Please build it first.
```

**Fix.**

1. Build it:

    === "Windows"
        ```powershell
        .\scripts\build_python.ps1
        ```
    === "Linux / macOS"
        ```bash
        ./scripts/build_python.sh
        ```

    Or `task build:go`.

2. Confirm the platform file exists in `python/gopptx/`:

    | Platform | File |
    |---|---|
    | Windows | `gopptx.dll` |
    | Linux | `libgopptx.so` |
    | macOS | `libgopptx.dylib` |

3. If it lives elsewhere, point the loader at it:

    ```bash
    export GOPPTX_LIB_PATH=/path/to/libgopptx.so
    ```

---

## `ModuleNotFoundError: No module named 'gopptx'`

Build the shared library first, then install from the repository root:

```bash
task build:go
pip install -e .
```

---

## Go changes have no effect in Python

The Python package binds a compiled shared library. Editing Go source does nothing until you
rebuild it:

```bash
task build:go
```

If the change touched `bindings/c/bridge.go`, the generated `python/gopptx/gopptx.h` changes too
and must be staged alongside it — otherwise the pre-commit hooks fail with a confusing error.

---

## `read operation '…' is not allowed inside batch()`

```python
with pres.batch() as batch:
    batch.add_slide("A")
    count = pres.slide_count      # GopptxError
```

**Cause.** Buffered writes have not executed yet, so a read would answer from stale state. The
error is deliberate — a silently wrong answer would be worse.

**Fix.** Move the read outside the block, or use `execute_batch()` with an explicit command list
where you control the ordering. See [Batch execution](guides/batch-execution.md).

---

## `unknown theme preset "…"`

`apply_theme` and `set_global_theme_preset` accept two vocabularies:

```python
from gopptx import THEME_CORPORATE
from gopptx.presentation.theme import get_theme

pres.apply_theme(THEME_CORPORATE)      # gopptx names: Corporate, Modern, Vibrant,
pres.apply_theme("integral")           #   Dark, Nature, Tech, Carbon, Office
pres.apply_theme(get_theme("ocean"))   # Office presets: office, office2013, facet,
                                       #   integral, ion, retrospect, slice, wisp
```

Matching ignores case and separators. The error message lists every accepted name.

`get_theme(...)` returns a `Theme` object and takes a different set — `aurora`, `ocean`,
`sunset`, `forest` — applied as a colour scheme plus a font scheme.

---

## `DeprecationWarning: chart_type should be a ChartType member`

```python
from gopptx import ChartType

pres.add_chart(0, ChartType.COLUMN, categories, values, bounds=...)
```

Bare strings still work today and will be rejected in a future release.

---

## `cannot import name '…' from 'gopptx'`

`ChartType`, `PDFOptions`, `HTMLOptions`, `SlideLayoutType`, `Inches`, `Point` and `Emu` are all
exported at the top level:

```python
from gopptx import ChartType, Emu, HTMLOptions, Inches, PDFOptions, Point, SlideLayoutType
```

If your editor disagrees, the installed package is older than the docs — reinstall with
`pip install -e .`.

---

## `DeprecationWarning: save_as_pdf() is deprecated`

Use `export_pdf`, which pairs with `export_html`:

```python
from gopptx import PDFOptions

pres.export_pdf("deck.pdf", PDFOptions(driver="auto"))
```

`save_as_pdf` still works and forwards to `export_pdf`.

---

## PDF export fails or looks wrong

| Symptom | Likely cause | Fix |
|---|---|---|
| `soffice not found` | LibreOffice not on `PATH` | Add `C:\Program Files\LibreOffice\program`, or use `driver="native"` |
| COM error on Windows | PowerPoint not installed, or no `powershell`/`pwsh` | Use `auto`, `libreoffice` or `native` |
| Fonts look wrong with `native` | The renderer could not find the deck's fonts and substituted | Pass `PDFOptions(font_paths=[...])` |
| `UserWarning: PDF driver 'native' is experimental` | Informational | The PDF is still produced; use `auto` if fidelity matters |

Rasterise the result with `pypdfium2` when comparing renders — LibreOffice drops clips and
substitutes fonts, so it disagrees with itself.

See the [export guide](guides/export.md).

---

## SmartArt shows `[Text]` placeholders

**Likely causes.**

- A stale PowerPoint process still holds the previous version of the file.
- You opened an older artifact, not the one just generated.

**Fix.**

1. Close every PowerPoint window before running the generator.
2. Open the newly generated file directly rather than reusing an open window.
3. If it persists, look at the XML:

    ```bash
    unzip -p output.pptx ppt/diagrams/drawing2.xml | grep -o '<a:t>[^<]*</a:t>'
    ```

gopptx generates SmartArt from the data model rather than filling fixed template slots, so
PowerPoint re-lays it out — an empty data tree is the one thing that produces `[Text]`.

---

## Python throughput is lower than expected

**Cause.** Each call crosses Python → C → Go. One operation per crossing is the slow path.

**Fix.**

```python
with pres.batch(stop_on_error=True) as batch:
    for row in rows:
        batch.add_slide(row.title)
```

Also: install `orjson`, and call `save()` once at the end rather than inside the loop —
`save()` serialises the whole package each time.

---

## Save fails

```
RuntimeError: save failed — handle may be closed or path is not writable
```

**Causes.** The handle was already closed; the destination directory does not exist or is not
writable; an earlier invalid mutation left the deck in a bad state.

**Fix.** Use the context manager so handle lifetime is deterministic, and check the destination
exists:

```python
with Presentation.new("Deck") as pres:
    pres.add_slide("Intro")
    pres.save("out/deck.pptx")     # out/ must exist
```

---

## PowerPoint offers to repair the file

Run validation and repair what it can:

```python
with Presentation("broken.pptx") as pres:
    for issue in pres.validate():
        print(issue["code"], issue["path"], issue["description"])
    print(pres.repair())
    pres.save("fixed.pptx")
```

`validate()` returns a list — empty when the package is clean — and checks the package as it
would be saved, so slides added since opening are not reported as orphans.

Validation covers OPC package rules, required parts, relationship targets and content types —
precisely the class of defect that triggers the repair dialog.

---

## Tests pass but PowerPoint disagrees

A deck can be structurally valid and visually wrong — the EMU trap is exactly that. Structural
tests cannot catch it. Open the output in real PowerPoint before trusting a layout, and
pixel-diff PDF renders when you need a regression signal.

---

## COM-related test failures

`pkg/pptx/presentation/protection` includes a scenario requiring COM, which is Windows-only and
skipped automatically elsewhere. If it fails on Windows, confirm Microsoft Office or the OpenXML
SDK is installed.
