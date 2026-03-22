"""Shared internal utilities for text-related facades."""

from __future__ import annotations


def as_optional_int(value: object) -> int | None:
    """Coerce value to an optional integer, supporting integer-like floats."""
    if value is None:
        return None
    if isinstance(value, int):
        return value
    if isinstance(value, float) and value.is_integer():
        return int(value)
    return None


def as_optional_bool(value: object) -> bool | None:
    """Coerce value to an optional boolean."""
    if value is None:
        return None
    if isinstance(value, bool):
        return value
    return None


def as_optional_string(value: object) -> str | None:
    """Coerce value to an optional string."""
    if value is None:
        return None
    if isinstance(value, str):
        return value
    return None


def as_optional_float(value: object) -> float | None:
    """Coerce value to an optional float."""
    if value is None:
        return None
    if isinstance(value, (int, float)):
        return float(value)
    return None
