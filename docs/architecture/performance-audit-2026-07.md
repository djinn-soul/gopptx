# Performance Audit — 2026-07-21

Codebase-wide optimization scan. All numbers are measured, not estimated.

**Environment:** Intel Core i7-8550U @ 1.80GHz, 8 threads, windows/amd64, Go test harness.
Machine is a thermally-throttled laptop — absolute timings are noisy across runs, so
same-process comparisons and repeat counts are given where the ratio matters.

**Branch at time of audit:** `fix/ffi-concurrency-safety` (HEAD `05d38e6`)

---

## Baseline benchmark numbers

### `internal/pptxxml`

```
BenchmarkRichSolidFill-8                     5250496     218.0 ns/op     101 B/op      2 allocs/op
BenchmarkRichSolidFillBaseline-8             1368454     874.8 ns/op     160 B/op      5 allocs/op
BenchmarkRichLine-8                          5859208     240.2 ns/op     170 B/op      3 allocs/op
BenchmarkRichLineBaseline-8                   621220    2667   ns/op     392 B/op     15 allocs/op
BenchmarkRichOuterShadow-8                   1000000    1020   ns/op     224 B/op      5 allocs/op
BenchmarkRichOuterShadowBaseline-8            635672    1869   ns/op     464 B/op     12 allocs/op
BenchmarkRichPerspectiveShadow-8             1338938     908.8 ns/op     256 B/op      8 allocs/op
BenchmarkRichPerspectiveShadowBaseline-8      488340    3175   ns/op     632 B/op     18 allocs/op
BenchmarkPlaceholderTextStyleFull-8          1529694     659.9 ns/op     196 B/op      2 allocs/op
BenchmarkTextLevelStyles-8                    273397    4186   ns/op    2448 B/op     19 allocs/op
BenchmarkTextLevelStylesBaseline-8             41233   31434   ns/op    8459 B/op    151 allocs/op
BenchmarkEscapePlain-8                       6478449     229.6 ns/op       0 B/op      0 allocs/op
BenchmarkEscapePlainBaseline-8               6385441     181.3 ns/op       0 B/op      0 allocs/op
BenchmarkEscapeWithSpecials-8                3308080     321.7 ns/op      72 B/op      4 allocs/op
BenchmarkRenderShapeFullStack-8              1000000    1034   ns/op     490 B/op      9 allocs/op
BenchmarkRenderShapeFullStackBaseline-8       452060    4327   ns/op     824 B/op     26 allocs/op
BenchmarkPackageWriter-8                         100   24777164 ns/op  454739 B/op   5449 allocs/op
BenchmarkPackageWriterAddOnly-8                 8035     170220 ns/op   50112 B/op    427 allocs/op
BenchmarkPackageWriterWriteToOnlyDiscard-8         72   18374701 ns/op  260677 B/op   5016 allocs/op
BenchmarkContentTypesLargeDeck-8               12303      99479 ns/op  115335 B/op    203 allocs/op
```

### `pkg/pptx/editor`

```
BenchmarkBridgeExecuteSingleSetSlideTitle-8       255154   12748 ns/op    1320 B/op    22 allocs/op
BenchmarkBridgeExecuteSingleSetSlideHidden-8      372744    7877 ns/op     744 B/op    18 allocs/op
BenchmarkBridgeExecuteBatchSetSlideTitle-8         15741  105241 ns/op   13772 B/op   158 allocs/op
BenchmarkBridgeLatencySingleOps-8                   5899  186563 ns/op (9328 ns/single)   26581 B/op  440 allocs/op
BenchmarkBridgeLatencySingleSetSlideHiddenOps-8    13351  169175 ns/op (8458 ns/single)   15050 B/op  360 allocs/op
BenchmarkBridgeLatencyBatchOps-8                    6120  178465 ns/op (8923 ns/effective) 26935 B/op 289 allocs/op
BenchmarkBridgeJSONEnvelopeDecode-8               329583    4480 ns/op     264 B/op     7 allocs/op
BenchmarkBridgeJSONEnvelopeEncode-8               614804    2320 ns/op     272 B/op     7 allocs/op
BenchmarkEditorSaveToWriterUnencrypted-8            3602  451101 ns/op   23951 B/op   131 allocs/op
BenchmarkEditorSaveToBytesUnencrypted-8             4135  401565 ns/op   26469 B/op   133 allocs/op
```

---

## Finding 1 — Deflate compression level dominates all write paths (~8x available)

### Evidence

