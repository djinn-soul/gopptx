"""Formatting properties shared by category and value chart axes."""

from __future__ import annotations

import math
from typing import TYPE_CHECKING, Protocol, cast

if TYPE_CHECKING:
    from ...schemas import ChartAxisState, ChartFormatUpdate, ChartState


class _AxisChartProtocol(Protocol):
    """Minimal chart surface used by axis formatting helpers."""

    def snapshot(self) -> ChartState:
        """Return the current chart state."""
        ...

    def apply_format(self, fmt: ChartFormatUpdate) -> None:
        """Apply a chart format update."""
        ...


_AXIS_NAMES = ("category", "value")


class ChartAxisFormatMixin:
    """Adds python-pptx-style axis title, scale, unit, and format controls."""

    def __init__(self, chart: _AxisChartProtocol, *, axis_name: str) -> None:
        """Initialize the shared axis formatting state."""
        super().__init__()
        # Validated here rather than per property: every axis lookup treats any
        # name other than "category" as the value axis, so a mis-specified name
        # would silently format the wrong axis.
        if axis_name not in _AXIS_NAMES:
            raise ValueError(f"axis_name must be one of: {', '.join(_AXIS_NAMES)}")
        self._chart = chart
        self._axis_name = axis_name

    def _payload(self) -> ChartAxisState:
        """Return the current axis state payload from the concrete proxy."""
        raise NotImplementedError

    @property
    def _format_prefix(self) -> str:
        return "category_axis_" if self._axis_name == "category" else "value_axis_"

    def _format_key(self, name: str) -> str:
        return self._format_prefix + name

    def _axis_state_value(self, key: str) -> object:
        payload = self._payload()
        return payload.get(key)

    def _apply_axis_format(self, name: str, value: object) -> None:
        self._chart.apply_format(
            cast("ChartFormatUpdate", {self._format_key(name): value})
        )

    @property
    def title(self) -> str | None:
        """Return the axis title text, if present."""
        value = self._axis_state_value("title")
        return str(value) if isinstance(value, str) and value else None

    @title.setter
    def title(self, value: str) -> None:
        self._apply_axis_format("title", value)

    @property
    def minimum_scale(self) -> float | None:
        """Return the explicit axis minimum, or ``None`` for automatic scaling."""
        return self._float_state("minimum_scale")

    @minimum_scale.setter
    def minimum_scale(self, value: float) -> None:
        self._set_scale("minimum_scale", value, "maximum_scale")

    @property
    def maximum_scale(self) -> float | None:
        """Return the explicit axis maximum, or ``None`` for automatic scaling."""
        return self._float_state("maximum_scale")

    @maximum_scale.setter
    def maximum_scale(self, value: float) -> None:
        self._set_scale("maximum_scale", value, "minimum_scale")

    def _set_scale(self, name: str, value: float, other_name: str) -> None:
        normalized = self._finite_number(value, name)
        payload: dict[str, object] = {self._format_key(name): normalized}
        other = self._float_state(other_name)
        if other is not None:
            payload[self._format_key(other_name)] = other
        self._chart.apply_format(cast("ChartFormatUpdate", payload))

    @property
    def major_unit(self) -> float | None:
        """Return the major tick interval, if explicitly configured."""
        return self._float_state("major_unit")

    @major_unit.setter
    def major_unit(self, value: float) -> None:
        self._set_unit("major_unit", value)

    @property
    def minor_unit(self) -> float | None:
        """Return the minor tick interval, if explicitly configured."""
        return self._float_state("minor_unit")

    @minor_unit.setter
    def minor_unit(self, value: float) -> None:
        self._set_unit("minor_unit", value)

    def _set_unit(self, name: str, value: float) -> None:
        normalized = self._finite_number(value, name)
        if normalized <= 0:
            raise ValueError(f"{name} must be greater than zero")
        self._apply_axis_format(name, normalized)

    @property
    def tick_label_number_format(self) -> str | None:
        """Return the OOXML number format for tick labels."""
        value = self._axis_state_value("number_format")
        return str(value) if isinstance(value, str) and value else None

    @tick_label_number_format.setter
    def tick_label_number_format(self, value: str) -> None:
        if not value.strip():
            raise ValueError("tick_label_number_format must not be empty")
        self._apply_axis_format("number_format", value)

    @property
    def tick_label_number_format_is_linked(self) -> bool | None:
        """Return whether the number format follows its linked source data."""
        value = self._axis_state_value("format_linked")
        return value if isinstance(value, bool) else None

    @tick_label_number_format_is_linked.setter
    def tick_label_number_format_is_linked(self, value: bool) -> None:
        self._apply_axis_format("format_linked", value)

    def _float_state(self, name: str) -> float | None:
        value = self._axis_state_value(name)
        # bool is a subclass of int, so it would otherwise become 1.0/0.0.
        if isinstance(value, bool) or not isinstance(value, int | float):
            return None
        return float(value)

    @staticmethod
    def _finite_number(value: float, name: str) -> float:
        normalized = value
        if not math.isfinite(normalized):
            raise ValueError(f"{name} must be finite")
        return normalized
