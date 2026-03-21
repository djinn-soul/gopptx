# C Bridge Guidance

Yes, you can handle your own C binding for Go. It is a valid path, but the complexity grows quickly once you move beyond basic demos.

This guide summarizes a practical approach for a **stateful PPTX session** (`open -> edit many times -> save`) and the main pitfalls to avoid.

## Opaque Handle API (recommended)

Do not expose Go structs across the boundary. Expose only an opaque handle.

Suggested C API shape:

- `DeckHandle deck_open(const char* path);`
- `int deck_add_slide(DeckHandle h, ...);`
- `int deck_add_chart(DeckHandle h, ...);`
- `int deck_save(DeckHandle h, const char* out_path);`
- `void deck_close(DeckHandle h);`
- `const char* deck_last_error(DeckHandle h);`
- `void deck_free_string(const char* s);`

Where `DeckHandle` is an integer/uintptr key into an internal registry.

## Handle Registry Safety

Use a global map `handle -> deck` protected by a mutex.

Key rules:

- Allocate handles monotonically.
- Delete on close.
- Return explicit error codes; do not panic across C boundary.

## Error and Memory Model

At C boundaries, convert failures to:

- integer return code (`0` success, non-zero error), and
- error string retrievable via `deck_last_error`.

If returning `char*` from Go:

- allocate via `C.CString`,
- free via exported `deck_free_string` (`C.free`).

## cgo Gotchas

1. Do not store C pointers in Go structs.
2. Do not let C retain pointers to Go memory for long-term use.
3. Assume Python may call from multiple threads; either lock per handle or clearly document single-thread use.
4. Packaging/build matrix (`.dll/.so/.dylib`) is a real maintenance cost.

## PPTX Editing Runtime Model

Use an in-memory session model:

- `Open()` loads central directory and parts (lazy where possible),
- edit operations mutate in-memory model/patches,
- `Save()` writes one new PPTX in one pass.

Avoid rewriting ZIP on every operation.

## Layering Strategy

Keep two layers:

- Go engine layer (`pkg/pptx`): real business logic.
- Thin cgo adapter layer (`bindings/c` or `pkg/cbridge`): type conversion + call-through.

Keep cgo layer boring and minimal.

## When This Path Is Worth It

Good fit when:

- in-process calls are required,
- you can maintain ABI + multi-platform CI,
- session semantics are required without an extra service.

Less ideal when:

- distribution simplicity matters most (`pip install` expectations),
- support burden must stay minimal,
- background RPC process is acceptable.

## Safe API Checklist

- [ ] Opaque handles (`uintptr`)
- [ ] `open/save/close` lifecycle
- [ ] `last_error` retrieval
- [ ] `free_string` function
- [ ] panic recovery at exported boundary
- [ ] synchronized handle map + per-handle locking (or documented single-thread model)
- [ ] fuzz/corruption tests for PPTX output

## Typical Go Binary Size (CLI/Tools)

- Minimal CLI: ~1.5-2.5 MB
- Small stdlib tool: ~3-6 MB
- Medium app (HTTP/JSON/ZIP/XML): ~6-12 MB
- Larger dependency-heavy app: ~12-25 MB

For a PPTX engine, expect roughly ~8-15 MB per platform build.

## Size Optimization

Use production flags:

```bash
go build -ldflags="-s -w"
```

## Python Wheel Packaging Impact (if bundling binaries)

- Linux only: +10-15 MB
- Linux + macOS: +20-30 MB
- Linux + macOS + Windows: +30-45 MB

A practical expectation for this type of package is often 12-20 MB per platform wheel.
