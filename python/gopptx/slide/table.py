"""Table and Cell proxy classes for gopptx library."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .. import ops
from ..utils import normalize_table_index
from .table_cells import Cell, CellRange

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ..presentation.presentation import Presentation

# Constants for table indexing
_TABLE_INDEX_DIMENSIONS = 2  # (row, col)
__all__ = ["Cell", "CellRange", "Table"]


class Table:
    """Pythonic Table API for gopptx.

    Provides grid-based access via slicing: table[row, col] or table[r1:r2, c1:c2].
    """

    def __init__(self, prs: Presentation, slide_index: int, shape_id: int) -> None:
        """Initialize the table proxy."""
        super().__init__()
        self.prs = prs  # Public for companion classes (Cell, CellRange)
        self.slide_index = slide_index  # Public for companion classes
        self.shape_id = shape_id  # Public for companion classes
        self._cache: dict[str, object] | None = None
        self._cell_map: dict[tuple[int, int], dict[str, object]] = {}
        self._row_count: int | None = None
        self._col_count: int | None = None

        # Pre-fetch table structure if not in batch mode to enable usage during batch
        if not getattr(self.prs, "_batch_active", False):
            self._ensure_cache()

    def _ensure_cache(self) -> None:
        if self._cache is None:
            res = self.prs.execute(
                ops.OP_GET_TABLE,
                {"slide_index": self.slide_index, "shape_id": self.shape_id},
            )
            self._cache = cast("dict[str, object]", res.get("table", {}))

            # Optimization: Build coordinate map for O(1) cell lookup
            self._cell_map = {}
            cells = cast("list[dict[str, object]]", self._cache.get("cells", []))
            for c in cells:
                try:
                    row_idx = normalize_table_index(c["row"])
                    col_idx = normalize_table_index(c["col"])
                except (KeyError, ValueError):
                    continue
                self._cell_map[row_idx, col_idx] = c

            # Permanent cache of dimensions
            self._row_count = int(cast("int", self._cache.get("row_count", 0)))
            self._col_count = int(cast("int", self._cache.get("col_count", 0)))

    def invalidate_cache(self) -> None:
        """Invalidate the table cache.

        If we are in a batch, we MUST NOT invalidate the cache as we cannot re-fetch it.
        This is safe because structural changes (merge/split) are rare and we can't
        properly support them in batch while also supporting batched cell updates.
        """
        if not getattr(self.prs, "_batch_active", False):
            self._cache = None
            self._cell_map = {}

    def get_cell_info(self, row: int, col: int) -> dict[str, object]:
        """Get cell information for the specified row and column."""
        self._ensure_cache()
        return self._cell_map.get((row, col), {})

    def update_cell(self, row: int, col: int, updates: dict[str, object]) -> None:
        """Updates a cell and its local cache representation.

        Note: The local cache update is a best-effort sync for properties like text.
        Structural changes or complex aggregate properties may still require a full
        re-fetch outside of batch mode.
        """
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
        # Update local cache to support reading updated values during batch
        if (row, col) in self._cell_map:
            self._cell_map[row, col].update(updates)

        # Do not invalidate the whole cache for simple cell updates,
        # especially during batch mode where re-fetching is forbidden.
        if not getattr(self.prs, "_batch_active", False):
            self.invalidate_cache()

    @property
    def row_count(self) -> int:
        """Get the number of rows in the table."""
        if self._row_count is not None:
            return self._row_count
        self._ensure_cache()
        return self._row_count or 0

    @property
    def col_count(self) -> int:
        """Get the number of columns in the table."""
        if self._col_count is not None:
            return self._col_count
        self._ensure_cache()
        return self._col_count or 0

    def __getitem__(self, idx: tuple[int | slice, int | slice]) -> Cell | CellRange:
        """Get a cell or range of cells by index."""
        if len(idx) != _TABLE_INDEX_DIMENSIONS:
            raise TypeError("Table indices must be a tuple of (row, col)")

        row_idx, col_idx = idx

        # Support negative indexing for integers (e.g., -1 for last row).
        # We use a single addition rather than modulo (%) to match standard Python
        # behavior where indices like [-len-1] raise IndexError instead of wrapping.
        if isinstance(row_idx, int) and row_idx < 0:
            row_idx += self.row_count
        if isinstance(col_idx, int) and col_idx < 0:
            col_idx += self.col_count

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
        """Get a cell by row and column index."""
        return self[row, col]  # type: ignore

    def iter_cells(self) -> Iterator[Cell]:
        """Iterate over all cells in the table."""
        for r in range(self.row_count):
            for c in range(self.col_count):
                yield Cell(self, r, c)

    # Style Flags
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
        """Check if the table has a header row."""
        self._ensure_cache()
        return self._cache.get("first_row", False) is True if self._cache else False

    @has_header_row.setter
    def has_header_row(self, value: bool) -> None:
        self._update_flags({"first_row": value})

    @property
    def has_banded_rows(self) -> bool:
        """Check if the table has banded rows."""
        self._ensure_cache()
        return self._cache.get("band_row", False) is True if self._cache else False

    @has_banded_rows.setter
    def has_banded_rows(self, value: bool) -> None:
        self._update_flags({"band_row": value})

    # python-pptx compatible aliases
    @property
    def first_row(self) -> bool:
        """python-pptx alias for has_header_row."""
        return self.has_header_row

    @first_row.setter
    def first_row(self, value: bool) -> None:
        self.has_header_row = value

    @property
    def horz_banding(self) -> bool:
        """python-pptx alias for has_banded_rows."""
        return self.has_banded_rows

    @horz_banding.setter
    def horz_banding(self, value: bool) -> None:
        self.has_banded_rows = value

    @property
    def first_col(self) -> bool:
        """When True, first column should have distinct formatting (side-heading column)."""
        self._ensure_cache()
        return self._cache.get("first_col", False) is True if self._cache else False

    @first_col.setter
    def first_col(self, value: bool) -> None:
        self._update_flags({"first_col": value})

    @property
    def last_col(self) -> bool:
        """When True, last column should have distinct formatting (totals column)."""
        self._ensure_cache()
        return self._cache.get("last_col", False) is True if self._cache else False

    @last_col.setter
    def last_col(self, value: bool) -> None:
        self._update_flags({"last_col": value})

    @property
    def last_row(self) -> bool:
        """When True, last row should have distinct formatting (totals row)."""
        self._ensure_cache()
        return self._cache.get("last_row", False) is True if self._cache else False

    @last_row.setter
    def last_row(self, value: bool) -> None:
        self._update_flags({"last_row": value})

    @property
    def vert_banding(self) -> bool:
        """When True, columns should have alternating shading."""
        self._ensure_cache()
        return self._cache.get("band_col", False) is True if self._cache else False

    @vert_banding.setter
    def vert_banding(self, value: bool) -> None:
        self._update_flags({"band_col": value})

    # Table Style
    def apply_style(self, style_guid: str) -> None:
        """Apply a table style to the table.

        Args:
            style_guid: A valid PowerPoint table style GUID, e.g.:
                "{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}" - Medium Style 2 - Accent 1
                "{B9AC3A68-259E-4EED-9050-4AE35E7F2B2D}" - Light Style 1
                "{5940675A-B579-460E-94D1-54222C63F5DA}" - Medium Style 1 - Accent 1
                "{3C2FF68D-0BFD-4DAC-8644-511001392665}" - Light Style 2 - Accent 1
                "{F370B699-2CC9-41A6-8B99-48FF35B8A595}" - Light Shading
                "{4EC70C17-9B8C-4085-A1A2-671F7FA0C7DD}" - Medium Shading 2 - Accent 1
        """
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
