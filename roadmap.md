# goppt Roadmap

## Current parity status
- Chart variant parity with `ppt-rs` is implemented for:
- `BarHorizontal`, `BarStacked`, `BarStacked100`
- `LineMarkers`, `LineStacked`
- `AreaStacked`, `AreaStacked100`
- `Bubble`
- `Radar`, `RadarFilled`
- `StockHLC`, `StockOHLC`
- `Combo`
- Deterministic chart parity signature report now covers 21 chart variants (`reports/chart_parity_report.md`, `21/21` pass).
- Table parity includes styled cell support:
- Bold cell text
- Cell background fill color (RGB hex)
- Cell horizontal alignment (`l|ctr|r|just`)
- Cell vertical alignment (`t|ctr|b`)
- Deterministic uniform cell borders (width + RGB hex color)
- Per-side cell border controls (`left|right|top|bottom`) with independent width/color
- Per-side border dash styles (`solid|dash|dot|dashDot|lgDash`)

## Active parity backlog
1. Add fixture parity checks for styled bullet and enhanced-text decks from `ppt-rs` examples.
2. Expand slide surface parity beyond layout selection (placeholder-level controls).

## Completed recently
- Expanded chart parity report coverage beyond the prior 6-signature slice:
- Added full variant catalog comparisons for `bar`, `line`, `area`, `pie`, `doughnut`, `scatter` (3 styles), `bubble`, `radar`, `stock`, and `combo`.
- Updated `tools/ppt-rs-chart-signatures` and `scripts/compare_chart_parity_with_ppt_rs` to emit/compare 21 deterministic chart signatures.
- Added explicit normalization for known `ppt-rs` bar token quirks so semantic checks remain strict and stable.

1. `StockHLC` / `StockOHLC` visual parity enhancement:
- Added `<c:hiLowLines>` for stock charts.
- Added `<c:upDownBars>` for OHLC charts.
- Manual PowerPoint screenshots confirm improved financial-style rendering.

2. Deterministic appearance regression coverage:
- Added Bubble/Radar/Stock/Combo XML signature tests in `pkg/pptx/chart_parity_signature_test.go`.

3. Table styled cell parity slice:
- Added `TableCell` style model in `pkg/pptx/table.go`.
- Added styled row API `AddStyledRow([]TableCell)`.
- Added table XML rendering for bold text and background fill in `internal/pptxxml/slide_table_xml.go`.
- Added validation and integration tests in `pkg/pptx/presentation_test.go`.

4. Table alignment/border parity slice:
- Added cell alignment APIs (`WithAlign*`, `WithVAlign*`) and validation in `pkg/pptx/table.go`.
- Added border API (`WithBorder(widthPt, color)`) with strict validation in `pkg/pptx/table.go`.
- Added OOXML rendering for `<a:pPr algn="...">`, `<a:tcPr anchor="...">`, and deterministic `lnL/lnR/lnT/lnB` borders in `internal/pptxxml/slide_table_xml.go`.
- Added integration coverage in `pkg/pptx/table_alignment_border_test.go`.

5. Slide layout parity slice:
- Added caller-selectable slide layouts in `pkg/pptx`:
- `SlideLayoutTitleAndContent` (default)
- `SlideLayoutTitleOnly`
- `SlideLayoutBlank`
- Added fluent APIs: `WithLayout(...)`, `WithTitleAndContentLayout()`, `WithTitleOnlyLayout()`, `WithBlankLayout()`.
- Added strict layout validation:
- `title_only` rejects bullets.
- `blank` requires empty title and rejects bullets.
- Added OOXML packaging/rendering updates:
- Multiple layout parts (`slideLayout1/2/3.xml`) and relationships.
- Layout-specific slide rel target selection.
- Layout-aware slide shape rendering.
- Added integration coverage in `pkg/pptx/slide_layout_test.go`.

