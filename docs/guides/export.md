# Export guide

gopptx exports a deck to PDF, HTML and flat XML, and can convert it to grayscale or compress it
in place.

## PDF

```python
from gopptx import PDFOptions, Presentation

with Presentation("deck.pptx") as pres:
    path = pres.export_pdf("deck.pdf", PDFOptions(driver="auto"))
```

With no `output_path` it writes `presentation.pdf` into the working directory. It returns the
absolute path written.

`export_pdf` pairs with `export_html`: `export_*` converts the deck into a format gopptx cannot
reopen, while `save` / `save_flat_xml` write something PowerPoint reads back. The older name
`save_as_pdf` still works but emits a `DeprecationWarning`.

### Choosing a driver

There are three real drivers, plus `auto` which picks between them.

| Driver | External requirement | Fidelity | Use when |
|---|---|---|---|
| `auto` *(default)* | none | best available | Production. Tries LibreOffice, then PowerPoint, then falls back to `native`. |
| `native` | none | good, still improving | Containers, CI, anywhere you cannot install Office or LibreOffice. |
| `libreoffice` | `soffice` on `PATH` | high | You want consistent output across platforms without Office. |
| `powerpoint` | Microsoft PowerPoint + PowerShell, Windows only | ground truth | Reference renders, and grading the native renderer against reality. |

```python
PDFOptions(driver="native", font_paths=["/usr/share/fonts/truetype/dejavu"])
```

`font_paths` only affects the native renderer — it is where that renderer looks for the fonts a
deck asks for. Without them it substitutes, and substituted metrics change line breaks.

### About the native renderer

The native renderer is pure Go: no subprocess, no Office, no browser. It draws shapes, preset
geometry, text with real font metrics, tables, charts, images with crops, gradients, SmartArt
from its cached layout, hyperlinks and speaker notes.

It emits a `UserWarning` when selected directly, because it is the one driver whose output can
still diverge from PowerPoint on layout-heavy decks. That warning is informational — the PDF is
produced either way.

If fidelity matters, compare against the PowerPoint driver on the same file rather than
eyeballing it.

### Platform setup

=== "Windows"

    - **PowerPoint driver** — needs Microsoft PowerPoint installed and `powershell` or `pwsh` on
      `PATH`; it drives PowerPoint over COM.
    - **LibreOffice driver** — if `soffice` is not recognised, add
      `C:\Program Files\LibreOffice\program` to `PATH`.

=== "Linux"

    ```bash
    sudo apt install libreoffice-impress
    ```

    The PowerPoint driver is unavailable. `auto` falls through to LibreOffice, then native.

=== "macOS"

    ```bash
    brew install --cask libreoffice
    ```

## HTML

```python
from gopptx import HTMLOptions

with Presentation("deck.pptx") as pres:
    pres.export_html("deck.html", HTMLOptions(
        embed_images=True,          # inline images as base64 data URIs
        include_navigation=True,    # prev/next buttons
        include_slide_numbers=True, # slide number overlay
        base_url="",                # asset prefix when embed_images is False
    ))
```

All four options default to the values shown. `embed_images=True` produces a single
self-contained file — convenient to email, larger on disk. Set it to `False` plus a `base_url`
when the images should be served alongside the page.

HTML export includes speaker notes, charts rendered as SVG, connectors and SmartArt.

## Flat XML

```python
pres.save_flat_xml("deck.xml")
```

A single-file XML serialisation of the package — useful for diffing two decks, or for feeding a
deck to an XML toolchain without unzipping it.

## Grayscale

```python
pres.convert_to_grayscale()                       # whole deck
pres.convert_to_grayscale(slides=[0, 2])          # selected slides
```

Scope can be narrowed to specific shapes, text runs or placeholders via the `shapes=`, `text=`
and `placeholders=` arguments.

## Compression

Package-level size reduction lives in Go and the CLI:

```bash
# Report where the bytes are, without changing anything
go run ./cmd/pptcli compress -in deck.pptx -analyze

# Compress in place or to a new file
go run ./cmd/pptcli compress -in deck.pptx -out small.pptx -level maximum
go run ./cmd/pptcli compress -in deck.pptx -target-size 5000000
```

| Flag | Default | Effect |
|---|---|---|
| `-level` | `balanced` | `light`, `balanced` or `maximum` |
| `-target-size` | `0` | Best-effort maximum output size in bytes |
| `-analyze` | `false` | Report the size breakdown instead of compressing |
| `-format` | `text` | `text` or `json` |
| `-remove-unused-media` | `true` | Drop media parts no relationship points at |
| `-optimize-xml` | `true` | Strip insignificant whitespace from XML parts |
| `-remove-properties` | `false` | Drop custom properties and the package thumbnail |
| `-remove-notes` | `false` | Drop notes slides |
| `-remove-comments` | `false` | Drop comments and comment authors |

`-out` defaults to overwriting the input, so pass it explicitly unless that is what you want.

## Go

```go
import "github.com/djinn-soul/gopptx/pkg/pptx/export"
```

The same three exporters are available directly — `pdf.go` (with `pdf_native.go` and
`pdf_libreoffice.go` behind it) and `html.go`. The Python methods are thin wrappers over these.

## Verifying output

Rendering agreement is the only real test.

- Rasterise the PDF with `pypdfium2` rather than reopening it in LibreOffice — LibreOffice drops
  clips and substitutes fonts, so it will disagree with itself.
- Grade the native renderer by pixel-diffing its pages against the PowerPoint COM driver's on
  the same deck.
- For the `.pptx` itself, open it in real PowerPoint. Structural validation passes on files
  PowerPoint refuses.

## Known gaps

| Missing | Notes |
|---|---|
| PPTX → Markdown export | Not implemented. |
| PPTX → image (PNG/JPEG) export from the library | Screenshots are produced by driving PowerPoint externally — see [the screenshot pipeline](pptx-png-export.md). |
| PDF import | Export only; gopptx does not read PDFs. |

See the [feature matrix](../reference/feature-matrix.md) for the full list.
