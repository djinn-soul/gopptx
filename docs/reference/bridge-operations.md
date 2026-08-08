# JSON bridge operations

Every gopptx operation — from Python, from C, from any other client — is one JSON envelope in
and one out. This page lists all **179** operations the dispatcher accepts.

**Authoritative source:** `pkg/pptx/editor/opspec.go`. The Python constants in
`python/gopptx/_ops_constants.py` are generated from it, so if this page and the code disagree,
the code is right.

## Envelope

**Request**

```json
{
  "api_version": 1,
  "request_id": "b1c2…",
  "op": "add_slide",
  "payload": {"title": "Agenda"}
}
```

`request_id` is optional and echoed back; `api_version` is currently `1`.

**Success**

```json
{"ok": true, "request_id": "b1c2…", "result": {"index": 2}}
```

**Failure**

```json
{
  "ok": false,
  "request_id": "b1c2…",
  "error": {"code": "INVALID_INDEX", "message": "slide index 9 out of range"}
}
```

### Error codes

`INVALID_JSON` · `UNSUPPORTED_VERSION` · `UNKNOWN_OP` · `INVALID_PAYLOAD` · `MISSING_FIELD` ·
`INVALID_TYPE` · `INVALID_VALUE` · `INVALID_INDEX` · `INVALID_HANDLE` · `INVALID_BATCH_ITEM` ·
`OP_FAILED` · `MARSHAL_ERROR` · `INTERNAL_ERROR`

## Calling an operation

=== "Python (typed)"

    ```python
    slide = pres.add_slide("Agenda")
    ```

=== "Python (raw)"

    ```python
    result = pres.execute("add_slide", {"title": "Agenda"})
    ```

=== "Go"

    ```go
    resp := editor.ExecuteCommand(e, `{"api_version":1,"op":"add_slide","payload":{"title":"Agenda"}}`)
    ```

Every typed Python method is a wrapper over one of these operations, so the raw form is always
available for anything not yet wrapped. Payload field names are the snake_case JSON tags on the
Go request structs in `pkg/pptx/editor`.

---

## Batch and session

| Operation | Purpose |
|---|---|
| `batch_execute` | Run an ordered list of commands in one crossing |
| `slide_count` | Number of slides |
| `get_metadata` | Deck metadata summary — size, counts, properties |
| `validate` | Structural findings for the package |
| `repair` | Fix repairable structural defects |

## Slides

| Operation | Purpose |
|---|---|
| `add_slide` | Append a slide |
| `remove_slide` | Delete a slide by index |
| `move_slide` | Reorder a slide |
| `duplicate_slide` | Copy a slide |
| `duplicate_slide_after` | Copy a slide and place it after a given index |
| `update_slide` | Update a slide's content in place |
| `set_slide_title` | Set the title placeholder text |
| `set_slide_hidden` | Hide or show a slide in the slide show |
| `list_slides` | Enumerate slides with their metadata |
| `copy_slides_from` | Copy selected slides in from another deck |
| `merge_from_file` | Append an entire deck from a file |
| `merge_from_editor` | Append an entire deck from another open editor |

## Layouts and masters

| Operation | Purpose |
|---|---|
| `add_slide_master` / `remove_slide_master` | Manage slide masters |
| `add_slide_layout` / `remove_slide_layout` | Manage layouts |
| `reorder_slide_layouts` | Change layout order within a master |
| `list_slide_masters` | Enumerate masters |
| `list_slide_layouts` | Enumerate layouts |
| `list_master_layouts` | Enumerate the layouts belonging to one master |
| `get_slide_layout_ref` | The layout a slide is bound to |
| `rebind_slide_layout` | Point a slide at a different layout part |
| `clone_layout_master_family` | Copy a master with all of its layouts |
| `import_layout_from` | Import a layout from another deck |
| `get_layout_shapes` / `get_master_shapes` | Read shapes on a layout or master |
| `add_layout_shape` / `add_master_shape` | Add a shape to a layout or master |
| `add_layout_textbox` / `add_master_textbox` | Add a textbox to a layout or master |
| `get_layout_placeholders` / `get_master_placeholders` | Read placeholder definitions |
| `get_slide_show_master_shapes` / `set_slide_show_master_shapes` | Whether a slide draws master shapes |
| `set_slide_follow_master_background` | Whether a slide inherits the master background |

