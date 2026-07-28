"""Per-label formatting proxies for individual data labels."""

from __future__ import annotations

import string
from typing import TYPE_CHECKING, Protocol, cast

if TYPE_CHECKING:
    from collections.abc import Iterator, Sequence

    from ...schemas import ChartDataLabelPointSpec, ChartFormatUpdate, ChartState

_FONT_SIZE_MIN_PT = 1
_FONT_SIZE_MAX_PT = 400
_HEX_COLOR_LENGTH = 6

# Optional fields copied straight through to the bridge payload.
_PASS_THROUGH_FIELDS = (
    "number_format",
    "format_linked",
    "font_color",
    "font_size_pt",
    "font_bold",
    "show_value",
    "show_category",
    "show_series_name",
    "show_percent",
    "show_legend_key",
    "delete",
)


class _DataLabelPointChartProtocol(Protocol):
    """Minimal chart surface used by per-label helpers."""

    def snapshot(self) -> ChartState:
        """Return the current chart state."""
        ...

    def apply_format(self, fmt: ChartFormatUpdate) -> None:
        """Apply a chart format update."""
        ...


class DataLabelPoint:
    """Read-only proxy for one persisted ``c:dLbl``."""

    def __init__(self, payload: dict[str, object]) -> None:
        """Initialize with the raw label payload."""
        super().__init__()
        self._payload = payload

    @property
    def series_index(self) -> int:
        """Return the zero-based series this label belongs to."""
        return self._int_field("series_index") or 0

    @property
    def point_index(self) -> int:
        """Return the zero-based index of the labelled point."""
        return self._int_field("point_index") or 0

    @property
    def number_format(self) -> str | None:
        """Return this label's own number format code."""
        return self._string_field("number_format")

    @property
    def format_linked(self) -> bool | None:
        """Return whether the label's number format follows the source data."""
        return self._bool_field("format_linked")

    @property
    def font_color(self) -> str | None:
        """Return the label font colour as a hex RGB string."""
        return self._string_field("font_color")

    @property
    def font_size_pt(self) -> int | None:
        """Return the label font size in points."""
        return self._int_field("font_size_pt")

    @property
    def font_bold(self) -> bool | None:
        """Return whether the label font is bold."""
        return self._bool_field("font_bold")

    @property
    def show_value(self) -> bool | None:
        """Return whether the label shows the point value."""
        return self._bool_field("show_value")

    @property
    def show_category(self) -> bool | None:
        """Return whether the label shows the category name."""
        return self._bool_field("show_category")

    @property
    def show_series_name(self) -> bool | None:
        """Return whether the label shows the series name."""
        return self._bool_field("show_series_name")

    @property
    def show_percent(self) -> bool | None:
        """Return whether the label shows the percentage."""
        return self._bool_field("show_percent")

    @property
    def deleted(self) -> bool:
        """Return whether this point's label is removed."""
        return self._bool_field("delete") or False

    def to_payload(self) -> dict[str, object]:
        """Return a copy of the raw payload, ready to re-send unchanged."""
        return dict(self._payload)

    def _string_field(self, key: str) -> str | None:
        value = self._payload.get(key)
        return str(value) if isinstance(value, str) and value else None

    def _bool_field(self, key: str) -> bool | None:
        value = self._payload.get(key)
        return value if isinstance(value, bool) else None

    def _int_field(self, key: str) -> int | None:
        value = self._payload.get(key)
        # bool is a subclass of int, so it would otherwise become 1/0.
        if isinstance(value, bool) or not isinstance(value, int):
            return None
        return value


class DataLabelPointCollection:
    """Sequence-like container of per-label proxies."""

    def __init__(self, payload: list[dict[str, object]]) -> None:
        """Initialize the collection with raw label payloads."""
        super().__init__()
        self._payload = payload

    def __len__(self) -> int:
        """Return the number of individually formatted labels."""
        return len(self._payload)

    def __getitem__(self, index: int) -> DataLabelPoint:
        """Return the formatted label at the given index."""
        if index < 0:
            index += len(self._payload)
        if index < 0 or index >= len(self._payload):
            raise IndexError("data label index out of range")
        return DataLabelPoint(self._payload[index])

    def __iter__(self) -> Iterator[DataLabelPoint]:
        """Iterate over per-label proxies."""
        for item in self._payload:
            yield DataLabelPoint(item)

    def for_series(self, series_index: int) -> list[DataLabelPoint]:
        """Return the formatted labels belonging to one series."""
        return [label for label in self if label.series_index == series_index]

    def for_point(
        self, point_index: int, *, series_index: int = 0
    ) -> DataLabelPoint | None:
        """Return the label of one point, or None when it has no own formatting."""
        for label in self:
            if label.series_index == series_index and label.point_index == point_index:
                return label
        return None


