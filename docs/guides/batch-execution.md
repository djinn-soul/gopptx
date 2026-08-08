# Batch execution

Batching is a Python concern. In Go you are already inside the engine and there is nothing to
batch.

## Why it matters

Every Python call travels Python → ctypes → C ABI → Go, and the answer travels back. The
document operation itself is cheap; the crossing is not. Adding 500 slides one call at a time
means 500 crossings, each with its own JSON encode, decode and handle lookup.

Batching sends many operations per crossing. The engine executes them in order and returns one
array of results.

## The fluent form

```python
from gopptx import Presentation

with Presentation.new("Batch Demo") as pres:
    with pres.batch(stop_on_error=True) as batch:
        for i in range(500):
            batch.add_slide(f"Slide {i}")
    pres.save("batch.pptx")
```

Calls on `batch` are buffered and flushed when the block exits. The API mirrors `Presentation`,
so the loop body reads the same as the unbatched version.

## The explicit form

```python
from gopptx import Presentation, ops

with Presentation.new("Batch Demo") as pres:
    commands = [
        {"op": ops.OP_ADD_SLIDE, "payload": {"title": f"Slide {i}"}}
        for i in range(500)
    ]
    results = pres.execute_batch(commands, stop_on_error=False)

    for i, item in enumerate(results):
        if not item.get("ok"):
            print(i, item.get("error"))

    pres.save("batch.pptx")
```

Each result is a mapping — `{"ok": True, "op": "add_slide", "result": {"index": 2}}` — or, on
failure, carries an `error` with `code` and `message`.

Use the explicit form when you need to:

- mix reads and writes in a controlled order
- issue an operation the typed API does not wrap
- inspect per-command results rather than relying on an exception

## Reads are rejected inside `batch()`

```python
with pres.batch() as batch:
    batch.add_slide("A")
    count = pres.slide_count
    # GopptxError: read operation 'slide_count' is not allowed inside batch()
```

This is deliberate. The buffered writes have not executed yet, so any read would answer from
state that is about to change — a silently wrong answer is worse than an error. Move the read
outside the block, or use `execute_batch()` where you control the ordering yourself.

## Error handling

| | `stop_on_error=True` | `stop_on_error=False` (default) |
|---|---|---|
| First failure | Aborts the batch | Recorded; the batch continues |
| Later commands | Not attempted | All attempted |
| You should | Wrap in `try/except GopptxError` | Inspect every result's `ok` |

`pres.abort_batch()` discards a batch you have begun manually with `begin_batch()` /
`end_batch()`; the context manager handles this for you on an exception.

## Write buffers

Some high-volume paths have their own buffers, flushed on `save()` or on demand:

```python
pres.queue_textbox_add(slide_index, {"left": ..., "top": ..., "width": ..., "height": ..., "text": ...})
pres.flush_pending_textbox_adds(slide_index)
pres.flush_all_pending_textbox_adds()

pres.queue_shape_run_text_update(...)
pres.flush_all_pending_slide_run_text_updates()
```

If you read shape state between queueing and flushing, flush first — the queue is not visible to
reads.

## What to batch, and what not to

**Batch**

- Loops that add slides, textboxes, shapes or table rows
- Bulk text substitution across a deck
- Anything generated from a data source, where the count scales with the data

**Do not bother**

- A handful of operations — the overhead is not worth the indirection
- Anything where the next call depends on the result of the previous one
- Read-heavy inspection passes

## Practical guidance

- One batch per logical unit of work; you do not need to batch the whole document.
- Prefer `stop_on_error=True` when the batch is a transaction, `False` when it is a best-effort
  bulk import you will reconcile afterwards.
- Install `orjson` — the batch payload is the largest JSON gopptx encodes.
- Save once at the end, not inside the loop. `save()` serialises the whole package each time.

## See also

- [Core concepts — batching](../concepts.md#batching)
- [Bridge operations](../reference/bridge-operations.md) — the operation names for the explicit form
- [Batch execute envelope](../architecture/batch_execute_envelope.md) — the wire format
