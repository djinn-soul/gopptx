# gopptx

High-performance PowerPoint (PPTX) engine powered by Go.

## Features
- **Fast**: Blazing fast PPTX generation and manipulation.
- **Concurrent**: Optimized for multi-threaded slide processing.
- **Lazy Loading**: Efficient memory usage for large presentations.
- **Cross-Language**: C-compatible bindings for Python and other languages.

## Installation (Go)
```bash
go get github.com/djinn-soul/gopptx
```

## Python Library

`gopptx` can be used as a high-performance Python library. It provides a high-level API similar to `python-pptx` but utilizes the Go engine for heavy-duty operations.

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
