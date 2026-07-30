"""Per-point formatting proxies and negative-value inversion helpers."""

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

if TYPE_CHECKING:
    from collections.abc import Iterator, Sequence

    from ...schemas import ChartDataPointSpec, ChartFormatUpdate, ChartState

_EXPLOSION_MAX = 400
_MARKER_SIZE_MIN = 2
_MARKER_SIZE_MAX = 72

# CT_MarkerStyle values.
MARKER_SYMBOLS = frozenset({
    "auto",
    "circle",
    "dash",
    "diamond",
    "dot",
    "none",
    "plus",
    "square",
    "star",
    "triangle",
    "x",
})

# Optional fields copied straight through to the bridge payload.
_PASS_THROUGH_FIELDS = (
    "fill_color",
    "line_color",
    "line_width_emu",
    "invert_if_negative",
    "bubble_3d",
    "explosion",
    "marker_fill_color",
    "marker_line_color",
    "marker_symbol",
    "marker_size",
)


class _DataPointChartProtocol(Protocol):
    """Minimal chart surface used by data-point helpers."""

    def snapshot(self) -> ChartState:
        """Return the current chart state."""
        ...

    def apply_format(self, fmt: ChartFormatUpdate) -> None:
        """Apply a chart format update."""
        ...


class DataPoint:
    """Read-only proxy for one persisted ``c:dPt``."""

    def __init__(self, payload: dict[str, object]) -> None:
        """Initialize with the raw data-point payload."""
        super().__init__()
        self._payload = payload

    @property
    def series_index(self) -> int:
        """Return the zero-based series this point belongs to."""
        return self._int_field("series_index") or 0

    @property
    def point_index(self) -> int:
        """Return the zero-based index of the formatted point."""
        return self._int_field("point_index") or 0

    @property
    def fill_color(self) -> str | None:
        """Return the point fill colour as a hex RGB string."""
        return self._string_field("fill_color")

    @property
    def line_color(self) -> str | None:
        """Return the point outline colour as a hex RGB string."""
        return self._string_field("line_color")

    @property
    def line_width_emu(self) -> int | None:
        """Return the point outline width in EMU."""
        return self._int_field("line_width_emu")

    @property
    def invert_if_negative(self) -> bool | None:
        """Return whether this point inverts its fill when negative."""
        value = self._payload.get("invert_if_negative")
        return value if isinstance(value, bool) else None

    @property
    def explosion(self) -> int | None:
        """Return the pie/doughnut explosion distance, if set."""
        return self._int_field("explosion")

    @property
    def marker_fill_color(self) -> str | None:
        """Return the marker fill colour as a hex RGB string."""
        return self._string_field("marker_fill_color")

    @property
    def marker_line_color(self) -> str | None:
        """Return the marker outline colour as a hex RGB string."""
        return self._string_field("marker_line_color")

    @property
    def marker_symbol(self) -> str | None:
        """Return the marker symbol, such as ``circle``."""
        return self._string_field("marker_symbol")

    @property
    def marker_size(self) -> int | None:
        """Return the marker size in points."""
        return self._int_field("marker_size")

    def to_payload(self) -> dict[str, object]:
        """Return a copy of the raw payload, ready to re-send unchanged."""
        return dict(self._payload)

    def _string_field(self, key: str) -> str | None:
        value = self._payload.get(key)
        return str(value) if isinstance(value, str) and value else None

    def _int_field(self, key: str) -> int | None:
        value = self._payload.get(key)
        # bool is a subclass of int, so it would otherwise become 1/0.
        if isinstance(value, bool) or not isinstance(value, int):
            return None
        return value


class DataPointCollection:
    """Sequence-like container of data-point proxies."""

    def __init__(self, payload: list[dict[str, object]]) -> None:
        """Initialize the collection with raw data-point payloads."""
        super().__init__()
        self._payload = payload

    def __len__(self) -> int:
        """Return the number of formatted points on the chart."""
        return len(self._payload)

    def __getitem__(self, index: int) -> DataPoint:
        """Return the formatted point at the given index."""
        if index < 0:
            index += len(self._payload)
        if index < 0 or index >= len(self._payload):
            raise IndexError("data point index out of range")
        return DataPoint(self._payload[index])

    def __iter__(self) -> Iterator[DataPoint]:
        """Iterate over data-point proxies."""
        for item in self._payload:
            yield DataPoint(item)

    def for_series(self, series_index: int) -> list[DataPoint]:
        """Return the formatted points belonging to one series."""
        return [point for point in self if point.series_index == series_index]


