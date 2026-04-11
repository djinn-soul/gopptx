# 20 - Missing Advanced and Unique Features

Scope: track missing advanced capabilities in `gopptx` from the provided feature audit screenshots.

## Missing in gopptx

| Feature | Status | Tracking |
| --- | --- | --- |
| HTML table -> PPTX | Missing | Tracked in [17-missing-table-features.md](17-missing-table-features.md). |
| Web URL -> PPTX | Missing | [ ] Add URL-to-PPTX import pipeline (fetch, parse, map to slides) with robust validation and limits. |
| HTML -> PPTX | Missing | [ ] Add HTML-to-PPTX import/conversion support with layout/text/image mapping controls. |
| PDF -> PPTX (import) | Missing | [ ] Add PDF-to-PPTX import support with page-to-slide mapping and conversion options. |
| Percentage-based positioning | Missing | Tracked in [15-missing-shape-features.md](15-missing-shape-features.md). |

## Verification Follow-up

- [ ] Verify listed gaps against current runtime behavior and keep cross-file tracking aligned.
- [ ] Add regression tests for URL/HTML/PDF import paths and related conversion options once implemented.
