# Python Presentation API

This is the detailed reference for the Python surface of `gopptx`.
It is organized like a library manual: classes, methods, signatures, and short contracts.

Canonical signature source:

- `python/gopptx/api.pyi`

Implementation source:

- `python/gopptx/presentation/runtime.py`
- `python/gopptx/presentation/slides/slides_mixin.py`
- `python/gopptx/presentation/layout_theme_mixin.py`
- `python/gopptx/presentation/shapes/*`
- `python/gopptx/presentation/tables/*`
- `python/gopptx/presentation/charts/*`
- `python/gopptx/presentation/notes/*`

## `Presentation`

Main entry point for working with decks.

```python
from gopptx import Presentation
```

### Construction

#### `Presentation(path: str | None = None)`

Open an existing deck when `path` is provided, or create an empty runtime object when it is not.

#### `Presentation.new(title: str) -> Presentation`

Create a new presentation with the given title.

```python
with Presentation.new("Deck Title") as pres:
    ...
```

### Lifecycle

#### `open(path: str) -> None`

Load a deck from disk into the current presentation object.

#### `save(path: str) -> None`

Write the current deck to disk.

#### `close() -> None`

Release the runtime handle and close the presentation.

#### `__enter__()`, `__exit__(...)`

Support `with Presentation.new(...) as pres:` style usage.

### Runtime

#### `execute(op: str, payload: dict[str, object] | None = None) -> dict[str, object]`

Send one bridge operation directly to the engine.

#### `execute_batch(commands: list[BatchCommand], *, stop_on_error: bool = False) -> list[BatchItemResult]`

Send a batch of operations and collect structured results.

#### `batch(stop_on_error: bool = False)`

Context manager for grouped bridge operations.

#### `invalidate_cache() -> None`

Refresh cached presentation state after low-level mutations.

### Properties

- `slide_count`
- `slides`
- `slides_metadata`
- `metadata`
- `sections`
- `slide_masters`
- `core_properties`
- `title`
- `author`
- `comments`
- `identifier`
- `language`
- `last_printed`
- `version`

## Slide management

### `add_slide(title: str, layout: str | None = None, bullets: list[str] | None = None) -> Slide`

Create a slide with an optional layout and bullet list.

### `remove_slide(index: int) -> None`

Remove a slide by position.

### `move_slide(from_index: int, to_index: int) -> None`

Reorder slides within the deck.

### `duplicate_slide(index: int, insert_at: int | None = None) -> int`

Clone an existing slide and return the new index.

### `duplicate_slide_after(index: int) -> int`

Clone a slide and insert the copy immediately after it.

### `update_slide(index: int, title: str | None = None, layout: str | None = None, bullets: list[str] | None = None) -> None`

Update title, layout, or bullet content in one call.

### `set_slide_title(index: int, title: str) -> None`

Update only the slide title.

### Convenience builders

#### `add_title_slide(title: str) -> Slide`

Create a title-only slide.

#### `add_bullet_slide(title: str, bullets: list[str]) -> Slide`

Create a title + bullet slide.

#### `add_paragraph_slide(title: str, paragraph: str, *, bounds: tuple[float, float, float, float] | None = None, layout: str | None = None) -> Slide`

Create a paragraph slide with optional bounds and layout.

### Import helpers

#### `add_slide_from_markdown(markdown: str, *, layout: str = "") -> int`

Parse Markdown and generate a slide.

#### `add_slide_from_url(url: str, *, layout: str = "") -> int`

Fetch a URL and render the extracted content into slides.

### Merging

#### `merge_from_file(path: str) -> None`

Merge content from another deck file.

#### `merge_from_editor(other: Presentation) -> None`

Merge slides from another live presentation object.

## Layout and theme

### `apply_theme(theme_name: str) -> None`

Apply a named theme to the presentation.

### `set_slide_size(width: int, height: int) -> None`

Set the slide canvas size in pixels or engine units used by the API.

### `list_slide_layouts() -> list[SlideLayoutInfo]`

Enumerate available layouts.

### `rebind_slide_layout(slide_index: int, layout_part: str) -> None`

Point a slide at a different layout part.

### `clone_layout_master_family(layout_part: str) -> SlideMasterCloneResult`

Clone a layout and its related master family.

### Theme customization

