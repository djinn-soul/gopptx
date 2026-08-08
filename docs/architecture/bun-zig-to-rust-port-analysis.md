# Bun Zig-to-Rust Port Analysis

Date: 2026-05-30

This note investigates Bun's May 2026 Zig-to-Rust rewrite and extracts practical lessons for a possible staged native-language port in `gopptx`.

## Executive Summary

Bun did not appear to approach the rewrite as a clean redesign. The public artifacts show a large, mechanically constrained translation: keep the same architecture, keep the same data structures, place Rust files next to Zig files, translate one file at a time, then make the result compile and optimize after the fact.

The important lesson is not "AI can safely rewrite a runtime." The useful lesson is that large ports become tractable when the team creates a strict migration contract:

- Preserve source structure during the first pass.
- Write explicit type, ownership, allocator, error, naming, and idiom maps.
- Batch files through a manifest-driven workflow.
- Keep TODO and PERF markers instead of guessing.
- Use the existing test suite as the behavior oracle.
- Expect a cleanup phase after tests pass.

For `gopptx`, this argues for a staged port with a stable C ABI or package-level behavior boundary, not a broad rewrite of the Go engine.

## What Happened

The primary public artifact is Bun PR `#30412`, titled "Rewrite Bun in Rust". GitHub shows it was merged by Jarred Sumner into `main` from `claude/phase-a-port` on 2026-05-14, with 6,755 commits in the PR.

In the PR opening comment, Sumner stated that the Rust rewrite passed Bun's pre-existing test suite on all platforms, fixed several memory leaks and flaky tests, reduced binary size by 3 MB to 8 MB, and produced benchmarks that were neutral to faster. He also wrote that the codebase was "otherwise largely the same": same architecture, same data structures, few third-party libraries, and no async Rust.

The earlier porting guide makes the method more concrete. `docs/PORTING.md` says Phase A is a draft `.rs` file next to the `.zig` file that captures logic faithfully and does not need to compile. Phase B then makes it compile crate by crate. The guide explicitly requires the port to match Zig structure, function names, field order, and control flow so reviewers can diff Zig and Rust side by side.

## How They Did It

### 1. They Wrote a Porting Contract First

The Bun porting guide is the critical artifact. It defines:

- File placement rules, including when a Zig file becomes `mod.rs` or `lib.rs`.
- A crate map from Zig namespaces to Rust crates.
- A type map for slices, optional pointers, errors, JavaScriptCore handles, allocators, packed structs, tagged unions, and opaque FFI handles.
- An idiom map for `defer`, `errdefer`, `comptime`, `@ptrCast`, `@intCast`, `@memcpy`, `switch`, `try`, `orelse`, loops, and other Zig patterns.
- Rules for when `unsafe` is acceptable.
- Rules for comments such as `TODO(port)` and `PERF(port)`.

That matters because the porting guide converts a vague rewrite into a repeatable translation process. It gives humans and agents the same review contract.

### 2. They Preserved Shape Before Improving Design

The guide explicitly says Phase A should keep the same structure and control flow. It even says Phase B reviewers should be able to diff `.zig` and `.rs` side by side.

This is a conservative migration pattern:

1. First preserve behavior.
2. Then make the target language compile.
3. Then reshape target-language ergonomics.
4. Then optimize.

The PR's later commit history supports this sequence. Recent commits include performance repairs, deduplication, arena and allocator work, hot-path inlining, and memory-safety fixes after the initial translation.

### 3. They Automated File Batching

`scripts/port-batch.ts` reads a manifest at `/tmp/port-manifest-filtered.tsv`, maps each `.zig` path to the expected `.rs` path, filters already-done files, and emits a JSON batch. The script supports head, tail, status, and numbered batch modes.

This is simple but important. The workflow is file-manifest driven, resumable, and measurable:

- total files
- done files
- pending files
- selected batch size
- selected batch LOC

For a port, this kind of boring accounting prevents drift.

### 4. They Used an Ownership/Lifetime Classification Step

The PR files page exposes comments around a workflow named `.claude/workflows/lifetime-classify.workflow.js`. The visible code classifies pointer fields, verifies uncertain or low-confidence classifications with multiple agent votes, synthesizes a TSV, and reports fields by ownership class.

The porting guide also references `docs/LIFETIMES.tsv` as precomputed cross-file analysis for pointer fields, with classes such as owned, shared, borrowed, static, intrusive, FFI, arena, and unknown.

That is the strongest technical signal in the process. The hard part of Zig-to-Rust is not syntax; it is ownership. Bun tried to make ownership decisions explicit data before relying on the compiler.

### 5. They Accepted Unsafe Rust Where the Source Was Unsafe

The guide says unsafe is acceptable when the Zig was already unsafe, but every unsafe block should carry a `SAFETY` comment that mirrors the Zig invariant.

This implies the port was not an immediate rewrite into idiomatic safe Rust. It was a translation into Rust as a systems language with many existing low-level invariants carried forward.

That choice probably made the port faster, but it also reduces the immediate safety benefit. A faithful unsafe port can preserve old memory hazards and add Rust-specific undefined behavior if raw pointers are translated incorrectly.

### 6. They Let Tests Be the Main Behavior Oracle

The PR comment emphasized passing the pre-existing test suite. The strategy appears to rely on existing tests to confirm compatibility while the architecture remains mostly unchanged.

That is pragmatic, but it means the quality of the port is bounded by the test suite. Passing tests is not proof of memory safety, lifetime correctness, API equivalence, or performance stability.

### 7. They Continued With Cleanup and Optimization

