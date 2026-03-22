"""Slide table mixin scoped to table-domain operations."""

from __future__ import annotations

from collections.abc import Callable
from typing import TYPE_CHECKING

from .table import Table

if TYPE_CHECKING:
    from ...schemas import TableCellInfo, TableInfo
    from ..contracts import SlidePresentationProtocol


class SlideTableMixin:
    """Mixin providing table manipulation methods for Slide objects."""

    if TYPE_CHECKING:
        _presentation: SlidePresentationProtocol  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int:
            """Slide index."""
            ...

    def add_table(
        self,
        rows: int,
        cols: int,
        bounds: tuple[int, int, int, int],
        data: list[list[str]] | None = None,
        **kwargs: object,
    ) -> int:
        """Add a table and invalidate shape/text caches when available.

        Args:
            rows: Number of rows.
            cols: Number of columns.
            bounds: (x, y, cx, cy) table position and size in EMU.
            data: Optional 2D array of cell text.
            **kwargs: Additional options (first_row, band_row, column_widths, etc).

        Returns:
            Shape ID of the created table.
        """
        slide_index = int(self.index)
        add_table: Callable[..., int] = self._presentation.add_table
        shape_id = add_table(
            slide=slide_index, rows=rows, cols=cols, bounds=bounds, data=data, **kwargs
        )
        invalidate = getattr(self, "_invalidate_shape_cache", None)
        if callable(invalidate):
            invalidate()
        invalidate_text = getattr(self, "_invalidate_text_state_cache", None)
        if callable(invalidate_text):
            invalidate_text()
        return shape_id

    def get_table(self, shape_id: int) -> TableInfo:
        """Fetch raw table payload by shape id."""
        return self._presentation.get_table(self.index, shape_id)

    def table(self, shape_id: int) -> Table:
        """Return table proxy object by shape id."""
        return Table(self._presentation, self.index, shape_id)

    def set_table_flags(self, shape_id: int, flags: dict[str, bool]) -> None:
        """Set table style flags (first_row, banded_rows, etc.)."""
        self._presentation.set_table_flags(self.index, shape_id, flags)

    def set_table_cell_text(self, shape_id: int, row: int, col: int, text: str) -> None:
        """Set text for a specific table cell."""
        self._presentation.set_table_cell_text(self.index, shape_id, row, col, text)

    def get_table_cell(self, shape_id: int, row: int, col: int) -> TableCellInfo:
        """Get payload for a specific table cell."""
        return self._presentation.get_table_cell(self.index, shape_id, row, col)

    def merge_table_cells(
        self, shape_id: int, cell_range: tuple[int, int, int, int]
    ) -> None:
        """Merge a rectangular table-cell range."""
        self._presentation.merge_table_cells(self.index, shape_id, cell_range)

    def split_table_cell(self, shape_id: int, row: int, col: int) -> None:
        """Split a previously merged table cell."""
        self._presentation.split_table_cell(self.index, shape_id, row, col)

    def set_table_style(self, shape_id: int, style_guid: str) -> None:
        """Apply table style by GUID."""
        self._presentation.set_table_style(self.index, shape_id, style_guid)

    def set_table_row_height(self, shape_id: int, row: int, height: int) -> None:
        """Set row height in EMUs."""
        self._presentation.set_table_row_height(self.index, shape_id, row, height)

    def set_table_column_width(self, shape_id: int, col: int, width: int) -> None:
        """Set column width in EMUs."""
        self._presentation.set_table_column_width(self.index, shape_id, col, width)
