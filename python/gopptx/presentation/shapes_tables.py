"""Shape, text, notes, and table mixins for the Presentation API."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .. import ops
from ..api_errors import GopptxError
from ..utils import normalize_table_index
from .helpers import PresentationProtocol

if TYPE_CHECKING:
    from ..schemas import (
        Shape,
        ShapeProps,
        ShapeSearchQuery,
        ShapeSearchResult,
        ShapeUpdate,
        TableCellInfo,
        TableInfo,
    )
else:

    class PresentationProtocol:
        """Runtime placeholder to avoid Protocol abstract behavior."""


class PresentationNotesMixin(PresentationProtocol):
    """Mixin providing speaker notes methods."""

    def get_notes(self, slide_index: int) -> str:
        """Get speaker notes for a slide."""
        result = self.execute(ops.OP_GET_NOTES, {"slide_index": slide_index})
        return str(cast("str", result.get("text", "")))

    def set_notes(self, slide_index: int, text: str) -> None:
        """Set speaker notes for a slide."""
        self.execute(ops.OP_SET_NOTES, {"slide_index": slide_index, "text": text})


class PresentationShapeMixin(PresentationProtocol):
    """Mixin providing shape manipulation methods."""

    def search_shapes(self, query: ShapeSearchQuery | str) -> list[ShapeSearchResult]:
        """Search for shapes matching a query."""
        if isinstance(query, str):
            query = {"text_contains": query}
        result = self.execute(ops.OP_SEARCH_SHAPES, cast("dict[str, object]", query))
        return cast("list[ShapeSearchResult]", result.get("results", []))

    def list_shapes(self, slide_index: int) -> list[Shape]:
        """List all shapes on a slide."""
        result = self.execute(ops.OP_LIST_SHAPES, {"slide_index": slide_index})
        return cast("list[Shape]", result.get("shapes", []))

    def add_shape(
        self,
        slide_index: int,
        shape_type: str,
        bounds: tuple[float, float, float, float],
        **kwargs: str | ShapeProps,
    ) -> int:
        """Add a shape to a slide."""
        x, y, w, h = bounds
        text = kwargs.get("text")
        properties = kwargs.get("properties")
        payload: dict[str, object] = {
            "slide_index": slide_index,
            "type": shape_type,
            "x": x,
            "y": y,
            "w": w,
            "h": h,
        }
        if text is not None:
            payload["text"] = text
        if properties is not None:
            payload["properties"] = properties
        result = self.execute(ops.OP_ADD_SHAPE, payload)
        return int(cast("int", result.get("shape_id", -1)))

    def add_image(
        self,
        slide_index: int,
        path: str,
        bounds: tuple[float, float, float, float],
    ) -> int:
        """Add an image to a slide."""
        x, y, w, h = bounds
        result = self.execute(
            ops.OP_ADD_IMAGE,
            {"slide_index": slide_index, "path": path, "x": x, "y": y, "w": w, "h": h},
        )
        return int(cast("int", result.get("shape_id", -1)))

    def remove_shape(self, slide_index: int, shape_id: int) -> None:
        """Remove a shape from a slide."""
        self.execute(
            ops.OP_REMOVE_SHAPE, {"slide_index": slide_index, "shape_id": shape_id}
        )

    def move_shape_to_front(self, slide_index: int, shape_id: int) -> None:
        """Move a shape to the front of the z-order."""
        self.execute(
            ops.OP_MOVE_SHAPE_TO_FRONT,
            {"slide_index": slide_index, "shape_id": shape_id},
        )

    def move_shape_to_back(self, slide_index: int, shape_id: int) -> None:
        """Move a shape to the back of the z-order."""
        self.execute(
            ops.OP_MOVE_SHAPE_TO_BACK,
            {"slide_index": slide_index, "shape_id": shape_id},
        )

    def update_shape(
        self, slide_index: int, shape_id: int, updates: ShapeUpdate
    ) -> None:
        """Update shape properties."""
        self.execute(
            ops.OP_UPDATE_SHAPE,
            {"slide_index": slide_index, "shape_id": shape_id, "updates": updates},
        )


class PresentationTextMixin(PresentationProtocol):
    """Mixin providing text search and replace methods."""

    def find_and_replace(self, find_text: str, replace_text: str) -> int:
        """Find and replace text in the presentation."""
        result = self.execute(
            ops.OP_FIND_AND_REPLACE, {"find": find_text, "replace": replace_text}
        )
        return int(cast("int", result.get("replacements", 0)))


class PresentationTableMixin(PresentationProtocol):
    """Mixin providing table creation and manipulation methods."""

    def add_table(
        self,
        slide_index: int,
        rows: int,
        cols: int,
        bounds: tuple[int, int, int, int],
    ) -> int:
        """Add a table to a slide."""
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
        """Get table information for a table shape."""
        result = self.execute(
            ops.OP_GET_TABLE, {"slide_index": slide_index, "shape_id": shape_id}
        )
        return cast("TableInfo", cast("dict[str, object]", result.get("table", {})))

    def set_table_flags(
        self, slide_index: int, shape_id: int, flags: dict[str, bool]
    ) -> None:
        """Set table style flags."""
        self.execute(
            ops.OP_UPDATE_TABLE_FLAGS,
            {"slide_index": slide_index, "shape_id": shape_id, "flags": flags},
        )

    def set_table_cell_text(
        self, slide_index: int, shape_id: int, row: int, col: int, text: str
    ) -> None:
        """Set the text of a table cell."""
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
        self, slide_index: int, shape_id: int, row: int, col: int
    ) -> TableCellInfo:
        """Get information about a table cell."""
        table = self.get_table(slide_index, shape_id)
        cells = table.get("cells", [])
        cell_map: dict[tuple[int, int], dict[str, object]] = {}
        for cell in cells:
            try:
                row_idx = normalize_table_index(cell["row"])
                col_idx = normalize_table_index(cell["col"])
            except (KeyError, ValueError):
                continue
            cell_map[row_idx, col_idx] = cast("dict[str, object]", cell)
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
        """Merge a range of table cells."""
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
        self, slide_index: int, shape_id: int, row: int, col: int
    ) -> None:
        """Split a merged table cell."""
        self.execute(
            ops.OP_SPLIT_TABLE_CELL,
            {"slide_index": slide_index, "shape_id": shape_id, "row": row, "col": col},
        )
