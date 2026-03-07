"""Live chart object model proxies."""
# ruff: noqa: D101,D102,D105,D107,SLF001,TC003
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false, reportUnknownMemberType=false, reportAttributeAccessIssue=false

from __future__ import annotations

from collections.abc import Iterator
from typing import TYPE_CHECKING, cast

if TYPE_CHECKING:
    from ..schemas import ChartDataUpdate, ChartFormatUpdate
    from .chart_data import CategoryChartData, XyChartData
    from .slide import Slide


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
    def has_data_labels(self) -> bool:
        return bool(self._chart._state.get("show_data_labels", False))

    @has_data_labels.setter
    def has_data_labels(self, value: bool) -> None:
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
    def has_title(self) -> bool:
        return self._chart_title.visible

    @has_title.setter
    def has_title(self, value: bool) -> None:
        self._chart_title.visible = value

    @property
    def legend(self) -> ChartLegend:
        return self._legend

    @property
    def plots(self) -> ChartPlots:
        return self._plots

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


class ChartCollection:
    def __init__(self, slide: Slide) -> None:
        self._slide = slide

    def _items(self) -> list[Chart]:
        refs = self._slide.list_charts()
        out: list[Chart] = []
        for item in refs:
            ref = item
            index = int(ref.get("Index", ref.get("index", 0)))
            rel_id = str(ref.get("RelID", ref.get("rel_id", "")))
            chart_part = str(ref.get("ChartPart", ref.get("chart_part", "")))
            out.append(Chart(self._slide, index, rel_id, chart_part))
        out.sort(key=lambda chart: chart.index)
        return out

    def __len__(self) -> int:
        return len(self._items())

    def __getitem__(self, index: int) -> Chart:
        items = self._items()
        if index < 0:
            index += len(items)
        if index < 0 or index >= len(items):
            raise IndexError("chart index out of range")
        return items[index]

    def __iter__(self) -> Iterator[Chart]:
        return iter(self._items())
