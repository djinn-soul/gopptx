"""Trendline proxies and chart-level trendline operations."""

from __future__ import annotations

from math import isfinite
from typing import TYPE_CHECKING, Protocol, cast

if TYPE_CHECKING:
    from collections.abc import Iterator, Sequence

    from ...schemas import ChartFormatUpdate, ChartState, ChartTrendlineSpec

TRENDLINE_TYPES = (
    "linear",
    "poly",
    "exp",
    "log",
    "movingAvg",
    "power",
)

_POLY_ORDER_MIN = 2
_POLY_ORDER_MAX = 6
_MOVING_AVG_PERIOD_MIN = 2

# Optional fields copied straight through to the bridge payload.
_PASS_THROUGH_FIELDS = (
    "name",
    "forward",
    "backward",
    "intercept",
    "display_r_squared",
    "display_equation",
    "line_color",
    "line_width_emu",
    "line_dash",
)


class _TrendlineChartProtocol(Protocol):
    """Minimal chart surface used by trendline helpers."""

    def snapshot(self) -> ChartState:
        """Return the current chart state."""
        ...

    def apply_format(self, fmt: ChartFormatUpdate) -> None:
        """Apply a chart format update."""
        ...


class Trendline:
    """Read-only proxy for one persisted ``c:trendline``."""

    def __init__(self, payload: dict[str, object]) -> None:
        """Initialize with the raw trendline payload."""
        super().__init__()
        self._payload = payload

    @property
    def series_index(self) -> int:
        """Return the zero-based series this trendline belongs to."""
        value = self._payload.get("series_index", 0)
        return value if isinstance(value, int) and not isinstance(value, bool) else 0

    @property
    def trendline_type(self) -> str | None:
        """Return the trendline type token."""
        value = self._payload.get("type")
        return str(value) if isinstance(value, str) and value else None

    @property
    def order(self) -> int | None:
        """Return the polynomial order, if this is a poly trendline."""
        return self._int_field("order")

    @property
    def period(self) -> int | None:
        """Return the averaging period, if this is a movingAvg trendline."""
        return self._int_field("period")

    @property
    def name(self) -> str | None:
        """Return the trendline label text, if set."""
        value = self._payload.get("name")
        return str(value) if isinstance(value, str) else None

    @property
    def display_r_squared(self) -> bool | None:
        """Return whether the R-squared value is shown on the chart."""
        value = self._payload.get("display_r_squared")
        return value if isinstance(value, bool) else None

    @property
    def display_equation(self) -> bool | None:
        """Return whether the fit equation is shown on the chart."""
        value = self._payload.get("display_equation")
        return value if isinstance(value, bool) else None

    def to_payload(self) -> dict[str, object]:
        """Return a copy of the raw payload, ready to re-send unchanged."""
        return dict(self._payload)

    def _int_field(self, key: str) -> int | None:
        value = self._payload.get(key)
        # bool is a subclass of int, so it would otherwise become 1/0.
        if isinstance(value, bool) or not isinstance(value, int):
            return None
        return value


class TrendlineCollection:
    """Sequence-like container of trendline proxies."""

    def __init__(self, payload: list[dict[str, object]]) -> None:
        """Initialize the collection with raw trendline payloads."""
        super().__init__()
        self._payload = payload

    def __len__(self) -> int:
        """Return the number of trendlines on the chart."""
        return len(self._payload)

    def __getitem__(self, index: int) -> Trendline:
        """Return the trendline at the given index."""
        if index < 0:
            index += len(self._payload)
        if index < 0 or index >= len(self._payload):
            raise IndexError("trendline index out of range")
        return Trendline(self._payload[index])

    def __iter__(self) -> Iterator[Trendline]:
        """Iterate over trendline proxies."""
        for item in self._payload:
            yield Trendline(item)

    def for_series(self, series_index: int) -> list[Trendline]:
        """Return the trendlines belonging to one series."""
        return [line for line in self if line.series_index == series_index]


