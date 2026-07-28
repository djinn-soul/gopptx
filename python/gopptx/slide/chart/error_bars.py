"""Error-bar proxies and chart-level error-bar operations."""

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

if TYPE_CHECKING:
    from collections.abc import Iterator, Sequence

    from ...schemas import ChartErrorBarSpec, ChartFormatUpdate, ChartState

ERROR_BAR_TYPES = ("both", "minus", "plus")
ERROR_BAR_VALUE_TYPES = ("cust", "fixedVal", "percentage", "stdDev", "stdErr")
ERROR_BAR_DIRECTIONS = ("x", "y")

_CUSTOM_VALUE_TYPE = "cust"

# Optional fields copied straight through to the bridge payload.
_PASS_THROUGH_FIELDS = (
    "direction",
    "value",
    "plus_reference",
    "minus_reference",
    "no_end_cap",
    "line_color",
)


class _ErrorBarChartProtocol(Protocol):
    """Minimal chart surface used by error-bar helpers."""

    def snapshot(self) -> ChartState:
        """Return the current chart state."""
        ...

    def apply_format(self, fmt: ChartFormatUpdate) -> None:
        """Apply a chart format update."""
        ...


class ErrorBars:
    """Read-only proxy for one persisted ``c:errBars``."""

    def __init__(self, payload: dict[str, object]) -> None:
        """Initialize with the raw error-bar payload."""
        super().__init__()
        self._payload = payload

    @property
    def series_index(self) -> int:
        """Return the zero-based series these error bars belong to."""
        value = self._payload.get("series_index", 0)
        return value if isinstance(value, int) and not isinstance(value, bool) else 0

    @property
    def bar_type(self) -> str | None:
        """Return the error-bar direction token (``both``/``minus``/``plus``)."""
        return self._string_field("bar_type")

    @property
    def value_type(self) -> str | None:
        """Return how the error amount is computed."""
        return self._string_field("value_type")

    @property
    def direction(self) -> str | None:
        """Return the ``x`` or ``y`` axis these bars apply to, if set."""
        return self._string_field("direction")

    @property
    def value(self) -> float | None:
        """Return the fixed error amount, for non-custom error bars."""
        value = self._payload.get("value")
        if isinstance(value, bool) or not isinstance(value, int | float):
            return None
        return float(value)

    @property
    def plus_reference(self) -> str | None:
        """Return the positive custom-value formula, if set."""
        return self._string_field("plus_reference")

    @property
    def minus_reference(self) -> str | None:
        """Return the negative custom-value formula, if set."""
        return self._string_field("minus_reference")

    @property
    def no_end_cap(self) -> bool | None:
        """Return whether the bars are drawn without end caps."""
        value = self._payload.get("no_end_cap")
        return value if isinstance(value, bool) else None

    def to_payload(self) -> dict[str, object]:
        """Return a copy of the raw payload, ready to re-send unchanged."""
        return dict(self._payload)

    def _string_field(self, key: str) -> str | None:
        value = self._payload.get(key)
        return str(value) if isinstance(value, str) and value else None


class ErrorBarCollection:
    """Sequence-like container of error-bar proxies."""

    def __init__(self, payload: list[dict[str, object]]) -> None:
        """Initialize the collection with raw error-bar payloads."""
        super().__init__()
        self._payload = payload

    def __len__(self) -> int:
        """Return the number of error-bar sets on the chart."""
        return len(self._payload)

    def __getitem__(self, index: int) -> ErrorBars:
        """Return the error-bar set at the given index."""
        if index < 0:
            index += len(self._payload)
        if index < 0 or index >= len(self._payload):
            raise IndexError("error bar index out of range")
        return ErrorBars(self._payload[index])

    def __iter__(self) -> Iterator[ErrorBars]:
        """Iterate over error-bar proxies."""
        for item in self._payload:
            yield ErrorBars(item)

    def for_series(self, series_index: int) -> list[ErrorBars]:
        """Return the error-bar sets belonging to one series."""
        return [bars for bars in self if bars.series_index == series_index]


