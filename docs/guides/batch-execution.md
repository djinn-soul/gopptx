# Batch Execution

Use batching when issuing many operations from Python.

## Direct Command List

```python
from gopptx import Presentation, ops

with Presentation.new("Batch Demo") as pres:
    commands = [
        {"op": ops.OP_ADD_SLIDE, "payload": {"title": f"Slide {i}"}}
        for i in range(200)
    ]
    results = pres.execute_batch(commands)
    assert all(item.get("ok", False) for item in results)
```

## Fluent Batch Context

```python
from gopptx import Presentation

with Presentation.new("Batch Context") as pres:
    with pres.batch(stop_on_error=False) as batch:
        batch.add_slide("A")
        batch.add_slide("B")
        batch.set_slide_title(0, "Updated")
```

## Limitations

- Read operations (e.g., `get_slide_title`) are blocked inside `batch()` contexts by design.
- Move read operations outside the batch block or use direct `execute_batch` with mixed commands.

## Rules of Thumb

- Batch write-heavy loops.
- Keep requests explicit and typed.
- Validate batch results per item when partial failures are allowed.
