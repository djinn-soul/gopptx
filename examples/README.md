# Examples

Task-focused and smoke examples live under this directory.

## Example Assets

Shared fixture assets used by smoke/task examples now live under task-tagged paths:

- `examples/assets/01/01_basic_pptx.pptx`
- `examples/assets/37/160070-labyrinth-template-16x9.pptx`
- `examples/assets/37/162301-moneybox-template-16x9.pptx`
- `examples/assets/55/repository-open-graph-template.png`

## Smoke Generators (Canonical)

Use `examples/smoke/*` as the preferred path for smoke/demo PPT generators:

- `go run ./examples/smoke/actions/generate_action_api_smoke`
- `go run ./examples/smoke/actions/generate_action_smoke`
- `go run ./examples/smoke/charts/generate_chart_smoke_samples`
- `go run ./examples/smoke/core/generate_animation_smoke`
- `go run ./examples/smoke/core/generate_background_smoke`
- `go run ./examples/smoke/core/generate_feature_showcase`
- `go run ./examples/smoke/core/generate_slide_props_smoke`
- `go run ./examples/smoke/core/generate_text_enhancements_smoke`
- `go run ./examples/smoke/core/generate_text_frame_smoke`
- `go run ./examples/smoke/core/generate_theme_master_smoke`
- `go run ./examples/smoke/editor/generate_complex_duplication_smoke`
- `go run ./examples/smoke/editor/generate_editor_chart_smoke`
- `go run ./examples/smoke/editor/generate_editor_image_smoke`
- `go run ./examples/smoke/editor/generate_editor_notes_smoke`
- `go run ./examples/smoke/editor/generate_editor_overwrite_smoke`
- `go run ./examples/smoke/editor/generate_editor_smoke`
- `go run ./examples/smoke/editor/generate_modular_sections_smoke`
- `go run ./examples/smoke/editor/generate_multi_template_smoke`
- `go run ./examples/smoke/editor/generate_slide_duplication_smoke`

## Output Location

All example generators write `.pptx` outputs to:

- `examples/output/`

## Task Examples

Task-specific examples are also kept here, e.g.:

- `examples/39-editor-chart-support`
- `examples/40-metadata-writer`
- `examples/41-deep-shape-editing`
- `examples/42-smart-merge-assets`
- `examples/43-presentation-props-editor`
