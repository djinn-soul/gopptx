"""Chart/layout related TypedDict definitions."""

from __future__ import annotations

try:
    from typing import NotRequired, TypedDict
except ImportError:  # pragma: no cover
    from typing_extensions import NotRequired, TypedDict

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from .schemas_chart_series import (
        ChartDataLabelPointSpec,
        ChartDataPointSpec,
        ChartDataTableSpec,
        ChartErrorBarSpec,
        ChartLineFormatSpec,
        ChartSeriesFormatSpec,
        ChartSeriesInvertSpec,
        ChartSeriesLinesSpec,
        ChartTrendlineSpec,
    )


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
    plot_area_line: ChartLineFormatSpec
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
    data_label_fill_color: str
    data_label_no_fill: bool
    data_label_border: ChartLineFormatSpec
    series_formats: list[ChartSeriesFormatSpec]
    series_lines: ChartSeriesLinesSpec
    chart_grouping: str
    gap_width: int
    overlap: int
    category_axis_tick_label_pos: str
    value_axis_tick_label_pos: str
    category_axis_major_gridlines: bool
    value_axis_major_gridlines: bool
    category_axis_minor_gridlines: bool
    value_axis_minor_gridlines: bool
    category_axis_major_gridline_format: ChartLineFormatSpec
    category_axis_minor_gridline_format: ChartLineFormatSpec
    value_axis_major_gridline_format: ChartLineFormatSpec
    value_axis_minor_gridline_format: ChartLineFormatSpec
    category_axis_crosses: str
    value_axis_crosses: str
    category_axis_has_title: bool
    value_axis_has_title: bool
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
    category_axis_visible: bool
    value_axis_visible: bool
    category_axis_tick_label_rotation: float
    value_axis_tick_label_rotation: float
    value_axis_cross_between: str
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
    has_title: bool
    title: str
    minimum_scale: float
    maximum_scale: float
    major_unit: float
    minor_unit: float
    number_format: str
    format_linked: bool
    tick_mark_skip: int
    label_alignment: str
    visible: bool
    tick_label_rotation: float
    cross_between: str
    major_gridline_format: ChartLineFormatSpec
    minor_gridline_format: ChartLineFormatSpec


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
    plot_area_line: NotRequired[ChartLineFormatSpec]
    series_formats: NotRequired[list[ChartSeriesFormatSpec]]
    series_lines: NotRequired[ChartSeriesLinesSpec]


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
    ShapeID: int


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
    LayoutID: NotRequired[int]
    Shapes: NotRequired[list[str]]
    Placeholders: NotRequired[list[PlaceholderInfo]]


class SlideMasterCloneResult(TypedDict):
    """Result of cloning a slide master."""

    MasterPart: str
    ThemePart: str
    LayoutMap: dict[str, str]