6. Deep table border semantics parity slice:
- Added per-side border APIs in `pkg/pptx/table.go`:
- `WithLeftBorder*`, `WithRightBorder*`, `WithTopBorder*`, `WithBottomBorder*`
- Added dashed border style APIs (`With*BorderStyle`) with allowed dash values:
- `solid`, `dash`, `dot`, `dashDot`, `lgDash`
- Preserved backward compatibility of `WithBorder(widthPt, color)` by mapping to explicit side borders.
- Extended table XML model/rendering in `internal/pptxxml/slide_table_xml.go` to emit only configured border sides with side-specific dash styles.
- Added strict validation for per-side width/color/dash semantics in `pkg/pptx/table_validation.go`.
- Added integration coverage for side overrides and dash rendering in `pkg/pptx/table_alignment_border_test.go`.

7. Markdown inline rich-text parity slice:
- Added markdown inline parsing for `**bold**`, `*italic*`, and `` `code` `` in `pkg/pptx/markdown.go`.
- Added additive rich text model/API for bullets in `pkg/pptx/text_runs.go` and `pkg/pptx/slide.go` (`TextRun`, `AddBulletRuns`).
- Added run-aware slide XML rendering in `internal/pptxxml/slide_xml.go` with code font mapping (`Consolas`).
- Wired run mapping in `pkg/pptx/presentation.go`.
- Added parser + XML integration tests in `pkg/pptx/markdown_test.go`.

8. Text formatting parity slice (run-level):
- Extended `TextRun` API in `pkg/pptx/text_runs.go` with `underline`, `color`, `font`, and `size` controls.
- Added run validation in `pkg/pptx/text_runs_validation.go` (RGB color + size bounds).
- Refined OOXML run rendering in `internal/pptxxml/slide_text_runs_xml.go` with deterministic run-property ordering and style emission.
- Added integration tests in `pkg/pptx/text_runs_test.go` and updated markdown OOXML assertions in `pkg/pptx/markdown_test.go`.

9. Text formatting parity slice (paragraph-level):
- Added paragraph style API in `pkg/pptx/text_paragraph.go`:
- Alignment (`l|ctr|r|just`)
- Spacing controls (`space-before`, `space-after`, `line-spacing`)
- Added paragraph style validation in `pkg/pptx/text_paragraph_validation.go`.
- Extended slide XML rendering with paragraph style tags in `internal/pptxxml/slide_text_paragraph_xml.go`.
- Wired paragraph style mapping in `pkg/pptx/presentation.go`.
- Added integration coverage in `pkg/pptx/text_paragraph_test.go`.

10. Bullet styles parity slice:
- Added bullet style model in `pkg/pptx/text_bullet_style.go`:
- `bullet`, `number`, `letter_lower`, `letter_upper`, `roman_lower`, `roman_upper`, `custom`, `none`
- Added helper bullet APIs in `pkg/pptx`:
- `WithBulletStyle`, `AddStyledBullet`, `AddNumbered`, `AddLettered`, `AddRomanLower`, `AddRomanUpper`, `AddCustomBullet`, `AddSubBullet`
- Extended paragraph style model and validation in `pkg/pptx/text_paragraph.go` + `pkg/pptx/text_paragraph_validation.go` for bullet style, custom char, and level (`0..8`).
- Extended OOXML paragraph rendering in `internal/pptxxml/slide_text_paragraph_xml.go` for `buAutoNum`, `buChar`, `buNone`, and level-based indentation.
- Added integration coverage in `pkg/pptx/bullet_styles_test.go`.
- Updated markdown ordered-list parsing in `pkg/pptx/markdown.go` to emit numbered bullet style.

11. Text enhancements parity slice:
- Extended `TextRun` API in `pkg/pptx/text_runs.go`:
- `WithStrikethrough`, `WithHighlight`, `WithSubscript`, `WithSuperscript`
- Added ppt-rs-comparable run-size presets in `pkg/pptx/text_size_presets.go`.
- Added strict validation in `pkg/pptx/text_runs_validation.go`:
- highlight color must be RGB hex
- subscript + superscript cannot both be true
- Extended OOXML run rendering in `internal/pptxxml/slide_text_runs_xml.go`:
- `strike="sngStrike"`
- `baseline="-25000"` / `baseline="30000"`
- `<a:highlight><a:srgbClr .../></a:highlight>`
- Added integration coverage in `pkg/pptx/text_runs_test.go`.

