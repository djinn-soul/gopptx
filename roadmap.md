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
- Table parity includes styled cell support:
- Bold cell text
- Cell background fill color (RGB hex)
- Cell horizontal alignment (`l|ctr|r|just`)
- Cell vertical alignment (`t|ctr|b`)
- Deterministic uniform cell borders (width + RGB hex color)
- Per-side cell border controls (`left|right|top|bottom`) with independent width/color
- Per-side border dash styles (`solid|dash|dot|dashDot|lgDash`)

## Active parity backlog
1. Next parity target (to choose):
- Expand slide surface parity beyond layout selection (placeholder-level controls).

## Completed recently
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
