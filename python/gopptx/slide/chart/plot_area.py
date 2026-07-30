"""Plot-area formatting proxies."""

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

from .line_format import normalize_hex_color

if TYPE_CHECKING:
    from ...schemas import ChartFormatUpdate, ChartLineFormatSpec, ChartState


class _PlotAreaChartProtocol(Protocol):
    def apply_format(self, fmt: ChartFormatUpdate) -> None: ...
    def snapshot(self) -> ChartState: ...


class ChartPlotAreaLine:
    """Line-format proxy for the outline around a chart plot area."""

    def __init__(self, chart: _PlotAreaChartProtocol) -> None:
        """Bind the outline proxy to a chart."""
        super().__init__()
        self._chart = chart

    def _state(self) -> ChartLineFormatSpec:
        raw: object = self._chart.snapshot().get("plot_area_line", {})
        if not isinstance(raw, dict):
            return {}
        return cast("ChartLineFormatSpec", raw)

    def _apply(self, **changes: object) -> None:
        merged: dict[str, object] = dict(self._state())
        if any(key in changes for key in ("color", "width_emu", "dash")):
            merged.pop("none", None)
        merged.update(changes)
        self._chart.apply_format(cast("ChartFormatUpdate", {"plot_area_line": merged}))

    @property
    def color(self) -> str | None:
        """Return the explicit RGB outline colour."""
        value = self._state().get("color")
        return str(value) if isinstance(value, str) else None

    @color.setter
    def color(self, value: str) -> None:
        self._apply(color=normalize_hex_color(value, "plot area line color"))

    @property
    def width(self) -> int | None:
        """Return the outline width in EMU."""
        value = self._state().get("width_emu")
        return int(value) if isinstance(value, int) else None

    @width.setter
    def width(self, value: int) -> None:
        if value < 0:
            raise ValueError("plot area line width must not be negative")
        self._apply(width_emu=int(value))

    @property
    def dash_style(self) -> str | None:
        """Return the explicit DrawingML dash token."""
        value = self._state().get("dash")
        return str(value) if isinstance(value, str) else None

    @dash_style.setter
    def dash_style(self, value: str) -> None:
        if not value.strip():
            raise ValueError("plot area line dash_style must not be empty")
        self._apply(dash=value.strip())

    def hide(self) -> None:
        """Draw no plot-area outline."""
        self._chart.apply_format(
            cast("ChartFormatUpdate", {"plot_area_line": {"none": True}})
        )


class ChartPlotAreaFormat:
    """Formatting facade for a chart plot area."""

    def __init__(self, chart: _PlotAreaChartProtocol) -> None:
        """Bind the format facade to a chart."""
        super().__init__()
        self._line = ChartPlotAreaLine(chart)

    @property
    def line(self) -> ChartPlotAreaLine:
        """Return the plot-area outline proxy."""
        return self._line


class ChartPlotArea:
    """Plot-area proxy exposing its format."""

    def __init__(self, chart: _PlotAreaChartProtocol) -> None:
        """Bind the plot-area proxy to a chart."""
        super().__init__()
        self._format = ChartPlotAreaFormat(chart)

    @property
    def format(self) -> ChartPlotAreaFormat:
        """Return plot-area formatting."""
        return self._format
