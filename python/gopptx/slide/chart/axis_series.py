"""Axis and series helpers for live chart proxies."""
# pyright: reportMissingSuperCall=false, reportPrivateUsage=false

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

from typing_extensions import override

from .axis_format import ChartAxisFormatMixin

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ...schemas import ChartAxisState, ChartFormatUpdate, ChartState

_VALID_TICK_LABEL_POSITIONS = {
    "high",
    "low",
    "nextTo",
    "none",
}

_VALID_CROSSES = {
    "autoZero",
    "max",
    "min",
}

_VALID_LABEL_ALIGNMENTS = {
    "ctr",
    "l",
    "r",
}


class ChartProtocol(Protocol):
    """Protocol for chart operations required by axis helpers."""

    def snapshot(self) -> ChartState:
        """Return a chart state snapshot."""
        ...

    def apply_format(self, fmt: ChartFormatUpdate) -> None:
        """Apply a formatting update to the chart."""
        ...


class ChartAxis(ChartAxisFormatMixin):
    """Proxy exposing chart axis state and formatting operations."""

    def __init__(self, chart: ChartProtocol, *, axis_name: str) -> None:
        """Initialize an axis proxy for the given axis name."""
        super().__init__(chart, axis_name=axis_name)

    @staticmethod
    def _major_gridline_state_key() -> str:
        return "has" + "_major_gridline"

    @staticmethod
    def _major_gridline_format_key(axis_name: str) -> str:
        prefix = "category_axis_" if axis_name == "category" else "value_axis_"
        return prefix + "major_gridlines"

    @staticmethod
    def _minor_gridline_format_key(axis_name: str) -> str:
        prefix = "category_axis_" if axis_name == "category" else "value_axis_"
        return prefix + "minor_gridlines"

    @override
    def _payload(self) -> ChartAxisState:
        snapshot = self._chart.snapshot()
        key = "category_axis" if self._axis_name == "category" else "value_axis"
        raw: object = snapshot.get(key, {})
        return cast("ChartAxisState", raw if isinstance(raw, dict) else {})

    @property
    def axis_kind(self) -> str:
        """Return the normalized axis kind."""
        return "category" if self._axis_name == "category" else "value"

    @property
    def is_category_axis(self) -> bool:
        """Return whether this proxy targets the category axis."""
        return self._axis_name == "category"

    @property
    def is_value_axis(self) -> bool:
        """Return whether this proxy targets the value axis."""
        return self._axis_name == "value"

    @property
    def present(self) -> bool:
        """Return whether this axis is present in chart state."""
        return bool(self._payload().get("present", False))

    @property
    def tick_label_position(self) -> str | None:
        """Return the current tick-label position token."""
        value = self._payload().get("tick_label_pos")
        return str(value) if isinstance(value, str) else None

    @tick_label_position.setter
    def tick_label_position(self, value: str) -> None:
        normalized = self._normalize_tick_label_position(value)
        key = (
            "category_axis_tick_label_pos"
            if self._axis_name == "category"
            else "value_axis_tick_label_pos"
        )
        self._chart.apply_format(cast("ChartFormatUpdate", {key: normalized}))

    @property
    def major_gridlines_visible(self) -> bool:
        """Return whether major gridlines are enabled."""
        payload = self._payload()
        state_key = self._major_gridline_state_key()
        return bool(payload.get("major_gridline", payload.get(state_key, False)))

    @major_gridlines_visible.setter
    def major_gridlines_visible(self, value: bool) -> None:
        key = self._major_gridline_format_key(self._axis_name)
        self._chart.apply_format(cast("ChartFormatUpdate", {key: bool(value)}))

    @property
    def minor_gridlines_visible(self) -> bool:
        """Return whether minor gridlines are enabled."""
        payload = self._payload()
        return bool(payload.get("minor_gridline", False))

    @minor_gridlines_visible.setter
    def minor_gridlines_visible(self, value: bool) -> None:
        key = self._minor_gridline_format_key(self._axis_name)
        self._chart.apply_format(cast("ChartFormatUpdate", {key: bool(value)}))

    @property
    def crosses(self) -> str | None:
        """Return the current axis crossing mode token."""
        value = self._payload().get("crosses")
        return str(value) if isinstance(value, str) else None

    @crosses.setter
    def crosses(self, value: str) -> None:
        normalized = self._normalize_crosses(value)
        key = (
            "category_axis_crosses"
            if self._axis_name == "category"
            else "value_axis_crosses"
        )
        self._chart.apply_format(cast("ChartFormatUpdate", {key: normalized}))

    @property
    def crosses_auto_zero(self) -> bool:
        """Return whether axis crosses at zero."""
        return self.crosses == "autoZero"

    @property
    def crosses_at_maximum(self) -> bool:
        """Return whether axis crosses at maximum."""
        return self.crosses == "max"

    @property
    def crosses_at_minimum(self) -> bool:
        """Return whether axis crosses at minimum."""
        return self.crosses == "min"

    def set_crosses_auto_zero(self) -> None:
        """Set axis crossing mode to `autoZero`."""
        self.crosses = "autoZero"

    def set_crosses_at_maximum(self) -> None:
        """Set axis crossing mode to `max`."""
        self.crosses = "max"

    def set_crosses_at_minimum(self) -> None:
        """Set axis crossing mode to `min`."""
        self.crosses = "min"

    @property
    def tick_mark_skip(self) -> int | None:
        """Return how many categories are skipped between tick marks."""
        value = self._payload().get("tick_mark_skip")
        # bool is a subclass of int, so it would otherwise become 1/0.
        if isinstance(value, bool) or not isinstance(value, int):
            return None
        return value

    @tick_mark_skip.setter
    def tick_mark_skip(self, value: int) -> None:
        self._require_category_axis("tick_mark_skip")
        if value < 1:
            raise ValueError("tick_mark_skip must be greater than or equal to 1")
        self._chart.apply_format(
            cast("ChartFormatUpdate", {"category_axis_tick_mark_skip": int(value)})
        )

    @property
    def label_alignment(self) -> str | None:
        """Return the tick-label alignment token (``ctr``, ``l``, or ``r``)."""
        value = self._payload().get("label_alignment")
        return str(value) if isinstance(value, str) and value else None

    @label_alignment.setter
    def label_alignment(self, value: str) -> None:
        self._require_category_axis("label_alignment")
        normalized = value.strip()
        if normalized not in _VALID_LABEL_ALIGNMENTS:
            raise ValueError("label_alignment must be one of: ctr, l, r")
        self._chart.apply_format(
            cast("ChartFormatUpdate", {"category_axis_label_alignment": normalized})
        )

    def _require_category_axis(self, name: str) -> None:
        # c:tickMarkSkip and c:lblAlgn are CT_CatAx-only children; writing them
        # onto a value axis makes PowerPoint repair the file.
        if not self.is_category_axis:
            raise ValueError(f"{name} is only supported on the category axis")

    @staticmethod
    def _normalize_tick_label_position(value: str) -> str:
        normalized = value.strip()
        if normalized not in _VALID_TICK_LABEL_POSITIONS:
            raise ValueError(
                "tick_label_position must be one of: high, low, nextTo, none"
            )
        return normalized

    @staticmethod
    def _normalize_crosses(value: str) -> str:
        normalized = value.strip()
        if normalized not in _VALID_CROSSES:
            raise ValueError("crosses must be one of: autoZero, max, min")
        return normalized


class ChartSeries:
    """Read-only proxy for a chart series payload."""

    def __init__(self, payload: dict[str, object]) -> None:
        """Initialize with raw series payload."""
        self._payload = payload

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


class ChartSeriesCollection:
    """Sequence-like container of chart-series proxies."""

    def __init__(self, payload: list[dict[str, object]]) -> None:
        """Initialize the collection with raw series payloads."""
        self._payload = payload

    def __len__(self) -> int:
        """Return the number of series."""
        return len(self._payload)

    def __getitem__(self, index: int) -> ChartSeries:
        """Return the series at the given index."""
        if index < 0:
            index += len(self._payload)
        if index < 0 or index >= len(self._payload):
            raise IndexError("series index out of range")
        return ChartSeries(self._payload[index])

    def __iter__(self) -> Iterator[ChartSeries]:
        """Iterate over series proxies."""
        for item in self._payload:
            yield ChartSeries(item)