12. Foundational drawings slice (shapes + connectors):
- Added public shape model/API in `pkg/pptx/shape.go`:
- baseline shape preset enum/constants
- fill, line, transparency, rotation, and text setters
- Added public connector model/API in `pkg/pptx/connector.go`:
- straight/elbow/curved connector types
- dash + arrowhead + arrow-size controls
- shape-anchor connection API (`ConnectStart`/`ConnectEnd`)
- Added strict drawing validation in `pkg/pptx/slide_drawings_validation.go` (types, geometry, colors, dash/arrow values, anchor range/site checks).
- Wired drawing specs in `pkg/pptx/presentation_xml_specs.go` and slide assembly in `pkg/pptx/presentation.go`.
- Added custom shape/connector OOXML rendering in `internal/pptxxml/slide_drawings_xml.go` and integrated it into `internal/pptxxml/slide_xml.go`.
- Added EMU helper building blocks in `pkg/pptx/units.go` (`Inches`, `Centimeters`, `Points`).
- Added integration tests in `pkg/pptx/shape_connector_test.go`.

13. Markdown parity expansion slice:
- Refactored markdown parsing into modular components:
- parser state machine in `pkg/pptx/markdown_parser.go`
- block parsers in `pkg/pptx/markdown_blocks.go`
- code/mermaid handlers in `pkg/pptx/markdown_code_mermaid.go`
- Added GFM table parsing into native `Table` model with styled header rows.
- Added fenced code-block parsing with language-tagged, no-bullet code paragraphs and deterministic line coloring.
- Added Mermaid fenced-block parsing for supported diagram directives with typed placeholder shapes.
- Added blockquote parsing into slide speaker notes (`SlideContent.Notes`).
- Added end-to-end fixture coverage using upstream `md2ppt_demo.md` in `pkg/pptx/testdata/ppt_rs/md2ppt_demo.md` and `pkg/pptx/markdown_blocks_test.go`.

14. Speaker-notes persistence slice:
- Added notes-slide XML generation in `internal/pptxxml/notes_xml.go`.
- Added notes package wiring in `pkg/pptx/presentation.go`:
- notes-slide content type overrides
- notes-master relationship and package parts
- per-slide `notesSlide` relationship attachment when notes exist
- Added concurrent notes rendering pipeline in `pkg/pptx/notes_parts.go` (goroutine fan-out, deterministic compaction back to slide order).
- Added notes parts writer in `pkg/pptx/notes_package.go`.
- Added integration tests in `pkg/pptx/presentation_notes_test.go` for both direct notes API and markdown blockquote persistence.

15. CLI command parity slice:
- Refactored `cmd/pptcli` into modular command handlers with explicit exit codes.
- Added subcommands: `create`, `md2ppt`, `info`, `validate`.
- Preserved legacy flag mode (`-out`, `-md`, `-title`) for backward compatibility.
- Added deterministic CLI integration tests in `cmd/pptcli/main_test.go`.
- Added command documentation and usage examples in `cmd/pptcli/README.md`.

16. Read/modify existing PPTX API slice:
- Added `PresentationEditor` API in `pkg/pptx` for opening existing PPTX packages with `archive/zip`.
- Added parsed metadata surfaces (`PresentationMetadata`, `SlideMetadata`) with concurrent slide-title extraction.
- Added editing operations:
- `AddSlide`
- `UpdateSlide` (fail-fast on unsupported external slide relationships)
- `RemoveSlide`
- `MergeFromFile` / `MergeFromEditor`
- Added round-trip save pipeline preserving untouched package parts while rewriting slide list/relationships/content-types entries.
- Added fixture-backed regression tests in `pkg/pptx/editor_test.go` for open/edit/save, merge, preservation, and corruption handling.

