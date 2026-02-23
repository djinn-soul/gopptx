"""Utility functions for gopptx library."""

from __future__ import annotations


def _normalize_table_index(value: float) -> int:
    """Normalize a table index value to an integer."""
    if isinstance(value, bool):
        raise ValueError("table index must be an integer")
    if isinstance(value, int):
        return value
    if isinstance(value, float):
        if not value.is_integer():
            raise ValueError("table index must be integral")
        return int(value)
    raise ValueError("table index must be an integer")
