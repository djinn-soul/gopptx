# gopptx

High-performance PowerPoint (PPTX) engine powered by Go.

## Quick Start

### Install (Go)

```bash
go get github.com/djinn-soul/gopptx
```

### Basic usage (Python)

```python
from gopptx import Presentation

with Presentation("input.pptx") as pres:
    pres.add_slide("Hello from Python")
    pres.save("output.pptx")
```

## Core Capabilities

- Fast PPTX generation and manipulation.
- Concurrency-oriented processing for large workloads.
- Lazy-loading behavior for lower memory pressure.
- Mermaid rendering (12+ diagram types) with theme support.
- Native pie-slice rendering for Mermaid pie charts.
- C-compatible bridge for Python and other languages.
- Stable, typed JSON command API.

## Python Library

`gopptx` can be used as a high-performance Python library with an explicit command-first API over the Go engine.

### Build and install

1. Build the Go shared library and bundle it into the Python package.
   ```powershell
   .\scripts\build_python.ps1
   ```
2. Install the package.
   ```bash
   pip install -e .
   ```

More examples: [`python/README.md`](python/README.md)

### Throughput tip: use batching

Use batch execution for write-heavy loops to reduce Python -> C -> Go boundary crossings.

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

You can also use the fluent batch context manager:

```python
with Presentation.new("Batch Context") as pres:
    with pres.batch(stop_on_error=False) as batch:
        batch.add_slide("A")
        batch.add_slide("B")
        batch.set_slide_title(0, "Updated")
```

Optional: install `orjson` to speed up Python-side JSON encode/decode.

Batch contract details: [`docs/architecture/batch_execute_envelope.md`](docs/architecture/batch_execute_envelope.md)

## JSON Command Bridge

All bridge operations use a JSON envelope.

### Request

```json
{
  "api_version": 1,
  "request_id": "uuid",
  "op": "add_slide",
  "payload": { "title": "New Slide" }
}
```

### Response

```json
{
  "ok": true,
  "result": { "index": 1 },
  "request_id": "uuid"
}
```

### Supported operations

41 operations are available across slides, metadata, shapes, images/charts, tables, sections, comments, notes, layouts/masters, placeholders, and utility commands.

Operation reference: [`docs/architecture/bridge-phase1-ops.md`](docs/architecture/bridge-phase1-ops.md)  
C API reference: [`bindings/c/README.md`](bindings/c/README.md)

## Benchmarks

- Go bridge microbench:
  `go test ./pkg/pptx/editor -run ^$ -bench "BenchmarkBridge(Execute|JSON)" -benchmem -count=3`
- Python benchmark script:
  `uv run python scripts/smoke/python_batch_latency_benchmark.py`
- JSON profile and transport decision record:
  [`docs/benchmarks/json_bridge_profile_2026-02-21.md`](docs/benchmarks/json_bridge_profile_2026-02-21.md)

## Contributing

Contribution guide: [`CONTRIBUTING.md`](CONTRIBUTING.md)

## License

This project is licensed under Apache License 2.0. See [`LICENSE`](LICENSE).
