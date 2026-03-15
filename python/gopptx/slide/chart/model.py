"""Live chart object model proxies."""
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false, reportUnknownMemberType=false, reportAttributeAccessIssue=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .axis_series import ChartAxis, ChartSeriesCollection
from .scene3d_area import ChartArea

if TYPE_CHECKING:
    from collections.abc import Iterator
    from typing import Protocol

    from ...schemas import ChartDataUpdate, ChartFormatUpdate, ChartState
    from ..slide import Slide
    from .data import CategoryChartData, XyChartData

    class _ChartAxisBulkOpsProto(Protocol):
        def _axes_for(self, axis: str) -> tuple[ChartAxis, ...]:
            ...


class ChartTitle:
    """Chart title proxy."""

    def __init__(self, chart: Chart) -> None:
        """Initialize a title proxy."""
        self._chart = chart

    @property
    def text(self) -> str:
        """Current chart title text."""
        return str(self._chart.state_get("title", default=""))

    @text.setter
    def text(self, value: str) -> None:
        self._chart.state_set("title", value)
        self._chart.apply_format({"show_title": True, "title": value})

    @property
    def visible(self) -> bool:
        """Whether the chart title is visible."""
        return bool(self._chart.state_get("show_title", default=False))

    @visible.setter
    def visible(self, value: bool) -> None:
        self._chart.state_set("show_title", value)
        self._chart.apply_format({"show_title": value})


class ChartLegend:
    """Chart legend proxy."""

    def __init__(self, chart: Chart) -> None:
        """Initialize a legend proxy."""
        self._chart = chart

    @property
    def visible(self) -> bool:
        """Whether the legend is visible."""
        return bool(self._chart.state_get("show_legend", default=True))

    @visible.setter
    def visible(self, value: bool) -> None:
        self._chart.state_set("show_legend", value)
        self._chart.apply_format({"show_legend": value})

    @property
    def position(self) -> str:
        """Legend position code."""
        return str(self._chart.state_get("legend_position", default="r"))

    @position.setter
    def position(self, value: str) -> None:
        self._chart.state_set("legend_position", value)
        self._chart.apply_format({"legend_position": value})

    @property
    def include_in_layout(self) -> bool:
        """Whether legend participates in layout."""
        overlay = bool(self._chart.state_get("legend_overlay", default=False))
        return not overlay

    @include_in_layout.setter
    def include_in_layout(self, value: bool) -> None:
        overlay = not value
        self._chart.state_set("legend_overlay", overlay)
        self._chart.apply_format({"legend_overlay": overlay})


class DataLabels:
    """Data labels proxy for the single chart plot."""

    def __init__(self, chart: Chart) -> None:
        """Initialize data labels proxy."""
        self._chart = chart

    @property
    def show_value(self) -> bool:
        """Whether value labels are shown."""
        return bool(self._chart.state_get("data_label_show_value", default=False))

    @show_value.setter
    def show_value(self, value: bool) -> None:
        self._chart.state_set("data_label_show_value", value)
        self._chart.apply_format({
            "show_data_labels": True,
            "data_label_show_value": value,
        })

    @property
    def show_category_name(self) -> bool:
        """Whether category-name labels are shown."""
        return bool(self._chart.state_get("data_label_show_category", default=False))

    @show_category_name.setter
    def show_category_name(self, value: bool) -> None:
        self._chart.state_set("data_label_show_category", value)
        self._chart.apply_format({
            "show_data_labels": True,
            "data_label_show_category": value,
        })


class ChartPlot:
    """Single plot proxy."""

    def __init__(self, chart: Chart) -> None:
        """Initialize plot proxy."""
        self._chart = chart
        self._data_labels = DataLabels(chart)

    @property
    def data_labels_visible(self) -> bool:
        """Whether plot-level data labels are visible."""
        return bool(self._chart.state_get("show_data_labels", default=False))

    @data_labels_visible.setter
    def data_labels_visible(self, value: bool) -> None:
        self._chart.state_set("show_data_labels", value)
        self._chart.apply_format({"show_data_labels": value})

    @property
    def data_labels(self) -> DataLabels:
        """Plot data-labels proxy."""
        return self._data_labels


class ChartPlots:
    """Collection-like proxy for chart plots."""

    def __init__(self, chart: Chart) -> None:
        """Initialize plots collection proxy."""
        self._chart = chart
        self._plot = ChartPlot(chart)

    def __len__(self) -> int:
        """Return number of plots."""
        return 1

    def __getitem__(self, index: int) -> ChartPlot:
        """Return the only plot for index ``0`` (or ``-1``)."""
        if index not in {0, -1}:
            raise IndexError("plot index out of range")
        return self._plot

    def __iter__(self) -> Iterator[ChartPlot]:
        """Iterate chart plots."""
        yield self._plot


class _ChartAxisBulkOpsMixin:
    """Bulk axis convenience operations."""

    def set_axis_crosses(
        self: _ChartAxisBulkOpsProto, *, crosses: str, axis: str = "both"
    ) -> None:
        """Set axis crossing mode for one or both axes."""
        for target in self._axes_for(axis):
            target.crosses = crosses

    def set_axis_gridlines(
        self: _ChartAxisBulkOpsProto,
        *,
        major: bool | None = None,
        minor: bool | None = None,
        axis: str = "both",
    ) -> None:
        """Set major/minor gridline visibility for one or both axes."""
        if major is None and minor is None:
            raise ValueError("at least one of major/minor must be provided")
        for target in self._axes_for(axis):
            if major is not None:
                target.major_gridlines_visible = major
            if minor is not None:
                target.minor_gridlines_visible = minor


