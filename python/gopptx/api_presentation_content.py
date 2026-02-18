from __future__ import annotations

from typing import Any, Dict, Optional, cast

from . import ops
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
        values: list[float],
        title: str = "Chart",
        x: int = 0,
        y: int = 0,
        w: int = 0,
        h: int = 0,
    ) -> None:
        self.execute(
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

    def find_and_replace(self, find_text: str, replace_text: str) -> int:
        result = self.execute(ops.OP_FIND_AND_REPLACE, {"find": find_text, "replace": replace_text})
        return int(result.get("replacements", 0))

    def search_shapes(self, query: ShapeSearchQuery) -> list[ShapeSearchResult]:
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
        return cast(list[Comment], result.get("comments", []))

    def add_comment(self, slide_index: int, author_id: int, text: str, x: int = 0, y: int = 0) -> None:
        self.execute(ops.OP_ADD_COMMENT, {"slide_index": slide_index, "author_id": author_id, "text": text, "x": x, "y": y})

    def remove_comment(self, slide_index: int, author_id: int, author_index: int) -> None:
        self.execute(
            ops.OP_REMOVE_COMMENT,
            {"slide_index": slide_index, "author_id": author_id, "author_index": author_index},
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
        payload: Dict[str, Any] = {"slide_index": slide_index, "type": shape_type, "x": x, "y": y, "w": w, "h": h}
        if text is not None:
            payload["text"] = text
        if properties is not None:
            payload["properties"] = properties
        result = self.execute(ops.OP_ADD_SHAPE, payload)
        return int(result.get("shape_id", -1))

    def add_image(self, slide_index: int, path: str, x: float, y: float, w: float, h: float) -> int:
        result = self.execute(
            ops.OP_ADD_IMAGE,
            {"slide_index": slide_index, "path": path, "x": x, "y": y, "w": w, "h": h},
        )
        return int(result.get("shape_id", -1))

    def remove_shape(self, slide_index: int, shape_id: int) -> None:
        self.execute(ops.OP_REMOVE_SHAPE, {"slide_index": slide_index, "shape_id": shape_id})

    def update_shape(self, slide_index: int, shape_id: int, updates: ShapeUpdate) -> None:
        self.execute(ops.OP_UPDATE_SHAPE, {"slide_index": slide_index, "shape_id": shape_id, "updates": updates})

    def get_notes(self, slide_index: int) -> str:
        result = self.execute(ops.OP_GET_NOTES, {"slide_index": slide_index})
        return str(result.get("text", ""))

    def set_notes(self, slide_index: int, text: str) -> None:
        self.execute(ops.OP_SET_NOTES, {"slide_index": slide_index, "text": text})

    def list_slide_charts(self, slide_index: int) -> list[SlideChartRef]:
        result = self.execute(ops.OP_LIST_SLIDE_CHARTS, {"slide_index": slide_index})
        return cast(list[SlideChartRef], result.get("charts", []))

    def update_chart_data(self, slide_index: int, chart_selector: ChartSelector, data: ChartDataUpdate) -> None:
        self.execute(ops.OP_UPDATE_CHART_DATA, {"slide_index": slide_index, "chart_selector": chart_selector, "data": data})

