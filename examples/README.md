# Examples

Task-focused examples and smoke demos for `gopptx`.

## Quick Start

- Run any example from repo root:
  - `go run ./examples/<example-folder>/<file>.go`
- Generated decks are written to:
  - `examples/output/`

## Sequential Index

1. `01-basic-pptx-generation` - What: Create a minimal hello-world deck - Run: `go run ./examples/01-basic-pptx-generation/basic_gen.go`
2. `02-slide-layouts` - What: Show built-in slide layout choices - Run: `go run ./examples/02-slide-layouts/main.go`
3. `03-markdown-to-pptx` - What: Generate slides from Markdown content - Run: `go run ./examples/03-markdown-to-pptx/main.go`
4. `04-text-formatting` - What: Apply rich text formatting and text-frame settings - Run: `go run ./examples/04-text-formatting/text_enhancements.go`; `go run ./examples/04-text-formatting/text_frame.go`
5. `04-text-styling` - What: Demonstrate title/body styling combinations - Run: `go run ./examples/04-text-styling/main.go`
6. `05-bullet-styles` - What: Use bullet, numbered, and nested bullet styles - Run: `go run ./examples/05-bullet-styles/main.go`
7. `06-text-enhancements` - What: Use advanced text effects like caps/highlight/superscript - Run: `go run ./examples/06-text-enhancements/main.go`
8. `07-tables` - What: Create and style basic tables - Run: `go run ./examples/07-tables/main.go`
9. `08-dynamic-report-generation` - What: Build a report deck from runtime data - Run: `python examples/08-dynamic-report-generation/main.py`
10. `08-table-cell-merge` - What: Merge table cells with row/column spans - Run: `go run ./examples/08-table-cell-merge/main.go`
11. `09-charts` - What: Generate standard chart types - Run: `go run ./examples/09-charts/chart_smoke.go`
12. `10-images` - What: Embed images from file and bytes - Run: `go run ./examples/10-images/main.go`
13. `11-image-advanced-sources` - What: Add images from bytes/base64/advanced sources - Run: `go run ./examples/11-image-advanced-sources/main.go`
14. `12-shapes` - What: Showcase shape creation and styling - Run: `go run ./examples/12-shapes/feature_showcase.go`
15. `13-connectors` - What: Draw straight/elbow/curved connectors - Run: `go run ./examples/13-connectors/main.go`
16. `14-table-api-improvements` - What: Python table API coverage including styling variants - Run: `python examples/14-table-api-improvements/main.py`; `python examples/14-table-api-improvements/styling_example.py`
17. `14-transitions` - What: Apply slide transition effects - Run: `go run ./examples/14-transitions/main.go`
18. `15-headers-footers` - What: Manage headers/footers in Python flows - Run: `python examples/15-headers-footers/main.py`
19. `15-cli-commands` - What: Reference and demo CLI command workflows - Run: `go run ./examples/15-cli-commands/main.go`
20. `16-chart-type-enum` - What: Use chart type enums in Python - Run: `python examples/16-chart-type-enum/main.py`
21. `16-templates` - What: Generate decks from reusable templates - Run: `go run ./examples/16-templates/main.go`
22. `17-custom-slide-layout-composition` - What: Compose custom slide layouts in Python - Run: `python examples/17-custom-slide-layout-composition/main.py`
23. `17-themes` - What: Apply theme and master-based styling - Run: `go run ./examples/17-themes/theme_master_smoke.go`
24. `18-theme-aware-presentation` - What: Create theme-aware presentations in Python - Run: `python examples/18-theme-aware-presentation/main.py`
25. `18-layout-helpers` - What: Use layout/unit helper utilities - Run: `go run ./examples/18-layout-helpers/main.go`
26. `19-read-modify-existing` - What: Open and modify existing PPTX files - Run: `go run ./examples/19-read-modify-existing/editor_smoke.go`; `go run ./examples/19-read-modify-existing/editor_overwrite.go`
27. `20-validation-repair` - What: Validate and repair malformed presentations - Run: `go run ./examples/20-validation-repair/main.go`
28. `21-export-html-pdf` - What: Export presentation content to HTML/PDF - Run: `go run ./examples/21-export-html-pdf/main.go`
29. `22-speaker-notes` - What: Read and write speaker notes - Run: `go run ./examples/22-speaker-notes/editor_notes_smoke.go`
30. `23-media-embed` - What: Embed audio/video media into slides - Run: `go run ./examples/23-media-embed/main.go`
31. `24-smartart` - What: Create SmartArt diagrams across layouts - Run: `go run ./examples/24-smartart/main.go`
32. `25-export-html` - What: Export slides to HTML (and PDF when enabled) - Run: `go run ./examples/25-export-html/main.go`
33. `26-vba-macros` - What: Create macro-enabled PPTM output - Run: `go run ./examples/26-vba-macros/main.go`
34. `27-custom-xml` - What: Attach and round-trip custom XML parts - Run: `go run ./examples/27-custom-xml/main.go`
35. `28-animations` - What: Add entrance/emphasis/exit animations - Run: `go run ./examples/28-animations/main.go`
36. `29-handout-master` - What: Configure handout master print layouts - Run: `go run ./examples/29-handout-master/main.go`
37. `30-embedded-fonts` - What: Embed fonts for consistent rendering - Run: `go run ./examples/30-embedded-fonts/main.go`
38. `31-hyperlinks` - What: Add external and internal hyperlinks - Run: `go run ./examples/31-hyperlinks/main.go`
39. `32-mermaid` - What: Render Mermaid diagrams into slides - Run: `go run ./examples/32-mermaid/main.go`
40. `32-mermaid-diagrams` - What: Render multiple Mermaid diagram types - Run: `go run ./examples/32-mermaid-diagrams/main.go`
41. `33-notes-master` - What: Work with notes master content - Run: `go run ./examples/33-notes-master/notes_master_smoke.go`
42. `34-urlfetch` - What: Generate slides from fetched web content - Run: `go run ./examples/34-urlfetch/main.go`
43. `34-web2ppt` - What: Convert webpage structures to PPTX slides - Run: `go run ./examples/34-web2ppt/main.go`
44. `35-layout-helpers` - What: Apply layout helper APIs in real decks - Run: `go run ./examples/35-layout-helpers/main.go`
45. `35-prelude-helpers` - What: Use presentation prelude helper patterns - Run: `go run ./examples/35-prelude-helpers/main.go`
46. `36-slide-master` - What: Use multiple slide masters/layout families - Run: `go run ./examples/36-slide-master/multi_master_smoke.go`
47. `37-slide-duplication` - What: Duplicate slides (simple and complex cases) - Run: `go run ./examples/37-slide-duplication/slide_duplication.go`; `go run ./examples/37-slide-duplication/complex_duplication.go`
48. `38-editor-image-support` - What: Edit and replace images via editor APIs - Run: `go run ./examples/38-editor-image-support/editor_image_smoke.go`
49. `39-editor-chart-support` - What: Edit chart data/formatting via editor APIs - Run: `go run ./examples/39-editor-chart-support/smoke_main.go`
50. `40-editor-notes-support` - What: Set and update notes using editor APIs - Run: `go run ./examples/40-editor-notes-support/main.go`
51. `40-metadata-writer` - What: Write and update presentation metadata - Run: `go run ./examples/40-metadata-writer/main.go`
52. `41-deep-shape-editing` - What: Perform deep shape-level mutations - Run: `go run ./examples/41-deep-shape-editing/main.go`
53. `42-smart-merge-assets` - What: Merge decks/assets with deduplication behavior - Run: `go run ./examples/42-smart-merge-assets/main.go`
54. `43-advanced-theme-management` - What: Advanced theme swapping and preset workflows - Run: `go run ./examples/43-advanced-theme-management/main.go`
55. `43-presentation-props-editor` - What: Edit presentation/core properties - Run: `go run ./examples/43-presentation-props-editor/main.go`
56. `44-section-management` - What: Create, rename, and manage slide sections - Run: `go run ./examples/44-section-management/smoke_main.go`
57. `45-commenting-api` - What: Add and manage comments - Run: `go run ./examples/45-commenting-api/comments_basic.go`
58. `46-presentation-protection` - What: Apply document/presentation protection options - Run: `go run ./examples/46-presentation-protection/main.go`
59. `47-morph-transitions` - What: Demonstrate morph-style transitions - Run: `go run ./examples/47-morph-transitions/main.go`
60. `48-accessibility-alt-text` - What: Set alt text and decorative flags for accessibility - Run: `go run ./examples/48-accessibility-alt-text/main.go`
61. `49-advanced-hyperlinks` - What: Use advanced hyperlink/action patterns - Run: `go run ./examples/49-advanced-hyperlinks/main.go`
62. `50-interactive-elements` - What: Build clickable interactive slide elements - Run: `go run ./examples/50-interactive-elements/main.go`
63. `51-language-localization` - What: Set language tags and localization metadata - Run: `go run ./examples/51-language-localization/main.go`
64. `51-localization` - What: Localization helpers and language-specific content - Run: `go run ./examples/51-localization/main.go`
65. `52-legacy-interop` - What: Interop flow for legacy formats and conversion - Run: `go run ./examples/52-legacy-interop/main.go`
66. `52-legacy-ppt-interop` - What: Legacy .ppt compatibility workflow - Run: `go run ./examples/52-legacy-ppt-interop/main.go`
67. `53-slide-properties` - What: Read/write per-slide properties - Run: `go run ./examples/53-slide-properties/slide_props_smoke.go`
68. `54-theme-master` - What: Theme master API and defaults - Run: `go run ./examples/54-theme-master/main.go`
69. `55-background-fills` - What: Solid/gradient/picture background fills - Run: `go run ./examples/55-background-fills/background_smoke.go`
70. `56-action-api` - What: Shape/text action API examples - Run: `go run ./examples/56-action-api/action_smoke.go`; `go run ./examples/56-action-api/action_api_smoke.go`
71. `57-bridge-command-api` - What: Execute editor actions via command bridge - Run: `go run ./examples/57-bridge-command-api/main.go`
72. `57-placeholder-overrides` - What: Override placeholder text/content behavior - Run: `go run ./examples/57-placeholder-overrides/placeholder_override_smoke.go`
73. `58-bridge-performance` - What: Benchmark bridge command throughput - Run: `go run ./examples/58-bridge-performance/main.go`
74. `58-gopptx-rich-slide` - What: Compose a rich mixed-content showcase slide - Run: `go run ./examples/58-gopptx-rich-slide/main.go`
75. `60-presentation-api-metadata` - What: Presentation API metadata read/write flow - Run: `go run ./examples/60-presentation-api-metadata/main.go`
76. `61-smartart-edit-hidden` - What: Edit SmartArt and hidden-slide states - Run: `go run ./examples/61-smartart-edit-hidden/main.go`
77. `62-smartart-edit-variations` - What: Apply many SmartArt mutation variations - Run: `go run ./examples/62-smartart-edit-variations/main.go`
78. `63-presentation-api` - What: Manage presentation lifecycle (create/open/save/properties) - Run: `go run ./examples/63-presentation-api/main.go`
79. `64-slides-api` - What: Manage slide layouts, ordering, numbering, and visibility - Run: `go run ./examples/64-slides-api/main.go`
80. `65-shapes-api` - What: Build and style shapes with geometry and text options - Run: `go run ./examples/65-shapes-api/main.go`
81. `66-placeholders-api` - What: Fill typed placeholders and override placeholder content - Run: `go run ./examples/66-placeholders-api/main.go`
82. `67-table-api` - What: Build tables and control rows, cells, and styling - Run: `go run ./examples/67-table-api/main.go`
83. `68-chart-data-api` - What: Build chart data sets and run data update flows - Run: `go run ./examples/68-chart-data-api/main.go`
84. `69-chart-api` - What: Create and format chart variants - Run: `go run ./examples/69-chart-api/main.go`
85. `70-text-api` - What: Compose text runs, bullets, and notes text - Run: `go run ./examples/70-text-api/main.go`
86. `71-action-api` - What: Attach click/hover actions and hyperlink triggers - Run: `go run ./examples/71-action-api/main.go`
87. `72-dml-api` - What: Apply DrawingML fills, lines, gradients, and effects - Run: `go run ./examples/72-dml-api/main.go`
88. `73-image-api` - What: Transform images with crop/effects and accessibility metadata - Run: `go run ./examples/73-image-api/main.go`
89. `74-util-api` - What: Convert units and use utility helpers - Run: `go run ./examples/74-util-api/main.go`
90. `75-media-ole-api` - What: Embed media and OLE objects - Run: `go run ./examples/75-media-ole-api/main.go`
91. `76-notes-api` - What: Author rich notes and update notes via editor APIs - Run: `go run ./examples/76-notes-api/main.go`
92. `77-background-api` - What: Configure slide backgrounds (solid/gradient/picture) - Run: `go run ./examples/77-background-api/main.go`
93. `78-enum-api` - What: Reference and validate enum constants by category - Run: `go run ./examples/78-enum-api/main.go`
94. `79-document-infra` - What: Demonstrate document infrastructure APIs (sections/comments/reorder/z-order) - Run: `go run ./examples/79-document-infra/main.go`; `uv run python examples/79-document-infra/main.py`

