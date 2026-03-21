# Batch Execute Envelope

This page documents the request/response envelope used by batch execution.

## Request Shape

```json
{
  "api_version": 1,
  "request_id": "uuid",
  "op": "add_slide",
  "payload": {"title": "Agenda"}
}
```

## Response Shape

```json
{
  "ok": true,
  "request_id": "uuid",
  "result": {"index": 1}
}
```

## Notes

- `api_version` keeps the envelope extensible.
- `request_id` helps correlate responses with batched requests.
- `op` selects the bridge operation.
- `payload` carries the operation-specific data.
