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

- `add_slide`
- `add_slide_layout`
- `add_slide_master`
- `apply_theme`
- `batch_execute`
- `clone_layout_master_family`
- `duplicate_slide`
- `duplicate_slide_after`
- `get_core_properties`
- `get_layout_placeholders`
- `get_layout_shapes`
- `get_master_placeholders`
- `get_master_shapes`
- `get_metadata`
- `get_slide_layout_ref`
- `get_theme_inventory`
- `list_master_layouts`
- `list_slide_layouts`
- `list_slide_masters`
- `list_slides`
- `merge_from_editor`
- `merge_from_file`
- `move_slide`
- `rebind_slide_layout`
- `remove_slide`
- `remove_slide_layout`
- `remove_slide_master`
- `set_core_properties`
- `set_global_theme_preset`
- `set_slide_hidden`
- `set_slide_size`
- `set_slide_title`
- `set_theme_color_scheme`
- `set_theme_font_scheme`
- `slide_count`
- `update_slide`

## Text, Shapes, and Placeholders Ops

- `add_connector`
- `add_connectors`
- `add_group_shape`
- `add_shape`
- `add_textbox`
- `add_textboxes`
- `append_shape_run`
- `build_freeform`
- `clear_shapes`
- `find_and_replace`
- `get_shape_runs`
- `get_shape_text_state`
- `get_slide_text_states`
- `group_shapes`
- `list_placeholders`
- `list_shapes`
- `move_shape_to_back`
- `move_shape_to_front`
- `move_shape_to_index`
- `remove_shape`
- `reserve_shape_ids`
- `search_shapes`
- `set_placeholder_content`
- `set_shape_runs`
- `set_slide_shape_runs`
- `ungroup_shapes`
- `update_deck_run_texts`
- `update_shape`
- `update_shape_run_text`
- `update_slide_run_texts`

## Table Ops

- `add_table`
- `add_table_column`
- `add_table_row`
- `define_table_style`
- `get_table`
- `insert_table_column`
- `insert_table_row`
- `list_table_styles`
- `merge_table_cells`
- `remove_table_column`
- `remove_table_row`
- `set_table_column_width`
- `set_table_row_height`
- `set_table_style`
- `split_table_cell`
- `update_table_cell`
- `update_table_cell_border`
- `update_table_flags`

## Chart Ops

- `add_chart`
- `get_chart_state`
- `list_slide_charts`
- `update_chart_data`
- `update_chart_data_batch`
- `update_chart_formatting`

## Notes, Comments, and Sections Ops

Sections:

- `add_section`
- `get_sections`
- `remove_section`
- `rename_section`

Comments:

- `add_author`
- `add_comment`
- `get_authors`
- `get_comments`
- `remove_comment`

Notes and handout:

- `get_handout_master`
- `get_notes`
- `list_notes_placeholders`
- `list_notes_shapes`
- `notes_slide_exists`
- `set_notes`
- `set_notes_shape_props`
- `set_notes_shape_text`
- `update_handout_master`
- `update_notes_master`

## Media, Smart Content, and Export Ops

Media:

- `add_audio`
- `add_image`
- `add_ole_object`
- `add_video`
- `get_image_metadata`
- `list_slide_images`
- `swap_image_by_index`
- `swap_image_by_rel_id`

Smart content:

- `add_animation`
- `add_mermaid_shape`
- `add_smartart`
- `change_smartart_layout`
- `delete_smartart`
- `get_slide_header_footer`
- `markdown_to_slides`
- `set_slide_background`
- `set_slide_header_footer`
- `set_slide_transition`
- `set_smartart_nodes`
- `set_smartart_style`
- `update_smartart`
- `url_fetch_to_slides`

Export:

- `export_html`
- `export_pdf`

## Template Ops

- `build_proposal_template`
- `build_simple_template`
- `build_status_template`
- `build_technical_template`
- `build_training_template`
- `render_template`

## Validation, Security, and Package Ops

- `add_custom_xml`
- `add_vba`
- `convert_to_grayscale`
- `is_digitally_signed`
- `list_custom_xml`
- `remove_custom_xml`
- `repair`
- `set_mark_as_final`
- `set_modify_password`
- `validate`
