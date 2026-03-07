"""Table proxy classes for gopptx."""
# ruff: noqa: D102,D105,D107
# pyright: reportMissingSuperCall=false, reportPrivateUsage=false, reportOptionalMemberAccess=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .. import ops
from ..utils import normalize_table_index
from .table_cells import Cell, CellRange
from .table_collections import TableColumn, TableColumns, TableRow, TableRows

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ..presentation.presentation import Presentation

_TABLE_INDEX_DIMENSIONS = 2

__all__ = [
    "Cell",
    "CellRange",
    "Table",
    "TableColumn",
    "TableColumns",
    "TableRow",
    "TableRows",
]


class Table:
    """Pythonic Table API supporting table[row, col] and slices."""

    def __init__(self, prs: Presentation, slide_index: int, shape_id: int) -> None:
        self.prs = prs
        self.slide_index = slide_index
        self.shape_id = shape_id
        self._cache: dict[str, object] | None = None
        self._cell_map: dict[tuple[int, int], dict[str, object]] = {}
        self._row_count: int | None = None
        self._col_count: int | None = None
        if not getattr(self.prs, "_batch_active", False):
            self._ensure_cache()

    def _ensure_cache(self) -> None:
        if self._cache is not None:
            return
        res = self.prs.execute(
            ops.OP_GET_TABLE,
            {"slide_index": self.slide_index, "shape_id": self.shape_id},
        )
        self._cache = cast("dict[str, object]", res.get("table", {}))
        cells = cast("list[dict[str, object]]", self._cache.get("cells", []))
        for cell in cells:
            try:
                row_idx = normalize_table_index(cell["row"])
                col_idx = normalize_table_index(cell["col"])
            except (KeyError, ValueError):
                continue
            self._cell_map[row_idx, col_idx] = cell
        self._row_count = int(cast("int", self._cache.get("row_count", 0)))
        self._col_count = int(cast("int", self._cache.get("col_count", 0)))

    def invalidate_cache(self) -> None:
        if getattr(self.prs, "_batch_active", False):
            return
        self._cache = None
        self._cell_map = {}

    def get_cell_info(self, row: int, col: int) -> dict[str, object]:
        self._ensure_cache()
        return self._cell_map.get((row, col), {})

    def update_cell(self, row: int, col: int, updates: dict[str, object]) -> None:
        self.prs.execute(
            ops.OP_UPDATE_TABLE_CELL,
            {
                "slide_index": self.slide_index,
                "shape_id": self.shape_id,
                "row": row,
                "col": col,
                "updates": updates,
            },
        )
        if (row, col) in self._cell_map:
            self._cell_map[row, col].update(updates)
        if not getattr(self.prs, "_batch_active", False):
            self.invalidate_cache()

    @property
    def row_count(self) -> int:
        if self._row_count is not None:
            return self._row_count
        self._ensure_cache()
        return self._row_count or 0

    @property
    def col_count(self) -> int:
        if self._col_count is not None:
            return self._col_count
        self._ensure_cache()
        return self._col_count or 0

    def __getitem__(self, idx: tuple[int | slice, int | slice]) -> Cell | CellRange:
        if len(idx) != _TABLE_INDEX_DIMENSIONS:
            raise TypeError("Table indices must be a tuple of (row, col)")
        row_idx, col_idx = idx
        if isinstance(row_idx, int) and row_idx < 0:
            row_idx += self.row_count
        if isinstance(col_idx, int) and col_idx < 0:
            col_idx += self.col_count

        if isinstance(row_idx, slice) or isinstance(col_idx, slice):
            r_start, r_end, r_step = (
                row_idx.indices(self.row_count)
                if isinstance(row_idx, slice)
                else (row_idx, row_idx + 1, 1)
            )
            c_start, c_end, c_step = (
                col_idx.indices(self.col_count)
                if isinstance(col_idx, slice)
                else (col_idx, col_idx + 1, 1)
            )
            if r_step != 1 or c_step != 1:
                raise ValueError("Table slicing does not support steps other than 1")
            return CellRange(self, r_start, r_end, c_start, c_end)

        if (
            row_idx < 0
            or row_idx >= self.row_count
            or col_idx < 0
            or col_idx >= self.col_count
        ):
            raise IndexError("Cell index out of range")
        return Cell(self, row_idx, col_idx)

    def cell(self, row: int, col: int) -> Cell:
        return self[row, col]  # type: ignore[return-value]

    def iter_cells(self) -> Iterator[Cell]:
        for row in range(self.row_count):
            for col in range(self.col_count):
                yield Cell(self, row, col)

    def _update_flags(self, flags: dict[str, bool]) -> None:
        self.prs.execute(
            ops.OP_UPDATE_TABLE_FLAGS,
            {
                "slide_index": self.slide_index,
                "shape_id": self.shape_id,
                "flags": flags,
            },
        )
        if not getattr(self.prs, "_batch_active", False):
            self.invalidate_cache()

    @property
    def has_header_row(self) -> bool:
        self._ensure_cache()
        return self._cache.get("first_row", False) is True if self._cache else False

    @has_header_row.setter
    def has_header_row(self, value: bool) -> None:
        self._update_flags({"first_row": value})

    @property
    def has_banded_rows(self) -> bool:
        self._ensure_cache()
        return self._cache.get("band_row", False) is True if self._cache else False

    @has_banded_rows.setter
    def has_banded_rows(self, value: bool) -> None:
        self._update_flags({"band_row": value})

    @property
    def first_row(self) -> bool:
        return self.has_header_row

    @first_row.setter
    def first_row(self, value: bool) -> None:
        self.has_header_row = value

    @property
    def horz_banding(self) -> bool:
        return self.has_banded_rows

    @horz_banding.setter
    def horz_banding(self, value: bool) -> None:
        self.has_banded_rows = value

    @property
    def first_col(self) -> bool:
        self._ensure_cache()
        return self._cache.get("first_col", False) is True if self._cache else False

    @first_col.setter
    def first_col(self, value: bool) -> None:
        self._update_flags({"first_col": value})

    @property
    def last_col(self) -> bool:
        self._ensure_cache()
        return self._cache.get("last_col", False) is True if self._cache else False

    @last_col.setter
    def last_col(self, value: bool) -> None:
        self._update_flags({"last_col": value})

    @property
    def last_row(self) -> bool:
        self._ensure_cache()
        return self._cache.get("last_row", False) is True if self._cache else False

    @last_row.setter
    def last_row(self, value: bool) -> None:
        self._update_flags({"last_row": value})

    @property
    def vert_banding(self) -> bool:
        self._ensure_cache()
        return self._cache.get("band_col", False) is True if self._cache else False

    @vert_banding.setter
    def vert_banding(self, value: bool) -> None:
        self._update_flags({"band_col": value})

    def apply_style(self, style_guid: str) -> None:
        self.prs.execute(
            ops.OP_SET_TABLE_STYLE,
            {
                "slide_index": self.slide_index,
                "shape_id": self.shape_id,
                "style_guid": style_guid,
            },
        )
        if not getattr(self.prs, "_batch_active", False):
            self.invalidate_cache()

    @property
    def rows(self) -> TableRows:
        return TableRows(self)

    @property
    def columns(self) -> TableColumns:
        return TableColumns(self)
