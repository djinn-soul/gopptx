"""Comment and chart mixins for the Presentation API."""

from __future__ import annotations

import warnings
from typing import TYPE_CHECKING, cast

from .. import ops
from ..api_errors import GopptxError
from .helpers import PresentationProtocol

if TYPE_CHECKING:
    from ..schemas import (
        Author,
        ChartDataUpdate,
        ChartSelector,
        Comment,
        SlideChartRef,
    )
else:

    class PresentationProtocol:
        """Runtime placeholder to avoid Protocol abstract behavior."""


class PresentationCommentMixin(PresentationProtocol):
    """Mixin providing comment and author management methods."""

    def get_authors(self) -> list[Author]:
        """Get all comment authors in the presentation."""
        result = self.execute(ops.OP_GET_AUTHORS, {})
        return cast("list[Author]", result.get("authors", []))

    def add_author(self, name: str, initials: str) -> int:
        """Add a comment author to the presentation."""
        result = self.execute(ops.OP_ADD_AUTHOR, {"name": name, "initials": initials})
        return int(cast("int", result.get("author_id", -1)))

    def get_comments(self, slide_index: int) -> list[Comment]:
        """Get all comments on a slide."""
        result = self.execute(ops.OP_GET_COMMENTS, {"slide_index": slide_index})
        raw_comments = result.get("comments")
        comments = cast(
            "list[Comment]", raw_comments if isinstance(raw_comments, list) else []
        )
        for item in cast("list[dict[str, object]]", comments):
            if "Index" in item and "index" not in item:
                item["index"] = item["Index"]
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
        for item in reversed(cast("list[dict[str, object]]", comments)):
            c_author = item.get("AuthorID", item.get("author_id"))
            if c_author is not None and int(str(c_author)) == author_id:
                c_idx = item.get("Index", item.get("index", 0))
                author_index = int(str(c_idx)) if c_idx is not None else 0
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


class PresentationChartMixin(PresentationProtocol):
    """Mixin providing chart creation and manipulation methods."""

    _BOUNDS_LEN = 4

    def add_chart(
        self,
        slide_index: int,
        chart_type: str,
        categories: list[str],
        values_or_series: list[float] | list[dict[str, str | list[float]]],
        **kwargs: str | tuple[float, float, float, float],
    ) -> int:
        """Add a chart to a slide."""
        bounds = kwargs.get("bounds", (0, 0, 0, 0))
        if len(bounds) != self._BOUNDS_LEN:
            raise ValueError("bounds must be a tuple of (x, y, w, h)")
        x, y, w, h = bounds
        title = kwargs.get("title", "Chart")

        values: list[float]
        if values_or_series and isinstance(values_or_series[0], dict):
            first = cast("dict[str, object]", values_or_series[0])
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
        return int(cast("int", result.get("shape_id") or result.get("chart_id", 0)))

    def list_slide_charts(self, slide_index: int) -> list[SlideChartRef]:
        """List all charts on a slide."""
        result = self.execute(ops.OP_LIST_SLIDE_CHARTS, {"slide_index": slide_index})
        return cast("list[SlideChartRef]", result.get("charts", []))

    def update_chart_data(
        self,
        slide_index: int,
        chart_selector: ChartSelector | list[str],
        data: ChartDataUpdate | list[dict[str, object]],
    ) -> None:
        """Update chart data for a chart on a slide."""
        if isinstance(chart_selector, dict):
            self.execute(
                ops.OP_UPDATE_CHART_DATA,
                {
                    "slide_index": slide_index,
                    "chart_selector": chart_selector,
                    "data": data,
                },
            )
            return

        charts = self.list_slide_charts(slide_index)
        selector: dict[str, object] = {"index": int(slide_index)}
        if charts:
            first = cast("dict[str, object]", charts[0])
            rel_id = first.get("RelID", first.get("rel_id"))
            if isinstance(rel_id, str) and rel_id:
                selector = {"rel_id": rel_id}
        categories = chart_selector
        series = cast("list[dict[str, object]]", data)
        normalized_series: list[dict[str, object]] = []
        for item in series:
            merged = dict(item)
            merged.setdefault("categories", categories)
            normalized_series.append(merged)
        payload = {
            "slide_index": slide_index,
            "chart_selector": selector,
            "data": {"categories": categories, "series": normalized_series},
        }
        try:
            self.execute(ops.OP_UPDATE_CHART_DATA, cast("dict[str, object]", payload))
        except GopptxError as error:
            warnings.warn(
                f"Chart data update failed (slide {slide_index}): {error}",
                UserWarning,
                stacklevel=2,
            )
