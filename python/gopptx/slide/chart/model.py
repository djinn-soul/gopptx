"""Live chart object model proxies."""
# ruff: noqa: D101,D102,D105,D107,SLF001
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false, reportUnknownMemberType=false, reportAttributeAccessIssue=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .axis_series import ChartAxis, ChartSeriesCollection
from .scene3d_area import ChartArea

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ...schemas import ChartDataUpdate, ChartFormatUpdate, ChartState
    from ..chart_data import CategoryChartData, XyChartData
    from ..slide import Slide


class ChartTitle:
    def __init__(self, chart: Chart) -> None:
        self._chart = chart

    @property
    def text(self) -> str:
        return str(self._chart._state.get("title", ""))

    @text.setter
    def text(self, value: str) -> None:
        self._chart._state["title"] = value
        self._chart._apply_format({"show_title": True, "title": value})

    @property
    def visible(self) -> bool:
        return bool(self._chart._state.get("show_title", False))

    @visible.setter
    def visible(self, value: bool) -> None:
        self._chart._state["show_title"] = value
        self._chart._apply_format({"show_title": value})


class ChartLegend:
    def __init__(self, chart: Chart) -> None:
        self._chart = chart

    @property
    def visible(self) -> bool:
        return bool(self._chart._state.get("show_legend", True))

    @visible.setter
    def visible(self, value: bool) -> None:
        self._chart._state["show_legend"] = value
        self._chart._apply_format({"show_legend": value})

    @property
    def position(self) -> str:
        return str(self._chart._state.get("legend_position", "r"))

    @position.setter
    def position(self, value: str) -> None:
        self._chart._state["legend_position"] = value
        self._chart._apply_format({"legend_position": value})

    @property
    def include_in_layout(self) -> bool:
        overlay = bool(self._chart._state.get("legend_overlay", False))
        return not overlay

    @include_in_layout.setter
    def include_in_layout(self, value: bool) -> None:
        overlay = not value
        self._chart._state["legend_overlay"] = overlay
        self._chart._apply_format({"legend_overlay": overlay})


class DataLabels:
    def __init__(self, chart: Chart) -> None:
        self._chart = chart

    @property
    def show_value(self) -> bool:
        return bool(self._chart._state.get("data_label_show_value", False))

    @show_value.setter
    def show_value(self, value: bool) -> None:
        self._chart._state["data_label_show_value"] = value
        self._chart._apply_format({
            "show_data_labels": True,
            "data_label_show_value": value,
        })

    @property
    def show_category_name(self) -> bool:
        return bool(self._chart._state.get("data_label_show_category", False))

    @show_category_name.setter
    def show_category_name(self, value: bool) -> None:
        self._chart._state["data_label_show_category"] = value
        self._chart._apply_format({
            "show_data_labels": True,
            "data_label_show_category": value,
        })


class ChartPlot:
    def __init__(self, chart: Chart) -> None:
        self._chart = chart
        self._data_labels = DataLabels(chart)

    @property
    def data_labels_visible(self) -> bool:
        return bool(self._chart._state.get("show_data_labels", False))

    @data_labels_visible.setter
    def data_labels_visible(self, value: bool) -> None:
        self._chart._state["show_data_labels"] = value
        self._chart._apply_format({"show_data_labels": value})

    @property
    def data_labels(self) -> DataLabels:
        return self._data_labels


class ChartPlots:
    def __init__(self, chart: Chart) -> None:
        self._chart = chart
        self._plot = ChartPlot(chart)

    def __len__(self) -> int:
        return 1

    def __getitem__(self, index: int) -> ChartPlot:
        if index not in {0, -1}:
            raise IndexError("plot index out of range")
        return self._plot

    def __iter__(self) -> Iterator[ChartPlot]:
        yield self._plot


class Chart:
    def __init__(
        self,
        slide: Slide,
        index: int,
        rel_id: str,
        chart_part: str,
    ) -> None:
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
        return self._slide._presentation.get_chart_state_by_index(
            self._slide.index,
            self._index,
        )

    @property
    def index(self) -> int:
        return self._index

    @property
    def rel_id(self) -> str:
        return self._rel_id

    @property
    def chart_part(self) -> str:
        return self._chart_part

    @property
    def chart_title(self) -> ChartTitle:
        return self._chart_title

    @property
    def title_visible(self) -> bool:
        return self._chart_title.visible

    @title_visible.setter
    def title_visible(self, value: bool) -> None:
        self._chart_title.visible = value

    @property
    def legend(self) -> ChartLegend:
        return self._legend

    @property
    def plots(self) -> ChartPlots:
        return self._plots

    @property
    def chart_area(self) -> ChartArea:
        return self._chart_area

    @property
    def category_axis(self) -> ChartAxis:
        return self._category_axis

    @property
    def value_axis(self) -> ChartAxis:
        return self._value_axis

    @property
    def axes(self) -> tuple[ChartAxis, ChartAxis]:
        return (self._category_axis, self._value_axis)

    def axis(self, axis_name: str) -> ChartAxis:
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
        value = "nextTo" if visible else "none"
        self._category_axis.tick_label_position = value
        self._value_axis.tick_label_position = value

    def set_axis_crosses(self, *, crosses: str, axis: str = "both") -> None:
        for target in self._axes_for(axis):
            target.crosses = crosses

    def set_axis_gridlines(
        self,
        *,
        major: bool | None = None,
        minor: bool | None = None,
        axis: str = "both",
    ) -> None:
        if major is None and minor is None:
            raise ValueError("at least one of major/minor must be provided")
        for target in self._axes_for(axis):
            if major is not None:
                target.major_gridlines_visible = major
            if minor is not None:
                target.minor_gridlines_visible = minor

    @property
    def chart_style(self) -> int | None:
        snapshot = self._snapshot()
        style = snapshot.get("chart_style")
        return int(style) if isinstance(style, int) else None

    @property
    def series(self) -> ChartSeriesCollection:
        snapshot = self._snapshot()
        return ChartSeriesCollection(
            cast("list[dict[str, object]]", snapshot.get("series", []))
        )

    def replace_data(
        self, data: CategoryChartData | XyChartData | ChartDataUpdate
    ) -> None:
        payload: ChartDataUpdate
        if hasattr(data, "to_update_payload"):
            payload = cast("ChartDataUpdate", data.to_update_payload())
        else:
            payload = cast("ChartDataUpdate", data)
        self._slide._presentation.update_chart_data_by_index(
            self._slide.index,
            self._index,
            payload,
        )

    def _apply_format(self, fmt: ChartFormatUpdate) -> None:
        self._slide._presentation.update_chart_formatting_by_index(
            self._slide.index,
            self._index,
            fmt,
        )
