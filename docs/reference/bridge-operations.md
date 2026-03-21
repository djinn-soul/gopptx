# JSON Bridge Operations

This page groups operation identifiers accepted by the bridge command dispatcher.

Primary source:

- `pkg/pptx/editor/opspec.go`
- `python/gopptx/ops.py`

## Envelope Format

Request:

```json
{
  "api_version": 1,
  "request_id": "uuid",
  "op": "add_slide",
  "payload": {"title": "Agenda"}
}
```

Response:

```json
{
  "ok": true,
  "request_id": "uuid",
  "result": {"index": 1}
}
```

## Slide, Layout, and Theme Ops

- `batch_execute`
- `slide_count`
- `add_slide`
- `remove_slide`
- `move_slide`
- `duplicate_slide`
- `duplicate_slide_after`
- `list_slides`
- `update_slide`
- `set_slide_title`
- `merge_from_file`
- `merge_from_editor`
- `get_metadata`
- `get_core_properties`
- `set_core_properties`
- `set_slide_size`
- `apply_theme`
- `set_global_theme_preset`
- `set_theme_color_scheme`
- `set_theme_font_scheme`
- `get_theme_inventory`
- `list_slide_layouts`
- `get_slide_layout_ref`
- `list_slide_masters`
- `list_master_layouts`
- `rebind_slide_layout`
- `clone_layout_master_family`
- `add_slide_master`
- `remove_slide_master`
- `add_slide_layout`
- `remove_slide_layout`
- `get_layout_shapes`
- `get_master_shapes`
- `get_layout_placeholders`
- `get_master_placeholders`

## Text, Shapes, and Placeholders Ops

- `find_and_replace`
- `search_shapes`
- `list_shapes`
- `add_shape`
- `update_shape`
- `remove_shape`
- `add_textbox`
- `add_textboxes`
- `add_connector`
- `add_connectors`
- `reserve_shape_ids`
- `add_group_shape`
- `group_shapes`
- `ungroup_shapes`
- `build_freeform`
- `move_shape_to_front`
- `move_shape_to_back`
- `move_shape_to_index`
- `get_slide_text_states`
- `get_shape_text_state`
- `get_shape_runs`
- `set_shape_runs`
- `set_slide_shape_runs`
- `update_deck_run_texts`
- `update_slide_run_texts`
- `update_shape_run_text`
- `append_shape_run`
- `list_placeholders`
- `set_placeholder_content`

## Table Ops

- `add_table`
- `get_table`
- `set_table_style`
- `define_table_style`
- `list_table_styles`
- `update_table_flags`
- `update_table_cell`
- `merge_table_cells`
- `split_table_cell`
- `set_table_row_height`
- `set_table_column_width`

## Chart Ops

- `add_chart`
- `list_slide_charts`
- `get_chart_state`
- `update_chart_data`
- `update_chart_data_batch`
- `update_chart_formatting`

## Notes, Comments, and Sections Ops

Sections:

- `get_sections`
- `add_section`
- `remove_section`
- `rename_section`

Comments:

- `get_authors`
- `add_author`
- `get_comments`
- `add_comment`
- `remove_comment`

Notes and handout:

- `get_notes`
- `notes_slide_exists`
- `set_notes`
- `set_notes_shape_text`
- `set_notes_shape_props`
- `list_notes_shapes`
- `list_notes_placeholders`
- `update_notes_master`
- `get_handout_master`
- `update_handout_master`

## Media, Smart Content, and Export Ops

Media:

- `add_image`
- `get_image_metadata`
- `list_slide_images`
- `swap_image_by_index`
- `swap_image_by_rel_id`
- `add_video`
- `add_audio`
- `add_ole_object`

Smart content:

- `markdown_to_slides`
- `url_fetch_to_slides`
- `add_mermaid_shape`
- `add_smartart`
- `update_smartart`
- `add_animation`
- `set_slide_transition`
- `set_slide_background`
- `set_slide_header_footer`

Export:

- `export_pdf`
- `export_html`

## Validation, Security, and Package Ops

- `validate`
- `repair`
- `set_modify_password`
- `set_mark_as_final`
- `has_digital_signature`
- `add_custom_xml`
- `list_custom_xml`
- `remove_custom_xml`
- `add_vba`