class ChartDataLabelPointMixin:
    """Adds per-label reads and writes to the live chart proxy."""

    @property
    def data_label_points(
        self: _DataLabelPointChartProtocol,
    ) -> DataLabelPointCollection:
        """Return every label that carries its own format, font or flags."""
        # The snapshot is a bridge payload, so the declared type is a promise
        # the runtime value has to be checked against.
        raw: object = self.snapshot().get("data_label_points", [])
        if not isinstance(raw, list):
            return DataLabelPointCollection([])
        items = cast("list[object]", raw)
        return DataLabelPointCollection([
            dict(cast("dict[str, object]", item))
            for item in items
            if isinstance(item, dict)
        ])

    def format_data_label(
        self: _DataLabelPointChartProtocol,
        point_index: int,
        *,
        series_index: int = 0,
        **options: object,
    ) -> None:
        """Format the label of one data point, keeping its other properties.

        A ``c:dLbl`` inherits none of the surrounding display flags, so the ones
        the label, series or plot already had are carried over unless a
        ``show_*`` option overrides them.

        Args:
            point_index: Zero-based index of the labelled point.
            series_index: Zero-based index of the owning series.
            **options: Any of ``number_format``, ``format_linked``,
                ``font_color``, ``font_size_pt``, ``font_bold``, ``show_value``,
                ``show_category``, ``show_series_name``, ``show_percent``,
                ``show_legend_key``, and ``delete``.

        Raises:
            ValueError: If an index or value is out of range.
        """
        spec = build_data_label_point_spec(
            point_index, series_index=series_index, **options
        )
        self.apply_format(cast("ChartFormatUpdate", {"data_label_points": [spec]}))

    def format_data_labels(
        self: _DataLabelPointChartProtocol,
        labels: Sequence[ChartDataLabelPointSpec],
        *,
        series_index: int = 0,
    ) -> None:
        """Format several individual labels in one round-trip."""
        specs = [dict(label) for label in labels]
        for spec in specs:
            spec.setdefault("series_index", series_index)
            validate_data_label_point_spec(spec)
        if not specs:
            return
        self.apply_format(cast("ChartFormatUpdate", {"data_label_points": specs}))

    def set_data_label_number_format(
        self: _DataLabelPointChartProtocol,
        point_index: int,
        number_format: str,
        *,
        series_index: int = 0,
    ) -> None:
        """Give one point's label a number format of its own."""
        spec = build_data_label_point_spec(
            point_index, series_index=series_index, number_format=number_format
        )
        self.apply_format(cast("ChartFormatUpdate", {"data_label_points": [spec]}))


def build_data_label_point_spec(
    point_index: int, *, series_index: int = 0, **options: object
) -> dict[str, object]:
    """Build and validate a bridge per-label payload."""
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
        raise ValueError(f"unknown data label option(s): {unknown}")
    validate_data_label_point_spec(spec)
    return spec


def validate_data_label_point_spec(spec: dict[str, object]) -> None:
    """Reject indexes, colours and sizes the Go side would refuse.

    Raises:
        ValueError: If an index, colour or font size is out of range.
    """
    for key in ("series_index", "point_index"):
        value = spec.get(key, 0)
        if not isinstance(value, int) or isinstance(value, bool) or value < 0:
            raise ValueError(f"data label {key} must be a non-negative integer")

    color = spec.get("font_color")
    if color is not None:
        if not isinstance(color, str):
            raise ValueError("data label font_color must be a hex RGB string")
        cleaned = color.lstrip("#")
        if len(cleaned) != _HEX_COLOR_LENGTH or any(
            character not in string.hexdigits for character in cleaned
        ):
            raise ValueError(
                f"data label font_color must be a hex RGB string, got {color!r}"
            )

    size = spec.get("font_size_pt")
    if size is not None:
        if not isinstance(size, int) or isinstance(size, bool):
            raise ValueError("data label font_size_pt must be an integer")
        if size < _FONT_SIZE_MIN_PT or size > _FONT_SIZE_MAX_PT:
            bounds = f"{_FONT_SIZE_MIN_PT} and {_FONT_SIZE_MAX_PT}"
            raise ValueError(
                f"data label font_size_pt must be between {bounds}, got {size}"
            )
