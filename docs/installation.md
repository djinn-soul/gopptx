# Installation

## Prerequisites

| | Requirement | Needed for |
|---|---|---|
| Go | 1.25 or later | Both paths |
| Python | 3.10 or later | Python package only |
| C toolchain | gcc/clang (Linux, macOS) or MinGW-w64 / MSVC (Windows) | Building the shared library for Python |
| Platforms | Windows, Linux, macOS | — |

The Go package alone needs no C toolchain. The toolchain is required only because the Python
package is served by a cgo `c-shared` build.

## Go

```bash
go get github.com/djinn-soul/gopptx
```

```go
import "github.com/djinn-soul/gopptx/pkg/pptx"
```

No native dependency, no build step. Jump to the [Quickstart](quickstart.md).

## Python

The Python package calls into a Go shared library. Build the library first, then install the
package.

=== "Windows (PowerShell)"

    ```powershell
    .\scripts\build_python.ps1
    pip install -e .
    ```

=== "Linux / macOS"

    ```bash
    ./scripts/build_python.sh
    pip install -e .
    ```

=== "Task runner"

    ```bash
    task build:go        # builds the shared library into python/gopptx/
    pip install -e .
    ```

Verify the install:

```bash
python -c "import gopptx; print(gopptx.__version__)"
```

### Optional: faster JSON

Every bridge call encodes a JSON envelope. `orjson` is used automatically when present:

```bash
pip install orjson
```

### Where the shared library lives

The build writes a platform-specific file into `python/gopptx/`:

| Platform | Filename |
|---|---|
| Windows | `gopptx.dll` |
| Linux | `libgopptx.so` |
| macOS | `libgopptx.dylib` |

If you keep the library elsewhere — a shared volume, a system path, a packaged wheel layout —
point the loader at it:

```bash
export GOPPTX_LIB_PATH=/opt/gopptx/libgopptx.so     # Linux / macOS
$env:GOPPTX_LIB_PATH = "C:\opt\gopptx\gopptx.dll"   # Windows
```

If the library cannot be found you get a `Could not find shared library` error — see
[Troubleshooting](troubleshooting.md#bridge-library-not-found).

### Rebuild after changing Go code

The Python package does **not** pick up Go changes until the library is rebuilt:

```bash
task build:go
```

If the change touched `bindings/c/bridge.go`, the generated `python/gopptx/gopptx.h` changes with
it and must be staged alongside it, or the pre-commit hooks fail with a confusing error.

## Optional external tools

These are only needed for specific features. Nothing about authoring a `.pptx` requires them.

| Tool | Enables | Install note |
|---|---|---|
| LibreOffice | `libreoffice` PDF driver | `soffice` must be on `PATH`. On Windows add `C:\Program Files\LibreOffice\program`. |
| Microsoft PowerPoint | `powerpoint` PDF driver | Windows only, via COM automation; needs `powershell` or `pwsh`. |
| Docker | `task docs:serve` / `task docs:build` | The docs site is built in the pinned MkDocs Material image. |

## Development setup

```bash
task setup:python    # uv sync of the dev environment
task build:go        # shared library
task test            # Go + Python suites
task lint            # golangci-lint, ruff, basedpyright, generated-code drift
```

Docs tooling is pinned — `mkdocs==1.6.1`, `mkdocs-material==9.7.5` — and runs through Docker so
the local build matches CI:

```bash
task docs:serve      # http://localhost:8000
task docs:build      # strict build; fails on broken links
```

See [`CONTRIBUTING.md`](https://github.com/djinn-soul/gopptx/blob/main/CONTRIBUTING.md) for the
full contributor workflow.