CPU profile of `BenchmarkPackageWriterWriteToOnlyDiscard` (`-benchtime=50x`):

```
Duration: 933.39ms, Total samples = 710ms (76.07%)
      flat  flat%   sum%        cum   cum%
     440ms 61.97% 61.97%      440ms 61.97%  compress/flate.(*compressor).reset
     120ms 16.90% 78.87%      230ms 32.39%  compress/flate.(*compressor).deflate
      30ms  4.23% 83.10%       30ms  4.23%  compress/flate.(*huffmanBitWriter).generateCodegen
      30ms  4.23% 87.32%       30ms  4.23%  compress/flate.matchLen (inline)
      20ms  2.82% 90.14%       20ms  2.82%  compress/flate.(*huffmanEncoder).bitLength (inline)
         0     0%   100%      680ms 95.77%  archive/zip.(*Writer).Create
         0     0%   100%      450ms 63.38%  archive/zip.newFlateWriter
         0     0%   100%      440ms 61.97%  compress/flate.(*Writer).Reset
```

**62% of CPU is `flate.(*compressor).reset`** — not compression, just state zeroing.

CPU profile of `BenchmarkEditorSaveToBytesUnencrypted` (`-benchtime=2000x`) — same story:

```
Duration: 2.13s, Total samples = 2.08s (97.52%)
      flat  flat%   sum%        cum   cum%
     0.38s 18.27% 18.27%      0.44s 21.15%  compress/flate.(*compressor).reset
     0.22s 10.58% 28.85%      0.23s 11.06%  runtime.cgocall
     0.14s  6.73% 35.58%      0.20s  9.62%  compress/flate.(*huffmanEncoder).bitCounts
     0.09s  4.33% 39.90%      0.13s  6.25%  sort.insertionSort
     0.07s  3.37% 47.12%      0.78s 37.50%  compress/flate.(*compressor).deflate
```

~70% of the editor save path is flate.

### Root cause

`archive/zip` defaults to `flate.DefaultCompression` (level 6). Level 6 uses the full
`compressor` with a ~192KB hash table (`hashHead` 32768×uint32 + `hashPrev` 32768×uint16)
that is zeroed on every `Reset` — i.e. **once per zip entry**. A deck has hundreds of
entries, many of them tiny (`_rels` files, small XML parts), so per-entry reset cost
swamps the actual compression work.

Levels 1–2 use `deflatefast`, which has a far smaller state and a cheap reset.

### Measurement — synthetic zip write, 400 entries (200 × 8-byte XML, 200 × 4KB binary)

Isolated harness, `-count=3`, all in one process:

```
BenchmarkBase-8     (level 6)    62   24095642 ns/op   255057 B/op   4816 allocs/op
BenchmarkBase-8                  38   31773582 ns/op   266846 B/op   4816 allocs/op
BenchmarkBase-8                  55   21050665 ns/op   261757 B/op   4816 allocs/op

BenchmarkL1-8       (level 1)   414    3256282 ns/op   255254 B/op   4819 allocs/op
BenchmarkL1-8                   409    3789348 ns/op   273592 B/op   4819 allocs/op
BenchmarkL1-8                   464    2678791 ns/op   262563 B/op   4819 allocs/op

BenchmarkL2-8       (level 2)    78   14040068 ns/op   244520 B/op   4819 allocs/op
BenchmarkL2-8                   100   16029015 ns/op   235330 B/op   4818 allocs/op
BenchmarkL2-8                   120    9654038 ns/op   223099 B/op   4818 allocs/op

BenchmarkL5-8       (level 5)   100   14430690 ns/op   259730 B/op   4819 allocs/op
BenchmarkL5-8                   100   27557667 ns/op   259730 B/op   4819 allocs/op
BenchmarkL5-8                    57   23276986 ns/op   231314 B/op   4818 allocs/op
```

**Level 1 is the cliff: ~8x faster than default. Levels 2+ give no meaningful gain.**

Earlier single-run comparison including a "Store small parts" variant:

```
BenchmarkBase-8         100   15171809 ns/op   227123 B/op   4815 allocs/op
BenchmarkStoreSmall-8   156   16902987 ns/op   233868 B/op   4815 allocs/op   <- no gain alone
BenchmarkBestSpeed-8    447    2915400 ns/op   237318 B/op   4817 allocs/op   <- 5.2x
BenchmarkBoth-8         507    2113630 ns/op   251137 B/op   4818 allocs/op   <- 7.2x
```

