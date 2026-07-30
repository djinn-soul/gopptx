"""Per-series chart annotation TypedDicts.

Trendlines, error bars, per-point formatting, series-level formatting and the
line styles they share. Split out of ``schemas_chart_layout`` to keep both
modules inside the repository's file-length ceiling.
"""

from __future__ import annotations

try:
    from typing import TypedDict
except ImportError:  # pragma: no cover
    from typing_extensions import TypedDict


class ChartTrendlineSpec(TypedDict, total=False):
    """One c:trendline fitted to a chart series.

    order applies only to the poly type and period only to movingAvg;
    PowerPoint repairs a file that carries either on any other type.
    """

    series_index: int
    type: str
    order: int
    period: int
    name: str
    forward: float
    backward: float
    intercept: float
    display_r_squared: bool
    display_equation: bool
    line_color: str
    line_width_emu: int
    line_dash: str


class ChartErrorBarSpec(TypedDict, total=False):
    """One c:errBars attached to a chart series.

    Custom error bars carry plus_reference/minus_reference formulas; every
    other value_type carries the scalar value instead. Scatter and bubble
    series take one entry per direction, so X and Y bars are separate sets.
    """

    series_index: int
    bar_type: str
    value_type: str
    direction: str
    value: float
    plus_reference: str
    minus_reference: str
    no_end_cap: bool
    line_color: str


class ChartDataPointSpec(TypedDict, total=False):
    """Per-point formatting for one c:dPt in a series.

    explosion applies to pie and doughnut series only; bubble_3d to bubble
    series. The marker_* fields recolour the point's marker on scatter, line
    and radar series, where the point's own fill formats the connecting
    segment rather than the marker.
    """

    series_index: int
    point_index: int
    fill_color: str
    line_color: str
    line_width_emu: int
    invert_if_negative: bool
    bubble_3d: bool
    explosion: int
    marker_fill_color: str
    marker_line_color: str
    marker_symbol: str
    marker_size: int


class ChartLineFormatSpec(TypedDict, total=False):
    """Line style of a chart line element.

    Used for axis gridlines, the series lines of a stacked bar chart, and the
    outline of a data label. ``none`` draws no line at all and wins over the
    other three keys.
    """

    color: str
    width_emu: int
    dash: str
    none: bool


class ChartDataLabelPointSpec(TypedDict, total=False):
    """Formatting for the label of a single data point, its c:dLbl.

    PowerPoint keeps a point's number format and label font on the label rather
    than on the point. A c:dLbl inherits nothing from the surrounding c:dLbls,
    so the flags of the label being replaced -- else the series, else the plot
    -- are carried over, and the show_* fields override them.
    """

    series_index: int
    point_index: int
    number_format: str
    format_linked: bool
    font_color: str
    font_size_pt: int
    font_bold: bool
    show_value: bool
    show_category: bool
    show_series_name: bool
    show_percent: bool
    show_legend_key: bool
    fill_color: str
    no_fill: bool
    border: ChartLineFormatSpec
    delete: bool


class ChartSeriesInvertSpec(TypedDict, total=False):
    """Negative-value inversion for one series.

    The flag alone leaves PowerPoint drawing negative points in the inverse of
    the series fill, which is usually white; negative_fill_color writes an
    explicit per-point fill instead.
    """

    series_index: int
    invert_if_negative: bool
    negative_fill_color: str


class ChartDataTableSpec(TypedDict, total=False):
    """The c:dTable grid drawn under a chart's plot area.

    show false removes it; the flags default to on when the chart has no data
    table yet, and otherwise keep whatever the chart already had.
    """

    show: bool
    show_horizontal_border: bool
    show_vertical_border: bool
    show_outline: bool
    show_keys: bool
    font_size_pt: int


class ChartSeriesLinesSpec(TypedDict, total=False):
    """The c:serLines connectors of a stacked bar or bar-of-pie chart.

    ``show`` false removes them; a ``line`` without ``show`` restyles the
    connectors already there.
    """

    show: bool
    line: ChartLineFormatSpec


class ChartSeriesFormatSpec(TypedDict, total=False):
    """Series-level fill, line and marker formatting.

    On a line, scatter or radar series the series line is what ``line_color``
    styles, and the markers carry their own fill and outline, so recolouring
    markers needs the ``marker_*`` keys.
    """

    series_index: int
    fill_color: str
    no_fill: bool
    line_color: str
    line_width_emu: int
    line_dash: str
    no_line: bool
    marker_symbol: str
    marker_size: int
    marker_fill_color: str
    marker_line_color: str
    marker_line_width_emu: int
    marker_no_fill: bool
    marker_no_line: bool
    smooth: bool