class _ChartStateMixin:
    """Read/write helpers for staged chart state."""

    def state_get(self, key: str, *, default: object) -> object:
        """Read a staged chart state value."""
        state = cast("dict[str, object]", self._state)
        return state.get(key, default)

    def state_set(self, key: str, value: object) -> None:
        """Write a staged chart state value."""
        state = cast("dict[str, object]", self._state)
        state[key] = value


class Chart(_ChartStateMixin, _ChartAxisBulkOpsMixin):
    """Live chart proxy with python-pptx-style accessors."""

    def __init__(
        self,
        slide: Slide,
        index: int,
        rel_id: str,
        chart_part: str,
    ) -> None:
        """Initialize a chart proxy bound to a slide chart relation."""
        self._slide = slide
        self._index = index
        self._rel_id = rel_id
        self._chart_part = chart_part
        self._state: dict[str, object] = {}
        self._chart_title = ChartTitle(self)
        self._legend = ChartLegend(self)
        self._plots = ChartPlots(self)
        self._chart_area = ChartArea(self)
        self._category_axis = ChartAxis(self, axis_name="category")
        self._value_axis = ChartAxis(self, axis_name="value")

    def _snapshot(self) -> ChartState:
        return self._slide.presentation.get_chart_state_by_index(
            self._slide.index,
            self._index,
        )

    def snapshot(self) -> ChartState:
        """Return a fresh chart state snapshot."""
        return self._snapshot()

    @property
    def index(self) -> int:
        """Zero-based chart index within the slide."""
        return self._index

    @property
    def rel_id(self) -> str:
        """Relationship ID for the chart part."""
        return self._rel_id

    @property
    def chart_part(self) -> str:
        """Package path for the chart XML part."""
        return self._chart_part

    @property
    def chart_title(self) -> ChartTitle:
        """Chart title proxy."""
        return self._chart_title

    @property
    def title_visible(self) -> bool:
        """Whether the chart title is visible."""
        return self._chart_title.visible

    @title_visible.setter
    def title_visible(self, value: bool) -> None:
        self._chart_title.visible = value

    @property
    def legend(self) -> ChartLegend:
        """Chart legend proxy."""
        return self._legend

    @property
    def plots(self) -> ChartPlots:
        """Plot collection proxy."""
        return self._plots

    @property
    def chart_area(self) -> ChartArea:
        """Chart-area formatting proxy."""
        return self._chart_area

    @property
    def category_axis(self) -> ChartAxis:
        """Category axis proxy."""
        return self._category_axis

    @property
    def value_axis(self) -> ChartAxis:
        """Value axis proxy."""
        return self._value_axis

    @property
    def axes(self) -> tuple[ChartAxis, ChartAxis]:
        """Tuple of ``(category_axis, value_axis)``."""
        return (self._category_axis, self._value_axis)

    def axis(self, axis_name: str) -> ChartAxis:
        """Resolve an axis by alias name."""
        normalized = axis_name.strip().lower()
        if normalized in {"category", "cat"}:
            return self._category_axis
        if normalized in {"value", "val"}:
            return self._value_axis
        raise ValueError("axis_name must be one of: category, value")

    def _axes_for(self, axis: str) -> tuple[ChartAxis, ...]:
        normalized = axis.strip().lower()
        if normalized in {"both", "all"}:
            return self.axes
        if normalized in {"category", "cat"}:
            return (self._category_axis,)
        if normalized in {"value", "val"}:
            return (self._value_axis,)
        raise ValueError("axis must be one of: both, category, value")

    def set_tick_labels_visibility(self, *, visible: bool) -> None:
        """Toggle tick-label visibility on both axes."""
        value = "nextTo" if visible else "none"
        self._category_axis.tick_label_position = value
        self._value_axis.tick_label_position = value

    @property
    def chart_style(self) -> int | None:
        """Chart style index when present."""
        snapshot = self._snapshot()
        style = snapshot.get("chart_style")
        return int(style) if isinstance(style, int) else None

    @property
    def series(self) -> ChartSeriesCollection:
        """Series collection snapshot proxy."""
        snapshot = self._snapshot()
        return ChartSeriesCollection(
            cast("list[dict[str, object]]", snapshot.get("series", []))
        )

    def replace_data(
        self, data: CategoryChartData | XyChartData | ChartDataUpdate
    ) -> None:
        """Replace chart data from either builder data or raw payload."""
        payload: ChartDataUpdate
        if hasattr(data, "to_update_payload"):
            payload = cast("ChartDataUpdate", data.to_update_payload())
        else:
            payload = cast("ChartDataUpdate", data)
        self._slide.presentation.update_chart_data_by_index(
            self._slide.index,
            self._index,
            payload,
        )

    def _apply_format(self, fmt: ChartFormatUpdate) -> None:
        self._slide.presentation.update_chart_formatting_by_index(
            self._slide.index,
            self._index,
            fmt,
        )

    def apply_format(self, fmt: ChartFormatUpdate) -> None:
        """Apply formatting updates to the chart."""
        self._apply_format(fmt)