- `set_global_theme_preset(name: str) -> None`
- `set_theme_font_scheme(major: str, minor: str) -> None`
- `set_theme_color_scheme(**colors: str) -> None`
- `get_theme_inventory() -> dict[str, object]`

### Layout and master inspection

- `get_layout_shapes(layout_part: str) -> list[str]`
- `get_master_shapes(master_part: str) -> list[str]`
- `get_layout_placeholders(layout_part: str) -> list[dict[str, object]]`
- `get_master_placeholders(master_part: str) -> list[dict[str, object]]`

## Placeholders

### `list_placeholders(slide_index: int) -> list[Placeholder]`

Return placeholder objects on a slide.

### `get_slide_layout_ref(slide_index: int)`

Return the layout reference for a slide.

### `set_placeholder_content(slide_index: int, ph_index: int, ph_type: str = "", **kwargs) -> None`

Update inherited placeholder content.

This method supports text, image, table, and chart payloads.

## Shapes and text

### `list_shapes(slide_index: int) -> list[Shape]`

List visible shapes on a slide.

### `search_shapes(query: ShapeSearchQuery | str) -> list[ShapeSearchResult]`

Search shapes by text or query object.

### `add_shape(slide_index: int, shape_type: ShapeType, bounds: tuple[float, float, float, float], **kwargs: str | ShapeProps) -> int`

Add a shape using a typed shape enum and explicit bounds.

### `add_textbox(slide_index: int, left: float, top: float, width: float, height: float, *, text: str = "", **kwargs: str | ShapeProps) -> int`

Add a text box to a slide.

### `add_textboxes(slide_index: int, textboxes: list[dict[str, float | str]]) -> list[int]`

Add multiple text boxes in one call.

### `add_connector(slide_index: int, connector_type: ConnectorType, begin_x: float, begin_y: float, end_x: float, end_y: float, **kwargs: str | ShapeProps) -> int`

Add a connector between two points.

### `add_connectors(slide_index: int, connectors)`

Add multiple connectors.

### Grouping

- `add_group_shape(slide_index: int, shapes: list[int] | None = None) -> int`
- `group_shapes(slide_index: int, shape_ids)`
- `ungroup_shapes(slide_index: int, shape_id)`

### Freeform

- `build_freeform(...)`
- `commit_freeform(...)`

### Shape mutation

- `remove_shape(slide_index: int, shape_id: int) -> None`
- `update_shape(slide_index: int, shape_id: int, updates: ShapeUpdate) -> None`
- `move_shape_to_front(slide_index: int, shape_id: int) -> None`
- `move_shape_to_back(slide_index: int, shape_id: int) -> None`
- `move_shape_to_index(slide_index: int, shape_id: int, target_index: int) -> None`

### Text runs

- `find_and_replace(find_text: str, replace_text: str) -> int`
- `get_slide_text_states(slide_index: int)`
- `get_shape_text_state(slide_index: int, shape_id: int)`
- `get_shape_runs(slide_index: int, shape_id: int)`
- `set_shape_runs(slide_index: int, shape_id: int, runs)`
- `update_shape_run_text(slide_index: int, shape_id: int, run_index: int, text: str) -> None`
- `append_shape_run(slide_index: int, shape_id: int, text: str, **style) -> None`
- `update_slide_run_texts(slide_index: int, updates) -> None`
- `update_deck_run_texts(updates) -> None`

## Tables

### `add_table(slide_index: int, rows: int, cols: int, bounds: tuple[int, int, int, int]) -> int`

Create a table shape.

### Table inspection and styling

- `get_table(slide_index: int, shape_id: int) -> TableInfo`
- `set_table_style(slide_index: int, shape_id: int, style_guid: str) -> None`
- `define_table_style(name: str, style_id: str | None = None) -> str`
- `list_table_styles() -> list[dict[str, str]]`
- `set_table_flags(slide_index: int, shape_id: int, flags: dict[str, bool]) -> None`

### Table content

- `set_table_cell_text(slide_index: int, shape_id: int, row: int, col: int, text: str) -> None`
- `get_table_cell(slide_index: int, shape_id: int, row: int, col: int) -> TableCellInfo`
- `merge_table_cells(slide_index: int, shape_id: int, cell_range: tuple[int, int, int, int]) -> None`
- `split_table_cell(slide_index: int, shape_id: int, row: int, col: int) -> None`
- `set_table_row_height(slide_index: int, shape_id: int, row: int, height) -> None`
- `set_table_column_width(slide_index: int, shape_id: int, col: int, width) -> None`

