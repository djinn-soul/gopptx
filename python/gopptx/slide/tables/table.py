"""Table proxy classes for gopptx."""
# pyright: reportMissingSuperCall=false, reportPrivateUsage=false, reportOptionalMemberAccess=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast, overload

from ... import api_errors, ops
from ...utils import normalize_table_index
from .table_cells import Cell, CellRange
from .table_collections import TableColumn, TableColumns, TableRow, TableRows
from .table_flags_mixin import TableFlagsMixin

if TYPE_CHECKING:
    from collections.abc import Iterator
    from typing import Protocol

    from ..contracts import SlidePresentationProtocol

    class _TableBandingProto(Protocol):
        _cache: dict[str, object] | None

        def _ensure_cache(self) -> None: ...

        def _update_flags(self, flags: dict[str, bool]) -> None: ...


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


class _TableBandingMixin:
    """Banding flag properties shared by table proxies."""

    @property
    def header_row_enabled(self: _TableBandingProto) -> bool:
        """Whether the first row is formatted as a header row."""
        self._ensure_cache()
        if self._cache is None:
            return False
        return self._cache.get("first_row", False) is True

    @header_row_enabled.setter
    def header_row_enabled(self: _TableBandingProto, value: bool) -> None:
        self._update_flags({"first_row": value})

    @property
    def banded_rows_enabled(self: _TableBandingProto) -> bool:
        """Whether alternating row banding is enabled."""
        self._ensure_cache()
        if self._cache is None:
            return False
        return self._cache.get("band_row", False) is True

    @banded_rows_enabled.setter
    def banded_rows_enabled(self: _TableBandingProto, value: bool) -> None:
        self._update_flags({"band_row": value})


class Table(_TableBandingMixin, TableFlagsMixin):
    """Pythonic Table API supporting ``table[row, col]`` and slices."""

    def __init__(
        self, prs: SlidePresentationProtocol, slide_index: int, shape_id: int
    ) -> None:
        """Initialize a table proxy bound to a slide shape."""
        self.prs = prs
        self.slide_index = slide_index
        self.shape_id = shape_id
        self._cache: dict[str, object] | None = None
        self._cell_map: dict[tuple[int, int], dict[str, object]] = {}
        self._row_count: int | None = None
        self._col_count: int | None = None
        if not getattr(self.prs, "_batch_active", False):
            self._ensure_cache()

    @overload
    def __getitem__(self, idx: tuple[int, int]) -> Cell: ...

    @overload
    def __getitem__(self, idx: tuple[int | slice, int | slice]) -> Cell | CellRange: ...

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
        """Clear local table caches when not inside a batch."""
        if getattr(self.prs, "_batch_active", False):
            return
        self._cache = None
        self._cell_map = {}

    def table_state(self) -> dict[str, object]:
        """Return cached table state, loading it if needed."""
        self._ensure_cache()
        return self._cache or {}

    def get_cell_info(self, row: int, col: int) -> dict[str, object]:
        """Return raw cell metadata for ``(row, col)`` when present."""
        self._ensure_cache()
        return self._cell_map.get((row, col), {})

    def update_cell(self, row: int, col: int, updates: dict[str, object]) -> None:
        """Apply updates to a single cell."""
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
        """Number of rows in the table."""
        if self._row_count is not None:
            return self._row_count
        self._ensure_cache()
        return self._row_count or 0

    @property
    def col_count(self) -> int:
        """Number of columns in the table."""
        if self._col_count is not None:
            return self._col_count
        self._ensure_cache()
        return self._col_count or 0

    def __getitem__(self, idx: tuple[int | slice, int | slice]) -> Cell | CellRange:
        """Return a cell or rectangular range for the given index tuple."""
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
        """Return the cell at the given zero-based row and column."""
        return self[row, col]

    def iter_cells(self) -> Iterator[Cell]:
        """Iterate all cells row-major."""
        for row in range(self.row_count):
            for col in range(self.col_count):
                yield Cell(self, row, col)

    def set_data(self, rows: list[list[str]]) -> None:
        """Replace all cell text in the table with new data."""
        if getattr(self.prs, "_batch_active", False):
            raise api_errors.GopptxError(
                "Cannot bulk-replace table data in batch context",
                code="BATCH_STRUCTURAL_CHANGE_NOT_ALLOWED",
            )

        if len(rows) != self.row_count:
            raise ValueError(
                f"Row count mismatch: expected {self.row_count}, got {len(rows)}"
            )

        for row_idx, row in enumerate(rows):
            if len(row) != self.col_count:
                raise ValueError(
                    f"Row {row_idx}: column count mismatch: expected {self.col_count}, got {len(row)}"
                )
            for col_idx, text in enumerate(row):
                self[row_idx, col_idx].text = text

        if not getattr(self.prs, "_batch_active", False):
            self.invalidate_cache()

    @property
    def rows(self) -> TableRows:
        """A row collection proxy."""
        return TableRows(self)

    @property
    def columns(self) -> TableColumns:
        """A column collection proxy."""
        return TableColumns(self)
