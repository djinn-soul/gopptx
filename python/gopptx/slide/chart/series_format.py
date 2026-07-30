"""Series-level fill, line and marker formatting, plus series lines."""

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

from .line_format import build_line_format, normalize_hex_color

if TYPE_CHECKING:
    from ...schemas import (
        ChartFormatUpdate,
        ChartSeriesFormatSpec,
        ChartState,
    )

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

_MARKER_SIZE_MIN = 2
_MARKER_SIZE_MAX = 72

_COLOR_FIELDS = (
    "fill_color",
    "line_color",
    "marker_fill_color",
    "marker_line_color",
)

_WIDTH_FIELDS = ("line_width_emu", "marker_line_width_emu")

_FLAG_FIELDS = ("no_fill", "no_line", "marker_no_fill", "marker_no_line", "smooth")

_OTHER_FIELDS = ("line_dash", "marker_symbol", "marker_size")

_PASS_THROUGH_FIELDS = _COLOR_FIELDS + _WIDTH_FIELDS + _FLAG_FIELDS + _OTHER_FIELDS


class _SeriesFormatChartProto(Protocol):
    """Chart surface used by the series-formatting helpers."""

    def apply_format(self, fmt: ChartFormatUpdate) -> None: ...
    def snapshot(self) -> ChartState: ...


class ChartSeriesFormatMixin:
    """Adds series-level formatting and series-line controls to a chart."""

    def format_series(
        self: _SeriesFormatChartProto,
        series_index: int = 0,
        **options: object,
    ) -> None:
        """Format a whole series: its fill, its line and its markers (#872).

        On a line, scatter or radar series ``line_color`` styles the line
        joining the points; the markers have their own fill and outline, so
        recolouring them needs the ``marker_*`` options.

        Args:
            series_index: Zero-based index of the series to format.
            **options: Any of ``fill_color``, ``no_fill``, ``line_color``,
                ``line_width_emu``, ``line_dash``, ``no_line``,
                ``marker_symbol``, ``marker_size``, ``marker_fill_color``,
                ``marker_line_color``, ``marker_line_width_emu``,
                ``marker_no_fill``, ``marker_no_line`` and ``smooth``.

        Raises:
            ValueError: If an option is unknown or its value out of range.
        """
        spec = build_series_format_spec(series_index, **options)
        self.apply_format(cast("ChartFormatUpdate", {"series_formats": [spec]}))

    def series_format(
        self: _SeriesFormatChartProto, series_index: int = 0
    ) -> ChartSeriesFormatSpec | None:
        """Return the explicit formatting of one series, if it carries any."""
        raw: object = self.snapshot().get("series_formats", [])
        if not isinstance(raw, list):
            return None
        for entry in cast("list[object]", raw):
            if not isinstance(entry, dict):
                continue
            payload = cast("ChartSeriesFormatSpec", entry)
            if payload.get("series_index", 0) == series_index:
                return payload
        return None

    def set_series_lines(
        self: _SeriesFormatChartProto,
        *,
        show: bool | None = None,
        color: str | None = None,
        width_emu: int | None = None,
        dash: str | None = None,
        none: bool | None = None,
    ) -> None:
        """Draw or style the connectors between stacked bars (#846).

        Only a stacked bar chart and a bar-of-pie chart have series lines;
        other plot types ignore the request.
        """
        payload: dict[str, object] = {}
        if show is not None:
            payload["show"] = bool(show)
        if color is not None or width_emu is not None or dash is not None or none:
            payload["line"] = build_line_format(
                color=color,
                width_emu=width_emu,
                dash=dash,
                none=none,
                name="series lines",
            )
        if not payload:
            raise ValueError("set_series_lines needs show or a line style")
        self.apply_format(cast("ChartFormatUpdate", {"series_lines": payload}))

    def series_lines(self: _SeriesFormatChartProto) -> dict[str, object] | None:
        """Return the series-line state of the first plot that draws them."""
        raw: object = self.snapshot().get("series_lines")
        if not isinstance(raw, dict):
            return None
        return cast("dict[str, object]", raw)


def build_series_format_spec(
    series_index: int = 0, **options: object
) -> ChartSeriesFormatSpec:
    """Build and validate a bridge series-format payload."""
    if series_index < 0:
        raise ValueError("series_index must not be negative")
    spec: dict[str, object] = {"series_index": int(series_index)}
    for key in _PASS_THROUGH_FIELDS:
        value = options.pop(key, None)
        if value is not None:
            spec[key] = value
    if options:
        unknown = ", ".join(sorted(options))
        raise ValueError(f"unknown series format option(s): {unknown}")
    if len(spec) == 1:
        raise ValueError("format_series needs at least one formatting option")
    _validate_series_format_spec(spec)
    return cast("ChartSeriesFormatSpec", spec)


def _validate_series_format_spec(spec: dict[str, object]) -> None:
    for key in _COLOR_FIELDS:
        value = spec.get(key)
        if value is not None:
            spec[key] = normalize_hex_color(str(value), key)
    for key in _WIDTH_FIELDS:
        value = spec.get(key)
        if value is not None and int(cast("int", value)) < 0:
            raise ValueError(f"{key} must not be negative")
    for key in _FLAG_FIELDS:
        value = spec.get(key)
        if value is not None:
            spec[key] = bool(value)
    _validate_series_marker(spec)
    dash = spec.get("line_dash")
    if dash is not None and not str(dash).strip():
        raise ValueError("line_dash must not be empty")


def _validate_series_marker(spec: dict[str, object]) -> None:
    symbol = spec.get("marker_symbol")
    if symbol is not None and str(symbol) not in MARKER_SYMBOLS:
        raise ValueError(
            "marker_symbol must be one of: " + ", ".join(sorted(MARKER_SYMBOLS))
        )
    size = spec.get("marker_size")
    if size is None:
        return
    if not _MARKER_SIZE_MIN <= int(cast("int", size)) <= _MARKER_SIZE_MAX:
        raise ValueError(
            f"marker_size must be between {_MARKER_SIZE_MIN} and {_MARKER_SIZE_MAX}"
        )
