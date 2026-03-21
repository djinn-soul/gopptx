# Core Concepts

## Architecture Layers

1. Go engine (`pkg/pptx`): PPTX editing and serialization logic
2. C bridge (`bindings/c`): handle-based ABI layer
3. Python runtime (`python/gopptx`): friendly API over bridge commands

## Session Lifecycle

- `Presentation.new(title)` for new decks
- `Presentation(path)` or `.open(path)` for existing files
- Execute many operations in memory
- `save(path)` once at logical checkpoints
- `close()` or context manager exit to release handle

## Command Envelope Contract

Each operation executes through a stable envelope:

Request:

```json
{
  "api_version": 1,
  "request_id": "uuid",
  "op": "add_slide",
  "payload": {"title": "Agenda"}
}
```

Response:

```json
{
  "ok": true,
  "request_id": "uuid",
  "result": {"index": 2}
}
```

If `ok` is `false`, error details include message and error code.

## Batch Model

- `execute_batch(commands)` sends many operations in one crossing
- `with pres.batch(...):` buffers mutating calls and flushes once
- Read operations are blocked inside `batch()` by design

Each Python -> C -> Go boundary crossing has overhead. Batching reduces crossings and improves throughput for write-heavy workloads.

See [Batch Execution](guides/batch-execution.md) for patterns.
