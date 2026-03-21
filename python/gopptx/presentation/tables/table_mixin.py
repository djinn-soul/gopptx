"""Presentation table mixin."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ..helpers import PresentationMixinBase
from .table_cell_mixin import PresentationTableCellMixin
from .table_style_mixin import PresentationTableStyleMixin

if TYPE_CHECKING:
    from ...schemas import TableInfo


class PresentationTableMixin(
    PresentationTableStyleMixin,
    PresentationTableCellMixin,
    PresentationMixinBase,
):
    """Mixin providing table creation and manipulation methods."""

    def add_table(  # noqa: C901, PLR0913
        self,
        slide: int | None = None,
        slide_index: int | None = None,
        rows: int | None = None,
        cols: int | None = None,
        bounds: tuple[int, int, int, int] | None = None,
        x: int | None = None,
        y: int | None = None,
        cx: int | None = None,
        cy: int | None = None,
        data: list[list[str]] | None = None,
        first_row: bool = False,
        first_col: bool = False,
        last_row: bool = False,
        last_col: bool = False,
        band_row: bool = False,
        band_col: bool = False,
        column_widths: list[int] | None = None,
    ) -> int:
        """Add a table shape to a slide and return its shape ID."""
        # Detect old positional calling pattern: add_table(slide_idx, rows, cols, bounds)
        if (
            slide is not None
            and isinstance(slide, int)
            and slide_index is not None
            and isinstance(slide_index, int)
            and rows is not None
            and isinstance(rows, int)
            and cols is not None
            and isinstance(cols, tuple)
            and len(cast("tuple", cols)) == 4
        ):
            slide_idx = slide
            rows_count = slide_index
            cols_count = rows
            bounds = cast("tuple", cols)
            rows = rows_count
            cols = cols_count
            if rows is None or cols is None:
                raise ValueError("'rows' and 'cols' parameters are required")
        else:
            slide_idx = slide if slide is not None else slide_index
            if slide_idx is None:
                raise ValueError("Either 'slide' or 'slide_index' must be provided")
            if rows is None:
                raise ValueError("'rows' parameter is required")
            if cols is None:
                raise ValueError("'cols' parameter is required")

        if bounds is not None:
            x, y, cx, cy = bounds
        elif x is None or y is None or cx is None or cy is None:
            raise ValueError("Either 'bounds' tuple or (x, y, cx, cy) must be provided")

        result = self.execute(
            ops.OP_ADD_TABLE,
            {
                "slide_index": slide_idx,
                "rows": rows,
                "cols": cols,
                "x": x,
                "y": y,
                "cx": cx,
                "cy": cy,
            },
        )
        shape_id = int(cast("int", result.get("shape_id", 0)))

        if data is not None:
            for row_idx, row in enumerate(data):
                for col_idx, text in enumerate(row):
                    self.set_table_cell_text(
                        slide_idx, shape_id, row_idx, col_idx, text
                    )

        if column_widths is not None:
            for col_idx, width in enumerate(column_widths):
                self.set_table_column_width(slide_idx, shape_id, col_idx, width)

        flags: dict[str, bool] = {
            "first_row": first_row,
            "first_col": first_col,
            "last_row": last_row,
            "last_col": last_col,
            "band_row": band_row,
            "band_col": band_col,
        }
        if any(flags.values()):
            self.set_table_flags(slide_idx, shape_id, flags)

        return shape_id

    def get_table(self, slide_index: int, shape_id: int) -> TableInfo:
        """Return serialized table information for a table shape."""
        result = self.execute(
            ops.OP_GET_TABLE,
            {"slide_index": slide_index, "shape_id": shape_id},
        )
        return cast("TableInfo", cast("dict[str, object]", result.get("table", {})))
