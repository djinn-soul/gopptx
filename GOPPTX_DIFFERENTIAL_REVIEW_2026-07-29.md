# gopptx Differential Review — 2026-07-29

## Executive Summary

| Severity | Count |
|----------|-------|
| Critical | 0 |
| High | 1 |
| Medium | 5 |
| Low | 0 |

**Overall risk:** High
**Recommendation:** Reject until the high- and medium-severity findings are fixed.

The review covered the complete working-tree product diff against `75ce2ff`,
excluding `.serena/`, `CONTINUITY.md`, and `scratch/`. That scope contains 40
modified tracked files and 20 relevant untracked Go/Python source or test files.
No security vulnerability was found; the blocking issues are correctness,
public-API behavior, validation, and regression coverage.

## Findings

### HIGH: A shape is not associated with its own chart

**Files:** `python/gopptx/slide/shapes/shape_proxy.py:117`, `:136-144`
**Test coverage:** Missing
**Blast radius:** Every `ShapeProxy` on a slide containing a chart

`has_chart` falls back to whether the slide contains any chart, so a textbox on
a chart-bearing slide reports `True`. The `chart` property then returns the first
slide chart for every shape. It also reads `rel_id` and `part`, while the bridge
contract returns `RelID` and `ChartPart`, and fabricates `rId1/chart1.xml` when no
chart exists.

Reproduction from the live tree:

```text
plain_shape_type= rect
plain_shape_has_chart= True
shape_chart= 3 0 rId1 ppt/charts/chart1.xml
shape_chart= 4 0 rId1 ppt/charts/chart1.xml
```

The slide actually returned two chart references (`rId2/chart1.xml` and
`rId3/chart2.xml`). Mutating `second_chart_shape.chart` therefore targets chart
index 0 and can silently modify the wrong chart.

**Recommendation:** Include the containing graphic-frame shape ID in
`SlideChartRef`, map by `ShapeProxy.id`, consume the actual `RelID` and
`ChartPart` keys, and raise `AttributeError` for non-chart shapes. Add a
two-chart-plus-textbox regression.

### MEDIUM: `add_paragraph()` does not create a paragraph

**Files:** `python/gopptx/slide/text/text_paragraph_model.py:151-176`,
`python/gopptx/slide/text/text_model.py:32-34`
**Test coverage:** Missing

The collection remains hard-coded to length one. `add_paragraph()` appends a run
whose text starts with `"\n"` and returns the existing paragraph proxy. It never
creates a second `<a:p>` and cannot preserve independent paragraph properties.

Live reproduction:

```text
paragraphs_before= 1
paragraphs_after= 1
returned_is_first= True
```

**Recommendation:** Extend the bridge model to read and write a list of
paragraphs, create a new `<a:p>` with its own runs/properties, and return a proxy
bound to the new paragraph index.

### MEDIUM: Axis-title visibility writes unsupported payload keys

**File:** `python/gopptx/slide/chart/axis_series.py:50-116`
**Test coverage:** Missing

The new setters send `category_axis_has_title` and
`value_axis_has_title`. Neither key exists in `ChartFormatUpdate`; the Go engine
only accepts `category_axis_title` and `value_axis_title`. Unknown JSON fields
are ignored, so `axis.has_title = False` silently does nothing.

Live reproduction after setting a value-axis title and then disabling it:

```text
axis_title_remains= True
```

This mismatch is also reported by basedpyright as incompatible
`ChartFormatUpdate` payloads.

**Recommendation:** Add an explicit title visibility/removal field through the
Go and Python contracts, or map `False` to a supported removal operation. Test
show, hide, and reopen behavior.

### MEDIUM: `_FontProxy.remove()` is a broken copy of run removal

**File:** `python/gopptx/slide/text/text_run_model.py:153-159`
**Test coverage:** Missing

`_FontProxy` owns only `_run_proxy`, but `remove()` accesses `_text_frame` and
`_run_index`. Calling it raises `AttributeError`. basedpyright identifies all of
these members as missing.

