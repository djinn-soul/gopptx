"""Utility functions for gopptx library."""

from __future__ import annotations


def normalize_table_index(value: object) -> int:
    """Normalize a table index value to an integer."""
    if isinstance(value, float) and value.is_integer():
        return int(value)
    if isinstance(value, int) and not isinstance(value, bool):
        return value
    raise ValueError("table index must be an integer")
