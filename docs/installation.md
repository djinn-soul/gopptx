# Installation

## Prerequisites

- **Go**: 1.25.9 or later
- **Python**: 3.10 or later (required for Python bindings)
- **Platforms**: Windows, Linux, macOS

## Choose Runtime

- Go-only use: install module and use Go APIs directly
- Python use: build bridge shared library and install Python package
- Docs development: use pinned MkDocs Material versions

## Go Setup

```bash
go get github.com/djinn-soul/gopptx
```

## Python Setup (from this repo)

The Python package uses a high-performance Go shared library for PPTX processing. You need to build this library before installing the package.

1. **Build the Go shared library:**
   - **Windows (PowerShell):**
     ```powershell
     .\scripts\build_python.ps1
     ```
   - **Linux/macOS (Bash):**
     ```bash
     ./scripts/build_python.sh
     ```

2. **Install the Python package:**
   ```bash
   pip install -e .
   ```

3. **Optional performance package:**
   ```bash
   pip install orjson
   ```

## Environment Notes

- Shared library filename is platform-specific:
  - Windows: `gopptx.dll`
  - Linux: `libgopptx.so`
  - macOS: `libgopptx.dylib`
- You can override lookup path with `GOPPTX_LIB_PATH`.

## Docs Tooling (Pinned)

This repo pins:

- `mkdocs==1.6.1`
- `mkdocs-material==9.7.5`

Use repo tasks:

- `task docs:serve`
- `task docs:build`
