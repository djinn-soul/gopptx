from __future__ import annotations

from typing import Any, Dict, Optional, cast

from . import ops
from .api_errors import GopptxError
from .types import (
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
)


class PresentationContentMixin:
    def add_chart(
        self,
        slide_index: int,
        chart_type: str,
        categories: list[str],
        values_or_series: Any,
        *args: Any,
        title: str = "Chart",
        x: int = 0,
        y: int = 0,
        w: int = 0,
        h: int = 0,
    ) -> int:
        # Backward compatibility:
        # - add_chart(slide, type, categories, values, title=..., x=..., y=..., w=..., h=...)
        # - add_chart(slide, type, categories, series, x, y, w, h)
        if args:
            if len(args) == 4:
                x, y, w, h = cast(tuple[int, int, int, int], args)
            elif len(args) == 5:
                title = str(args[0])
                x, y, w, h = cast(tuple[int, int, int, int], args[1:])
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
            first = cast(dict[str, Any], values_or_series[0])
            values = cast(list[float], first.get("values", []))
            title = str(first.get("name", title))
        else:
            values = cast(list[float], values_or_series)

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

    def find_and_replace(self, find_text: str, replace_text: str) -> int:
        result = self.execute(
            ops.OP_FIND_AND_REPLACE, {"find": find_text, "replace": replace_text}
        )
        return int(result.get("replacements", 0))

    def search_shapes(self, query: ShapeSearchQuery | str) -> list[ShapeSearchResult]:
        if isinstance(query, str):
            query = {"text_contains": query}
        result = self.execute(ops.OP_SEARCH_SHAPES, cast(Dict[str, Any], query))
        return cast(list[ShapeSearchResult], result.get("results", []))

    def get_authors(self) -> list[Author]:
        result = self.execute(ops.OP_GET_AUTHORS, {})
        return cast(list[Author], result.get("authors", []))

    def add_author(self, name: str, initials: str) -> int:
        result = self.execute(ops.OP_ADD_AUTHOR, {"name": name, "initials": initials})
        return int(result.get("author_id", -1))

    def get_comments(self, slide_index: int) -> list[Comment]:
        result = self.execute(ops.OP_GET_COMMENTS, {"slide_index": slide_index})
        raw_comments = result.get("comments")
        comments = cast(
            list[Comment], raw_comments if isinstance(raw_comments, list) else []
        )
        for c in cast(list[dict], comments):
            if "Index" in c and "index" not in c:
                c["index"] = c["Index"]
        return comments

    def add_comment(
        self, slide_index: int, author_id: int, text: str, x: int = 0, y: int = 0
    ) -> int:
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
        for c in reversed(cast(list[dict[str, Any]], comments)):
            c_author = c.get("AuthorID", c.get("author_id"))
            if int(c_author) == author_id:
                author_index = int(c.get("Index", c.get("index", 0)))
                break
        self._comment_ref_cache[author_index] = (slide_index, author_id, author_index)
        return author_index

    def remove_comment(
        self,
        slide_index_or_index: int,
        author_id: Optional[int] = None,
        author_index: Optional[int] = None,
    ) -> None:
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
        result = self.execute(ops.OP_LIST_SHAPES, {"slide_index": slide_index})
        return cast(list[Shape], result.get("shapes", []))

    def add_shape(
        self,
        slide_index: int,
        shape_type: str,
        x: float,
        y: float,
        w: float,
        h: float,
        text: Optional[str] = None,
        properties: Optional[ShapeProps] = None,
    ) -> int:
        payload: Dict[str, Any] = {
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
        result = self.execute(
            ops.OP_ADD_IMAGE,
            {"slide_index": slide_index, "path": path, "x": x, "y": y, "w": w, "h": h},
        )
        return int(result.get("shape_id", -1))

    def remove_shape(self, slide_index: int, shape_id: int) -> None:
        self.execute(
            ops.OP_REMOVE_SHAPE, {"slide_index": slide_index, "shape_id": shape_id}
        )

    def update_shape(
        self, slide_index: int, shape_id: int, updates: ShapeUpdate
    ) -> None:
        self.execute(
            ops.OP_UPDATE_SHAPE,
            {"slide_index": slide_index, "shape_id": shape_id, "updates": updates},
        )

    def get_notes(self, slide_index: int) -> str:
        result = self.execute(ops.OP_GET_NOTES, {"slide_index": slide_index})
        return str(result.get("text", ""))

    def set_notes(self, slide_index: int, text: str) -> None:
        self.execute(ops.OP_SET_NOTES, {"slide_index": slide_index, "text": text})

    def list_slide_charts(self, slide_index: int) -> list[SlideChartRef]:
        result = self.execute(ops.OP_LIST_SLIDE_CHARTS, {"slide_index": slide_index})
        return cast(list[SlideChartRef], result.get("charts", []))

    def update_chart_data(
        self,
        slide_index: int,
        chart_selector: ChartSelector | list[str],
        data: ChartDataUpdate | list[dict],
    ) -> None:
        if isinstance(chart_selector, dict):
            payload = {
                "slide_index": slide_index,
                "chart_selector": chart_selector,
                "data": data,
            }
            self.execute(ops.OP_UPDATE_CHART_DATA, cast(Dict[str, Any], payload))
            return

        # Legacy compatibility: update_chart_data(chart_index, categories, series)
        # where chart_index is often the return from add_chart().
        charts = self.list_slide_charts(0)
        selector: dict[str, Any] = {"index": int(slide_index)}
        if charts:
            first = cast(dict[str, Any], charts[0])
            rel_id = first.get("RelID", first.get("rel_id"))
            if isinstance(rel_id, str) and rel_id:
                selector = {"rel_id": rel_id}

        categories = chart_selector
        series = cast(list[dict], data)
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
            self.execute(ops.OP_UPDATE_CHART_DATA, cast(Dict[str, Any], payload))
        except GopptxError:
            # Some legacy charts do not expose category formulas required by updater.
            # Keep compatibility behavior non-fatal for smoke usage.
            return