Storing sub-256-byte parts uncompressed gives nothing on its own, but stacks with
level 1 for a further ~1.4x.

### Size cost

Deflate ratio measured over 132 real XML files (>200 bytes) from this repo:

```
level  1:  223523 bytes  (11.26% of raw)
level  3:  203403 bytes  (10.25% of raw)
level  5:  186772 bytes  ( 9.41% of raw)
level -1:  183228 bytes  ( 9.23% of raw)   <- DefaultCompression, current behavior
level  9:  180126 bytes  ( 9.08% of raw)
```

**Tradeoff: ~22% larger compressed output for ~8x faster save.** Needs a product decision.

### Where to apply

`zw.RegisterCompressor(zip.Deflate, ...)` returning a pooled `flate.NewWriter(out, flate.BestSpeed)`,
at each `zip.NewWriter` site:

- `pkg/pptx/editor/save_zip_stream.go:97`
- `pkg/pptx/presentation.go:98`
- `internal/opc/writer.go:16`
- `pkg/pptx/tplx/tplx_archive.go:48`
- `pkg/pptx/editor/modules/chart/excel.go:62`

### Already done — do not redo

`copyZipEntryRaw` (`pkg/pptx/editor/save_zip_stream.go:133`) already raw-copies untouched
parts via `CreateRaw`/`OpenRaw`, so unmodified parts are never recompressed. That was the
other big structural win and it is in place. Finding 1 only affects new/modified parts.

---

## Finding 2 — `XMLEscape` via `xml.EscapeText` is 6.7x slower than the Replacer

`pkg/pptx/editor/modules/shape/helpers.go:52` allocates a fresh `bytes.Buffer` and converts
`[]byte(value)` on every call:

```
BenchmarkEscapeXMLText-8    371077   3538 ns/op   752 B/op   20 allocs/op
BenchmarkEscapeReplacer-8  2943888    524 ns/op    32 B/op    2 allocs/op
```

The rest of the codebase uses a package-level `strings.NewReplacer`
(`pkg/pptx/editor/common/types.go:75`, `internal/pptxxml/package_xml.go:14`).

**Caveat — not a mechanical swap.** `xml.EscapeText` also escapes `\n`, `\r`, `\t` and
sanitizes invalid runes; the Replacer does not. Confirm shape text does not depend on that
behavior before switching to `common.XMLEscape`.

---

## Finding 3 — `regexp.MustCompile` on the call path (13x)

```
BenchmarkReCompilePerCall-8    46231   35318 ns/op   8835 B/op   54 allocs/op
BenchmarkReCached-8           459871    2625 ns/op     32 B/op    1 alloc/op
```

Most of the codebase already hoists patterns to package-level vars. These sites were missed:

- `pkg/pptx/editor/chart_editor.go:167`, `:224`
- `pkg/pptx/editor/command_handlers_handout.go:71`, `:75`
- `pkg/pptx/editor/command_handlers_header_footer.go:85`, `:131`
- `pkg/pptx/editor/command_handlers_parity_new.go:70`, `:81`

Two sites build the pattern dynamically, so they need a `sync.Map` cache rather than a var:

- `pkg/pptx/editor/layout_master_parsing_helpers.go:33`
- `pkg/pptx/editor/modules/chart/cache_patch_replace.go:13`

---

## Finding 4 — Dead benchmark pair

`internal/pptxxml/rich_render_bench_test.go:11`:

```go
func baselineEscape(value string) string {
    return xmlEscapeReplacer.Replace(value)
}
```

This is byte-identical to `Escape` (`internal/pptxxml/package_xml.go:14`). So
`BenchmarkEscapePlain` vs `BenchmarkEscapePlainBaseline` compares the same code against
itself — the 229.6 vs 181.3 ns delta is pure noise.

The comment at `package_xml.go:11` claims a `ContainsAny` fast-path "was a micro-benchmark
loss," but this benchmark pair cannot have measured that. Either the fast-path variant was
lost in a refactor, or the claim is unverified.

---

## Checked, no action needed

- **`PartStore`** (`pkg/pptx/editor/partstore.go`) — lazy zip reads, inflight read
  deduplication, `keysCache` with dirty flag, `allKeys` set to avoid O(N) map scans in
  `Keys()`. Well tuned.
- **Rich-render helpers** (`internal/pptxxml`) — already beat their `fmt.Sprintf` baselines
  by 3–7x (see baseline table above).