The PR comment explicitly says optimization and cleanup work remained before the Rust version would land outside canary. Later commits listed in the PR include many `perf:` and `dedup:` changes.

So the rewrite was not "done" at merge. It was moved into the mainline/canary feedback loop with known follow-up work.

## Risks and Criticism

The public artifacts also show real risk.

One issue, `#30719`, reported a Miri undefined-behavior case involving a dangling reference in a `PathString::slice` example. The issue body asks Bun to add Miri to CI and shows Rust detecting undefined behavior from constructing an invalid `&[u8]`.

This is exactly the class of failure a mechanical systems-language port can introduce:

- a source-language pointer pattern is copied too directly;
- the target language gives references stronger validity assumptions;
- the code compiles and may pass ordinary tests;
- dynamic UB tooling catches what normal tests miss.

The takeaway is not that the Bun rewrite is invalid. The takeaway is that a port needs specialized verification beyond unit tests.

## Lessons for `gopptx`

`gopptx` is Go core plus C bridge plus Python bindings. A Zig port would differ from Bun's direction: we would be moving from Go to Zig, not Zig to Rust. That changes the hard parts.

### What Transfers Well

- Use a porting guide before porting code.
- Preserve behavior and public APIs first.
- Port by small slices.
- Keep source and target modules easy to compare.
- Use a manifest to track progress.
- Treat tests and generated `.pptx` files as the behavior oracle.
- Add explicit markers for uncertain translations.
- Separate correctness from optimization.

### What Does Not Transfer Directly

- Rust ownership maps are not Zig ownership maps. For Zig, the hard problem is allocator discipline and lifetime documentation, not satisfying a borrow checker.
- Bun carried low-level JSC and runtime constraints. `gopptx` has a different core risk: valid Office Open XML packaging, relationships, content types, ZIP structure, and Python/Go API parity.
- Bun could accept a large unsafe Rust translation because it already had a runtime test suite. For `gopptx`, the decisive evidence should be generated PPTX compatibility, XML/package diffs, and Python binding parity.

## Recommended `gopptx` Port Strategy

Do not start with a whole-engine rewrite. Start with a proof-oriented staged port.

### Phase 0: Port Contract

Create `docs/architecture/go-to-zig-porting-guide.md` before writing production Zig. It should define:

- Go package to Zig module mapping.
- Public API boundaries that must not move.
- C ABI expectations.
- Error model.
- Allocator policy.
- String and byte ownership policy.
- ZIP/XML writer policy.
- Test oracle and golden file rules.
- Allowed third-party Zig libraries, if any.
- `TODO(port)` and `PERF(port)` marker rules.

### Phase 1: Minimal Zig PPTX Spike

Build a separate Zig spike that writes one valid `.pptx`:

- `[Content_Types].xml`
- `_rels/.rels`
- `ppt/presentation.xml`
- `ppt/_rels/presentation.xml.rels`
- one slide XML
- one text box

Success means PowerPoint/LibreOffice can open it and the package structure matches a known-good Go-generated file.

### Phase 2: Golden Compatibility Harness

Generate small decks from the Go implementation and compare:

- ZIP entry names
- required relationship IDs
- XML normalized output
- visible slide content
- metadata behavior
- images and media once those features are reached

Do this before broad translation. Without golden tests, a port becomes guesswork.

### Phase 3: Leaf Modules

Port low-risk, dependency-light modules first:

- EMU/unit conversions
- XML escaping
- content type constants
- relationship constants
- path helpers
- simple ZIP package writer abstractions

Each module should have Go-vs-Zig parity tests or golden output checks.

### Phase 4: C ABI Boundary

If the end goal is Python using a Zig native backend, define the C ABI early. Do not let Python bind directly to unstable Zig internals.

The ABI should expose stable operations, not internal structs. That lets Go and Zig coexist behind the same Python surface during migration.

### Phase 5: Feature Slices

Port visible features one by one:

1. blank presentation
2. slide creation
3. text boxes
4. basic shapes
5. images
6. tables
7. charts
8. templates/markdown/html/PDF-adjacent features later

Each feature should be merged only when it has a golden deck and a parity assertion.

## Practical Recommendation

Use Bun as a process model, not as proof that a huge rewrite is safe.

For `gopptx`, the best next step is a small Zig spike plus a written porting guide. Do not attempt a broad Go-to-Zig rewrite until the project has:

- a stable native ABI plan;
- golden deck comparison tests;
- a module manifest;
- explicit allocator and ownership rules;
- a decision on whether Go and Zig coexist during migration.

If we do this, a step-by-step Zig port is realistic. If we skip those controls, the port will likely create two engines with unclear parity.

## Sources

- Bun PR `#30412`, "Rewrite Bun in Rust": https://github.com/oven-sh/bun/pull/30412
- Bun porting guide at commit `3157cb14b5970b69532a47800504a28ef5963e22`: https://github.com/oven-sh/bun/blob/3157cb14b5970b69532a47800504a28ef5963e22/docs/PORTING.md
- Raw Bun porting guide: https://raw.githubusercontent.com/oven-sh/bun/3157cb14b5970b69532a47800504a28ef5963e22/docs/PORTING.md
- Bun file batching script: https://github.com/oven-sh/bun/blob/3157cb14b5970b69532a47800504a28ef5963e22/scripts/port-batch.ts
- Bun PR files page showing lifetime classification workflow review context: https://github.com/oven-sh/bun/pull/30412/files
- Bun issue `#30719`, Miri/dangling-reference UB report: https://github.com/oven-sh/bun/issues/30719
- Zig Code of Conduct, strict no LLM/no AI policy: https://ziglang.org/code-of-conduct/
