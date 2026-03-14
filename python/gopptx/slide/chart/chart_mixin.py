"""Slide chart mixin scoped to chart-domain operations."""

from __future__ import annotations

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from ...presentation.presentation import Presentation
    from ...schemas import SlideChartRef
    from ..chart_data import CategoryChartData, XyChartData


class SlideChartMixin:
    """Mixin providing chart-related methods for Slide objects."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int:
            """Slide index."""
            ...

    def list_charts(self) -> list[SlideChartRef]:
        """List chart references on this slide."""
        return self._presentation.list_slide_charts(self.index)

    def add_chart(
        self,
        chart_type: str,
        categories: list[str] | CategoryChartData | XyChartData,
        values_or_series: list[float]
        | list[dict[str, str | list[float]]]
        | None = None,
        **kwargs: str | tuple[float, float, float, float],
    ) -> int:
        """Add a chart and invalidate shape/text caches when available."""
        chart_id = self._presentation.add_chart(
            self.index, chart_type, categories, values_or_series, **kwargs
        )
        invalidate = getattr(self, "_invalidate_shape_cache", None)
        if callable(invalidate):
            invalidate()
        invalidate_text = getattr(self, "_invalidate_text_state_cache", None)
        if callable(invalidate_text):
            invalidate_text()
        return chart_id
