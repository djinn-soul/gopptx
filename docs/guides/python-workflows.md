# Python Workflows

Use this page for practical Python patterns with `gopptx.Presentation`.

## 1) Create New Decks

```python
from gopptx import Presentation

with Presentation.new("Weekly Report") as pres:
    pres.add_slide("Overview")
    pres.add_bullet_slide("Highlights", ["Revenue +8%", "NPS +3"])
    pres.save("weekly_report.pptx")
```

## 2) Edit Existing Decks

```python
from gopptx import Presentation

with Presentation("input.pptx") as pres:
    pres.set_slide_title(0, "Updated Title")
    pres.find_and_replace("Draft", "Final")
    pres.save("output_edited.pptx")
```

## 3) High-Throughput Writes

```python
from gopptx import Presentation, ops

with Presentation.new("Batch") as pres:
    commands = [
        {"op": ops.OP_ADD_SLIDE, "payload": {"title": f"Slide {i}"}}
        for i in range(50)
    ]
    pres.execute_batch(commands, stop_on_error=True)
    pres.save("batch_output.pptx")
```

## Related

- [Python Library](python-library.md)
- [JSON Bridge Operations](../reference/bridge-operations.md)
