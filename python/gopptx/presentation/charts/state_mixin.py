"""Chart state read helpers for Presentation."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ...schemas import ChartSelector, ChartState


class PresentationChartStateMixin(PresentationMixinBase):
    """Mixin providing chart state lookup helpers."""

    def get_chart_state(
        self,
        slide_index: int,
        chart_selector: ChartSelector,
    ) -> ChartState:
        """Return chart state selected by a chart selector on a slide."""
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
        """Return chart state by zero-based chart index on a slide."""
        return self.get_chart_state(slide_index, {"index": chart_index})

    def get_chart_state_by_rel_id(self, slide_index: int, rel_id: str) -> ChartState:
        """Return chart state by relationship id on a slide."""
        return self.get_chart_state(slide_index, {"rel_id": rel_id})
