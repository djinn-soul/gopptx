# Examples map

The repository carries **96 runnable examples** under `examples/`, one directory per topic.
Most contain both a Go program and a Python equivalent, so you can compare the two APIs on the
same task.

```
examples/<topic>/
├── *.go        # Go version
└── main.py     # Python version
```

## Running one

```bash
# Go
go run ./examples/01-basic-pptx-generation

# Python (build the shared library first — see Installation)
python examples/01-basic-pptx-generation/main.py
```

Generated files land in `examples/output/`; rendered screenshots, where a topic has them, in
`examples/output/png/`.

## Where to start

| If you want to… | Read |
|---|---|
| See the minimum viable deck | `01-basic-pptx-generation` |
| Format text | `04-text-formatting`, `06-text-enhancements`, `70-text-api` |
| Draw shapes | `12-shapes`, `13-connectors`, `65-shapes-api` |
| Build tables | `07-tables`, `14-table-api-improvements`, `67-table-api` |
| Add charts | `09-charts`, `16-chart-type-enum`, `68-chart-data-api`, `69-chart-api` |
| Place images | `10-images`, `11-image-advanced-sources`, `73-image-api` |
| Edit an existing deck | `19-read-modify-existing`, `41-deep-shape-editing` |
| Work with layouts and masters | `02-slide-layouts`, `36-slide-master`, `54-theme-master`, `17-custom-slide-layout-composition` |
| Apply themes | `17-themes`, `43-advanced-theme-management`, `18-theme-aware-presentation` |
| Author SmartArt | `24-smartart`, `61-smartart-edit-hidden`, `62-smartart-edit-variations` |
| Render Mermaid | `32-mermaid`, `32-mermaid-diagrams` |
| Import Markdown or the web | `03-markdown-to-pptx`, `34-urlfetch`, `34-web2ppt` |
| Export | `21-export-html-pdf`, `25-export-html`, `81-native-pdf-showcase` |
| Add notes, comments, sections | `22-speaker-notes`, `45-commenting-api`, `44-section-management` |
| Animate | `14-transitions`, `28-animations`, `47-morph-transitions` |
| Protect or sign a deck | `46-presentation-protection`, `20-validation-repair` |
| Metadata and reskinning | `40-metadata-writer`, `43-presentation-props-editor` |
| Use the bridge directly | `57-bridge-command-api`, `58-bridge-performance` |
| See everything at once | `58-gopptx-rich-slide` |

## Grouped by number

| Range | Theme |
|---|---|
| `01`–`06` | Generation basics, layouts, Markdown, text |
| `07`–`11` | Tables, charts, images |
| `12`–`18` | Shapes, connectors, transitions, CLI, templates, themes, layout helpers |
| `19`–`23` | Editing, validation, export, notes, media |
| `24`–`33` | SmartArt, HTML export, VBA, custom XML, animations, handout and notes masters, fonts, hyperlinks, Mermaid |
| `34`–`42` | URL fetch, layout helpers, masters, duplication, editor image/chart/notes support, deep editing, asset merge |
| `43`–`52` | Theme management, sections, comments, protection, morph, accessibility, hyperlinks, interactivity, localisation, legacy interop |
| `53`–`59` | Slide properties, backgrounds, actions, bridge API, placeholder overrides, performance |
| `60`–`79` | The API-surface series — one directory per typed API area |
| `80`–`81` | Chart-type export, native PDF showcase |

## Annotated recipes

The [showcase](showcase/usages/index.md) presents 30 of these as narrated pages with
screenshots and downloadable `.pptx` output, graded simple → intermediate → complex:

- [Simple](showcase/usages/simple/index.md) — hello world, text frames, images, backgrounds,
  shapes, transitions
- [Intermediate](showcase/usages/intermediate/index.md) — tables, notes, headers and footers,
  groups, charts, custom layouts, themes, dynamic reports, sections, animation
- [Complex](showcase/usages/complex/index.md) — Markdown + Mermaid decks, URL-to-slides, code
  rendering, cloning, template injection, Jinja2, export pipelines, storyboards, SmartArt,
  protection

## Docs gallery cross-links

- [Basic Generation](showcase/basic-generation.md)
- [Basic Usage](showcase/basic-usage.md)
- [Rich Slide](showcase/rich-slide.md)
- [Text Styling](showcase/text-styling.md)
- [Radar Chart](showcase/chart-radar.md)

## A caveat on older examples

Examples are compiled and run in CI, so they work — but a few predate the current preferred
idioms. Where an example and this documentation disagree on style:

- prefer `SlideBuilder` over chained `SlideContent` in Go
- prefer `ChartType` members over bare chart-type strings in Python
- wrap every coordinate in `Inches()` / `Point()` / `Emu()`
