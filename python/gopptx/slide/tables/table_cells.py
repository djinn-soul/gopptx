"""Table cell proxy objects."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ...api_errors import GopptxError
from ..shapes.picture_fill import picture_fill_payload
from .table_cell_margins import CellMarginMixin

if TYPE_CHECKING:
    import os

    from ._protocols import TableWriteProto


class Cell(CellMarginMixin):
    """Proxy object for a table cell."""

    def __init__(self, table: TableWriteProto, row: int, col: int) -> None:
        """Initialize the cell proxy."""
        super().__init__()
        self._table = table
        self.row = row
        self.col = col

    @property
    def row_idx(self) -> int:
        """python-pptx alias for :attr:`row` (zero-based row index)."""
        return self.row

    @property
    def col_idx(self) -> int:
        """python-pptx alias for :attr:`col` (zero-based column index)."""
        return self.col

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
    def span_height(self) -> int:
        """python-pptx alias for row_span (number of rows spanned by this cell)."""
        return self.row_span

    @property
    def span_width(self) -> int:
        """python-pptx alias for col_span (number of columns spanned by this cell)."""
        return self.col_span

    @property
    def text(self) -> str:
        """Get the text content of this cell."""
        return str(self._table.get_cell_info(self.row, self.col).get("text", ""))

    @text.setter
    def text(self, value: str) -> None:
        self._table.update_cell(self.row, self.col, {"text": str(value)})

    @property
    def size_pt(self) -> float | None:
        """Get the font size in points for this cell, or None if not set."""
        val = self._table.get_cell_info(self.row, self.col).get("size_pt")
        return float(val) if val is not None else None  # type: ignore[arg-type]

    @size_pt.setter
    def size_pt(self, value: float) -> None:
        self._table.update_cell(self.row, self.col, {"size_pt": float(value)})

    @property
    def font_name(self) -> str | None:
        """Get the font family name for this cell, or None if not set."""
        val = self._table.get_cell_info(self.row, self.col).get("font_name")
        return str(val) if val is not None else None

    @font_name.setter
    def font_name(self, value: str) -> None:
        self._table.update_cell(self.row, self.col, {"font_name": str(value)})

    def _get_border(self, side: str) -> dict[str, object] | None:
        val = self._table.get_cell_info(self.row, self.col).get(f"border_{side}")
        return cast("dict[str, object]", val) if isinstance(val, dict) else None

    def _set_border(self, side: str, value: dict[str, object] | None) -> None:
        self._table.prs.execute(
            ops.OP_UPDATE_TABLE_CELL_BORDER,
            {
                "slide_index": self._table.slide_index,
                "shape_id": self._table.shape_id,
                "row": self.row,
                "col": self.col,
                "side": side,
                "border": value,
            },
        )
        self._table.invalidate_cache()

    @property
    def border_left(self) -> dict[str, object] | None:
        """Get the left border dict, or None if unset.

        Keys: width, color, dash, cap, join, compound, inset, miter_limit.
        """
        return self._get_border("left")

    @border_left.setter
    def border_left(self, value: dict[str, object] | None) -> None:
        self._set_border("left", value)

    @property
    def border_right(self) -> dict[str, object] | None:
        """Get the right border dict, or None if unset.

        Keys: width, color, dash, cap, join, compound, inset, miter_limit.
        """
        return self._get_border("right")

    @border_right.setter
    def border_right(self, value: dict[str, object] | None) -> None:
        self._set_border("right", value)

    @property
    def border_top(self) -> dict[str, object] | None:
        """Get the top border dict, or None if unset.

        Keys: width, color, dash, cap, join, compound, inset, miter_limit.
        """
        return self._get_border("top")

    @border_top.setter
    def border_top(self, value: dict[str, object] | None) -> None:
        self._set_border("top", value)

    @property
    def border_bottom(self) -> dict[str, object] | None:
        """Get the bottom border dict, or None if unset.

        Keys: width, color, dash, cap, join, compound, inset, miter_limit.
        """
        return self._get_border("bottom")

    @border_bottom.setter
    def border_bottom(self, value: dict[str, object] | None) -> None:
        self._set_border("bottom", value)

    def set_border(
        self,
        color: str | None = None,
        width: float | None = None,
        side: str = "all",
    ) -> None:
        """Set border color and width for this table cell (Issue #1084).

        Args:
            color: Hex RGB color string (e.g. "FF0000").
            width: Line width in EMU or points.
            side: "left", "right", "top", "bottom", or "all".
        """
        border_data: dict[str, object] = {}
        if color is not None:
            border_data["color"] = str(color).lstrip("#")
        if width is not None:
            border_data["width"] = int(width)

        sides = (
            ["left", "right", "top", "bottom"]
            if side.lower() == "all"
            else [side.lower()]
        )
        for s in sides:
            self._set_border(s, border_data)

    def set_picture_fill(
        self,
        image: str | bytes | os.PathLike[str],
        *,
        tile: bool = False,
        crop: tuple[float, float, float, float] | None = None,
    ) -> None:
        """Fill this cell's background with an image (Issue #234).

        Args:
            image: Image file path, ``os.PathLike``, or raw image bytes.
            tile: Repeat the image at its natural size instead of stretching a
                single copy over the cell.
            crop: ``(left, top, right, bottom)`` fractions in ``[0, 1)`` trimmed
                off the source image before it is drawn.

        Writes an ``<a:blipFill>`` into the cell's ``<a:tcPr>``, replacing any
        solid or gradient fill it had.
        """
        payload = picture_fill_payload(image, tile=tile, crop=crop)
        request: dict[str, object] = {
            "slide_index": self._table.slide_index,
            "shape_id": self._table.shape_id,
            "row": self.row,
            "column": self.col,
        }
        request.update(payload)
        self._table.prs.execute(ops.OP_SET_PICTURE_FILL, request)
        self._table.invalidate_cache()

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
        self,
        table: TableWriteProto,
        row_start: int,
        row_end: int,
        col_start: int,
        col_end: int,
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