class ChartErrorBarMixin:
    """Adds error-bar reads and writes to the live chart proxy."""

    @property
    def error_bars(self: _ErrorBarChartProtocol) -> ErrorBarCollection:
        """Return every error-bar set currently persisted on the chart."""
        return ErrorBarCollection(_error_bar_payloads(self))

    def add_error_bars(
        self: _ErrorBarChartProtocol,
        bar_type: str,
        value_type: str,
        *,
        series_index: int = 0,
        **options: object,
    ) -> None:
        """Add one error-bar set to a series, keeping the series' existing ones.

        Args:
            bar_type: One of ``both``, ``minus``, or ``plus``.
            value_type: One of ``cust``, ``fixedVal``, ``percentage``,
                ``stdDev``, or ``stdErr``.
            series_index: Zero-based index of the series to annotate.
            **options: Any of ``direction``, ``value``, ``plus_reference``,
                ``minus_reference``, ``no_end_cap``, and ``line_color``.

        Raises:
            ValueError: If the combination is one PowerPoint would repair.
        """
        spec = build_error_bar_spec(
            bar_type, value_type, series_index=series_index, **options
        )
        # The bridge replaces every error bar on an addressed series, so the
        # existing ones are re-sent alongside the new set.
        existing = [
            payload
            for payload in _error_bar_payloads(self)
            if payload.get("series_index", 0) == series_index
        ]
        self.apply_format(cast("ChartFormatUpdate", {"error_bars": [*existing, spec]}))

    def set_error_bars(
        self: _ErrorBarChartProtocol,
        error_bars: Sequence[ChartErrorBarSpec],
        *,
        series_index: int = 0,
    ) -> None:
        """Replace every error-bar set on one series with the given specs."""
        specs = [dict(spec) for spec in error_bars]
        for spec in specs:
            spec.setdefault("series_index", series_index)
            validate_error_bar_spec(spec)
        if not specs:
            self.apply_format(
                cast("ChartFormatUpdate", {"clear_error_bar_series": [series_index]})
            )
            return
        self.apply_format(cast("ChartFormatUpdate", {"error_bars": specs}))

    def clear_error_bars(
        self: _ErrorBarChartProtocol, *, series_index: int = 0
    ) -> None:
        """Remove every error-bar set from one series."""
        self.apply_format(
            cast("ChartFormatUpdate", {"clear_error_bar_series": [series_index]})
        )


def _error_bar_payloads(chart: _ErrorBarChartProtocol) -> list[dict[str, object]]:
    """Return the raw error-bar payloads from a fresh chart snapshot."""
    # The snapshot is a bridge payload, so the declared type is a promise the
    # runtime value has to be checked against.
    raw: object = chart.snapshot().get("error_bars", [])
    if not isinstance(raw, list):
        return []
    items = cast("list[object]", raw)
    return [
        dict(cast("dict[str, object]", item))
        for item in items
        if isinstance(item, dict)
    ]


def build_error_bar_spec(
    bar_type: str, value_type: str, *, series_index: int = 0, **options: object
) -> dict[str, object]:
    """Build and validate a bridge error-bar payload."""
    spec: dict[str, object] = {
        "series_index": series_index,
        "bar_type": bar_type.strip(),
        "value_type": value_type.strip(),
    }
    for key in _PASS_THROUGH_FIELDS:
        value = options.pop(key, None)
        if value is not None:
            spec[key] = value
    if options:
        unknown = ", ".join(sorted(options))
        raise ValueError(f"unknown error bar option(s): {unknown}")
    validate_error_bar_spec(spec)
    return spec


def validate_error_bar_spec(spec: dict[str, object]) -> None:
    """Reject error-bar payloads PowerPoint would repair on open.

    Raises:
        ValueError: If a token is unknown, the series index is negative, or the
            custom/scalar value fields do not match the value type.
    """
    _require_token(spec.get("bar_type"), ERROR_BAR_TYPES, "bar_type")
    value_type = spec.get("value_type")
    _require_token(value_type, ERROR_BAR_VALUE_TYPES, "value_type")
    direction = spec.get("direction")
    if direction is not None:
        _require_token(direction, ERROR_BAR_DIRECTIONS, "direction")
    series_index = spec.get("series_index", 0)
    if isinstance(series_index, int) and series_index < 0:
        raise ValueError("series_index must not be negative")
    _validate_error_bar_values(cast("str", value_type), spec)


def _validate_error_bar_values(value_type: str, spec: dict[str, object]) -> None:
    has_reference = (
        spec.get("plus_reference") is not None
        or spec.get("minus_reference") is not None
    )
    value = spec.get("value")
    if value_type == _CUSTOM_VALUE_TYPE:
        if not has_reference:
            raise ValueError(
                "value_type 'cust' requires plus_reference and/or minus_reference"
            )
        if value is not None:
            raise ValueError("value is not valid with value_type 'cust'")
        return
    if has_reference:
        raise ValueError("plus_reference/minus_reference require value_type 'cust'")
    if value is not None and (not isinstance(value, int | float) or value < 0):
        raise ValueError("value must be a non-negative number")


def _require_token(value: object, allowed: tuple[str, ...], name: str) -> None:
    if not isinstance(value, str) or value not in allowed:
        joined = ", ".join(allowed)
        raise ValueError(f"{name} must be one of: {joined}")
