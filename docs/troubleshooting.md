# Troubleshooting

## Bridge Library Not Found

Error pattern:

- `Could not find shared library ... Please build it first.`

Fix:

1. Build library with `.\scripts\build_python.ps1`
2. Confirm platform library exists:
   - Windows: `gopptx.dll`
   - Linux: `libgopptx.so`
   - macOS: `libgopptx.dylib`
3. Set `GOPPTX_LIB_PATH` if library is not in default lookup paths.

## Batch Errors

Error pattern:

- read operations blocked in `batch()`

Cause:

- `batch()` only buffers mutating operations.

Fix:

- move read operations outside the batch block
- or issue explicit mixed command lists using `execute_batch`

## SmartArt Shows `[Text]`

Likely causes:

- older file opened in PowerPoint,
- stale PowerPoint process lock,
- validating wrong generated artifact.

Checks:

```powershell
go run ./tmp_smartart_all_v2.go
# Open the generated PPTX manually in Microsoft PowerPoint and verify text is rendered.
tar -xOf examples/output/smartart_all_layouts_v2_random.pptx ppt/diagrams/drawing2.xml | Select-String -Pattern '<a:t>[^<]*</a:t>'
```

## Python Throughput Is Slower Than Expected

- Batch operations instead of issuing one command at a time.
- Minimize cross-boundary call count.
- Optionally install `orjson` for faster Python-side encode/decode.

## Save Fails

Common causes:

- invalid/closed handle,
- path permission issues,
- earlier invalid mutation state.

Checks:

- use context manager lifecycle (`with Presentation...`)
- run `validate()` before saving in risky workflows
- verify destination path is writable

## Known Environment-Dependent Test

`pkg/pptx/presentation/protection` has a known COM-dependent scenario that can fail on some hosts. See `CONTINUITY.md` incident `I001` for current status and mitigation.
