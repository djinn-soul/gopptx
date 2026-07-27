"""Guards on chart axis name validation and numeric state coercion."""

from __future__ import annotations

import math
from typing import TYPE_CHECKING, cast

import pytest
from gopptx.slide.chart.axis_format import ChartAxisFormatMixin

if TYPE_CHECKING:
    from gopptx.schemas import ChartAxisState


class _StubAxis(ChartAxisFormatMixin):
    """Minimal axis exposing a fixed state payload."""

    def __init__(self, axis_name: str, payload: dict[str, object]) -> None:
        """Store the canned payload the mixin reads back."""
        super().__init__(chart=cast("object", None), axis_name=axis_name)  # pyright: ignore[reportArgumentType]
        self._payload_data = payload

    def _payload(self) -> ChartAxisState:
        return cast("ChartAxisState", self._payload_data)


def test_unknown_axis_name_is_rejected() -> None:
    """A mis-specified name would silently format the value axis instead."""
    with pytest.raises(ValueError, match="axis_name must be one of"):
        _StubAxis("diagonal", {})


@pytest.mark.parametrize("axis_name", ["category", "value"])
def test_known_axis_names_are_accepted(axis_name: str) -> None:
    """Both supported axis names construct and read their payload."""
    axis = _StubAxis(axis_name, {"maximum_scale": 12.5})
    value = axis.maximum_scale
    assert value is not None
    assert math.isclose(value, 12.5)


def test_bool_state_is_not_coerced_to_float() -> None:
    """bool subclasses int, so True would otherwise read back as 1.0."""
    axis = _StubAxis("value", {"maximum_scale": True})
    assert axis.maximum_scale is None


def test_int_state_is_returned_as_float() -> None:
    """An int state value is widened to float."""
    axis = _StubAxis("value", {"maximum_scale": 200})
    value = axis.maximum_scale
    assert value is not None
    assert math.isclose(value, 200.0)


def test_unit_state_is_read_through_the_same_guard() -> None:
    """The bool guard applies to every numeric axis property."""
    axis = _StubAxis("value", {"major_unit": False, "minor_unit": 5})
    assert axis.major_unit is None
    minor = axis.minor_unit
    assert minor is not None
    assert math.isclose(minor, 5.0)
