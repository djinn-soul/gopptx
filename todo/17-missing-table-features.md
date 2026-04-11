# 17 - Missing Table Features

Scope: track missing table capabilities in `gopptx` from the provided feature audit screenshots.

## Missing in gopptx

| Feature | Status | TODO |
| --- | --- | --- |
| Bullets in cells | Missing | [ ] Add bullet/list rendering support inside table cells, including numbered/bulleted runs. |
| Auto-paging on overflow | Missing | [ ] Add table overflow pagination API and bridge behavior for multi-slide continuation. |
| Repeat header row | Missing | [ ] Add table header-row repeat controls for multi-slide or paginated table rendering. |
| HTML table -> PPTX conversion | Missing | [ ] Add HTML-table ingestion and conversion helpers for PPTX table generation. |
| TableMergeMap (advanced) | Missing | [ ] Add advanced merge-map table construction API for complex merged-cell layouts. |

## Capability Depth Gap

| Feature | Status | TODO |
| --- | --- | --- |
| BorderStyle variations | Partial | [ ] Expand border-style variant coverage and document supported style matrix for table cell borders. |

## Verification Follow-up

- [ ] Verify listed gaps against current runtime behavior and update this file if needed.
- [ ] Add regression tests for bullets in cells, repeat header row, overflow pagination, HTML-table conversion, advanced merge-map behavior, and border-style variants once implemented.
