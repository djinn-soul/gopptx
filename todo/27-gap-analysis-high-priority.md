# 27 - Gap Analysis (High Priority)

Scope: prioritized gaps for `gopptx` based on the provided section `6.1 High Priority Gaps (High Impact, Feasible)`.

## Prioritized Gaps

| Gap | Impact | Notes | Tracking |
| --- | --- | --- | --- |
| Strikethrough text | High | Common formatting; missing from run-level props | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| Superscript / subscript | High | Scientific and technical content | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| Font highlight color | High | Common in annotation workflows | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| Vertical text alignment (anchor) | High | Required for proper table/shape layout | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| Bullet list support | High | Standard presentation feature | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| Line spacing control | High | Essential for layout density control | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| Character spacing | Medium | Typography control | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| Paragraph space before/after | Medium | Layout control | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| Auto-fit text (shrink/resize) | Medium | Overflow handling | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| Text wrap control | Medium | Layout control | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| Image crop | High | Very common image operation | [18-missing-images-media-features.md](18-missing-images-media-features.md) |
| Image rotation | High | Common visual treatment | [18-missing-images-media-features.md](18-missing-images-media-features.md) |
| Image opacity | Medium | Common visual treatment | [18-missing-images-media-features.md](18-missing-images-media-features.md) |
| Image shadow / glow | Medium | Polish effects | [18-missing-images-media-features.md](18-missing-images-media-features.md) |
| SVG image support | Medium | Icon and graphic libraries | [18-missing-images-media-features.md](18-missing-images-media-features.md) |
| Horizontal bar chart | High | Common chart variant | [16-missing-chart-features.md](16-missing-chart-features.md) |
| Stacked bar chart | High | Very common in business decks | [16-missing-chart-features.md](16-missing-chart-features.md) |
| Scatter / XY chart | High | Data visualization | [16-missing-chart-features.md](16-missing-chart-features.md) |
| Doughnut chart | Medium | Common KPI display | [16-missing-chart-features.md](16-missing-chart-features.md) |
| Area chart | Medium | Time-series visualization | [16-missing-chart-features.md](16-missing-chart-features.md) |
| Table auto-pagination | Medium | Long data tables | [17-missing-table-features.md](17-missing-table-features.md) |
| Table repeat header | Medium | Multi-page tables | [17-missing-table-features.md](17-missing-table-features.md) |
| Bullets in table cells | Medium | Rich table content | [17-missing-table-features.md](17-missing-table-features.md) |
| Shape shadow | Medium | Visual polish | [23-missing-shapes-advanced-operations-features.md](23-missing-shapes-advanced-operations-features.md) |
| Shape flip (H/V) | Low | Mirror transformations | [23-missing-shapes-advanced-operations-features.md](23-missing-shapes-advanced-operations-features.md) |

## Execution Order Suggestion

- [ ] Phase 1 (high impact): text core + chart core + image core.
- [ ] Phase 2 (layout polish): table pagination/header + shape shadow/flip.
- [ ] Phase 3 (quality): tests and docs for all delivered gaps.
