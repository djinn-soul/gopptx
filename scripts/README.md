# Scripts

## What These Scripts Do

- `build_python.*`
  - Builds the Go shared library and copies it into `python/gopptx/`.
  - Use this for local Python package development/testing.
- `build_cbridge.*`
  - Builds the same C-shared bridge into `bindings/c/build/`.
  - Also syncs the produced library into `python/gopptx/`.
- `architectural_guardrails.py`
  - Enforces project structural limits/rules used in CI.
- `architectural_guardrails.json`
  - Baseline/config consumed by the guardrails script.

## Python Bridge Build

- Dev build (default, debug-friendly size):
  - Windows: `./scripts/build_python.ps1`
  - Linux/macOS: `./scripts/build_python.sh`

- Release build (smaller binary):
  - PowerShell: `$env:GOPPTX_RELEASE_BUILD="1"; ./scripts/build_python.ps1`
  - Bash: `GOPPTX_RELEASE_BUILD=1 ./scripts/build_python.sh`

Release mode enables:
- `-trimpath`
- `-buildvcs=false`
- `-ldflags "-s -w"`

Expected Python artifact names:
- Windows: `python/gopptx/gopptx.dll`
- Linux: `python/gopptx/libgopptx.so`
- macOS: `python/gopptx/libgopptx.dylib`

Tip:
- Use dev build for debugging/local iteration.
- Use release build for packaging/publishing.

## CI Guardrails

- `uv run python scripts/ci/architectural_guardrails.py`
- `uv run python scripts/ci/architectural_guardrails.py --write-current-baseline`