17. Advanced table controls + parity reporting slice:
- Added explicit row-height controls via `Table.WithRowHeights([]int64)` with strict validation.
- Added per-cell text layout controls in `TableCell`:
- `WithMarginsPt`
- `WithMarginLeftPt` / `WithMarginRightPt` / `WithMarginTopPt` / `WithMarginBottomPt`
- `WithWrap(bool)`
- Extended OOXML table rendering to emit:
- per-row `<a:tr h="...">` from explicit row heights
- per-cell `marL/marR/marT/marB` on `<a:tcPr>`
- per-cell wrap mode on `<a:bodyPr wrap="...">`
- Added integration validation coverage in `pkg/pptx/table_alignment_border_test.go` for row heights, margins, wrap output, and invalid margin/row-height cases.
- Added deterministic table parity report pipeline in `scripts/compare_table_parity_with_ppt_rs/main.go` with generated output at `reports/table_parity_report.md`.

18. Shape gradient fill parity slice:
- Added gradient fill model/API in `pkg/pptx`:
- `ShapeGradientFill` and `ShapeGradientStop`
- gradient types: `linear`, `radial`, `rectangular`, `path`
- fluent API: `WithGradientFill(...)` and `WithLinearAngle(...)`
- Added strict shape validation in `pkg/pptx/slide_drawings_validation.go`:
- fail-fast for unsupported gradient types
- requires at least two strictly-increasing stops in `[0..100]`
- validates stop colors/transparency and linear-angle constraints
- Added OOXML rendering in `internal/pptxxml`:
- emits `<a:gradFill>` + `<a:gsLst>` stop nodes
- emits linear vector (`<a:lin>`) and non-linear path modes (`circle|rect|shape`)
- Added integration coverage in `pkg/pptx/shape_gradient_test.go`.

19. Connector auto-site parity slice:
- Added connector auto-anchor APIs in `pkg/pptx`:
- `ConnectStartAuto(shapeIndex)` and `ConnectEndAuto(shapeIndex)`
- Added relative-position site inference for anchored connectors:
- horizontal dominance -> `left|right`
- vertical dominance -> `top|bottom`
- Wired auto-site resolution into XML spec conversion so empty anchored sites emit deterministic `stCxn/endCxn` indices.
- Updated connector validation to allow anchored connectors with empty sites (auto mode) while preserving invalid-site fail-fast behavior.
- Added flow-layout integration tests in `pkg/pptx/connector_auto_sites_test.go` for horizontal and vertical anchor inference.

20. Shape text-contrast parity slice:
- Added automatic text-color contrast selection for shape text based on fill luminance and WCAG-style contrast ratio.
- Supports both solid fills (including transparency blending) and gradient fills (average stop color with per-stop transparency).
- Wired contrast-aware text run rendering into shape XML output (`a:rPr` + `a:solidFill`).
- Added integration tests in `pkg/pptx/shape_text_contrast_test.go` for dark/light solid fills and dark-gradient fills.

21. Shape text auto-fit parity slice:
- Added automatic shape text sizing based on shape bounds and text length heuristics.
- Wired shape text bodies to emit `<a:spAutoFit/>` and dynamic run-size values in OOXML.
- Added integration coverage in `pkg/pptx/shape_text_autofit_test.go` to assert long text renders with smaller font size than short text.

22. Shape XML parallel rendering slice:
- Added deterministic concurrent shape rendering helper in `internal/pptxxml` for large slides.
- Preserved stable shape ID ordering while computing each `p:sp` block in parallel goroutines.
- Wired `SlideWithLayout` to use the parallel helper before connector emission.
- Added focused deterministic-order coverage in `internal/pptxxml/slide_drawings_parallel_xml_test.go`.

23. Extensible shape interface parity slice:
- Added `ShapeDefinition` interface in `pkg/pptx` to support pluggable shape builders.
- Implemented `Shape.ToShape()` for backward compatibility with existing fluent usage.
- Updated `SlideContent.AddShape(...)` to accept interface-backed shape definitions.
- Added integration coverage in `pkg/pptx/shape_definition_test.go` using a custom external shape type.

24. Connector site-selection optimization slice:
- Replaced dominant-axis heuristic with geometry-aware nearest-anchor selection across edge/corner/center sites.
- Optimized anchor resolution to use shape-local candidate points and squared-distance comparisons (no floating-point cost).
- Preserved deterministic behavior and backward compatibility with explicit site overrides.
- Added diagonal-flow parity coverage in `pkg/pptx/connector_auto_sites_test.go` to validate corner-site inference.
