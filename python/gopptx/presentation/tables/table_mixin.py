"""Presentation table mixin."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ...api_errors import GopptxError
from ...utils import normalize_table_index
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ...schemas import TableCellInfo, TableInfo


class PresentationTableMixin(PresentationMixinBase):
    """Mixin providing table creation and manipulation methods."""

    def add_table(
        self,
        slide_index: int,
        rows: int,
        cols: int,
        bounds: tuple[int, int, int, int],
    ) -> int:
        """Add a table shape to a slide and return its shape ID."""
        x, y, cx, cy = bounds
        result = self.execute(
            ops.OP_ADD_TABLE,
            {
                "slide_index": slide_index,
                "rows": rows,
                "cols": cols,
                "x": x,
                "y": y,
                "cx": cx,
                "cy": cy,
            },
        )
        return int(cast("int", result.get("shape_id", 0)))

    def get_table(self, slide_index: int, shape_id: int) -> TableInfo:
        """Return serialized table information for a table shape."""
        result = self.execute(
            ops.OP_GET_TABLE,
            {"slide_index": slide_index, "shape_id": shape_id},
        )
        return cast("TableInfo", cast("dict[str, object]", result.get("table", {})))

    def set_table_style(self, slide_index: int, shape_id: int, style_guid: str) -> None:
        """Apply a table style GUID to an existing table shape."""
        self.execute(
            ops.OP_SET_TABLE_STYLE,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "style_guid": style_guid,
            },
        )

    def define_table_style(self, name: str, style_id: str | None = None) -> str:
        """Define a custom table style and return its resolved style ID."""
        payload: dict[str, object] = {"name": name}
        if style_id is not None:
            payload["style_id"] = style_id
        result = self.execute(ops.OP_DEFINE_TABLE_STYLE, payload)
        return str(result.get("style_id", ""))

    def list_table_styles(self) -> list[dict[str, str]]:
        """List available table styles visible to the presentation."""
        result = self.execute(ops.OP_LIST_TABLE_STYLES, {})
        return cast("list[dict[str, str]]", result.get("styles", []))

    def set_table_flags(
        self,
        slide_index: int,
        shape_id: int,
        flags: dict[str, bool],
    ) -> None:
        """Set table display flags such as header-row or banded options."""
        self.execute(
            ops.OP_UPDATE_TABLE_FLAGS,
            {"slide_index": slide_index, "shape_id": shape_id, "flags": flags},
        )

    def set_table_cell_text(
        self,
        slide_index: int,
        shape_id: int,
        row: int,
        col: int,
        text: str,
    ) -> None:
        """Update the text value for one table cell."""
        self.execute(
            ops.OP_UPDATE_TABLE_CELL,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "row": row,
                "col": col,
                "updates": {"text": text},
            },
        )

    def get_table_cell(
        self,
        slide_index: int,
        shape_id: int,
        row: int,
        col: int,
    ) -> TableCellInfo:
        """Return one table cell payload by zero-based row and column."""
        table = self.get_table(slide_index, shape_id)
        cells = cast("list[dict[str, object]]", table.get("cells", []))
        cell_map: dict[tuple[int, int], dict[str, object]] = {}
        for cell in cells:
            try:
                row_idx = normalize_table_index(cell["row"])
                col_idx = normalize_table_index(cell["col"])
            except (KeyError, ValueError):
                continue
            cell_map[row_idx, col_idx] = cell
        found = cell_map.get((row, col))
        if found is not None:
            return cast("TableCellInfo", found)
        raise GopptxError(f"table cell [{row},{col}] not found", code="OP_FAILED")

    def merge_table_cells(
        self,
        slide_index: int,
        shape_id: int,
        cell_range: tuple[int, int, int, int],
    ) -> None:
        """Merge a rectangular range of table cells."""
        row1, col1, row2, col2 = cell_range
        self.execute(
            ops.OP_MERGE_TABLE_CELLS,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "row1": row1,
                "col1": col1,
                "row2": row2,
                "col2": col2,
            },
        )

    def split_table_cell(
        self,
        slide_index: int,
        shape_id: int,
        row: int,
        col: int,
    ) -> None:
        """Split a merged table cell back into its original cells."""
        self.execute(
            ops.OP_SPLIT_TABLE_CELL,
            {"slide_index": slide_index, "shape_id": shape_id, "row": row, "col": col},
        )

    def set_table_row_height(
        self,
        slide_index: int,
        shape_id: int,
        row: int,
        height: int,
    ) -> None:
        """Set the height of a specific table row."""
        self.execute(
            ops.OP_SET_TABLE_ROW_HEIGHT,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "row": row,
                "height": height,
            },
        )

    def set_table_column_width(
        self,
        slide_index: int,
        shape_id: int,
        col: int,
        width: int,
    ) -> None:
        """Set the width of a specific table column."""
        self.execute(
            ops.OP_SET_TABLE_COLUMN_WIDTH,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "col": col,
                "width": width,
            },
        )
