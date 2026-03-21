# Python Library

## Main Entry Points

- `from gopptx import Presentation`
- `Presentation.new(title)` for new files
- `Presentation(path)` for existing files

## Core Session Pattern

```python
from gopptx import Presentation

with Presentation("input.pptx") as pres:
    pres.add_slide("Hello")
    pres.save("output.pptx")
```

## Common Slide Operations

```python
with Presentation.new("Slides") as pres:
    s = pres.add_slide("Intro")
    pres.set_slide_title(0, "Updated Intro")
    pres.duplicate_slide_after(0)
    pres.move_slide(1, 0)
```

## Common Shape/Table/Chart Operations

```python
with Presentation.new("Content") as pres:
    pres.add_slide("Data")
    shape_id = pres.add_textbox(0, 40, 40, 400, 50, text="KPI")
    table_id = pres.add_table(0, 3, 4, (40, 120, 600, 220))
    chart_id = pres.add_chart(
        0,
        "bar",
        ["Q1", "Q2", "Q3", "Q4"],
        [12.0, 15.5, 18.0, 21.0],
        bounds=(40, 360, 500, 240),
        title="Quarterly Trend",
    )
    pres.save("content.pptx")
```

## Edit Existing Presentations

```python
with Presentation("existing.pptx") as pres:
    for slide in pres.slides:
        _ = slide  # navigate slide proxies
    pres.find_and_replace("Draft", "Final")
    pres.save("existing_edited.pptx")
```

## Recommended Practices

- Always use context manager (`with`) for deterministic cleanup.
- Batch bulk writes with `execute_batch()` or `batch()`.
- Validate generated files in critical flows with `pres.validate()`.
- Keep payloads explicit; avoid hidden defaults during development.

## Next References

- [Batch Execution](batch-execution.md)
- [Python Presentation API](../reference/python-presentation-api.md)
- [JSON Bridge Operations](../reference/bridge-operations.md)
