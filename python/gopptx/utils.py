"""Utility functions for gopptx library."""

from __future__ import annotations

import math
from typing import TYPE_CHECKING, cast

if TYPE_CHECKING:
    from typing_extensions import TypeGuard

_FOUR_BOUNDS_COMPONENTS = 4


def normalize_table_index(value: object) -> int:
    """Normalize a table index value to an integer."""
    if isinstance(value, float) and value.is_integer():
        return int(value)
    if isinstance(value, int) and not isinstance(value, bool):
        return value
    raise ValueError("table index must be an integer")


def is_four_number_bounds(
    value: object,
) -> TypeGuard[tuple[float, float, float, float]]:
    """Check if value is a tuple of four finite numbers (bounds)."""
    if not isinstance(value, tuple):
        return False
    components = cast("tuple[object, ...]", value)
    if len(components) != _FOUR_BOUNDS_COMPONENTS:
        return False
    for component in components:
        if isinstance(component, bool) or not isinstance(component, int | float):
            return False
        if not math.isfinite(float(component)):
            return False
    return True
