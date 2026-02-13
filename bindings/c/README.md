# gopptx C Bindings

This package provides C-compatible bindings for the `gopptx` engine, allowing it to be used from Python, C, C++, and other languages that support the C ABI.

## Features
- **Opaque Handles**: Safe memory management using `DeckHandle` (uintptr).
- **Thread Safety**: Global handle registry is thread-safe.
- **Python Support**: Easily used via `ctypes`.
- **In-Memory Editing**: Open, modify, and save without manual XML manipulation.

## Building

### Windows (PowerShell)
```powershell
.\scripts\build_cbridge.ps1
```

### Linux / macOS (Bash)
```bash
./scripts/build_cbridge.sh
```

This will generate a shared library (`gopptx.dll`, `libgopptx.so`, or `libgopptx.dylib`) and a header file in `bindings/c/build/`.

## C API Reference

### Lifecycle
- `DeckHandle deck_open(const char* path)`: Opens a PPTX file and returns a handle. Returns 0 on failure.
- `void deck_close(DeckHandle h)`: Closes the deck and frees associated resources.
- `int deck_save(DeckHandle h, const char* path)`: Saves the deck to the specified path. Returns 0 on success.

### Operations
- `const char* deck_execute_json(DeckHandle h, const char* json_input)`: Executes a JSON command envelope and returns a JSON response string. **Caller must free the returned string via `deck_free_string`**.

### Error Handling
- `const char* deck_last_error(DeckHandle h)`: Returns the last error message for the handle. **Caller must free the returned string via `deck_free_string`**.
- `void deck_free_string(const char* s)`: Frees a string allocated by the Go bridge.

## Python Example (ctypes)

```python
import ctypes
import json
import uuid

lib = ctypes.CDLL("./gopptx.dll")
h = lib.deck_open(b"example.pptx")
if h:
    req = {
        "api_version": 1,
        "request_id": str(uuid.uuid4()),
        "op": "add_slide",
        "payload": {"title": "New Slide"},
    }
    raw = lib.deck_execute_json(h, json.dumps(req).encode("utf-8"))
    lib.deck_save(h, b"modified.pptx")
    lib.deck_close(h)
```