## Placeholders

| Operation | Purpose |
|---|---|
| `list_placeholders` | Placeholders on a slide |
| `set_placeholder_content` | Fill a placeholder with text, an image, a table or a chart |

## Shapes

| Operation | Purpose |
|---|---|
| `add_shape` | Add a preset-geometry shape |
| `add_textbox` / `add_textboxes` | Add one or many textboxes |
| `add_connector` / `add_connectors` | Add one or many connectors |
| `add_group_shape` | Add a group |
| `group_shapes` / `ungroup_shapes` | Group and ungroup existing shapes |
| `build_freeform` | Build a freeform path |
| `update_shape` | Change geometry, fill, line, rotation, text |
| `remove_shape` | Delete a shape |
| `clear_shapes` | Remove every shape from a slide |
| `list_shapes` | Enumerate shapes with ids and geometry |
| `search_shapes` | Find shapes by text or property |
| `move_shape_to_front` / `move_shape_to_back` / `move_shape_to_index` | Z-order |
| `set_shape_adjustments` | Preset-geometry adjustment handles |
| `set_picture_fill` | Fill a shape with a picture |
| `reserve_shape_ids` | Pre-allocate shape ids for bulk insertion |
| `get_effective_shape_style` | Resolved style after theme and layout inheritance |

## Text

| Operation | Purpose |
|---|---|
| `get_shape_runs` / `set_shape_runs` | Read and replace a shape's runs |
| `append_shape_run` | Append one run |
| `update_shape_run_text` | Change the text of a single run |
| `set_slide_shape_runs` | Replace runs on several shapes of one slide |
| `update_slide_run_texts` | Bulk run-text edits on one slide |
| `update_deck_run_texts` | Bulk run-text edits across slides |
| `get_shape_text_state` / `get_slide_text_states` | Read text state |
| `find_and_replace` | Deck-wide text substitution |
| `fit_shape_text` / `fit_shape_to_text` | Fit text to the shape, or the shape to its text |
| `add_equation` | Insert an OMML equation from LaTeX |

## Tables

| Operation | Purpose |
|---|---|
| `add_table` | Insert a table |
| `get_table` | Read table structure and cell contents |
| `add_table_row` / `add_table_column` | Append |
| `insert_table_row` / `insert_table_column` | Insert at a position |
| `remove_table_row` / `remove_table_column` | Delete |
| `update_table_cell` | Set cell text and formatting |
| `update_table_cell_border` | Cell border styling |
| `merge_table_cells` / `split_table_cell` | Merge and split |
| `set_table_row_height` / `set_table_column_width` | Sizing |
| `update_table_flags` | Header row, total row, banding |
| `set_table_style` | Apply a table style |
| `define_table_style` | Define a custom style |
| `list_table_styles` | Enumerate available styles |

## Charts

| Operation | Purpose |
|---|---|
| `add_chart` | Insert a chart |
| `list_slide_charts` | Enumerate charts on a slide |
| `get_chart_state` | Full chart state — series, axes, formatting |
| `get_chart_data_source` | Where the chart's data comes from |
| `update_chart_data` | Replace categories and series |
| `update_chart_data_batch` | Update several charts at once |
| `update_chart_cached_values` | Refresh the cached values without touching the workbook |
| `update_chart_formatting` | Title, legend, axes, labels, colours |
| `add_chart_user_shapes` | Overlay shapes on a chart |

## Images and media

| Operation | Purpose |
|---|---|
| `add_image` | Insert an image from path, bytes or base64 |
| `add_video` / `add_audio` | Embed media |
| `add_online_video` | Reference a hosted video |
| `add_ole_object` | Embed an OLE object |
| `get_image_metadata` | Dimensions, format, relationship id |
| `list_slide_images` / `list_slide_media` | Enumerate media on a slide |
| `swap_image_by_index` / `swap_image_by_rel_id` | Replace an image without re-laying out |
| `extract_media` | Read a media part's bytes |

## SmartArt