class ChartDataPointMixin:
    """Adds per-point formatting reads and writes to the live chart proxy."""

    @property
    def data_points(self: _DataPointChartProtocol) -> DataPointCollection:
        """Return every explicitly formatted point on the chart."""
        # The snapshot is a bridge payload, so the declared type is a promise
        # the runtime value has to be checked against.
        raw: object = self.snapshot().get("data_points", [])
        if not isinstance(raw, list):
            return DataPointCollection([])
        items = cast("list[object]", raw)
        return DataPointCollection([
            dict(cast("dict[str, object]", item))
            for item in items
            if isinstance(item, dict)
        ])

    def format_data_point(
        self: _DataPointChartProtocol,
        point_index: int,
        *,
        series_index: int = 0,
        **options: object,
    ) -> None:
        """Merge formatting into one data point, keeping its other properties.

        Args:
            point_index: Zero-based index of the point to format.
            series_index: Zero-based index of the owning series.
            **options: Any of ``fill_color``, ``line_color``, ``line_width_emu``,
                ``invert_if_negative``, ``bubble_3d``, and ``explosion``.

        Raises:
            ValueError: If an index or value is out of range.
        """
        spec = build_data_point_spec(point_index, series_index=series_index, **options)
        self.apply_format(cast("ChartFormatUpdate", {"data_points": [spec]}))

    def format_data_points(
        self: _DataPointChartProtocol,
        points: Sequence[ChartDataPointSpec],
        *,
        series_index: int = 0,
    ) -> None:
        """Merge formatting into several data points in one round-trip."""
        specs = [dict(point) for point in points]
        for spec in specs:
            spec.setdefault("series_index", series_index)
            validate_data_point_spec(spec)
        if not specs:
            return
        self.apply_format(cast("ChartFormatUpdate", {"data_points": specs}))

    def clear_data_point_formatting(
        self: _DataPointChartProtocol, *, series_index: int = 0
    ) -> None:
        """Drop every per-point override from one series."""
        self.apply_format(
            cast("ChartFormatUpdate", {"clear_data_point_series": [series_index]})
        )

    def set_invert_if_negative(
        self: _DataPointChartProtocol,
        *,
        enabled: bool = True,
        series_index: int = 0,
        negative_fill_color: str | None = None,
    ) -> None:
        """Invert negative points, optionally giving them an explicit colour.

        The flag on its own leaves PowerPoint drawing negative points in the
        inverse of the series fill, which is usually white; passing
        ``negative_fill_color`` writes a per-point fill instead.
        """
        spec: dict[str, object] = {
            "series_index": series_index,
            "invert_if_negative": enabled,
        }
        if negative_fill_color is not None:
            spec["negative_fill_color"] = negative_fill_color
        self.apply_format(
            cast("ChartFormatUpdate", {"series_invert_if_negative": [spec]})
        )


def build_data_point_spec(
    point_index: int, *, series_index: int = 0, **options: object
) -> dict[str, object]:
    """Build and validate a bridge data-point payload."""
    spec: dict[str, object] = {
        "series_index": series_index,
        "point_index": point_index,
    }
    for key in _PASS_THROUGH_FIELDS:
        value = options.pop(key, None)
        if value is not None:
            spec[key] = value
    if options:
        unknown = ", ".join(sorted(options))
        raise ValueError(f"unknown data point option(s): {unknown}")
    validate_data_point_spec(spec)
    return spec


def validate_data_point_spec(spec: dict[str, object]) -> None:
    """Reject data-point payloads the chart schema would not accept.

    Raises:
        ValueError: If an index is negative or the explosion is out of range.
    """
    for key in ("series_index", "point_index"):
        value = spec.get(key, 0)
        if isinstance(value, int) and value < 0:
            raise ValueError(f"{key} must not be negative")
    width = spec.get("line_width_emu")
    if isinstance(width, int) and width < 0:
        raise ValueError("line_width_emu must not be negative")
    explosion = spec.get("explosion")
    if explosion is not None and (
        not isinstance(explosion, int) or not 0 <= explosion <= _EXPLOSION_MAX
    ):
        raise ValueError(f"explosion must be between 0 and {_EXPLOSION_MAX}")
    symbol = spec.get("marker_symbol")
    if symbol is not None and symbol not in MARKER_SYMBOLS:
        raise ValueError(
            f"marker_symbol must be one of {', '.join(sorted(MARKER_SYMBOLS))}"
        )
    size = spec.get("marker_size")
    if size is not None and (
        not isinstance(size, int)
        or isinstance(size, bool)
        or not _MARKER_SIZE_MIN <= size <= _MARKER_SIZE_MAX
    ):
        raise ValueError(
            f"marker_size must be between {_MARKER_SIZE_MIN} and {_MARKER_SIZE_MAX}"
        )
