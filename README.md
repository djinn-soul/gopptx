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

## License

Apache 2.0