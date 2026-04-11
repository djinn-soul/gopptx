# Python Presentation API

This reference is auto-generated from the Python source and `.pyi` stubs in `python/gopptx/`.
It updates automatically whenever `python/**` changes are pushed.

Canonical signature source: `python/gopptx/api.pyi`

## Task-Oriented Method Map

Use this map for quick navigation before the full class docs below.

### Session lifecycle

- `Presentation.new(title)`
- `Presentation(path)`
- `open(path)`, `open_bytes(data)`, `save(path)`, `to_bytes()`, `close()`

### Slide management

- `add_slide(...)`, `remove_slide(index)`, `move_slide(from_index, to_index)`
- `duplicate_slide(index, insert_at=None)`, `duplicate_slide_after(index)`
- `set_slide_title(index, title)`, `set_slide_hidden(index, hidden)`
- `list_slide_layouts()`, `rebind_slide_layout(slide_index, layout_part)`

### Masters and layouts

- `slide_masters()`
- `add_slide_master()`, `remove_slide_master(master_part)`
- `add_slide_layout(master_part, layout_name="Custom Layout")`, `remove_slide_layout(layout_part)`
- `clone_layout_master_family(source_slide_index, target_slide_index)`

### Shapes and text

- `add_shape(...)`, `add_textbox(...)`, `add_textboxes(...)`, `add_connector(...)`
- `add_group_shape(...)`, `build_freeform(...)`
- `update_shape(...)`, `remove_shape(slide_index, shape_id)`, `clear_shapes(slide_index)`
- `update_slide_run_texts(...)`, `update_deck_run_texts(...)`, `find_and_replace(...)`

### Tables and charts

- Tables: `add_table(...)`, `get_table(...)`, `set_table_cell_text(...)`, `merge_table_cells(...)`, `split_table_cell(...)`
- Charts: `add_chart(...)`, `get_chart_state(...)`, `update_chart_data(...)`, `replace_chart_data_by_index(...)`

### SmartArt, transitions, and animation

- `add_smartart(...)`, `update_smartart(shape_id, items)`
- `delete_smartart(shape_id)`, `change_smartart_layout(shape_id, layout)`, `set_smartart_style(...)`, `set_smartart_nodes(...)`
- `add_animation(...)`, `set_transition(...)`

### Export and validation

- `validate()`, `repair()`, `convert_to_grayscale(...)`

### Notes, comments, sections, and metadata

- Notes: `get_notes(slide_index)`, `set_notes(slide_index, text)`, `list_notes_shapes(slide_index)`
- Comments: `get_comments(slide_index)`, `add_comment(...)`, `remove_comment(...)`
- Sections: `sections()`, `add_section(name, slide_indices)`, `rename_section(...)`, `remove_section(name)`
- Core properties: `core_properties`, `title`, `author`, `comments`, `identifier`, `language`, `last_printed`, `version`

---

## Presentation

Main entry point. Open an existing deck or create a new one.

```python
from gopptx import Presentation

with Presentation.new("Deck Title") as pres:
    pres.add_title_slide("Hello")
    pres.save("out.pptx")
```

::: gopptx.presentation.presentation.Presentation
    options:
      inherited_members: true
      show_signature_annotations: true
      show_source: false

---

## PresentationBuilder

Fluent construction API for building decks step by step.

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

::: gopptx.builder.PresentationBuilder
    options:
      show_source: false

---

## Slide

Represents a single slide within a presentation.

::: gopptx.slide.slide.Slide
    options:
      inherited_members: true
      show_source: false

---

## Table

Fluent table interface returned by `add_table()`.

::: gopptx.api.Table
    options:
      show_source: false

---

## Cell / CellRange

::: gopptx.api.Cell
    options:
      show_source: false

::: gopptx.api.CellRange
    options:
      show_source: false

---

## Freeform

::: gopptx.api.FreeformBuilder
    options:
      show_source: false

---

## Chart data helpers

::: gopptx.api.CategoryChartData
    options:
      show_source: false

::: gopptx.api.CategorySeries
    options:
      show_source: false

::: gopptx.api.XySeries
    options:
      show_source: false

---

## Errors

::: gopptx.api_errors.GopptxError
    options:
      show_source: false

---

## See also

- [API Reference](../api-reference.md)
- [Bridge Operations](bridge-operations.md)
