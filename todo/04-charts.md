# 04 - Charts

Source: `TODO.md` plus `web_sid/docs.aspose.com/slides/python-net/{create-chart,powerpoint-charts,chart-types,chart-series,chart-calculations,chart-axis,chart-legend,chart-data-label,chart-data-marker,chart-formatting,chart-plot-area,chart-entities,bubble-chart,pie-chart,doughnut-chart,data-points-of-treemap-and-sunburst-chart,error-bar,trend-line,chart-data-table,chart-workbook,chart-worksheet-formulas}`.

## Supported

| Area | Docs | Status | Notes |
| --- | --- | --- | --- |
| Core chart creation and update | `create-chart/`, `powerpoint-charts/`, `chart-types/`, `chart-series/`, `chart-calculations/` | Supported | `add_chart`, `add_combo_chart`, chart data builders, `update_chart_data`, batch update, replace-by-index/rel-id, and chart listing are marked `[x]`. |
| Axes, legends, labels, markers, formatting, plot area | `chart-axis/`, `chart-legend/`, `chart-data-label/`, `chart-data-marker/`, `chart-formatting/`, `chart-plot-area/`, `chart-entities/` | Supported | The basic proxies and common formatting controls are marked `[x]`. |

## Not Supported / Partial

| Area | Docs | Status | Notes |
| --- | --- | --- | --- |
| Bubble charts | `bubble-chart/` | Partial | Bubble chart creation and x/y/size values are supported, but scale factor, negative bubble visibility, and border color are still gaps. |
| Pie charts | `pie-chart/` | Partial | Base pie + value/category labels are supported, but first-slice angle, explosion, slice fill, and leader lines are gaps. |
| Doughnut charts | `doughnut-chart/` | Partial | Base doughnut type is supported, but hole size, first-slice angle, explosion, and per-segment fill are gaps. |
| Treemap and sunburst | `data-points-of-treemap-and-sunburst-chart/` | Not supported | Treemap/sunburst types and point fill or label-position customization are unchecked. |
| Error bars | `error-bar/` | Not supported | Error bar creation and styling remain unchecked. |
| Trend lines | `trend-line/` | Not supported | Trend line creation, equation display, R-squared, and styling remain unchecked. |
| Chart data tables | `chart-data-table/` | Not supported | Show/hide, legend keys, border, and font controls remain unchecked. |
| Chart workbook and worksheet formulas | `chart-workbook/`, `chart-worksheet-formulas/` | Not supported | Embedded workbook read/write/extract and worksheet formula binding remain unchecked. |
