from __future__ import annotations

from typing import TYPE_CHECKING, Any, Dict, Iterator, Optional, cast, Union, Tuple

from . import ops

if TYPE_CHECKING:
    from .api_presentation import Presentation


class Cell:
    def __init__(self, table: Table, row: int, col: int):
        self._table = table
        self.row = row
        self.col = col

    @property
    def is_merge_origin(self) -> bool:
        return self._table._get_cell_info(self.row, self.col).get(
            "is_merge_origin", False
        )

    @property
    def is_spanned(self) -> bool:
        return self._table._get_cell_info(self.row, self.col).get("is_spanned", False)

    @property
    def row_span(self) -> int:
        return self._table._get_cell_info(self.row, self.col).get("row_span", 1)

    @property
    def col_span(self) -> int:
        return self._table._get_cell_info(self.row, self.col).get("col_span", 1)

    @property
    def text(self) -> str:
        return str(self._table._get_cell_info(self.row, self.col).get("text", ""))

    @text.setter
    def text(self, value: str) -> None:
        self._table._update_cell(self.row, self.col, {"text": str(value)})

    def split(self) -> None:
        """Splits a merged cell back into a 1x1 cell."""
        self._table._prs.execute(
            ops.OP_SPLIT_TABLE_CELL,
            {
                "slide_index": self._table._slide_index,
                "shape_id": self._table._shape_id,
                "row": self.row,
                "col": self.col,
            },
        )
        self._table._invalidate_cache()

    def __repr__(self) -> str:
        return f"<Cell [{self.row}, {self.col}] text={self.text!r}>"


class CellRange:
    """Represents a 2D slice of cells in a table, allowing bulk operations like merging."""

    def __init__(
        self, table: Table, row_start: int, row_end: int, col_start: int, col_end: int
    ):
        self._table = table
        self.row_start = max(0, row_start)
        self.row_end = min(table.row_count, row_end)
        self.col_start = max(0, col_start)
        self.col_end = min(table.col_count, col_end)

    def merge(self) -> None:
        """
        Merges all cells in this range into a single spanned cell.
        Follows Python exclusive slice convention for end indices.
        """
        if self.row_end <= self.row_start + 1 and self.col_end <= self.col_start + 1:
            return  # Nothing to merge

        self._table._prs.execute(
            ops.OP_MERGE_TABLE_CELLS,
            {
                "slide_index": self._table._slide_index,
                "shape_id": self._table._shape_id,
                # Python slices are exclusive at the end, calculate last index inclusive for Go backend  
                "row1": self.row_start,
                "col1": self.col_start,
                "row2": self.row_end - 1,
                "col2": self.col_end - 1,
            },
        )
        self._table._invalidate_cache()


class Table:
    """
    Pythonic Table API for gopptx.
    Provides grid-based access via slicing: table[row, col] or table[r1:r2, c1:c2].
    """

    def __init__(self, prs: Presentation, slide_index: int, shape_id: int):
        self._prs = prs
        self._slide_index = slide_index
        self._shape_id = shape_id
        self._cache: Optional[Dict[str, Any]] = None
        self._cell_map: Dict[Tuple[int, int], Dict[str, Any]] = {}

    def _ensure_cache(self) -> None:
        if self._cache is None:
            res = self._prs.execute(
                ops.OP_GET_TABLE,
                {"slide_index": self._slide_index, "shape_id": self._shape_id},
            )
            self._cache = cast(Dict[str, Any], res.get("table", {}))
            
            # Optimization: Build coordinate map for O(1) cell lookup
            self._cell_map = {}
            cells = self._cache.get("cells", [])
            for c in cells:
                row = c.get("row")
                col = c.get("col")
                if row is not None and col is not None:
                    self._cell_map[(row, col)] = c

    def _invalidate_cache(self) -> None:
        self._cache = None
        self._cell_map = {}

    def _get_cell_info(self, row: int, col: int) -> Dict[str, Any]:
        self._ensure_cache()
        return self._cell_map.get((row, col), {})

    def _update_cell(self, row: int, col: int, updates: Dict[str, Any]) -> None:
        self._prs.execute(
            ops.OP_UPDATE_TABLE_CELL,
            {
                "slide_index": self._slide_index,
                "shape_id": self._shape_id,
                "row": row,
                "col": col,
                "updates": updates,
            },
        )
        self._invalidate_cache()

    @property
    def row_count(self) -> int:
        self._ensure_cache()
        return int(self._cache.get("row_count", 0))  # type: ignore

    @property
    def col_count(self) -> int:
        self._ensure_cache()
        return int(self._cache.get("col_count", 0))  # type: ignore

    def __getitem__(
        self, idx: Tuple[Union[int, slice], Union[int, slice]]
    ) -> Union[Cell, CellRange]:
        if not isinstance(idx, tuple) or len(idx) != 2:
            raise TypeError("Table indices must be a tuple of (row, col)")

        row_idx, col_idx = idx

        if isinstance(row_idx, slice) or isinstance(col_idx, slice):
            # Resolve slices
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

        # Single cell access
        if (
            row_idx < 0
            or row_idx >= self.row_count
            or col_idx < 0
            or col_idx >= self.col_count
        ):
            raise IndexError("Cell index out of range")
        return Cell(self, row_idx, col_idx)

    # Alias for compatibility with traditional APIs
    def cell(self, row: int, col: int) -> Cell:
        return self[row, col]  # type: ignore

    def iter_cells(self) -> Iterator[Cell]:
        for r in range(self.row_count):
            for c in range(self.col_count):
                yield Cell(self, r, c)

    # Style Flags
    def _update_flags(self, flags: Dict[str, bool]) -> None:
        self._prs.execute(
            ops.OP_UPDATE_TABLE_FLAGS,
            {
                "slide_index": self._slide_index,
                "shape_id": self._shape_id,
                "flags": flags,
            },
        )
        self._invalidate_cache()

    @property
    def has_header_row(self) -> bool:
        self._ensure_cache()
        return bool(self._cache.get("first_row", False))  # type: ignore

    @has_header_row.setter
    def has_header_row(self, value: bool) -> None:
        self._update_flags({"first_row": value})

    @property
    def has_banded_rows(self) -> bool:
        self._ensure_cache()
        return bool(self._cache.get("band_row", False))  # type: ignore

    @has_banded_rows.setter
    def has_banded_rows(self, value: bool) -> None:
        self._update_flags({"band_row": value})
