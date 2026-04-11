# Troubleshooting

## Bridge Library Not Found

**Error pattern:**

```
Could not find shared library ... Please build it first.
```

**Fix:**

1. Build the shared library for your platform:

    === "Windows"
        ```powershell
        .\scripts\build_python.ps1
        ```
    === "Linux / macOS"
        ```bash
        ./scripts/build_python.sh
        ```

2. Confirm the platform library exists in the repo root:

    | Platform | File |
    |---|---|
    | Windows | `gopptx.dll` |
    | Linux | `libgopptx.so` |
    | macOS | `libgopptx.dylib` |

3. If the library is in a non-default location, set:
    ```bash
    export GOPPTX_LIB_PATH=/path/to/library
    ```

---

## Import Error: `gopptx` Not Found

**Error pattern:**

```
ModuleNotFoundError: No module named 'gopptx'
```

**Fix:**

Install the Python package from the repo root after building the shared library:

```bash
pip install -e .
```

---

## Batch Errors

**Error pattern:**

```
RuntimeError: read operations are not allowed inside a batch block
```

**Cause:**

`batch()` only buffers mutating operations. Read calls (e.g. `get_slide_count()`) are blocked inside a batch block.

**Fix:**

- Move read operations outside the `with pres.batch():` block.
- Or use `execute_batch()` with an explicit mixed command list for interleaved reads and writes.

---

## SmartArt Shows `[Text]`

**Likely causes:**

- Opened an older file in PowerPoint before the engine finished writing.
- A stale PowerPoint process still holds a file lock.
- Validating the wrong generated artifact.

**Fix:**

1. Close all PowerPoint windows before running the generator.
2. Open the newly generated PPTX file directly — do not reuse a previously open window.
3. If the problem persists, inspect the XML:
    ```bash
    # Extract the diagram drawing XML to check for text nodes
    unzip -p output.pptx ppt/diagrams/drawing2.xml | grep -o '<a:t>[^<]*</a:t>'
    ```

---

## Python Throughput Is Slower Than Expected

**Cause:**

Each Python → C → Go boundary crossing has overhead. Issuing one operation at a time compounds this for large decks.

**Fix:**

- Use `execute_batch()` or the `with pres.batch():` context manager to send many operations in a single crossing.
- Minimize the total number of cross-boundary calls in write-heavy workflows.
- Optionally install `orjson` for faster Python-side JSON encode/decode:
    ```bash
    pip install orjson
    ```

See [Batch Execution](guides/batch-execution.md) for patterns.

---

## Save Fails

**Error pattern:**

```
RuntimeError: save failed — handle may be closed or path is not writable
```

**Common causes:**

- Handle was already closed (e.g. `close()` called before `save()`).
- Destination path does not exist or is not writable.
- An earlier invalid mutation left the presentation in a bad state.

**Fix:**

- Always use the context manager so handle lifetime is deterministic:
    ```python
    with Presentation.new("My Deck") as pres:
        pres.add_slide("Intro")
        pres.save("output.pptx")
    ```
- Run `pres.validate()` before saving in workflows that do conditional mutations.
- Verify the destination directory exists and is writable.

---

## Protection / COM Errors on Some Hosts

`pkg/pptx/presentation/protection` includes a scenario that requires COM (Windows-only). This test is skipped automatically on Linux and macOS. If you see unexpected failures on Windows, ensure the host has Microsoft Office or the OpenXML SDK installed.
