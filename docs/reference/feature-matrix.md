# Feature matrix

An honest account of what gopptx does, what it does not, and how it compares to the library
people usually weigh it against: **python-pptx**, the Python default.

Verification method: public API and type names, plus the OOXML each feature actually emits —
not documentation claims from either side.

## At a glance

| | gopptx | python-pptx |
|---|---|---|
| Language | Go + Python | Python |
| Create decks | ✅ | ✅ |
| Read and edit existing decks | ✅ full shape tree | ✅ |
| Chart types | **24**, with embedded data workbook | 15+ |
| Shape presets | **~200** — 206 Go constants over 203 distinct OOXML `prst` values; 202 members in the Python `ShapeType` enum | full MSO enum |
| Tables incl. merge/split, custom styles | ✅ | partial |
| SmartArt | ✅ author, edit, restyle | ❌ |
| Mermaid diagrams | ✅ native | ❌ |
| Markdown / HTML / URL import | ✅ | ❌ |
| PDF export | ✅ native Go + LibreOffice + PowerPoint COM | ❌ |
| HTML export | ✅ incl. notes | ❌ |
| Markdown export | ❌ | ❌ |
| PDF import | ❌ | ❌ |
| Image (PNG/JPEG) export | ❌ from the library | ❌ |
| MCP server | ❌ | ❌ |
| Low-level OPC package API | ✅ `pkg/pptx/opc` | ✅ |
| Package validation and repair | ✅ | ❌ |
| Package compression / size analysis | ✅ | ❌ |
| Cross-language bridge | ✅ JSON over a C ABI | n/a |
| Batch execution | ✅ | n/a |

## By area

| Area | State |
|---|---|
| Slides | Add, insert, remove, move, duplicate, hide, copy between decks, merge decks, rebind layouts |
| Text | Runs with full character formatting, paragraph properties, autofit, anchors, inline formatting outside the Markdown importer, RTL and complex-script handling incl. Persian, Urdu, Pashto and Kurdish |
| Shapes | ~200 presets, freeforms, groups, connectors with auto-routing, gradient/pattern/picture fills, shadows, glow, soft edges, rotation, flips, z-order |
| Tables | Rows, columns, merge, split, sizing, borders, banding, built-in and custom styles; overflow tables beyond the first are written |
| Charts | 24 types, combo, **multiple charts per slide**, embedded Excel workbook so data stays editable, axes, legends, labels, trendlines, error bars |
| Images | Local, bytes, base64, URL, natural-size placement, crops, effects (inner shadow, glow, soft edges, blur, reflection) |
| Media | Video, audio, online video and OLE on both the editor and generator paths |
| Diagrams | SmartArt generated from the data model — every layout PowerPoint offers, node-level edits, quick styles, colour styles — and native Mermaid |
| Layouts | Every generated layout carries a `p:sldLayout/@type`; the vertical-text layouts ship |
| Themes | Built-in presets plus custom colour and font schemes; a colour type with arithmetic |
| Units | Absolute lengths plus relative dimensions (`PercentOf`, `Ratio`) |
| Slide show | Complete `p:showPr` settings, plus custom shows |
| Package | `presProps`, `viewProps`, `tableStyles`, a default text style and the chart-space defaults Office writes; handout master declared in `presentation.xml` |
| Document properties | Core and app properties written from what the caller actually set |
| Validation | Package rules beyond the shallow part-existence check, plus repair |
| Export | PDF (three drivers), HTML, flat XML, grayscale, compression and size analysis; speaker notes in both HTML and PDF |

## What gopptx does not do

Four things, all verified absent in the current tree rather than inferred:

| Missing | Detail | Workaround |
|---|---|---|
| **PPTX → Markdown export** | No exporter. | Read the deck and format the text yourself; `list_shapes` and `get_shape_runs` give you the content. |
| **PDF → PPTX import** | Export only; there is no PDF reader. | None in-library. |
| **Image (PNG/JPEG) export and thumbnails** | The library renders to PDF and HTML but not to raster images. | Drive PowerPoint or LibreOffice externally — see [the screenshot pipeline](../guides/pptx-png-export.md) — or rasterise the exported PDF with `pypdfium2`. |
| **MCP server** | No binary exposing presentation tools over the Model Context Protocol. | The JSON bridge is a stable contract; an MCP wrapper over it would be thin. |

## Known rough edges

Not missing features — things that work but will surprise you.

| Behaviour | Detail |
|---|---|
| **Geometry is EMU, with no guard rail** | A shape at `(40, 120, 600, 220)` is smaller than a pixel. Nothing warns you; the file saves and opens. Use `Inches()`. |
| **Bare chart-type strings are deprecated** | `add_chart(0, "bar", …)` warns and will be rejected later. Pass a `ChartType` member. |
| **Result keys come in two spellings** | Responses carry the original Go key *and* a snake_case alias — `shape["ID"]` and `shape["id"]` both work. Handlers whose Go structs are JSON-tagged only ever had the snake_case form. |

### Fixed since the last revision

| Was | Now |
|---|---|
| `THEME_*` constants rejected by `apply_theme` | `apply_theme` and `set_global_theme_preset` share one vocabulary — the gopptx theme names *and* the Office preset names, case- and separator-insensitive, resolved by `styling.ResolveTheme`. |
| `validate()` returned `None` on an opened deck, and reported spurious `ORPHAN_SLIDE` findings on an edited one | Always returns a list, and validates the package **as it would be saved**, so slides added since opening are no longer reported as orphans. |
| `ChartType`, `PDFOptions`, `HTMLOptions` needed deep module paths | Re-exported at the top level, along with `Inches`, `Point`, `Emu` and `SlideLayoutType`. |
| Read results mixed key casing with no aliasing | Aliased centrally, so every response is readable in snake_case. |

## Where gopptx is ahead

- **Two languages, one engine.** Python gets Go's performance and Go's exact output bytes.
- **Batching.** Write-heavy Python workloads cross the FFI boundary once per batch.
- **SmartArt and Mermaid.** python-pptx does not author SmartArt at all.
- **PDF fidelity.** Three drivers, including a pure-Go renderer that needs no external process,
  graded against PowerPoint's own output.
- **Editing depth.** The editor exposes the full shape tree — 179 operations — over a stable
  cross-language contract.

## Keeping this page honest

Entries here are verified against the tree: element-level, attribute-level, OPC-identity and
numeric-defaults checks on the OOXML that gopptx actually writes.

When you close a gap, update this table. When you find a new one, record it here rather than
leaving the comparison flattering.
