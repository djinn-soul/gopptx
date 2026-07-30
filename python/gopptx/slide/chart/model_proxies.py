"""Small chart proxy classes: title, legend, data labels, and plot proxies."""
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol

from .line_format import build_line_format, normalize_hex_color

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ...schemas import ChartFormatUpdate, ChartState


class _ChartProto(Protocol):
    """Structural protocol for the chart object used by proxy classes."""

    def state_get(self, key: str, *, default: object) -> object: ...
    def state_set(self, key: str, value: object) -> None: ...
    def apply_format(self, fmt: ChartFormatUpdate) -> None: ...
    def snapshot(self) -> ChartState: ...


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
        val_map = {
            "top_right": "tr",
            "topright": "tr",
            "top": "t",
            "bottom": "b",
            "left": "l",
            "right": "r",
        }
        norm = val_map.get(value.lower(), value)
        self._chart.state_set("legend_position", norm)
        self._chart.apply_format({"legend_position": norm})

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

    def _value(self, staged_key: str, persisted_key: str, *, default: object) -> object:
        staged = self._chart.state_get(staged_key, default=None)
        if staged is not None:
            return staged
        labels = self._chart.snapshot().get("data_labels", {})
        return labels.get(persisted_key, default)

    @property
    def show_value(self) -> bool:
        """Whether value labels are shown."""
        return bool(self._value("data_label_show_value", "show_value", default=False))

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
        return bool(
            self._value("data_label_show_category", "show_category", default=False)
        )

    @show_category_name.setter
    def show_category_name(self, value: bool) -> None:
        self._chart.state_set("data_label_show_category", value)
        self._chart.apply_format({
            "show_data_labels": True,
            "data_label_show_category": value,
        })

    @property
    def position(self) -> str | None:
        """Return the label position token, if explicitly configured."""
        value = self._value("data_label_position", "position", default=None)
        return value if isinstance(value, str) else None

    @position.setter
    def position(self, value: str) -> None:
        self._chart.state_set("data_label_position", value)
        self._chart.apply_format({
            "show_data_labels": True,
            "data_label_position": value,
        })

    @property
    def show_series_name(self) -> bool:
        """Whether series names are included in data labels."""
        return bool(
            self._value(
                "data_label_show_series_name", "show_series_name", default=False
            )
        )

    @show_series_name.setter
    def show_series_name(self, value: bool) -> None:
        self._chart.state_set("data_label_show_series_name", value)
        self._chart.apply_format({
            "show_data_labels": True,
            "data_label_show_series_name": value,
        })

    @property
    def number_format(self) -> str | None:
        """Return the OOXML number format used by data labels."""
        value = self._value("data_label_number_format", "number_format", default=None)
        return value if isinstance(value, str) else None

    @number_format.setter
    def number_format(self, value: str) -> None:
        self._chart.state_set("data_label_number_format", value)
        self._chart.apply_format({
            "show_data_labels": True,
            "data_label_number_format": value,
        })

    @property
    def word_wrap(self) -> bool | None:
        """Whether data-label text wraps, or ``None`` when not specified."""
        value = self._value("data_label_word_wrap", "word_wrap", default=None)
        return value if isinstance(value, bool) else None

    @word_wrap.setter
    def word_wrap(self, value: bool) -> None:
        self._chart.state_set("data_label_word_wrap", value)
        self._chart.apply_format({
            "show_data_labels": True,
            "data_label_word_wrap": value,
        })

    def set_fill(self, color: str | None = None, *, none: bool = False) -> None:
        """Set the background of every data label (issue #662).

        ``none`` clears the fill, which is how a label is made transparent
        again after one has been set.
        """
        if none:
            self._chart.apply_format({
                "show_data_labels": True,
                "data_label_no_fill": True,
            })
            return
        if color is None:
            raise ValueError("set_fill needs a color unless none is true")
        self._chart.apply_format({
            "show_data_labels": True,
            "data_label_fill_color": normalize_hex_color(color, "fill color"),
        })

    def set_border(
        self,
        *,
        color: str | None = None,
        width_emu: int | None = None,
        dash: str | None = None,
        none: bool | None = None,
    ) -> None:
        """Set the outline drawn around every data label (issue #716)."""
        spec = build_line_format(
            color=color,
            width_emu=width_emu,
            dash=dash,
            none=none,
            name="data label border",
        )
        self._chart.apply_format({
            "show_data_labels": True,
            "data_label_border": spec,
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
        staged = self._chart.state_get("show_data_labels", default=None)
        if staged is not None:
            return bool(staged)
        labels = self._chart.snapshot().get("data_labels", {})
        return bool(labels.get("present", False))

    @data_labels_visible.setter
    def data_labels_visible(self, value: bool) -> None:
        self._chart.state_set("show_data_labels", value)
        self._chart.apply_format({"show_data_labels": value})

    @property
    def data_labels(self) -> DataLabels:
        """Plot data-labels proxy."""
        return self._data_labels

    def set_bar_options(self, *, grouping: str, gap_width: int, overlap: int) -> None:
        """Set grouping and spacing options for bar or column charts."""
        self._chart.state_set("chart_grouping", grouping)
        self._chart.state_set("gap_width", gap_width)
        self._chart.state_set("overlap", overlap)
        self._chart.apply_format({
            "chart_grouping": grouping,
            "gap_width": gap_width,
            "overlap": overlap,
        })


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
