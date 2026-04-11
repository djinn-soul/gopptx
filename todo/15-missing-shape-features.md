# 15 - Missing Shape Features

Scope: track missing shape capabilities in `gopptx` from the provided feature audit screenshot.

## Missing in gopptx

| Feature | Status | TODO |
| --- | --- | --- |
| Position (percentages) | Missing | [ ] Add percentage-based shape positioning API and bridge payload mapping (relative-to-slide coordinates). |
| Glow effects | Missing | [ ] Add shape glow effects API and bridge support with documented options. |
| Scheme / theme-aware colors | Missing | [ ] Add theme/scheme color tokens for shape fill/line APIs and bridge mapping. |

## Capability Depth Gap

| Feature | Status | TODO |
| --- | --- | --- |
| Preset shape types coverage | Limited | [ ] Audit preset shape coverage and expand enum/token support where feasible. |

## Verification Follow-up

- [ ] Verify listed gaps against current runtime behavior and update this file if needed.
- [ ] Add regression tests for percentage positioning, shape glow, and theme-aware color flows once implemented.
