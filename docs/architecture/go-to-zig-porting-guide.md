# Go to Zig Porting Guide

Date: 2026-05-30

Status: Draft v0.1

This guide is the migration contract for any experimental `gopptx` Go-to-Zig port. It is intentionally modeled after the process lessons in `docs/architecture/bun-zig-to-rust-port-analysis.md`, but adapted to this project: Go core engine, C bridge, Python bindings, and PPTX/OPC correctness.

## Goal

Port `gopptx` incrementally to Zig without breaking the current Go and Python user experience.

Phase A is not a full rewrite. Phase A creates a small, testable Zig backend slice that can write a valid minimal PPTX and prove parity against Go-generated golden files.

## Non-Goals

- Do not rewrite the whole Go engine in one pass.
- Do not replace the Python API during the port.
- Do not expose Zig internals directly to Python.
- Do not change the public Go API as part of the port.
- Do not optimize before correctness and parity are proven.
- Do not add silent fallback behavior. If a ported operation is unsupported, return an explicit error.

## Hard Invariants

- Existing Go implementation remains the source of truth until a Zig feature passes parity tests.
- Python bindings must keep the same user-facing API.
- The native boundary must stay handle-based and C-compatible.
- Generated PPTX files must be valid OPC ZIP packages.
- XML output must be deterministic enough for normalized golden comparison.
- Every allocation crossing an ABI boundary must have a matching free function.
- No panic, trap, or allocator failure may cross the C ABI.
- No code file added for the port should exceed 300 lines; split modules early.

## Phase Plan

### Phase 0: Contract and Harness

Deliverables:

- This guide.
- A port manifest format.
- Golden PPTX comparison strategy.
- Decision on initial Zig directory layout.

No engine porting starts until these are clear.

### Phase 1: Minimal Zig PPTX Spike

Create a separate Zig backend that writes one valid one-slide deck with one text box.

Required package parts:

- `[Content_Types].xml`
- `_rels/.rels`
- `docProps/core.xml`
- `docProps/app.xml`
- `ppt/presentation.xml`
- `ppt/_rels/presentation.xml.rels`
- `ppt/slides/slide1.xml`
- `ppt/slides/_rels/slide1.xml.rels`
- `ppt/slideLayouts/slideLayout1.xml`
- `ppt/slideMasters/slideMaster1.xml`
- `ppt/theme/theme1.xml`

Success criteria:

- PowerPoint or LibreOffice opens the file without repair.
- Normalized ZIP entry list matches the expected minimal fixture.
- Normalized XML comparison passes against the Go-generated minimal fixture, except for explicitly allowed differences.

### Phase 2: Golden Compatibility Harness

Build tests that compare Go and Zig outputs.

Compare:

- ZIP entry set.
- OPC relationship targets and types.
- Content type defaults and overrides.
- Normalized XML.
- Media file presence and checksums once media is supported.
- Slide count, dimensions, and visible text.

The harness must fail closed. Unknown differences are failures until documented.

### Phase 3: Leaf Module Ports

Port low-risk helpers before user-facing features:

- EMU/unit conversions.
- XML escaping.
- Relationship constants.
- Content type constants.
- Package path helpers.
- Deterministic XML writer helpers.
- Minimal ZIP writer wrapper.

Each module needs unit tests or golden output tests.

### Phase 4: Stable C ABI

Introduce a Zig C ABI only after the minimal file writer is proven.

The ABI must follow the existing bridge style:

- Opaque handles only.
- Integer error codes.
- `last_error` style error retrieval.
- Explicit `free_string` / `free_bytes`.
- No exposed Go or Zig structs.
- Per-handle lifecycle: open/new, mutate, save, close.

### Phase 5: Feature Slices

Port visible features in this order:

1. Blank presentation.
2. Slide creation.
3. Plain text boxes.
4. Basic shapes.
5. Images.
6. Tables.
7. Charts.
8. Notes/comments/sections.
9. Templates.
10. Export/repair/security features.

Each slice must have a golden fixture and parity assertions before moving on.

## Directory Layout

Use a separate tree while Zig is experimental:

```text
zig/
  build.zig
  build.zig.zon
  src/
    main.zig
    gopptx.zig
    abi/
    opc/
    pptx/
    xml/
    zip/
    testing/
```

Suggested mapping:

