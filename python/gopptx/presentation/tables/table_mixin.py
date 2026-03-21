"""Presentation table mixin."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ..helpers import PresentationMixinBase
from .table_cell_mixin import PresentationTableCellMixin
from .table_style_mixin import PresentationTableStyleMixin

_BOUNDS_TUPLE_LEN = 4

if TYPE_CHECKING:
    from ...schemas import TableInfo


def _is_legacy_positional_call(
    slide: object,
    slide_index: object,
    rows: object,
    cols: object,
) -> bool:
    """Return True when called with the old positional (slide_idx, rows, cols, bounds) pattern."""
    return (
        slide is not None
        and isinstance(slide, int)
        and slide_index is not None
        and isinstance(slide_index, int)
        and rows is not None
        and isinstance(rows, int)
        and cols is not None
        and isinstance(cols, tuple)
        and len(cast("tuple", cols)) == _BOUNDS_TUPLE_LEN
    )


def _resolve_table_identity(
    slide: int | None,
    slide_index: int | None,
    rows: int | None,
    cols: int | None,
) -> tuple[int, int, int, tuple[int, int, int, int] | None]:
    """Resolve slide index, row/col counts, and optional bounds from raw args.

    Returns (slide_idx, rows, cols, bounds_or_None).
    """
    if _is_legacy_positional_call(slide, slide_index, rows, cols):
        bounds = cast("tuple[int, int, int, int]", cols)
        rows_val = slide_index
        cols_val = rows
        slide_idx = slide
    else:
        slide_idx = slide if slide is not None else slide_index
        if slide_idx is None:
            raise ValueError("Either 'slide' or 'slide_index' must be provided")
        rows_val = rows
        cols_val = cols
        bounds = None
    if rows_val is None:
        raise ValueError("'rows' parameter is required")
    if cols_val is None:
        raise ValueError("'cols' parameter is required")
    return cast("int", slide_idx), rows_val, cols_val, bounds


def _resolve_bounds(
    bounds: tuple[int, int, int, int] | None,
    kwargs: dict[str, object],
) -> tuple[int, int, int, int]:
    """Return (x, y, cx, cy) from bounds or individual x/y/cx/cy kwargs."""
    if bounds is not None:
        return bounds
    x = kwargs.get("x")
    y = kwargs.get("y")
    cx = kwargs.get("cx")
    cy = kwargs.get("cy")
    if x is None or y is None or cx is None or cy is None:
        raise ValueError("Either 'bounds' tuple or (x, y, cx, cy) must be provided")
    return cast("int", x), cast("int", y), cast("int", cx), cast("int", cy)


def _populate_table(
    mixin: object,
    slide_idx: int,
    shape_id: int,
    content: tuple[list[list[str]] | None, list[int] | None],
    flags: dict[str, bool],
) -> None:
    """Populate table data, column widths, and style flags.

    Note: Unlike the previous implementation using `any(flags.values())`,
    this now sends all flag values to the bridge, including explicit False.
    This allows users to disable features like header_row or banded_rows
    by explicitly passing first_row=False or band_row=False.
    """
    data, column_widths = content
    if data is not None:
        for row_idx, row in enumerate(data):
            for col_idx, text in enumerate(row):
                mixin.set_table_cell_text(slide_idx, shape_id, row_idx, col_idx, text)  # type: ignore[union-attr]
    if column_widths is not None:
        for col_idx, width in enumerate(column_widths):
            mixin.set_table_column_width(slide_idx, shape_id, col_idx, width)  # type: ignore[union-attr]
    if flags:
        mixin.set_table_flags(slide_idx, shape_id, flags)  # type: ignore[union-attr]


class PresentationTableMixin(
    PresentationTableStyleMixin,
    PresentationTableCellMixin,
    PresentationMixinBase,
):
    """Mixin providing table creation and manipulation methods."""

    def add_table(
        self,
        slide: int | None = None,
        slide_index: int | None = None,
        rows: int | None = None,
        cols: int | None = None,
        **kwargs: object,
    ) -> int:
        """Add a table shape to a slide and return its shape ID.

        Supports both the legacy positional API and the new named-parameter API.
        Keyword options: bounds, x, y, cx, cy, data, first_row, first_col,
        last_row, last_col, band_row, band_col, column_widths.
        """
        slide_idx, rows_val, cols_val, resolved_bounds = _resolve_table_identity(
            slide, slide_index, rows, cols
        )
        x, y, cx, cy = _resolve_bounds(
            resolved_bounds
            or cast("tuple[int, int, int, int] | None", kwargs.get("bounds")),
            kwargs,
        )

        result = self.execute(
            ops.OP_ADD_TABLE,
            {
                "slide_index": slide_idx,
                "rows": rows_val,
                "cols": cols_val,
                "x": x,
                "y": y,
                "cx": cx,
                "cy": cy,
            },
        )
        shape_id = result.get("shape_id")
        if not isinstance(shape_id, int):
            msg = "bridge response shape_id must be an int"
            raise TypeError(msg)

        flags: dict[str, bool] = {
            k: bool(kwargs[k])
            for k in (
                "first_row",
                "first_col",
                "last_row",
                "last_col",
                "band_row",
                "band_col",
            )
            if k in kwargs
        }
        _populate_table(
            self,
            slide_idx,
            shape_id,
            (
                cast("list[list[str]] | None", kwargs.get("data")),
                cast("list[int] | None", kwargs.get("column_widths")),
            ),
            flags,
        )
        return shape_id

    def get_table(self, slide_index: int, shape_id: int) -> TableInfo:
        """Return serialized table information for a table shape."""
        result = self.execute(
            ops.OP_GET_TABLE,
            {"slide_index": slide_index, "shape_id": shape_id},
        )
        return cast("TableInfo", cast("dict[str, object]", result.get("table", {})))
