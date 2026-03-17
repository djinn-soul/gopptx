"""Small chart proxy classes: title, legend, data labels, and plot proxies."""
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ...schemas import ChartFormatUpdate


class _ChartProto(Protocol):
    """Structural protocol for the chart object used by proxy classes."""

    def state_get(self, key: str, *, default: object) -> object: ...
    def state_set(self, key: str, value: object) -> None: ...
    def apply_format(self, fmt: ChartFormatUpdate) -> None: ...


class ChartTitle:
    """Chart title proxy."""

    def __init__(self, chart: _ChartProto) -> None:
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

    def __init__(self, chart: _ChartProto) -> None:
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

    def __init__(self, chart: _ChartProto) -> None:
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

    def __init__(self, chart: _ChartProto) -> None:
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

    def __init__(self, chart: _ChartProto) -> None:
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
