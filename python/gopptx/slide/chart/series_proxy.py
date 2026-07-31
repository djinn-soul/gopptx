"""Live chart-series proxies and their formatting facade."""
# pyright: reportMissingSuperCall=false

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ...schemas import ChartSeriesFormatSpec


class _SeriesChartProto(Protocol):
    def format_series(self, series_index: int = 0, **options: object) -> None: ...

    def series_format(self, series_index: int = 0) -> ChartSeriesFormatSpec | None: ...


class ChartSeriesLineFormat:
    """Line formatting for one chart series."""

    def __init__(self, chart: _SeriesChartProto, series_index: int) -> None:
        """Bind line formatting to one series index."""
        self._chart = chart
        self._series_index = series_index

    @property
    def dash_style(self) -> str | None:
        """Return the explicit DrawingML dash token, if any."""
        spec = self._chart.series_format(self._series_index)
        if spec is None:
            return None
        value = spec.get("line_dash")
        return str(value) if isinstance(value, str) else None

    @dash_style.setter
    def dash_style(self, value: str) -> None:
        """Set a series dash, including ``MSO_LINE.ROUND_DOT`` (Issue #332)."""
        self._chart.format_series(self._series_index, line_dash=str(value))


class ChartSeriesFormat:
    """Formatting facade for one chart series."""

    def __init__(self, chart: _SeriesChartProto, series_index: int) -> None:
        """Bind the facade to one series index."""
        self._line = ChartSeriesLineFormat(chart, series_index)

    @property
    def line(self) -> ChartSeriesLineFormat:
        """Return line formatting for this series."""
        return self._line


class ChartSeries:
    """Live proxy for a chart series payload."""

    def __init__(
        self,
        chart: _SeriesChartProto,
        series_index: int,
        payload: dict[str, object],
    ) -> None:
        """Bind a series payload and its chart formatting surface."""
        self._payload = payload
        self._format = ChartSeriesFormat(chart, series_index)

    @property
    def name(self) -> str | None:
        """Return the series name, if present."""
        value = self._payload.get("name")
        return str(value) if isinstance(value, str) else None

    @property
    def values(self) -> list[float]:
        """Return numeric series values."""
        raw = self._payload.get("values")
        if not isinstance(raw, list):
            return []
        values_raw = cast("list[object]", raw)
        return [float(item) for item in values_raw if isinstance(item, int | float)]

    @property
    def format(self) -> ChartSeriesFormat:
        """Return this series' formatting facade."""
        return self._format


class ChartSeriesCollection:
    """Sequence-like container of live chart-series proxies."""

    def __init__(
        self,
        chart: _SeriesChartProto,
        payload: list[dict[str, object]],
    ) -> None:
        """Initialize a live series collection."""
        self._chart = chart
        self._payload = payload

    def __len__(self) -> int:
        """Return the number of chart series."""
        return len(self._payload)

    def __getitem__(self, index: int) -> ChartSeries:
        """Return one live series proxy by index."""
        if index < 0:
            index += len(self._payload)
        if index < 0 or index >= len(self._payload):
            raise IndexError("series index out of range")
        return ChartSeries(self._chart, index, self._payload[index])

    def __iter__(self) -> Iterator[ChartSeries]:
        """Iterate live series proxies in chart order."""
        for index, item in enumerate(self._payload):
            yield ChartSeries(self._chart, index, item)
