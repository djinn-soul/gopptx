# Python library guide

This is the working guide for `gopptx.Presentation` — the surface most Python code touches.
For an exhaustive list of all 208 methods see the
[Python `Presentation` API reference](../reference/python-presentation-api.md).

!!! warning "Coordinates are EMU"

    `from gopptx import Emu, Inches, Point`. A bare `600` is 0.00066 inch. See
    [Units](../concepts.md#units-and-geometry).

## Two ways to write the same thing

gopptx gives you a **flat API** on `Presentation` and a **navigable object layer** through
`pres.slides`. They call the same engine; pick per situation.

```python
# Flat — you hold the slide index
pres.add_textbox(1, Inches(1), Inches(1), Inches(3), Inches(0.5), text="Hello")

# Object layer — you hold the slide
slide = pres.slides[1]
slide.add_textbox(Inches(1), Inches(1), Inches(3), Inches(0.5), text="Hello")
```

The flat API is better for generated code and batching. The object layer is better for reading
and editing an existing deck, where you are iterating over what is already there.

## Opening and saving

```python
from gopptx import Presentation

Presentation.new("Title")              # new deck — already has a title slide at index 0
Presentation("input.pptx")             # open a file
Presentation.open_deck("input.pptx")   # same, named constructor
Presentation.from_template("brand.pptx", {"client": "Acme"})
```

```python
with Presentation.new("Deck") as pres:
    ...
    pres.save("out.pptx")     # write to disk
    data = pres.to_bytes()    # or serialise in memory
```

Always use the context manager — the handle is Go-owned memory, and leaking it keeps that memory
alive for the life of the process.

## Slides

```python
slide = pres.add_slide("Title")                 # append; returns a Slide
pres.insert_slide(2, "Inserted")
pres.duplicate_slide(0, insert_at=3)
pres.move_slide(3, 1)
pres.remove_slide(4)
pres.set_slide_title(0, "New Title")
pres.set_slide_hidden(2, hidden=True)

pres.slide_count                                # property, not a call
for slide in pres.slides:
    print(slide.index, slide.title)
```

Convenience constructors that build a whole slide:

```python
pres.add_title_slide("FY26 Q3")
pres.add_bullet_slide("Highlights", ["Revenue +12%", "Retention +4%"])
pres.add_paragraph_slide("Summary", "We beat plan in every region.")
pres.add_slide_from_markdown("# Agenda\n- One\n- Two")
pres.add_slide_from_url("https://example.com/post")
```

Pulling slides in from another deck:

```python
pres.copy_slides_from("other.pptx", [0, 2, 5])   # selected slides
pres.copy_slides_from("other.pptx")              # all of them
pres.merge_from_file("appendix.pptx")            # append the whole deck
```

## Text

```python
shape_id = pres.add_textbox(
    0, Inches(0.5), Inches(1.0), Inches(4.0), Inches(0.8),
    text="Revenue by quarter",
)
```

Rich text is expressed as **runs** — a run is a span with uniform formatting:

```python
pres.set_shape_runs(0, shape_id, [
    {"text": "Revenue ", "bold": True, "size": 24, "color": "1F4E78"},
    {"text": "+12%", "bold": True, "size": 24, "color": "2E7D32"},
])

runs = pres.get_shape_runs(0, shape_id)
pres.append_shape_run(0, shape_id, {"text": " YoY", "italic": True})
```

Deck-wide text edits:

```python
count = pres.find_and_replace("Draft", "Final")           # returns replacements made

# Rewrite specific runs across several slides in one request
pres.update_deck_run_texts([
    {"slide_index": 0, "updates": [{"shape_id": shape_id, "run_index": 0, "text": "Acme"}]},
])
```

## Shapes

```python
from gopptx import ShapeType

shape_id = pres.add_shape(
    0,
    ShapeType.ROUNDED_RECTANGLE,
    (Inches(0.5), Inches(1.2), Inches(2.0), Inches(0.8)),
    text="On plan",
    properties={
        "fill": {"solid": "DCE6F2"},
        "line": {"color": "1F4E78", "width_emu": 12700},
    },
)
```

`ShapeType` has 202 members. Beyond placement:

```python
pres.list_shapes(0)
pres.update_shape(0, shape_id, {"rotation": 15})
pres.move_shape_to_front(0, shape_id)
pres.group_shapes(0, [id_a, id_b])
pres.remove_shape(0, shape_id)
pres.search_shapes({"text_contains": "plan"})
```

Connectors and freeforms:

```python
from gopptx import ConnectorType

pres.add_connector(0, ConnectorType.ELBOW,
                   Inches(2.6), Inches(2.7), Inches(4.6), Inches(2.7))

pres.commit_freeform(
    0,
    [(Inches(1), Inches(1)), (Inches(2), Inches(2)), (Inches(1), Inches(2))],
    close=True,
)
```

`build_freeform(slide_index, start_x, start_y, scale)` returns a `FreeformBuilder` if you would
rather accumulate points fluently; `commit_freeform` is the one-shot form and needs `close` as a
keyword.

## Tables

```python
table_id = pres.add_table_from_rows(
    0,
    [["Region", "Revenue"], ["EMEA", "4.1M"], ["APAC", "2.8M"]],
    bounds=(Inches(0.5), Inches(1.5), Inches(4.0), Inches(2.0)),
    first_row=True,
    band_row=True,
)

pres.add_table_from_dicts(0, [{"Region": "EMEA", "Revenue": "4.1M"}], bounds=...)
```

`add_table_from_rows` infers the row and column counts from the data; `add_table()` itself needs
explicit `rows` and `cols`.

Then edit it:

```python
pres.set_table_cell_text(0, table_id, 1, 1, "4.4M")
pres.merge_table_cells(0, table_id, (0, 0, 0, 1))     # row1, col1, row2, col2
pres.split_table_cell(0, table_id, 1, 0)
pres.set_table_column_width(0, table_id, 0, Inches(1.5))
pres.set_table_row_height(0, table_id, 0, Inches(0.4))
pres.set_table_style(0, table_id, "MediumStyle2Accent1")
pres.list_table_styles()
```

## Charts

```python
from gopptx import ChartType

chart_id = pres.add_chart(
    0,
    ChartType.COLUMN,
    ["Q1", "Q2", "Q3", "Q4"],
    [12.0, 15.5, 18.0, 21.0],
    bounds=(Inches(5.0), Inches(1.2), Inches(4.0), Inches(3.0)),
    title="Quarterly trend",
)
```

Pass a `ChartType` member, not a string — bare strings are deprecated and will be rejected in a
future release. There are 24 types, from `COLUMN` and `LINE` to `RADAR_FILLED`, `STOCK_OHLC` and
`BUBBLE`, plus `add_combo_chart` for mixed bar/line.

Generated charts embed an Excel workbook, so the data stays editable inside PowerPoint.

Updating a chart in an existing deck:

```python
pres.list_slide_charts(2)
pres.update_chart_data_by_index(2, 0, {
    "categories": ["Q1", "Q2"],
    "series": [{"name": "Revenue", "values": [10.0, 14.0]}],
})
pres.update_chart_formatting_by_index(2, 0, {"title": "Restated"})
state = pres.get_chart_state_by_index(2, 0)
```

Multi-series data is a list of series dicts, each with `name` and `values`; scatter and bubble
series use `x_values` / `y_values` / `sizes`.

## Images and media

```python
pres.add_image(0, "logo.png", bounds=(Inches(8), Inches(0.3), Inches(1.5), Inches(0.6)))
pres.add_image(0, image_bytes, bounds=...)          # bytes work too
pres.add_picture(0, "photo.jpg", Inches(1), Inches(1), Inches(4), Inches(3))

pres.list_slide_images(0)
pres.swap_image_by_index(0, 0, new_bytes, "png")    # reskin without re-layout
data = pres.extract_media("ppt/media/image1.png")
pres.save_media("ppt/media/image1.png", "out.png")

pres.add_video(0, "clip.mp4", bounds=...)
pres.add_audio(0, "voice.m4a", bounds=...)
pres.add_online_video(0, "https://youtu.be/…", bounds=...)
pres.add_ole_object(0, "model.xlsx", bounds=...)
```

## SmartArt and Mermaid

Both live on the slide object:

```python
from gopptx.smartart import SMARTART_BASIC_PROCESS

slide = pres.slides[1]

slide.add_smartart(
    SMARTART_BASIC_PROCESS,
    ["Collect", "Model", "Ship"],
    bounds=(Inches(1), Inches(1.5), Inches(8), Inches(3)),
)

slide.add_mermaid("graph LR; A[Ingest] --> B[Transform] --> C[Publish]")
```

SmartArt is generated from the data model, the way PowerPoint does it — so PowerPoint re-lays it
out correctly rather than showing empty `[Text]` placeholders. Nodes, quick styles and colour
styles can all be edited afterwards (`set_smartart_nodes`, `set_smartart_style`,
`change_smartart_layout`).

## Notes, comments and sections

```python
pres.set_notes(0, "Open with the revenue number.")
pres.get_notes(0)
pres.notes_slide_exists(0)

author_id = pres.add_author("Ada Lovelace", "AL")
pres.add_comment(0, author_id, "Check this figure before Friday.")
pres.get_comments(0)

pres.add_section("Results", [1, 2, 3])      # a section names a set of slides
pres.get_sections()
pres.rename_section("Results", "Q3 Results")
```

## Themes, layouts and masters

`apply_theme` accepts either a `Theme` object or a theme name:

```python
from gopptx.presentation.theme import get_theme, list_themes

list_themes()                       # ['aurora', 'ocean', 'sunset', 'forest']
pres.apply_theme(get_theme("ocean"))
```

```python
from gopptx import THEME_CORPORATE

pres.apply_theme(THEME_CORPORATE)   # gopptx theme names
pres.apply_theme("integral")        # Office preset names
```

Accepted names are the gopptx themes — `Corporate`, `Modern`, `Vibrant`, `Dark`, `Nature`,
`Tech`, `Carbon`, `Office`, matching the exported `THEME_*` constants — and the Office presets
`office`, `office2013`, `facet`, `integral`, `ion`, `retrospect`, `slice`, `wisp`. Matching
ignores case and separators. `pres.set_global_theme_preset(name)` takes the same set.

Custom schemes and the rest of the layout surface:

```python
pres.set_theme_color_scheme(accent1="1F4E78", accent2="2E7D32")   # keyword arguments
pres.set_theme_font_scheme("Segoe UI", "Segoe UI")                # major, minor — positional
pres.get_theme_inventory()

pres.set_slide_size(12192000, 6858000)     # 16:9, in EMU
layouts = pres.list_slide_layouts()
pres.rebind_slide_layout(2, layouts[5]["part"])   # takes a layout part path, not an index
pres.list_placeholders(2)
pres.set_placeholder_content(2, 1, text="Body text")
```

`set_placeholder_content` requires exactly one of `text=`, `image_path=`, `table=` or `chart=`.

## Document properties and protection

```python
pres.set_core_properties({"title": "FY26 Q3", "creator": "Reporting Bot",
                          "subject": "Quarterly review", "keywords": "qbr,finance"})
pres.core_properties
pres.author, pres.title, pres.language, pres.version

pres.set_mark_as_final(final=True)     # `final` is keyword-only
pres.set_modify_password("s3cret")
pres.is_digitally_signed()
```

## Motion

```python
slide = pres.slides[1]
slide.set_transition("fade", duration_ms=400)
slide.add_animation(shape_id, "fadeIn", trigger="onClick", duration_ms=500)
```

## Export

```python
from gopptx import HTMLOptions, PDFOptions

pres.export_pdf("deck.pdf", PDFOptions(driver="auto"))
pres.export_html("deck.html", HTMLOptions(embed_images=True, include_navigation=True))
pres.save_flat_xml("deck.xml")
pres.convert_to_grayscale()
```

`export_pdf` pairs with `export_html`; `save_as_pdf` is a deprecated alias for it. See the
[export guide](export.md).

## Throughput

```python
with pres.batch(stop_on_error=True) as batch:
    for i in range(500):
        batch.add_slide(f"Slide {i}")
```

Reads are rejected inside a `batch()` block. See [Batch execution](batch-execution.md).

## Escape hatch

Anything the typed API does not wrap is one call away:

```python
result = pres.execute("set_slide_hidden", {"slide_index": 2, "hidden": True})
```

All 179 operations are listed in [Bridge operations](../reference/bridge-operations.md).

## Reading results

Bridge results are JSON from Go structs, so some carry Go field names. The Python layer adds
snake_case aliases to every response while keeping the originals, so either spelling works:

```python
shape = pres.list_shapes(0)[0]
shape["id"], shape["ID"]              # both present
shape["alt_text"], shape["AltText"]   # both present
```

## Practices that pay off

- Use the context manager, always.
- Wrap every coordinate in `Inches()` / `Point()` / `Emu()`.
- Batch write-heavy loops; keep reads outside the batch.
- Open the result in real PowerPoint before you trust a layout. Tests pass on files that
  PowerPoint refuses.
