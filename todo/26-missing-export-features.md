# 26 - Missing Export Features

Scope: track missing export/output capabilities in `gopptx` from the provided feature audit screenshot.

## Missing in gopptx

| Feature | Status | Tracking |
| --- | --- | --- |
| PNG per slide | Missing | Tracked in [10-conversion-import-export.md](10-conversion-import-export.md). |
| Browser download (Blob/base64) | Missing | [ ] Add browser-targeted download helpers (blob/base64 packaging) for web integrations. |
| ArrayBuffer / Uint8Array output | Missing | [ ] Add raw binary output helpers for JS runtimes (ArrayBuffer/Uint8Array equivalents). |
| HTTP stream output | Missing | [ ] Add streaming output API for HTTP response pipelines. |

## Verification Follow-up

- [ ] Verify listed gaps against current runtime behavior and update this file if needed.
- [ ] Add regression tests for browser/binary/stream output flows once implemented.
