# Scripts Layout

This repository groups scripts by intent:

- `scripts/smoke/`: smoke/demo generators for gopptx features
- `scripts/parity/`: parity checks against `ppt-rs` fixtures/signatures
- `scripts/tasks/`: batch task sample generators

## Smoke Scripts

### Actions

- `go run ./scripts/smoke/actions/generate_action_api_smoke`
- `go run ./scripts/smoke/actions/generate_action_smoke`

### Charts

- `go run ./scripts/smoke/charts/generate_chart_smoke_samples`

### Core

- `go run ./scripts/smoke/core/generate_animation_smoke`
- `go run ./scripts/smoke/core/generate_background_smoke`
- `go run ./scripts/smoke/core/generate_feature_showcase`
- `go run ./scripts/smoke/core/generate_slide_props_smoke`
- `go run ./scripts/smoke/core/generate_text_enhancements_smoke`
- `go run ./scripts/smoke/core/generate_text_frame_smoke`
- `go run ./scripts/smoke/core/generate_theme_master_smoke`

### Editor

- `go run ./scripts/smoke/editor/generate_complex_duplication_smoke`
- `go run ./scripts/smoke/editor/generate_editor_chart_smoke`
- `go run ./scripts/smoke/editor/generate_editor_image_smoke`
- `go run ./scripts/smoke/editor/generate_editor_notes_smoke`
- `go run ./scripts/smoke/editor/generate_editor_overwrite_smoke`
- `go run ./scripts/smoke/editor/generate_editor_smoke`
- `go run ./scripts/smoke/editor/generate_modular_sections_smoke`
- `go run ./scripts/smoke/editor/generate_multi_template_smoke`
- `go run ./scripts/smoke/editor/generate_slide_duplication_smoke`

## Parity Scripts

- `go run ./scripts/parity/compare_chart_parity_with_ppt_rs`
- `go run ./scripts/parity/compare_table_parity_with_ppt_rs`

## Task Samples

- `go run ./scripts/tasks/generate_task_samples`
