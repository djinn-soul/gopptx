# Installation

## Choose Runtime

- Go-only use: install module and use Go APIs directly
- Python use: build bridge shared library and install Python package
- Docs development: use pinned MkDocs Material versions

## Go Setup

```bash
go get github.com/djinn-soul/gopptx
```

## Python Setup (from this repo)

1. Build shared library:

```powershell
.\scripts\build_python.ps1
```

2. Install package:

```bash
pip install -e .
```

3. Optional performance package:

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