## Charts

### `add_chart(slide_index: int, chart_type: str, categories: list[str], values_or_series: list[float] | list[dict[str, str | list[float]]], **kwargs) -> int`

Create a chart shape.

### Chart state

- `list_slide_charts(slide_index: int) -> list[SlideChartRef]`
- `get_chart_state(slide_index: int, chart_selector) -> ChartState`
- `get_chart_state_by_index(slide_index: int, chart_index: int) -> ChartState`
- `get_chart_state_by_rel_id(slide_index: int, rel_id: str) -> ChartState`

### Data updates

- `update_chart_data(slide_index: int, chart_selector, data) -> None`
- `update_chart_data_batch(slide_index: int, updates) -> None`
- `update_chart_data_by_index(slide_index: int, chart_index: int, data) -> None`
- `update_chart_data_by_rel_id(slide_index: int, rel_id: str, data) -> None`
- `replace_chart_data_by_index(slide_index: int, chart_index: int, categories: list[str], values: list[float]) -> None`
- `replace_chart_data_by_rel_id(slide_index: int, rel_id: str, categories: list[str], values: list[float]) -> None`

### Formatting

- `update_chart_formatting(slide_index: int, chart_selector, fmt) -> None`
- `update_chart_formatting_by_index(...)`
- `update_chart_formatting_by_rel_id(...)`

## Notes, comments, sections

### Notes

- `get_notes(slide_index: int) -> str`
- `set_notes(slide_index: int, text: str) -> None`
- `notes_slide_exists(slide_index: int) -> bool`
- `get_notes_payload(slide_index: int)`
- `set_notes_shape_text(slide_index: int, shape_id: int, text: str) -> None`
- `set_notes_shape_props(slide_index: int, shape_id: int, **props) -> None`
- `list_notes_shapes(slide_index: int) -> list[dict[str, object]]`
- `list_notes_placeholders(slide_index: int) -> list[dict[str, object]]`
- `update_notes_master(**kwargs) -> None`

### Comments

- `get_authors() -> list[Author]`
- `add_author(name: str, initials: str) -> int`
- `get_comments(slide_index: int) -> list[Comment]`
- `add_comment(slide_index: int, author_id: int, text: str, x: int = 0, y: int = 0) -> int`
- `remove_comment(slide_index_or_index: int, author_id: int | None = None, author_index: int | None = None) -> None`

### Sections

- `get_sections() -> list[Section]`
- `add_section(name: str, slide_indices: list[int]) -> None`
- `remove_section(name: str) -> None`
- `rename_section(old_name: str, new_name: str) -> None`

## Media and export

### Media

- `add_image(slide_index: int, source: str | bytes, bounds: tuple[float, float, float, float], *, name: str = "", crop: ImageCrop | None = None, rotation: float | None = None, flip_h: bool | None = None, flip_v: bool | None = None, **kwargs) -> int`
- `get_image_metadata(slide_index: int, shape_id: int) -> ImageMetadata`
- `list_slide_images(slide_index: int) -> list[SlideImageRef]`
- `swap_image_by_index(slide_index: int, image_index: int, data: bytes, img_format: str) -> None`
- `swap_image_by_rel_id(slide_index: int, rel_id: str, data: bytes, img_format: str) -> None`
- `add_video(...)`
- `add_audio(...)`
- `add_ole_object(...)`

### Export

- `save_as_pdf(path: str, **kwargs) -> None`
- `export_html(path: str, **kwargs) -> None`

### Packaging

- `add_custom_xml(...)`
- `list_custom_xml()`
- `remove_custom_xml(index)`
- `add_vba_project(data)`

## Validation and protection

- `validate() -> list[dict[str, object]]`
- `repair() -> dict[str, object]`
- `set_modify_password(password: str) -> None`
- `set_mark_as_final(final: bool = True) -> None`
- `has_digital_signature() -> bool`

## Advanced features

- `add_smartart(layout: str, items: list[str], bounds: tuple[float, float, float, float] | None = None) -> int`
- `add_animation(shape_id: int, effect: str, *, trigger: str = ..., duration_ms: int = ..., delay_ms: int = ...) -> None`
- `set_transition(transition_type: str, *, duration_ms: int = ..., advance_ms: int = ...) -> None`

## Builder helpers

### `PresentationBuilder`

