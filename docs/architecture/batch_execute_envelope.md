# Batch execute envelope

The wire format of the `batch_execute` operation — the one command that carries other commands.

**Source of truth:** `pkg/pptx/editor/modules/command/batch.go`.

## Request

`batch_execute` is an ordinary operation, so it travels inside the normal envelope. Its payload
is a list of commands.

```json
{
  "api_version": 1,
  "request_id": "outer-1",
  "op": "batch_execute",
  "payload": {
    "commands": [
      {"op": "add_slide", "payload": {"title": "One"}, "request_id": "c1"},
      {"op": "add_slide", "payload": {"title": "Two"}, "request_id": "c2"},
      {"op": "set_slide_title", "payload": {"slide_index": 0, "title": "Updated"}}
    ],
    "stop_on_error": true
  }
}
```

| Field | Type | Meaning |
|---|---|---|
| `commands` | array | Executed in order against the same open presentation |
| `commands[].op` | string | Any of the 179 operations except `batch_execute` itself |
| `commands[].payload` | object | The operation's own payload, unchanged |
| `commands[].request_id` | string, optional | Echoed back on the matching result |
| `stop_on_error` | bool, optional | Default `false` |

## Response

The outer envelope succeeds as long as the batch itself was well-formed. Per-command outcomes
live in `result.results`, one entry per command, in request order.

```json
{
  "ok": true,
  "request_id": "outer-1",
  "result": {
    "results": [
      {"ok": true, "op": "add_slide", "request_id": "c1", "result": {"index": 1}},
      {"ok": true, "op": "add_slide", "request_id": "c2", "result": {"index": 2}},
      {
        "ok": false,
        "op": "set_slide_title",
        "error": {
          "code": "INVALID_INDEX",
          "message": "payload validation failed",
          "details": {"cause": ["slide_index 99 out of bounds [0, 2)"], "index": 2}
        }
      }
    ]
  }
}
```

| Field | Meaning |
|---|---|
| `results[].ok` | Whether that command succeeded |
| `results[].op` | Echo of the operation name |
| `results[].request_id` | Echo, when the request carried one |
| `results[].result` | The operation's own result, on success |
| `results[].error` | `{code, message, details?}`, on failure. `details.index` is the command's position in the batch, and `details.cause` lists the validation failures. |

## Error semantics

| | `stop_on_error: false` (default) | `stop_on_error: true` |
|---|---|---|
| On the first failure | Recorded, execution continues | Execution stops |
| `results` length | One entry per command | Entries up to and including the failure |
| Outer `ok` | `true` | `true` — the batch ran; check the entries |

**The outer `ok` is not a summary.** A batch in which every command failed still returns
`"ok": true` at the envelope level, because the batch operation itself worked. Always inspect
each entry.

A malformed batch payload — not a failing command, but an unparseable request — fails the outer
envelope with `INVALID_PAYLOAD`.

## Constraints

- `batch_execute` cannot nest. A `batch_execute` inside `commands` fails that entry with
  `INVALID_BATCH_ITEM` — `"nested batch_execute is not supported"` — rather than failing the
  whole request.
- Commands share one presentation and execute sequentially, so later commands see the effects of
  earlier ones. Index-based payloads must account for slides added earlier in the same batch.
- An unknown `op` yields `UNKNOWN_OP` for that entry, not for the whole batch.

## From Python

The typed API builds this envelope for you:

```python
from gopptx import Presentation, ops

with Presentation.new("Deck") as pres:
    results = pres.execute_batch(
        [{"op": ops.OP_ADD_SLIDE, "payload": {"title": f"Slide {i}"}} for i in range(100)],
        stop_on_error=False,
    )
    failures = [r for r in results if not r.get("ok")]
```

`execute_batch` returns the `results` array directly. The `pres.batch()` context manager wraps
the same operation, buffering calls and flushing on exit.

## See also

- [Batch execution guide](../guides/batch-execution.md) — when and how to use it
- [Bridge operations](../reference/bridge-operations.md) — the operations you can batch
- [Core concepts](../concepts.md#the-command-envelope) — the single-command envelope