| Go area | Zig area | Notes |
| --- | --- | --- |
| `internal/opc` | `zig/src/opc` | OPC package model, relationships, content types. |
| `internal/pptxxml` | `zig/src/pptx/xml_parts` | Generated/static XML part writers. |
| `pkg/pptx/common` | `zig/src/pptx/common` | Units, colors, dimensions, shared types. |
| `pkg/pptx/presentation` | `zig/src/pptx/presentation` | Presentation-level model and save flow. |
| `pkg/pptx/text` | `zig/src/pptx/text` | Text runs, paragraphs, rich text later. |
| `pkg/pptx/shapes` | `zig/src/pptx/shapes` | Shape type mapping and geometry. |
| `pkg/pptx/images` / `media` | `zig/src/pptx/media` | Media inventory and content-type rules. |
| `pkg/pptx/charts` | `zig/src/pptx/charts` | Later phase; high XML complexity. |
| `pkg/pptx/editor` | `zig/src/pptx/editor` | Later phase; stateful mutation semantics. |
| `bindings/c` | `zig/src/abi` | C ABI wrappers, handle registry. |

Do not mirror every Go package automatically. Keep Zig modules grouped by runtime responsibility.

## Port Manifest

Track work with a tab-separated manifest:

```text
go_path	zig_path	phase	status	owner	notes
internal/opc/...	zig/src/opc/...	3	pending	UNASSIGNED	UNCONFIRMED
```

Allowed status values:

- `pending`
- `draft`
- `compiles`
- `unit-tested`
- `golden-tested`
- `accepted`
- `blocked`

The manifest should be small at first. Do not inventory all 973 Go files before proving the approach.

## Translation Rules

### Preserve Behavior, Not Syntax

Go and Zig have different strengths. Do not translate syntax line-by-line when a small Zig-native helper preserves behavior more clearly.

However, during a feature slice, preserve:

- public behavior;
- error conditions;
- generated OPC structure;
- generated XML semantics;
- order-sensitive output when tests depend on it.

### Use Explicit Markers

Use these exact comments in Zig port files:

- `// TODO(port): reason`
- `// PERF(port): reason`
- `// SAFETY(port): invariant`
- `// PARITY(port): Go path or golden fixture`

Rules:

- Never leave an unexplained `TODO`.
- Never use `SAFETY(port)` to justify convenience. It must state the invariant.
- Any `PERF(port)` item needs a benchmark or profiling note before removal.

## Type Map

| Go | Zig | Rule |
| --- | --- | --- |
| `string` | `[]const u8` | UTF-8 by convention. Duplicate if retained beyond caller lifetime. |
| `[]byte` | `[]u8` / `[]const u8` | Mutable only when mutation is required. |
| `[]T` | `[]T` / `std.ArrayList(T)` | Use slice for borrowed views, ArrayList for owned growable buffers. |
| `map[string]T` | `std.StringHashMap(T)` | Deterministic output requires sorted keys before writing XML/ZIP. |
| `error` | error union | Define narrow error sets per module where practical. |
| `io.Reader` | explicit reader interface/wrapper | Do not build a broad abstraction until two call sites need it. |
| `io.Writer` | explicit writer interface/wrapper | XML and ZIP writers should own their write contract. |
| `time.Time` | explicit ISO/string helper | Avoid broad time dependency in Phase 1. |
| `sync.Mutex` | `std.Thread.Mutex` | Only in ABI/session registry or real shared state. |
| `context.Context` | explicit parameter or unsupported | Do not create fake cancellation. |
| struct pointer | pointer or owned value | Document owner in the field comment for long-lived pointers. |

## Allocator Policy

Every public Zig function that allocates must accept an allocator or be owned by a struct that stores one.

Rules:

- No hidden global allocator.
- No leaking by design.
- No returning allocator-owned memory without a matching free API at the ABI boundary.
- Use arena allocation only for short-lived build/write phases.
- If an arena owns a value, do not free that value individually.
- If a buffer crosses C ABI, allocate it through the ABI allocator and free it through an exported function.

Preferred pattern:

```zig
pub fn createMinimalDeck(allocator: std.mem.Allocator, title: []const u8) ![]u8 {
    // caller owns returned bytes
}
```

## Error Policy

No default fallback behavior during development.

Rules:

- Return errors explicitly.
- Keep errors narrow enough to diagnose the failing layer.
- Include path/part context in error values where practical.
- Convert all ABI-facing errors to integer code plus last-error string.
- Never panic across the ABI.

Example error sets:

```zig
pub const OpcError = error{
    InvalidPartName,
    DuplicatePart,
    MissingRelationship,
    ZipWriteFailed,
};
```

## XML Policy

PPTX correctness is mostly XML correctness.

Rules:

