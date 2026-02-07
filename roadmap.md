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

## Active parity backlog
1. Next parity target (to choose):
- Deepen table visual parity (advanced border semantics, mixed-edge styling).
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