The builder is useful when you want a fluent construction API.

```python
from gopptx import PresentationBuilder

deck = (
    PresentationBuilder("Deck Title")
    .with_author("Alice")
    .with_theme("aurora")
    .add_title_slide("Intro")
    .build()
)
```

Methods:

- `__init__(title: str) -> None`
- `with_author(author: str) -> PresentationBuilder`
- `with_subject(subject: str) -> PresentationBuilder`
- `with_keywords(keywords: str) -> PresentationBuilder`
- `with_description(description: str) -> PresentationBuilder`
- `with_theme(theme: str) -> PresentationBuilder`
- `with_slide_size(width_inches: float, height_inches: float) -> PresentationBuilder`
- `with_modify_password(password: str) -> PresentationBuilder`
- `with_mark_as_final(final: bool = ...) -> PresentationBuilder`
- `add_title_slide(title: str, layout: str = ...) -> PresentationBuilder`
- `add_bullet_slide(title: str, bullets: list[str], layout: str = ...) -> PresentationBuilder`
- `build() -> Presentation`
- `save(path: str) -> None`

### `ShapeBuilder`

Use this helper to define a shape before adding it to a slide.

Methods:

- `of(shape_type: ShapeType, x: float, y: float, w: float, h: float) -> ShapeBuilder`
- `rectangle(x: float, y: float, w: float, h: float) -> ShapeBuilder`
- `rounded_rectangle(x: float, y: float, w: float, h: float) -> ShapeBuilder`
- `ellipse(x: float, y: float, w: float, h: float) -> ShapeBuilder`
- `triangle(x: float, y: float, w: float, h: float) -> ShapeBuilder`
- `right_triangle(x: float, y: float, w: float, h: float) -> ShapeBuilder`
- `diamond(x: float, y: float, w: float, h: float) -> ShapeBuilder`
- `pentagon(x: float, y: float, w: float, h: float) -> ShapeBuilder`
- `hexagon(x: float, y: float, w: float, h: float) -> ShapeBuilder`
- `parallelogram(x: float, y: float, w: float, h: float) -> ShapeBuilder`
- `cloud(x: float, y: float, w: float, h: float) -> ShapeBuilder`
- `heart(x: float, y: float, w: float, h: float) -> ShapeBuilder`
- `star5(x: float, y: float, w: float, h: float) -> ShapeBuilder`
- `star6(x: float, y: float, w: float, h: float) -> ShapeBuilder`
- `with_text(text: str) -> ShapeBuilder`
- `with_fill(color: str) -> ShapeBuilder`
- `with_no_fill() -> ShapeBuilder`
- `with_line(color: str, *, width_emu: int | None = ..., dash_style: str | None = ...) -> ShapeBuilder`
- `with_no_line() -> ShapeBuilder`
- `with_shadow(color: str = ..., *, blur_emu: int | None = ..., distance_emu: int | None = ..., angle_deg: float = ...) -> ShapeBuilder`
- `with_rotation(degrees: float) -> ShapeBuilder`
- `flip_horizontal() -> ShapeBuilder`
- `flip_vertical() -> ShapeBuilder`
- `shape_type`
- `bounds`
- `to_kwargs() -> dict[str, object]`

### `RunBuilder`

Use this helper to build a formatted text run.

Methods:

- `text(value: str) -> RunBuilder`
- `bold(value: bool = True) -> RunBuilder`
- `italic(value: bool = True) -> RunBuilder`
- `underline(style: str = "sng") -> RunBuilder`
- `strikethrough(style: str = "sng") -> RunBuilder`
- `subscript(value: bool = True) -> RunBuilder`
- `superscript(value: bool = True) -> RunBuilder`
- `color(hex_color: str) -> RunBuilder`
- `highlight(hex_color: str) -> RunBuilder`
- `font(name: str) -> RunBuilder`
- `size_pt(points: int) -> RunBuilder`
- `code(value: bool = True) -> RunBuilder`
- `all_caps(value: bool = True) -> RunBuilder`
- `small_caps(value: bool = True) -> RunBuilder`
- `hyperlink(address: str, *, tooltip: str = "") -> RunBuilder`
- `hover_action(address: str) -> RunBuilder`
- `build() -> TextRun`
- `to_payload() -> dict[str, object]`

## See also

- [API Reference](../api-reference.md)
- [Bridge Operations](bridge-operations.md)
