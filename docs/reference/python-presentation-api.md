# Python Presentation API

This reference is auto-generated from the Python source and `.pyi` stubs in `python/gopptx/`.
It updates automatically whenever `python/**` changes are pushed.

Canonical signature source: `python/gopptx/api.pyi`

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
