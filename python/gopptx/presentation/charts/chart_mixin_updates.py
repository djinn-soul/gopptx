"""Chart update/list helpers for the Presentation API."""

from __future__ import annotations

import warnings
from typing import TYPE_CHECKING, cast

from ... import ops
from ...api_errors import GopptxError
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ...schemas import (
        ChartDataUpdate,
        ChartFormatUpdate,
        ChartSelector,
        SlideChartRef,
    )


class PresentationChartUpdatesMixin(PresentationMixinBase):
    """Mixin for chart list/update operations."""

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
        selector: dict[str, object] = {"index": 0}
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

    def update_chart_data_batch(
        self,
        slide_index: int,
        updates: list[dict[str, object]],
    ) -> None:
        """Update multiple charts on one slide in a single bridge call."""
        if not updates:
            return
        self.execute(
            ops.OP_UPDATE_CHART_DATA_BATCH,
            {
                "slide_index": slide_index,
                "updates": [dict(item) for item in updates],
            },
        )

    def update_chart_formatting(
        self,
        slide_index: int,
        chart_selector: ChartSelector,
        fmt: ChartFormatUpdate,
    ) -> None:
        """Update chart formatting for a chart on a slide."""
        self.execute(
            ops.OP_UPDATE_CHART_FORMATTING,
            {
                "slide_index": slide_index,
                "chart_selector": chart_selector,
                "format": fmt,
            },
        )

    def update_chart_formatting_by_index(
        self,
        slide_index: int,
        chart_index: int,
        fmt: ChartFormatUpdate,
    ) -> None:
        """Update chart formatting by chart index."""
        self.update_chart_formatting(slide_index, {"index": chart_index}, fmt)

    def update_chart_formatting_by_rel_id(
        self,
        slide_index: int,
        rel_id: str,
        fmt: ChartFormatUpdate,
    ) -> None:
        """Update chart formatting by chart relationship id."""
        self.update_chart_formatting(slide_index, {"rel_id": rel_id}, fmt)

    def update_chart_data_by_index(
        self,
        slide_index: int,
        chart_index: int,
        data: ChartDataUpdate,
    ) -> None:
        """Update chart data by slide-local chart index."""
        self.update_chart_data(slide_index, {"index": chart_index}, data)

    def update_chart_data_by_rel_id(
        self,
        slide_index: int,
        rel_id: str,
        data: ChartDataUpdate,
    ) -> None:
        """Update chart data by chart relationship id."""
        self.update_chart_data(slide_index, {"rel_id": rel_id}, data)

    def replace_chart_data_by_index(
        self,
        slide_index: int,
        chart_index: int,
        categories: list[str],
        values: list[float],
    ) -> None:
        """Replace category/value chart data by slide-local chart index."""
        payload: ChartDataUpdate = {
            "categories": categories,
            "series": [{"values": values}],
        }
        self.update_chart_data_by_index(slide_index, chart_index, payload)

    def replace_chart_data_by_rel_id(
        self,
        slide_index: int,
        rel_id: str,
        categories: list[str],
        values: list[float],
    ) -> None:
        """Replace category/value chart data by chart relationship id."""
        payload: ChartDataUpdate = {
            "categories": categories,
            "series": [{"values": values}],
        }
        self.update_chart_data_by_rel_id(slide_index, rel_id, payload)
