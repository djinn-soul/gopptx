# Python `Presentation` API

Every public member of `gopptx.Presentation` — **209** methods and properties — grouped by
subject. Generated from the live class, so it matches the installed package rather than an
older intention.

For task-oriented usage read the [Python library guide](../guides/python-library.md) first;
this page is for looking up a specific call.

!!! warning "Coordinates are EMU"

    Any argument named `left`, `top`, `width`, `height`, `bounds`, `x`, `y` or `*_emu` is in
    English Metric Units — 914 400 per inch. Use `Inches()`, `Point()` or `Emu()` from
    `gopptx`. See [Units](../concepts.md#units-and-geometry).

## Conventions

| | |
|---|---|
| Indices | Zero-based, everywhere |
| `slide_index` | Position of the slide in the deck |
| `shape_id` | The numeric id returned when the shape was created — not its position |
| Chart identity | `chart_index` (position) or `rel_id` (relationship id); most chart calls have both variants |
| Colours | RGB hex without `#`, e.g. `"1F4E78"` |
| Result keys | Responses carry both the original Go key and a snake_case alias |
| `self` | Omitted from the signatures below |

## The object layer

Beyond `Presentation`, the navigable objects carry their own methods:

| Class | Members | Reach it via |
|---|---|---|
| `Slide` | 90 | `pres.slides[i]` |
| `Chart` | 47 | `slide.chart(i)` / `slide.charts` |
| `Table` | 23 | `slide.get_table(shape_id)` |
| `ShapeBuilder` | 26 | `from gopptx import ShapeBuilder` |
| `RunBuilder` | 18 | `from gopptx import RunBuilder` |
| `PresentationBuilder` | 13 | `from gopptx import PresentationBuilder` |

---

## `Presentation` members

### Lifecycle and session

| Method | Description |
|---|---|
| `abort_batch()` | Abort the current batch operation. |
| `batch(*, stop_on_error: bool = False)` | Context manager for buffered mutating operations. |
| `batch_active` *(property)* | Return whether the presentation is currently buffering a batch. |
| `begin_batch(*, stop_on_error: bool = False)` | Begin buffering operations for a batch execute. |
| `close()` | Discard queued textbox state when closing the deck. |
| `end_batch()` | End the batch operation and execute queued commands. |
| `execute(op: str, payload: dict[str, object] \| None = None)` | Flush queued textbox inserts before incompatible bridge operations. |
| `execute_batch(commands: list[dict[str, object]], *, stop_on_error: bool = False)` | Execute multiple bridge commands in one boundary crossing. |
| `flush_all_pending_textbox_adds()` | Flush queued textboxes for every slide. |
| `flush_pending_textbox_adds(slide_index: int)` | Flush queued textboxes for one slide using the bulk bridge op. |
| `from_template(path: str, context: dict[str, object], *, undefined: str = keep)` | Open a .pptx file as a Jinja2 template and render tags with context data. |
| `handle` *(property)* | The internal handle to the Go engine deck. |
| `invalidate_cache()` | Clear all cached data for the presentation. |
| `metadata` *(property)* | The presentation metadata. |
| `new(title: str)` | Create a new presentation with the given title. |
| `open(file_like_or_path: str \| os.PathLike[str] \| bytes \| object)` | Discard pending textbox state before opening another deck. |
| `open_bytes(data: bytes)` | Discard pending textbox state before opening a deck from bytes. |
| `open_deck(file_like_or_path: str \| os.PathLike[str] \| bytes \| object)` | Open a presentation from a path or file-like object. |
| `pending_textbox_adds(slide_index: int)` | Return whether the slide has queued textbox inserts. |
| `queue_textbox_add(slide_index: int, textbox: Mapping[str, object])` | Queue a simple textbox insert and return its real reserved shape ID. |
| `render_template(context: dict[str, object])` | Render all Jinja2 template expressions in slide shapes using context values. |
| `save(path: str)` | Flush queued textboxes before the regular save pipeline runs. |
| `slide_count` *(property)* | The number of slides in the presentation. |
| `slides_metadata` *(property)* | List of metadata for all slides in the presentation. |
| `to_bytes()` | Flush queued textboxes then serialize the deck to bytes. |

### Slides

| Method | Description |
|---|---|
| `add_bullet_slide(title: str, bullets: list[str])` | Add a bullet slide to the presentation. |
| `add_paragraph_slide(title: str, paragraph: str, *, bounds: tuple[float, float, float, float] \| None = None, layout: str \| None = None)` | Add a slide with one paragraph textbox using sensible default bounds. |
| `add_slide(title: str = , layout: str \| SlideLayout \| None = None, bullets: list[str] \| None = None, index: int \| None = None)` | Add a new slide to the presentation, optionally inserting at index. |
| `add_slide_from_markdown(markdown: str, *, layout: str = )` | Append slides generated from a Markdown string. |
| `add_slide_from_url(url: str, *, layout: str = )` | Fetch a web page and append slides generated from its content. |
| `add_title_slide(title: str)` | Add a title slide to the presentation. |
| `copy_slide_from(source: PresentationProtocol \| str \| os.PathLike[str], slide_index: int)` | Copy one slide from another presentation. |
| `copy_slides_from(source: PresentationProtocol \| str \| os.PathLike[str], slide_indices: Sequence[int] \| int \| None = None)` | Copy selected slides from another presentation. |
| `duplicate_slide(index: int, insert_at: int \| None = None)` | Duplicate a slide and return the new slide index. |
| `duplicate_slide_after(index: int)` | Duplicate slide at *index* and insert it immediately after the original. |
| `insert_slide(index: int, layout: str \| None = None, title: str = , bullets: list[str] \| None = None)` | Insert a slide at index. |
| `merge_from_editor(other: PresentationProtocol)` | Merge all slides from *other* into this presentation. |
| `merge_from_file(path: str)` | Merge slides from another presentation file. |
| `move_slide(from_index: int, to_index: int)` | Move a slide to a new position. |
| `remove_slide(index: int)` | Remove a slide from the presentation. |
| `set_slide_hidden(index: int, *, hidden: bool = True)` | Mark or unmark a slide as hidden. |
| `set_slide_title(index: int, title: str)` | Set the title of a slide. |
| `slide_index_for_id(slide_id: int)` | Resolve a slide index by stable slide ID using a cached lookup map. |
| `slide_layouts` *(property)* | Return slide layouts of the primary slide master. |
| `slide_master` *(property)* | Return the primary slide master. |
| `slide_masters` *(property)* | Get the slide masters collection. |
| `slides` *(property)* | List of cached slide proxies for all slides in the presentation. |
| `update_slide(index: int, title: str \| None = None, layout: str \| None = None, bullets: list[str] \| None = None)` | Update slide properties. |

### Layouts, masters, placeholders and themes

| Method | Description |
|---|---|
| `add_layout_shape(layout_part: str, shape_type: str, bounds: tuple[float, float, float, float])` | Add an autoshape to a slide layout and return its shape id. |
| `add_layout_textbox(layout_part: str, text: str, bounds: tuple[float, float, float, float])` | Add a text box to a slide layout and return its shape id. |
| `add_master_shape(master_part: str, shape_type: str, bounds: tuple[float, float, float, float])` | Add an autoshape to a slide master and return its shape id. |
| `add_master_textbox(master_part: str, text: str, bounds: tuple[float, float, float, float])` | Add a text box to a slide master and return its shape id. |
| `add_slide_layout(master_part: str, layout_name: str = Custom Layout)` | Add a slide layout under a slide master and return the layout part path. |
| `add_slide_master()` | Add a new slide master and return its part path. |
| `apply_theme(self: _PresentationThemeOps, theme: Theme \| str)` | Apply a theme to the presentation. |
| `clone_layout_master_family(layout_part: str)` | Clone a layout and its master family. |
| `get_layout_placeholders(layout_part: str)` | Return placeholder metadata for a slide layout. |
| `get_layout_shapes(layout_part: str)` | Return the shape names defined in a slide layout. |
| `get_master_placeholders(master_part: str)` | Return placeholder metadata for a slide master. |
| `get_master_shapes(master_part: str)` | Return the shape names defined in a slide master. |
| `get_slide_layout_ref(slide_index: int)` | Return (layout_part, master_part) for a slide. |
| `get_theme_inventory()` | Return all theme parts and master/theme bindings in the package. |
| `list_placeholders(slide_index: int)` | Bridge op: list all placeholders on a slide. |
| `list_slide_layouts()` | List all available slide layouts. |
| `rebind_slide_layout(slide_index: int, layout_part: str)` | Rebind a slide to a different layout. |
| `remove_slide_layout(layout_part: str)` | Remove a slide layout by part path. |
| `remove_slide_master(master_part: str)` | Remove a slide master by part path. |
| `reorder_slide_layouts(layout_parts: list[str], master_part: str \| None = None)` | Reorder slide layouts within a slide master. |
| `set_global_theme_preset(name: str)` | Apply a named built-in theme preset (e.g. 'facet', 'ion', 'office'). |
| `set_placeholder_content(slide_index: int, ph_index: int, ph_type: str = , **kwargs: object)` | Bridge op: insert rich content into a placeholder. |
| `set_slide_size(width: int, height: int)` | Set the slide size. |
| `set_theme_color_scheme(**colors: str)` | Update one or more standard theme color slots. |
| `set_theme_font_scheme(major: str, minor: str)` | Update major/minor latin typefaces across all theme parts. |
| `slide_height` *(property)* | Return slide height in EMU. |
| `slide_width` *(property)* | Return slide width in EMU. |

### Shapes

| Method | Description |
|---|---|
| `add_connector(slide_index: int, connector_type: ConnectorType, *points: float, **kwargs: object)` | Add a connector-like shape to a slide. |
| `add_connectors(slide_index: int, connectors: list[Mapping[str, object]])` | Add multiple connectors to one slide in a single bridge call. |
| `add_group_shape(slide_index: int, shapes: list[int] \| None = None)` | Add a group shape to a slide. |
| `add_shape(slide_index: int, shape_type: ShapeType, bounds: tuple[float, float, float, float], **kwargs: object)` | Add a shape to a slide. |
| `add_textbox(slide_index: int, *bounds: float, text: str = , **kwargs: object)` | Add a textbox-like shape to a slide. |
| `add_textboxes(slide_index: int, textboxes: list[Mapping[str, float \| str]])` | Add multiple textboxes to one slide in a single bridge call. |
| `build_freeform(slide_index: int, start_x: float = 0, start_y: float = 0, scale: tuple[float, float] \| float = 1.0)` | Create a freeform builder for this slide. |
| `clear_shapes(slide_index: int)` | Remove all shapes from a slide and return count removed. |
| `commit_freeform(slide_index: int, points: list[tuple[float, float]], *, close: bool, options: dict[str, object] \| None = None)` | Create a freeform shape from prepared points. |
| `get_shape_center(slide_index: int, shape_id: int)` | Return the absolute center (cx, cy) of a shape in EMU. |
| `group_shapes(slide_index: int, shape_ids: list[int])` | Group multiple shapes on a slide into a group shape. |
| `list_shapes(slide_index: int)` | List all shapes on a slide. |
| `move_shape_to_back(slide_index: int, shape_id: int)` | Move a shape to the back of the z-order. |
| `move_shape_to_front(slide_index: int, shape_id: int)` | Move a shape to the front of the z-order. |
| `move_shape_to_index(slide_index: int, shape_id: int, target_index: int)` | Move a shape to a specific z-index within a slide. |
| `remove_shape(slide_index: int, shape_id: int)` | Remove a shape from a slide. |
| `search_shapes(query: ShapeSearchQuery \| str)` | Search for shapes matching a query. |
| `set_shape_adjustments(slide_index: int, shape_id: int, adjustments: list[ShapeAdjustmentValue])` | Set preset-geometry adjustment values on a shape. |
| `ungroup_shapes(slide_index: int, shape_id: int)` | Ungroup a group shape, returning the ID of the first member shape. |
| `update_shape(slide_index: int, shape_id: int, updates: ShapeUpdate)` | Update shape properties. |

### Text

| Method | Description |
|---|---|
| `append_shape_run(slide_index: int, shape_id: int, run: TextRun)` | Append a run to a shape. |
| `find_and_replace(find_text: str, replace_text: str, scope: str = slides)` | Replace exact text matches in the parts named by ``scope``. |
| `flush_all_pending_shape_runs_replacements()` | Flush buffered full run replacements for all slides. |
| `flush_all_pending_slide_run_text_updates()` | Flush buffered run-text updates for all slides. |
| `flush_pending_shape_runs_replacements(*, slide_index: int, shape_id: int \| None = None)` | Flush buffered full run replacements for one slide or shape. |
| `flush_pending_slide_run_text_updates(slide_index: int)` | Flush buffered run-text updates for one slide. |
| `get_effective_shape_style(slide_index: int, shape_id: int)` | Resolve how a shape actually looks, following the inheritance chain. |
| `get_shape_runs(slide_index: int, shape_id: int)` | Get text runs for a shape. |
| `get_shape_text_state(slide_index: int, shape_id: int)` | Get text/runs/text-frame/paragraph state for a shape. |
| `get_slide_text_states(slide_index: int)` | Get text/runs/text-frame/paragraph state for all shapes on a slide. |
| `pending_shape_run_replacement(slide_index: int, shape_id: int)` | Return whether one shape has a pending full run replacement. |
| `pending_slide_run_text_updates(slide_index: int)` | Return whether a slide has buffered run-text updates. |
| `queue_shape_run_text_update(slide_index: int, shape_id: int, run_index: int, text: str)` | Buffer a run-text update until a flush boundary is reached. |
| `queue_shape_runs_replace(slide_index: int, shape_id: int, runs: list[dict[str, object]])` | Buffer full run replacement for one shape until a flush boundary. |
| `set_shape_runs(slide_index: int, shape_id: int, runs: list[TextRun])` | Replace all text runs on a shape. |
| `update_deck_run_texts(slide_updates: list[dict[str, object]])` | Apply run-text updates across multiple slides in one request. |
| `update_shape_run_text(slide_index: int, shape_id: int, run_index: int, text: str)` | Update text for one run by run index. |
| `update_slide_run_texts(slide_index: int, updates: list[dict[str, object]])` | Apply run-text updates for all targeted shapes on one slide. |

### Tables

| Method | Description |
|---|---|
| `add_table(slide: int \| None = None, slide_index: int \| None = None, rows: int \| None = None, cols: int \| None = None, **kwargs: object)` | Add a table shape to a slide and return its shape ID. |
| `add_table_from_dicts(slide: int, rows: list[dict[str, str]], bounds: tuple[int, int, int, int] \| None = None, *, column_names: list[str] \| None = None, first_row: bool = True, band_row: bool = True, **kwargs: object)` | Create a table from a list of dictionaries. |
| `add_table_from_rows(slide: int, rows: list[list[str]], bounds: tuple[int, int, int, int] \| None = None, *, first_row: bool = True, band_row: bool = True, column_widths: list[int] \| None = None, **kwargs: object)` | Create a table from a list of row data. |
| `define_table_style(name: str, style_id: str \| None = None)` | Define a custom table style and return its resolved style ID. |
| `get_all_table_style_names()` | Get all available table style names in the presentation. |
| `get_table(slide_index: int, shape_id: int)` | Return serialized table information for a table shape. |
| `get_table_cell(slide_index: int, shape_id: int, row: int, col: int)` | Return one table cell payload by zero-based row and column. |
| `get_table_style_by_name(name: str)` | Find a presentation table style GUID by name. |
| `list_table_styles()` | List available table styles visible to the presentation. |
| `merge_table_cells(slide_index: int, shape_id: int, cell_range: tuple[int, int, int, int])` | Merge a rectangular range of table cells. |
| `set_table_cell_text(slide_index: int, shape_id: int, row: int, col: int, text: str)` | Update the text value for one table cell. |
| `set_table_column_width(slide_index: int, shape_id: int, col: int, width: int)` | Set the width of a specific table column. |
| `set_table_flags(slide_index: int, shape_id: int, flags: dict[str, bool])` | Set table display flags such as header-row or banded options. |
| `set_table_row_height(slide_index: int, shape_id: int, row: int, height: int)` | Set the height of a specific table row. |
| `set_table_style(slide_index: int, shape_id: int, style: str)` | Apply a table style by name or GUID. |
| `split_table_cell(slide_index: int, shape_id: int, row: int, col: int)` | Split a merged table cell back into its original cells. |

### Charts

| Method | Description |
|---|---|
| `add_chart(slide_index: int, chart_type: ChartType, categories: Sequence[str] \| CategoryChartData \| XyChartData, values_or_series: Sequence[float] \| Sequence[dict[str, object]] \| None = None, *, title: str = Chart, bounds: tuple[float, float, float, float] = (0, 0, 0, 0))` | Add a chart to a slide. |
| `add_combo_chart(slide_index: int, categories: list[str], bar_series: list[dict[str, object]], line_series: list[dict[str, object]], *, title: str = Chart, bounds: tuple[float, float, float, float] = (0, 0, 0, 0))` | Add a combo (bar + line) chart to a slide. |
| `get_chart_data_source(slide_index: int, chart_selector: ChartSelector)` | Report whether a chart's data is embedded, externally linked, or absent. |
| `get_chart_state(slide_index: int, chart_selector: ChartSelector)` | Return chart state selected by a chart selector on a slide. |
| `get_chart_state_by_index(slide_index: int, chart_index: int)` | Return chart state by zero-based chart index on a slide. |
| `get_chart_state_by_rel_id(slide_index: int, rel_id: str)` | Return chart state by relationship id on a slide. |
| `list_slide_charts(slide_index: int)` | List all charts on a slide. |
| `replace_chart_data_by_index(slide_index: int, chart_index: int, categories: list[str], values: list[float])` | Replace category/value chart data by slide-local chart index. |
| `replace_chart_data_by_rel_id(slide_index: int, rel_id: str, categories: list[str], values: list[float])` | Replace category/value chart data by chart relationship id. |
| `update_chart_cached_values(slide_index: int, chart_selector: ChartSelector, data: ChartDataUpdate)` | Refresh the numbers a chart displays without touching its link. |
| `update_chart_data(slide_index: int, chart_selector: ChartSelector \| list[str], data: ChartDataUpdate \| list[dict[str, object]])` | Update chart data for a chart on a slide. |
| `update_chart_data_batch(slide_index: int, updates: list[dict[str, object]])` | Update multiple charts on one slide in a single bridge call. |
| `update_chart_data_by_index(slide_index: int, chart_index: int, data: ChartDataUpdate)` | Update chart data by slide-local chart index. |
| `update_chart_data_by_rel_id(slide_index: int, rel_id: str, data: ChartDataUpdate)` | Update chart data by chart relationship id. |
| `update_chart_formatting(slide_index: int, chart_selector: ChartSelector, fmt: ChartFormatUpdate)` | Update chart formatting for a chart on a slide. |
| `update_chart_formatting_by_index(slide_index: int, chart_index: int, fmt: ChartFormatUpdate)` | Update chart formatting by chart index. |
| `update_chart_formatting_by_rel_id(slide_index: int, rel_id: str, fmt: ChartFormatUpdate)` | Update chart formatting by chart relationship id. |

### Images and media

| Method | Description |
|---|---|
| `add_audio(slide_index: int, source: str \| bytes, bounds: tuple[float, float, float, float], **kwargs: object)` | Add an audio file to a slide and return the created shape ID. |
| `add_image(slide_index: int, source: str \| bytes \| None = None, bounds: tuple[float, float, float, float] = (0, 0, 0, 0), **kwargs: object)` | Add an image to a slide and return the created shape ID. |
| `add_ole_object(slide_index: int, source: str \| bytes, bounds: tuple[float, float, float, float], **kwargs: object)` | Add an OLE object to a slide and return the created shape ID. |
| `add_online_video(slide_index: int, url: str, bounds: tuple[float, float, float, float], **kwargs: object)` | Link a slide to a hosted video and return the created shape ID. |
| `add_picture(slide_index: int, source: str \| bytes \| None = None, left: float = 0, top: float = 0, width: float = 0, height: float = 0, **kwargs: object)` | Add a picture shape to a slide (python-pptx compatible API). |
| `add_video(slide_index: int, source: str \| bytes, bounds: tuple[float, float, float, float], **kwargs: object)` | Add a video to a slide and return the created shape ID. |
| `extract_media(part_path: str)` | Return the bytes of one media part. |
| `get_image_metadata(slide_index: int, shape_id: int)` | Get dimensions and format metadata for an image shape. |
| `list_slide_images(slide_index: int)` | List all images embedded in a slide. |
| `list_slide_media(slide_index: int)` | List every media relationship on a slide: images, sounds and movies. |
| `save_media(part_path: str, destination: str \| PathLike[str])` | Write one media part out to a file, and return the bytes written. |
| `swap_image_by_index(slide_index: int, image_index: int, data: bytes, img_format: str)` | Replace an image at a given position within a slide. |
| `swap_image_by_rel_id(slide_index: int, rel_id: str, data: bytes, img_format: str)` | Replace an image identified by its relationship ID. |

### Notes, comments and sections

| Method | Description |
|---|---|
| `add_author(name: str, initials: str)` | Add a comment author to the presentation. |
| `add_comment(slide_index: int, author_id: int, text: str, x: int = 0, y: int = 0)` | Add a comment to a slide. |
| `add_section(name: str, slide_indices: list[int])` | Add a section to the presentation. |
| `get_authors()` | Get all comment authors in the presentation. |
| `get_comments(slide_index: int)` | Get all comments on a slide. |
| `get_handout_master()` | Return handout master information. |
| `get_notes(slide_index: int)` | Return speaker notes plain text for a slide index. |
| `get_notes_payload(slide_index: int)` | Return raw notes payload for a slide index. |
| `get_sections()` | Get all sections in the presentation. |
| `is_digitally_signed()` | Return whether the presentation has a digital signature. |
| `list_notes_placeholders(slide_index: int)` | List all placeholders in the notes pane of a slide. |
| `list_notes_shapes(slide_index: int)` | List all shapes in the notes pane of a slide. |
| `notes_slide_exists(slide_index: int)` | Return True if a notes slide exists for the given slide index. |
| `remove_comment(slide_index_or_index: int, author_id: int \| None = None, author_index: int \| None = None)` | Remove a comment from a slide. |
| `remove_section(name: str)` | Remove a section from the presentation. |
| `rename_section(old_name: str, new_name: str)` | Rename a section. |
| `sections` *(property)* | Get all sections in the presentation. |
| `set_notes(slide_index: int, text: str)` | Set speaker notes plain text for a slide index. |
| `set_notes_shape_props(slide_index: int, shape_id: int, updates: ShapeUpdate)` | Patch style/geometry properties for one notes shape by shape ID. |
| `set_notes_shape_text(slide_index: int, shape_id: int, text: str)` | Set text for one notes shape by shape ID. |
| `update_handout_master(*, orientation: str = , slides_per_page: int = 0)` | Configure the handout master. |
| `update_notes_master(*, header: str = , footer: str = , show_date_time: bool = True, show_slide_num: bool = True)` | Configure the global notes master. |

### Document properties, validation and protection

| Method | Description |
|---|---|
| `author` *(property)* | The author/creator of the presentation (python-pptx: author). |
| `comments` *(property)* | The comments/description of the presentation (python-pptx: comments). |
| `core_properties` *(property)* | Get the core properties of the presentation. |
| `get_core_properties()` | Get the core properties of the presentation. |
| `identifier` *(property)* | The identifier of the presentation. |
| `language` *(property)* | The language of the presentation. |
| `last_printed` *(property)* | The last printed date of the presentation. |
| `repair()` | Attempt to automatically repair structural issues. |
| `set_core_properties(props: CoreProperties)` | Set the core properties of the presentation. |
| `set_mark_as_final(*, final: bool = True)` | Mark the presentation as final. |
| `set_modify_password(password: str)` | Set the modify password for the presentation. |
| `title` *(property)* | The title of the presentation. |
| `validate()` | Run structural validation and return a list of issues. |
| `version` *(property)* | The version of the presentation. |

### Colour

| Method | Description |
|---|---|
| `resolve_color(color: ColorFormat)` | Resolve a colour to hex RGB, looking up theme colours in this deck. |
| `theme_color_rgb(theme_color: ThemeColor \| str, *, brightness: float = 0.0)` | Return one theme slot as hex RGB, with an optional brightness tweak. |
| `theme_color_scheme` *(property)* | Return the theme's colour slots as ``{slot: hex}``. |

### Headers and footers

| Method | Description |
|---|---|
| `get_header_footer(slide_index: int)` | Get header/footer configuration for a specific slide. |
| `set_header_footer(footer: str = , *, show_footer: bool = False, show_slide_num: bool = False, show_date_time: bool = False, date_time_text: str = )` | Set header/footer for ALL slides in the presentation. |

### Export and conversion

| Method | Description |
|---|---|
| `convert_to_grayscale(*, slides: list[int] \| None = None, shapes: list[GrayscaleShapeRef] \| None = None, text: list[GrayscaleTextRef] \| None = None, placeholders: list[GrayscalePlaceholderRef] \| None = None, scope: GrayscaleScope \| None = None)` | Convert selected slides, shapes, placeholders, runs, images, and backgrounds to grayscale. |
| `export_html(output_path: str \| None = None, options: HTMLOptions \| None = None)` | Export the presentation to an HTML document. |
| `export_pdf(output_path: str \| None = None, options: PDFOptions \| None = None)` | Export the presentation to a PDF file. |
| `save_as_pdf(output_path: str \| None = None, options: PDFOptions \| None = None)` | Export the presentation to a PDF file. |
| `save_flat_xml(output_path: str)` | Save the presentation as a single PowerPoint XML Presentation file. |

### Package extras

| Method | Description |
|---|---|
| `add_custom_xml(content: str, root_element: str \| None = None, namespace: str \| None = None, properties: dict[str, str] \| None = None)` | Embed a custom XML part in the presentation. |
| `add_vba_project(data: bytes)` | Add a VBA project binary blob to the presentation. |
| `list_custom_xml()` | Return all custom XML parts embedded in the presentation. |
| `remove_custom_xml(index: int)` | Remove a custom XML part by its index. |

---

## Anything not listed here

Every typed method wraps one JSON bridge operation, and the raw form is always available:

```python
result = pres.execute("set_slide_hidden", {"slide_index": 2, "hidden": True})
```

All 179 operations are catalogued in [Bridge operations](bridge-operations.md).

## Regenerating this page

This page is written from introspection of the installed `gopptx.Presentation`. If you add or
rename a method, regenerate rather than hand-editing, so signatures and docstrings stay true.
