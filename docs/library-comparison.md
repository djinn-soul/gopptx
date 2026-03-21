# PPTX Library Comparison: gopptx vs ppt-rs vs PptxGenJS

> **Report Date**: March 2026
> **Versions Analyzed**: gopptx (latest main), ppt-rs v0.2.7, PptxGenJS v4.0.1
> **Scope**: Feature inventory, gap analysis, use-case guidance, and roadmap suggestions

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Library Profiles](#2-library-profiles)
3. [Feature Comparison Matrix](#3-feature-comparison-matrix)
4. [Similarity Analysis](#4-similarity-analysis)
5. [Differentiator Analysis](#5-differentiator-analysis)
6. [Gap Analysis: gopptx](#6-gap-analysis-gopptx)
7. [Gap Analysis: ppt-rs](#7-gap-analysis-ppt-rs)
8. [Gap Analysis: PptxGenJS](#8-gap-analysis-pptxgenjs)
9. [Use Case Recommendations](#9-use-case-recommendations)
10. [Feature Category Deep-Dives](#10-feature-category-deep-dives)
    - [10.1 Charts](#101-charts)
    - [10.2 Text & Formatting](#102-text--formatting)
    - [10.3 Tables](#103-tables)
    - [10.4 Images & Media](#104-images--media)
    - [10.5 Themes](#105-themes)
    - [10.6 Export & Output](#106-export--output)
    - [10.7 Animations & Transitions](#107-animations--transitions)
11. [API Style Comparison](#11-api-style-comparison)
12. [Conclusion & Roadmap Suggestions](#12-conclusion--roadmap-suggestions)

---

## 1. Executive Summary

The PPTX ecosystem has three meaningfully different libraries serving distinct developer audiences. **gopptx** is a Go backend library with Python bindings, designed for server-side automation and programmatic modification of existing presentations. **ppt-rs** is a Rust library with a focus on compile-time safety and ergonomic builder patterns, targeting systems programmers who need to generate and manipulate PPTX files. **PptxGenJS** is a JavaScript/TypeScript library built for client-side and server-side web environments, offering zero-dependency PPTX creation in browsers and Node.js with no support for opening or editing existing files.

The three libraries occupy non-overlapping niches in one important axis: **read/write vs create-only**. Both gopptx and ppt-rs can open, modify, and save existing PPTX files, enabling document automation workflows, template filling, and round-trip editing. PptxGenJS is strictly a creation library — it has no ability to open existing files, reorder slides, or apply transitions and animations. This single constraint defines its suitability almost entirely to greenfield presentation generation in web environments.

In terms of raw feature breadth, **gopptx leads significantly**. Its 145+ JSON-bridged operations cover areas that neither ppt-rs nor PptxGenJS address: Jinja2 template rendering, Mermaid diagram import, animated transitions, slide sections, header/footer management, document protection, VBA/macros, custom XML, OLE objects, comments, audio/video embedding, SmartArt, and a Python binding layer with type-safe enums. Its architecture — a Go core exposed via a C bridge to Python — enables the library to be used from both Go and Python codebases, effectively serving two language ecosystems simultaneously with a single implementation.

**ppt-rs** is the most technically ambitious from a Rust-language perspective, offering features that neither competitor has: embedded fonts, digital signatures, 3D model embedding, sub-bullets with hierarchical indentation, code blocks with syntax highlighting, and rich image effects (shadow, reflection, glow, soft edges, blur). Its fluent builder API and macro system (`pptx!()`, `shape!()`) provide an ergonomic authoring experience. Despite being at v0.2.7, ppt-rs has impressive feature coverage for a younger library. The main gaps versus gopptx are in the editing/automation category: no find-and-replace, no per-shape text state introspection, no placeholder API, and no animations.

**PptxGenJS** excels in its target niche: it is the go-to choice for JavaScript developers who need to generate PPTX files without any server-side runtime or native dependencies. Its browser-native output modes (Blob, base64, ArrayBuffer), HTML table import, and YouTube embed support have no equivalent in the other two libraries. The zero-dependency design makes it uniquely suited to serverless and edge-computing environments. However, any workflow requiring existing file modification, slide transitions, animations, SmartArt, or rich theme control will quickly hit its hard boundaries.

---

## 2. Library Profiles

### Quick Reference Table

| Attribute | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| **Primary Language** | Go | Rust | JavaScript / TypeScript |
| **Binding Languages** | Go, Python | Rust (no bindings) | JS, TS, Node.js, browsers |
| **Version Analyzed** | latest main | v0.2.7 | v4.0.1 |
| **Maturity** | Active / Production | Active / Beta | Stable / Widely Used |
| **License** | (project default) | (project default) | MIT |
| **Paradigm** | JSON command bridge + OOP API | Fluent builder + macros | Fluent builder (OOP) |
| **Read Existing PPTX** | ✅ Yes | ✅ Yes | ❌ No |
| **Modify Existing PPTX** | ✅ Yes | ✅ Yes | ❌ No |
| **Create New PPTX** | ✅ Yes | ✅ Yes | ✅ Yes |
| **Server-Side Use** | ✅ Yes | ✅ Yes | ✅ Yes (Node.js) |
| **Client-Side (Browser)** | ❌ No | ❌ No | ✅ Yes |
| **Zero Dependencies** | ❌ No | ❌ No | ✅ Yes |
| **Python API** | ✅ Full bindings | ❌ None | ❌ None |
| **Primary Use Case** | Automation / server-side editing | Systems / embedded generation | Web / serverless generation |
| **Batch Operations** | ✅ Yes (BatchExecute) | ❌ No | ❌ No |
| **CLI Tool** | ✅ pptcli | ✅ pptcli md2ppt | ❌ No (library only) |

### Language & Runtime Requirements

| Library | Runtime | Build Dependencies | Deployment Model |
|---|---|---|---|
| gopptx | Go 1.21+, Python 3.9+ (optional) | CGO for Python bridge | Single Go binary or Python package |
| ppt-rs | Rust stable | Cargo + optional system libs for PDF export | Compiled Rust binary or library crate |
| PptxGenJS | Node.js 14+ or modern browser | npm/yarn, zero runtime deps | npm package, CDN, or bundled |

---

## 3. Feature Comparison Matrix

The table below rates each library against every major feature category using the following scale:

| Symbol | Meaning |
|---|---|
| ✅ Full | Complete, production-ready support |
| ⚡ Partial | Supported but with notable limitations |
| ❌ None | Not supported |
| 🔧 External | Requires external tool or plugin |

### 3.1 Core Presentation Management

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Create new presentation | ✅ Full | ✅ Full | ✅ Full |
| Open existing PPTX | ✅ Full | ✅ Full | ❌ None |
| Save / Save As | ✅ Full | ✅ Full | ✅ Full |
| Output to bytes/stream | ✅ Full | ✅ Full | ✅ Full |
| Output to base64 | ✅ Full | ✅ Full | ✅ Full |
| Metadata (title, author, etc.) | ✅ Full (12 fields) | ⚡ Partial (title) | ✅ Full (6 fields) |
| Slide size / layout dimensions | ✅ Full | ⚡ Partial | ✅ Full |
| Preset layout dimensions | ✅ Full | ✅ Full | ✅ Full (4 presets) |
| Custom layout dimensions | ✅ Full | ✅ Full | ✅ Full |
| RTL mode | ❌ None | ✅ Full | ✅ Full |
| Merge presentations | ✅ Full | ✅ Full | ❌ None |
| Validate & repair | ✅ Full | ✅ Full (8 categories) | ❌ None |
| Batch operations | ✅ Full | ❌ None | ❌ None |

### 3.2 Slide Operations

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Add slide | ✅ Full | ✅ Full | ✅ Full |
| Remove slide | ✅ Full | ✅ Full | ❌ None |
| Move / reorder slide | ✅ Full | ❌ None | ❌ None |
| Duplicate slide | ✅ Full | ❌ None | ❌ None |
| Slide count | ✅ Full | ✅ Full | ❌ None |
| List slides | ✅ Full | ✅ Full | ❌ None |
| Hidden slide | ❌ None | ❌ None | ✅ Full |
| Slide background (solid) | ✅ Full | ✅ Full | ✅ Full |
| Slide background (gradient) | ✅ Full | ✅ Full | ❌ None |
| Slide background (image) | ✅ Full | ✅ Full | ✅ Full |
| Slide background (theme) | ✅ Full | ✅ Full | ⚡ Partial |
| Slide number display | ✅ Full | ✅ Full | ✅ Full |
| Sections | ✅ Full | ✅ Full | ✅ Full |
| Update slide (title, layout) | ✅ Full | ✅ Full | ❌ None |
| Rebind slide layout | ✅ Full | ❌ None | ❌ None |

### 3.3 Shapes

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| 100+ shape types | ✅ Full | ✅ Full | ✅ Full (~200) |
| Add shape | ✅ Full | ✅ Full | ✅ Full |
| Remove shape | ✅ Full | ❌ None | ❌ None |
| Update / reposition shape | ✅ Full | ⚡ Partial | ❌ None |
| Move shape to front/back | ✅ Full | ❌ None | ❌ None |
| Move shape to Z-index | ✅ Full | ❌ None | ❌ None |
| Group / ungroup shapes | ✅ Full | ❌ None | ❌ None |
| Shape fill (solid) | ✅ Full | ✅ Full | ✅ Full |
| Shape fill (gradient) | ✅ Full | ✅ Full | ❌ None |
| Shape line / border | ✅ Full | ✅ Full | ✅ Full |
| Hyperlinks on shapes | ✅ Full | ✅ Full | ✅ Full |
| Shape rotation | ✅ Full | ✅ Full | ✅ Full |
| Shape flip (H/V) | ❌ None | ❌ None | ✅ Full |
| Shadow on shapes | ❌ None | ✅ Full | ✅ Full |
| Freeform / custom geometry | ✅ Full | ⚡ Partial | ✅ Full (points array) |
| Connectors | ✅ Full | ❌ None | ❌ None |
| Layout helpers (grid, center) | ❌ None | ✅ Full | ❌ None |
| Search shapes | ✅ Full | ❌ None | ❌ None |

### 3.4 Text & Formatting

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Bold / italic / underline | ✅ Full | ✅ Full | ✅ Full |
| Strikethrough | ❌ None | ✅ Full | ✅ Full (single/double) |
| Superscript / subscript | ❌ None | ✅ Full | ✅ Full |
| Font family | ✅ Full | ✅ Full | ✅ Full |
| Font size | ✅ Full | ✅ Full | ✅ Full |
| Font color | ✅ Full | ✅ Full | ✅ Full |
| Font highlight color | ❌ None | ✅ Full | ✅ Full |
| Text transparency | ❌ None | ❌ None | ✅ Full |
| Text alignment (H) | ✅ Full | ✅ Full | ✅ Full |
| Text alignment (V / anchor) | ❌ None | ✅ Full | ✅ Full |
| Line spacing | ❌ None | ❌ None | ✅ Full (points/multiple) |
| Character spacing | ❌ None | ❌ None | ✅ Full |
| Paragraph space before/after | ❌ None | ❌ None | ✅ Full |
| Bullets (standard) | ❌ None | ✅ Full | ✅ Full |
| Bullets (numbered) | ❌ None | ✅ Full | ✅ Full |
| Bullets (custom / Unicode) | ❌ None | ✅ Full | ✅ Full |
| Hierarchical sub-bullets | ❌ None | ✅ Full | ⚡ Partial (indent levels) |
| RTL text | ❌ None | ✅ Full | ✅ Full |
| Vertical text | ❌ None | ❌ None | ✅ Full (7 modes) |
| Auto-fit (shrink/resize) | ❌ None | ❌ None | ✅ Full |
| Text wrap control | ❌ None | ❌ None | ✅ Full |
| Tab stops | ❌ None | ❌ None | ✅ Full |
| Code block with syntax highlighting | ❌ None | ✅ Full | ❌ None |
| Text shadow / glow / outline | ❌ None | ❌ None | ✅ Full |
| Hyperlinks in text | ✅ Full | ✅ Full | ✅ Full |
| Get/set text runs (editing) | ✅ Full | ❌ None | ❌ None |
| Find and replace | ✅ Full | ❌ None | ❌ None |
| Run-level formatting inspection | ✅ Full | ❌ None | ❌ None |
| Language tag | ❌ None | ⚡ Partial (RTL lang) | ✅ Full |

### 3.5 Tables

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Add table | ✅ Full | ✅ Full | ✅ Full |
| Cell text & formatting | ✅ Full | ✅ Full | ✅ Full |
| Cell fill / background | ✅ Full | ✅ Full | ✅ Full |
| Cell borders | ✅ Full | ✅ Full | ✅ Full |
| Merge cells | ✅ Full | ✅ Full | ✅ Full |
| Split cells | ✅ Full | ⚡ Partial | ❌ None |
| Row height control | ✅ Full | ✅ Full | ✅ Full |
| Column width control | ✅ Full | ✅ Full | ✅ Full |
| Table style (named) | ✅ Full | ❌ None | ❌ None |
| Define custom table style | ✅ Full | ❌ None | ❌ None |
| Colspan / rowspan | ✅ Full | ✅ Full | ✅ Full |
| Cell alignment (H/V) | ✅ Full | ✅ Full | ✅ Full |
| Bullets in cells | ❌ None | ❌ None | ✅ Full |
| Auto-paginate table | ❌ None | ❌ None | ✅ Full |
| Repeat header row | ❌ None | ❌ None | ✅ Full |
| HTML table → PPTX | ❌ None | ❌ None | 🔧 Browser only |
| Get existing table state | ✅ Full | ❌ None | ❌ None |
| TableMergeMap (advanced) | ❌ None | ✅ Full | ❌ None |

### 3.6 Charts

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Bar / Column chart | ✅ Full | ✅ Full | ✅ Full |
| Bar horizontal | ❌ None | ✅ Full | ✅ Full |
| Bar stacked | ❌ None | ✅ Full | ✅ Full |
| Bar stacked 100% | ❌ None | ✅ Full | ✅ Full |
| Bar 3D | ❌ None | ❌ None | ✅ Full |
| Line chart | ✅ Full | ✅ Full | ✅ Full |
| Line with markers | ❌ None | ✅ Full | ✅ Full |
| Line stacked | ❌ None | ✅ Full | ❌ None |
| Area chart | ❌ None | ✅ Full | ✅ Full |
| Area stacked / 100% | ❌ None | ✅ Full | ✅ Full |
| Pie chart | ✅ Full | ✅ Full | ✅ Full |
| Doughnut chart | ❌ None | ✅ Full | ✅ Full |
| Scatter / XY chart | ❌ None | ✅ Full | ✅ Full |
| Scatter smooth/lines | ❌ None | ✅ Full | ✅ Full |
| Bubble chart | ❌ None | ✅ Full | ✅ Full |
| Bubble 3D | ❌ None | ❌ None | ✅ Full |
| Radar chart | ❌ None | ✅ Full | ✅ Full |
| Radar filled | ❌ None | ✅ Full | ❌ None |
| Stock chart (HLC/OHLC) | ❌ None | ✅ Full | ❌ None |
| Combo chart | ❌ None | ✅ Full | ✅ Full |
| Chart title | ✅ Full | ✅ Full | ✅ Full |
| Chart legend | ✅ Full | ⚡ Partial | ✅ Full |
| Axis labels / formatting | ✅ Full | ⚡ Partial | ✅ Full |
| Value / data labels | ✅ Full | ⚡ Partial | ✅ Full |
| Per-series colors | ✅ Full | ✅ Full | ✅ Full |
| Logarithmic scale | ❌ None | ❌ None | ✅ Full |
| Data table in chart | ❌ None | ❌ None | ✅ Full |
| Update existing chart data | ✅ Full | ❌ None | ❌ None |
| Get chart state | ✅ Full | ❌ None | ❌ None |
| List charts on slide | ✅ Full | ❌ None | ❌ None |

### 3.7 Images

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| From file path | ✅ Full | ✅ Full | ✅ Full |
| From base64 | ✅ Full | ✅ Full | ✅ Full |
| From URL | ❌ None | ❌ None | ✅ Full |
| From bytes | ❌ None | ✅ Full | ⚡ Partial |
| JPEG / PNG / GIF | ✅ Full | ✅ Full | ✅ Full |
| SVG | ❌ None | ❌ None | ✅ Full |
| Auto-detect format | ✅ Full | ✅ Full | ✅ Full |
| SHA256 deduplication | ✅ Full | ❌ None | ❌ None |
| Crop | ❌ None | ✅ Full (percentage) | ✅ Full |
| Rotation | ❌ None | ✅ Full | ✅ Full |
| Opacity / transparency | ❌ None | ❌ None | ✅ Full |
| Circle / rounding crop | ❌ None | ❌ None | ✅ Full |
| Shadow effect | ❌ None | ✅ Full | ✅ Full |
| Reflection effect | ❌ None | ✅ Full | ❌ None |
| Glow effect | ❌ None | ✅ Full | ✅ Full |
| Soft edges | ❌ None | ✅ Full | ❌ None |
| Blur effect | ❌ None | ✅ Full | ❌ None |
| Inner shadow | ❌ None | ✅ Full | ❌ None |
| Sizing modes (contain/cover) | ❌ None | ❌ None | ✅ Full |
| Hyperlink on image | ✅ Full | ✅ Full | ✅ Full |
| Alt text | ❌ None | ❌ None | ✅ Full |
| Swap existing image | ✅ Full | ❌ None | ❌ None |
| List images on slide | ✅ Full | ❌ None | ❌ None |
| Get image metadata | ✅ Full | ❌ None | ❌ None |

### 3.8 Media (Video, Audio, 3D)

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Video embedding (MP4 etc.) | ✅ Full | ✅ Full | ✅ Full |
| Video auto-play / loop / muted | ✅ Full | ❌ None | ❌ None |
| Audio embedding (MP3/WAV) | ✅ Full | ✅ Full | ✅ Full |
| Audio with icon | ✅ Full | ❌ None | ❌ None |
| YouTube embed | ❌ None | ❌ None | ✅ Full (M365) |
| 3D model embedding (GLB/GLTF) | ❌ None | ✅ Full | ❌ None |
| OLE object embedding | ✅ Full | ❌ None | ❌ None |

### 3.9 Themes

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Built-in themes (named) | ✅ Full (7) | ✅ Full (7) | ❌ None |
| Apply theme to presentation | ✅ Full | ✅ Full | ❌ None |
| Custom color scheme | ✅ Full | ✅ Full | ⚡ Partial (SchemeColor enum) |
| Custom font scheme | ✅ Full | ✅ Full | ⚡ Partial (heading/body faces) |
| Color manipulation (lighter/darker) | ❌ None | ✅ Full | ❌ None |
| Theme inventory / list | ✅ Full | ❌ None | ❌ None |
| .thmx file support | ❌ None | ❌ None | ❌ None |
| Embedded fonts | ❌ None | ✅ Full | ❌ None |
| Material Design color palette | ❌ None | ✅ Full | ❌ None |
| IBM Carbon color palette | ❌ None | ✅ Full | ❌ None |

### 3.10 Slide Layouts & Masters

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Slide layout types (enum) | ✅ Full (6 types) | ✅ Full (12 types) | ❌ No native API |
| List slide layouts | ✅ Full | ✅ Full | ❌ None |
| List slide masters | ✅ Full | ✅ Full | ❌ None |
| Add slide master | ✅ Full | ✅ Full | ✅ Full (defineSlideMaster) |
| Remove slide master | ✅ Full | ❌ None | ❌ None |
| Add slide layout | ✅ Full | ✅ Full | ✅ Full (as master object) |
| Remove slide layout | ✅ Full | ❌ None | ❌ None |
| Clone layout+master family | ✅ Full | ❌ None | ❌ None |
| Get layout shapes | ✅ Full | ❌ None | ❌ None |
| Get layout placeholders | ✅ Full | ❌ None | ❌ None |
| Get master shapes | ✅ Full | ❌ None | ❌ None |
| Get master placeholders | ✅ Full | ❌ None | ❌ None |
| Handout master | ✅ Full | ✅ Full | ❌ None |
| HandoutLayout type | ❌ None | ✅ Full | ❌ None |

### 3.11 Placeholders

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| List placeholders | ✅ Full | ❌ None | ❌ None |
| Set placeholder content (text) | ✅ Full | ❌ None | ✅ Full (via master def) |
| Set placeholder content (image) | ✅ Full | ❌ None | ✅ Full |
| Set placeholder content (chart) | ✅ Full | ❌ None | ✅ Full |
| Set placeholder content (table) | ✅ Full | ❌ None | ✅ Full |
| Standard placeholder types | ✅ Full (8 types) | ❌ None | ✅ Full (5 types) |
| Inspect existing placeholders | ✅ Full | ❌ None | ❌ None |
| Notes placeholders | ✅ Full | ❌ None | ❌ None |

### 3.12 Notes

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Get slide notes | ✅ Full | ✅ Full | ❌ None |
| Set slide notes | ✅ Full | ✅ Full | ✅ Full (plain text) |
| Rich-text notes | ✅ Full | ✅ Full (Markdown) | ❌ None |
| Notes slide exists check | ✅ Full | ❌ None | ❌ None |
| List notes shapes | ✅ Full | ❌ None | ❌ None |
| Set notes shape text | ✅ Full | ❌ None | ❌ None |
| Set notes shape props | ✅ Full | ❌ None | ❌ None |
| Update notes master | ✅ Full | ❌ None | ❌ None |

### 3.13 Animations & Transitions

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Slide transitions | ✅ Full (9 types) | ✅ Full (9 types) | ❌ None |
| Shape animations | ✅ Full | ❌ None | ❌ None |
| Animation triggers (onClick/with/after) | ✅ Full | ❌ None | ❌ None |
| Animation types (fade, fly, etc.) | ✅ Full | ❌ None | ❌ None |

### 3.14 SmartArt

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Add SmartArt | ✅ Full (25 layouts) | ✅ Full (25 layouts) | ❌ None |
| Update SmartArt | ✅ Full | ❌ None | ❌ None |
| Custom layout URIs | ✅ Full | ⚡ Partial | ❌ None |

### 3.15 Special Input/Import Features

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Markdown → slides | ✅ Full | ✅ Full (CLI) | ❌ None |
| Mermaid diagram → PPTX shapes | ✅ Full (12 types) | ✅ Full (12 types) | ❌ None |
| Jinja2 template rendering | ✅ Full | ❌ None | ❌ None |
| HTML table → PPTX | ❌ None | ❌ None | 🔧 Browser only |
| Web URL → PPTX | ❌ None | ✅ Full (feature flag) | ❌ None |
| HTML → PPTX | ❌ None | ✅ Full (feature flag) | ❌ None |
| PDF → PPTX (import) | ❌ None | ✅ Full | ❌ None |
| Template builders | ✅ Full (5 templates) | ❌ None | ❌ None |

### 3.16 Document Features

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Comments (get/add/remove) | ✅ Full | ❌ None | ❌ None |
| Headers / footers | ✅ Full | ❌ None | ❌ None |
| Document protection / password | ✅ Full | ❌ None | ❌ None |
| Mark as final | ✅ Full | ❌ None | ❌ None |
| Digital signature detection | ✅ Full | ✅ Full | ❌ None |
| Digital signature creation | ❌ None | ✅ Full | ❌ None |
| VBA / macros | ✅ Full | ✅ Full | ❌ None |
| Custom XML parts | ✅ Full | ✅ Full | ❌ None |
| Find and replace | ✅ Full | ❌ None | ❌ None |

### 3.17 Export

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| PPTX file output | ✅ Full | ✅ Full | ✅ Full |
| PDF export (native) | ✅ Full | 🔧 LibreOffice | ❌ None |
| PDF export (LibreOffice) | ✅ Full | 🔧 LibreOffice | ❌ None |
| PDF export (PowerPoint) | ✅ Full | ❌ None | ❌ None |
| HTML export | ✅ Full | ✅ Full | ❌ None |
| PNG per slide | ❌ None | 🔧 pdftoppm | ❌ None |
| Browser download (Blob/base64) | ❌ None | ❌ None | ✅ Full |
| ArrayBuffer / Uint8Array output | ❌ None | ❌ None | ✅ Full |
| HTTP stream output | ❌ None | ❌ None | ✅ Full |

---

## 4. Similarity Analysis

All three libraries share a meaningful common ground across several core feature categories. Understanding what is universal helps clarify which requirements are baseline (solved by any of the three) versus which require careful library selection.

### 4.1 Universal Core: What All Three Libraries Do Well

**PPTX Creation from Scratch**
Every library provides solid support for creating new PPTX presentations, adding slides, inserting shapes and images, and writing the file to disk or a byte buffer. This is the lowest common denominator and works reliably across all three.

**Basic Shape Support**
All three support 100+ PowerPoint shape types with fill, line/border, and hyperlink support. The shape type enum coverage is functionally equivalent for common presentation shapes (rectangles, circles, arrows, callouts, flowchart elements).

**Core Text Formatting**
Bold, italic, underline, font size, font family, and font color work across all three libraries. Horizontal text alignment is universal. Basic rich text with mixed inline styles (combining font weight, color, and size within a single paragraph) is supported by all three.

**Image Embedding**
All three support PNG, JPEG, and GIF from file path or base64. Positioning and sizing are handled consistently. All three support hyperlinks on images.

**Chart Basics**
Bar/column, line, and pie charts with titles and per-series colors are available in all three libraries. These cover the overwhelming majority of presentation chart use cases.

**Table Construction**
Adding tables with cell text formatting, fill colors, borders, and merged cells is fully supported across all three. Row height and column width control is universal.

**Sections**
All three include a section/grouping model for organizing slides.

**Output Formats**
All three can produce PPTX files as binary output. Both gopptx and ppt-rs offer HTML export. PptxGenJS offers the broadest set of JavaScript-native output types (Blob, base64, Buffer, ArrayBuffer).

**7 Built-in Themes**
Both gopptx and ppt-rs share the exact same set of 7 named built-in themes (Corporate, Modern, Vibrant, Dark, Nature, Tech, Carbon), suggesting shared design intent or lineage. PptxGenJS does not have equivalent built-in themes.

**SmartArt with 25 Layouts**
Both gopptx and ppt-rs support the same 25 SmartArt layouts. This is a non-trivial feature that PptxGenJS entirely lacks.

**Mermaid Diagram Import**
Both gopptx and ppt-rs convert Mermaid diagram syntax into native PPTX shapes and connectors across 12 diagram types. This is a distinctive shared capability absent from PptxGenJS.

**VBA / Custom XML**
Both gopptx and ppt-rs support embedding VBA projects and custom XML parts — enterprise document features absent from PptxGenJS.

**Slide Transitions**
Both gopptx and ppt-rs implement the same 9 transition types (Fade, Cut, Push, Wipe, Split, Reveal, Cover, Zoom, None). PptxGenJS has no transition support.

---

## 5. Differentiator Analysis

### 5.1 What Makes gopptx Unique

**Full Python Bindings with Type-Safe Enums**
The most commercially significant differentiator is the complete Python binding layer. Python is the dominant language for data science, analytics automation, and enterprise scripting. gopptx exposes every Go operation to Python via a C bridge, with type-safe enums (`ShapeType`, `ChartType`, `SlideLayoutType`), a `RunBuilder` for text composition, and a `GopptxError` exception hierarchy. No competing library offers this cross-language availability.

**JSON Command Bridge with BatchExecute**
The 145+ operation JSON bridge is a stable, version-controlled API surface usable from any language that can call a C function or make subprocess calls. The `BatchExecute` capability reduces bridge round-trips dramatically for bulk operations, which is critical for generating large decks (200+ slides) with acceptable latency.

**Deep Editing & Introspection Operations**
gopptx is the most capable editing library. It offers fine-grained text run introspection (`GetShapeRuns`, `SetShapeRuns`, `UpdateShapeRunText`), per-shape text state capture (`GetSlideTextStates`), chart state retrieval (`GetChartState`), image swap by index or relationship ID, placeholder listing and content injection, comment management, and header/footer control. These operations are necessary for template-filling and document automation workflows.

**Jinja2 Template Rendering**
The `RenderTemplate` operation enables Python Jinja2 syntax inside PPTX files, effectively making presentations templatable data documents. This feature has no equivalent in ppt-rs or PptxGenJS.

**5 Built-in Presentation Template Builders**
`BuildStatusTemplate`, `BuildSimpleTemplate`, `BuildProposalTemplate`, `BuildTrainingTemplate`, and `BuildTechnicalTemplate` produce complete, styled multi-slide presentations from structured data. These are high-level content generation helpers absent from other libraries.

**Comprehensive Protection & Password Features**
`SetModifyPassword`, `SetMarkAsFinal`, and `HasDigitalSignature` cover document protection requirements common in enterprise PowerPoint automation.

**PDF Export with Multiple Drivers**
The four PDF export modes (auto, native, LibreOffice, PowerPoint) give operators flexibility depending on what rendering infrastructure is available.

**Shape Z-Order Control**
`MoveShapeToFront`, `MoveShapeToBack`, and `MoveShapeToIndex` give precise control over overlapping shapes — a capability absent from both competitors.

**Connectors**
`AddConnector` creates line connectors between shapes. This is necessary for flowcharts and diagrams and is not available in ppt-rs or PptxGenJS as a first-class feature.

### 5.2 What Makes ppt-rs Unique

**Digital Signatures (Creation)**
ppt-rs can create digital signatures (`DigitalSignature`, `SignerInfo`, `HashAlgorithm`), not just detect them. This is a compliance-critical feature for regulated industries that gopptx and PptxGenJS do not provide.

**Embedded Fonts**
`EmbeddedFontList` and `EmbeddedFontPart` allow font binaries to be embedded inside the PPTX file, ensuring consistent rendering on machines that do not have the fonts installed. Neither gopptx nor PptxGenJS support this.

**3D Model Embedding**
GLB, GLTF, OBJ, FBX, and STL models can be embedded as native 3D objects with camera presets and rotation. This is a modern PowerPoint feature with no equivalent in the other two libraries.

**Rich Image Effects Chain**
Shadow, reflection, glow, soft edges, inner shadow, and blur effects can all be chained on images. The `syntect`-powered code block with Solarized Dark syntax highlighting is unique to ppt-rs.

**Color Manipulation API**
The `.lighter()`, `.darker()`, `.opacity()`, and `.mix()` methods on color values enable programmatic color derivation, which is essential for generating visually coherent decks from data without hardcoding every color.

**Sub-Bullets and Hierarchical Lists**
`sub_bullets` enables proper multi-level bullet hierarchies including LetterLower, LetterUpper, RomanLower, RomanUpper, Custom, and None — the full PowerPoint list model. PptxGenJS supports indent levels but not the full bullet style taxonomy.

**Web-to-PPTX and HTML-to-PPTX**
The `url_to_pptx()` and `html_to_pptx()` functions (behind a feature flag) can ingest arbitrary web content and convert it to slides. This has no equivalent in either competitor.

**PDF Import**
`from_pdf()` imports each PDF page as an image slide, enabling round-trip conversion from PDF to PPTX without external tools.

**Macro System**
The `pptx!()` and `shape!()` macros provide a concise declarative syntax for building presentations that reads more like a DSL than function calls.

**Shape Layout Helpers**
`center()`, `grid()`, `stack_horizontal()`, and `distribute_horizontal()` automate common slide layout tasks, reducing the manual math required to position groups of shapes.

### 5.3 What Makes PptxGenJS Unique

**Browser-Native Execution**
PptxGenJS is the only library that runs natively in a web browser without any server-side component, WebAssembly wrapper, or native binary. It generates PPTX files entirely in JavaScript, making it deployable as a static HTML page, a Chrome extension, or a serverless function.

**Zero Runtime Dependencies**
No native dependencies, no FFI, no shared libraries — the entire library ships as pure JavaScript. This eliminates an entire class of deployment problems (shared library versions, CGO compilation, Rust toolchain availability).

**Multiple Browser-Native Output Formats**
The ability to produce a file as a `Blob`, `base64` string, `ArrayBuffer`, `Uint8Array`, or binary string covers every integration pattern a web developer might need — direct download, upload to API, embed in form, or store in IndexedDB.

**HTTP Stream Output**
`stream()` enables chunked PPTX delivery via HTTP response streams, useful in Node.js API servers without writing intermediate files.

**HTML Table → PPTX Slides**
`tableToSlides()` converts an existing HTML `<table>` element into paginated PPTX slides. This browser-only feature is unique and valuable for data portal or BI dashboard export scenarios.

**YouTube Embed Support**
Embedding YouTube videos (Microsoft 365 compatible) is a unique feature unavailable in either Go or Rust library.

**SVG Image Support**
PptxGenJS supports SVG images directly, which is important for scalable graphics and icon libraries. Neither gopptx nor ppt-rs support SVG.

**Image URL Support**
Images can be referenced by URL and fetched at generation time, removing the need to download and re-encode media before calling the library.

**Vertical Text (7 Modes)**
Seven vertical text orientation modes are supported, a typographic feature needed for East Asian language presentations.

**Advanced Text Effects**
Shadow, glow, and outline effects on text runs are uniquely available in PptxGenJS among the three libraries.

**Logarithmic Scale**
Logarithmic axis scaling for charts is only available in PptxGenJS.

---

## 6. Gap Analysis: gopptx

This section identifies features available in ppt-rs or PptxGenJS that gopptx currently lacks or could improve.

### 6.1 High Priority Gaps (High Impact, Feasible)

| Gap | Source Library | Impact | Notes |
|---|---|---|---|
| Strikethrough text | PptxGenJS, ppt-rs | High | Common formatting; missing from run-level props |
| Superscript / subscript | PptxGenJS, ppt-rs | High | Scientific and technical content |
| Font highlight color | PptxGenJS, ppt-rs | High | Common in annotation workflows |
| Vertical text alignment (anchor) | PptxGenJS, ppt-rs | High | Required for proper table/shape layout |
| Bullet list support | PptxGenJS, ppt-rs | High | Standard presentation feature |
| Line spacing control | PptxGenJS | High | Essential for layout density control |
| Character spacing | PptxGenJS | Medium | Typography control |
| Paragraph space before/after | PptxGenJS | Medium | Layout control |
| Auto-fit text (shrink/resize) | PptxGenJS | Medium | Overflow handling |
| Text wrap control | PptxGenJS | Medium | Layout control |
| Image crop | ppt-rs, PptxGenJS | High | Very common image operation |
| Image rotation | ppt-rs, PptxGenJS | High | Common visual treatment |
| Image opacity | PptxGenJS | Medium | Common visual treatment |
| Image shadow / glow | ppt-rs, PptxGenJS | Medium | Polish effects |
| SVG image support | PptxGenJS | Medium | Icon and graphic libraries |
| Horizontal bar chart | ppt-rs, PptxGenJS | High | Common chart variant |
| Stacked bar chart | ppt-rs, PptxGenJS | High | Very common in business decks |
| Scatter / XY chart | ppt-rs, PptxGenJS | High | Data visualization |
| Doughnut chart | ppt-rs, PptxGenJS | Medium | Common KPI display |
| Area chart | ppt-rs, PptxGenJS | Medium | Time-series visualization |
| Table auto-pagination | PptxGenJS | Medium | Long data tables |
| Table repeat header | PptxGenJS | Medium | Multi-page tables |
| Bullets in table cells | PptxGenJS | Medium | Rich table content |
| Shape shadow | ppt-rs, PptxGenJS | Medium | Visual polish |
| Shape flip (H/V) | PptxGenJS | Low | Mirror transformations |

### 6.2 Medium Priority Gaps

| Gap | Source Library | Impact | Notes |
|---|---|---|---|
| Color manipulation API (.lighter/.darker/.opacity) | ppt-rs | Medium | Useful for algorithmic theming |
| Embedded fonts | ppt-rs | Medium | Consistent rendering in locked environments |
| Sub-bullets / hierarchical bullets | ppt-rs | Medium | Complex list content |
| RTL text support | ppt-rs, PptxGenJS | Medium | Arabic, Hebrew, Persian content |
| Language tag on runs | PptxGenJS | Medium | Spell check and accessibility |
| Code block with syntax highlighting | ppt-rs | Low | Technical presentations |
| 3D model embedding | ppt-rs | Low | Modern PowerPoint feature |
| Digital signature creation | ppt-rs | Medium | Compliance workflows |
| Image from URL | PptxGenJS | Medium | Simplifies image pipelines |
| Image reflection / soft edges / blur | ppt-rs | Low | Visual polish |
| Shape layout helpers (grid/center) | ppt-rs | Low | Ergonomic positioning |
| Hidden slide flag | PptxGenJS | Low | Presenter-mode specific |

### 6.3 Low Priority / Niche Gaps

| Gap | Source Library | Impact | Notes |
|---|---|---|---|
| YouTube embed | PptxGenJS | Low | M365-only feature |
| HTML table → PPTX | PptxGenJS | Low | Browser-specific use case |
| Logarithmic chart scale | PptxGenJS | Low | Specialized chart need |
| Vertical text (7 modes) | PptxGenJS | Low | East Asian language use |
| Web URL → PPTX | ppt-rs | Low | Niche import feature |
| PDF → PPTX import | ppt-rs | Low | Conversion workflow |
| HTML → PPTX | ppt-rs | Low | Conversion workflow |
| Tab stops in text | PptxGenJS | Low | Advanced typography |
| Data table inside chart | PptxGenJS | Low | Chart annotation |
| TableMergeMap advanced API | ppt-rs | Low | Complex table merging |
| HandoutLayout type | ppt-rs | Low | Handout-specific |
| IBM Carbon / Material Design color constants | ppt-rs | Low | Design system shortcuts |

---

## 7. Gap Analysis: ppt-rs

This section identifies features in gopptx or PptxGenJS that ppt-rs currently lacks.

### 7.1 Editing & Document Automation Gaps

These are the most significant gaps, as they define the "edit existing file" capability that ppt-rs partially supports but does not expose through a rich API.

| Gap | Source Library | Impact | Notes |
|---|---|---|---|
| Find and replace | gopptx | High | Core document automation |
| Get/set text runs with formatting | gopptx | High | Template filling, run-level edit |
| List shapes with metadata | gopptx | High | Document introspection |
| Remove shape | gopptx | High | Editing existing slides |
| Move shape to front/back/index | gopptx | Medium | Z-order control |
| Group / ungroup shapes | gopptx | Medium | Shape management |
| Move / reorder slides | gopptx | High | Slide management |
| Duplicate slide | gopptx | High | Common editing operation |
| Update slide title programmatically | gopptx | High | Template filling |
| Rebind slide layout | gopptx | Medium | Template management |
| List / inspect placeholders | gopptx | High | Template introspection |
| Set placeholder content directly | gopptx | High | Template filling |
| Get chart state / data | gopptx | High | Chart editing |
| Update existing chart data | gopptx | High | Data refresh workflows |
| List charts on slide | gopptx | Medium | Document introspection |
| Get existing table state | gopptx | High | Table editing |
| List images on slide | gopptx | Medium | Asset management |
| Swap image by index / relID | gopptx | Medium | Image replacement workflows |
| Get image metadata | gopptx | Medium | Asset inventory |
| Comments (get/add/remove) | gopptx | Medium | Collaborative annotation |
| Header / footer management | gopptx | Medium | Print layout |
| Document protection / password | gopptx | Medium | Security |
| Mark as final | gopptx | Low | Document lifecycle |
| Notes shape management | gopptx | Medium | Rich notes editing |
| Update notes master | gopptx | Low | Notes template |
| Section management (full CRUD) | gopptx | Medium | Presentation structure |
| Shape animations | gopptx | High | Presentation delivery |
| Animation triggers | gopptx | High | Presentation delivery |
| Clone layout+master family | gopptx | Medium | Theme propagation |
| Remove slide master / layout | gopptx | Medium | Cleanup operations |
| Get layout/master shapes | gopptx | Medium | Introspection |
| OLE object embedding | gopptx | Low | Legacy document embedding |

### 7.2 Content Generation Gaps

| Gap | Source Library | Impact | Notes |
|---|---|---|---|
| Jinja2 template rendering | gopptx | High | Data-driven presentation generation |
| 5 built-in template builders | gopptx | Medium | Rapid complete-deck generation |
| Batch operations (single bridge call) | gopptx | High | Performance-critical bulk generation |
| Python bindings | gopptx | High | Python ecosystem access |
| HTML table → PPTX | PptxGenJS | Low | Web content conversion |

### 7.3 Output Format Gaps

| Gap | Source Library | Impact | Notes |
|---|---|---|---|
| PDF export (native / multiple drivers) | gopptx | High | ppt-rs requires LibreOffice |
| PowerPoint-driven PDF export | gopptx | Medium | Highest fidelity export |
| Browser-native output | PptxGenJS | Low | Not applicable for Rust |
| HTTP stream output | PptxGenJS | Low | Not applicable for Rust CLI |

### 7.4 Text Formatting Gaps in ppt-rs

| Gap | Source Library | Impact | Notes |
|---|---|---|---|
| Vertical text (7 modes) | PptxGenJS | Low | East Asian layouts |
| Auto-fit text | PptxGenJS | Medium | Overflow handling |
| Text wrap control | PptxGenJS | Medium | Layout control |
| Tab stops | PptxGenJS | Low | Advanced typography |
| Text shadow / glow / outline effects | PptxGenJS | Low | Visual effects |
| Find and replace | gopptx | High | Document automation |
| Run-level text inspection | gopptx | High | Editing existing text |

---

## 8. Gap Analysis: PptxGenJS

This section identifies features in gopptx or ppt-rs that PptxGenJS currently lacks. Given PptxGenJS is creation-only, many of these reflect fundamental architectural decisions rather than simple omissions.

### 8.1 Fundamental Architectural Gaps

| Gap | Source Library | Impact | Notes |
|---|---|---|---|
| Open existing PPTX | gopptx, ppt-rs | Critical | Core editing capability |
| Modify existing PPTX | gopptx, ppt-rs | Critical | Requires open support |
| Remove slides | gopptx, ppt-rs | High | Basic slide management |
| Reorder slides | gopptx | High | Slide organization |
| Duplicate slides | gopptx | High | Template reuse |
| Find and replace | gopptx | High | Document automation |
| Validate & repair | gopptx, ppt-rs | High | Document health |
| Merge presentations | gopptx, ppt-rs | Medium | Combining decks |

### 8.2 Presentation Feature Gaps

| Gap | Source Library | Impact | Notes |
|---|---|---|---|
| Slide transitions | gopptx, ppt-rs | High | Presentation delivery |
| Shape animations | gopptx | High | Presentation delivery |
| SmartArt | gopptx, ppt-rs | High | Diagram content |
| Mermaid diagrams | gopptx, ppt-rs | Medium | Technical diagram import |
| Built-in themes (named) | gopptx, ppt-rs | High | Consistent visual design |
| Full theme XML (color + font scheme) | gopptx, ppt-rs | High | Brand consistency |
| Embedded fonts | ppt-rs | Medium | Cross-device fidelity |
| Digital signatures | ppt-rs | Medium | Compliance |
| 3D model embedding | ppt-rs | Low | Modern presentations |
| OLE object embedding | gopptx | Low | Legacy office objects |
| VBA / macros | gopptx, ppt-rs | Medium | Office automation |
| Custom XML parts | gopptx, ppt-rs | Low | Enterprise metadata |
| Document protection | gopptx | Medium | Security |
| Comments | gopptx | Medium | Collaboration metadata |
| Headers / footers | gopptx | Medium | Print layout |
| Sections CRUD | gopptx | Medium | (PptxGenJS has sections, limited CRUD) |
| Handout master | gopptx, ppt-rs | Low | Print handout |

### 8.3 Chart Gaps in PptxGenJS

| Gap | Source Library | Impact | Notes |
|---|---|---|---|
| Stock chart (HLC/OHLC) | ppt-rs | Medium | Financial data |
| Radar filled | ppt-rs | Low | Data visualization variant |
| Line stacked | ppt-rs | Low | Chart variant |
| Get / update existing chart | gopptx | High | (N/A - creation only) |

### 8.4 Image & Media Gaps in PptxGenJS

| Gap | Source Library | Impact | Notes |
|---|---|---|---|
| Image from bytes (raw buffer) | ppt-rs | Low | Programmatic image generation |
| SHA256 deduplication | gopptx | Medium | Reduces file size |
| Reflection / soft edges / blur effects | ppt-rs | Low | Polish effects |
| 3D models | ppt-rs | Low | Modern media |
| Audio with icon | gopptx | Low | Audio slides |
| Video auto-play / loop / muted | gopptx | Medium | Kiosk/auto-run presentations |

### 8.5 Text Gaps in PptxGenJS

| Gap | Source Library | Impact | Notes |
|---|---|---|---|
| Code block with syntax highlighting | ppt-rs | Low | Technical presentations |
| Sub-bullets / hierarchical bullets | ppt-rs | Medium | Complex list structure |
| RTL language-level (RtlLanguage enum) | ppt-rs | Medium | Script-level RTL control |

---

## 9. Use Case Recommendations

### 9.1 Choose gopptx When…

| Scenario | Reason |
|---|---|
| Python-based automation pipeline | Only library with full Python bindings |
| Template-based document generation | Jinja2, placeholder API, find-and-replace |
| Opening and modifying existing PPTX files | Deepest editing API with 145+ operations |
| Go microservice generating presentations | Native Go API, JSON bridge for multi-language |
| Batch-generating 100s of slides | BatchExecute minimizes bridge overhead |
| PDF export with multiple fallback drivers | Flexible export (native/LibreOffice/PowerPoint) |
| Mermaid + Markdown → PPTX pipeline | Both supported with same quality as ppt-rs |
| Presentation auditing / content inventory | Shape listing, text extraction, chart inspection |
| Chart data refresh in existing files | `UpdateChartData`, `GetChartState` |
| Image management (swap, dedup, audit) | SHA256 dedup, swap by index/relID |
| Data portal with scheduled report generation | Python + BatchExecute + themes |
| Corporate template system | Layout/master CRUD, placeholder injection |
| SmartArt-heavy corporate decks | 25 layouts, updateable content |
| Slide animations for live demos | Full animation + transition support |

### 9.2 Choose ppt-rs When…

| Scenario | Reason |
|---|---|
| Rust-native application | No FFI overhead, compile-time safety |
| Compliance documents with digital signatures | Only library supporting signature creation |
| Presentations requiring embedded fonts | Guaranteed font rendering on any machine |
| Documents with 3D model content | GLB/GLTF/OBJ/FBX/STL support |
| Technical presentations with code samples | Syntax-highlighted code blocks |
| Financial reports with Stock/OHLC charts | Stock chart types not in other libraries |
| Programmatic color derivation (lighter/darker) | Color manipulation API |
| PDF import → PPTX | Convert PDF pages to image slides |
| Web scraping → PPTX pipeline | url_to_pptx(), html_to_pptx() |
| Broad chart type coverage | 22+ chart types vs gopptx's 3 |
| Image-heavy decks with polish effects | Shadow, reflection, glow, blur on images |
| Hierarchical bullet content | Sub-bullets with 6+ numbering styles |
| Maximum Rust ergonomics | Fluent builder API, pptx!() macros |

### 9.3 Choose PptxGenJS When…

| Scenario | Reason |
|---|---|
| In-browser PPTX generation | Only library that runs client-side |
| Zero-dependency requirement | No native binaries, no build tools |
| Serverless / edge function deployment | Pure JS, no shared library dependencies |
| JavaScript/TypeScript codebase | Native ecosystem integration |
| HTML table → PPTX export feature | tableToSlides() is unique |
| YouTube video embedding | M365 YouTube embed support |
| Report download button in web app | Blob download, no server required |
| Node.js API server with stream output | stream() for HTTP chunked delivery |
| SVG-heavy presentations | SVG image support |
| Multi-format output (base64/Blob/Buffer) | All JS-native output formats |
| BI dashboard export | Browser-native generation from live data |
| Simple greenfield presentation generation | Excellent DX for basic creation use cases |
| TypeScript-first projects | First-class .d.ts definitions |

### 9.4 Decision Matrix

```
Start here:
│
├─ Do you need to open / modify existing PPTX files?
│   ├─ Yes → gopptx or ppt-rs
│   │   ├─ Python integration needed?     → gopptx
│   │   ├─ Deep chart editing needed?     → gopptx
│   │   ├─ Find/replace or template fill? → gopptx
│   │   ├─ Animations required?           → gopptx
│   │   ├─ Digital signatures?            → ppt-rs
│   │   ├─ Embedded fonts?                → ppt-rs
│   │   └─ Pure Rust codebase?            → ppt-rs
│   │
│   └─ No (create-only)
│       ├─ Browser / serverless target?   → PptxGenJS
│       ├─ Zero dependency requirement?   → PptxGenJS
│       ├─ JavaScript codebase?           → PptxGenJS
│       ├─ Python codebase?               → gopptx
│       ├─ Rust codebase?                 → ppt-rs
│       ├─ Broad chart types needed?      → ppt-rs
│       ├─ Code block slides?             → ppt-rs
│       ├─ HTML table import?             → PptxGenJS
│       └─ Mermaid / SmartArt?            → gopptx or ppt-rs
```

---

## 10. Feature Category Deep-Dives

### 10.1 Charts

Charts are one of the most significant differentiators across the three libraries. The breadth of supported chart types, the depth of formatting control, and the ability to edit existing chart data vary substantially.

#### Chart Type Coverage

| Chart Type | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Bar (clustered vertical) | ✅ | ✅ | ✅ |
| Bar (clustered horizontal) | ❌ | ✅ | ✅ |
| Bar (stacked) | ❌ | ✅ | ✅ |
| Bar (stacked 100%) | ❌ | ✅ | ✅ |
| Bar 3D | ❌ | ❌ | ✅ |
| Line | ✅ | ✅ | ✅ |
| Line with markers | ❌ | ✅ | ✅ |
| Line stacked | ❌ | ✅ | ❌ |
| Line smooth | ❌ | ❌ | ✅ |
| Area | ❌ | ✅ | ✅ |
| Area stacked | ❌ | ✅ | ✅ |
| Area stacked 100% | ❌ | ✅ | ✅ |
| Pie | ✅ | ✅ | ✅ |
| Doughnut | ❌ | ✅ | ✅ |
| Scatter | ❌ | ✅ | ✅ |
| Scatter with lines | ❌ | ✅ | ✅ |
| Scatter smooth | ❌ | ✅ | ✅ |
| Bubble | ❌ | ✅ | ✅ |
| Bubble 3D | ❌ | ❌ | ✅ |
| Radar | ❌ | ✅ | ✅ |
| Radar filled | ❌ | ✅ | ❌ |
| Stock HLC | ❌ | ✅ | ❌ |
| Stock OHLC | ❌ | ✅ | ❌ |
| Combo | ❌ | ✅ | ✅ |
| **Total types** | **3** | **22** | **18** |

**Key observation**: gopptx has the most restricted chart type set among the three libraries. For data-rich presentations, this is a significant gap. ppt-rs leads with 22 types including the unique Stock HLC/OHLC charts. PptxGenJS covers most common business chart types except stock charts.

#### Chart Formatting Depth

| Formatting Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Chart title | ✅ | ✅ | ✅ |
| Chart legend (position, font) | ✅ Full | ⚡ Position only | ✅ Full |
| Value axis (title, gridlines) | ✅ Full | ⚡ Partial | ✅ Full |
| Category axis (title, rotation) | ✅ Full | ⚡ Partial | ✅ Full |
| Axis number format | ✅ | ❌ | ✅ |
| Logarithmic scale | ❌ | ❌ | ✅ |
| Data labels (value/percent/name) | ✅ Full | ⚡ Partial | ✅ Full |
| Leader lines | ❌ | ❌ | ✅ |
| Data table in chart | ❌ | ❌ | ✅ |
| Per-series colors | ✅ | ✅ | ✅ |
| Chart area formatting | ❌ | ❌ | ✅ |
| Plot area formatting | ❌ | ❌ | ✅ |
| First slice angle (pie) | ❌ | ❌ | ✅ |
| Bar gap width | ❌ | ❌ | ✅ |
| Display units on axis | ❌ | ❌ | ✅ |

#### Chart Editing (Existing Files)

| Editing Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Get chart state (data + formatting) | ✅ | ❌ | ❌ |
| Update chart data | ✅ | ❌ | ❌ |
| Batch update multiple charts | ✅ | ❌ | ❌ |
| List charts on slide | ✅ | ❌ | ❌ |

gopptx is the **only library** supporting chart editing in existing files. For data refresh pipelines (e.g., weekly report decks where only the data changes), gopptx is the only viable option.

#### Roadmap Suggestion
gopptx should prioritize expanding chart types to at least cover: horizontal bar, stacked bar/column, stacked 100%, doughnut, scatter, area, and combo. These account for the vast majority of real-world presentation charts. The underlying Excel embedded data model (`excel_gen.go`) already supports multi-series data; extending to more chart types is primarily an XML schema exercise.

---

### 10.2 Text & Formatting

Text formatting is a core presentation capability where PptxGenJS leads in typographic depth, ppt-rs leads in list/bullet support, and gopptx leads in text editing and run-level inspection.

#### Run-Level Properties

| Property | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Bold | ✅ | ✅ | ✅ |
| Italic | ✅ | ✅ | ✅ |
| Underline | ✅ | ✅ | ✅ (style + color) |
| Strikethrough | ❌ | ✅ | ✅ (single/double) |
| Superscript | ❌ | ✅ | ✅ |
| Subscript | ❌ | ✅ | ✅ |
| Font family | ✅ | ✅ | ✅ |
| Font size | ✅ | ✅ (7 presets) | ✅ |
| Font color | ✅ | ✅ | ✅ |
| Highlight color | ❌ | ✅ | ✅ |
| Transparency | ❌ | ❌ | ✅ |
| Shadow | ❌ | ❌ | ✅ |
| Glow | ❌ | ❌ | ✅ |
| Outline | ❌ | ❌ | ✅ |
| Language tag | ❌ | ⚡ RTL only | ✅ |
| Hyperlink | ✅ | ✅ | ✅ |

#### Paragraph-Level Properties

| Property | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Horizontal alignment | ✅ | ✅ | ✅ |
| Vertical alignment (anchor) | ❌ | ✅ | ✅ |
| Line spacing | ❌ | ❌ | ✅ (points/multiple) |
| Space before | ❌ | ❌ | ✅ |
| Space after | ❌ | ❌ | ✅ |
| Indent | ❌ | ✅ (bullets) | ✅ |
| Tab stops | ❌ | ❌ | ✅ |
| Bullet (standard) | ❌ | ✅ | ✅ |
| Bullet (numbered) | ❌ | ✅ | ✅ |
| Bullet (custom/Unicode) | ❌ | ✅ | ✅ |
| Sub-bullets / hierarchy | ❌ | ✅ | ⚡ Indent levels |
| RTL paragraph | ❌ | ✅ | ✅ |
| Vertical text (East Asian) | ❌ | ❌ | ✅ (7 modes) |

#### Text Introspection (Editing Existing Text)

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Get text runs from shape | ✅ | ❌ | ❌ |
| Get run formatting | ✅ | ❌ | ❌ |
| Set run text | ✅ | ❌ | ❌ |
| Set run formatting | ✅ | ❌ | ❌ |
| Append run to shape | ✅ | ❌ | ❌ |
| Update entire slide texts | ✅ | ❌ | ❌ |
| Update deck-wide texts | ✅ | ❌ | ❌ |
| Find and replace text | ✅ | ❌ | ❌ |

**Summary**: Text introspection and editing is exclusively in gopptx. For read/write workflows this is critical. For creation workflows, PptxGenJS offers the deepest typographic control (line spacing, paragraph spacing, auto-fit, text wrap, effects). ppt-rs leads on list/bullet depth and is the only library with hierarchical sub-bullets and code block syntax highlighting.

---

### 10.3 Tables

Tables are well-supported across all three libraries for basic use cases. The differentiators are in advanced scenarios.

#### Table Construction

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Basic table creation | ✅ | ✅ | ✅ |
| QuickTable helper | ❌ | ✅ | ❌ |
| Cell text formatting | ✅ | ✅ | ✅ |
| Cell fill / background | ✅ | ✅ | ✅ |
| Cell border styles | ✅ | ✅ | ✅ |
| BorderStyle variations | ⚡ | ✅ | ✅ |
| Cell alignment (H/V) | ✅ | ✅ | ✅ |
| Colspan | ✅ | ✅ | ✅ |
| Rowspan | ✅ | ✅ | ✅ |
| TableMergeMap | ❌ | ✅ | ❌ |
| Row height | ✅ | ✅ | ✅ |
| Column width | ✅ | ✅ | ✅ |
| Named table style | ✅ | ❌ | ❌ |
| Custom table style definition | ✅ | ❌ | ❌ |
| Bullets in cells | ❌ | ❌ | ✅ |
| Auto-pagination | ❌ | ❌ | ✅ |
| Repeat header row | ❌ | ❌ | ✅ |
| HTML table import | ❌ | ❌ | 🔧 Browser |

#### Table Editing (Existing Files)

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Get table state | ✅ | ❌ | ❌ |
| Update cell content | ✅ | ❌ | ❌ |
| Merge cells in existing table | ✅ | ❌ | ❌ |
| Split cells in existing table | ✅ | ❌ | ❌ |
| Set row height in existing table | ✅ | ❌ | ❌ |
| Set column width in existing table | ✅ | ❌ | ❌ |
| Apply named style to existing table | ✅ | ❌ | ❌ |
| List named styles | ✅ | ❌ | ❌ |

**Notable**: gopptx is uniquely positioned for data table updates in existing presentations (e.g., a monthly board deck where slide layouts are fixed but data changes). The ability to update individual cells, apply named styles, and split/merge cells in existing files covers a wide range of real-world needs.

ppt-rs has the `TableMergeMap` API for handling complex merged cell configurations declaratively — useful when building tables programmatically and merging is complex enough to need a map-based specification.

PptxGenJS's `tableToSlides()` auto-pagination is genuinely useful for the common scenario of exporting a long data table from a web application where the content might not fit on one slide.

---

### 10.4 Images & Media

#### Image Capabilities

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| File path source | ✅ | ✅ | ✅ |
| Base64 source | ✅ | ✅ | ✅ |
| URL source | ❌ | ❌ | ✅ |
| Bytes source | ❌ | ✅ | ⚡ |
| JPEG | ✅ | ✅ | ✅ |
| PNG | ✅ | ✅ | ✅ |
| GIF | ✅ | ✅ | ✅ |
| SVG | ❌ | ❌ | ✅ |
| Auto-detect format | ✅ | ✅ | ✅ |
| Deduplication (SHA256) | ✅ | ❌ | ❌ |
| Positioning | ✅ | ✅ | ✅ |
| Sizing | ✅ | ✅ | ✅ |
| Contain / cover sizing | ❌ | ❌ | ✅ |
| Crop | ❌ | ✅ | ✅ |
| Rotation | ❌ | ✅ | ✅ |
| Opacity | ❌ | ❌ | ✅ |
| Circle crop / rounding | ❌ | ❌ | ✅ |
| Shadow | ❌ | ✅ | ✅ |
| Reflection | ❌ | ✅ | ❌ |
| Glow | ❌ | ✅ | ✅ |
| Soft edges | ❌ | ✅ | ❌ |
| Blur | ❌ | ✅ | ❌ |
| Inner shadow | ❌ | ✅ | ❌ |
| Hyperlink | ✅ | ✅ | ✅ |
| Alt text | ❌ | ❌ | ✅ |
| Get image metadata | ✅ | ❌ | ❌ |
| Swap image by index | ✅ | ❌ | ❌ |
| Swap image by relID | ✅ | ❌ | ❌ |
| List images on slide | ✅ | ❌ | ❌ |

SHA256 deduplication in gopptx is a practical file-size optimization: when the same image asset is used on multiple slides (e.g., a logo), only one copy is stored in the PPTX ZIP. This can reduce file sizes significantly for asset-heavy decks.

ppt-rs has the most comprehensive image effect pipeline, supporting the complete set of PowerPoint picture format effects. If visual polish matters (e.g., a design agency's presentation tool), ppt-rs is clearly superior for images.

#### Video & Audio

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| MP4 video | ✅ | ✅ | ✅ |
| WMV / AVI / MOV / MKV | ❌ | ✅ | ✅ |
| Video auto-play | ✅ | ❌ | ❌ |
| Video loop | ✅ | ❌ | ❌ |
| Video muted | ✅ | ❌ | ❌ |
| Video volume | ✅ | ❌ | ❌ |
| MP3 / WAV audio | ✅ | ✅ | ✅ |
| WMA / OGG / FLAC / AAC | ❌ | ✅ | ❌ |
| Audio display icon | ✅ | ❌ | ❌ |
| YouTube embed | ❌ | ❌ | ✅ |
| 3D models (GLB/GLTF/OBJ) | ❌ | ✅ | ❌ |
| 3D camera preset | ❌ | ✅ | ❌ |
| OLE object | ✅ | ❌ | ❌ |

gopptx's video playback control properties (auto_play, loop, muted, volume) are important for kiosk-style presentations that run without user interaction. ppt-rs supports broader codec formats. Neither supports 3D models except ppt-rs.

---

### 10.5 Themes

All three libraries share the same 7 named built-in themes (Corporate, Modern, Vibrant, Dark, Nature, Tech, Carbon) in gopptx and ppt-rs. PptxGenJS has no equivalent theme system.

#### Theme Capabilities

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| 7 named built-in themes | ✅ | ✅ | ❌ |
| Apply theme to presentation | ✅ | ✅ | ❌ |
| Custom color scheme (6 accent colors) | ✅ | ✅ | ⚡ SchemeColor ref |
| Custom font scheme (heading/body) | ✅ | ✅ | ⚡ Font face only |
| Get theme inventory | ✅ | ❌ | ❌ |
| Set global theme preset | ✅ | ❌ | ❌ |
| Color manipulation (.lighter/.darker) | ❌ | ✅ | ❌ |
| Color opacity | ❌ | ✅ | ❌ |
| Color mix | ❌ | ✅ | ❌ |
| Embedded fonts | ❌ | ✅ | ❌ |
| .thmx file import | ❌ | ❌ | ❌ |
| Material Design colors | ❌ | ✅ | ❌ |
| IBM Carbon colors | ❌ | ✅ | ❌ |

Neither gopptx nor ppt-rs support `.thmx` file import, which would allow applying a fully designed Office theme from an external file. This is a gap for enterprise environments where design teams maintain branded `.thmx` files.

ppt-rs's color manipulation API (`.lighter(0.2)`, `.darker(0.3)`, `.opacity(0.7)`, `.mix(other, 0.5)`) is a significant ergonomic advantage when generating decks from data. It enables generating consistent color families from a single base color without hardcoding every shade.

---

### 10.6 Export & Output

| Export Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| PPTX file (disk) | ✅ | ✅ | ✅ |
| PPTX bytes / buffer | ✅ | ✅ | ✅ |
| PPTX base64 | ✅ | ✅ | ✅ |
| PPTX Blob (browser) | ❌ | ❌ | ✅ |
| PPTX ArrayBuffer | ❌ | ❌ | ✅ |
| PPTX Uint8Array | ❌ | ❌ | ✅ |
| PPTX HTTP stream | ❌ | ❌ | ✅ |
| PDF (native) | ✅ | ❌ | ❌ |
| PDF (LibreOffice) | ✅ | ✅ | ❌ |
| PDF (PowerPoint) | ✅ | ❌ | ❌ |
| PDF (auto-detect driver) | ✅ | ❌ | ❌ |
| HTML export | ✅ | ✅ | ❌ |
| HTML with embedded images | ✅ | ✅ | ❌ |
| HTML with navigation | ✅ | ❌ | ❌ |
| HTML with slide numbers | ✅ | ❌ | ❌ |
| PNG per slide | ❌ | 🔧 pdftoppm | ❌ |
| ZIP compression control | ❌ | ❌ | ✅ |

#### PDF Export Notes

gopptx's multi-driver PDF export is the most flexible approach. The "auto" driver selects the best available renderer (native → LibreOffice → PowerPoint), enabling graceful degradation across environments. The "native" driver requires no external tools but may have fidelity limitations. The "LibreOffice" driver requires `soffice` on PATH but produces high-fidelity PDFs. The "PowerPoint" driver requires a Windows machine with PowerPoint installed and produces pixel-perfect PDFs.

ppt-rs's PDF export is LibreOffice-only, which is the most common production approach but requires the LibreOffice runtime.

PptxGenJS has no PDF export — conversion must be done externally if needed.

#### HTML Export Notes

gopptx's HTML export supports navigation controls and slide numbers in the output, making it suitable for creating browser-viewable presentation slides without a presentation runtime. ppt-rs produces HTML with embedded images. Neither supports interactive presenter mode.

---

### 10.7 Animations & Transitions

This is one of the starkest differentiators: PptxGenJS has no animation or transition support whatsoever.

#### Transitions

| Transition Type | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| None (instant) | ✅ | ✅ | ❌ |
| Fade | ✅ | ✅ | ❌ |
| Cut | ✅ | ✅ | ❌ |
| Push | ✅ | ✅ | ❌ |
| Wipe | ✅ | ✅ | ❌ |
| Split | ✅ | ✅ | ❌ |
| Reveal | ✅ | ✅ | ❌ |
| Cover | ✅ | ✅ | ❌ |
| Zoom | ✅ | ✅ | ❌ |

Both gopptx and ppt-rs support the same 9 slide transition types. Transitions are essential for presentations that will be delivered live. PptxGenJS presentations will always play with no transitions unless added manually in PowerPoint after the fact.

#### Shape Animations

| Feature | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Fade animation | ✅ | ❌ | ❌ |
| Appear animation | ✅ | ❌ | ❌ |
| Fly Out animation | ✅ | ❌ | ❌ |
| Additional animation types | ✅ | ❌ | ❌ |
| onClick trigger | ✅ | ❌ | ❌ |
| withPrev trigger | ✅ | ❌ | ❌ |
| afterPrev trigger | ✅ | ❌ | ❌ |

gopptx is the **only library** with shape animation support. For live-delivered presentations where shapes should animate in on click or in sequence, gopptx is the only programmatic option. This is a meaningful differentiator for use cases like:

- Automated generation of training materials with step-by-step reveals
- Sales decks with animated charts and data reveals
- Marketing presentations with entrance and exit effects

---

## 11. API Style Comparison

### 11.1 API Paradigm Overview

| Library | Primary Paradigm | Auxiliary Patterns |
|---|---|---|
| gopptx | Command/response (JSON bridge) + OOP `Presentation` API | Python bindings with builders |
| ppt-rs | Fluent builder + macros | Unit helpers (inch/cm/pt) |
| PptxGenJS | Fluent method chaining (OOP) | TypeScript interfaces |

### 11.2 gopptx API Style

gopptx exposes two API surfaces:

**Go API** — A traditional OOP API via `PresentationEditor` and the higher-level `Presentation` wrapper in `presentation_api.go`. This is straightforward and idiomatic for Go developers.

```go
// Go API example
prs, _ := pptx.Open("template.pptx")
defer prs.Close()
prs.SetTitle("Q4 Report")
slide := prs.Slides()[0]
slide.SetTitle("Revenue Overview")
prs.SaveAs("q4_report.pptx")
```

**Python API** — The Python bindings layer wraps the Go JSON bridge, providing a Pythonic interface with type-safe enums and builders.

```python
# Python API example
from gopptx import Presentation, ops

with Presentation.open("template.pptx") as pres:
    pres.set_title("Q4 Report")
    results = pres.execute_batch([
        {"op": ops.OP_SET_SLIDE_TITLE, "payload": {"slide_index": 0, "title": "Revenue Overview"}},
        {"op": ops.OP_UPDATE_CHART_DATA, "payload": {"slide_index": 0, "chart_index": 0, "data": {...}}}
    ])
    pres.save_as("q4_report.pptx")
```

**JSON Bridge API** — The lowest-level API, usable from any language via the C bridge. Operations are sent as `RequestEnvelope` JSON objects with `op`, `payload`, `api_version`, and `request_id` fields. This enables calling gopptx from languages other than Go and Python (Ruby, PHP, .NET, etc. via subprocess or C FFI).

The bridge's `BatchExecute` pattern is the performance-critical feature for bulk generation: a single call containing 200 add-slide operations has far lower overhead than 200 individual bridge calls. This is especially important in the Python layer where each round-trip incurs marshaling overhead.

**Strengths**: The Go API is clean and idiomatic. The Python binding layer with type-safe enums (`ShapeType.ROUNDED_RECTANGLE`, `ChartType.LINE`, `SlideLayoutType.TITLE_AND_CONTENT`) reduces runtime errors in Python code that would otherwise pass invalid string constants. The `RunBuilder` and `ShapeBuilder` patterns in Python align with modern Python conventions.

**Weaknesses**: The JSON bridge adds indirection and makes the Go API slightly less natural than a pure Go library. Error messages from deep in the bridge can be verbose. The two-layer architecture (Go + Python) increases maintenance surface.

### 11.3 ppt-rs API Style

ppt-rs uses Rust's fluent builder pattern extensively, leveraging ownership and the type system for correctness.

```rust
// ppt-rs example - fluent builder
let pptx = Pptx::builder()
    .with_title("Q4 Report")
    .add_slide(
        Slide::builder()
            .add_shape(
                Shape::rect(2.0, 2.0, 6.0, 3.0)
                    .with_text("Revenue Overview")
                    .with_fill(ShapeFill::solid(Color::CORPORATE_BLUE))
                    .with_shadow()
            )
            .add_chart(
                ChartBuilder::bar()
                    .add_series("Q3", vec![100.0, 200.0, 150.0])
                    .add_series("Q4", vec![120.0, 210.0, 180.0])
            )
    )
    .build();
```

The `pptx!()` and `shape!()` macros provide a more declarative syntax:

```rust
let pptx = pptx! {
    title: "Q4 Report",
    slides: [
        slide! {
            shapes: [
                shape!(rect [2, 2, 6, 3] "Revenue Overview"),
            ]
        }
    ]
};
```

**Unit helpers** (`inch(2.0)`, `cm(5.0)`, `pt(12.0)`) eliminate the mental overhead of converting between measurement systems.

**Strengths**: The fluent builder pattern is ergonomic and Rust-idiomatic. Compile-time type checking prevents many categories of invalid input. The macro system enables concise declaration of simple presentations. Unit helpers and color manipulation methods reduce boilerplate.

**Weaknesses**: The builder pattern can become deeply nested for complex slides. Rust's ownership model can make dynamic presentation construction (building slides in a loop from external data) more verbose than Python or JavaScript equivalents. No reflection/introspection API for existing files.

### 11.4 PptxGenJS API Style

PptxGenJS uses a method-chaining OOP style familiar to JavaScript developers.

```javascript
// PptxGenJS example
const pptx = new PptxGenJS();
pptx.title = "Q4 Report";

const slide = pptx.addSlide();
slide.addText("Revenue Overview", {
    x: 1, y: 1, w: 8, h: 1.5,
    fontSize: 24, bold: true, color: "363636"
});
slide.addChart(pptx.ChartType.bar, [
    { name: "Q3", labels: ["Jan","Feb","Mar"], values: [100, 200, 150] },
    { name: "Q4", labels: ["Jan","Feb","Mar"], values: [120, 210, 180] }
], { x: 1, y: 2, w: 8, h: 4 });

await pptx.writeFile({ fileName: "q4_report.pptx" });
```

**Strengths**: Familiar to any JavaScript developer. Async/await for file output is natural in Node.js. TypeScript definitions provide IntelliSense. Zero dependencies makes it instantly usable in any JS project. The option objects are flexible and forgiving (unknown keys are ignored).

**Weaknesses**: Option objects are large and use string constants rather than type-safe enums, leading to potential runtime errors for typos (mitigated in TypeScript but not eliminated). No read capability fundamentally limits use to creation workflows. The large option object pattern can be verbose for heavily formatted content.

### 11.5 Python Bindings Consideration

gopptx's Python bindings are a significant strategic asset. The Python data ecosystem (pandas, numpy, matplotlib, sqlalchemy, Jinja2) is the de facto standard for data analysis and reporting pipelines. Having a PPTX library with:

- Native Python package installation (`pip install gopptx`)
- Type-safe enums compatible with IDEs and type checkers (basedpyright/mypy)
- Batch operations for high-throughput generation
- Jinja2 integration for template rendering
- pandas-compatible data input for charts and tables

...positions gopptx uniquely for automated reporting and business intelligence workflows. No competing library — ppt-rs or PptxGenJS — offers Python bindings of comparable depth.

The architectural trade-off is the CGO build requirement: Python users must have the compiled Go shared library (`libgopptx.so` / `gopptx.dll` / `libgopptx.dylib`) available. This adds complexity to installation compared to a pure-Python package, but task runners and pre-built binary distribution mitigate this for most deployment scenarios.

---

## 12. Conclusion & Roadmap Suggestions

### 12.1 Summary Positioning

| Library | Best For | Weakest At |
|---|---|---|
| **gopptx** | Editing, Python automation, template filling, animations | Chart type breadth, image effects, text typographic depth |
| **ppt-rs** | Rust systems, visual quality, digital signatures, broad chart coverage | Editing existing files, Python/JS ecosystem, no animations |
| **PptxGenJS** | Browser/web generation, zero-dep deployment, JS ecosystem | Everything requiring file editing, transitions, themes |

### 12.2 Roadmap Suggestions for gopptx

The following are prioritized recommendations based on the gap analysis, ordered by expected impact on the widest range of real-world use cases.

#### Priority 1: Chart Type Expansion (High Impact, Medium Effort)

The jump from 3 chart types to 8-10 would address the most visible functional gap. The recommended additions:

1. **Horizontal bar chart** — Used in almost every business deck for ranking comparisons
2. **Stacked column / stacked bar** — Standard for part-to-whole comparisons
3. **Stacked 100% column / bar** — Percentage breakdowns
4. **Doughnut chart** — KPI dashboards
5. **Area chart + stacked area** — Time series trends
6. **Scatter / XY chart** — Correlation analysis

The `excel_gen.go` embedded workbook generator and `chart_xml.go` template infrastructure already exist; adding chart types is primarily writing XML templates and mapping data structures to the correct OOXML chart schema elements.

#### Priority 2: Text Run Property Completeness (High Impact, Low Effort)

Adding the following run-level properties to the text formatting API would close most practical gaps:

1. **Strikethrough** — Single property addition to run XML
2. **Superscript / subscript** — `<a:vertAlign>` attribute
3. **Font highlight color** — `<a:highlight>` element
4. **Vertical text alignment / anchor** — `<a:bodyPr>` `anchor` attribute

These are all small XML attribute additions to existing code paths in `command_handlers_shapes_notes.go` and the run formatting structures.

#### Priority 3: Bullet List Support (High Impact, Medium Effort)

Bullet lists are fundamental presentation content. Adding a `BulletProperties` struct to the text run model supporting:

- Standard bullets (•, ○, ▪)
- Numbered lists (1., 2., 3.)
- Letter lists (a., b., c.)
- Indent level
- Custom bullet characters

This is a meaningful content authoring gap — many real-world slides consist primarily of bulleted text.

#### Priority 4: Image Enhancement (Medium Impact, Low-Medium Effort)

1. **Image crop** — Crop percentages map to `<a:srcRect>` in OOXML
2. **Image rotation** — `rot` attribute on `<a:xfrm>`
3. **Image shadow** — `<a:effectLst><a:outerShdw>` element; already used for shapes potentially
4. **Image from URL** — Fetch + cache in mediaInventory; already has SHA256 dedup

#### Priority 5: Paragraph-Level Formatting (Medium Impact, Low Effort)

1. **Line spacing** — `<a:lnSpc>` in paragraph properties
2. **Space before / after paragraph** — `<a:spcBef>` / `<a:spcAft>`
3. **Auto-fit text** — `<a:bodyPr>` `autofit` attribute

These affect layout quality significantly for content-heavy slides.

#### Priority 6: Table Auto-Pagination (Medium Impact, High Effort)

`tableToSlides()` in PptxGenJS is a uniquely useful feature for data portal use cases. Implementing equivalent functionality in gopptx would enable Python analytics pipelines to export large DataFrames directly to multi-slide table decks. This requires measuring text height, calculating how many rows fit per slide, and auto-splitting table content across additional slides.

#### Priority 7: Color Manipulation API (Medium Impact, Low Effort)

Adding `.Lighter(factor float64)`, `.Darker(factor float64)`, and `.WithOpacity(alpha float64)` methods to the color type would enable algorithmic color scheme generation. This is especially useful for chart series colors, table header/row alternation, and theme-consistent shape fills.

#### Priority 8: Digital Signature Detection Improvement (Low Impact, Low Effort)

`HasDigitalSignature` detects existing signatures. Adding signature creation (using Go's `crypto/x509` and OOXML XML-Signature spec) would close the gap with ppt-rs for compliance workflows. This is a medium-effort feature with niche but high-value application.

#### Priority 9: .thmx File Import (Low Impact, Medium Effort)

Supporting `.thmx` files would allow designers to create brand themes in PowerPoint and have developers apply them programmatically. The format is a ZIP containing `theme/theme1.xml` — it could be imported and injected into the presentation's theme parts.

#### Priority 10: SVG Image Support (Medium Impact, Medium Effort)

SVG images are increasingly important for icon libraries, logos, and scalable graphics. PowerPoint supports SVG natively since 2016. Adding SVG to the media inventory with appropriate relationship types would close a meaningful gap with PptxGenJS.

### 12.3 Architectural Observations

**The JSON Bridge is a Strategic Asset**: The JSON command bridge design gives gopptx a unique capability — it can be called from any language that can invoke a C function or spawn a subprocess. This is a meaningful long-term advantage as new language ecosystems want PPTX support. Consider documenting the bridge protocol more formally to enable third-party bindings (Ruby, PHP, Elixir, .NET, etc.).

**Python Bindings Competitive Moat**: The Python binding layer is gopptx's clearest competitive advantage. Neither ppt-rs nor PptxGenJS supports Python, yet Python is the dominant language for the analytics and automation use cases where PPTX generation matters most. Continued investment in the Python API (better error messages, pandas integration helpers, more convenience methods) has high leverage.

**BatchExecute Performance Advantage**: The `BatchExecute` pattern enables scale that neither competitor supports. For generating decks of 100-500 slides from database queries or API data, this is the difference between a 2-second and a 200-second operation. Documenting and showcasing this capability more prominently would differentiate gopptx clearly.

**Chart Type Gap is the Most Visible Weakness**: For a developer evaluating PPTX libraries by chart type support, gopptx's 3 types vs ppt-rs's 22 and PptxGenJS's 18 looks like a significant limitation. Expanding to 8-10 types would address the vast majority of real-world use cases and remove this as an evaluation blocker.

---

## Appendix: Feature Count Summary

| Category | gopptx | ppt-rs | PptxGenJS |
|---|---|---|---|
| Chart types | 3 | 22 | 18 |
| SmartArt layouts | 25 | 25 | 0 |
| Built-in themes | 7 | 7 | 0 |
| Slide transition types | 9 | 9 | 0 |
| Animation types | Multiple | 0 | 0 |
| Shape types | 100+ | 100+ | ~200 |
| JSON bridge operations | 145+ | N/A | N/A |
| PDF export drivers | 4 | 1 (LibreOffice) | 0 |
| Output formats | 4 | 3 | 7 |
| Placeholder types | 8 | 0 | 5 |
| Language bindings | Go, Python | Rust | JS, TS |
| Text run properties | ~8 | ~14 | ~18 |
| Image effects | 0 | 6 | 4 |

---

*Report generated from feature inventories of gopptx (main branch, March 2026), ppt-rs v0.2.7, and PptxGenJS v4.0.1. Feature coverage is based on public API documentation and source code inspection. Accuracy depends on the completeness of the feature inventories provided.*
