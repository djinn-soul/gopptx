"""Chart/layout related TypedDict definitions."""

from __future__ import annotations

try:
    from typing import NotRequired, TypedDict
except ImportError:  # pragma: no cover
    from typing_extensions import NotRequired, TypedDict


class ChartDataSource(TypedDict, total=False):
    """Where a chart's numbers come from.

    ``kind`` is ``"embedded"``, ``"external"`` or ``"none"``.
    """

    chart_part: str
    kind: str
    rel_id: str
    target: str
    part_path: str
    auto_update: bool


class ChartSelector(TypedDict, total=False):
    """Chart selector for identifying charts."""

    index: int
    rel_id: str


class ChartSeriesData(TypedDict, total=False):
    """Chart series data for updates."""

    name: str
    categories: list[str]
    values: list[float]
    x_values: list[float]
    y_values: list[float]
    sizes: list[float]


class ChartDataUpdate(TypedDict, total=False):
    """Chart data update payload."""

    categories: list[str]
    series: list[ChartSeriesData]


class DataLabelOffset(TypedDict, total=False):
    """Manual position for one data label.

    x and y are fractions of the chart area, offset from the label's default
    spot. Doughnut and pie charts reject most data_label_position values, so
    this is how one of their labels is moved.
    """

    series_index: int
    point_index: int
    x: float
    y: float


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


class ChartFormatUpdate(TypedDict, total=False):
    """Chart formatting update payload."""

    show_title: bool
    title: str
    title_overlay: bool
    title_x: float
    title_y: float
    data_label_offsets: list[DataLabelOffset]
    trendlines: list[ChartTrendlineSpec]
    append_trendlines: list[ChartTrendlineSpec]
    clear_trendline_series: list[int]
    error_bars: list[ChartErrorBarSpec]
    clear_error_bar_series: list[int]
    data_points: list[ChartDataPointSpec]
    data_label_points: list[ChartDataLabelPointSpec]
    clear_data_point_series: list[int]
    series_invert_if_negative: list[ChartSeriesInvertSpec]
    data_table: ChartDataTableSpec
    plot_visible_only: bool
    show_legend: bool
    legend_position: str
    legend_overlay: bool
    show_data_labels: bool
    data_label_position: str
    data_label_show_legend_key: bool
    data_label_show_value: bool
    data_label_show_category: bool
    data_label_show_series_name: bool
    data_label_show_percent: bool
    data_label_show_bubble_size: bool
    data_label_number_format: str
    data_label_format_linked: bool
    data_label_word_wrap: bool
    chart_grouping: str
    gap_width: int
    overlap: int
    category_axis_tick_label_pos: str
    value_axis_tick_label_pos: str
    category_axis_major_gridlines: bool
    value_axis_major_gridlines: bool
    category_axis_minor_gridlines: bool
    value_axis_minor_gridlines: bool
    category_axis_crosses: str
    value_axis_crosses: str
    category_axis_title: str
    value_axis_title: str
    category_axis_minimum_scale: float
    category_axis_maximum_scale: float
    value_axis_minimum_scale: float
    value_axis_maximum_scale: float
    category_axis_major_unit: float
    category_axis_minor_unit: float
    value_axis_major_unit: float
    value_axis_minor_unit: float
    category_axis_number_format: str
    value_axis_number_format: str
    category_axis_tick_mark_skip: int
    category_axis_label_alignment: str
    category_axis_format_linked: bool
    value_axis_format_linked: bool
    camera_preset: str
    camera_field_of_view: int
    light_rig: str
    light_direction: str
    light_rig_revolution: bool


class ChartAxisState(TypedDict, total=False):
    """Chart axis state snapshot."""

    present: bool
    tick_label_pos: str
    major_gridline: bool
    minor_gridline: bool
    crosses: str
    title: str
    minimum_scale: float
    maximum_scale: float
    major_unit: float
    minor_unit: float
    number_format: str
    format_linked: bool
    tick_mark_skip: int
    label_alignment: str


class ChartDataLabelState(TypedDict, total=False):
    """Persisted data-label state for the first chart plot."""

    present: bool
    position: str
    show_value: bool
    show_category: bool
    show_series_name: bool
    number_format: str
    format_linked: bool
    word_wrap: bool


class ChartState(TypedDict, total=False):
    """Chart traversal state snapshot."""

    chart_style: int
    category_axis: ChartAxisState
    value_axis: ChartAxisState
    series: list[ChartSeriesData]
    scene3d: NotRequired[ChartScene3DState]
    data_labels: NotRequired[ChartDataLabelState]
    trendlines: NotRequired[list[ChartTrendlineSpec]]
    error_bars: NotRequired[list[ChartErrorBarSpec]]
    data_points: NotRequired[list[ChartDataPointSpec]]
    data_label_points: NotRequired[list[ChartDataLabelPointSpec]]
    data_table: NotRequired[ChartDataTableSpec]


class ChartScene3DState(TypedDict, total=False):
    """Chart scene3d state snapshot."""

    camera_preset: str
    camera_field_of_view: int
    light_rig: str
    light_direction: str
    light_rig_revolution: bool


class SlideChartRef(TypedDict):
    """Reference to a chart on a slide."""

    Index: int
    RelID: str
    ChartPart: str


class PlaceholderInfo(TypedDict):
    """Placeholder information on layouts/masters."""

    Type: str
    Index: int
    Name: str
    X: float
    Y: float
    CX: float
    CY: float


class SlideLayoutInfo(TypedDict):
    """Information about a slide layout."""

    Part: str
    Name: str
    MasterPart: str
    Shapes: NotRequired[list[str]]
    Placeholders: NotRequired[list[PlaceholderInfo]]


class SlideMasterCloneResult(TypedDict):
    """Result of cloning a slide master."""

    MasterPart: str
    ThemePart: str
    LayoutMap: dict[str, str]