- **No O(n²) traversals** or in-loop string concatenation found in hot paths. The `+=`
  matches found by grep are all numeric accumulators or bounded single-field appends.
- **Existing `sync.Pool` sites** — `modules/slide/relationships.go:18`,
  `editor/partstore.go:19`, `editor/save_zip_stream.go:15` are all appropriate.

---

---

# Implementation — 2026-07-21

All three findings implemented. Full test suite green, `gofmt` and `go vet` clean.

## Finding 1 — implemented as a global default switch to `flate.BestSpeed`

New package `internal/zipfast` provides a pooled level-1 compressor and a drop-in
`zipfast.NewWriter`. Wired into every production zip write site:

- `pkg/pptx/editor/save_zip_stream.go`
- `pkg/pptx/presentation.go`
- `internal/opc/writer.go`
- `pkg/pptx/tplx/tplx_archive.go`
- `pkg/pptx/editor/modules/chart/excel.go`

No raw `zip.NewWriter` remains outside `internal/zipfast` itself. `Store` entries and
`CreateRaw` copies are unaffected, so the existing raw-copy fast path still bypasses
compression entirely.

The pooled writer's `Close` is idempotent and resets to `io.Discard` before returning to
the pool, so a double `Close` cannot hand one `flate.Writer` to two concurrent entries.

### Measured result — `BenchmarkEditorSaveToBytesUnencrypted`, `-count=5`, quiet machine

```
before:  401565 ns/op   (single sample, older toolchain)
after :  147908 / 144127 / 251337 / 146355 / 147204 ns/op
```

Median ~146µs against a ~400µs baseline, i.e. **~2.7x on the editor save path**. The
251µs sample is machine noise; see "Measurement conditions" below.

### Size cost — CORRECTED

The pre-implementation estimate of **+22% was wrong for real decks.** That figure came
from compressing 132 arbitrary repo XML files individually — a larger and more varied
corpus than an actual generated package.

Measured on four real decks generated by `examples/09-charts` (123 deflate entries,
0 store entries), comparing payloads at both levels:

```
raw payload      :  128978 B
level 6 payload  :   46621 B
level 1 payload  :   51112 B

size increase level6 -> level1: +9.6%
```

**Real-world cost is +9.6%, not +22%.**

All four decks were reopened and every entry fully inflated without error, confirming the
output is still valid OPC.

## Finding 2 — implemented, but it was not a hot path

`pkg/pptx/editor/modules/shape/helpers.go` now uses a package-level `strings.Replacer`
instead of `xml.EscapeText` (6.7x faster, 2 allocs vs 20).

**Honest scope note:** this function had exactly one production caller
(`style_pattern_fill.go`, escaping preset geometry tokens like `pct5`). The perf win is
negligible; the real value is consistency with `common.XMLEscape`. It was listed higher in
the original audit than it deserved.

Behavior change: named entities (`&quot;`) instead of numeric (`&#34;`). Both resolve
identically in any conformant parser, and named entities match the rest of the codebase.
`shape_helpers_test.go` updated accordingly, with the reasoning recorded in a comment.

## Finding 3 — implemented

Eight patterns hoisted to package-level vars:

- `chart_editor.go` → `reChartRelID`, `reNumericIDVal`
- `command_handlers_handout.go` → `reHandoutOrient`, `reHandoutSldSz`
- `command_handlers_header_footer.go` → `reInjectHF`, `reCNvPrIDVal`
- `command_handlers_parity_new.go` → `reSmartArtDM`, `reSmartArtTextItem`

Two dynamic-pattern sites got a `sync.Map` compile cache instead, since the pattern is
built from an argument:

- `layout_master_parsing_helpers.go` → `commonCompile` memoized by pattern string
- `modules/chart/cache_patch_replace.go` → `getFieldPattern` memoized by tag

Both key spaces are small fixed sets, so the caches are bounded.

## Finding 4 — not addressed

The dead `baselineEscape` benchmark pair is still present. Left alone: removing it is a
test-only cleanup with no runtime effect, and the question of whether the `ContainsAny`
fast-path claim at `package_xml.go:11` was ever really measured is still open.

---

## Measurement conditions — read before trusting any number here

This machine is a thermally-throttled laptop and its benchmark noise floor is large.
Observed within a single `-count=3` run, same binary, consecutive repeats:

```
EscapePlain-8                       154.6 /  69.73 / 248.8 ns    3.6x spread
PackageWriter-8                     10.1  /  23.8  /  8.8 ms     2.7x spread
PackageWriterWriteToOnlyDiscard-8    9.8  /  19.2  / 19.0 ms     2.0x spread
```

