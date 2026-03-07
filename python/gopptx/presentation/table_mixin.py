"""Presentation table mixin."""
# ruff: noqa: D102

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .. import ops
from ..api_errors import GopptxError
from ..utils import normalize_table_index
from .helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ..schemas import TableCellInfo, TableInfo


class PresentationTableMixin(PresentationMixinBase):
    """Mixin providing table creation and manipulation methods."""

    def add_table(
        self,
        slide_index: int,
        rows: int,
        cols: int,
        bounds: tuple[int, int, int, int],
    ) -> int:
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
        result = self.execute(
            ops.OP_GET_TABLE,
            {"slide_index": slide_index, "shape_id": shape_id},
        )
        return cast("TableInfo", cast("dict[str, object]", result.get("table", {})))

    def set_table_style(self, slide_index: int, shape_id: int, style_guid: str) -> None:
        self.execute(
            ops.OP_SET_TABLE_STYLE,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "style_guid": style_guid,
            },
        )

    def set_table_flags(
        self,
        slide_index: int,
        shape_id: int,
        flags: dict[str, bool],
    ) -> None:
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
        self.execute(
            ops.OP_SET_TABLE_COLUMN_WIDTH,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "col": col,
                "width": width,
            },
        )
