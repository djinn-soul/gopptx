"""Table cell proxy objects."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .. import ops
from ..api_errors import GopptxError

if TYPE_CHECKING:
    from .table import Table


class Cell:
    """Proxy object for a table cell."""

    def __init__(self, table: Table, row: int, col: int) -> None:
        """Initialize the cell proxy."""
        super().__init__()
        self._table = table
        self.row = row
        self.col = col

    @property
    def is_merge_origin(self) -> bool:
        """Check if this cell is the origin of a merged cell range."""
        return (
            self._table.get_cell_info(self.row, self.col).get("is_merge_origin", False)
            is True
        )

    @property
    def is_spanned(self) -> bool:
        """Check if this cell is spanned by another merged cell."""
        return (
            self._table.get_cell_info(self.row, self.col).get("is_spanned", False)
            is True
        )

    @property
    def row_span(self) -> int:
        """Get the number of rows this cell spans."""
        return int(
            cast(
                "int", self._table.get_cell_info(self.row, self.col).get("row_span", 1)
            )
        )

    @property
    def col_span(self) -> int:
        """Get the number of columns this cell spans."""
        return int(
            cast(
                "int", self._table.get_cell_info(self.row, self.col).get("col_span", 1)
            )
        )

    @property
    def text(self) -> str:
        """Get the text content of this cell."""
        return str(self._table.get_cell_info(self.row, self.col).get("text", ""))

    @text.setter
    def text(self, value: str) -> None:
        self._table.update_cell(self.row, self.col, {"text": str(value)})

    def split(self) -> None:
        """Split a merged cell back into a 1x1 cell."""
        if getattr(self._table.prs, "_batch_active", False):
            raise GopptxError(
                "structural changes (split) are not allowed inside a batch",
                code="BATCH_STRUCTURAL_CHANGE_NOT_ALLOWED",
            )
        self._table.prs.execute(
            ops.OP_SPLIT_TABLE_CELL,
            {
                "slide_index": self._table.slide_index,
                "shape_id": self._table.shape_id,
                "row": self.row,
                "col": self.col,
            },
        )
        self._table.invalidate_cache()

    def __repr__(self) -> str:  # pyright: ignore[reportImplicitOverride]
        """Return a string representation of the cell."""
        return f"<Cell [{self.row}, {self.col}] text={self.text!r}>"


class CellRange:
    """Represents a 2D slice of cells in a table for bulk operations."""

    def __init__(
        self, table: Table, row_start: int, row_end: int, col_start: int, col_end: int
    ) -> None:
        """Initialize the cell range with bounds."""
        super().__init__()
        self._table = table
        self.row_start = max(0, row_start)
        self.row_end = min(table.row_count, row_end)
        self.col_start = max(0, col_start)
        self.col_end = min(table.col_count, col_end)

    def merge(self) -> None:
        """Merge all cells in this range into a single spanned cell."""
        if getattr(self._table.prs, "_batch_active", False):
            raise GopptxError(
                "structural changes (merge) are not allowed inside a batch",
                code="BATCH_STRUCTURAL_CHANGE_NOT_ALLOWED",
            )
        if self.row_end <= self.row_start + 1 and self.col_end <= self.col_start + 1:
            return

        self._table.prs.execute(
            ops.OP_MERGE_TABLE_CELLS,
            {
                "slide_index": self._table.slide_index,
                "shape_id": self._table.shape_id,
                "row1": self.row_start,
                "col1": self.col_start,
                "row2": self.row_end - 1,
                "col2": self.col_end - 1,
            },
        )
        self._table.invalidate_cache()
