"""Table row/column collection proxies."""
# ruff: noqa: D102,D105,D107,SLF001
# pyright: reportMissingSuperCall=false, reportPrivateUsage=false, reportOptionalMemberAccess=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .. import ops
from .table_cells import Cell

if TYPE_CHECKING:
    from collections.abc import Iterator

    from .table import Table


class TableRow:
    """Row proxy with height accessor."""

    def __init__(self, table: Table, index: int) -> None:
        self._table = table
        self.index = index

    @property
    def height(self) -> int:
        self._table._ensure_cache()
        rows = cast("list[dict[str, object]]", self._table._cache.get("rows", []))
        if self.index >= len(rows):
            return 0
        value = rows[self.index].get("height", 0)
        return int(cast("int", value))

    @height.setter
    def height(self, value: int) -> None:
        self._table.prs.execute(
            ops.OP_SET_TABLE_ROW_HEIGHT,
            {
                "slide_index": self._table.slide_index,
                "shape_id": self._table.shape_id,
                "row": self.index,
                "height": int(value),
            },
        )
        self._table.invalidate_cache()

    @property
    def cells(self) -> list[Cell]:
        return [
            Cell(self._table, self.index, col) for col in range(self._table.col_count)
        ]


class TableRows:
    """Row collection proxy."""

    def __init__(self, table: Table) -> None:
        self._table = table

    def __len__(self) -> int:
        return self._table.row_count

    def __getitem__(self, index: int) -> TableRow:
        if index < 0:
            index += len(self)
        if index < 0 or index >= len(self):
            raise IndexError("row index out of range")
        return TableRow(self._table, index)

    def __iter__(self) -> Iterator[TableRow]:
        for i in range(len(self)):
            yield TableRow(self._table, i)


class TableColumn:
    """Column proxy with width accessor."""

    def __init__(self, table: Table, index: int) -> None:
        self._table = table
        self.index = index

    @property
    def width(self) -> int:
        self._table._ensure_cache()
        cols = cast("list[dict[str, object]]", self._table._cache.get("columns", []))
        if self.index >= len(cols):
            return 0
        value = cols[self.index].get("width", 0)
        return int(cast("int", value))

    @width.setter
    def width(self, value: int) -> None:
        self._table.prs.execute(
            ops.OP_SET_TABLE_COLUMN_WIDTH,
            {
                "slide_index": self._table.slide_index,
                "shape_id": self._table.shape_id,
                "col": self.index,
                "width": int(value),
            },
        )
        self._table.invalidate_cache()

    @property
    def cells(self) -> list[Cell]:
        return [
            Cell(self._table, row, self.index) for row in range(self._table.row_count)
        ]


class TableColumns:
    """Column collection proxy."""

    def __init__(self, table: Table) -> None:
        self._table = table

    def __len__(self) -> int:
        return self._table.col_count

    def __getitem__(self, index: int) -> TableColumn:
        if index < 0:
            index += len(self)
        if index < 0 or index >= len(self):
            raise IndexError("column index out of range")
        return TableColumn(self._table, index)

    def __iter__(self) -> Iterator[TableColumn]:
        for i in range(len(self)):
            yield TableColumn(self._table, i)
