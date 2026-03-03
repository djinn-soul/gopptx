# gopptx

High-performance PowerPoint (PPTX) engine powered by Go.

## Features

- **Fast**: Blazing fast PPTX generation and manipulation.
- **Concurrent**: Optimized for multi-threaded slide processing.
- **Lazy Loading**: Efficient memory usage for large presentations.
- **12+ Mermaid Diagram Types**: Native shape-based rendering for all common Mermaid types (Flowchart, Sequence, Gantt, etc.).
- **Mermaid Themes**: Support for Mermaid themes (Forest, Dark, Neutral) and custom initialization blocks.
- **True Pie Slices**: Native PowerPoint pie slice rendering for Mermaid pie charts.
- **Cross-Language**: C-compatible bindings for Python and other languages.
- **JSON Command API**: Stable, typed JSON bridge for C/Python clients.

## Installation (Go)

```bash
go get github.com/djinn-soul/gopptx
```

## Python Library

`gopptx` can be used as a high-performance Python library. It provides an explicit command-first API over the Go engine.

### Build and Install

1. Build the Go shared library and bundle it into the Python package:
   ```powershell
   .\scripts\build_python.ps1
   ```
2. Install the package:
   ```bash
   pip install -e .
   ```

### Usage

```python
from gopptx import Presentation

with Presentation("input.pptx") as pres:
    pres.add_slide("Hello from Python")
    pres.save("output.pptx")
```

See also: [`python/README.md`](python/README.md) for batch patterns and table command API examples.

### Bridge Throughput Tips

Use batching for write-heavy loops to reduce Python -> C -> Go boundary crossings:

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

Or use the context manager for fluent code:

```python
with Presentation.new("Batch Context") as pres:
    with pres.batch(stop_on_error=False) as batch:
        batch.add_slide("A")
        batch.add_slide("B")
        batch.set_slide_title(0, "Updated")
```

Optional: install `orjson` to speed up Python-side bridge JSON encode/decode.

Batch request/response contract details: [`docs/architecture/batch_execute_envelope.md`](docs/architecture/batch_execute_envelope.md)

## JSON Command Bridge

The bridge exposes a stable JSON API for C/Python clients. All operations use a JSON envelope format:

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

### Supported Operations (41 ops)

**Slide Operations**: `slide_count`, `add_slide`, `remove_slide`, `move_slide`, `duplicate_slide`, `update_slide`, `list_slides`, `set_slide_title`

**Metadata**: `get_metadata`, `get_core_properties`, `set_core_properties`, `set_slide_size`, `apply_theme`, `set_modify_password`, `set_mark_as_final`

**Shapes**: `list_shapes`, `add_shape`, `remove_shape`, `update_shape`, `search_shapes`, `find_and_replace`, `move_shape_to_front`, `move_shape_to_back`

**Images/Charts**: `add_image`, `add_chart`, `list_slide_charts`, `update_chart_data`, `get_image_metadata`

**Tables**: `add_table`, `get_table`, `set_table_style`, `merge_table_cells`, `split_table_cell`, `update_table_flags`, `update_table_cell`

**Sections**: `get_sections`, `add_section`, `remove_section`, `rename_section`

**Comments**: `get_authors`, `add_author`, `get_comments`, `add_comment`, `remove_comment`

**Notes**: `get_notes`, `set_notes`

**Layout/Master**: `list_slide_layouts`, `rebind_slide_layout`, `clone_layout_master_family`

**Placeholders**: `list_placeholders`, `set_placeholder_content`

**Other**: `merge_from_file`, `batch_execute`, `add_custom_xml`, `list_custom_xml`, `remove_custom_xml`, `add_vba`

See [`docs/architecture/bridge-phase1-ops.md`](docs/architecture/bridge-phase1-ops.md) for complete specifications and [`bindings/c/README.md`](bindings/c/README.md) for C API details.

### Performance Benchmarks

- Go bridge microbench:
  - `go test ./pkg/pptx/editor -run ^$ -bench "BenchmarkBridge(Execute|JSON)" -benchmem -count=3`
- Python bridge benchmark script:
  - `uv run python scripts/smoke/python_batch_latency_benchmark.py`
- JSON profile and transport decision record:
  - [`docs/benchmarks/json_bridge_profile_2026-02-21.md`](docs/benchmarks/json_bridge_profile_2026-02-21.md)

## SmartArt Troubleshooting

If SmartArt opens but some shapes still show `[Text]` in PowerPoint:

- This is usually an older/open-locked output file, not current XML generation.
- `phldrT="[Text]"` in `data.xml` is normal placeholder metadata.
- What must be non-placeholder is text runs in `drawing.xml` (`<a:t>...</a:t>`).

### Verify a generated deck

1. Generate a fresh deck:
   ```powershell
   go run ./tmp_smartart_all_v2.go
   ```
2. Validate openability in Microsoft PowerPoint:
   ```powershell
   ./scripts/smoke/validate_with_powerpoint.ps1 -Files @('examples/output/smartart_all_layouts_v2_random.pptx')
   ```
3. Check rendered SmartArt text runs (example slide 3 / drawing2):
   ```powershell
   tar -xOf examples/output/smartart_all_layouts_v2_random.pptx ppt/diagrams/drawing2.xml | Select-String -Pattern '<a:t>[^<]*</a:t>'
   ```

If `[Text]` still appears in UI, close all PowerPoint windows and reopen only the latest generated file.

## License

Apache 2.0
