"""Type-only bridge contract for slide proxy objects."""

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import Protocol

if TYPE_CHECKING:
    from ...schemas import TableCellInfo, TableInfo


class TableOperationsProtocol(Protocol):
    """Table management operations."""

    def add_table(
        self,
        slide: int | None = None,
        slide_index: int | None = None,
        rows: int | None = None,
        cols: int | None = None,
        **kwargs: object,
    ) -> int:
        """Protocol member."""
        ...

    def get_table(self, slide_index: int, shape_id: int) -> TableInfo:
        """Protocol member."""
        ...

    def set_table_flags(
        self, slide_index: int, shape_id: int, flags: dict[str, bool]
    ) -> None:
        """Protocol member."""
        ...

    def set_table_cell_text(
        self, slide_index: int, shape_id: int, row: int, col: int, text: str
    ) -> None:
        """Protocol member."""
        ...

    def get_table_cell(
        self, slide_index: int, shape_id: int, row: int, col: int
    ) -> TableCellInfo:
        """Protocol member."""
        ...

    def set_table_style(self, slide_index: int, shape_id: int, style_guid: str) -> None:
        """Protocol member."""
        ...

    def merge_table_cells(
        self, slide_index: int, shape_id: int, cell_range: tuple[int, int, int, int]
    ) -> None:
        """Protocol member."""
        ...

    def split_table_cell(
        self, slide_index: int, shape_id: int, row: int, col: int
    ) -> None:
        """Protocol member."""
        ...

    def set_table_row_height(
        self, slide_index: int, shape_id: int, row: int, height: int
    ) -> None:
        """Protocol member."""
        ...

    def set_table_column_width(
        self, slide_index: int, shape_id: int, col: int, width: int
    ) -> None:
        """Protocol member."""
        ...
