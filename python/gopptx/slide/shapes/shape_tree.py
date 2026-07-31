"""Shared traversal helpers for nested slide shape records."""

from __future__ import annotations

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ...schemas import Shape


def iter_shape_records(records: list[Shape]) -> Iterator[Shape]:
    """Yield every record in depth-first document order."""
    for record in records:
        yield record
        yield from iter_shape_records(record.get("Shapes", []))