**Recommendation:** Remove this method from `_FontProxy`; run removal already
exists on `_ShapeRunProxy` and `_ShapeRunCollection`.

### MEDIUM: Shape alt text and title cannot be cleared

**File:** `pkg/pptx/editor/shape_editor_mutation.go:297-335`
**Test coverage:** Missing for clearing

The updater applies `descr` and `title` only when the new strings are non-empty.
Consequently, assigning `shape.alt_text = ""` or `shape.title = ""` reports
success while retaining the old OOXML attributes.

**Recommendation:** Treat a non-nil empty string as an explicit clear and remove
the corresponding attribute. Keep nil as “leave unchanged.”

### MEDIUM: The mandatory validation gates are not green

**Files:** Multiple changed Python files; principal failures include
`python/gopptx/slide/shapes/shape_proxy.py`,
`python/gopptx/slide/chart/axis_series.py`,
`python/gopptx/slide/shapes/picture_image.py`, and
`python/gopptx/slide/text/text_run_model.py`.

- Ruff: 35 errors.
- basedpyright: 51 errors.
- Architectural guardrails: new over-ceiling
  `shape_proxy.py` at 463 lines versus the configured 400-line ceiling.
- The workspace instruction is stricter (300 lines); eleven changed code files
  exceed 300 lines.

Examples include undefined `os` in a public annotation, incompatible axis-title
property overrides, unknown protocol members, invalid `_FontProxy` members,
private-access violations, imports inside functions, and excessive class size.

**Recommendation:** Split the added shape/chart/text facades into focused
modules, repair the protocols and TypedDict payloads, and require Ruff,
basedpyright, and the architecture guardrail to pass before merge.

## Test Coverage Analysis

`uv run pytest python/tests` reports 459 passing tests, and
`go test ./pkg/pptx/editor/...` passes. Those results do not exercise the
failure cases above.

The changed `python/tests/test_table_cell_borders.py` also replaces three
round-trip/schema regressions with one broad XML-presence assertion, removing
coverage for cap/join/inset persistence, border-before-fill schema order, and
border clearing. The new SVG and picture-background tests similarly accept
generic `blip` substrings rather than proving the feature-specific contract.

Required regressions:

- Textbox plus two chart shapes: exact `has_chart` and chart association.
- Actual paragraph count and independent formatting after reopen.
- Axis-title show/hide/remove after reopen.
- Clearing shape alt text/title.
- `run.remove()` remains available while `run.font` has no removal method.
- Restore the removed table-border round-trip and ordering assertions.

## Historical Context

Baseline `75ce2ff` did not contain the reviewed shape-chart, paragraph-add,
axis-title proxy, font proxy, layout-reorder, or metadata-update APIs. No
security-related validation or authorization code was removed.

## Validation

| Command | Result |
|---------|--------|
| `git diff --check` | Passed |
| `go test ./pkg/pptx/editor/...` | Passed |
| `uv run pytest python/tests` | 459 passed |
| `uv run ruff check python/gopptx python/tests scripts/ci hatch_build.py` | Failed: 35 errors |
| `uv run basedpyright` | Failed: 51 errors |
| `uv run python scripts/ci/architectural_guardrails.py` | Failed: `shape_proxy.py` has 463 lines |

## Methodology and Limits

**Strategy:** Focused differential review of a medium/large working-tree change.

The review inspected all changed file groups, traced the public Python APIs into
the Go bridge contracts, checked relevant history and call sites, ran the
project's focused Go and full Python tests, ran static and architecture gates,
and executed isolated reproductions for the highest-risk API paths.

PowerPoint rendering was not repeated during this review. The chart line-style
OOXML paths have prior render evidence in `CONTINUITY.md`, but the newly reviewed
shape/chart association, paragraph creation, and axis-title visibility paths do
not.

**Confidence:** High for the reported findings; medium for the entire 60-file
scope because this was a focused differential review rather than a line-by-line
audit of every low-risk test and stub change.