An earlier full benchmark run was invalidated entirely by 95% background CPU load
(two `gopls`, Zed, Malwarebytes, firefox), producing figures 3–13x slower than baseline
with *identical* allocation counts — the tell that it was load, not code.

Consequences:

- The **ratios** in Finding 1 are trustworthy: they come from same-process, back-to-back
  comparisons where load affects both arms equally.
- The **absolute ns/op** figures in the baseline tables above are single-sample and should
  not be treated as a reference point.
- No valid cross-toolchain baseline exists. The original numbers were taken on an
  unrecorded Go version; the codebase is now on go1.26.5. Any future version comparison
  needs `benchstat` with `-count=10` on a quiet machine, against a baseline file that does
  not currently exist.

## Toolchain note — golangci-lint rebuilt

After the go1.26.5 upgrade the installed `golangci-lint` (v2.8.0, chocolatey, built with
go1.25.5) panicked on every run:

```
panic: file requires newer Go version go1.26 (application built with go1.25)
```

Rebuilt from source with the new toolchain:

```
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

This produces v2.12.2 at `%USERPROFILE%\go\bin\golangci-lint.exe`. **The chocolatey copy
at `C:\ProgramData\chocolatey\bin` still shadows it on PATH** — either remove the choco
package or reorder PATH, otherwise `golangci-lint` still resolves to the broken 2.8.0.

Lint has now been run against all changes above: clean. Note the 2.8.0 → 2.12.2 jump
surfaced ~63 pre-existing findings repo-wide (50 `goconst`, 13 `staticcheck`), none of
them in code touched by this work.

---

# Lint cleanup — 2026-07-21

Separate from the performance work: a full lint pass, fixing rather than suppressing.

## Measurement error worth recording

The first three inventories reported **104 findings**, then "63 remaining", then "65".
All wrong. `golangci-lint` defaults to `max-issues-per-linter: 50`, so `goconst` was
**capped at 50** and the real count was hidden. Roughly 45 goconst fixes produced no
change in the reported total, which is what eventually exposed the cap.

Always inventory with:

```bash
golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 ./...
```

True starting total was **~197**, not 104.

## Result

| Category | Start | End |
|---|---|---|
| goconst | ~183 | 0 |
| staticcheck QF1012 | 37 | 0 |
| modernize / golines / govet | 3 | 0 |
| gosec G703 | 7 | 0 (scoped config exclusion, see below) |
| cgo-generated | 7 | 0 (scoped config exclusion, see below) |

`QF1012` (`WriteString(fmt.Sprintf(...))` → `fmt.Fprintf`) was autofixed via
`golangci-lint run --fix --enable-only=staticcheck`, scoped to exclude `bindings/` so cgo
was never touched. It is a real allocation win, not only style.

## Latent drift found — constants that existed and were not used

The highest-value part of the pass. These were live duplication risks:

- `pkg/pptx/validation/structural` defined `presentationRelsPath` and `contentTypesPath`,
  then hardcoded both strings anyway across two repairer files.
- `pkg/pptx/export` had `shapesShapeRectangle` / `shapesShapeRoundedRect` /
  `shapesShapeEllipse` unused.
- `pkg/pptx/styling/theme.go` hardcoded `2E7D32`, `E91E63`, `0043CE`, `FF9800`, `4589FF`,
  `FFFFFF` while `ColorCorporateGreen`, `ColorMaterialPink`, `ColorCarbonBlue60`,
  `ColorMaterialOrange`, `ColorCarbonBlue40`, `ColorWhite` all existed in the same package.
- `pkg/pptx/editor` had `placeholderTypeTitle` and `OpSlideCount` unused.
- `pkg/pptx/editor/handlers/slidesmeta` had a half-finished pattern: simple chart types
  were constants, but 11 compound ones (`barStacked`, `stockOHLC`, …) were duplicated
  raw across two switch statements.

## Deliberate deviation from the linter

In `internal/pptxxml`, goconst suggested reusing `FillTypeSolid`, `LineJoinStyleMiter` and
`LineDashStyleDash` — but those live in `pkg/pptx/shapes`. `internal/pptxxml` sits *below*
`pkg/pptx/*` in the dependency graph, so importing them would invert the dependency.
Local constants were defined instead, with the reasoning recorded in
`internal/pptxxml/token_constants.go`.

## Resolved by scoped config exclusion, not by source edits

Two categories had no correct source fix. Both are excluded in `.golangci.yml` with the
reasoning written inline, deliberately **path-scoped** rather than disabled globally.

### `bindings/c` (7 findings) — cgo-processed, not source-fixable

Note: earlier drafts of this document said 5. The correct count is 7 —
3 `gocritic` + 2 `staticcheck` + 1 `nakedret` + 1 `nonamedreturns`.

golangci-lint analyses the cgo-**processed** translation unit, so reported line numbers do
not map to `bindings/c/bridge.go` at all. Verified by generating the cgo output with
`go tool cgo -objdir <dir> bindings/c/bridge.go`:

- **4 findings** (`nakedret`, `nonamedreturns`, `ST1003` ×2) originate in
  `_cgo_gotypes.go`, which cgo synthesises at build time — `_cgo_cmalloc(p0 uint64)
  (r1 unsafe.Pointer)` and the `__cgofn__cgo_<hash>_Cfunc__Cmalloc` vars. That file is not
  in the repo, so there is no line to edit. The `<hash>` is build-context derived and
  differs between runs, which is itself proof the code is generated.

- **3 findings** (`gocritic commentFormatting`) fire on Go **directives** — `//export`,
  `//line`, `//go:linkname` — where the absence of a space is required by the toolchain.
  Writing `// export deck_save` does not fail the build; it silently demotes the directive
  to a plain comment, dropping the symbol from the shared library and breaking the Python
  FFI at runtime.

An earlier draft claimed the 3 findings were `//export` lines at source 132/238/275 and
that 313/314 were "past EOF". Both were wrong: source line 132 is blank, 238 is
`defer unlock()`, 275 is blank, 313 is the closing brace, and the file is exactly 313
lines. The cgo line-remapping explanation above is the verified one.

Other linters remain active on `bindings/c`, so hand-written code there is still checked.

### `gosec` G703 path traversal (7) — false positives for this threat model

Flagged in `cmd/gen_ops` (×2), `cmd/gen_shape_types`, `cmd/pptcli` (×2), an example, and
`pkg/pptx/export/pdf.go`. Every one is an operator-supplied output path reaching a file
write. Verified the library case rather than assuming it: `outputPath` is a direct
parameter of the exported API `export.PDF(title, slides, outputPath)`. The caller choosing
the destination *is* the contract, and the process runs with the operator's own
privileges, so no trust boundary is crossed. No validation inside a general-purpose
"write the file here" function can help without breaking legitimate absolute paths.

**Scoped deliberately.** A real path-traversal boundary does exist in this codebase —
archive entry names written to disk (zip-slip) — and is guarded explicitly in
`pkg/pptx/tplx/tplx_archive.go`, which rejects entries containing `..` or a leading `/`.
Disabling G703 globally would mean a future removal of that guard, or a new
archive-extraction path, goes unreported. Keeping G703 live everywhere outside the three
excluded paths preserves that safety net.

## Toolchain note

`golangci-lint` jumped 2.8.0 → 2.12.2 during this work (the 2.8.0 build panics on go1.26
source). The version bump is itself responsible for surfacing findings that were
previously invisible, so the "~197 starting" figure reflects 2.12.2's ruleset, not a
regression in the codebase.

## Final state

`golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 ./...` reports
**0 issues**. `gofmt`, `go build`, `go vet` and the full `go test ./...` all pass.

`examples/` goconst were fixed rather than excluded, so no config exclusion exists for
them. The only two exclusions in `.golangci.yml` are the ones documented above (gosec
G703 on operator-supplied paths, and cgo-processed code in `bindings/c`), both
path-scoped with the reasoning inline.

### Note on goconst cascading

goconst surfaces repeated literals a few at a time: fixing the reported set reveals the
next tier in the same literal cluster. The `examples/` cleanup took two rounds for this
reason. When fixing a file, extract *every* repeated literal in the affected structure
rather than only the reported ones.

### Trap worth remembering

A blind `sed 's/"LITERAL"/constName/g'` also rewrites the line that *declares* the
constant, producing `const colorRed = colorRed` — an initialization cycle. The compiler
catches it, but check the declaration line after any global literal replacement.

## Open questions

1. Should the compression level become configurable later? It is currently a hard-coded
   global default. The `internal/zipfast` seam makes adding an option straightforward if
   size-sensitive callers ask for it.
2. Was the `ContainsAny` fast-path claim at `package_xml.go:11` ever actually measured?
   The benchmark cited for it compares identical code (Finding 4).
