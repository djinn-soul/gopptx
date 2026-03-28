"""Table row/column collection proxies."""
# pyright: reportMissingSuperCall=false, reportPrivateUsage=false, reportOptionalMemberAccess=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from .table_cells import Cell

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ._protocols import TableWriteProto


class TableRow:
    """Row proxy with height accessor."""

    def __init__(self, table: TableWriteProto, index: int) -> None:
        """Initialize a row proxy for a table row index."""
        self._table = table
        self.index = index

    @property
    def height(self) -> int:
        """Return row height in EMUs."""
        rows = cast(
            "list[dict[str, object]]", self._table.table_state().get("rows", [])
        )
        if self.index >= len(rows):
            return 0
        return _as_int(rows[self.index].get("height"))

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
        """Return cells in this row."""
        return [
            Cell(self._table, self.index, col) for col in range(self._table.col_count)
        ]


class TableRows:
    """Row collection proxy."""

    def __init__(self, table: TableWriteProto) -> None:
        """Initialize a row collection for a table."""
        self._table = table

    def __len__(self) -> int:
        """Return total row count."""
        return self._table.row_count

    def __getitem__(self, index: int) -> TableRow:
        """Return row proxy by index, supporting negative indices."""
        if index < 0:
            index += len(self)
        if index < 0 or index >= len(self):
            raise IndexError("row index out of range")
        return TableRow(self._table, index)

    def __iter__(self) -> Iterator[TableRow]:
        """Iterate row proxies in order."""
        for i in range(len(self)):
            yield TableRow(self._table, i)

    def add(self, height: int = 0) -> TableRow:
        """Append a new empty row to the table and return a proxy for it.

        Args:
            height: Row height in EMUs. Pass 0 to let PowerPoint auto-size.
        """
        self._table.prs.execute(
            ops.OP_ADD_TABLE_ROW,
            {
                "slide_index": self._table.slide_index,
                "shape_id": self._table.shape_id,
                "height": int(height),
            },
        )
        self._table.invalidate_cache()
        return TableRow(self._table, len(self) - 1)

    def insert(self, at_index: int, height: int = 0) -> TableRow:
        """Insert a new empty row before at_index and return a proxy for it.

        Args:
            at_index: Position to insert before. Use len(rows) to append.
            height: Row height in EMUs. Pass 0 to let PowerPoint auto-size.
        """
        self._table.prs.execute(
            ops.OP_INSERT_TABLE_ROW,
            {
                "slide_index": self._table.slide_index,
                "shape_id": self._table.shape_id,
                "at_index": int(at_index),
                "height": int(height),
            },
        )
        self._table.invalidate_cache()
        return TableRow(self._table, at_index)

    def remove(self, at_index: int) -> None:
        """Remove the row at at_index.

        Args:
            at_index: Zero-based row index to remove.
        """
        self._table.prs.execute(
            ops.OP_REMOVE_TABLE_ROW,
            {
                "slide_index": self._table.slide_index,
                "shape_id": self._table.shape_id,
                "at_index": int(at_index),
            },
        )
        self._table.invalidate_cache()


class TableColumn:
    """Column proxy with width accessor."""

    def __init__(self, table: TableWriteProto, index: int) -> None:
        """Initialize a column proxy for a table column index."""
        self._table = table
        self.index = index

    @property
    def width(self) -> int:
        """Return column width in EMUs."""
        cols = cast(
            "list[dict[str, object]]", self._table.table_state().get("columns", [])
        )
        if self.index >= len(cols):
            return 0
        return _as_int(cols[self.index].get("width"))

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
        """Return cells in this column."""
        return [
            Cell(self._table, row, self.index) for row in range(self._table.row_count)
        ]


class TableColumns:
    """Column collection proxy."""

    def __init__(self, table: TableWriteProto) -> None:
        """Initialize a column collection for a table."""
        self._table = table

    def __len__(self) -> int:
        """Return total column count."""
        return self._table.col_count

    def __getitem__(self, index: int) -> TableColumn:
        """Return column proxy by index, supporting negative indices."""
        if index < 0:
            index += len(self)
        if index < 0 or index >= len(self):
            raise IndexError("column index out of range")
        return TableColumn(self._table, index)

    def __iter__(self) -> Iterator[TableColumn]:
        """Iterate column proxies in order."""
        for i in range(len(self)):
            yield TableColumn(self._table, i)

    def add(self, width: int) -> TableColumn:
        """Append a new empty column to the table and return a proxy for it.

        Args:
            width: Column width in EMUs (must be > 0).
        """
        self._table.prs.execute(
            ops.OP_ADD_TABLE_COLUMN,
            {
                "slide_index": self._table.slide_index,
                "shape_id": self._table.shape_id,
                "width": int(width),
            },
        )
        self._table.invalidate_cache()
        return TableColumn(self._table, len(self) - 1)

    def insert(self, at_index: int, width: int) -> TableColumn:
        """Insert a new empty column before at_index and return a proxy for it.

        Args:
            at_index: Position to insert before. Use len(columns) to append.
            width: Column width in EMUs (must be > 0).
        """
        self._table.prs.execute(
            ops.OP_INSERT_TABLE_COLUMN,
            {
                "slide_index": self._table.slide_index,
                "shape_id": self._table.shape_id,
                "at_index": int(at_index),
                "width": int(width),
            },
        )
        self._table.invalidate_cache()
        return TableColumn(self._table, at_index)

    def remove(self, at_index: int) -> None:
        """Remove the column at at_index.

        Args:
            at_index: Zero-based column index to remove.
        """
        self._table.prs.execute(
            ops.OP_REMOVE_TABLE_COLUMN,
            {
                "slide_index": self._table.slide_index,
                "shape_id": self._table.shape_id,
                "at_index": int(at_index),
            },
        )
        self._table.invalidate_cache()


def _as_int(value: object) -> int:
    if isinstance(value, int):
        return value
    return 0
