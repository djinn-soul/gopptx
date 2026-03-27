# Examples

This directory contains organized, task-focused examples and smoke tests, aligned with the project's task numbering system.

## Organized Examples

Use the following commands to run the examples:

### Core Presentation Generation

- `01-basic-pptx-generation`: `go run ./examples/01-basic-pptx-generation/basic_gen.go` -> `01_hello_world.pptx`
- `04-text-formatting` (Enhancements): `go run ./examples/04-text-formatting/text_enhancements.go` -> `04_text_enhancements.pptx`
- `04-text-formatting` (Text Frame): `go run ./examples/04-text-formatting/text_frame.go` -> `04_text_frame_smoke.pptx`
- `09-charts`: `go run ./examples/09-charts/chart_smoke.go` -> `09_charts_*.pptx`
- `12-shapes`: `go run ./examples/12-shapes/feature_showcase.go` -> `12_feature_showcase.pptx`
- `16-templates`: `go run ./examples/16-templates/main.go` -> `16_invoice_template.pptx`, `16_template_invoice.pptx`
- `17-themes`: `go run ./examples/17-themes/theme_master_smoke.go` -> `17_theme_master_smoke.pptx`
- `28-animations`: `go run ./examples/28-animations/main.go` -> `28_animations.pptx`
- `31-hyperlinks`: `go run ./examples/31-hyperlinks/main.go` -> `31_advanced_hyperlink_smoke.pptx`
- `53-slide-properties`: `go run ./examples/53-slide-properties/slide_props_smoke.go` -> `53_slide_properties.pptx`
- `55-background-fills`: `go run ./examples/55-background-fills/background_smoke.go` -> `55_background_fills.pptx`
- `57-placeholder-overrides`: `go run ./examples/57-placeholder-overrides/placeholder_override_smoke.go` -> `57_placeholder_override_smoke.pptx`
- `58-gopptx-rich-slide`: `go run ./examples/58-gopptx-rich-slide/main.go` -> `58_gopptx_rich_slide.pptx`

### Masters & Layouts

- `33-notes-master`: `go run ./examples/33-notes-master/notes_master_smoke.go` -> `33_notes_master_smoke.pptx`
- `34-urlfetch`: `go run ./examples/34-urlfetch/main.go` -> `34_urlfetch_*.pptx` (HTML to PPTX with custom CSS selectors and image embedding)
- `36-slide-master`: `go run ./examples/36-slide-master/multi_master_smoke.go` -> `36_multi_master_smoke.pptx`

### Editor & Modification

- `19-read-modify-existing` (Basic): `go run ./examples/19-read-modify-existing/editor_smoke.go` -> `19_editor_modified.pptx`
- `19-read-modify-existing` (Overwrite): `go run ./examples/19-read-modify-existing/editor_overwrite.go` -> `19_editor_overwrite.pptx`
- `22-speaker-notes`: `go run ./examples/22-speaker-notes/editor_notes_smoke.go` -> `22_editor_notes_smoke.pptx`
- `23-media-embed`: `go run ./examples/23-media-embed/main.go` -> `23_media_embed_editor.pptx`
- `37-slide-duplication` (Basic): `go run ./examples/37-slide-duplication/slide_duplication.go` -> `37_slide_duplication.pptx`
- `37-slide-duplication` (Complex): `go run ./examples/37-slide-duplication/complex_duplication.go` -> `37_complex_duplication.pptx`
- `38-editor-image-support`: `go run ./examples/38-editor-image-support/editor_image_smoke.go` -> `38_editor_image_smoke.pptx`
- `39-editor-chart-support`: `go run ./examples/39-editor-chart-support/smoke_main.go`
- `44-section-management`: `go run ./examples/44-section-management/smoke_main.go`

### Advanced APIS

- `40-metadata-writer`: `go run ./examples/40-metadata-writer/main.go`
- `41-deep-shape-editing`: `go run ./examples/41-deep-shape-editing/main.go`
- `42-smart-merge-assets`: `go run ./examples/42-smart-merge-assets/main.go`
- `43-presentation-props-editor`: `go run ./examples/43-presentation-props-editor/main.go` -> `43_presentation_props_editor.pptx`, `43_brand_reskin_theme_swap.pptx`
- `45-commenting-api`: `go run ./examples/45-commenting-api/comments_basic.go`
- `49-advanced-hyperlinks`: `go run ./examples/49-advanced-hyperlinks/main.go`
- `56-action-api` (Smoke): `go run ./examples/56-action-api/action_smoke.go` -> `56_action_smoke.pptx`
- `56-action-api` (API): `go run ./examples/56-action-api/action_api_smoke.go` -> `56_action_api_smoke.pptx`
- `60-presentation-api-metadata`: `go run ./examples/60-presentation-api-metadata/main.go`

### Python Examples

- Scripts: `examples/python/scripts/`
- Tests/verification scripts: `examples/python/tests/`
- Grayscale targeting demo: `uv run python examples/python/scripts/61_grayscale_targeted.py` -> `61_grayscale_targeted_source.pptx`, `61_grayscale_targeted_result.pptx` (prints slide indices, shape IDs, and placeholder types; uses `PlaceholderType.TITLE` / `PlaceholderType.FOOTER`)
- Run script examples via `python examples/python/scripts/<script>.py` (from root).
- Run verification scripts via `python examples/python/tests/<script>.py` (from root).

## Example Assets

Shared fixture assets live under: `examples/assets/[task-number]/`

## Output Location

All example generators write `.pptx` outputs to: `examples/output/`
