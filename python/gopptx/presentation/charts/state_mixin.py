"""Chart state read helpers for Presentation."""
# ruff: noqa: D101,D102

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ...schemas import ChartSelector, ChartState


class PresentationChartStateMixin(PresentationMixinBase):
    def get_chart_state(
        self,
        slide_index: int,
        chart_selector: ChartSelector,
    ) -> ChartState:
        result = self.execute(
            ops.OP_GET_CHART_STATE,
            {
                "slide_index": slide_index,
                "chart_selector": chart_selector,
            },
        )
        return cast("ChartState", result.get("state", {}))

    def get_chart_state_by_index(
        self, slide_index: int, chart_index: int
    ) -> ChartState:
        return self.get_chart_state(slide_index, {"index": chart_index})

    def get_chart_state_by_rel_id(self, slide_index: int, rel_id: str) -> ChartState:
        return self.get_chart_state(slide_index, {"rel_id": rel_id})