| Operation | Purpose |
|---|---|
| `add_smartart` | Insert a diagram from a layout and a list of nodes |
| `get_smartart` / `list_smartart` | Read diagrams |
| `update_smartart` | Replace the node text |
| `set_smartart_nodes` | Replace the whole node tree |
| `add_smartart_node` / `update_smartart_node` / `remove_smartart_node` | Node-level edits |
| `change_smartart_layout` | Switch layout, keeping the data |
| `set_smartart_style` | Quick style and colour style |
| `delete_smartart` | Remove a diagram |

## Diagrams and generated content

| Operation | Purpose |
|---|---|
| `add_mermaid_shape` | Render a Mermaid diagram onto a slide |
| `markdown_to_slides` | Convert Markdown into slides |
| `url_fetch_to_slides` | Fetch a URL and convert its content into slides |

## Notes, comments and sections

| Operation | Purpose |
|---|---|
| `set_notes` / `get_notes` | Speaker notes |
| `notes_slide_exists` | Whether a slide has a notes slide |
| `list_notes_shapes` / `list_notes_placeholders` | Notes-slide contents |
| `set_notes_shape_text` / `set_notes_shape_props` | Edit notes-slide shapes |
| `update_notes_master` | Edit the notes master |
| `get_handout_master` / `update_handout_master` | Handout master |
| `add_author` / `get_authors` | Comment authors |
| `add_comment` / `get_comments` / `remove_comment` | Comments |
| `add_section` / `remove_section` / `rename_section` / `get_sections` | Sections |

## Themes and appearance

| Operation | Purpose |
|---|---|
| `apply_theme` | Apply a theme by name |
| `set_global_theme_preset` | Apply a built-in Office preset |
| `set_theme_color_scheme` / `get_theme_color_scheme` | Colour scheme |
| `set_theme_font_scheme` | Major and minor fonts |
| `get_theme_inventory` | Themes present in the package |
| `set_slide_size` | Slide dimensions in EMU |
| `set_slide_background` / `get_slide_background` | Slide background |
| `set_slide_header_footer` / `get_slide_header_footer` | Header, footer, slide number, date |
| `convert_to_grayscale` | Convert deck, slides or shapes to grayscale |

`apply_theme` and `set_global_theme_preset` share one vocabulary: the gopptx theme names
(`Corporate`, `Modern`, `Vibrant`, `Dark`, `Nature`, `Tech`, `Carbon`, `Office`) and the Office
preset names (`office`, `office2013`, `facet`, `integral`, `ion`, `retrospect`, `slice`,
`wisp`). Matching ignores case and separators.

## Motion

| Operation | Purpose |
|---|---|
| `set_slide_transition` | Slide transition, including morph |
| `add_animation` | Entrance, emphasis or exit animation on a shape |

## Templates

| Operation | Purpose |
|---|---|
| `render_template` | Render a template with data |
| `build_simple_template` | Simple deck template |
| `build_status_template` | Status-report template |
| `build_proposal_template` | Proposal template |
| `build_training_template` | Training template |
| `build_technical_template` | Technical template |

## Document properties and protection

| Operation | Purpose |
|---|---|
| `get_core_properties` / `set_core_properties` | Title, creator, subject, keywords, dates |
| `set_mark_as_final` | Mark the deck read-only-by-convention |
| `set_modify_password` | Set the modify password |
| `is_digitally_signed` | Whether the package carries a signature |

## Package extras

| Operation | Purpose |
|---|---|
| `add_custom_xml` / `list_custom_xml` / `remove_custom_xml` | Custom XML parts |
| `add_vba` | Attach a VBA project |
| `export_pdf` | Export to PDF |
| `export_html` | Export to HTML |
| `save_flat_xml` | Serialise the package as one XML file |

---

## Notes on payloads

- Indices are **zero-based** everywhere.
- Geometry is **EMU**. See [Units](../concepts.md#units-and-geometry).
- Shape identity is the numeric `shape_id` returned when the shape was created, not its position.
- Chart identity can be either `chart_index` (position on the slide) or `rel_id` (relationship
  id); the `*_by_index` and `*_by_rel_id` Python wrappers correspond to these.
- Colours are RGB hex strings without `#`, e.g. `"1F4E78"`.
- Results are JSON from Go structs, so key casing follows the struct. The Python layer adds
  snake_case aliases while keeping the originals.

For the exact fields of any payload, read the request struct in `pkg/pptx/editor` — the JSON tags
are the field names — or the `.pyi` stubs alongside the Python wrapper.