- XML writers must escape text and attributes through shared helpers.
- Do not hand-concatenate unescaped user text.
- Preserve namespace declarations required by Office.
- Preserve deterministic element ordering.
- Use explicit part writer functions for major PPTX parts.
- Do not introduce generic XML DOM machinery in Phase 1 unless needed by a real feature.

Allowed Phase 1 shape:

```zig
pub fn writeSlideXml(writer: anytype, slide: SlideModel) !void {
    // deterministic, direct XML writer
}
```

## ZIP and OPC Policy

PPTX is an OPC package. ZIP entry correctness is part of the API.

Rules:

- Use forward slashes in package paths.
- Do not write absolute paths.
- Do not write duplicate entries.
- Content type overrides must match part names.
- Relationship IDs must be deterministic.
- Relationship targets must be relative to the relationship part.
- ZIP timestamps should be deterministic for golden tests.

If the Zig standard library ZIP support is insufficient for deterministic PPTX writing, evaluate a small self-hosted Zig ZIP dependency before implementing a custom ZIP writer.

## C ABI Policy

The ABI is the long-term seam between Python and native code.

Use:

```c
typedef uintptr_t DeckHandle;
```

Required functions:

- `gopptx_new`
- `gopptx_open`
- `gopptx_save`
- `gopptx_save_bytes`
- `gopptx_close`
- `gopptx_last_error`
- `gopptx_free_string`
- `gopptx_free_bytes`

Rules:

- Return `0` handles on failure.
- Return non-zero error codes for mutation/save failures.
- Store handle state in a registry protected by a mutex.
- Delete handles on close.
- Treat double-close as an error, not a crash.
- Python must never receive a pointer to internal Zig memory.

## Python Binding Policy

Python should bind to a stable operation surface, not internal models.

Rules:

- Keep existing Python names and behavior.
- Add backend selection only after Go and Zig both satisfy the same API operation.
- Backend selection must be explicit during development.
- Do not silently fall back from Zig to Go when a Zig operation fails.

Development backend selection can be environment-based, for example:

```text
GOPPTX_NATIVE_BACKEND=go
GOPPTX_NATIVE_BACKEND=zig
```

If Zig is selected and unsupported, return an explicit unsupported-operation error.

## Testing Policy

Minimum validation per ported feature:

- Zig unit tests for leaf logic.
- Go-generated golden fixture.
- ZIP entry comparison.
- Normalized XML comparison.
- Openability check with at least one external renderer when feasible.
- Python API parity test once exposed to Python.

Normalization may remove:

- ZIP timestamps.
- XML formatting whitespace.
- Relationship IDs only if the tested feature declares IDs nondeterministic. Prefer deterministic IDs.

Normalization must not remove:

- element names;
- attribute names;
- relationship targets;
- content type declarations;
- visible text;
- slide dimensions.

## Security Rules

PPTX files are ZIP/XML inputs and may be attacker-controlled.

Rules:

- No path traversal when reading ZIP entries.
- No writing outside the intended output path.
- No unbounded decompression without limits.
- No XML entity expansion.
- No fetching remote resources in the core writer.
- URL fetching remains isolated in URL-specific modules and must retain existing security checks.

If a security weakness is found during porting, add a `WARNING` comment at the risky source and document the safer design before implementing.

## First Candidate Slice

Start here:

```text
Go behavior: pptx.Create(title, 1)
Zig target: create a one-slide deck with one title/text shape
Output: []u8 containing a PPTX package
Validation: golden package diff against Go output
```

Initial Zig modules:

- `zig/src/gopptx.zig`
- `zig/src/opc/package.zig`
- `zig/src/opc/content_types.zig`
- `zig/src/opc/relationships.zig`
- `zig/src/xml/escape.zig`
- `zig/src/pptx/minimal.zig`
- `zig/src/testing/golden.zig`

Keep every module small. If a file approaches 250 lines, split it before continuing.

## Review Checklist

Before accepting a ported slice:

- [ ] The Go behavior being ported is named.
- [ ] The golden fixture exists.
- [ ] The Zig code has no unexplained `TODO(port)`.
- [ ] Every allocation has an owner.
- [ ] Every ABI allocation has a free function.
- [ ] XML user text is escaped.
- [ ] ZIP paths are normalized and relative.
- [ ] Unsupported behavior returns an explicit error.
- [ ] Tests pass in the container workflow.
- [ ] `CONTINUITY.md` records the accepted slice.

## Open Decisions

- Zig version to target.
- Zig package/dependency policy for ZIP writing.
- Whether to keep Zig as a separate backend forever or eventually replace Go.
- Whether the first ABI should be bytes-only or stateful-handle based.
- Exact golden normalization rules.
