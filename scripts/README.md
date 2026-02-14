# Scripts Layout

This repository groups scripts by intent:

- `scripts/smoke/`: smoke/demo generators for gopptx features
- `scripts/parity/`: parity checks against `ppt-rs` fixtures/signatures
- `scripts/tasks/`: batch task sample generators

## Smoke Scripts

Smoke generators are now in `examples/smoke/*` (canonical).
Under `scripts/smoke/`, only the validator remains.

### Actions

- `go run ./examples/smoke/actions/generate_action_api_smoke`
- `go run ./examples/smoke/actions/generate_action_smoke`

### Charts

- `go run ./examples/smoke/charts/generate_chart_smoke_samples`

### Core

- `go run ./examples/smoke/core/generate_animation_smoke`
- `go run ./examples/smoke/core/generate_background_smoke`
- `go run ./examples/smoke/core/generate_feature_showcase`
- `go run ./examples/smoke/core/generate_slide_props_smoke`
- `go run ./examples/smoke/core/generate_text_enhancements_smoke`
- `go run ./examples/smoke/core/generate_text_frame_smoke`
- `go run ./examples/smoke/core/generate_theme_master_smoke`

### Editor

- `go run ./examples/smoke/editor/generate_complex_duplication_smoke`
- `go run ./examples/smoke/editor/generate_editor_chart_smoke`
- `go run ./examples/smoke/editor/generate_editor_image_smoke`
- `go run ./examples/smoke/editor/generate_editor_notes_smoke`
- `go run ./examples/smoke/editor/generate_editor_overwrite_smoke`
- `go run ./examples/smoke/editor/generate_editor_smoke`
- `go run ./examples/smoke/editor/generate_modular_sections_smoke`
- `go run ./examples/smoke/editor/generate_multi_template_smoke`
- `go run ./examples/smoke/editor/generate_slide_duplication_smoke`

### Validation

- `go run ./scripts/smoke/validate_smoke_outputs` (scan `smoke_samples/` for `.pptx` and validate package structure)
- `go run ./scripts/smoke/validate_smoke_outputs -file examples/assets/01/01_basic_pptx.pptx`
- `go run ./scripts/smoke/validate_multi_master_against_powerpoint -baseline examples/output/pp_multi_master_reference.pptx -candidate examples/output/36_multi_master_smoke.pptx`

## Parity Scripts

- `go run ./scripts/parity/compare_chart_parity_with_ppt_rs`
- `go run ./scripts/parity/compare_table_parity_with_ppt_rs`

## Task Samples

- `go run ./scripts/tasks/generate_task_samples`

## Python Bindings

- Windows (PowerShell): `./scripts/build_python.ps1`
- Linux/macOS (bash): `./scripts/build_python.sh`
