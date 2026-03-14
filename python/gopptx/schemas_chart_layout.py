"""Chart/layout related TypedDict definitions."""

from __future__ import annotations

try:
    from typing import NotRequired, TypedDict
except ImportError:  # pragma: no cover
    from typing_extensions import NotRequired, TypedDict


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


class ChartFormatUpdate(TypedDict, total=False):
    """Chart formatting update payload."""

    show_title: bool
    title: str
    title_overlay: bool
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
    category_axis_tick_label_pos: str
    value_axis_tick_label_pos: str
    category_axis_major_gridlines: bool
    value_axis_major_gridlines: bool
    category_axis_crosses: str
    value_axis_crosses: str
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
    crosses: str


class ChartState(TypedDict, total=False):
    """Chart traversal state snapshot."""

    chart_style: int
    category_axis: ChartAxisState
    value_axis: ChartAxisState
    series: list[ChartSeriesData]
    scene3d: NotRequired[ChartScene3DState]


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