## Python Examples

- Python examples index: `examples/python/README.md`
- Starter example: `uv run python examples/01-basic-pptx-generation/main.py`
- Numbered examples:
  - `uv run python examples/02-slide-layouts/main.py`
  - `uv run python examples/03-markdown-to-pptx/main.py`
  - `uv run python examples/04-text-formatting/main.py`
  - `uv run python examples/05-bullet-styles/main.py`
  - `uv run python examples/06-text-enhancements/main.py`
  - `uv run python examples/07-tables/main.py`
  - `uv run python examples/08-dynamic-report-generation/main.py`
  - `uv run python examples/09-charts/main.py`
  - `uv run python examples/10-images/main.py`
  - `uv run python examples/11-image-advanced-sources/main.py`
  - `uv run python examples/12-shapes/main.py`
  - `uv run python examples/13-connectors/main.py`
  - `uv run python examples/14-table-api-improvements/main.py`
  - `uv run python examples/15-headers-footers/main.py`
  - `uv run python examples/16-chart-type-enum/main.py`
  - `uv run python examples/17-custom-slide-layout-composition/main.py`
  - `uv run python examples/18-theme-aware-presentation/main.py`
  - `uv run python examples/19-read-modify-existing/main.py`
  - `uv run python examples/20-validation-repair/main.py`
  - `uv run python examples/21-export-html-pdf/main.py`
  - `uv run python examples/22-speaker-notes/main.py`
  - `uv run python examples/23-media-embed/main.py`
  - `uv run python examples/24-smartart/main.py`
  - `uv run python examples/25-export-html/main.py`
  - `uv run python examples/26-vba-macros/main.py`
  - `uv run python examples/27-custom-xml/main.py`
  - `uv run python examples/28-animations/main.py`
  - `uv run python examples/29-handout-master/main.py`
  - `uv run python examples/30-embedded-fonts/main.py`
  - `uv run python examples/31-hyperlinks/main.py`
  - `uv run python examples/32-mermaid/main.py`
  - `uv run python examples/33-notes-master/main.py`
  - `uv run python examples/35-layout-helpers/main.py`
  - `uv run python examples/36-slide-master/main.py`
  - `uv run python examples/37-slide-duplication/main.py`
  - `uv run python examples/38-editor-image-support/main.py`
  - `uv run python examples/39-editor-chart-support/main.py`
  - `uv run python examples/40-editor-notes-support/main.py`
- Grayscale demo: `uv run python examples/python/scripts/61_grayscale_targeted.py`
- Tests and verification: `uv run python examples/python/tests/python_bridge_smoke_test.py`

## Example Assets

Shared fixtures live under: `examples/assets/[task-number]/`
