"""Shared structural protocols for table-domain proxies."""

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol

if TYPE_CHECKING:
    from ..contracts import SlidePresentationProtocol


class TableReadProto(Protocol):
    def table_state(self) -> dict[str, object]: ...

    @property
    def row_count(self) -> int: ...

    @property
    def col_count(self) -> int: ...


class TableWriteProto(TableReadProto, Protocol):
    prs: SlidePresentationProtocol
    slide_index: int
    shape_id: int

    def invalidate_cache(self) -> None: ...

    def get_cell_info(self, row: int, col: int) -> dict[str, object]: ...

    def update_cell(self, row: int, col: int, updates: dict[str, object]) -> None: ...
