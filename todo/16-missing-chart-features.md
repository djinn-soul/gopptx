# 16 - Missing Chart Features

Scope: track missing chart capabilities in `gopptx` from the provided feature audit screenshots.

## Missing in gopptx

| Feature | Status | TODO |
| --- | --- | --- |
| Bar horizontal | Missing | [ ] Add horizontal bar chart type support and mapping in chart APIs. |
| Bar stacked | Missing | [ ] Add stacked bar chart type support and serialization. |
| Bar stacked 100% | Missing | [ ] Add 100% stacked bar chart support and normalization behavior. |
| 3D Bar / 3D Bubble | Missing | [ ] Add 3D chart API surface and bridge mapping for supported 3D chart families. |
| Line with markers | Missing | [ ] Add line-with-markers chart type support. |
| Line stacked | Missing | [ ] Add stacked line chart type support. |
| Line smooth | Missing | [ ] Add smooth-line chart variant support for line series. |
| Area chart | Missing | [ ] Add area chart type support. |
| Area stacked / 100% | Missing | [ ] Add stacked and 100% stacked area chart support. |
| Doughnut chart | Missing | [ ] Add doughnut chart type support and key configuration options. |
| Scatter / XY chart | Missing | [ ] Add scatter/XY chart type support. |
| Scatter with lines | Missing | [ ] Add scatter-with-lines chart variant support. |
| Scatter smooth | Missing | [ ] Add smooth-scatter chart variant support. |
| Bubble chart | Missing | [ ] Add bubble chart type support with size mapping. |
| Bubble 3D | Missing | [ ] Add 3D bubble chart variant support. |
| Radar chart | Missing | [ ] Add radar chart type support. |
| Radar filled | Missing | [ ] Add filled-radar chart variant support. |
| Stock HLC | Missing | [ ] Add stock HLC chart type support with correct series mapping. |
| Stock OHLC | Missing | [ ] Add stock OHLC chart type support with correct series mapping. |
| Combo chart | Missing | [ ] Add combo/mixed-series chart type support and axis binding controls. |
| Logarithmic axis scaling | Missing | [ ] Add logarithmic axis controls (base/min/max behavior) in chart axis APIs and bridge payloads. |
| Data table in chart | Missing | [ ] Add chart data-table display controls (show/hide, formatting, legend keys) in chart APIs. |
| Leader lines | Missing | [ ] Add chart data-label leader-line controls for supported chart types. |
| Chart area formatting | Missing | [ ] Add chart area formatting controls (fill/line/effects) in chart APIs. |
| Plot area formatting | Missing | [ ] Add plot area formatting controls (fill/line/effects) in chart APIs. |
| First slice angle (pie) | Missing | [ ] Add pie first-slice-angle controls and serialization. |
| Bar gap width | Missing | [ ] Add bar-gap-width controls for bar/column chart spacing. |
| Display units on axis | Missing | [ ] Add axis display-units controls for value-axis formatting. |

## Verification Follow-up

- [ ] Verify listed gaps against current runtime behavior and update this file if needed.
- [ ] Add regression tests for all newly added chart-type rows plus logarithmic-axis behavior once implemented.

## Coverage Snapshot (Screenshot-Derived)

| Metric | gopptx | Notes |
| --- | --- | --- |
| Total chart types currently supported | 3 | Screenshot-derived summary row; verify against runtime and update this value if stale. |