class ChartTrendlineMixin:
    """Adds trendline reads and writes to the live chart proxy."""

    @property
    def trendlines(self: _TrendlineChartProtocol) -> TrendlineCollection:
        """Return every trendline currently persisted on the chart."""
        return TrendlineCollection(_trendline_payloads(self))

    def add_trendline(
        self: _TrendlineChartProtocol,
        trendline_type: str,
        *,
        series_index: int = 0,
        **options: object,
    ) -> None:
        """Add one trendline to a series, keeping the series' existing ones.

        Args:
            trendline_type: One of ``linear``, ``poly``, ``exp``, ``log``,
                ``movingAvg``, or ``power``.
            series_index: Zero-based index of the series to fit.
            **options: Any of ``order``, ``period``, ``name``, ``forward``,
                ``backward``, ``intercept``, ``display_r_squared``,
                ``display_equation``, ``line_color``, ``line_width_emu``, and
                ``line_dash``.

        Raises:
            ValueError: If the type, order, or period is not valid.
        """
        spec = build_trendline_spec(
            trendline_type, series_index=series_index, **options
        )
        self.apply_format(cast("ChartFormatUpdate", {"append_trendlines": [spec]}))

    def set_trendlines(
        self: _TrendlineChartProtocol,
        trendlines: Sequence[ChartTrendlineSpec],
        *,
        series_index: int = 0,
    ) -> None:
        """Replace every trendline on one series with the given specs."""
        specs = [dict(spec) for spec in trendlines]
        for spec in specs:
            spec.setdefault("series_index", series_index)
            validate_trendline_spec(spec)
        if not specs:
            self.apply_format(
                cast("ChartFormatUpdate", {"clear_trendline_series": [series_index]})
            )
            return
        self.apply_format(cast("ChartFormatUpdate", {"trendlines": specs}))

    def clear_trendlines(
        self: _TrendlineChartProtocol, *, series_index: int = 0
    ) -> None:
        """Remove every trendline from one series."""
        self.apply_format(
            cast("ChartFormatUpdate", {"clear_trendline_series": [series_index]})
        )


def _trendline_payloads(chart: _TrendlineChartProtocol) -> list[dict[str, object]]:
    """Return the raw trendline payloads from a fresh chart snapshot."""
    # The snapshot is a bridge payload, so the declared type is a promise the
    # runtime value has to be checked against.
    raw: object = chart.snapshot().get("trendlines", [])
    if not isinstance(raw, list):
        return []
    items = cast("list[object]", raw)
    return [
        dict(cast("dict[str, object]", item))
        for item in items
        if isinstance(item, dict)
    ]


def build_trendline_spec(
    trendline_type: str, *, series_index: int = 0, **options: object
) -> dict[str, object]:
    """Build and validate a bridge trendline payload."""
    spec: dict[str, object] = {
        "series_index": series_index,
        "type": trendline_type.strip(),
    }
    for key in ("order", "period", *_PASS_THROUGH_FIELDS):
        value = options.pop(key, None)
        if value is not None:
            spec[key] = value
    if options:
        unknown = ", ".join(sorted(options))
        raise ValueError(f"unknown trendline option(s): {unknown}")
    validate_trendline_spec(spec)
    return spec


def validate_trendline_spec(spec: dict[str, object]) -> None:
    """Reject trendline payloads PowerPoint would repair on open.

    Raises:
        ValueError: If the type is unknown, the series index is negative, or
            ``order``/``period`` is paired with the wrong trendline type.
    """
    kind = spec.get("type")
    if not isinstance(kind, str) or kind not in TRENDLINE_TYPES:
        joined = ", ".join(TRENDLINE_TYPES)
        raise ValueError(f"trendline type must be one of: {joined}")
    series_index = spec.get("series_index", 0)
    if isinstance(series_index, int) and series_index < 0:
        raise ValueError("series_index must not be negative")
    _validate_order(kind, spec.get("order"))
    _validate_period(kind, spec.get("period"))
    for field in ("forward", "backward", "intercept"):
        value = spec.get(field)
        if value is None:
            continue
        if (
            isinstance(value, bool)
            or not isinstance(value, int | float)
            or not isfinite(float(value))
        ):
            raise ValueError(f"{field} must be a finite number")


def _validate_order(kind: str, order: object) -> None:
    if order is None:
        return
    if kind != "poly":
        raise ValueError("order is only valid for the poly trendline type")
    if not isinstance(order, int) or not _POLY_ORDER_MIN <= order <= _POLY_ORDER_MAX:
        raise ValueError(
            f"order must be between {_POLY_ORDER_MIN} and {_POLY_ORDER_MAX}"
        )


def _validate_period(kind: str, period: object) -> None:
    if period is None:
        return
    if kind != "movingAvg":
        raise ValueError("period is only valid for the movingAvg trendline type")
    if not isinstance(period, int) or period < _MOVING_AVG_PERIOD_MIN:
        raise ValueError(
            f"period must be greater than or equal to {_MOVING_AVG_PERIOD_MIN}"
        )
