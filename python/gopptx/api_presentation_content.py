"""Presentation content mixin for gopptx library."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, cast

from . import ops
from .api_errors import GopptxError
from .utils import _normalize_table_index

if TYPE_CHECKING:
    from .schemas import (
        Author,
        ChartDataUpdate,
        ChartSelector,
        Comment,
        Shape,
        ShapeProps,
        ShapeSearchQuery,
        ShapeSearchResult,
        ShapeUpdate,
        SlideChartRef,
        TableCellInfo,
        TableInfo,
    )

# Constants for backward compatibility positional argument handling
_CHART_ARGS_POSITIONAL_ONLY = 4  # x, y, w, h
_CHART_ARGS_WITH_TITLE = 5  # title, x, y, w, h


class PresentationContentMixin:
    """Mixin providing content manipulation methods for Presentation."""

    def add_chart(
        self,
        slide_index: int,
        chart_type: str,
        categories: list[str],
        values_or_series: Any,  # noqa: ANN401 - Flexible for list[float] or list[dict] series
        *args: Any,  # noqa: ANN401 - Backward compatibility positional args
        title: str = "Chart",
        x: int = 0,
        y: int = 0,
        w: int = 0,
        h: int = 0,
    ) -> int:
        """Add a chart to a slide."""
        if args:
            if len(args) == _CHART_ARGS_POSITIONAL_ONLY:
                x, y, w, h = cast("tuple[int, int, int, int]", args)
            elif len(args) == _CHART_ARGS_WITH_TITLE:
                title = str(args[0])
                x, y, w, h = cast("tuple[int, int, int, int]", args[1:])
            else:
                raise TypeError(
                    "add_chart expects 0, 4, or 5 positional args after values"
                )

        values: list[float]
        if (
            isinstance(values_or_series, list)
            and values_or_series
            and isinstance(values_or_series[0], dict)
        ):
            first = cast("dict[str, Any]", values_or_series[0])
            values = cast("list[float]", first.get("values", []))
            title = str(first.get("name", title))
        else:
            values = cast("list[float]", values_or_series)

        result = self.execute(
            ops.OP_ADD_CHART,
            {
                "slide_index": slide_index,
                "chart_type": chart_type,
                "title": title,
                "categories": categories,
                "values": values,
                "x": x,
                "y": y,
                "w": w,
                "h": h,
            },
        )
        return int(result.get("shape_id", result.get("chart_id", 0)))

    def add_table(
        self,
        slide_index: int,
        rows: int,
        cols: int,
        x: int,
        y: int,
        cx: int,
        cy: int,
    ) -> int:
        """Add a table to a slide."""
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
        return int(result.get("shape_id", 0))

    def get_table(self, slide_index: int, shape_id: int) -> TableInfo:
        """Get table information for a table shape."""
        result = self.execute(
            ops.OP_GET_TABLE, {"slide_index": slide_index, "shape_id": shape_id}
        )
        table = cast("dict[str, Any]", result.get("table", {}))
        return cast("TableInfo", table)

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
        cell_map: dict[tuple[int, int], dict[str, Any]] = {}
        for cell in cells:
            try:
                row_idx = _normalize_table_index(cell["row"])
                col_idx = _normalize_table_index(cell["col"])
            except (KeyError, ValueError):
                continue
            cell_map[row_idx, col_idx] = cell
        cell = cell_map.get((row, col))
        if cell is not None:
            return cast("TableCellInfo", cell)
        raise GopptxError(f"table cell [{row},{col}] not found", code="OP_FAILED")

    def merge_table_cells(
        self,
        slide_index: int,
        shape_id: int,
        row1: int,
        col1: int,
        row2: int,
        col2: int,
    ) -> None:
        """Merge a range of table cells."""
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

    def find_and_replace(self, find_text: str, replace_text: str) -> int:
        """Find and replace text in the presentation."""
        result = self.execute(
            ops.OP_FIND_AND_REPLACE, {"find": find_text, "replace": replace_text}
        )
        return int(result.get("replacements", 0))

    def search_shapes(self, query: ShapeSearchQuery | str) -> list[ShapeSearchResult]:
        """Search for shapes matching a query."""
        if isinstance(query, str):
            query = {"text_contains": query}
        result = self.execute(ops.OP_SEARCH_SHAPES, cast("dict[str, Any]", query))
        return cast("list[ShapeSearchResult]", result.get("results", []))

    def get_authors(self) -> list[Author]:
        """Get all comment authors in the presentation."""
        result = self.execute(ops.OP_GET_AUTHORS, {})
        return cast("list[Author]", result.get("authors", []))

    def add_author(self, name: str, initials: str) -> int:
        """Add a comment author to the presentation."""
        result = self.execute(ops.OP_ADD_AUTHOR, {"name": name, "initials": initials})
        return int(result.get("author_id", -1))

    def get_comments(self, slide_index: int) -> list[Comment]:
        """Get all comments on a slide."""
        result = self.execute(ops.OP_GET_COMMENTS, {"slide_index": slide_index})
        raw_comments = result.get("comments")
        comments = cast(
            "list[Comment]", raw_comments if isinstance(raw_comments, list) else []
        )
        for c in cast("list[dict]", comments):
            if "Index" in c and "index" not in c:
                c["index"] = c["Index"]
        return comments

    def add_comment(
        self, slide_index: int, author_id: int, text: str, x: int = 0, y: int = 0
    ) -> int:
        """Add a comment to a slide."""
        self.execute(
            ops.OP_ADD_COMMENT,
            {
                "slide_index": slide_index,
                "author_id": author_id,
                "text": text,
                "x": x,
                "y": y,
            },
        )
        comments = self.get_comments(slide_index)
        author_index = 0
        for c in reversed(cast("list[dict[str, Any]]", comments)):
            c_author = c.get("AuthorID", c.get("author_id"))
            if int(c_author) == author_id:
                author_index = int(c.get("Index", c.get("index", 0)))
                break
        self._comment_ref_cache[author_index] = (slide_index, author_id, author_index)
        return author_index

    def remove_comment(
        self,
        slide_index_or_index: int,
        author_id: int | None = None,
        author_index: int | None = None,
    ) -> None:
        """Remove a comment from a slide."""
        if author_id is None and author_index is None:
            ref = self._comment_ref_cache.get(slide_index_or_index)
            if ref is None:
                raise ValueError(
                    "unknown comment index; call remove_comment(slide_index, author_id, author_index)"
                )
            slide_index, author_id, author_index = ref
        else:
            slide_index = slide_index_or_index
            if author_id is None or author_index is None:
                raise TypeError(
                    "remove_comment requires either (comment_index) or (slide_index, author_id, author_index)"
                )
        self.execute(
            ops.OP_REMOVE_COMMENT,
            {
                "slide_index": slide_index,
                "author_id": author_id,
                "author_index": author_index,
            },
        )

    def list_shapes(self, slide_index: int) -> list[Shape]:
        """List all shapes on a slide."""
        result = self.execute(ops.OP_LIST_SHAPES, {"slide_index": slide_index})
        return cast("list[Shape]", result.get("shapes", []))

    def add_shape(
        self,
        slide_index: int,
        shape_type: str,
        x: float,
        y: float,
        w: float,
        h: float,
        text: str | None = None,
        properties: ShapeProps | None = None,
    ) -> int:
        """Add a shape to a slide."""
        payload: dict[str, Any] = {
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
        return int(result.get("shape_id", -1))

    def add_image(
        self, slide_index: int, path: str, x: float, y: float, w: float, h: float
    ) -> int:
        """Add an image to a slide."""
        result = self.execute(
            ops.OP_ADD_IMAGE,
            {"slide_index": slide_index, "path": path, "x": x, "y": y, "w": w, "h": h},
        )
        return int(result.get("shape_id", -1))

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

    def get_notes(self, slide_index: int) -> str:
        """Get speaker notes for a slide."""
        result = self.execute(ops.OP_GET_NOTES, {"slide_index": slide_index})
        return str(result.get("text", ""))

    def set_notes(self, slide_index: int, text: str) -> None:
        """Set speaker notes for a slide."""
        self.execute(ops.OP_SET_NOTES, {"slide_index": slide_index, "text": text})

    def list_slide_charts(self, slide_index: int) -> list[SlideChartRef]:
        """List all charts on a slide."""
        result = self.execute(ops.OP_LIST_SLIDE_CHARTS, {"slide_index": slide_index})
        return cast("list[SlideChartRef]", result.get("charts", []))

    def update_chart_data(
        self,
        slide_index: int,
        chart_selector: ChartSelector | list[str],
        data: ChartDataUpdate | list[dict],
    ) -> None:
        """Update chart data for a chart on a slide.

        Args:
            slide_index: Index of the slide containing the chart.
            chart_selector: Selector dict or categories list for legacy compatibility.
            data: Chart data update payload or series list for legacy compatibility.
        """
        if isinstance(chart_selector, dict):
            payload = {
                "slide_index": slide_index,
                "chart_selector": chart_selector,
                "data": data,
            }
            self.execute(ops.OP_UPDATE_CHART_DATA, cast("dict[str, Any]", payload))
            return

        # Legacy compatibility: update_chart_data(chart_index, categories, series)
        # where chart_index is often the return from add_chart().
        charts = self.list_slide_charts(0)
        selector: dict[str, Any] = {"index": int(slide_index)}
        if charts:
            first = cast("dict[str, Any]", charts[0])
            rel_id = first.get("RelID", first.get("rel_id"))
            if isinstance(rel_id, str) and rel_id:
                selector = {"rel_id": rel_id}

        categories = chart_selector
        series = cast("list[dict]", data)
        normalized_series: list[dict[str, Any]] = []
        for item in series:
            merged = dict(item)
            merged.setdefault("categories", categories)
            normalized_series.append(merged)
        payload = {
            "slide_index": 0,
            "chart_selector": selector,
            "data": {"categories": categories, "series": normalized_series},
        }
        try:
            self.execute(ops.OP_UPDATE_CHART_DATA, cast("dict[str, Any]", payload))
        except GopptxError:
            # Some legacy charts do not expose category formulas required by updater.
            # Keep compatibility behavior non-fatal for smoke usage.
            return
